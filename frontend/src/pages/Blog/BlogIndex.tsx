import React from 'react';
import {
  Box,
  Container,
  Typography,
  Card,
  CardContent,
  CardActionArea,
  Grid,
  Stack,
  Chip,
  Button,
  alpha,
  useTheme,
} from '@mui/material';
import {
  Home as HomeIcon,
  Groups as GroupsIcon,
  SwapHoriz as SwapIcon,
  Bolt as BoltIcon,
  Shield as ShieldIcon,
  Favorite as FavoriteIcon,
  Casino as CasinoIcon,
  Balance as BalanceIcon,
  Psychology as PsychologyIcon,
  GridView as GridViewIcon,
  PanTool as PanToolIcon,
  FlightTakeoff as FlightTakeoffIcon,
  LocalHospital as LocalHospitalIcon,
  Gavel as GavelIcon,
} from '@mui/icons-material';
import { Link } from 'react-router-dom';
import SEO from '../../components/SEO';

interface BlogPostMeta {
  slug: string;
  title: string;
  subtitle: string;
  date: string;
  series: string;
  icon: React.ReactNode;
  color: string;
  excerpt: string;
  isLive?: boolean;
}

const posts: BlogPostMeta[] = [
  {
    slug: 'social-gravity',
    title: 'Social Gravity',
    subtitle: 'ETP Post 1 of 12',
    date: 'March 2026',
    series: 'ETP Deep-Dive',
    icon: <GroupsIcon />,
    color: '#1565c0',
    excerpt:
      'Your brain has a battery. Social interaction either charges it or drains it. This isn\'t introversion vs. extroversion — that\'s a blunt instrument.',
  },
  {
    slug: 'energy-directionality',
    title: 'Energy Directionality',
    subtitle: 'ETP Post 2 of 12',
    date: 'March 2026',
    series: 'ETP Deep-Dive',
    icon: <SwapIcon />,
    color: '#6a1b9a',
    excerpt:
      'Your brain is a processor. But it runs in one of two directions. Social Gravity is about who charges you. Energy Directionality is about where your processing happens.',
  },
  {
    slug: 'voltage-sensitivity',
    title: 'Voltage Sensitivity',
    subtitle: 'ETP Post 3 of 12',
    date: 'March 2026',
    series: 'ETP Deep-Dive',
    icon: <BoltIcon />,
    color: '#f57c00',
    excerpt:
      'Your circuit breaker has a setting. You didn\'t choose it. How much emotional current your system can carry before it trips the breaker.',
  },
  {
    slug: 'threat-response',
    title: 'Threat Response',
    subtitle: 'ETP Post 4 of 12',
    date: 'March 2026',
    series: 'ETP Deep-Dive',
    icon: <ShieldIcon />,
    color: '#c62828',
    excerpt:
      'Your system just overloaded. What happens next? We treat Freeze as good behaviour and Fight as bad behaviour. In classrooms, this is catastrophic.',
  },
  {
    slug: 'care-response',
    title: 'Care Response',
    subtitle: 'ETP Post 5 of 12',
    date: 'April 2026',
    series: 'ETP Deep-Dive',
    icon: <FavoriteIcon />,
    color: '#ad1457',
    excerpt:
      'Whose pain do you feel first? Yours — or everyone else\'s? We moralise Care Response harder than almost any other spectrum.',
  },
  {
    slug: 'risk-tolerance',
    title: 'Risk Tolerance',
    subtitle: 'ETP Post 6 of 12',
    date: 'May 2026',
    series: 'ETP Deep-Dive',
    icon: <CasinoIcon />,
    color: '#2e7d32',
    excerpt:
      'Part of you wants to leap. Part of you won\'t let you. We have called those two parts angel and demon for thousands of years. We were wrong.',
  },
  {
    slug: 'integrity-logic',
    title: 'Integrity Logic',
    subtitle: 'ETP Post 7 of 12',
    date: 'May 2026',
    series: 'ETP Deep-Dive',
    icon: <BalanceIcon />,
    color: '#00695c',
    excerpt:
      'If there is no angel and no demon, then why does choosing the right thing sometimes feel like abandoning yourself?',
  },
  {
    slug: 'mirror-neuron-tuning',
    title: 'Mirror Neuron Tuning',
    subtitle: 'ETP Post 8 of 12',
    date: 'May 2026',
    series: 'ETP Deep-Dive',
    icon: <PsychologyIcon />,
    color: '#4527a0',
    excerpt:
      'What happens when the feeling you\'re carrying isn\'t yours? Mirror neurons simulate other people\'s emotional states in your own nervous system — not observe them. Simulate them.',
  },
  {
    slug: 'orderliness',
    title: 'Orderliness',
    subtitle: 'ETP Post 9 of 12',
    date: 'May 2026',
    series: 'ETP Deep-Dive',
    icon: <GridViewIcon />,
    color: '#455a64',
    excerpt:
      'The student who can\'t start until their desk is clear isn\'t procrastinating. Structure is not the enemy of creativity. It is its engine.',
  },
  {
    slug: 'pilot-strength',
    title: 'Pilot Strength',
    subtitle: 'ETP Post 10 of 12',
    date: 'June 2026',
    series: 'ETP Deep-Dive',
    icon: <FlightTakeoffIcon />,
    color: '#5d4037',
    excerpt:
      'When you know exactly what to do but cannot make yourself do it, the issue is usually capacity, not character. This post maps the knowing-to-doing gap.',
    isLive: true,
  },
  {
    slug: 'current-load',
    title: 'Current Load',
    subtitle: 'ETP Post 11 of 12',
    date: 'June 2026',
    series: 'ETP Deep-Dive',
    icon: <LocalHospitalIcon />,
    color: '#546e7a',
    excerpt:
      'The system under strain reads everything differently. Current Load separates state from wiring and explains why behaviour shifts as buffer runs out.',
    isLive: true,
  },
  {
    slug: 'libido',
    title: 'Libido',
    subtitle: 'ETP Post 12 of 12',
    date: 'June 2026',
    series: 'ETP Deep-Dive',
    icon: <PanToolIcon />,
    color: '#00838f',
    excerpt:
      'The one drive every school is afraid to name — and the one running loudest in the room. The brakes were teachable years earlier. We just flinched. A safeguarding spectrum first.',
    isLive: true,
  },
];

