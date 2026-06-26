import React from 'react';
import {
  Box, Button, Chip, Container, Grid, Paper, Stack,
  Typography, alpha, useTheme, Divider,
} from '@mui/material';
import { Link } from 'react-router-dom';
import {
  Psychology as CHISGIcon,
  Schema as SchemaIcon,
  MenuBook as MenuBookIcon,
  RecordVoiceOver as AudioIcon,
  Verified as VerifiedIcon,
  Link as LinkIcon,
  AutoAwesome as AutoAwesomeIcon,
  Science as ScienceIcon,
} from '@mui/icons-material';
import ContactReveal from '../components/ContactReveal';
import SEO from '../components/SEO';

// ── Shared nav (matches CHISGLanding pattern) ──────────────────────────────
const ESPPilotNav: React.FC = () => (
  <Box sx={{ position: 'absolute', top: 0, left: 0, right: 0, zIndex: 10, py: 2, px: 3 }}>
    <Container maxWidth="lg">
      <Stack direction="row" justifyContent="space-between" alignItems="center">
        <Box
          component="img"
          src={`${process.env.PUBLIC_URL}/ESPLogoLong.png`}
          alt="ESP Thinking"
          sx={{ height: 56, width: 'auto' }}
        />
        <Stack direction="row" spacing={1}>
          <Button
            component={Link} to="/chisg" size="small"
            sx={{ color: 'white', textTransform: 'none', borderColor: 'rgba(255,255,255,0.4)',
              bgcolor: 'rgba(15,23,42,0.35)', '&:hover': { bgcolor: 'rgba(148,163,184,0.25)' } }}
            variant="outlined"
          >
            CHISG
          </Button>
          <Button
            component={Link} to="/" size="small"
            sx={{ color: 'white', textTransform: 'none', borderColor: 'rgba(255,255,255,0.6)',
              bgcolor: 'rgba(15,23,42,0.35)', '&:hover': { borderColor: 'white', bgcolor: 'rgba(148,163,184,0.35)' } }}
            variant="outlined"
          >
            Back to Portfolio
          </Button>
        </Stack>
      </Stack>
    </Container>
  </Box>
);

// ── Step card ──────────────────────────────────────────────────────────────
const Step: React.FC<{ n: string; title: string; body: string }> = ({ n, title, body }) => (
  <Box sx={{ display: 'flex', gap: 3, alignItems: 'flex-start' }}>
    <Box sx={{
      flexShrink: 0, width: 44, height: 44, borderRadius: '50%',
      bgcolor: 'rgba(124,158,245,0.15)', border: '2px solid rgba(124,158,245,0.5)',
      color: '#7c9ef5', fontWeight: 800, fontSize: '1.1rem',
      display: 'flex', alignItems: 'center', justifyContent: 'center',
    }}>{n}</Box>
    <Box>
      <Typography variant="subtitle1" fontWeight={700} mb={0.5} color="white">{title}</Typography>
      <Typography variant="body2" sx={{ color: 'rgba(255,255,255,0.65)' }} lineHeight={1.8}>{body}</Typography>
    </Box>
  </Box>
);

// ── Feature card ───────────────────────────────────────────────────────────
const Feature: React.FC<{ icon: React.ReactNode; title: string; body: string }> = ({ icon, title, body }) => {
  const theme = useTheme();
  return (
    <Paper elevation={0} sx={{
      p: 3, border: '1px solid', borderColor: 'divider', borderRadius: 2,
      transition: 'box-shadow 0.15s',
      '&:hover': { boxShadow: theme.shadows[3] },
    }}>
      <Box sx={{ color: theme.palette.primary.main, mb: 1.5 }}>{icon}</Box>
      <Typography variant="subtitle2" fontWeight={700} mb={0.5}>{title}</Typography>
      <Typography variant="body2" color="text.secondary" lineHeight={1.7}>{body}</Typography>
    </Paper>
  );
};

