import React from 'react';
import { Box, Button, Card, CardActions, CardContent, Chip, Container, Grid, Paper, Stack, Typography, alpha, useTheme } from '@mui/material';
import { Link } from 'react-router-dom';
import { Psychology as CHISGIcon, Storage as WeaviateIcon, Hub as GraphIcon, Verified as VerifiedIcon, TravelExplore as ExploreIcon, ListAlt as ListIcon } from '@mui/icons-material';
import ContactReveal from '../../components/ContactReveal';
import SEO from '../../components/SEO';

const CHISGNav: React.FC = () => {
  return (
    <Box
      sx={{
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        zIndex: 10,
        py: 2,
        px: 3,
      }}
    >
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
                '&:hover': {
                  borderColor: 'white',
                  bgcolor: 'rgba(148,163,184,0.35)',
                },
              }}
            >
              Back to Portfolio
            </Button>
          </Stack>
        </Stack>
      </Container>
    </Box>
  );
};

const CHISGLanding: React.FC = () => {
  const theme = useTheme();

  return (
    <Box>
      <SEO
        title="CHISG — Contextualised Hierarchical Iterative Semantic Groupings"
        description="CHISG is the knowledge-graph methodology that structures educational content, maps semantic relationships between skills, and eliminates AI hallucinations through grounded knowledge."
        path="/chisg"
      />
      {/* Hero */}
      <Box
        sx={{
          position: 'relative',
          background: 'linear-gradient(135deg,#020617 0%, #0f172a 45%, #4f46e5 100%)',
          color: 'white',
          pt: { xs: 10, md: 14 },
          pb: { xs: 8, md: 12 },
          overflow: 'hidden',
        }}
      >
        <CHISGNav />
        <Container maxWidth="lg">
          <Grid container spacing={6} alignItems="center">
            <Grid item xs={12} md={7}>
              <Typography variant="overline" sx={{ opacity: 0.85 }}>
                Information Quality • Semantic Grounding • Safe Connection
              </Typography>
              <Typography variant="h2" component="h1" fontWeight="bold" gutterBottom>
                CHISG
              </Typography>
              <Typography variant="h5" sx={{ opacity: 0.92, mb: 3 }}>
                A disciplined way to model knowledge so that AI can connect ideas without inventing links that shouldn’t exist.
              </Typography>
              <Typography variant="body1" sx={{ opacity: 0.9, maxWidth: 760 }}>
                As models get stronger, they also get better at “joining dots” — sometimes like a really intelligent conspiracy theorist:
                pattern‑matching across weak signals, implying causation, and confidently weaving narratives.
                CHISG is about defining what <strong>information quality</strong> means, and how ideas connect in a mechanism that is as efficient as possible.
              </Typography>

              <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} sx={{ mt: 4 }}>
                <Button component={Link} to="/semantic-query" variant="contained" color="secondary" size="large">
                  Open CHISG Demo (Semantic Query)
                </Button>
                <Button component={Link} to="/chisg/graph" variant="outlined" size="large" sx={{ color: 'white', borderColor: 'rgba(255,255,255,0.35)' }}>
                  Explore Knowledge Graph
                </Button>
                <Button component={Link} to="/ai-chat" variant="outlined" size="large" sx={{ color: 'white', borderColor: 'rgba(255,255,255,0.35)' }}>
                  Try Grounded Chat
                </Button>
              </Stack>
            </Grid>

            <Grid item xs={12} md={5}>
              <Paper
                sx={{
                  p: 3,
                  borderRadius: 3,
                  bgcolor: 'rgba(255,255,255,0.04)',
                  border: '1px solid rgba(148,163,184,0.22)',
                }}
              >
                <Typography variant="subtitle1" gutterBottom>
                  What CHISG Optimizes For
                </Typography>
                <Stack spacing={1.25}>
                  <Chip
                    icon={<VerifiedIcon />}
                    label="High‑precision connections (not maximum connections)"
                    sx={{ bgcolor: 'rgba(255,255,255,0.03)', color: 'white' }}
                  />
                  <Chip
                    icon={<GraphIcon />}
                    label="Explicit relationships that can be inspected"
                    sx={{ bgcolor: 'rgba(255,255,255,0.03)', color: 'white' }}
                  />
                  <Chip
                    icon={<WeaviateIcon />}
                    label="Semantic discovery backed by a curated knowledge base"
                    sx={{ bgcolor: 'rgba(255,255,255,0.03)', color: 'white' }}
                  />
                  <Chip
                    icon={<CHISGIcon />}
                    label="Grounded synthesis (reduce hallucinations and slop)"
                    sx={{ bgcolor: 'rgba(255,255,255,0.03)', color: 'white' }}
                  />
                </Stack>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* Why */}
      <Box sx={{ py: 8, bgcolor: 'rgba(15,23,42,0.02)' }}>
        <Container maxWidth="lg">
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            The Problem: Powerful Models, Weak Structure
          </Typography>
          <Typography variant="body1" color="text.secondary" textAlign="center" sx={{ mb: 4, maxWidth: 900, mx: 'auto' }}>
            When the underlying information is messy, ambiguous, or over‑connected, a better model doesn’t fix the problem — it amplifies it.
            The result is confident outputs built on fragile assumptions.
          </Typography>

          <Grid container spacing={3}>
            {[
              {
                title: 'False Cohesion',
                body: 'The model finds a “story” that feels coherent, even when the evidence is thin.',
              },
              {
                title: 'Over‑Connection',
                body: 'Everything becomes related to everything else, so signal gets drowned in noise.',
              },
              {
                title: 'Uninspectable Links',
                body: 'You can’t tell which claim came from where, so errors become untraceable.',
              },
              {
                title: 'Efficiency vs. Accuracy',
                body: 'Fast retrieval is meaningless if the retrieved context is low‑quality or wrongly connected.',
              },
            ].map((item) => (
              <Grid item xs={12} sm={6} md={3} key={item.title}>
                <Card elevation={2} sx={{ height: '100%', borderTop: `4px solid ${theme.palette.primary.main}` }}>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      {item.title}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {item.body}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* What it is (without methodology) */}
      <Box sx={{ py: 8 }}>
        <Container maxWidth="lg">
          <Grid container spacing={4} alignItems="stretch">
            <Grid item xs={12} md={6}>
              <Typography variant="h4" component="h2" fontWeight="bold" gutterBottom>
                What CHISG Is (and Isn’t)
              </Typography>
              <Typography variant="body1" paragraph>
                CHISG is a <strong>quality layer</strong> for knowledge: it focuses on what a concept <em>is</em>, what it is <em>not</em>, and how it can legitimately relate to other concepts.
                The goal is to make AI outputs more dependable by constraining which connections are allowed.
              </Typography>
              <Typography variant="body1" paragraph>
                It’s not a “prompt trick” and it’s not just a vector database.
                It’s a way of keeping semantic search and synthesis anchored to a curated, inspectable structure.
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Note: the exact construction rules and procedures are intentionally not published here.
                This page is about what the system enables, not how it is built.
              </Typography>
            </Grid>

            <Grid item xs={12} md={6}>
              <Paper
                elevation={0}
                sx={{
                  p: 3,
                  borderRadius: 3,
                  border: '1px solid',
                  borderColor: 'divider',
                  bgcolor: alpha(theme.palette.background.paper, 0.9),
                }}
              >
                <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
                  A Practical Definition of Information Quality
                </Typography>
                <Stack spacing={1.5}>
                  {[
                    { title: 'Bounded meaning', desc: 'Terms have scoped definitions that prevent “drift”.' },
                    { title: 'Typed relationships', desc: 'Connections have explicit kinds (not vague similarity).' },
                    { title: 'Traceable claims', desc: 'You can inspect why a link exists and what supports it.' },
                    { title: 'Minimal connections', desc: 'Only the most informative links survive; the rest are noise.' },
                  ].map((row) => (
                    <Box key={row.title}>
                      <Typography variant="body2" fontWeight="bold">
                        {row.title}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        {row.desc}
                      </Typography>
                    </Box>
                  ))}
                </Stack>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* Demo links */}
      <Box sx={{ py: 8, bgcolor: 'grey.50' }}>
        <Container maxWidth="lg">
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            See It Working
          </Typography>
          <Typography variant="body1" color="text.secondary" textAlign="center" sx={{ mb: 4, maxWidth: 900, mx: 'auto' }}>
            These demos are powered by the <strong>weaviate-assist</strong> knowledge base and related tooling.
            They’re designed to show grounded discovery and safe connection — without exposing the underlying build methodology.
          </Typography>

          <Grid container spacing={3}>
            <Grid item xs={12} md={4}>
              <Card elevation={2} sx={{ height: '100%', borderTop: '4px solid #4f46e5' }}>
                <CardContent>
                  <Stack direction="row" spacing={1.5} alignItems="center" sx={{ mb: 1 }}>
                    <ExploreIcon sx={{ color: '#4f46e5' }} />
                    <Typography variant="h6">Semantic Query</Typography>
                  </Stack>
                  <Typography variant="body2" color="text.secondary">
                    Search skills and concepts; compare exact matches vs semantic matches; view a grounded synthesis.
                  </Typography>
                </CardContent>
                <CardActions sx={{ px: 2, pb: 2 }}>
                  <Button component={Link} to="/semantic-query" variant="contained" size="small">
                    Open Demo
                  </Button>
                </CardActions>
              </Card>
            </Grid>

            <Grid item xs={12} md={4}>
              <Card elevation={2} sx={{ height: '100%', borderTop: '4px solid #0ea5e9' }}>
                <CardContent>
                  <Stack direction="row" spacing={1.5} alignItems="center" sx={{ mb: 1 }}>
                    <WeaviateIcon sx={{ color: '#0ea5e9' }} />
                    <Typography variant="h6">Grounded Chat</Typography>
                  </Stack>
                  <Typography variant="body2" color="text.secondary">
                    Ask questions against the knowledge base and get answers that prioritize grounded context.
                  </Typography>
                </CardContent>
                <CardActions sx={{ px: 2, pb: 2 }}>
                  <Button component={Link} to="/ai-chat" variant="outlined" size="small">
                    Open Chat
                  </Button>
                </CardActions>
              </Card>
            </Grid>

            <Grid item xs={12} md={4}>
              <Card elevation={2} sx={{ height: '100%', borderTop: '4px solid #16a34a' }}>
                <CardContent>
                  <Stack direction="row" spacing={1.5} alignItems="center" sx={{ mb: 1 }}>
                    <GraphIcon sx={{ color: '#16a34a' }} />
                    <Typography variant="h6">Skills Map</Typography>
                  </Stack>
                  <Typography variant="body2" color="text.secondary">
                    Visualize a deterministic skill graph and explore how concepts connect in a structured way.
                  </Typography>
                </CardContent>
                <CardActions sx={{ px: 2, pb: 2 }}>
                  <Button component={Link} to="/skillstree" variant="outlined" size="small">
                    Open Map
                  </Button>
                </CardActions>
              </Card>
            </Grid>

            <Grid item xs={12} md={4}>
              <Card elevation={2} sx={{ height: '100%', borderTop: '4px solid #22d3ee' }}>
                <CardContent>
                  <Stack direction="row" spacing={1.5} alignItems="center" sx={{ mb: 1 }}>
                    <GraphIcon sx={{ color: '#22d3ee' }} />
                    <Typography variant="h6">Knowledge Graph</Typography>
                  </Stack>
                  <Typography variant="body2" color="text.secondary">
                    Explore the CHISG semantic graph — concepts as nodes, relations as edges, coloured by domain.
                  </Typography>
                </CardContent>
                <CardActions sx={{ px: 2, pb: 2 }}>
                  <Button component={Link} to="/chisg/graph" variant="contained" size="small" sx={{ bgcolor: '#22d3ee', color: '#0a0f1e', '&:hover': { bgcolor: '#06b6d4' } }}>
                    Explore Graph
                  </Button>
                  <Button component={Link} to="/ntm" variant="outlined" size="small" sx={{ color: '#22d3ee', borderColor: 'rgba(34,211,238,0.4)', '&:hover': { borderColor: '#22d3ee' } }}>
                    NTM Research
                  </Button>
                </CardActions>
              </Card>
            </Grid>

            <Grid item xs={12} md={4}>
              <Card elevation={2} sx={{ height: '100%', borderTop: '4px solid #16a34a' }}>
                <CardContent>
                  <Stack direction="row" spacing={1.5} alignItems="center" sx={{ mb: 1 }}>
                    <ListIcon sx={{ color: '#16a34a' }} />
                    <Typography variant="h6">Knowledge Base</Typography>
                  </Stack>
                  <Typography variant="body2" color="text.secondary">
                    Browse the extracted academic papers, view validated semantic triples, and explore the link builder.
                  </Typography>
                </CardContent>
                <CardActions sx={{ px: 2, pb: 2 }}>
                  <Button component={Link} to="/chisg/papers" variant="contained" size="small" sx={{ bgcolor: '#16a34a', '&:hover': { bgcolor: '#15803d' } }}>
                    View Papers
                  </Button>
                </CardActions>
              </Card>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* ESP Research Pilot call-out */}
      <Box sx={{ py: 5, bgcolor: '#0f172a' }}>
        <Container maxWidth="md">
          <Stack direction={{ xs: 'column', sm: 'row' }} alignItems={{ sm: 'center' }} spacing={3}
            sx={{ p: 3, border: '1px solid rgba(99,102,241,0.35)', borderRadius: 2,
              background: 'linear-gradient(135deg, rgba(99,102,241,0.08) 0%, rgba(49,46,129,0.2) 100%)' }}>
            <Box sx={{ flex: 1 }}>
              <Typography variant="overline" sx={{ color: '#818cf8', fontWeight: 700, letterSpacing: 2 }}>
                Research Pilot
              </Typography>
              <Typography variant="h6" fontWeight={700} color="white" mt={0.5} mb={0.75}>
                CHISG applied to academic research
              </Typography>
              <Typography variant="body2" sx={{ color: 'rgba(255,255,255,0.55)', lineHeight: 1.7 }}>
                ESP uses the CHISG methodology to extract causal relationships from scientific papers
                and build a personal, queryable knowledge oracle for PhD researchers and academics.
                We're looking for pilot participants.
              </Typography>
            </Box>
            <Button
              component={Link} to="/esp-pilot"
              variant="contained"
              sx={{ flexShrink: 0, bgcolor: '#6366f1', '&:hover': { bgcolor: '#4f46e5' },
                fontWeight: 700, textTransform: 'none', borderRadius: 2 }}
            >
              Learn more →
            </Button>
          </Stack>
        </Container>
      </Box>

      {/* CTA */}
      <Box sx={{ py: 6, background: 'linear-gradient(90deg,#eef2ff, #dbeafe)' }}>
        <Container maxWidth="md" sx={{ textAlign: 'center' }}>
          <Typography variant="h5" gutterBottom>
            Want the “why” without the proprietary “how”?
          </Typography>
          <Typography variant="body2" sx={{ mb: 3 }}>
            I can walk through the demo flows, what is being validated, and how the structure prevents bad dot‑joining.
          </Typography>
          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} justifyContent="center">
            <ContactReveal email="world@espthinkers.co.uk" label="Get in touch" />
            <Button variant="outlined" component={Link} to="/semantic-query">
              Open Semantic Query
            </Button>
          </Stack>
        </Container>
      </Box>
    </Box>
  );
};

export default CHISGLanding;