const appliedPosts: BlogPostMeta[] = [
  {
    slug: 'the-governor',
    title: 'The Governor',
    subtitle: 'Applied ETP · Evidence & Disagreement',
    date: 'June 2026',
    series: 'Applied ETP',
    icon: <GavelIcon />,
    color: '#283593',
    excerpt:
      'An AI role that audits both sides of a disagreement — checking logic, evidence, and sourcing — and is built never to declare a winner.',
    isLive: true,
  },
];

const BlogIndex: React.FC = () => {
  const theme = useTheme();

  return (
    <Box>
      <SEO
        title="Blog — ESP Thinking"
        description="Articles on Emergent Tendency Profiles, emotional readiness, and the architecture of learning."
        path="/blog"
      />

      {/* Hero */}
      <Box
        sx={{
          background: `linear-gradient(135deg, ${theme.palette.primary.dark} 0%, ${theme.palette.primary.main} 50%, ${theme.palette.secondary.main} 100%)`,
          borderTop: '6px solid #1a365d',
          color: 'white',
          pt: { xs: 4, md: 8 },
          pb: { xs: 6, md: 10 },
          px: 2,
          position: 'relative',
          overflow: 'hidden',
          '&::after': {
            content: '""',
            position: 'absolute',
            bottom: 0,
            left: 0,
            right: 0,
            height: '60px',
            background: 'linear-gradient(to top right, #f5f7f9 50%, transparent 50%)',
          },
        }}
      >
        <Container maxWidth="lg">
          <Stack spacing={2}>
            <Button
              component={Link}
              to="/"
              startIcon={<HomeIcon />}
              sx={{
                color: 'rgba(255,255,255,0.8)',
                textTransform: 'none',
                alignSelf: 'flex-start',
                '&:hover': { color: 'white', bgcolor: 'rgba(255,255,255,0.1)' },
              }}
            >
              Back to Portfolio
            </Button>
            <Box
              component="img"
              src={`${process.env.PUBLIC_URL}/ESPLogoLong.png`}
              alt="ESP Thinking"
              sx={{ height: { xs: 80, md: 120 }, width: 'auto', objectFit: 'contain' }}
            />
            <Typography
              variant="h2"
              component="h1"
              fontWeight="bold"
              sx={{ fontSize: { xs: '1.8rem', md: '2.8rem' } }}
            >
              Blog
            </Typography>
            <Typography variant="h5" sx={{ opacity: 0.92, maxWidth: 600 }}>
              Deep-dives into the 12 biological spectra that shape how we learn, connect, and grow.
            </Typography>
          </Stack>
        </Container>
      </Box>

      {/* ETP Series */}
      <Box sx={{ py: { xs: 4, md: 6 }, bgcolor: '#f5f7f9' }}>
        <Container maxWidth="lg">
          <Typography variant="overline" sx={{ letterSpacing: 2, color: 'text.secondary' }}>
            Series
          </Typography>
          <Typography variant="h4" fontWeight="bold" gutterBottom>
            ETP Deep-Dive: 12 Spectra
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 4, maxWidth: 700 }}>
            One post per spectrum. Each one explains what it is, why it matters for teachers,
            and how AI could respond to it. No labels — settings.
          </Typography>

          <Grid container spacing={3}>
            {posts.map((post) => (
              <Grid item xs={12} sm={6} md={4} key={post.slug}>
                <Card
                  sx={{
                    height: '100%',
                    borderTop: `4px solid ${post.color}`,
                    transition: 'transform 0.2s, box-shadow 0.2s',
                    opacity: post.isLive === false ? 0.92 : 1,
                    '&:hover': {
                      transform: post.isLive === false ? 'none' : 'translateY(-4px)',
                      boxShadow: post.isLive === false ? undefined : `0 8px 24px ${alpha(post.color, 0.2)}`,
                    },
                  }}
                >
                  <CardActionArea
                    component={post.isLive === false ? 'div' : Link}
                    to={post.isLive === false ? undefined : `/blog/${post.slug}`}
                    sx={{
                      height: '100%',
                      display: 'flex',
                      flexDirection: 'column',
                      alignItems: 'stretch',
                      cursor: post.isLive === false ? 'default' : 'pointer',
                    }}
                    disableRipple={post.isLive === false}
                  >
                    <CardContent sx={{ flexGrow: 1 }}>
                      <Stack spacing={1.5}>
                        <Stack direction="row" spacing={1} alignItems="center">
                          <Box sx={{ color: post.color }}>{post.icon}</Box>
                          <Chip
                            label={post.subtitle}
                            size="small"
                            sx={{
                              bgcolor: alpha(post.color, 0.1),
                              color: post.color,
                              fontWeight: 600,
                            }}
                          />
                          {post.isLive === false && (
                            <Chip
                              label="Coming soon"
                              size="small"
                              variant="outlined"
                              sx={{ borderColor: alpha(post.color, 0.35), color: post.color }}
                            />
                          )}
                        </Stack>
                        <Typography variant="h5" fontWeight="bold">
                          {post.title}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" sx={{ lineHeight: 1.7 }}>
                          {post.excerpt}
                        </Typography>
                        <Typography variant="caption" color="text.disabled">
                          {post.date}
                        </Typography>
                      </Stack>
                    </CardContent>
                  </CardActionArea>
                </Card>
              </Grid>
            ))}
          </Grid>

          {/* Coming Soon */}
          <Box
            sx={{
              mt: 4,
              p: 3,
              bgcolor: alpha(theme.palette.primary.main, 0.04),
              borderRadius: 2,
              border: `1px dashed ${alpha(theme.palette.primary.main, 0.2)}`,
            }}
          >
            <Typography variant="body2" color="text.secondary">
              <strong>The full set:</strong> all twelve spectra are now live. Each one stands alone —
              read the one that sounds like someone you know.
            </Typography>
          </Box>
        </Container>
      </Box>

      {/* Applied ETP */}
      <Box sx={{ py: { xs: 4, md: 6 }, bgcolor: 'white' }}>
        <Container maxWidth="lg">
          <Typography variant="overline" sx={{ letterSpacing: 2, color: 'text.secondary' }}>
            Series
          </Typography>
          <Typography variant="h4" fontWeight="bold" gutterBottom>
            Applied ETP
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 4, maxWidth: 700 }}>
            Where the twelve spectra get pointed outward — at arguments, evidence, and disagreement,
            not just at one nervous system trying to regulate itself.
          </Typography>

          <Grid container spacing={3}>
            {appliedPosts.map((post) => (
              <Grid item xs={12} sm={6} md={4} key={post.slug}>
                <Card
                  sx={{
                    height: '100%',
                    borderTop: `4px solid ${post.color}`,
                    transition: 'transform 0.2s, box-shadow 0.2s',
                    opacity: post.isLive === false ? 0.92 : 1,
                    '&:hover': {
                      transform: post.isLive === false ? 'none' : 'translateY(-4px)',
                      boxShadow: post.isLive === false ? undefined : `0 8px 24px ${alpha(post.color, 0.2)}`,
                    },
                  }}
                >
                  <CardActionArea
                    component={post.isLive === false ? 'div' : Link}
                    to={post.isLive === false ? undefined : `/blog/${post.slug}`}
                    sx={{
                      height: '100%',
                      display: 'flex',
                      flexDirection: 'column',
                      alignItems: 'stretch',
                      cursor: post.isLive === false ? 'default' : 'pointer',
                    }}
                    disableRipple={post.isLive === false}
                  >
                    <CardContent sx={{ flexGrow: 1 }}>
                      <Stack spacing={1.5}>
                        <Stack direction="row" spacing={1} alignItems="center">
                          <Box sx={{ color: post.color }}>{post.icon}</Box>
                          <Chip
                            label={post.subtitle}
                            size="small"
                            sx={{
                              bgcolor: alpha(post.color, 0.1),
                              color: post.color,
                              fontWeight: 600,
                            }}
                          />
                        </Stack>
                        <Typography variant="h5" fontWeight="bold">
                          {post.title}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" sx={{ lineHeight: 1.7 }}>
                          {post.excerpt}
                        </Typography>
                        <Typography variant="caption" color="text.disabled">
                          {post.date}
                        </Typography>
                      </Stack>
                    </CardContent>
                  </CardActionArea>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>
    </Box>
  );
};

export default BlogIndex;
