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
  School as SchoolIcon,
  Psychology as PsychologyIcon,
  MenuBook as BookIcon,
  EmojiPeople as PeopleIcon,
  Forest as ForestIcon,
  Explore as ExploreIcon,
  AutoFixHigh as AutoFixIcon,
  Timeline as TimelineIcon,
  Groups as GroupsIcon,
  ChildCare as ChildIcon,
  ConnectWithoutContact as ConnectIcon,
  Visibility as VisibilityIcon,
  Balance as BalanceIcon,
  Favorite as FavoriteIcon,
  Security as SecurityIcon,
  SelfImprovement as SelfIcon,
  Speed as SpeedIcon,
  SwapHoriz as SwapIcon,
  Sensors as SensorsIcon,
  Straighten as StraightenIcon,
  Lightbulb as LightbulbIcon,
  Verified as VerifiedIcon,
  Warning as WarningIcon,
  LockOpen as LockOpenIcon,
  AccountTree as TreeIcon,
  Assignment as JournalIcon,
  EmojiEvents as AchievementIcon,
  Flag as FlagIcon,
  Route as RouteIcon,
  TrendingUp as GrowthIcon,
  FamilyRestroom as ParentIcon,
} from '@mui/icons-material';
import { Link } from 'react-router-dom';
import ContactReveal from '../components/ContactReveal';
import SEO from '../components/SEO';

// ─── ETP Spectra for PrimaryOS ─────────────────────────────
const primarySpectra = [
  {
    name: 'Social Gravity',
    left: 'Independent',
    right: 'Cohesive',
    description: 'How you recharge — alone time or people time? Neither is wrong, but different situations need different settings.',
    skill: 'Social Switching',
    skillDescription: 'The ability to work independently when needed and collaborate when useful, rather than being stuck in one mode.',
    color: '#1565c0',
    classroom: 'Knows when to focus alone vs. when to ask a partner',
  },
  {
    name: 'Energy Directionality',
    left: 'Inward',
    right: 'Outward',
    description: 'Do you think best in silence, or by talking it through?',
    skill: 'Processing Flexibility',
    skillDescription: 'The ability to pause and process internally OR think out loud, depending on what the situation requires.',
    color: '#6a1b9a',
    classroom: 'Can work through problems quietly OR explain thinking aloud',
  },
  {
    name: 'Voltage Sensitivity',
    left: 'Insulated',
    right: 'Conductive',
    description: 'How much do you notice and absorb other people\'s feelings?',
    skill: 'Empathy Regulation',
    skillDescription: 'The ability to tune in to others when helpful, and maintain your own boundaries when needed.',
    color: '#f57c00',
    classroom: 'Notices when peers need help, but doesn\'t get overwhelmed by classroom energy',
  },
  {
    name: 'Threat Response',
    left: 'Passive',
    right: 'Aggressive',
    description: 'How do you react when you feel scared or cornered?',
    skill: 'Response Selection',
    skillDescription: 'The ability to choose your response rather than reacting automatically — stepping back OR standing ground with intention.',
    color: '#c62828',
    classroom: 'Can pause before reacting, and choose a helpful response',
  },
  {
    name: 'Care Response',
    left: 'Detached',
    right: 'Nurturing',
    description: 'How do you respond when someone is upset?',
    skill: 'Comfort Calibration',
    skillDescription: 'Knowing when to offer support and when to give space — both are valid ways to care.',
    color: '#ad1457',
    classroom: 'Offers help appropriately, respects peers\' need for space',
  },
  {
    name: 'Risk Tolerance',
    left: 'Averse',
    right: 'Seeking',
    description: 'Are you cautious, or do you love a challenge?',
    skill: 'Risk Assessment',
    skillDescription: 'The ability to spot danger AND spot opportunity — and choose which risks are worth taking.',
    color: '#2e7d32',
    classroom: 'Tries challenging work, knows when to ask for help',
  },
  {
    name: 'Integrity Logic',
    left: 'Relativistic',
    right: 'Absolutist',
    description: 'Are rules flexible guidelines or unbreakable truths?',
    skill: 'Rule Fluidity',
    skillDescription: 'Understanding that some rules keep us safe (and must be followed) while others are situational (and can be adapted).',
    color: '#4527a0',
    classroom: 'Follows core expectations, adapts to different classroom contexts',
  },
  {
    name: 'Mirror Neuron Tuning',
    left: 'Selective',
    right: 'Absorbent',
    description: 'Do you stay steady in group energy, or do you match the room?',
    skill: 'Energy Management',
    skillDescription: 'The ability to regulate how much of the group\'s energy you absorb — connecting without losing yourself.',
    color: '#00838f',
    classroom: 'Contributes to group energy without being overwhelmed by it',
  },
  {
    name: 'Orderliness',
    left: 'Flexible',
    right: 'Ordered',
    description: 'Do you prefer structure or spontaneity?',
    skill: 'Structure Adaptability',
    skillDescription: 'The ability to work within routines AND handle unexpected changes without distress.',
    color: '#5d4037',
    classroom: 'Follows routines, adapts when plans change',
  },
];

