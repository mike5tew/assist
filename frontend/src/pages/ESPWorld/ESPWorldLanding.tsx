import React from 'react';
import { Box, Container, Typography, Grid, Button, Stack, Paper, Chip } from '@mui/material';
import { Link } from 'react-router-dom';
import { AutoAwesome as MagicIcon, IntegrationInstructions as BridgeIcon, Psychology as CHISGIcon, Storage as WeaviateIcon } from '@mui/icons-material';
import ContactReveal from '../../components/ContactReveal';

const ESPWorldNav: React.FC = () => {
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
          <Button
            component={Link}
            to="/"
            variant="outlined"
            size="small"
            sx={{
              color: 'white',
              borderColor: 'rgba(255,255,255,0.6)',
              textTransform: 'none',
              bgcolor: 'rgba(15,23,42,0.4)',
              '&:hover': {
                borderColor: 'white',
                bgcolor: 'rgba(148,163,184,0.4)',
              },
            }}
          >
            Back to Portfolio
          </Button>
        </Stack>
      </Container>
    </Box>
  );
};

const ESPWorldLanding: React.FC = () => {
  return (
    <Box data-esp-marker="ESP_WORLD_MARKER_20260210">
      {/* Hero */}
      <Box
        sx={{
          position: 'relative',
          background: 'linear-gradient(135deg,#020617 0%, #0f172a 40%, #0e7490 100%)',
          color: 'white',
          pt: { xs: 10, md: 14 },
          pb: { xs: 8, md: 14 },
          overflow: 'hidden',
        }}
      >
        <ESPWorldNav />
        <Container maxWidth="lg">
          <Grid container spacing={6} alignItems="center">
            <Grid item xs={12} md={7}>
              <Typography variant="h2" component="h1" gutterBottom fontWeight="bold">ESP World</Typography>
              <Typography variant="h5" sx={{ opacity: 0.9, mb: 4 }}>A unified ecosystem of skills, semantic search, AI assistants and curriculum-aware tooling — bringing CHISG, Weaviate, the Skills Map and CTF standards together.</Typography>
              <Stack direction="row" spacing={2}>
                <Button component={Link} to="/" variant="outlined" sx={{ color: 'white', borderColor: 'rgba(255,255,255,0.25)' }}>VIEW DEMO</Button>
              </Stack>
            </Grid>
            <Grid item xs={12} md={5}>
              <Paper sx={{ p: 3, borderRadius: 2, bgcolor: 'rgba(255,255,255,0.04)', color: 'white' }}>
                <Typography variant="subtitle1" gutterBottom>Core Capabilities</Typography>
                <Stack spacing={1}>
                  <Chip icon={<WeaviateIcon />} label="Weaviate (Semantic Search)" sx={{ bgcolor: 'rgba(255,255,255,0.03)', color: 'white' }} />
                  <Chip icon={<CHISGIcon />} label="CHISG Skill Graphs" sx={{ bgcolor: 'rgba(255,255,255,0.03)', color: 'white' }} />
                  <Chip icon={<BridgeIcon />} label="Skills Map Integration" sx={{ bgcolor: 'rgba(255,255,255,0.03)', color: 'white' }} />
                  <Chip icon={<MagicIcon />} label="AI Assistants & Tooling" sx={{ bgcolor: 'rgba(255,255,255,0.03)', color: 'white' }} />
                </Stack>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* Features */}
      <Box sx={{ py: 8, bgcolor: 'rgba(15,23,42,0.02)' }}>
        <Container maxWidth="lg">
          <Typography variant="overline" color="text.secondary" textAlign="center" display="block">
            System Design
          </Typography>
          <Typography variant="h4" gutterBottom textAlign="center">The ESP World Architecture</Typography>
          <Grid container spacing={4} sx={{ mt: 2 }}>
            <Grid item xs={12} md={3}>
              <Paper sx={{ p: 3, height: '100%', borderRadius: 3, bgcolor: 'background.paper', boxShadow: '0 10px 30px rgba(15,23,42,0.06)' }}>
                <Typography variant="h6" gutterBottom>Weaviate (The Semantic Core)</Typography>
                <Typography variant="body2" color="text.secondary">
                  The <strong>weaviate-assist</strong> instance stores the deeper structural logic and philosphical "Project docs" of the ecosystem. It enables semantic discovery of curriculum gaps and the recursive logic that powers our AI assistants.
                </Typography>
              </Paper>
            </Grid>
            <Grid item xs={12} md={3}>
              <Paper sx={{ p: 3, height: '100%', borderRadius: 3, bgcolor: 'background.paper', boxShadow: '0 10px 30px rgba(15,23,42,0.06)' }}>
                <Typography variant="h6" gutterBottom>CHISG & Skills Map</Typography>
                <Typography variant="body2" color="text.secondary">
                  A deterministic skill graph that resolves the enigma of personal growth within a rigid curriculum. It tracks the "Sticker Album" of student achievements across five levels of complexity.
                </Typography>
              </Paper>
            </Grid>
            <Grid item xs={12} md={3}>
              <Paper sx={{ p: 3, height: '100%', borderRadius: 3, bgcolor: 'background.paper', boxShadow: '0 10px 30px rgba(15,23,42,0.06)' }}>
                <Typography variant="h6" gutterBottom>Microsoft Graph Fabric</Typography>
                <Typography variant="body2" color="text.secondary">
                  Deep integration with <strong>Microsoft Graph</strong> allows ESP World to ride on top of existing OneDrive/SharePoint file structures, keeping lesson plans, schemes of work and artefacts in place while adding live synchronisation from mobile capture back to the desktop knowledge fabric.
                </Typography>
              </Paper>
            </Grid>
            <Grid item xs={12} md={3}>
              <Paper sx={{ p: 3, height: '100%', borderRadius: 3, bgcolor: 'background.paper', boxShadow: '0 10px 30px rgba(15,23,42,0.06)' }}>
                <Typography variant="h6" gutterBottom>AI Assistants & CTF</Typography>
                <Typography variant="body2" color="text.secondary">
                  Using Content Transfer Format (CTF) to transmit ground-truth knowledge. These assistants apply "Vanquishing Principles" to turn learning friction into a sense of competence and mastery.
                </Typography>
              </Paper>
            </Grid>
          </Grid>

          <Box sx={{ mt: 4, textAlign: 'center' }}>
            <Button component={Link} to="/skillstree" variant="outlined" sx={{ mr: 2 }}>Open Skills Map</Button>
            <Button component={Link} to="/semantic-query" variant="contained">Explore Weaviate Tools</Button>
          </Box>
        </Container>
      </Box>

      {/* Philosophy Section */}
      <Box sx={{ py: 8, bgcolor: 'linear-gradient(180deg,#f9fafb 0%, #e5f3ff 100%)' }}>
        <Container maxWidth="lg">
          <Typography variant="overline" color="text.secondary" textAlign="center" display="block">
            Theory Layer
          </Typography>
          <Typography variant="h4" gutterBottom textAlign="center">Philosophical Foundation</Typography>
          <Typography variant="subtitle1" textAlign="center" sx={{ mb: 6, color: 'text.secondary', maxWidth: '800px', mx: 'auto' }}>
            ESP World is not just a collection of tools; it is a technological implementation of <strong>Hegelian Sublation</strong> (<em>Aufheben</em>) applied to educational enigmas.
          </Typography>

          <Grid container spacing={4}>
            <Grid item xs={12} md={6}>
              <Stack direction="row" spacing={1.5} alignItems="center" sx={{ mb: 1 }}>
                <CHISGIcon color="primary" />
                <Typography variant="h6" gutterBottom color="primary" sx={{ mb: 0 }}>
                  The Hegelian Dialectic in Education
                </Typography>
              </Stack>
              <Typography variant="body1" paragraph>
                Traditional education is often trapped in binary conflicts: Standardisation vs. Individualisation, Discipline vs. Relationship, or Data Security vs. Access. In Hegelian terms, these are <strong>Theses</strong> and <strong>Antitheses</strong> that create stagnation.
              </Typography>
              <Typography variant="body1">
                ESP World provides the <strong>Synthesis</strong>—the sublation where the opposing forces are not merely compromised, but preserved and lifted into a higher-order structure.
              </Typography>
            </Grid>
            <Grid item xs={12} md={6}>
              <Stack spacing={3}>
                <Paper sx={{ p: 3, borderLeft: '4px solid', borderColor: 'secondary.main', bgcolor: 'rgba(15,23,42,0.02)' }}>
                  <Typography variant="subtitle2" color="secondary" gutterBottom>Enigma: Standardisation vs. Individualisation</Typography>
                  <Typography variant="body2">
                    <strong>The Synthesis:</strong> A deterministic, standardized Skills Map (CHISG) that provides the common language required for radical, AI-driven individualisation of learning pathways.
                  </Typography>
                </Paper>
                <Paper sx={{ p: 3, borderLeft: '4px solid', borderColor: 'primary.main', bgcolor: 'rgba(15,23,42,0.02)' }}>
                  <Typography variant="subtitle2" color="primary" gutterBottom>Enigma: Discipline vs. Relationship</Typography>
                  <Typography variant="body2">
                    <strong>The Synthesis:</strong> Systemic focus management through digital tooling reduces the friction of "maintaining order," liberating the teacher to focus purely on high-value individual relationships and support.
                  </Typography>
                </Paper>
              </Stack>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* CTA */}
      <Box sx={{ py: 6, background: 'linear-gradient(90deg,#e6fffa, #bfe6ff)' }}>
        <Container maxWidth="md" sx={{ textAlign: 'center' }}>
          <Typography variant="h5" gutterBottom>Want a demo or to discuss integration?</Typography>
          <Typography variant="body2" sx={{ mb: 3 }}>I can show how Weaviate, CHISG and the Skills Map combine to deliver curriculum-aware assistants and reproducible knowledge graphs.</Typography>
          <Stack direction="row" spacing={2} justifyContent="center">
            <ContactReveal email="world@espthinkers.co.uk" label="Get in touch" />
            <Button variant="outlined" component={Link} to="/etp-profile">See demos</Button>
          </Stack>
        </Container>
      </Box>
    </Box>
  );
};

export default ESPWorldLanding;