// ── Main page ──────────────────────────────────────────────────────────────
const ESPPilotLanding: React.FC = () => {
  const theme = useTheme();

  return (
    <Box>
      <SEO
        title="ESP Research Pilot — Personal Scientific Knowledge Oracle"
        description="ESP extracts causal relationships from research papers with full provenance, building a queryable knowledge graph that grows with your research career. Join the pilot study."
        path="/esp-pilot"
      />

      {/* ── Hero ── */}
      <Box sx={{
        position: 'relative',
        background: 'linear-gradient(135deg, #020617 0%, #0f172a 45%, #312e81 100%)',
        color: 'white',
        pt: { xs: 12, md: 16 },
        pb: { xs: 8, md: 12 },
        overflow: 'hidden',
        textAlign: 'center',
      }}>
        <ESPPilotNav />

        {/* Background glow */}
        <Box sx={{
          position: 'absolute', top: '20%', left: '50%', transform: 'translateX(-50%)',
          width: 600, height: 600, borderRadius: '50%',
          background: 'radial-gradient(circle, rgba(99,102,241,0.15) 0%, transparent 70%)',
          pointerEvents: 'none',
        }} />

        <Container maxWidth="md" sx={{ position: 'relative' }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 3, justifyContent: 'center' }}>
            <AutoAwesomeIcon sx={{ color: '#818cf8', fontSize: 26 }} />
            <Typography variant="overline" sx={{ color: '#818cf8', letterSpacing: 2.5, fontWeight: 700 }}>
              ESP · Expert Semantic Pipeline · Research Pilot
            </Typography>
          </Box>

          <Typography variant="h2" fontWeight={800} lineHeight={1.15} mb={3}
            sx={{ fontSize: { xs: '2rem', md: '3rem' } }}>
            Your personal scientific<br />knowledge oracle
          </Typography>

          <Typography variant="h6" sx={{
            color: 'rgba(255,255,255,0.65)', fontWeight: 400,
            maxWidth: 560, lineHeight: 1.8, mb: 5, mx: 'auto',
          }}>
            As you read research papers, ESP extracts and stores the causal relationships —
            with full provenance — into a queryable graph that grows with your career.
          </Typography>

          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} justifyContent="center" mb={4}>
            <Button
              variant="contained" size="large"
              href="#pilot"
              sx={{
                bgcolor: '#6366f1', '&:hover': { bgcolor: '#4f46e5' },
                fontWeight: 700, px: 4, borderRadius: 2,
              }}
            >
              Join the pilot study
            </Button>
            <Button
              variant="outlined" size="large"
              href="#how-it-works"
              sx={{
                borderColor: 'rgba(255,255,255,0.35)', color: 'white', borderRadius: 2,
                '&:hover': { borderColor: 'white', bgcolor: 'rgba(255,255,255,0.05)' },
              }}
            >
              How it works
            </Button>
          </Stack>

          <Stack direction="row" spacing={1} justifyContent="center" flexWrap="wrap" gap={1}>
            <Chip label="Research pilot · Free access" size="small"
              sx={{ bgcolor: 'rgba(99,102,241,0.2)', color: '#a5b4fc', border: '1px solid rgba(99,102,241,0.4)' }} />
            <Chip label="Built on CHISG knowledge architecture" size="small"
              sx={{ bgcolor: 'rgba(255,255,255,0.05)', color: 'rgba(255,255,255,0.5)', border: '1px solid rgba(255,255,255,0.1)' }} />
            <Chip label="Your corpus stays yours" size="small"
              sx={{ bgcolor: 'rgba(255,255,255,0.05)', color: 'rgba(255,255,255,0.5)', border: '1px solid rgba(255,255,255,0.1)' }} />
          </Stack>
        </Container>
      </Box>

      {/* ── Problem ── */}
      <Box sx={{ py: { xs: 8, md: 12 }, bgcolor: 'white' }}>
        <Container maxWidth="md">
          <Typography variant="overline" color="primary" fontWeight={700} letterSpacing={2}>
            The problem
          </Typography>
          <Typography variant="h4" fontWeight={800} mt={1} mb={3}>
            You read hundreds of papers.<br />The knowledge lives only in your head.
          </Typography>
          <Typography variant="body1" color="text.secondary" lineHeight={1.9} maxWidth={600}>
            A PhD researcher reads 200–400 papers over their candidature. The causal relationships,
            contested claims, and methodological caveats buried in that literature become
            tacit knowledge — impossible to query, impossible to hand off, and prone
            to the same gaps as human memory.
          </Typography>
          <Divider sx={{ my: 4 }} />
          <Typography variant="body1" color="text.secondary" lineHeight={1.9} maxWidth={600}>
            Existing tools summarise prose (losing structure) or require manual curation
            (which nobody has time for). Neither tracks <em>who established a claim</em>,{' '}
            <em>under what conditions</em>, or <em>whether it has been cited or contested</em>.
            ESP captures this structure automatically, as a byproduct of your normal reading workflow.
          </Typography>
        </Container>
      </Box>

      {/* ── How it works ── */}
      <Box id="how-it-works" sx={{
        py: { xs: 8, md: 12 },
        background: 'linear-gradient(135deg, #0f172a 0%, #1e1b4b 100%)',
      }}>
        <Container maxWidth="md">
          <Typography variant="overline" sx={{ color: '#818cf8', fontWeight: 700, letterSpacing: 2 }}>
            How it works
          </Typography>
          <Typography variant="h4" fontWeight={800} color="white" mt={1} mb={6}>
            Structured knowledge, as a byproduct of reading
          </Typography>
          <Stack spacing={4.5}>
            <Step n="1" title="Upload a paper (PDF or text)"
              body="Drag in any research paper. AWS Textract handles OCR for scanned PDFs. The full text is automatically chunked into reviewable paragraphs." />
            <Step n="2" title="AI extracts causal relationships"
              body="Each paragraph is sent to Claude. It identifies typed semantic links — 'X causes Y under condition Z' — along with the verbatim source sentence and, critically, whether the claim belongs to this paper or is attributed to a cited source (e.g. 'Li et al., 2022')." />
            <Step n="3" title="You review, correct, approve"
              body="A clean review interface shows source text alongside proposed links. Approve what's right, reject noise, edit anything incorrect. Each decision also improves the extraction model." />
            <Step n="4" title="Your corpus grows"
              body="Approved links are stored in your personal knowledge graph — typed, attributed, queryable, and exportable. The oracle layer (Stage 2) will let you ask how a new paper fits what you already know." />
          </Stack>
        </Container>
      </Box>

      {/* ── Stage 2 features ── */}
      <Box sx={{ py: { xs: 8, md: 12 }, bgcolor: alpha(theme.palette.primary.main, 0.03) }}>
        <Container maxWidth="md">
          <Typography variant="overline" color="primary" fontWeight={700} letterSpacing={2}>
            What's being built
          </Typography>
          <Typography variant="h4" fontWeight={800} mt={1} mb={5}>
            Stage 2 capabilities
          </Typography>
          <Grid container spacing={2}>
            {[
              { icon: <SchemaIcon />, title: '"How does this paper fit what I already know?"',
                body: 'Graph traversal across your reviewed corpus — overlapping entities, conflicting claims, novel connections — narrated by an LLM grounded only in your verified links.' },
              { icon: <AudioIcon />, title: 'Audio summaries',
                body: 'Spoken synthesis of what a new paper adds to your existing knowledge graph. Listen on a commute rather than reading from scratch.' },
              { icon: <VerifiedIcon />, title: 'Epistemic provenance',
                body: "Every link carries attribution: this paper's own experimental finding, a cited claim from another group, or a contested result. Confidence is first-class metadata." },
              { icon: <LinkIcon />, title: 'Citation graph resolution',
                body: "References cited in a paper are matched against your corpus. 'Gourse et al., 1996' becomes a traversable node if you've reviewed that paper." },
              { icon: <MenuBookIcon />, title: 'Portable corpus',
                body: 'Your knowledge graph is exportable at any time. It belongs to you — not the platform. Postdoc, new institution, career change: it travels with you.' },
              { icon: <ScienceIcon />, title: 'Federated model improvement',
                body: 'Review corrections are pooled (anonymised) to fine-tune the extraction model. Every participant benefits from a better extractor without sharing their corpus.' },
            ].map((f, i) => (
              <Grid item xs={12} sm={6} key={i}>
                <Feature {...f} />
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* ── Pilot ── */}
      <Box id="pilot" sx={{
        py: { xs: 8, md: 12 },
        background: 'linear-gradient(135deg, #0f172a 0%, #1e1b4b 100%)',
      }}>
        <Container maxWidth="md">
          <Typography variant="overline" sx={{ color: '#818cf8', fontWeight: 700, letterSpacing: 2 }}>
            Research pilot
          </Typography>
          <Typography variant="h4" fontWeight={800} color="white" mt={1} mb={3}>
            We're looking for researchers to test this with
          </Typography>
          <Typography sx={{ color: 'rgba(255,255,255,0.65)', lineHeight: 1.9, maxWidth: 580, mb: 4 }}>
            We're running a small pilot study — ideally 5–10 PhD students or early-career
            researchers in a science domain (biology, medicine, chemistry, or similar).
            Participation is free. API costs are covered centrally. Your corpus is yours.
          </Typography>

          <Typography variant="subtitle2" fontWeight={700} color="white" mb={2}>
            What participation involves:
          </Typography>
          <Stack spacing={1.5} mb={5}>
            {[
              'Upload 10–20 papers from your literature review over 3–6 months',
              'Spend ~5 minutes per paper reviewing extracted links in the web interface',
              'Complete a short survey on usefulness and extraction accuracy',
              'Optional: short interview at the end of the study period',
            ].map((item, i) => (
              <Box key={i} sx={{ display: 'flex', gap: 2, alignItems: 'flex-start' }}>
                <Box sx={{
                  width: 7, height: 7, borderRadius: '50%', bgcolor: '#818cf8',
                  mt: '6px', flexShrink: 0,
                }} />
                <Typography variant="body2" sx={{ color: 'rgba(255,255,255,0.65)' }}>{item}</Typography>
              </Box>
            ))}
          </Stack>

          <ContactReveal
            email="michael.stewart@espthinking.co.uk"
            label="Express interest in the ESP pilot"
          />

          <Typography variant="caption" sx={{ display: 'block', mt: 3, color: 'rgba(255,255,255,0.3)' }}>
            michael.stewart@espthinking.co.uk · Response within 1–2 working days
          </Typography>
        </Container>
      </Box>

      {/* ── CHISG link ── */}
      <Box sx={{ py: { xs: 6, md: 8 }, bgcolor: 'white' }}>
        <Container maxWidth="md">
          <Paper elevation={0} sx={{
            p: 4, border: '1px solid', borderColor: 'divider', borderRadius: 3,
            display: 'flex', flexDirection: { xs: 'column', sm: 'row' },
            alignItems: { sm: 'center' }, gap: 3,
          }}>
            <CHISGIcon sx={{ fontSize: 44, color: theme.palette.primary.main, flexShrink: 0 }} />
            <Box sx={{ flex: 1 }}>
              <Typography variant="subtitle1" fontWeight={700} mb={0.5}>
                Built on the CHISG knowledge architecture
              </Typography>
              <Typography variant="body2" color="text.secondary" lineHeight={1.7}>
                ESP is the academic research application of CHISG — Contextualised Hierarchical
                Iterative Semantic Groupings — the same methodology that powers the skills graph,
                semantic query engine, and knowledge extraction pipeline across the ESP portfolio.
              </Typography>
            </Box>
            <Button
              component={Link} to="/chisg"
              variant="outlined" size="small"
              sx={{ flexShrink: 0, textTransform: 'none' }}
            >
              Learn about CHISG →
            </Button>
          </Paper>
        </Container>
      </Box>

      {/* ── Footer ── */}
      <Box sx={{ bgcolor: '#0a0a12', py: 4 }}>
        <Container maxWidth="lg">
          <Stack direction={{ xs: 'column', sm: 'row' }} justifyContent="space-between" alignItems="center" spacing={2}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <AutoAwesomeIcon sx={{ color: '#818cf8', fontSize: 18 }} />
              <Typography variant="caption" sx={{ color: 'rgba(255,255,255,0.3)' }}>
                ESP Research Pilot · Part of the ESP Thinking portfolio · espthinking.co.uk
              </Typography>
            </Box>
            <Stack direction="row" spacing={2}>
              <Button component={Link} to="/" size="small"
                sx={{ color: 'rgba(255,255,255,0.4)', textTransform: 'none', fontSize: '0.75rem' }}>
                Portfolio
              </Button>
              <Button component={Link} to="/chisg" size="small"
                sx={{ color: 'rgba(255,255,255,0.4)', textTransform: 'none', fontSize: '0.75rem' }}>
                CHISG
              </Button>
            </Stack>
          </Stack>
        </Container>
      </Box>
    </Box>
  );
};

export default ESPPilotLanding;