// ─── Foundational Skill Domains ────────────────────────────
const foundationalDomains = [
  {
    name: 'Focus & Attention',
    skills: ['Sustained focus', 'Selective attention', 'Task initiation', 'Distraction management'],
    icon: <VisibilityIcon />,
    color: '#1565c0',
  },
  {
    name: 'Self-Regulation',
    skills: ['Emotion identification', 'Impulse control', 'Frustration tolerance', 'Recovery after setback'],
    icon: <BalanceIcon />,
    color: '#6a1b9a',
  },
  {
    name: 'Social Cognition',
    skills: ['Perspective taking', 'Reading social cues', 'Conflict resolution', 'Group contribution'],
    icon: <PeopleIcon />,
    color: '#2e7d32',
  },
  {
    name: 'Executive Function',
    skills: ['Planning', 'Organising', 'Prioritising', 'Self-monitoring', 'Flexible thinking'],
    icon: <PsychologyIcon />,
    color: '#c62828',
  },
  {
    name: 'Motor Planning',
    skills: ['Fine motor control', 'Gross motor coordination', 'Spatial awareness', 'Tool use'],
    icon: <StraightenIcon />,
    color: '#f57c00',
  },
  {
    name: 'Metacognition',
    skills: ['Knowing what you know', 'Choosing strategies', 'Reflecting on learning', 'Setting goals'],
    icon: <SelfIcon />,
    color: '#00838f',
  },
];

// ─── Sticker Book Journey Stages ───────────────────────────
const journeyStages = [
  {
    stage: 'Stage 1',
    title: 'Discovering My Settings',
    age: 'Reception – Year 1',
    description: 'Students learn that everyone has different "settings" and that no setting is wrong. They begin to notice their own preferences.',
    stickers: ['I noticed my feeling', 'We are all different', 'My favourite way to learn'],
    teacherRole: 'Introduce vocabulary. Name what you observe without judgment.',
    studentRole: 'Notice and wonder. "I noticed I like..."',
    color: '#2e7d32',
  },
  {
    stage: 'Stage 2',
    title: 'Naming My Settings',
    age: 'Year 2 – Year 3',
    description: 'Students can identify their own spectrum positions and recognise when they\'re in different states.',
    stickers: ['My Social setting today', 'I needed quiet', 'I took a brave risk'],
    teacherRole: 'Provide structured reflection opportunities. Validate all positions.',
    studentRole: 'Identify and name your settings. "Right now I need..."',
    color: '#1565c0',
  },
  {
    stage: 'Stage 3',
    title: 'Moving My Sliders',
    age: 'Year 4 – Year 5',
    description: 'Students learn that they can shift their position with effort and support. They develop strategies for sticky sliders.',
    stickers: ['I shifted my setting', 'I used a strategy', 'I helped someone else shift'],
    teacherRole: 'Teach strategies. Scaffold flexibility. Celebrate effort to shift.',
    studentRole: 'Try moving your slider. Notice what helps. "I can..."',
    color: '#6a1b9a',
  },
  {
    stage: 'Stage 4',
    title: 'Piloting Myself',
    age: 'Year 6 +',
    description: 'Students can independently choose appropriate settings for different contexts. They self-regulate with decreasing adult support.',
    stickers: ['I chose my setting', 'I matched the situation', 'I am the pilot'],
    teacherRole: 'Transfer control. Consult on challenges. Step back.',
    studentRole: 'Choose your setting deliberately. "This situation needs..."',
    color: '#00838f',
  },
];

