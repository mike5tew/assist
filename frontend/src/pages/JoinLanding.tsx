import React from 'react';
import {
  Box,
  Container,
  Typography,
  TextField,
  Button,
  Paper,
  Stack,
  MenuItem,
  Alert,
  useTheme,
  alpha,
} from '@mui/material';
import { Link } from 'react-router-dom';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';
import SEO from '../components/SEO';
import { trackEvent } from '../hooks/useAnalytics';

const brandColor = '#1a365d';

const roles = [
  { value: 'parent', label: 'Parent / Carer' },
  { value: 'teacher', label: 'Teacher / Educator' },
  { value: 'clinician', label: 'Clinician / Therapist' },
  { value: 'researcher', label: 'Researcher' },
  { value: 'other', label: 'Other / Just Curious' },
];

const JoinLanding: React.FC = () => {
  const theme = useTheme();
  const [email, setEmail] = React.useState('');
  const [role, setRole] = React.useState('');
  const [source] = React.useState(() => {
    const params = new URLSearchParams(window.location.search);
    return params.get('src') || document.referrer || 'direct';
  });
  const [status, setStatus] = React.useState<'idle' | 'submitting' | 'success' | 'already' | 'error'>('idle');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email.trim()) return;

    setStatus('submitting');
    try {
      const resp = await fetch('/api/join', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: email.trim(), source, role }),
      });
      if (resp.ok) {
        const data = await resp.json();
        if (data.status === 'already_subscribed') {
          setStatus('already');
        } else {
          setStatus('success');
        }
        trackEvent('join_subscribed', window.location.pathname, { role, source });
      } else {
        setStatus('error');
      }
    } catch {
      setStatus('error');
    }
  };

  return (
    <>
      <SEO
        title="Join — ESP Thinking"
        description="Follow the thinking. A biological profile and skills map that replaces blanket labels with specific, actionable patterns."
        path="/join"
      />

      {/* Hero */}
      <Box
        sx={{
          background: `linear-gradient(135deg, ${brandColor} 0%, #2c5282 100%)`,
          color: 'white',
          pt: { xs: 6, md: 10 },
          pb: { xs: 10, md: 14 },
          position: 'relative',
          '&::after': {
            content: '""',
            position: 'absolute',
            bottom: 0,
            left: 0,
            right: 0,
            height: '60px',
            background: '#fff',
            clipPath: 'polygon(0 100%, 100% 100%, 100% 0)',
          },
        }}
      >
        <Container maxWidth={false} sx={{ maxWidth: '900px', position: 'relative', zIndex: 1 }}>
          {/* Back link */}
          <Box sx={{ mb: 4 }}>
            <Link to="/" style={{ color: 'rgba(255,255,255,0.7)', textDecoration: 'none', display: 'inline-flex', alignItems: 'center', gap: 4, fontSize: '0.9rem' }}>
              <ArrowBackIcon sx={{ fontSize: '1rem' }} /> espthinking.co.uk
            </Link>
          </Box>

          <Typography variant="h3" sx={{ fontWeight: 800, mb: 2, fontSize: { xs: '1.8rem', md: '2.6rem' } }}>
            What if labels told you what to <em>do</em>?
          </Typography>
          <Typography variant="h6" sx={{ fontWeight: 400, opacity: 0.9, mb: 0, lineHeight: 1.6, maxWidth: '700px' }}>
            We're building a framework that replaces blanket diagnostic terms with a biological profile and a skills map.
            Every child gets the same map. The "conditions" are just recognisable clusters on it.
          </Typography>
        </Container>
      </Box>

      {/* Main content */}
      <Box sx={{ py: { xs: 6, md: 8 }, background: '#fff' }}>
        <Container maxWidth={false} sx={{ maxWidth: '900px' }}>
          <Stack spacing={6}>

            {/* The pitch */}
            <Box>
              <Typography variant="h5" sx={{ fontWeight: 700, mb: 3, color: brandColor }}>
                The problem
              </Typography>
              <Typography sx={{ fontSize: '1.05rem', lineHeight: 1.8, color: '#333', mb: 2 }}>
                "Your child has autism" tells you <strong>what</strong>. It doesn't tell you <strong>where</strong> — 
                which skills are at precursor level, which biological settings are running hot, 
                or what to do on Monday morning.
              </Typography>
              <Typography sx={{ fontSize: '1.05rem', lineHeight: 1.8, color: '#333', mb: 2 }}>
                Two diagnoses? The advice contradicts itself. Autism says routine. ADHD says novelty. 
                The parent stands between two manuals that disagree.
              </Typography>
              <Typography sx={{ fontSize: '1.05rem', lineHeight: 1.8, color: '#333' }}>
                We've built a different approach: <strong>12 biological spectra</strong> that describe how every 
                person processes the world, and a <strong>skills graph</strong> with 200+ nodes connected by prerequisite links. 
                A "condition" is a recognisable cluster on the same map everyone uses.
              </Typography>
            </Box>

            {/* What you get */}
            <Box>
              <Typography variant="h5" sx={{ fontWeight: 700, mb: 3, color: brandColor }}>
                What you'll get
              </Typography>
              <Stack spacing={1.5}>
                {[
                  'The thinking behind the framework — one email at a time, not a textbook',
                  'Early access to the ETP profile tool when it launches',
                  'First look at ToddlerOS — the physical activity book built on these spectra',
                  'Direct replies from the person building it (just me — no marketing team)',
                ].map((item, i) => (
                  <Box key={i} sx={{ display: 'flex', gap: 1.5, alignItems: 'flex-start' }}>
                    <Box sx={{ width: 6, height: 6, borderRadius: '50%', bgcolor: brandColor, mt: 1, flexShrink: 0 }} />
                    <Typography sx={{ fontSize: '1rem', lineHeight: 1.6, color: '#444' }}>{item}</Typography>
                  </Box>
                ))}
              </Stack>
            </Box>

            {/* Signup form */}
            <Paper
              elevation={0}
              sx={{
                p: { xs: 3, md: 4 },
                borderRadius: 3,
                border: `1px solid ${alpha(brandColor, 0.15)}`,
                background: alpha(brandColor, 0.02),
              }}
            >
              {status === 'success' || status === 'already' ? (
                <Box sx={{ textAlign: 'center', py: 3 }}>
                  <CheckCircleOutlineIcon sx={{ fontSize: 48, color: theme.palette.secondary.main, mb: 2 }} />
                  <Typography variant="h6" sx={{ fontWeight: 600, mb: 1 }}>
                    {status === 'already' ? "You're already on the list." : "You're in."}
                  </Typography>
                  <Typography sx={{ color: '#666' }}>
                    {status === 'already'
                      ? "We've got your email — you'll hear from us soon."
                      : "First email coming soon. No spam. Unsubscribe any time."}
                  </Typography>
                </Box>
              ) : (
                <form onSubmit={handleSubmit}>
                  <Typography variant="h6" sx={{ fontWeight: 700, mb: 1, color: brandColor }}>
                    Follow the thinking
                  </Typography>
                  <Typography sx={{ color: '#666', mb: 3, fontSize: '0.95rem' }}>
                    Leave your email. I'll send the framework — not a sales pitch.
                  </Typography>

                  {status === 'error' && (
                    <Alert severity="error" sx={{ mb: 2 }}>
                      Something went wrong. Try again or email world@espthinking.co.uk directly.
                    </Alert>
                  )}

                  <Stack spacing={2}>
                    <TextField
                      label="Email"
                      type="email"
                      required
                      fullWidth
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      size="medium"
                    />
                    <TextField
                      label="I'm a..."
                      select
                      fullWidth
                      value={role}
                      onChange={(e) => setRole(e.target.value)}
                      size="medium"
                      helperText="Optional — helps me tailor what I send"
                    >
                      {roles.map((r) => (
                        <MenuItem key={r.value} value={r.value}>{r.label}</MenuItem>
                      ))}
                    </TextField>
                    <Button
                      type="submit"
                      variant="contained"
                      size="large"
                      disabled={status === 'submitting' || !email.trim()}
                      sx={{
                        bgcolor: brandColor,
                        fontWeight: 700,
                        py: 1.5,
                        fontSize: '1rem',
                        '&:hover': { bgcolor: '#2c5282' },
                      }}
                    >
                      {status === 'submitting' ? 'Joining...' : 'Join'}
                    </Button>
                  </Stack>

                  <Typography sx={{ mt: 2, fontSize: '0.8rem', color: '#999' }}>
                    No spam. No selling your email. Unsubscribe any time.
                  </Typography>
                </form>
              )}
            </Paper>

            {/* Social proof / context */}
            <Box sx={{ textAlign: 'center', py: 2 }}>
              <Typography sx={{ fontSize: '0.95rem', color: '#888', fontStyle: 'italic' }}>
                "Diagnosis opens the door. The profile opens the right page."
              </Typography>
            </Box>

          </Stack>
        </Container>
      </Box>

      {/* Footer */}
      <Box sx={{ py: 4, background: alpha(brandColor, 0.03), borderTop: '1px solid #eee', textAlign: 'center' }}>
        <Typography variant="body2" sx={{ color: '#999' }}>
          ESP Thinking · espthinking.co.uk · Built by{' '}
          <a href="https://www.linkedin.com/in/michael-stewart-esp/" target="_blank" rel="noopener noreferrer" style={{ color: '#888' }}>
            Michael Stewart
          </a>
        </Typography>
      </Box>
    </>
  );
};

export default JoinLanding;
