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
  useTheme,
  Divider,
  Chip,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Avatar,
  LinearProgress,
} from '@mui/material';
import {
  Home as HomeIcon,
  Work as WorkIcon,
  Psychology as PsychologyIcon,
  AccountTree as TreeIcon,
  Verified as VerifiedIcon,
  Badge as BadgeIcon,
  Timeline as TimelineIcon,
  Groups as GroupsIcon,
  Business as BusinessIcon,
  Person as PersonIcon,
  ConnectWithoutContact as ConnectIcon,
  Visibility as VisibilityIcon,
  Balance as BalanceIcon,
  Security as SecurityIcon,
  AutoFixHigh as AutoFixIcon,
  Lightbulb as LightbulbIcon,
  School as SchoolIcon,
  EmojiEvents as AchievementIcon,
  Flag as FlagIcon,
  Route as RouteIcon,
  TrendingUp as GrowthIcon,
  CompareArrows as CompareIcon,
  DataObject as DataIcon,
  AccountBalance as TrustIcon,
  Assignment as DocumentIcon,
  Download as DownloadIcon,
  Upload as UploadIcon,
  Search as SearchIcon,
  FilterAlt as FilterIcon,
  SwapHoriz as SwapIcon,
  Hub as HubIcon,
} from '@mui/icons-material';
import { Link } from 'react-router-dom';
import ContactReveal from '../components/ContactReveal';

// ─── Role Profile Example ──────────────────────────────────
const roleExample = {
  title: 'Senior Frontend Developer',
  skills: [
    { name: 'React/TypeScript', level: 4, required: true },
    { name: 'State Management', level: 3, required: true },
    { name: 'Component Library Design', level: 3, required: false },
    { name: 'Team Collaboration', level: 3, required: true },
    { name: 'Mentoring', level: 2, required: false },
    { name: 'Performance Optimisation', level: 3, required: true },
  ],
};

// ─── Applicant Match Example ───────────────────────────────
const matchExample = {
  name: 'Alex Chen',
  matchScore: 87,
  strengths: ['React/TypeScript', 'Team Collaboration'],
  gaps: ['Mentoring', 'Performance Optimisation'],
};

// ─── Key Capabilities ──────────────────────────────────────
const capabilities = [
  {
    title: 'Role Architect',
    description: 'Define the precise skill profile for any position — not vague "5 years experience" but verified competency levels mapped to the CHISG framework.',
    icon: <BusinessIcon />,
    color: '#2563eb',
    features: [
      'Skill-based job descriptions, not keyword lists',
      'Weighted importance scoring per skill',
      'Benchmarking across similar roles',
      'Salary benchmarking by verified capability'
    ]
  },
  {
    title: 'Skills Passport',
    description: 'Every individual holds their own verified skills map — a portable record of demonstrated competency that travels with them from education to employment and between roles.',
    icon: <BadgeIcon />,
    color: '#7c3aed',
    features: [
      'CHISG-aligned skill definitions',
      'Verifiable evidence attachments',
      'Training and qualification integration',
      'Portable across employers'
    ]
  },
  {
    title: 'Intelligent Matching',
    description: 'Stop guessing. CareerOS calculates the semantic distance between a candidate\'s verified skills and a role\'s requirements — revealing not just who qualifies, but who\'s ready to grow.',
    icon: <CompareIcon />,
    color: '#10b981',
    features: [
      'Semantic gap analysis, not keyword matching',
      'Potential prediction — distance to competency',
      'Blind matching reduces unconscious bias',
      'Explanatory matches, not black boxes'
    ]
  },
  {
    title: 'Evidence Engine',
    description: 'Every claimed skill requires evidence. Work samples, peer endorsements, certifications, or observed competency — all stored with provenance, all auditable.',
    icon: <VerifiedIcon />,
    color: '#dc2626',
    features: [
      'Multi-format evidence attachment',
      'Endorsement chains',
      'Expiry and renewal tracking',
      'Fraud-resistant verification'
    ]
  },
];

