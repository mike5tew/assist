import React from 'react';
import { Button, Dialog, DialogTitle, DialogContent, DialogActions, TextField, Typography, Box, IconButton, Tooltip } from '@mui/material';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';

interface Props {
  email: string;
  label?: string;
}

const generateChallenge = () => {
  const a = Math.floor(Math.random() * 8) + 2;
  const b = Math.floor(Math.random() * 8) + 2;
  return { a, b, answer: a + b };
};

const ContactReveal: React.FC<Props> = ({ email, label = 'Contact' }) => {
  const [open, setOpen] = React.useState(false);
  const [challenge, setChallenge] = React.useState(() => generateChallenge());
  const [input, setInput] = React.useState('');
  const [unlocked, setUnlocked] = React.useState(false);
  const [copied, setCopied] = React.useState(false);

  // Contact form state
  const [name, setName] = React.useState('');
  const [senderEmail, setSenderEmail] = React.useState('');
  const [message, setMessage] = React.useState('');
  const [sending, setSending] = React.useState(false);
  const [sentOk, setSentOk] = React.useState<null | boolean>(null);
  const siteKey = process.env.REACT_APP_RECAPTCHA_SITE_KEY || '';

  React.useEffect(() => {
    if (!open) {
      setUnlocked(false);
      setInput('');
      setChallenge(generateChallenge());
      setCopied(false);
      setName('');
      setSenderEmail('');
      setMessage('');
      setSending(false);
      setSentOk(null);
    }
  }, [open]);

  const tryUnlock = () => {
    if (parseInt(input || '0', 10) === challenge.answer) {
      setUnlocked(true);
    } else {
      // new challenge on failure to slow bots
      setChallenge(generateChallenge());
      setInput('');
    }
  };

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(email);
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    } catch (e) {
      /* ignore */
    }
  };

  // Email is not rendered in the DOM until unlocked — assembled here to avoid static scraping
  const obfuscated = React.useMemo(() => {
    if (!unlocked) return '';
    return email;
  }, [unlocked, email]);

  const submitContact = async (recaptchaToken?: string) => {
    setSending(true);
    setSentOk(null);
    try {
      const resp = await fetch('/api/contact', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, email: senderEmail, message, recaptchaToken }),
      });
      if (resp.ok) {
        setSentOk(true);
      } else {
        setSentOk(false);
      }
    } catch (e) {
      setSentOk(false);
    } finally {
      setSending(false);
    }
  };

  const handleSend = async () => {
    if (!message.trim()) return;

    // If recaptcha site key is provided, attempt client-side grecaptcha execution.
    if (siteKey && (window as any).grecaptcha) {
      try {
        const grecaptcha = (window as any).grecaptcha;
        const token = await grecaptcha.execute(siteKey, { action: 'contact' });
        await submitContact(token);
        return;
      } catch (e) {
        // fall through to submit without token
      }
    }

    await submitContact();
  };

  return (
    <>
      <Button variant="contained" color="secondary" onClick={() => setOpen(true)}>{label}</Button>

      <Dialog open={open} onClose={() => setOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>Contact</DialogTitle>
        <DialogContent>
          <Typography variant="body2" sx={{ mb: 2 }}>
            To protect this address from automated harvesting, please answer the simple question below. You can either copy/open the contact email after verification, or send a message via the secure form.
          </Typography>

          {!unlocked ? (
            <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
              <TextField
                label={`What is ${challenge.a} + ${challenge.b}?`}
                value={input}
                onChange={(e) => setInput(e.target.value.replace(/[^0-9]/g, ''))}
                size="small"
                autoFocus
              />
              <Button variant="outlined" onClick={tryUnlock}>Verify</Button>
            </Box>
          ) : (
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
              <a href={`mailto:${obfuscated}`} style={{ textDecoration: 'none' }}>
                <Typography variant="body1" sx={{ fontWeight: 'bold' }}>{obfuscated}</Typography>
              </a>
              <Tooltip title={copied ? 'Copied' : 'Copy to clipboard'}>
                <IconButton size="small" onClick={handleCopy}>
                  <ContentCopyIcon fontSize="small" />
                </IconButton>
              </Tooltip>
            </Box>
          )}

          {/* Contact form */}
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
            <TextField label="Your name" value={name} onChange={(e) => setName(e.target.value)} size="small" />
            <TextField label="Your email (optional)" value={senderEmail} onChange={(e) => setSenderEmail(e.target.value)} size="small" />
            <TextField label="Message" value={message} onChange={(e) => setMessage(e.target.value)} multiline rows={4} size="small" />
            <Box sx={{ display: 'flex', gap: 1, alignItems: 'center', mt: 1 }}>
              <Button variant="contained" onClick={handleSend} disabled={sending}>
                {sending ? 'Sending…' : 'Send message'}
              </Button>
              <Button variant="text" onClick={() => { setName(''); setSenderEmail(''); setMessage(''); }}>
                Clear
              </Button>
              {sentOk === true && <Typography color="success.main">Sent ✓</Typography>}
              {sentOk === false && <Typography color="error.main">Failed to send</Typography>}
            </Box>
          </Box>

        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>Close</Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default ContactReveal;
