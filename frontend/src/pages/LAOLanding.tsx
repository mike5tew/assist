import React from 'react';
import {
  Box,
  Container,
  Typography,
  Button,
  Grid,
  Card,
  CardContent,
  Stack,
  Paper,
  alpha,
  Link as MuiLink,
} from '@mui/material';
import {
  Speed as SpeedIcon,
  Headphones as AudioIcon,
  TextSnippet as TextIcon,
  VpnKey as KeyIcon, // Keywords
  Repeat as RepeatIcon,
  TrendingUp as TrendingUpIcon, // Improvement
  Bolt as BoltIcon,
  AccessTime as AccessTimeIcon,
} from '@mui/icons-material';
import { Link } from 'react-router-dom';
import NavHeader from '../components/NavHeader';
import SEO from '../components/SEO';


const features = [
  {
    title: 'Keywords',
    icon: <KeyIcon fontSize="large" />,
    description: 'The core vocabulary required for marks. Stripped back to the absolute essentials.',
    color: '#e91e63',
  },
  {
    title: 'Speed Reading',
    icon: <SpeedIcon fontSize="large" />,
    description: 'Rapid serial visual presentation (RSVP) of content to maximize intake speed and focus.',
    color: '#2196f3',
  },
  {
    title: 'Audio Summaries',
    icon: <AudioIcon fontSize="large" />,
    description: 'Listen on the go. Turn "dead time" (bus rides, walking) into productive revision time.',
    color: '#ff9800',
  },
  {
    title: 'Plain Text',
    icon: <TextIcon fontSize="large" />,
    description: 'The full lesson content for deep reading when you have more than 5 minutes.',
    color: '#4caf50',
  },
];

const methodologySteps = [
  {
    title: 'Micro-Learning',
    icon: <AccessTimeIcon />,
    description: 'Practice for as little as 5 minutes. The "pain of starting" is removed when the commitment is so small.',
  },
  {
    title: 'Self-Assessment',
    icon: <TrendingUpIcon />,
    description: 'Students rate their own confidence immediately after a session. Simple, fast, honest tracking.',
  },
  {
    title: 'Adaptive Targeting',
    icon: <BoltIcon />,
    description: 'The system remembers weak spots and prioritizes them in the next session.',
  },
  {
    title: 'Repetition & Habit',
    icon: <RepeatIcon />,
    description: 'Frequency beats intensity. Small daily actions bring facts to the surface better than last-minute cramming.',
  },
];

const appStoreUrl = 'https://apps.apple.com/us/app/little-and-often-lao/id6758043113';