// ─── Component ──────────────────────────────────────────────
const PrimaryOSLanding: React.FC = () => {
  const theme = useTheme();
  const brandColor = '#f59e0b';

  return (
    <Box>
      <SEO
        title="PrimaryOS — The Operating System for Growing Learners"
        description="PrimaryOS adapts the 9 ETP biological spectra for primary-age children. Neuron Navigators guide books, foundational skill domains, and teacher tools that replace shame with strategy."
        path="/primary-os"
      />
      {/* ── Hero Section ── */}
      <Box
        sx={{
          background: `linear-gradient(135deg, #f59e0b 0%, #fbbf24 50%, #fcd34d 100%)`,
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
                  src={`${process.env.PUBLIC_URL}/ESPLogoShort.png`}
                  alt="ESP Thinking Portfolio"
                  sx={{
                    maxWidth: '100%',
                    width: { xs: '200px', md: '280px' },
                    height: 'auto',
                    alignSelf: 'flex-start',
                  }}
                />

                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                  <SchoolIcon sx={{ fontSize: { xs: 48, md: 64 } }} />
                  <Box>
                    <Typography variant="h2" component="h1" fontWeight="bold" sx={{ lineHeight: 1.1 }}>
                      PrimaryOS
                    </Typography>
                    <Typography variant="subtitle1" sx={{ opacity: 0.85, mt: 0.5 }}>
                      The Operating System for Growing Learners
                    </Typography>
                  </Box>
                </Box>

                <Typography variant="h5" sx={{ opacity: 0.92, maxWidth: 620 }}>
                  A classroom framework that teaches children how their own minds work — 
                  and gives them the skills to become the pilot of their own learning.
                </Typography>

                <Stack direction="row" spacing={1.5} flexWrap="wrap" useFlexGap>
                  <Chip label="9 ETP Spectra" sx={{ bgcolor: 'rgba(255,255,255,0.15)', color: 'white' }} />
                  <Chip label="36+ Foundational Skills" sx={{ bgcolor: 'rgba(255,255,255,0.15)', color: 'white' }} />
                  <Chip label="Neuron Navigators Guide" sx={{ bgcolor: 'rgba(255,255,255,0.15)', color: 'white' }} />
                  <Chip label="CHISG-Aligned" sx={{ bgcolor: 'rgba(255,255,255,0.15)', color: 'white' }} />
                </Stack>

                <Stack direction="row" spacing={2} sx={{ pt: 1 }}>
                  {/* <Button
                    variant="contained"
                    size="large"
                    component={Link}
                    to="/primaryos/dashboard"
                    sx={{
                      bgcolor: 'white',
                      color: '#f59e0b',
                      fontWeight: 'bold',
                      px: 4,
                      '&:hover': { bgcolor: alpha('#ffffff', 0.9) },
                    }}
                  >
                    Launch Teacher Dashboard
                  </Button>
                  <Button
                    variant="outlined"
                    size="large"
                    component={Link}
                    to="/primaryos/neuron-navigators"
                    sx={{
                      color: 'white',
                      borderColor: 'white',
                      fontWeight: 'bold',
                      px: 4,
                      '&:hover': { borderColor: 'white', bgcolor: 'rgba(255,255,255,0.1)' },
                    }}
                  >
                    Explore Neuron Navigators
                  </Button> */}
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
                  <JournalIcon sx={{ color: '#f59e0b' }} />
                  <Typography variant="h6" color="text.primary" fontWeight="bold">
                    The Neuron Navigators Guide Book
                  </Typography>
                </Box>
                
                <Typography variant="body2" color="text.secondary" paragraph>
                  <strong>Your map to understanding how your brain works best.</strong> Every student 
                  receives their own interactive sticker book — a personal journal that transforms 
                  abstract neuroscience into tangible self-discovery.
                </Typography>
                
                <Divider sx={{ my: 1.5 }} />
                
                <Stack spacing={1.5}>
                  {[
                    { icon: <FlagIcon fontSize="small" />, text: 'Track your spectrum positions across the school year' },
                    { icon: <AchievementIcon fontSize="small" />, text: 'Collect stickers for shifting sliders and building skills' },
                    { icon: <RouteIcon fontSize="small" />, text: 'Chart your personal learning journey' },
                    { icon: <ConnectIcon fontSize="small" />, text: 'Share insights with teachers and parents' },
                  ].map((item, i) => (
                    <Box key={i} sx={{ display: 'flex', alignItems: 'flex-start', gap: 1.5 }}>
                      <Box sx={{ color: '#f59e0b', mt: 0.3 }}>{item.icon}</Box>
                      <Typography variant="body2">{item.text}</Typography>
                    </Box>
                  ))}
                </Stack>
                
                <Box sx={{ mt: 2, p: 1.5, bgcolor: alpha('#f59e0b', 0.06), borderRadius: 1 }}>
                  <Typography variant="caption" sx={{ fontStyle: 'italic', display: 'block' }}>
                    "I used to think I was bad at group work. Now I know my Social Gravity is just 'Independent' — 
                    and I can learn to switch settings when I need to."
                    <br />— Year 4 student
                  </Typography>
                </Box>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>

{/* ── The Core Mission: From Shame to Agency ── */}
      <Box sx={{ py: { xs: 6, md: 8 }, bgcolor: '#f5f7f9' }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Grid container spacing={4} alignItems="center">
            <Grid item xs={12} md={6}>
              <Typography variant="overline" sx={{ color: brandColor, fontWeight: 'bold' }}>
                The Central Aim
              </Typography>
              <Typography variant="h4" component="h2" fontWeight="bold" gutterBottom>
                From Shame to Agency
              </Typography>
              <Typography variant="body1" paragraph sx={{ fontSize: '1.1rem', lineHeight: 1.7 }}>
                Most classroom behaviour systems are built on shame — public charts, colour changes, 
                and the implicit message that some children are "good" and others are "bad."
              </Typography>
              <Typography variant="body1" paragraph sx={{ fontSize: '1.1rem', lineHeight: 1.7 }}>
                PrimaryOS replaces compliance with <strong>strategy</strong>. A child who interrupts isn't 
                "naughty" — they are running high voltage without an exhaust. A child who refuses group work isn't 
                "uncooperative" — they are preserving their battery.
              </Typography>
              <Typography variant="body1" sx={{ fontSize: '1.1rem', lineHeight: 1.7 }}>
                <strong>Shame has a place</strong> — regulating genuine harm. But for everything else, we move 
                from "Don't do that" to "How do we make this work?" We give children the manual to their own machine.
              </Typography>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Card elevation={4} sx={{ borderRadius: 3, bgcolor: 'white' }}>
                <CardContent sx={{ p: 4 }}>
                  <Typography variant="h6" gutterBottom fontWeight="bold" sx={{ color: brandColor }}>
                    The Vocabulary of Agency
                  </Typography>
                  <Stack spacing={2}>
                    {[
                      { 
                        before: 'Stop interrupting.', 
                        after: 'I get the enthusiasm—keep that energy. But you can\'t get anyone to hear you properly unless you listen first. That\'s how you win.' 
                      },
                      { 
                        before: 'Why can\'t you just join in?', 
                        after: 'I know the pack drains your battery. Just drop your idea in so we don\'t miss it, then you can pull back to recharge.' 
                      },
                      { 
                        before: 'You\'re being too sensitive.', 
                        after: 'You\'re picking up details everyone else is missing. That’s a superpower, but right now the signal is too loud. Let\'s step back.' 
                      },
                      { 
                        before: 'Stop making a fuss.', 
                        after: 'The plan is broken and I know that feels like chaos. We can\'t fix the schedule, but you can pick the one thing we keep predictable.' 
                      },
                    ].map((item, i) => (
                      <Paper
                        key={i}
                        elevation={0}
                        sx={{
                          p: 2,
                          bgcolor: alpha(brandColor, 0.04),
                          borderRadius: 2,
                          borderLeft: `4px solid ${brandColor}` // Added visual anchor for the 'Strategy'
                        }}
                      >
                        <Typography variant="caption" color="text.secondary" sx={{ textDecoration: 'line-through', display: 'block', mb: 1, opacity: 0.7 }}>
                          ❌ {item.before}
                        </Typography>
                        <Typography variant="body2" fontWeight="medium" sx={{ color: '#0f172a', lineHeight: 1.6 }}>
                          ✅ {item.after}
                        </Typography>
                      </Paper>
                    ))}
                  </Stack>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* ── The Two Skill Families ── */}
      <Box sx={{ py: { xs: 6, md: 8 } }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            Two Kinds of Skills, One Journey
          </Typography>
          <Typography
            variant="body1"
            color="text.secondary"
            textAlign="center"
            sx={{ mb: 5, maxWidth: 780, mx: 'auto' }}
          >
            PrimaryOS separates <strong>what you learn</strong> from <strong>how you learn it</strong>. 
            Both are tracked, celebrated, and developed through the same framework.
          </Typography>

          <Grid container spacing={4}>
            {/* ETP Skills Card */}
            <Grid item xs={12} md={6}>
              <Card
                elevation={3}
                sx={{
                  height: '100%',
                  borderTop: `8px solid ${brandColor}`,
                  borderRadius: 3,
                }}
              >
                <CardContent sx={{ p: 4 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 3 }}>
                    <Avatar sx={{ bgcolor: alpha(brandColor, 0.1), color: brandColor, width: 56, height: 56 }}>
                      <SwapIcon fontSize="large" />
                    </Avatar>
                    <Box>
                      <Typography variant="h5" fontWeight="bold">
                        Spectrum Skills
                      </Typography>
                      <Typography variant="subtitle2" color="text.secondary">
                        Learning to move your sliders
                      </Typography>
                    </Box>
                  </Box>

                  <Typography variant="body1" paragraph>
                    Each of the 9 ETP spectra has a corresponding <strong>spectrum skill</strong> — the ability 
                    to access different positions on that spectrum when the situation requires it.
                  </Typography>
                  
                  <Typography variant="body2" color="text.secondary" paragraph>
                    <strong>Example:</strong> A child with a strongly Independent Social Gravity setting 
                    isn't "wrong" — but they may need to access Cohesive settings during collaborative 
                    projects. The skill is <strong>Social Switching</strong>: knowing when to work alone 
                    and when to engage with others.
                  </Typography>

                  <Divider sx={{ my: 2 }} />

                  <Typography variant="subtitle2" gutterBottom fontWeight="bold">
                    What progress looks like:
                  </Typography>
                  <Stack spacing={1}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <LinearProgress 
                        variant="determinate" 
                        value={33} 
                        sx={{ flexGrow: 1, height: 8, borderRadius: 4, bgcolor: alpha(brandColor, 0.1) }}
                      />
                      <Typography variant="caption">Stage 1: Notices the setting exists</Typography>
                    </Box>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <LinearProgress 
                        variant="determinate" 
                        value={66} 
                        sx={{ flexGrow: 1, height: 8, borderRadius: 4, bgcolor: alpha(brandColor, 0.1) }}
                      />
                      <Typography variant="caption">Stage 2: Can shift with support</Typography>
                    </Box>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <LinearProgress 
                        variant="determinate" 
                        value={100} 
                        sx={{ flexGrow: 1, height: 8, borderRadius: 4, bgcolor: alpha(brandColor, 0.1) }}
                      />
                      <Typography variant="caption">Stage 3: Shifts independently when appropriate</Typography>
                    </Box>
                  </Stack>
                </CardContent>
              </Card>
            </Grid>

            {/* Foundational Skills Card */}
            <Grid item xs={12} md={6}>
              <Card
                elevation={3}
                sx={{
                  height: '100%',
                  borderTop: `8px solid #0284c7`,
                  borderRadius: 3,
                }}
              >
                <CardContent sx={{ p: 4 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 3 }}>
                    <Avatar sx={{ bgcolor: alpha('#0284c7', 0.1), color: '#0284c7', width: 56, height: 56 }}>
                      <ForestIcon fontSize="large" />
                    </Avatar>
                    <Box>
                      <Typography variant="h5" fontWeight="bold">
                        Foundational Skills
                      </Typography>
                      <Typography variant="subtitle2" color="text.secondary">
                        The building blocks of all learning
                      </Typography>
                    </Box>
                  </Box>

                  <Typography variant="body1" paragraph>
                    These are the <strong>fundamental competencies</strong> that underpin every subject — 
                    focus, self-regulation, social cognition, executive function, motor planning, and metacognition.
                  </Typography>
                  
                  <Typography variant="body2" color="text.secondary" paragraph>
                    <strong>Dual pathway:</strong> These skills can be developed through <strong>explicit exercises</strong> 
                    (a "focus gym" session) or <strong>embedded in curriculum</strong> (practising planning through 
                    a history project). Both count. Both are tracked.
                  </Typography>

                  <Divider sx={{ my: 2 }} />

                  <Typography variant="subtitle2" gutterBottom fontWeight="bold">
                    CHISG-aligned skill architecture:
                  </Typography>
                  <Grid container spacing={1} sx={{ mt: 1 }}>
                    {foundationalDomains.slice(0, 3).map((domain) => (
                      <Grid item xs={12} sm={4} key={domain.name}>
                        <Paper
                          elevation={0}
                          sx={{
                            p: 1.5,
                            bgcolor: alpha(domain.color, 0.04),
                            borderLeft: `3px solid ${domain.color}`,
                            height: '100%',
                          }}
                        >
                          <Typography variant="caption" fontWeight="bold" display="block">
                            {domain.name}
                          </Typography>
                          <Typography variant="caption" color="text.secondary">
                            {domain.skills.length} skills
                          </Typography>
                        </Paper>
                      </Grid>
                    ))}
                  </Grid>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* ── The 9 ETP Spectra with Skills ── */}
      <Box sx={{ py: { xs: 6, md: 8 }, bgcolor: '#f5f7f9' }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            The 9 Spectra — And the Skills to Navigate Them
          </Typography>
          <Typography
            variant="body1"
            color="text.secondary"
            textAlign="center"
            sx={{ mb: 2, maxWidth: 780, mx: 'auto' }}
          >
            Every child has default positions. The goal isn't to change their defaults — 
            it's to give them the ability to <strong>choose different settings when needed</strong>.
          </Typography>

          <Grid container spacing={2.5} sx={{ mt: 2 }}>
            {primarySpectra.map((s) => (
              <Grid item xs={12} md={6} lg={4} key={s.name}>
                <Card
                  elevation={2}
                  sx={{
                    height: '100%',
                    borderTop: `4px solid ${s.color}`,
                    borderRadius: 2,
                  }}
                >
                  <CardContent sx={{ p: 3 }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 1.5 }}>
                      <Typography variant="subtitle1" fontWeight="bold">
                        {s.name}
                      </Typography>
                      <Chip 
                        label={s.skill} 
                        size="small" 
                        sx={{ bgcolor: alpha(s.color, 0.1), color: s.color, fontWeight: 'medium' }} 
                      />
                    </Box>

                    {/* Spectrum bar */}
                    <Box sx={{ mb: 2 }}>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
                        <Typography variant="caption">{s.left}</Typography>
                        <Typography variant="caption">{s.right}</Typography>
                      </Box>
                      <Box
                        sx={{
                          height: 6,
                          borderRadius: 3,
                          background: `linear-gradient(90deg, ${alpha(s.color, 0.3)}, ${s.color})`,
                        }}
                      />
                    </Box>

                    <Typography variant="body2" color="text.secondary" sx={{ mb: 1.5 }}>
                      {s.description}
                    </Typography>

                    <Box sx={{ mt: 1, p: 1.5, bgcolor: alpha(s.color, 0.04), borderRadius: 1 }}>
                      <Typography variant="caption" fontWeight="bold" display="block" gutterBottom>
                        🎯 The skill: {s.skill}
                      </Typography>
                      <Typography variant="caption" color="text.secondary" display="block">
                        {s.skillDescription}
                      </Typography>
                      <Typography variant="caption" color="text.secondary" display="block" sx={{ mt: 1, fontStyle: 'italic' }}>
                        📋 Classroom: {s.classroom}
                      </Typography>
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* ── Neuron Navigators Journey ── */}
      <Box sx={{ py: { xs: 6, md: 8 } }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Box sx={{ textAlign: 'center', mb: 4 }}>
            <JournalIcon sx={{ fontSize: 56, color: brandColor, mb: 2 }} />
            <Typography variant="h4" component="h2" fontWeight="bold" gutterBottom>
              The Neuron Navigators Guide Book
            </Typography>
            <Typography variant="h6" color="text.secondary" sx={{ fontWeight: 'normal', maxWidth: 700, mx: 'auto' }}>
              A personal journal that grows with the child — from Reception to Year 6
            </Typography>
          </Box>

          <Grid container spacing={3} sx={{ mb: 5 }}>
            {journeyStages.map((stage) => (
              <Grid item xs={12} md={3} key={stage.stage}>
                <Card
                  elevation={2}
                  sx={{
                    height: '100%',
                    borderTop: `4px solid ${stage.color}`,
                    borderRadius: 2,
                  }}
                >
                  <CardContent>
                    <Typography variant="overline" sx={{ color: stage.color, fontWeight: 'bold' }}>
                      {stage.stage}
                    </Typography>
                    <Typography variant="h6" fontWeight="bold" gutterBottom>
                      {stage.title}
                    </Typography>
                    <Typography variant="caption" display="block" sx={{ mb: 1 }}>
                      {stage.age}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" paragraph>
                      {stage.description}
                    </Typography>
                    
                    <Divider sx={{ my: 1.5 }} />
                    
                    <Typography variant="caption" fontWeight="bold" display="block" gutterBottom>
                      Sticker examples:
                    </Typography>
                    <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, mb: 1.5 }}>
                      {stage.stickers.map((sticker) => (
                        <Chip
                          key={sticker}
                          label={sticker}
                          size="small"
                          sx={{ fontSize: '0.6rem', bgcolor: alpha(stage.color, 0.1) }}
                        />
                      ))}
                    </Box>
                    
                    <Typography variant="caption" display="block">
                      <strong>👩‍🏫 Teacher:</strong> {stage.teacherRole}
                    </Typography>
                    <Typography variant="caption" display="block">
                      <strong>🧑‍🎓 Student:</strong> {stage.studentRole}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            ))}
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
            <Grid container spacing={3} alignItems="center">
              <Grid item xs={12} md={8}>
                <Typography variant="h6" gutterBottom fontWeight="bold">
                  Every sticker tells a story
                </Typography>
                <Typography variant="body2" paragraph>
                  The Neuron Navigators Guide isn't a reward chart. It's a <strong>scientific field journal</strong>. 
                  Students collect stickers when they notice their own settings, successfully shift a slider, 
                  master a foundational skill, or help a peer understand their own wiring.
                </Typography>
                <Typography variant="body2">
                  Over seven years, the journal becomes a complete map of the child's development — 
                  not just what they learned, but <strong>how they learned to learn</strong>.
                </Typography>
              </Grid>
              <Grid item xs={12} md={4}>
                <Box sx={{ display: 'flex', justifyContent: 'center', gap: 1, flexWrap: 'wrap' }}>
                  <Avatar sx={{ bgcolor: '#2e7d32', width: 48, height: 48 }}>🧠</Avatar>
                  <Avatar sx={{ bgcolor: '#1565c0', width: 48, height: 48 }}>⚡</Avatar>
                  <Avatar sx={{ bgcolor: '#6a1b9a', width: 48, height: 48 }}>🧭</Avatar>
                  <Avatar sx={{ bgcolor: '#c62828', width: 48, height: 48 }}>🎯</Avatar>
                  <Avatar sx={{ bgcolor: '#f57c00', width: 48, height: 48 }}>📊</Avatar>
                </Box>
              </Grid>
            </Grid>
          </Paper>
        </Container>
      </Box>

      {/* ── Teacher as Facilitator ── */}
      <Box sx={{ py: { xs: 6, md: 8 }, bgcolor: alpha(brandColor, 0.04) }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Grid container spacing={4} alignItems="center">
            <Grid item xs={12} md={6}>
              <Typography variant="h4" component="h2" fontWeight="bold" gutterBottom>
                From Judge to Guide
              </Typography>
              <Typography variant="h6" color="text.secondary" sx={{ fontWeight: 'normal', mb: 3 }}>
                PrimaryOS shifts the teacher's role from behaviour adjudicator to learning facilitator.
              </Typography>
              
              <Stack spacing={2}>
                {[
                  {
                    icon: <VisibilityIcon />,
                    title: 'Observer',
                    description: 'Your primary task is noticing patterns. Which settings does this child default to? Which situations cause friction? Where are their sliders sticky?',
                  },
                  {
                    icon: <LightbulbIcon />,
                    title: 'Vocabulary Provider',
                    description: 'You give students the words to describe their own experience. "It sounds like your Voltage Sensitivity is high right now."',
                  },
                  {
                    icon: <ConnectIcon />,
                    title: 'Strategy Coach',
                    description: 'You don\'t fix problems — you equip students with tools to fix them themselves. "When your Risk Aversion is blocking you, what\'s one small step you could try?"',
                  },
                  {
                    icon: <AutoFixIcon />,
                    title: 'Scaffold Remover',
                    description: 'Your goal is to work yourself out of a job. Every student should leave your classroom more capable of piloting themselves than when they arrived.',
                  },
                ].map((item, i) => (
                  <Box key={i} sx={{ display: 'flex', gap: 2 }}>
                    <Avatar sx={{ bgcolor: alpha(brandColor, 0.1), color: brandColor }}>
                      {item.icon}
                    </Avatar>
                    <Box>
                      <Typography variant="subtitle1" fontWeight="bold">{item.title}</Typography>
                      <Typography variant="body2" color="text.secondary">{item.description}</Typography>
                    </Box>
                  </Box>
                ))}
              </Stack>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Card elevation={4} sx={{ borderRadius: 3 }}>
                <CardContent sx={{ p: 4 }}>
                  <Typography variant="h6" gutterBottom fontWeight="bold" sx={{ color: brandColor }}>
                    Workload, Not Increased — Transformed
                  </Typography>
                  
                  <Typography variant="body2" paragraph>
                    PrimaryOS doesn't add to teacher workload — it redirects it. The hours spent on 
                    behaviour management, parental complaints about "unfair" treatment, and repetitive 
                    low-level corrections become focused, strategic observations.
                  </Typography>
                  
                  <Divider sx={{ my: 2 }} />
                  
                  <Typography variant="subtitle2" gutterBottom fontWeight="bold">
                    How this changes the classroom:
                  </Typography>
                  
                  <Stack spacing={2}>
                    <Paper elevation={0} sx={{ p: 2, bgcolor: alpha(brandColor, 0.04), borderRadius: 2 }}>
                      <Typography variant="body2" sx={{ fontStyle: 'italic' }}>
                        Instead of spending time on "Stop calling out," teachers can say 
                        "I can see you're processing out loud — jot that thought down and we'll come back to it." 
                        The behaviour hasn't changed — the response has. And the PrimaryOS framework is designed to make that shift natural.
                      </Typography>
                    </Paper>
                    
                    <Paper elevation={0} sx={{ p: 2, bgcolor: alpha(brandColor, 0.04), borderRadius: 2 }}>
                      <Typography variant="body2" sx={{ fontStyle: 'italic' }}>
                        The Neuron Navigators book is designed as a morning check-in tool. Children arrive, open their 
                        journal, and place their sticker. Teachers can see at a glance who's regulated and who needs 
                        support — changing how they start their day.
                      </Typography>
                    </Paper>
                  </Stack>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* ── ParentOS Bridge ── */}
      <Box sx={{ py: { xs: 6, md: 8 } }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Paper
            elevation={3}
            sx={{
              p: 4,
              borderRadius: 4,
              background: `linear-gradient(135deg, ${alpha('#0f766e', 0.04)} 0%, ${alpha(brandColor, 0.04)} 100%)`,
              border: `1px solid ${alpha('#0f766e', 0.2)}`,
            }}
          >
            <Grid container spacing={4} alignItems="center">
              <Grid item xs={12} md={8}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                  <SchoolIcon sx={{ color: '#0f766e', fontSize: 32 }} />
                  <Typography variant="overline" sx={{ color: '#0f766e', fontWeight: 'bold' }}>
                    Connected by Design
                  </Typography>
                </Box>
                <Typography variant="h5" fontWeight="bold" gutterBottom>
                  One Child. One Nervous System. One Framework.
                </Typography>
                <Typography variant="body1" paragraph>
                  PrimaryOS and ParentOS speak the same language. The Social Switching skill a child practises 
                  in a Year 4 group project is the same Social Switching skill their parents see at the dinner table.
                </Typography>
                <Typography variant="body1">
                  When school and home share vocabulary, children don't have to maintain two separate selves. 
                  Insights from the classroom inform parenting strategies. Observations from home inform teacher 
                  support. The child is seen as whole.
                </Typography>
              </Grid>
              <Grid item xs={12} md={4}>
                <Box sx={{ textAlign: 'center' }}>
                  <Avatar sx={{ bgcolor: '#0f766e', width: 64, height: 64, mx: 'auto', mb: 2 }}>
                    <ParentIcon />
                  </Avatar>
                  <Typography variant="subtitle1" fontWeight="bold">ParentOS</Typography>
                  <Typography variant="caption" display="block" color="text.secondary" sx={{ mb: 2 }}>
                    The Foundational Code Your Family Runs On
                  </Typography>
                  <Button
                    variant="outlined"
                    size="small"
                    component={Link}
                    to="/parent-os"
                    sx={{ color: '#0f766e', borderColor: '#0f766e' }}
                  >
                    Learn More
                  </Button>
                </Box>
              </Grid>
            </Grid>
          </Paper>
        </Container>
      </Box>

      {/* ── What Makes PrimaryOS Different ── */}
      <Box sx={{ py: { xs: 6, md: 8 }, bgcolor: '#f5f7f9' }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            What Makes PrimaryOS Different
          </Typography>

          <Grid container spacing={2} sx={{ mt: 2 }}>
            {[
              'Replaces shame with vocabulary — behaviour is communication, not defiance',
              'Teaches spectrum skills, not just content knowledge',
              'Student-owned journal tracks development across 7 years',
              'Teacher shifts from judge to facilitator',
              'CHISG-aligned skill architecture prevents fragmentation',
              'Explicit foundational skills + curriculum-embedded practice',
              '9 ETP spectra explain why learning feels different for different children',
              'Home-school bridge through ParentOS alignment',
              'Sticker book is a scientific journal, not a reward chart',
              'No child is "bad" — only mismatched or stuck',
              'Built on 12 years of classroom observation, not theory',
              'Reduces teacher workload by transforming behaviour management',
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
          background: `linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)`,
          color: 'white',
        }}
      >
        <Container maxWidth="md">
          <SchoolIcon sx={{ fontSize: 56, mb: 2 }} />
          <Typography variant="h4" fontWeight="bold" gutterBottom>
            Give Every Child the Manual to Their Own Mind
          </Typography>
          <Typography variant="body1" sx={{ mb: 4, maxWidth: 600, mx: 'auto', opacity: 0.95 }}>
            PrimaryOS is being developed for schools as a complete framework — including teacher training, 
            Neuron Navigators journals, and CHISG-aligned skill tracking.
          </Typography>
          
          <Stack direction="row" spacing={2} justifyContent="center" flexWrap="wrap" useFlexGap>
            <ContactReveal email="world@espthinking.co.uk" label="Request School Demo" />
            {/* <Button
              variant="outlined"
              size="large"
              component={Link}
              to="/primaryos/neuron-navigators"
              sx={{
                px: 5,
                py: 1.5,
                fontWeight: 'bold',
                color: 'white',
                borderColor: 'white',
                '&:hover': { borderColor: 'white', bgcolor: alpha('#ffffff', 0.1) },
              }}
            >
              Explore Neuron Navigators
            </Button> */}
          </Stack>
          
          <Typography variant="caption" sx={{ display: 'block', mt: 4, opacity: 0.8 }}>
            Single-school and MAT-wide licensing planned. Neuron Navigators journals will be available in class sets.
          </Typography>
        </Container>
      </Box>

      {/* ── Footer ── */}
      <Box sx={{ py: 3, textAlign: 'center', bgcolor: '#f59e0b', color: 'rgba(255,255,255,0.9)' }}>
        <Typography variant="body2">
          PrimaryOS — The Operating System for Growing Learners — Part of the ESP Thinking Portfolio
        </Typography>
        <Typography variant="caption" sx={{ display: 'block', mt: 0.5, opacity: 0.8 }}>
          In partnership with Neuron Navigators Guide Book | Connected to ParentOS | Grounded in CHISG
        </Typography>
      </Box>
    </Box>
  );
};

export default PrimaryOSLanding;