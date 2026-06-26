import React from 'react';
import {
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  Container,
  Grid,
  Paper,
  Stack,
  Typography,
  alpha,
  useTheme,
} from '@mui/material';
import { Link } from 'react-router-dom';
import {
  SettingsEthernet as LoadBalanceIcon,
  Shield as SafeguardIcon,
  Speed as SpeedIcon,
  Storage as CacheIcon,
  Hub as PipelineIcon,
  Memory as ReplicaIcon,
} from '@mui/icons-material';
import ContactReveal from '../../components/ContactReveal';
import SEO from '../../components/SEO';

const ClassifierNav: React.FC = () => (
  <Box sx={{ position: 'absolute', top: 0, left: 0, right: 0, zIndex: 10, py: 2, px: 3 }}>
    <Container maxWidth="lg">
      <Stack direction="row" justifyContent="space-between" alignItems="center">
        <Box
          component="img"
          src={`${process.env.PUBLIC_URL}/ESPLogoLong.png`}
          alt="ESP Thinking Portfolio"
          sx={{ height: 56, width: 'auto' }}
        />
        <Stack direction="row" spacing={1}>
          <Button
            component={Link}
            to="/"
            variant="outlined"
            size="small"
            sx={{
              color: 'white',
              borderColor: 'rgba(255,255,255,0.6)',
              textTransform: 'none',
              bgcolor: 'rgba(15,23,42,0.35)',
              '&:hover': { borderColor: 'white', bgcolor: 'rgba(148,163,184,0.35)' },
            }}
          >
            Back to Portfolio
          </Button>
        </Stack>
      </Stack>
    </Container>
  </Box>
);