const LAOLanding: React.FC = () => {
  return (
    <Box>
      <SEO
        title="LAO — Adaptive GCSE Science Revision"
        description="Little and Often revision tool for GCSE Science. Speed reading, audio summaries, flashcards, and CHISG-driven gap analysis — designed to reduce revision resistance."
        path="/lao"
      />
      <NavHeader />
      {/* ── Hero Section ── */}
      <Box
        sx={{
          background: `linear-gradient(135deg, #1a237e 0%, #283593 100%)`, // Deep Blue for "Academic/Science" feel
          color: 'white',
          pt: { xs: 8, md: 12 },
          pb: { xs: 8, md: 12 },
          position: 'relative',
          overflow: 'hidden',
        }}
      >
        <Container maxWidth="lg">
          <Grid container spacing={6} alignItems="center">
            <Grid item xs={12} md={7}>
              <Box sx={{ mb: 2 }}>
                <Box
                    component="img"
                    src="LAOicon.png"
                    alt="DRB Ignite MAT"
                    sx={{
                      maxWidth: '100%',
                      width: { xs: '280px', md: '380px' },
                      height: 'auto',
                      alignSelf: 'flex-start',
                      borderRadius: 2,
                    }}
                  />
                
              </Box>
              
              <Typography variant="h2" component="h1" fontWeight="bold" gutterBottom sx={{ lineHeight: 1.1 }}>
                Little and Often
              </Typography>
              <Typography variant="h5" sx={{ opacity: 0.9, mb: 4, maxWidth: 600, lineHeight: 1.6 }}>
                Reducing the barriers to revision. How do you get the details a student needs to them as quickly as possible?
              </Typography>
              
              <Stack direction="row" spacing={2} alignItems="center">
                <Button
                  component="a"
                  href={appStoreUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  variant="contained"
                  color="secondary"
                >
                  Download on the App Store
                </Button>
                <Typography variant="caption" sx={{ opacity: 0.7 }}>
                  Live now on the App Store.
                </Typography>
              </Stack>
            </Grid>
            <Grid item xs={12} md={5} sx={{ display: { xs: 'none', md: 'block' } }}>
               {/* Abstract Phone/App Representation */}
               <Box 
                sx={{ 
                  position: 'relative',
                  width: '300px',
                  height: '600px',
                  bgcolor: '#fff',
                  borderRadius: '40px',
                  boxShadow: '0 20px 50px rgba(0,0,0,0.5)',
                  margin: '0 auto',
                  border: '8px solid #333',
                  overflow: 'hidden'
                }}
              >
                <Box sx={{ bgcolor: '#1a237e', height: '100%', color: 'white', p: 3, pt: 8 }}>
                    <Typography variant="h6" fontWeight="bold">Biology Paper 1</Typography>
                    <Typography variant="body2" sx={{ opacity: 0.7, mb: 4 }}>Cell Biology</Typography>
                    
                    <Paper sx={{ bgcolor: 'rgba(255,255,255,0.1)', p: 2, mb: 2, color: 'white' }}>
                        <Stack direction="row" alignItems="center" spacing={2}>
                            <KeyIcon sx={{ color: '#e91e63' }} />
                            <Box>
                                <Typography variant="subtitle2">Keywords</Typography>
                                <Typography variant="caption" sx={{ opacity: 0.7 }}>2 mins</Typography>
                            </Box>
                        </Stack>
                    </Paper>

                    <Paper sx={{ bgcolor: 'rgba(255,255,255,0.1)', p: 2, mb: 2, color: 'white' }}>
                        <Stack direction="row" alignItems="center" spacing={2}>
                            <SpeedIcon sx={{ color: '#2196f3' }} />
                            <Box>
                                <Typography variant="subtitle2">Speed Read</Typography>
                                <Typography variant="caption" sx={{ opacity: 0.7 }}>1 min</Typography>
                            </Box>
                        </Stack>
                    </Paper>

                     <Paper sx={{ bgcolor: 'rgba(255,255,255,0.1)', p: 2, mb: 2, color: 'white' }}>
                        <Stack direction="row" alignItems="center" spacing={2}>
                            <AudioIcon sx={{ color: '#ff9800' }} />
                            <Box>
                                <Typography variant="subtitle2">Audio Summary</Typography>
                                <Typography variant="caption" sx={{ opacity: 0.7 }}>3 mins</Typography>
                            </Box>
                        </Stack>
                    </Paper>
                </Box>
              </Box>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* ── The Problem Section ── */}
      <Box sx={{ py: 10, bgcolor: '#f5f7f9' }}>
        <Container maxWidth="md">
            <Box sx={{ textAlign: 'center', mb: 6 }}>
                <Typography variant="overline" color="error" fontWeight="bold" sx={{ letterSpacing: 2 }}>
                    THE FRICTION PROBLEM
                </Typography>
                <Typography variant="h3" fontWeight="bold" gutterBottom sx={{ mt: 1 }}>
                    "Too Little, Too Late"
                </Typography>
                <Typography variant="body1" color="text.secondary" sx={{ fontSize: '1.2rem' }}>
                    Students know they need to start revising earlier. But the sheer friction of getting started — 
                    finding the notes, planning the schedule, committing to an hour — paralysis them.
                </Typography>
            </Box>

            <Grid container spacing={4}>
                <Grid item xs={12} md={6}>
                    <Paper elevation={0} sx={{ p: 4, height: '100%', borderLeft: '4px solid #f44336' }}>
                        <Typography variant="h6" fontWeight="bold" gutterBottom>The Setup Pain</Typography>
                        <Typography variant="body2" color="text.secondary">
                            By the time a student has cleared their desk, found their book, and opened the right page, 
                            they've already spent 15 minutes of willpower.
                        </Typography>
                    </Paper>
                </Grid>
                <Grid item xs={12} md={6}>
                    <Paper elevation={0} sx={{ p: 4, height: '100%', borderLeft: '4px solid #f44336' }}>
                        <Typography variant="h6" fontWeight="bold" gutterBottom>The Schedule Trap</Typography>
                        <Typography variant="body2" color="text.secondary">
                             Big planned sessions feel daunting. When life gets busy, the "1-hour session" is the first thing dropped.
                             Consistency breaks, and the habit never forms.
                        </Typography>
                    </Paper>
                </Grid>
            </Grid>
        </Container>
      </Box>

      {/* ── The Solution: 4 Forms ── */}
      <Box sx={{ py: 10 }}>
        <Container maxWidth="lg">
            <Typography variant="h4" textAlign="center" fontWeight="bold" gutterBottom>
                The 4 Forms of Content
            </Typography>
            <Typography variant="body1" textAlign="center" color="text.secondary" sx={{ mb: 8, maxWidth: 700, mx: 'auto' }}>
                How do we get the necessary details to the student as quickly as possible? We strip away the overhead and deliver content in four distinct modes.
            </Typography>

            <Grid container spacing={4}>
                {features.map((f) => (
                    <Grid item xs={12} sm={6} md={3} key={f.title}>
                        <Card 
                            elevation={0}
                            sx={{
                                height: '100%',
                                textAlign: 'center',
                                border: `1px solid ${alpha(f.color, 0.2)}`,
                                '&:hover': { transform: 'translateY(-5px)', boxShadow: 4, transition: 'all 0.3s' }
                            }}
                        >
                            <CardContent sx={{ p: 4 }}>
                                <Box sx={{ color: f.color, mb: 2 }}>{f.icon}</Box>
                                <Typography variant="h6" fontWeight="bold" gutterBottom>{f.title}</Typography>
                                <Typography variant="body2" color="text.secondary">{f.description}</Typography>
                            </CardContent>
                        </Card>
                    </Grid>
                ))}
            </Grid>
        </Container>
      </Box>

      {/* ── Methodology ── */}
      <Box sx={{ py: 10, bgcolor: '#1a237e', color: 'white' }}>
        <Container maxWidth="lg">
            <Grid container spacing={6} alignItems="center">
                <Grid item xs={12} md={5}>
                    <Typography variant="h3" fontWeight="bold" gutterBottom>
                        Habit Over Intensity
                    </Typography>
                    <Typography variant="body1" sx={{ opacity: 0.8, mb: 4, fontSize: '1.1rem' }}>
                        The aim is to practice for as little as five minutes at a time. This simple change repurposes 
                        "dead time" into retrieval practice, bringing facts and stories to the surface through sheer frequency.
                    </Typography>
                    <Button 
                      component="a"
                      href={appStoreUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                        variant="outlined" 
                        color="inherit" 
                        size="large"
                        sx={{ borderColor: 'rgba(255,255,255,0.5)' }}
                    >
                      View on the App Store
                    </Button>
                </Grid>
                <Grid item xs={12} md={7}>
                    <Grid container spacing={3}>
                        {methodologySteps.map((step) => (
                            <Grid item xs={12} sm={6} key={step.title}>
                                <Paper sx={{ p: 3, bgcolor: 'rgba(255,255,255,0.05)', color: 'white' }} elevation={0}>
                                    <Box sx={{ display: 'flex', gap: 2, mb: 1 }}>
                                        <Box sx={{ color: '#64ffda' }}>{step.icon}</Box>
                                        <Typography variant="subtitle1" fontWeight="bold">{step.title}</Typography>
                                    </Box>
                                    <Typography variant="body2" sx={{ opacity: 0.7 }}>
                                        {step.description}
                                    </Typography>
                                </Paper>
                            </Grid>
                        ))}
                    </Grid>
                </Grid>
            </Grid>
        </Container>
      </Box>

      {/* ── Footer ── */}
      <Box sx={{ py: 4, textAlign: 'center', borderTop: '1px solid #eee' }}>
        <Typography variant="body2" color="text.secondary">
          Part of the ESP Thinking Portfolio — Little and Often (LAO)
        </Typography>
        <Stack direction="row" spacing={2} justifyContent="center" sx={{ mt: 1 }}>
          <MuiLink component={Link} to="/lao/privacy" variant="body2" color="text.secondary">
            Privacy Policy
          </MuiLink>
          <MuiLink component={Link} to="/lao/terms" variant="body2" color="text.secondary">
            Terms of Use
          </MuiLink>
        </Stack>
      </Box>
    </Box>
  );
};

export default LAOLanding;