// ─── Benefits Grid ─────────────────────────────────────────
const benefits = [
  {
    for: 'Employers',
    icon: <BusinessIcon />,
    points: [
      'Reduce time-to-hire by 40% with precision matching',
      'Eliminate "keyword bingo" applications',
      'Identify internal mobility candidates you\'d otherwise miss',
      'Auditable, defensible hiring decisions',
      'Skills gap analysis across your workforce'
    ]
  },
  {
    for: 'Individuals',
    icon: <PersonIcon />,
    points: [
      'Your CV becomes obsolete — your verified skills speak',
      'No more tailoring applications for every role',
      'See exactly where you match and where you need development',
      'Portable record that survives employer changes',
      'Evidence of capability, not years served'
    ]
  },
  {
    for: 'Educators & Trainers',
    icon: <SchoolIcon />,
    points: [
      'Map courses directly to CHISG skill definitions',
      'Graduates leave with verified competency profiles',
      'Close the feedback loop from employer outcomes',
      'Demonstrate ROI of training provision'
    ]
  },
];

// ─── Component ──────────────────────────────────────────────
const CareerOSLanding: React.FC = () => {
  const theme = useTheme();
  const brandColor = '#2563eb';

  return (
    <Box>
      {/* ── Hero Section ── */}
      <Box
        sx={{
          background: `linear-gradient(135deg, #2563eb 0%, #4f46e5 60%, #7c3aed 100%)`,
          color: 'white',
          pt: { xs: 6, md: 10 },
          pb: { xs: 8, md: 12 },
          px: 2,
          position: 'relative',
          overflow: 'hidden',
          '&::after': {
            content: '""',
            position: 'absolute',
            bottom: 0,
            left: 0,
            right: 0,
            height: '80px',
            background: 'linear-gradient(to top right, #f5f7f9 50%, transparent 50%)',
          },
        }}
      >
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Grid container spacing={4} alignItems="center">
            <Grid item xs={12} md={7}>
              <Stack spacing={3}>
                {/* Back to portfolio */}
                <Button
                  component="a"
                  href="/"
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
                  alt="ESP Thinking Portfolio"
                  sx={{
                    maxWidth: '100%',
                    width: { xs: '200px', md: '280px' },
                    height: 'auto',
                    alignSelf: 'flex-start',
                    mb: 1,
                  }}
                />

                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                  <BadgeIcon sx={{ fontSize: { xs: 48, md: 64 } }} />
                  <Box>
                    <Typography variant="h2" component="h1" fontWeight="bold" sx={{ lineHeight: 1.1 }}>
                      CareerOS
                    </Typography>
                    <Typography variant="subtitle1" sx={{ opacity: 0.85, mt: 0.5 }}>
                      Your Skills, Verified. Your Future, Mapped.
                    </Typography>
                  </Box>
                </Box>

                <Typography variant="h5" sx={{ opacity: 0.92, maxWidth: 620 }}>
                  The CV is obsolete. CareerOS replaces keyword-stuffed documents with 
                  verified skills maps — creating a common language between talent, employers, 
                  and educators. Every capability traced to evidence. Every match explained.
                </Typography>

                <Stack direction="row" spacing={1.5} flexWrap="wrap" useFlexGap>
                  <Chip label="CHISG-Aligned" sx={{ bgcolor: 'rgba(255,255,255,0.15)', color: 'white' }} />
                  <Chip label="Skills Passport" sx={{ bgcolor: 'rgba(255,255,255,0.15)', color: 'white' }} />
                  <Chip label="Semantic Matching" sx={{ bgcolor: 'rgba(255,255,255,0.15)', color: 'white' }} />
                  <Chip label="Evidence-Backed" sx={{ bgcolor: 'rgba(255,255,255,0.15)', color: 'white' }} />
                </Stack>

                <Stack direction="row" spacing={2} sx={{ pt: 1 }}>
                  <Button
                    variant="contained"
                    size="large"
                    component={Link}
                    to="/careeros/demo"
                    sx={{
                      bgcolor: 'white',
                      color: '#2563eb',
                      fontWeight: 'bold',
                      px: 4,
                      '&:hover': { bgcolor: alpha('#ffffff', 0.9) },
                    }}
                  >
                    Request Enterprise Demo
                  </Button>
                 
                </Stack>
              </Stack>
            </Grid>

            <Grid item xs={12} md={5}>
              <Paper
                elevation={8}
                sx={{
                  p: 3,
                  borderRadius: 3,
                  background: alpha(theme.palette.background.paper, 0.95),
                }}
              >
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 2 }}>
                  <CompareIcon sx={{ color: '#2563eb' }} />
                  <Typography variant="h6" color="text.primary" fontWeight="bold">
                    Match, Don't Guess
                  </Typography>
                </Box>
                
                {/* Role Profile Example */}
                <Typography variant="subtitle2" gutterBottom>
                  Role: {roleExample.title}
                </Typography>
                <Box sx={{ mb: 2 }}>
                  {roleExample.skills.map((skill) => (
                    <Box key={skill.name} sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}>
                      <Typography variant="caption" sx={{ width: 140 }}>{skill.name}</Typography>
                      <Box sx={{ flexGrow: 1, display: 'flex', gap: 0.5 }}>
                        {[1,2,3,4].map((level) => (
                          <Box
                            key={level}
                            sx={{
                              height: 8,
                              width: '25%',
                              borderRadius: 1,
                              bgcolor: level <= skill.level ? '#2563eb' : alpha('#94a3b8', 0.3),
                            }}
                          />
                        ))}
                      </Box>
                      <Typography variant="caption" color="text.secondary">
                        {skill.required ? 'Required' : 'Preferred'}
                      </Typography>
                    </Box>
                  ))}
                </Box>

                <Divider sx={{ my: 2 }} />

                {/* Applicant Match Example */}
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1.5 }}>
                  <Typography variant="subtitle2">Candidate: {matchExample.name}</Typography>
                  <Chip 
                    label={`${matchExample.matchScore}% Match`}
                    size="small"
                    sx={{ bgcolor: alpha('#10b981', 0.1), color: '#10b981', fontWeight: 'bold' }}
                  />
                </Box>

                <Grid container spacing={1} sx={{ mb: 1 }}>
                  <Grid item xs={6}>
                    <Typography variant="caption" color="success.main" fontWeight="bold">✓ Strengths</Typography>
                    {matchExample.strengths.map((s) => (
                      <Typography key={s} variant="caption" display="block" sx={{ color: '#10b981' }}>
                        {s}
                      </Typography>
                    ))}
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="caption" color="warning.main" fontWeight="bold">○ Development areas</Typography>
                    {matchExample.gaps.map((g) => (
                      <Typography key={g} variant="caption" display="block" sx={{ color: '#f59e0b' }}>
                        {g}
                      </Typography>
                    ))}
                  </Grid>
                </Grid>

                <Typography variant="body2" sx={{ mt: 2, fontStyle: 'italic', color: 'text.secondary' }}>
                  "Not just who qualifies today — but who could qualify with targeted development."
                </Typography>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* ── The Problem: Why the CV Fails ── */}
      <Box sx={{ py: { xs: 6, md: 8 }, bgcolor: '#f5f7f9' }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Grid container spacing={4} alignItems="center">
            <Grid item xs={12} md={6}>
              <Typography variant="overline" sx={{ color: brandColor, fontWeight: 'bold' }}>
                The Document That Lied
              </Typography>
              <Typography variant="h4" component="h2" fontWeight="bold" gutterBottom>
                The CV Was Never Designed for This
              </Typography>
              <Typography variant="body1" paragraph sx={{ fontSize: '1.1rem', lineHeight: 1.7 }}>
                The curriculum vitae — "course of life" — was invented in 1482. It describes what 
                someone has done, not what they can do. Five centuries later, we're still judging 
                capability by counting years and scanning keywords.
              </Typography>
              <Typography variant="body1" paragraph sx={{ fontSize: '1.1rem', lineHeight: 1.7 }}>
                <strong>The result:</strong> Applicants optimise for applicant tracking systems, not 
                actual competence. Hiring managers drown in irrelevant applications. Unconventional 
                talent is filtered out before a human ever sees them. And every organisation 
                competes for the same 20% of candidates while overlooking capable people who don't 
                know how to game the system.
              </Typography>
              <Typography variant="body1" sx={{ fontSize: '1.1rem', lineHeight: 1.7 }}>
                <strong>CareerOS doesn't digitise the CV — it replaces it.</strong>
              </Typography>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Card elevation={4} sx={{ borderRadius: 3 }}>
                <CardContent sx={{ p: 4 }}>
                  <Typography variant="h6" gutterBottom fontWeight="bold" sx={{ color: brandColor }}>
                    The CV vs. The Skills Passport
                  </Typography>
                  
                  <Grid container spacing={2}>
                    <Grid item xs={12} sm={6}>
                      <Paper elevation={0} sx={{ p: 2, bgcolor: alpha('#ef4444', 0.04), height: '100%' }}>
                        <Typography variant="subtitle2" gutterBottom sx={{ color: '#ef4444' }}>
                          ❌ CV / Resume
                        </Typography>
                        <List dense>
                          <ListItem sx={{ px: 0 }}>
                            <ListItemIcon sx={{ minWidth: 30 }}><Typography color="#ef4444">✗</Typography></ListItemIcon>
                            <ListItemText primary="Self-declared, unverified" />
                          </ListItem>
                          <ListItem sx={{ px: 0 }}>
                            <ListItemIcon sx={{ minWidth: 30 }}><Typography color="#ef4444">✗</Typography></ListItemIcon>
                            <ListItemText primary="Optimised for robots, not humans" />
                          </ListItem>
                          <ListItem sx={{ px: 0 }}>
                            <ListItemIcon sx={{ minWidth: 30 }}><Typography color="#ef4444">✗</Typography></ListItemIcon>
                            <ListItemText primary="Starts over with every application" />
                          </ListItem>
                          <ListItem sx={{ px: 0 }}>
                            <ListItemIcon sx={{ minWidth: 30 }}><Typography color="#ef4444">✗</Typography></ListItemIcon>
                            <ListItemText primary="Years = competence (they aren't)" />
                          </ListItem>
                          <ListItem sx={{ px: 0 }}>
                            <ListItemIcon sx={{ minWidth: 30 }}><Typography color="#ef4444">✗</Typography></ListItemIcon>
                            <ListItemText primary="Gaps are shameful" />
                          </ListItem>
                        </List>
                      </Paper>
                    </Grid>
                    
                    <Grid item xs={12} sm={6}>
                      <Paper elevation={0} sx={{ p: 2, bgcolor: alpha('#10b981', 0.04), height: '100%' }}>
                        <Typography variant="subtitle2" gutterBottom sx={{ color: '#10b981' }}>
                          ✅ Skills Passport
                        </Typography>
                        <List dense>
                          <ListItem sx={{ px: 0 }}>
                            <ListItemIcon sx={{ minWidth: 30 }}><Typography color="#10b981">✓</Typography></ListItemIcon>
                            <ListItemText primary="Verified with evidence" />
                          </ListItem>
                          <ListItem sx={{ px: 0 }}>
                            <ListItemIcon sx={{ minWidth: 30 }}><Typography color="#10b981">✓</Typography></ListItemIcon>
                            <ListItemText primary="Read by semantic matching" />
                          </ListItem>
                          <ListItem sx={{ px: 0 }}>
                            <ListItemIcon sx={{ minWidth: 30 }}><Typography color="#10b981">✓</Typography></ListItemIcon>
                            <ListItemText primary="Portable, lifetime record" />
                          </ListItem>
                          <ListItem sx={{ px: 0 }}>
                            <ListItemIcon sx={{ minWidth: 30 }}><Typography color="#10b981">✓</Typography></ListItemIcon>
                            <ListItemText primary="Demonstrated competency" />
                          </ListItem>
                          <ListItem sx={{ px: 0 }}>
                            <ListItemIcon sx={{ minWidth: 30 }}><Typography color="#10b981">✓</Typography></ListItemIcon>
                            <ListItemText primary="Gaps = development opportunities" />
                          </ListItem>
                        </List>
                      </Paper>
                    </Grid>
                  </Grid>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* ── The Skills Passport ── */}
      <Box sx={{ py: { xs: 6, md: 8 } }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Box sx={{ textAlign: 'center', mb: 5 }}>
            <BadgeIcon sx={{ fontSize: 56, color: brandColor, mb: 2 }} />
            <Typography variant="h4" component="h2" fontWeight="bold" gutterBottom>
              The Skills Passport
            </Typography>
            <Typography variant="h6" color="text.secondary" sx={{ fontWeight: 'normal', maxWidth: 700, mx: 'auto' }}>
              One verified record of your capability — from your first qualification to your executive role. 
              You own it. It travels with you. It never needs rewriting.
            </Typography>
          </Box>

          <Grid container spacing={3} sx={{ mb: 5 }}>
            <Grid item xs={12} md={4}>
              <Paper elevation={2} sx={{ p: 3, height: '100%', borderTop: `4px solid #2563eb` }}>
                <Typography variant="h6" fontWeight="bold" gutterBottom>
                  🎓 Education
                </Typography>
                <Typography variant="body2" color="text.secondary" paragraph>
                  Qualifications, courses, and training mapped directly to CHISG skill definitions. 
                  Not just "passed" — demonstrably competent in specific capabilities.
                </Typography>
                <Box sx={{ bgcolor: alpha('#2563eb', 0.04), p: 1.5, borderRadius: 1 }}>
                  <Typography variant="caption" display="block" fontWeight="bold">
                    Example:
                  </Typography>
                  <Typography variant="caption" display="block">
                    "A-Level Mathematics" → Algebraic Reasoning L4, Statistical Analysis L3
                  </Typography>
                </Box>
              </Paper>
            </Grid>
            
            <Grid item xs={12} md={4}>
              <Paper elevation={2} sx={{ p: 3, height: '100%', borderTop: `4px solid #7c3aed` }}>
                <Typography variant="h6" fontWeight="bold" gutterBottom>
                  💼 Employment
                </Typography>
                <Typography variant="body2" color="text.secondary" paragraph>
                  Demonstrated competency through real work. Peer endorsements, project artefacts, 
                  and manager verification — all attached as evidence to specific skills.
                </Typography>
                <Box sx={{ bgcolor: alpha('#7c3aed', 0.04), p: 1.5, borderRadius: 1 }}>
                  <Typography variant="caption" display="block" fontWeight="bold">
                    Example:
                  </Typography>
                  <Typography variant="caption" display="block">
                    "Led migration to React" → React L4, Team Leadership L3, Project Planning L3
                  </Typography>
                </Box>
              </Paper>
            </Grid>
            
            <Grid item xs={12} md={4}>
              <Paper elevation={2} sx={{ p: 3, height: '100%', borderTop: `4px solid #10b981` }}>
                <Typography variant="h6" fontWeight="bold" gutterBottom>
                  🔄 Lifelong Learning
                </Typography>
                <Typography variant="body2" color="text.secondary" paragraph>
                  Certifications, micro-credentials, self-study, and non-formal learning. 
                  Every skill acquisition captured, verified, and added to your permanent record.
                </Typography>
                <Box sx={{ bgcolor: alpha('#10b981', 0.04), p: 1.5, borderRadius: 1 }}>
                  <Typography variant="caption" display="block" fontWeight="bold">
                    Example:
                  </Typography>
                  <Typography variant="caption" display="block">
                    "AWS Certified Cloud Practitioner" → Cloud Architecture L2, Security Fundamentals L2
                  </Typography>
                </Box>
              </Paper>
            </Grid>
          </Grid>

          <Paper
            elevation={3}
            sx={{
              p: 4,
              borderRadius: 3,
              bgcolor: alpha(brandColor, 0.04),
              border: `1px solid ${alpha(brandColor, 0.2)}`,
            }}
          >
            <Typography variant="h6" gutterBottom fontWeight="bold" sx={{ color: brandColor }}>
              You don't update your passport at every border — and you shouldn't have to update your skills record for every application.
            </Typography>
            <Typography variant="body1">
              The Skills Passport is a living document. Every new competency, every verified skill, every piece of evidence 
              accumulates over your entire career. When you apply for a role, you grant access — you don't rebuild from scratch.
            </Typography>
          </Paper>
        </Container>
      </Box>

      {/* ── Core Capabilities ── */}
      <Box sx={{ py: { xs: 6, md: 8 }, bgcolor: '#f5f7f9' }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            Four Capabilities, One Platform
          </Typography>
          <Typography
            variant="body1"
            color="text.secondary"
            textAlign="center"
            sx={{ mb: 5, maxWidth: 780, mx: 'auto' }}
          >
            CareerOS connects the entire skills ecosystem — employers, individuals, and educators — 
            through a shared language of verified competency.
          </Typography>

          <Grid container spacing={4}>
            {capabilities.map((capability) => (
              <Grid item xs={12} md={6} key={capability.title}>
                <Card
                  elevation={2}
                  sx={{
                    height: '100%',
                    borderTop: `4px solid ${capability.color}`,
                    borderRadius: 2,
                  }}
                >
                  <CardContent sx={{ p: 3 }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                      <Avatar sx={{ bgcolor: alpha(capability.color, 0.1), color: capability.color }}>
                        {capability.icon}
                      </Avatar>
                      <Typography variant="h6" fontWeight="bold">
                        {capability.title}
                      </Typography>
                    </Box>
                    
                    <Typography variant="body2" color="text.secondary" paragraph>
                      {capability.description}
                    </Typography>
                    
                    <List dense>
                      {capability.features.map((feature) => (
                        <ListItem key={feature} sx={{ px: 0, py: 0.25 }}>
                          <ListItemIcon sx={{ minWidth: 30 }}>
                            <VerifiedIcon sx={{ fontSize: 16, color: capability.color }} />
                          </ListItemIcon>
                          <ListItemText 
                            primary={feature} 
                            primaryTypographyProps={{ variant: 'body2' }}
                          />
                        </ListItem>
                      ))}
                    </List>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* ── Benefits by Audience ── */}
      <Box sx={{ py: { xs: 6, md: 8 } }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            Designed for the Whole Ecosystem
          </Typography>
          <Typography
            variant="body1"
            color="text.secondary"
            textAlign="center"
            sx={{ mb: 5, maxWidth: 700, mx: 'auto' }}
          >
            CareerOS isn't just an HR tool — it's a shared infrastructure for talent discovery, 
            development, and deployment.
          </Typography>

          <Grid container spacing={4}>
            {benefits.map((benefit) => (
              <Grid item xs={12} md={4} key={benefit.for}>
                <Card
                  elevation={2}
                  sx={{
                    height: '100%',
                    borderRadius: 2,
                  }}
                >
                  <CardContent sx={{ p: 3 }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 2 }}>
                      <Avatar sx={{ bgcolor: alpha(brandColor, 0.1), color: brandColor }}>
                        {benefit.icon}
                      </Avatar>
                      <Typography variant="h6" fontWeight="bold">
                        For {benefit.for}
                      </Typography>
                    </Box>
                    
                    <List dense>
                      {benefit.points.map((point) => (
                        <ListItem key={point} sx={{ px: 0, py: 0.5 }}>
                          <ListItemIcon sx={{ minWidth: 30 }}>
                            <VerifiedIcon sx={{ fontSize: 16, color: '#10b981' }} />
                          </ListItemIcon>
                          <ListItemText 
                            primary={point} 
                            primaryTypographyProps={{ variant: 'body2' }}
                          />
                        </ListItem>
                      ))}
                    </List>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* ── Matching in Action ── */}
      <Box sx={{ py: { xs: 6, md: 8 }, bgcolor: alpha(brandColor, 0.04) }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Grid container spacing={4} alignItems="center">
            <Grid item xs={12} md={6}>
              <Typography variant="overline" sx={{ color: brandColor, fontWeight: 'bold' }}>
                Semantic Matching
              </Typography>
              <Typography variant="h4" component="h2" fontWeight="bold" gutterBottom>
                Beyond Keyword Counting
              </Typography>
              <Typography variant="body1" paragraph sx={{ fontSize: '1.1rem', lineHeight: 1.7 }}>
                Traditional ATS checks if the word "React" appears in your CV. CareerOS understands 
                that "Vue" + "component architecture" + "state management" indicates transferable 
                capability — even if you've never written a line of React.
              </Typography>
              <Typography variant="body1" paragraph sx={{ fontSize: '1.1rem', lineHeight: 1.7 }}>
                <strong>Why this matters:</strong> The best candidate for a React role might be a 
                Vue developer with strong fundamentals. The best manager might come from a different 
                industry entirely. CareerOS finds the signal in the noise.
              </Typography>
              
              <Paper elevation={1} sx={{ p: 2, bgcolor: 'white', mt: 2 }}>
                <Typography variant="subtitle2" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <HubIcon sx={{ color: brandColor }} /> How it works:
                </Typography>
                <Typography variant="body2">
                  CHISG's semantic graph understands that <strong>React</strong> and <strong>Vue</strong> 
                  share 78% of prerequisite skills. A candidate with Vue L4 is <strong>semantically close</strong> 
                  to a React L3 requirement — and CareerOS explains exactly why.
                </Typography>
              </Paper>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Card elevation={4} sx={{ borderRadius: 3 }}>
                <CardContent sx={{ p: 4 }}>
                  <Typography variant="h6" gutterBottom fontWeight="bold" sx={{ color: brandColor }}>
                    Distance to Competency
                  </Typography>
                  
                  <Box sx={{ mb: 3 }}>
                    <Typography variant="body2" paragraph>
                      Instead of "doesn't meet requirements," CareerOS shows <strong>how far</strong> 
                      a candidate is from competency — and what development would close the gap.
                    </Typography>
                  </Box>
                  
                  <Box sx={{ mb: 2 }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
                      <Typography variant="caption">Current capability</Typography>
                      <Typography variant="caption" fontWeight="bold">63%</Typography>
                    </Box>
                    <LinearProgress 
                      variant="determinate" 
                      value={63} 
                      sx={{ height: 10, borderRadius: 5, bgcolor: alpha('#94a3b8', 0.2) }}
                    />
                  </Box>
                  
                  <Box sx={{ mb: 3 }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
                      <Typography variant="caption">Required level</Typography>
                      <Typography variant="caption" fontWeight="bold">85%</Typography>
                    </Box>
                    <LinearProgress 
                      variant="determinate" 
                      value={85} 
                      sx={{ height: 10, borderRadius: 5, bgcolor: alpha('#94a3b8', 0.2) }}
                    />
                  </Box>
                  
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, p: 2, bgcolor: alpha('#10b981', 0.1), borderRadius: 2 }}>
                    <Typography variant="h4" sx={{ color: '#10b981', fontWeight: 'bold' }}>22%</Typography>
                    <Box>
                      <Typography variant="body2" fontWeight="bold">Distance to competency</Typography>
                      <Typography variant="caption" color="text.secondary">
                        Estimated 4-6 weeks of targeted development
                      </Typography>
                    </Box>
                  </Box>
                  
                  <Divider sx={{ my: 2 }} />
                  
                  <Typography variant="body2" sx={{ fontStyle: 'italic' }}>
                    "Not ready today — but would be with the right support. Compare this to rejecting 
                    an applicant who could succeed with minimal development."
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* ── Integration with Ecosystem ── */}
      <Box sx={{ py: { xs: 6, md: 8 } }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Paper
            elevation={3}
            sx={{
              p: 4,
              borderRadius: 4,
              background: `linear-gradient(135deg, ${alpha('#0f766e', 0.04)} 0%, ${alpha('#2563eb', 0.04)} 100%)`,
              border: `1px solid ${alpha('#2563eb', 0.2)}`,
            }}
          >
            <Grid container spacing={4} alignItems="center">
              <Grid item xs={12} md={8}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                  <TreeIcon sx={{ color: '#15803d', fontSize: 32 }} />
                  <Typography variant="overline" sx={{ color: '#15803d', fontWeight: 'bold' }}>
                    Built on CHISG
                  </Typography>
                </Box>
                <Typography variant="h5" fontWeight="bold" gutterBottom>
                  The Same Language, From Classroom to Boardroom
                </Typography>
                <Typography variant="body1" paragraph>
                  CareerOS uses the same CHISG skill architecture as PrimaryOS and ParentOS. 
                  A skill a child develops in Year 5 is the same skill an employer searches for 
                  20 years later. No translation. No loss. No starting over.
                </Typography>
                <Typography variant="body1">
                  <strong>This is the closing of the loop.</strong> Education knows what employment 
                  needs. Employment validates what education taught. And individuals carry their 
                  verified capability with them — always accumulating, never restarting.
                </Typography>
              </Grid>
              <Grid item xs={12} md={4}>
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                    <Avatar sx={{ bgcolor: alpha('#15803d', 0.1), color: '#15803d' }}>
                      <SchoolIcon />
                    </Avatar>
                    <Box>
                      <Typography variant="subtitle2">PrimaryOS</Typography>
                      <Typography variant="caption" color="text.secondary">Skill development</Typography>
                    </Box>
                  </Box>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                    <Avatar sx={{ bgcolor: alpha('#0f766e', 0.1), color: '#0f766e' }}>
                      <BadgeIcon />
                    </Avatar>
                    <Box>
                      <Typography variant="subtitle2">ParentOS</Typography>
                      <Typography variant="caption" color="text.secondary">Skill observation</Typography>
                    </Box>
                  </Box>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                    <Avatar sx={{ bgcolor: alpha('#2563eb', 0.1), color: '#2563eb' }}>
                      <WorkIcon />
                    </Avatar>
                    <Box>
                      <Typography variant="subtitle2">CareerOS</Typography>
                      <Typography variant="caption" color="text.secondary">Skill deployment</Typography>
                    </Box>
                  </Box>
                </Box>
              </Grid>
            </Grid>
          </Paper>
        </Container>
      </Box>

      {/* ── What Makes CareerOS Different ── */}
      <Box sx={{ py: { xs: 6, md: 8 }, bgcolor: '#f5f7f9' }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            What Makes CareerOS Different
          </Typography>

          <Grid container spacing={2} sx={{ mt: 2 }}>
            {[
              'Replaces the CV with verified, portable skills records',
              'Semantic matching, not keyword scanning',
              'Distance to competency — not just pass/fail',
              'Explanatory matches, not black boxes',
              'CHISG-aligned — one skill language from education to employment',
              'Evidence-backed verification, not self-declaration',
              'Reduces unconscious bias through blind matching',
              'Reveals internal mobility candidates hidden in your workforce',
              'Individuals own their passport — it travels with them',
              'Gaps become development pathways, not rejection reasons',
              'Education providers close the feedback loop',
              'Built on 12+ years of developmental research',
            ].map((point) => (
              <Grid item xs={12} sm={6} md={4} key={point}>
                <Box sx={{ display: 'flex', gap: 1.5, p: 1.5 }}>
                  <Box sx={{ color: brandColor, mt: 0.25 }}>✓</Box>
                  <Typography variant="body2">{point}</Typography>
                </Box>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* ── CTA Section ── */}
      <Box
        sx={{
          py: { xs: 6, md: 8 },
          textAlign: 'center',
          background: `linear-gradient(135deg, #2563eb 0%, #4f46e5 100%)`,
          color: 'white',
        }}
      >
        <Container maxWidth="md">
          <BadgeIcon sx={{ fontSize: 56, mb: 2 }} />
          <Typography variant="h4" fontWeight="bold" gutterBottom>
            Replace the CV. Unlock Hidden Talent.
          </Typography>
          <Typography variant="body1" sx={{ mb: 4, maxWidth: 600, mx: 'auto', opacity: 0.95 }}>
            CareerOS is available for enterprise deployment, skills passport pilots, 
            and education-partner integration.
          </Typography>
          
          <Stack direction="row" spacing={2} justifyContent="center" flexWrap="wrap" useFlexGap>
            <Button
              variant="contained"
              size="large"
              component={Link}
              to="/careeros/enterprise"
              sx={{
                px: 5,
                py: 1.5,
                fontWeight: 'bold',
                bgcolor: 'white',
                color: '#2563eb',
                '&:hover': { bgcolor: alpha('#ffffff', 0.9) },
              }}
            >
              Book Enterprise Demo
            </Button>
            <ContactReveal email="world@espthinking.co.uk" label="Get in Touch" />
            <Button
              variant="outlined"
              size="large"
              component={Link}
              to="/careeros/passport"
              sx={{
                px: 5,
                py: 1.5,
                fontWeight: 'bold',
                color: 'white',
                borderColor: 'white',
                '&:hover': { borderColor: 'white', bgcolor: alpha('#ffffff', 0.1) },
              }}
            >
              Create Your Skills Passport
            </Button>
          </Stack>
          
          <Typography variant="caption" sx={{ display: 'block', mt: 4, opacity: 0.8 }}>
            For individuals: free foundational passport. For organisations: enterprise licensing available.
          </Typography>
        </Container>
      </Box>

      {/* ── Footer ── */}
      <Box sx={{ py: 3, textAlign: 'center', bgcolor: '#2563eb', color: 'rgba(255,255,255,0.9)' }}>
        <Typography variant="body2">
          CareerOS — Your Skills, Verified. Your Future, Mapped. — Part of the ESP Thinking Portfolio
        </Typography>
        <Typography variant="caption" sx={{ display: 'block', mt: 0.5, opacity: 0.8 }}>
          Powered by CHISG | Connected to PrimaryOS & ParentOS | The Verified Alternative to the CV
        </Typography>
      </Box>
    </Box>
  );
};

export default CareerOSLanding;