const ClassifierLanding: React.FC = () => {
  const theme = useTheme();

  return (
    <Box>
      <SEO
        title="CHISG Classifier — Load Balanced Classification Pipeline"
        description="A horizontally-scaled Go microservice demonstrating production load balancing, concurrent pipelines, Redis caching, and safeguarding-first design for real-time student observation classification."
        path="/classifier-demo"
      />

      {/* Hero */}
      <Box
        sx={{
          position: 'relative',
          background: 'linear-gradient(135deg, #0c0a1a 0%, #1a1040 40%, #3730a3 100%)',
          color: 'white',
          pt: { xs: 10, md: 14 },
          pb: { xs: 8, md: 12 },
          overflow: 'hidden',
        }}
      >
        <ClassifierNav />
        <Container maxWidth="lg">
          <Grid container spacing={6} alignItems="center">
            <Grid item xs={12} md={7}>
              <Typography variant="overline" sx={{ opacity: 0.85 }}>
                Systems Architecture Demo
              </Typography>
              <Typography
                variant="h2"
                component="h1"
                fontWeight="bold"
                sx={{ fontSize: { xs: '2rem', md: '3rem' }, lineHeight: 1.15, mb: 2 }}
              >
                CHISG Classifier
              </Typography>
              <Typography
                variant="h5"
                sx={{ opacity: 0.9, mb: 3, fontWeight: 400, lineHeight: 1.5 }}
              >
                A load-balanced Go microservice that classifies child observations across 12 emotional spectra in real time.
                Three replicas, one cache, zero single points of failure.
              </Typography>
              <Typography variant="body1" sx={{ opacity: 0.8, mb: 4, maxWidth: 560 }}>
                If your data pipeline has a single point of failure, it's not a pipeline — it's a prayer.
                This demo shows how production-grade infrastructure handles concurrent classification at scale.
              </Typography>

              <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
                <Button
                  variant="contained"
                  size="large"
                  href="/classifier/"
                  sx={{
                    bgcolor: '#6366f1',
                    textTransform: 'none',
                    fontWeight: 600,
                    px: 4,
                    '&:hover': { bgcolor: '#4f46e5' },
                  }}
                >
                  Open Live Dashboard
                </Button>
                <Button
                  component={Link}
                  to="/"
                  variant="outlined"
                  size="large"
                  sx={{
                    color: 'white',
                    borderColor: 'rgba(255,255,255,0.5)',
                    textTransform: 'none',
                    '&:hover': { borderColor: 'white', bgcolor: 'rgba(255,255,255,0.08)' },
                  }}
                >
                  Back to Portfolio
                </Button>
              </Stack>
            </Grid>

            <Grid item xs={12} md={5}>
              <Paper
                elevation={0}
                sx={{
                  p: 3,
                  bgcolor: 'rgba(255,255,255,0.06)',
                  border: '1px solid rgba(255,255,255,0.12)',
                  borderRadius: 3,
                }}
              >
                <Typography variant="subtitle2" sx={{ opacity: 0.7, mb: 2, textTransform: 'uppercase', letterSpacing: 1 }}>
                  Architecture at a Glance
                </Typography>
                <Box
                  sx={{
                    fontFamily: 'monospace',
                    fontSize: '0.78rem',
                    lineHeight: 1.7,
                    color: '#a5b4fc',
                    whiteSpace: 'pre',
                    overflowX: 'auto',
                  }}
                >
{`Request ──▶ Nginx (least-conn)
             ├──▶ chisg-1 (Go)
             ├──▶ chisg-2 (Go) ──▶ Redis
             └──▶ chisg-3 (Go)

Pipeline per replica:
  1. Cache Check  (SHA-256 → Redis)
  2. Safeguarding (regex scan)
  3. Fan-Out      (spectra + barriers)`}
                </Box>
                <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap sx={{ mt: 2 }}>
                  {['Go 1.24', 'Nginx', 'Redis 7', 'Docker Compose', 'Goroutines'].map((t) => (
                    <Chip key={t} label={t} size="small" sx={{ bgcolor: 'rgba(255,255,255,0.1)', color: 'white', fontSize: '0.72rem' }} />
                  ))}
                </Stack>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* Why it matters — framed for MAT / school context */}
      <Box sx={{ py: 8, bgcolor: 'rgba(15,23,42,0.03)' }}>
        <Container maxWidth="lg">
          <Typography variant="h4" component="h2" gutterBottom textAlign="center" fontWeight="bold">
            Why This Matters for Trust Data
          </Typography>
          <Typography variant="body1" color="text.secondary" textAlign="center" sx={{ mb: 5, maxWidth: 700, mx: 'auto' }}>
            Multi-Academy Trusts pull data from dozens of schools, each with different MIS systems, attendance platforms, and assessment cycles.
            The classifier demonstrates the exact patterns needed to handle that complexity reliably.
          </Typography>

          <Grid container spacing={3}>
            {[
              {
                icon: <ReplicaIcon fontSize="large" />,
                title: 'Horizontal Scaling',
                desc: 'Three identical replicas share the load. If one fails, the others continue. This is how you build trust-wide dashboards that don\'t go down during census week.',
                color: '#3b82f6',
              },
              {
                icon: <CacheIcon fontSize="large" />,
                title: 'Shared Cache Layer',
                desc: 'Redis stores results so identical queries return instantly. In a school context: the same attendance report pulled by 12 heads of year doesn\'t hit the database 12 times.',
                color: '#22c55e',
              },
              {
                icon: <SafeguardIcon fontSize="large" />,
                title: 'Safeguarding-First Design',
                desc: 'The pipeline short-circuits on safeguarding concerns — classification stops, the alert fires. The architecture principle: safety checks before analytics, always.',
                color: '#ef4444',
              },
              {
                icon: <SpeedIcon fontSize="large" />,
                title: 'Concurrent Processing',
                desc: 'Go goroutines run spectrum analysis and barrier detection simultaneously with context timeouts. This is how you process 10,000 pupil records without making staff wait.',
                color: '#f59e0b',
              },
            ].map((item) => (
              <Grid item xs={12} sm={6} key={item.title}>
                <Card
                  elevation={1}
                  sx={{
                    height: '100%',
                    borderTop: `4px solid ${item.color}`,
                    transition: 'transform 0.2s',
                    '&:hover': { transform: 'translateY(-2px)', boxShadow: theme.shadows[4] },
                  }}
                >
                  <CardContent>
                    <Box sx={{ color: item.color, mb: 1.5 }}>{item.icon}</Box>
                    <Typography variant="h6" gutterBottom fontWeight="bold">{item.title}</Typography>
                    <Typography variant="body2" color="text.secondary">{item.desc}</Typography>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* The Three-Stage Pipeline */}
      <Box sx={{ py: 8 }}>
        <Container maxWidth="lg">
          <Typography variant="h4" component="h2" gutterBottom textAlign="center" fontWeight="bold">
            Three-Stage Classification Pipeline
          </Typography>
          <Typography variant="body1" color="text.secondary" textAlign="center" sx={{ mb: 5, maxWidth: 650, mx: 'auto' }}>
            Every request follows the same deterministic path. No randomness, no black boxes — just structured, auditable processing.
          </Typography>

          <Grid container spacing={4}>
            {[
              {
                stage: '1',
                title: 'Cache Check',
                detail: 'SHA-256 hash of input text checked against Redis. 5-minute TTL. On hit, return cached result instantly — no compute wasted.',
                color: '#22c55e',
                icon: <CacheIcon />,
              },
              {
                stage: '2',
                title: 'Safeguarding Scan',
                detail: 'Five compiled regex patterns scan for indicators of harm: sexual abuse, violence, self-harm, neglect, domestic violence. Severity escalates for children under 8. Pipeline short-circuits on detection.',
                color: '#ef4444',
                icon: <SafeguardIcon />,
              },
              {
                stage: '3',
                title: 'Fan-Out Classification',
                detail: 'Two goroutines run concurrently: spectrum analysis (12 bipolar ETP spectra) and barrier detection (5 engagement blockers). Results collected via Go channels with 500ms context timeout.',
                color: '#6366f1',
                icon: <PipelineIcon />,
              },
            ].map((s) => (
              <Grid item xs={12} md={4} key={s.stage}>
                <Paper
                  elevation={0}
                  sx={{
                    p: 3,
                    height: '100%',
                    border: `1px solid ${alpha(s.color, 0.3)}`,
                    borderRadius: 3,
                    position: 'relative',
                    overflow: 'hidden',
                  }}
                >
                  <Box
                    sx={{
                      position: 'absolute',
                      top: -10,
                      right: -10,
                      fontSize: '6rem',
                      fontWeight: 900,
                      color: alpha(s.color, 0.06),
                      lineHeight: 1,
                    }}
                  >
                    {s.stage}
                  </Box>
                  <Box sx={{ color: s.color, mb: 1.5 }}>{s.icon}</Box>
                  <Typography variant="h6" fontWeight="bold" gutterBottom>
                    Stage {s.stage}: {s.title}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">{s.detail}</Typography>
                </Paper>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* 12 ETP Spectra */}
      <Box sx={{ py: 8, bgcolor: 'rgba(15,23,42,0.03)' }}>
        <Container maxWidth="lg">
          <Typography variant="h4" component="h2" gutterBottom textAlign="center" fontWeight="bold">
            12 ETP Spectra Classified
          </Typography>
          <Typography variant="body1" color="text.secondary" textAlign="center" sx={{ mb: 4, maxWidth: 650, mx: 'auto' }}>
            Each spectrum has two opposing poles. The engine scores input text against both poles and returns the stronger match with a confidence score.
          </Typography>

          <Grid container spacing={1.5}>
            {[
              { name: 'Social Gravity', poles: 'isolated ↔ orbiting' },
              { name: 'Energy Directionality', poles: 'inward ↔ outward' },
              { name: 'Voltage Sensitivity', poles: 'hypo ↔ hyper' },
              { name: 'Threat Response', poles: 'freeze ↔ fight' },
              { name: 'Care Response', poles: 'self-focused ↔ other-focused' },
              { name: 'Risk Tolerance', poles: 'avoidant ↔ seeking' },
              { name: 'Integrity Logic', poles: 'deflective ↔ ownership' },
              { name: 'Mirror Neuron Tuning', poles: 'flat ↔ resonant' },
              { name: 'Orderliness', poles: 'chaotic ↔ rigid' },
              { name: 'Responsibility Threshold', poles: 'avoidant ↔ over-responsible' },
              { name: 'Loss Sensitivity', poles: 'hoarding ↔ detached' },
              { name: 'Libido', poles: 'withdrawn ↔ expressive' },
            ].map((s) => (
              <Grid item xs={12} sm={6} md={4} key={s.name}>
                <Paper
                  elevation={0}
                  sx={{
                    p: 2,
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    border: '1px solid',
                    borderColor: 'divider',
                    borderRadius: 2,
                  }}
                >
                  <Typography variant="body2" fontWeight="bold">{s.name}</Typography>
                  <Typography variant="caption" color="text.secondary">{s.poles}</Typography>
                </Paper>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* Technology Stack */}
      <Box sx={{ py: 8 }}>
        <Container maxWidth="lg">
          <Typography variant="h4" component="h2" gutterBottom textAlign="center" fontWeight="bold">
            Technology Stack
          </Typography>

          <Grid container spacing={3} sx={{ mt: 2 }}>
            {[
              { label: 'Go 1.24', desc: 'Classifier service with goroutine fan-out and context timeouts', color: '#00ADD8' },
              { label: 'Nginx Alpine', desc: 'L7 load balancer with least-connections routing', color: '#009639' },
              { label: 'Redis 7 Alpine', desc: 'Classification result cache with SHA-256 keying', color: '#DC382D' },
              { label: 'Docker Compose', desc: 'Orchestration of all 5 containers from a single config', color: '#2496ED' },
              { label: 'Multi-stage Build', desc: 'golang:1.24-alpine build → alpine:3.19 runtime (~15MB image)', color: '#6366f1' },
            ].map((t) => (
              <Grid item xs={12} sm={6} md={4} key={t.label}>
                <Paper elevation={0} sx={{ p: 2.5, border: '1px solid', borderColor: 'divider', borderRadius: 2, height: '100%' }}>
                  <Chip label={t.label} size="small" sx={{ bgcolor: alpha(t.color, 0.12), color: t.color, fontWeight: 700, mb: 1 }} />
                  <Typography variant="body2" color="text.secondary">{t.desc}</Typography>
                </Paper>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* CTA */}
      <Box
        sx={{
          py: 8,
          background: 'linear-gradient(135deg, #0c0a1a 0%, #1a1040 40%, #3730a3 100%)',
          color: 'white',
          textAlign: 'center',
        }}
      >
        <Container maxWidth="sm">
          <Typography variant="h4" fontWeight="bold" gutterBottom>
            Try it Live
          </Typography>
          <Typography variant="body1" sx={{ opacity: 0.85, mb: 4 }}>
            The full dashboard is running on the same production server as this portfolio. Classify observations, run a 50-request load test, and watch requests distribute across replicas in real time.
          </Typography>
          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} justifyContent="center">
            <Button
              variant="contained"
              size="large"
              href="/classifier/"
              sx={{
                bgcolor: '#6366f1',
                textTransform: 'none',
                fontWeight: 600,
                px: 4,
                '&:hover': { bgcolor: '#4f46e5' },
              }}
            >
              Open Live Dashboard
            </Button>
            <ContactReveal email="mike@espthinking.co.uk" label="Get in Touch" />
          </Stack>
        </Container>
      </Box>
    </Box>
  );
};

export default ClassifierLanding;
