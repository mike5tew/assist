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
} from '@mui/material';
import {
  Home as HomeIcon,
  FamilyRestroom as FamilyIcon,
  ChildCare as ChildIcon,
  Psychology as PsychologyIcon,
  School as SchoolIcon,
  Visibility as VisibilityIcon,
  Balance as BalanceIcon,
  Favorite as FavoriteIcon,
  Security as SecurityIcon,
  ConnectWithoutContact as ConnectIcon,
  Lightbulb as LightbulbIcon,
  AutoFixHigh as AutoFixIcon,
  Timeline as TimelineIcon,
  Groups as GroupsIcon,
  EmojiPeople as PeopleIcon,
  Forest as ForestIcon,
  Explore as ExploreIcon,
  MenuBook as BookIcon,
  QuestionMark as QuestionIcon,
  Warning as WarningIcon,
  ArrowForward as ArrowIcon,
  LockOpen as LockOpenIcon,
  SelfImprovement as SelfIcon,
  FamilyRestroom as ParentChildIcon,
  Verified as VerifiedIcon,
} from '@mui/icons-material';
import { Link } from 'react-router-dom';
import ContactReveal from '../../components/ContactReveal';

// ─── The 9 Spectra for Parents ──────────────────────────────
const parentSpectra = [
  {
    name: 'Social Gravity',
    left: 'Independent',
    right: 'Cohesive',
    description: 'Does your child recharge in solitude or with others? Neither is "better" — they\'re different charging modes.',
    color: '#1565c0',
    ageHint: 'Observable from ~3 years',
    question: 'After a busy day, does your child seek quiet time or immediately want to connect?',
  },
  {
    name: 'Energy Directionality',
    left: 'Inward',
    right: 'Outward',
    description: 'Does your child think by processing internally, or by talking it through?',
    color: '#6a1b9a',
    ageHint: 'Observable from ~4 years',
    question: 'When solving a problem, does your child go quiet or talk through options aloud?',
  },
  {
    name: 'Voltage Sensitivity',
    left: 'Insulated',
    right: 'Conductive',
    description: 'How much emotional current does your child absorb from others?',
    color: '#f57c00',
    ageHint: 'Observable from ~3 years',
    question: 'Does your child "catch" emotions easily from peers, or stay steady in chaos?',
  },
  {
    name: 'Threat Response',
    left: 'Passive',
    right: 'Aggressive',
    description: 'How does your child react when they feel threatened or unsafe?',
    color: '#c62828',
    ageHint: 'Observable from ~2 years',
    question: 'When frightened, does your child freeze/retreat or fight back loudly?',
  },
  {
    name: 'Care Response',
    left: 'Detached',
    right: 'Nurturing',
    description: 'How does your child respond to others\' vulnerability?',
    color: '#ad1457',
    ageHint: 'Observable from ~3 years',
    question: 'When someone is upset, does your child offer comfort or give space?',
  },
  {
    name: 'Risk Tolerance',
    left: 'Averse',
    right: 'Seeking',
    description: 'How does your child approach uncertainty and physical challenge?',
    color: '#2e7d32',
    ageHint: 'Observable from ~2 years',
    question: 'Does your child approach new playground equipment cautiously or dive straight in?',
  },
  {
    name: 'Integrity Logic',
    left: 'Relativistic',
    right: 'Absolutist',
    description: 'How does your child process rules and fairness?',
    color: '#4527a0',
    ageHint: 'Observable from ~4 years',
    question: 'Does your child see rules as flexible guidelines or unbreakable truths?',
  },
  {
    name: 'Mirror Neuron Tuning',
    left: 'Selective',
    right: 'Absorbent',
    description: 'How much does your child mirror the emotional states of others?',
    color: '#00838f',
    ageHint: 'Observable from ~3 years',
    question: 'Does your child stay grounded in group chaos, or absorb the room\'s energy?',
  },
  {
    name: 'Orderliness',
    left: 'Flexible',
    right: 'Ordered',
    description: 'How much structure does your child need to feel secure?',
    color: '#5d4037',
    ageHint: 'Observable from ~2 years',
    question: 'Does your child adapt easily to schedule changes, or need predictability?',
  },
];

// ─── Slider States ──────────────────────────────────────────
const sliderStates = [
  {
    state: 'Fluid',
    icon: <LockOpenIcon />,
    description: 'Moves freely — healthy. Can adapt to context.',
    color: '#2e7d32',
    parentHint: 'Your child can access different positions depending on the situation.',
  },
  {
    state: 'Sticky',
    icon: <TimelineIcon />,
    description: 'Moves with difficulty — developing. Needs practice to shift.',
    color: '#f57c00',
    parentHint: 'Your child prefers one position but can shift with support.',
  },
  {
    state: 'Stuck',
    icon: <WarningIcon />,
    description: 'Cannot move — frozen at one position regardless of context.',
    color: '#c62828',
    parentHint: 'Your child seems unable to access other settings, even when clearly needed.',
  },
  {
    state: 'Hijacked',
    icon: <AutoFixIcon />,
    description: 'Moves involuntarily — triggered. External stimuli move the slider without consent.',
    color: '#6a1b9a',
    parentHint: 'Your child\'s reactions are driven by triggers, not choice.',
  },
];

// ─── Age Stage Summaries ───────────────────────────────────
const ageStages = [
  {
    range: '3–5 years',
    title: 'Emerging Patterns',
    description: 'Spectrum positions become observable but remain highly fluid. The goal is curiosity, not diagnosis.',
    focus: 'Notice without judging. What settings does your child naturally prefer?',
    color: '#2e7d32',
  },
  {
    range: '5–7 years',
    title: 'Consolidating Tendencies',
    description: 'Consistent patterns emerge. Children develop preferences and can begin to recognise their own settings.',
    focus: 'Introduce spectrum vocabulary. "You seem to need quiet time right now."',
    color: '#1565c0',
  },
  {
    range: '7–9 years',
    title: 'Flexibility Training',
    description: 'Children can learn to shift positions with support. Sticky sliders become targets for gentle practice.',
    focus: 'Co-regulation and scaffolded flexibility. "Let\'s try the brave setting together."',
    color: '#6a1b9a',
  },
  {
    range: '9–11 years',
    title: 'Self-Awareness',
    description: 'Children can articulate their own settings and begin self-regulating with adult guidance.',
    focus: 'Transfer Pilot control from parent to child. "What setting do you need right now?"',
    color: '#00838f',
  },
  {
    range: '11–13 years',
    title: 'Peer Gravity Shift',
    description: 'Social gravity shifts toward peers. Voltage sensitivity amplifies. Hormonal modulation begins.',
    focus: 'Maintain connection while respecting autonomy. Normalise fluctuating settings.',
    color: '#ad1457',
  },
  {
    range: '13–16 years',
    title: 'Pilot Development',
    description: 'Executive function capacity expands. Teens can consciously choose positions with decreasing support.',
    focus: 'Transfer full agency gradually. "You\'re the pilot now — I\'m your co-pilot."',
    color: '#4527a0',
  },
];

// ─── Component ──────────────────────────────────────────────
const ParentOSLanding: React.FC = () => {
  const theme = useTheme();
  const brandColor = '#0f766e';

  return (
    <Box>
      {/* ── Hero Section ── */}
      <Box
        sx={{
          background: `linear-gradient(135deg, #0f766e 0%, #14b8a6 50%, #2dd4bf 100%)`,
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
              <Box
                        component="img"
                        src={`${process.env.PUBLIC_URL}/ESPLogoLong.png`}
                        alt="ESP Thinking Portfolio"
                        sx={{ height: 56, width: 'auto' }}
                      />
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

                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                  <FamilyIcon sx={{ fontSize: { xs: 48, md: 64 } }} />
                  <Box>
                    <Typography variant="h2" component="h1" fontWeight="bold" sx={{ lineHeight: 1.1 }}>
                      ParentOS
                    </Typography>
                    <Typography variant="subtitle1" sx={{ opacity: 0.85, mt: 0.5 }}>
                      The Foundational Code Your Family Runs On
                    </Typography>
                  </Box>
                </Box>

                <Typography variant="h5" sx={{ opacity: 0.92, maxWidth: 620 }}>
                  A developmental framework for parents who want to understand their child's wiring — 
                  not fix it. ParentOS translates 12 years of classroom observation into a practical 
                  lens for seeing behaviour as communication, not defiance.
                </Typography>

                <Stack direction="row" spacing={1.5} flexWrap="wrap" useFlexGap>
                  <Chip label="9 Biological Spectra" sx={{ bgcolor: 'rgba(255,255,255,0.15)', color: 'white' }} />
                  <Chip label="Ages 3–16" sx={{ bgcolor: 'rgba(255,255,255,0.15)', color: 'white' }} />
                  <Chip label="Free Foundation" sx={{ bgcolor: 'rgba(255,255,255,0.15)', color: 'white' }} />
                  <Chip label="CHISG-Aligned" sx={{ bgcolor: 'rgba(255,255,255,0.15)', color: 'white' }} />
                </Stack>

                <Stack direction="row" spacing={2} sx={{ pt: 1 }}>
                  <Button
                    variant="contained"
                    size="large"
                    component={Link}
                    to="/etp-profile"
                    sx={{
                      bgcolor: 'white',
                      color: '#0f766e',
                      fontWeight: 'bold',
                      px: 4,
                      '&:hover': { bgcolor: alpha('#ffffff', 0.9) },
                    }}
                  >
                    Launch Spectrum Explorer
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
                  <WarningIcon sx={{ color: '#f57c00' }} />
                  <Typography variant="h6" color="text.primary" fontWeight="bold">
                    How to Use This Framework
                  </Typography>
                </Box>
                
                <Typography variant="body2" color="text.secondary" paragraph>
                  <strong>ParentOS is not a diagnostic tool.</strong> It is a lens — one of many. 
                  Children develop at different rates, and every spectrum position has strengths 
                  and challenges. The goal is never to "fix" your child's settings, but to understand 
                  them so you can work with their wiring, not against it.
                </Typography>
                
                <Divider sx={{ my: 1.5 }} />
                
                <Stack spacing={1.5}>
                  {[
                    { icon: <VisibilityIcon fontSize="small" />, text: 'Observe, don\'t diagnose' },
                    { icon: <BalanceIcon fontSize="small" />, text: 'No position is "good" or "bad"' },
                    { icon: <ConnectIcon fontSize="small" />, text: 'Behaviour is communication, not defiance' },
                    { icon: <SchoolIcon fontSize="small" />, text: 'Bridges home observation to school data' },
                  ].map((item, i) => (
                    <Box key={i} sx={{ display: 'flex', alignItems: 'flex-start', gap: 1.5 }}>
                      <Box sx={{ color: '#0f766e', mt: 0.3 }}>{item.icon}</Box>
                      <Typography variant="body2">{item.text}</Typography>
                    </Box>
                  ))}
                </Stack>
                
                <Typography variant="body2" sx={{ mt: 2, fontStyle: 'italic', color: 'text.secondary' }}>
                  "If you have genuine concerns about your child's development, speak to a health visitor, GP, or specialist teacher. ParentOS complements professional judgment — it doesn't replace it."
                </Typography>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* ── The Core Philosophy: Curiosity, Not Correction ── */}
      <Box sx={{ py: { xs: 6, md: 8 }, bgcolor: '#f5f7f9' }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            Curiosity, Not Correction
          </Typography>
          <Typography
            variant="body1"
            color="text.secondary"
            textAlign="center"
            sx={{ mb: 5, maxWidth: 780, mx: 'auto' }}
          >
            Most parenting advice tells you what to <em>do</em>. ParentOS helps you understand 
            what you're <em>seeing</em>. Behaviour is the output of an internal system — when you 
            understand the system, the "what to do" becomes obvious.
          </Typography>

          <Grid container spacing={3}>
            {[
              {
                icon: <ChildIcon fontSize="large" />,
                title: 'See the Child, Not the Behaviour',
                description: 'A meltdown isn\'t defiance — it\'s a system overload. The question shifts from "How do I make this stop?" to "Which spectrum is maxed out and why?"',
              },
              {
                icon: <ForestIcon fontSize="large" />,
                title: 'Range of Motion, Not "Correct" Position',
                description: 'The goal is never to push your child to the "healthy" end of a spectrum. It\'s to expand their range — so they can choose the right position for the right situation.',
              },
              {
                icon: <LockOpenIcon fontSize="large" />,
                title: 'Stuck Isn\'t Broken',
                description: 'A stuck slider isn\'t a defect. It\'s a pattern that made sense once. The question is whether it still serves your child — and if not, how to gently expand their options.',
              },
              {
                icon: <ConnectIcon fontSize="large" />,
                title: 'The Home-School Bridge',
                description: 'Your child has one nervous system. It doesn\'t switch off at the school gate. ParentOS gives you a common language with teachers — so home and school can work together, not in parallel.',
              },
              {
                icon: <SelfIcon fontSize="large" />,
                title: 'Your Wiring Matters Too',
                description: 'ParentOS includes tools for parents to understand their own spectra. Compatibility friction happens when your settings clash with your child\'s — and that\'s nobody\'s fault.',
              },
              {
                icon: <SchoolIcon fontSize="large" />,
                title: 'From PrimaryOS to ParentOS',
                description: 'The same ETP framework used in classrooms now available for families. Your child\'s teacher may already be observing the same patterns — ParentOS helps you see what they see.',
              },
            ].map((card) => (
              <Grid item xs={12} md={6} key={card.title}>
                <Card
                  elevation={0}
                  sx={{
                    height: '100%',
                    border: `1px solid ${alpha(brandColor, 0.12)}`,
                    transition: 'transform 0.2s, box-shadow 0.2s',
                    '&:hover': {
                      transform: 'translateY(-3px)',
                      boxShadow: theme.shadows[4],
                    },
                  }}
                >
                  <CardContent sx={{ p: 3 }}>
                    <Box sx={{ display: 'flex', gap: 2, mb: 1.5 }}>
                      <Box
                        sx={{
                          p: 1.5,
                          borderRadius: 2,
                          bgcolor: alpha(brandColor, 0.08),
                          color: brandColor,
                          display: 'flex',
                          alignItems: 'center',
                        }}
                      >
                        {card.icon}
                      </Box>
                      <Typography variant="h6" fontWeight="bold" sx={{ display: 'flex', alignItems: 'center' }}>
                        {card.title}
                      </Typography>
                    </Box>
                    <Typography variant="body2" color="text.secondary" sx={{ lineHeight: 1.7 }}>
                      {card.description}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* ── The 9 Spectra for Parents ── */}
      <Box sx={{ py: { xs: 6, md: 8 } }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            The 9 Biological Spectra
          </Typography>
          <Typography
            variant="body1"
            color="text.secondary"
            textAlign="center"
            sx={{ mb: 2, maxWidth: 740, mx: 'auto' }}
          >
            Every child is born with these nine sliders. Their positions are not flaws — they're 
            your child's unique operating system. Your job isn't to reprogram them; it's to learn 
            the code so you can communicate effectively.
          </Typography>
          <Typography
            variant="body2"
            color="text.secondary"
            textAlign="center"
            sx={{ mb: 5, maxWidth: 600, mx: 'auto', fontStyle: 'italic' }}
          >
            These patterns become observable from age 2–3. Before that, focus on connection, not categorization.
          </Typography>

          <Grid container spacing={2.5}>
            {parentSpectra.map((s) => (
              <Grid item xs={12} sm={6} md={4} key={s.name}>
                <Paper
                  elevation={2}
                  sx={{
                    p: 3,
                    height: '100%',
                    borderTop: `4px solid ${s.color}`,
                    borderRadius: 2,
                    display: 'flex',
                    flexDirection: 'column',
                  }}
                >
                  <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1.5 }}>
                    <Typography variant="subtitle1" fontWeight="bold">
                      {s.name}
                    </Typography>
                    <Chip label={s.ageHint} size="small" sx={{ fontSize: '0.6rem', bgcolor: alpha(s.color, 0.1), color: s.color }} />
                  </Box>

                  {/* Spectrum visual bar */}
                  <Box sx={{ mb: 1.5 }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
                      <Typography variant="caption" fontWeight="medium">{s.left}</Typography>
                      <Typography variant="caption" fontWeight="medium">{s.right}</Typography>
                    </Box>
                    <Box
                      sx={{
                        height: 6,
                        borderRadius: 3,
                        background: `linear-gradient(90deg, ${alpha(s.color, 0.3)}, ${s.color})`,
                      }}
                    />
                  </Box>

                  <Typography variant="body2" color="text.secondary" sx={{ lineHeight: 1.6, mb: 1.5 }}>
                    {s.description}
                  </Typography>
                  
                  <Box sx={{ mt: 'auto', pt: 1, borderTop: `1px solid ${alpha(theme.palette.divider, 0.5)}` }}>
                    <Typography variant="caption" color="text.secondary" sx={{ fontStyle: 'italic', display: 'block' }}>
                      <strong>Try asking:</strong> {s.question}
                    </Typography>
                  </Box>
                </Paper>
              </Grid>
            ))}
          </Grid>

          <Box sx={{ textAlign: 'center', mt: 5 }}>
            <Button
              variant="contained"
              size="large"
              component={Link}
              to="/etp-profile"
              startIcon={<ExploreIcon />}
              sx={{ px: 5, py: 1.5, fontWeight: 'bold', bgcolor: brandColor }}
            >
              Explore the Spectra in Detail →
            </Button>
          </Box>
        </Container>
      </Box>

      {/* ── Slider States for Parents ── */}
      <Box sx={{ py: { xs: 6, md: 8 }, bgcolor: '#f5f7f9' }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            Understanding Slider States
          </Typography>
          <Typography
            variant="body1"
            color="text.secondary"
            textAlign="center"
            sx={{ mb: 5, maxWidth: 700, mx: 'auto' }}
          >
            Where your child sits on a spectrum matters less than whether they can <em>move</em>. 
            A child who can access different positions for different contexts is developing healthy 
            self-regulation.
          </Typography>

          <Grid container spacing={3} justifyContent="center">
            {sliderStates.map((s) => (
              <Grid item xs={12} sm={6} md={3} key={s.state}>
                <Paper
                  elevation={1}
                  sx={{
                    p: 3,
                    height: '100%',
                    textAlign: 'center',
                    borderRadius: 2,
                    borderTop: `4px solid ${s.color}`,
                  }}
                >
                  <Box sx={{ color: s.color, mb: 1 }}>{s.icon}</Box>
                  <Typography variant="h6" fontWeight="bold" gutterBottom>
                    {s.state}
                  </Typography>
                  <Typography variant="body2" color="text.secondary" sx={{ mb: 1.5 }}>
                    {s.description}
                  </Typography>
                  <Typography variant="caption" color="text.secondary" sx={{ fontStyle: 'italic' }}>
                    {s.parentHint}
                  </Typography>
                </Paper>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* ── Age Stage Roadmap ── */}
      <Box sx={{ py: { xs: 6, md: 8 } }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            Development as a Journey, Not a Checklist
          </Typography>
          <Typography
            variant="body1"
            color="text.secondary"
            textAlign="center"
            sx={{ mb: 2, maxWidth: 700, mx: 'auto' }}
          >
            These are not deadlines. They're general patterns — your child will arrive on their own timeline.
            Use this as a roadmap, not a stopwatch.
          </Typography>

          <Grid container spacing={2} sx={{ mt: 2 }}>
            {ageStages.map((stage) => (
              <Grid item xs={12} sm={6} md={4} key={stage.range}>
                <Box
                  sx={{
                    p: 2.5,
                    borderLeft: `4px solid ${stage.color}`,
                    bgcolor: alpha(stage.color, 0.04),
                    borderRadius: '0 8px 8px 0',
                    height: '100%',
                  }}
                >
                  <Typography variant="overline" sx={{ color: stage.color, fontWeight: 'bold' }}>
                    {stage.range}
                  </Typography>
                  <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
                    {stage.title}
                  </Typography>
                  <Typography variant="body2" color="text.secondary" sx={{ mb: 1, lineHeight: 1.6 }}>
                    {stage.description}
                  </Typography>
                  <Typography variant="caption" sx={{ display: 'block', mt: 1, fontStyle: 'italic' }}>
                    <strong>Focus:</strong> {stage.focus}
                  </Typography>
                </Box>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* ── Parent Self-Regulation ── */}
      <Box sx={{ py: { xs: 6, md: 8 }, bgcolor: alpha(brandColor, 0.04), borderY: `1px solid ${alpha(brandColor, 0.1)}` }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Grid container spacing={4} alignItems="center">
            <Grid item xs={12} md={6}>
              <Typography variant="h4" component="h2" fontWeight="bold" gutterBottom>
                Your Wiring Matters Too
              </Typography>
              <Typography variant="body1" paragraph sx={{ fontSize: '1.1rem', lineHeight: 1.7 }}>
                Parent-child friction isn't usually because one of you is "wrong." It's because your 
                biological settings are incompatible for that specific situation.
              </Typography>
              <Typography variant="body1" sx={{ fontSize: '1.1rem', lineHeight: 1.7, mb: 3 }}>
                An orderly parent with a flexible child isn't a parenting failure — it's a compatibility 
                challenge. ParentOS helps you identify where your settings clash, so you can design 
                protocols that work for both of you.
              </Typography>
              
              <Paper elevation={1} sx={{ p: 3, bgcolor: 'background.paper', borderRadius: 2 }}>
                <Typography variant="subtitle1" fontWeight="bold" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <BalanceIcon sx={{ color: brandColor }} /> Compatibility, Not Correction
                </Typography>
                <Typography variant="body2">
                  When your social gravity setting is "Independent" and your child's is "Cohesive," 
                  nobody needs fixing. You need a protocol: scheduled connection time, clear recharge 
                  signals, and shared vocabulary about what each of you needs.
                </Typography>
              </Paper>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Paper
                elevation={3}
                sx={{
                  p: 4,
                  borderRadius: 3,
                  bgcolor: 'background.paper',
                  border: `1px solid ${alpha(brandColor, 0.2)}`,
                }}
              >
                <Typography variant="h6" gutterBottom fontWeight="bold" sx={{ color: brandColor }}>
                  The Pilot & The Animal
                </Typography>
                <Typography variant="body2" paragraph>
                  <strong>Your child's Animal</strong> is their biological hardware — the 9 sliders that 
                  respond automatically to stimuli. <strong>Their Pilot</strong> is their developing 
                  executive function, the capacity to choose where to position each slider.
                </Typography>
                <Typography variant="body2" paragraph>
                  <strong>Your role</strong> is to be their co-pilot. You provide the executive function 
                  they haven't yet developed. You help them recognise their settings, understand when 
                  those settings serve them, and gently expand their range of motion.
                </Typography>
                <Divider sx={{ my: 2 }} />
                <Typography variant="body2" sx={{ fontStyle: 'italic' }}>
                  "The goal isn't to take over the controls. It's to teach them to fly."
                </Typography>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* ── Home-School Bridge ── */}
      <Box sx={{ py: { xs: 6, md: 8 } }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            The Home-School Bridge
          </Typography>
          <Typography
            variant="body1"
            color="text.secondary"
            textAlign="center"
            sx={{ mb: 5, maxWidth: 740, mx: 'auto' }}
          >
            Your child has one nervous system. It doesn't switch off at the school gate. 
            ParentOS connects to PrimaryOS, creating a shared language between home and school.
          </Typography>

          <Grid container spacing={4} alignItems="center">
            <Grid item xs={12} md={5}>
              <Stack spacing={3}>
                {[
                  {
                    icon: <VisibilityIcon />,
                    title: 'See What Teachers See',
                    description: 'Your child\'s teacher may already be observing ETP patterns in the classroom. ParentOS gives you access to the same framework.',
                  },
                  {
                    icon: <ConnectIcon />,
                    title: 'Continuity, Not Fragmentation',
                    description: 'When home and school use the same vocabulary, your child doesn\'t have to maintain two separate selves.',
                  },
                  {
                    icon: <SchoolIcon />,
                    title: 'Informed Conversations',
                    description: 'Parent-teacher meetings shift from "Your child sometimes struggles to focus" to "Your child\'s Voltage Sensitivity seems high during group work — let\'s design a protocol."',
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
            
            <Grid item xs={12} md={7}>
              <Paper
                elevation={3}
                sx={{
                  p: 3,
                  borderRadius: 3,
                  bgcolor: alpha(brandColor, 0.04),
                  border: `1px solid ${alpha(brandColor, 0.2)}`,
                }}
              >
                <Typography variant="h6" gutterBottom fontWeight="bold" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <ChildIcon /> How It Works
                </Typography>
                
                <Grid container spacing={2} sx={{ mt: 1 }}>
                  <Grid item xs={12} sm={6}>
                    <Paper elevation={0} sx={{ p: 2, bgcolor: 'background.paper', height: '100%' }}>
                      <Typography variant="subtitle2" fontWeight="bold" gutterBottom sx={{ color: brandColor }}>
                        📱 PrimaryOS (School)
                      </Typography>
                      <Typography variant="body2">
                        Teachers log observations against the 9 spectra. Patterns emerge across 
                        classroom contexts. The system suggests strategies that work for similar profiles.
                      </Typography>
                    </Paper>
                  </Grid>
                  
                  <Grid item xs={12} sm={6}>
                    <Paper elevation={0} sx={{ p: 2, bgcolor: 'background.paper', height: '100%' }}>
                      <Typography variant="subtitle2" fontWeight="bold" gutterBottom sx={{ color: brandColor }}>
                        🏠 ParentOS (Home)
                      </Typography>
                      <Typography variant="body2">
                        You receive insights, not raw data. "Your child shows Independent social 
                        preferences at school — here are strategies that have worked for similar children."
                      </Typography>
                    </Paper>
                  </Grid>
                  
                  <Grid item xs={12}>
                    <Paper elevation={0} sx={{ p: 2, bgcolor: alpha(brandColor, 0.06), mt: 1 }}>
                      <Typography variant="body2" sx={{ fontStyle: 'italic' }}>
                        <strong>Coming soon:</strong> Parent-requested observations. Flag patterns you see at home 
                        and discover whether teachers observe the same. No more guessing if it's "just at home" 
                        or "just at school."
                      </Typography>
                    </Paper>
                  </Grid>
                </Grid>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* ── Free Foundation Tier ── */}
      <Box sx={{ py: { xs: 6, md: 8 }, bgcolor: alpha('#0f766e', 0.06) }}>
        <Container maxWidth="md">
          <Paper
            elevation={3}
            sx={{
              p: 4,
              borderRadius: 4,
              textAlign: 'center',
              background: `linear-gradient(135deg, white 0%, ${alpha('#f0fdfa', 0.8)} 100%)`,
              border: `1px solid ${alpha('#0f766e', 0.2)}`,
            }}
          >
            <LockOpenIcon sx={{ fontSize: 56, color: '#0f766e', mb: 2 }} />
            <Typography variant="h4" component="h2" fontWeight="bold" gutterBottom>
              Foundation Access — Always Free
            </Typography>
            <Typography variant="body1" sx={{ mb: 3, maxWidth: 550, mx: 'auto' }}>
              Every parent deserves to understand their child's wiring. The core ParentOS framework — 
              the 9 spectra, slider states, and age-stage guidance — will always be free.
            </Typography>
            
            <Grid container spacing={2} sx={{ mb: 3 }}>
              {[
                '9 Spectra Explorer',
                'Observation Question Bank',
                'Age-Stage Roadmaps',
                'Basic Compatibility Tools',
              ].map((feature) => (
                <Grid item xs={12} sm={6} key={feature}>
                  <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 1 }}>
                    <VerifiedIcon sx={{ color: '#0f766e', fontSize: 18 }} />
                    <Typography variant="body2">{feature}</Typography>
                  </Box>
                </Grid>
              ))}
            </Grid>
            
            <Divider sx={{ my: 3 }} />
            
            <Typography variant="body2" color="text.secondary">
              <strong>Premium features</strong> (home-school integration, personalised protocols, 
              detailed progress tracking) are available through school partnerships. 
              ParentOS exists to serve families — not to profit from them.
            </Typography>
          </Paper>
        </Container>
      </Box>

      {/* ── The Origin: 12 Years of Classroom Evidence ── */}
      <Box sx={{ py: { xs: 6, md: 8 } }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Grid container spacing={4} alignItems="center">
            <Grid item xs={12} md={6}>
              <Typography variant="overline" sx={{ color: brandColor, fontWeight: 'bold' }}>
                Built on Evidence
              </Typography>
              <Typography variant="h4" component="h2" fontWeight="bold" gutterBottom>
                From the Classroom to Your Living Room
              </Typography>
              <Typography variant="body1" paragraph sx={{ fontSize: '1.1rem', lineHeight: 1.7 }}>
                ParentOS wasn't designed in a Silicon Valley boardroom. It was forged in 12 years 
                of classroom observation — watching thousands of children reveal their biological 
                wiring through their behaviour.
              </Typography>
              <Typography variant="body1" paragraph sx={{ fontSize: '1.1rem', lineHeight: 1.7 }}>
                The 9 spectra emerged from asking, "Why do two children in the same environment 
                have completely opposite reactions?" The answer wasn't parenting, personality, or 
                defiance. It was biology.
              </Typography>
              <Typography variant="body1" sx={{ fontSize: '1.1rem', lineHeight: 1.7 }}>
                ParentOS brings those 12 years of pattern recognition to your family — not as 
                prescriptive advice, but as a lens for seeing what was always there.
              </Typography>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Card elevation={4} sx={{ borderRadius: 3 }}>
                <CardContent sx={{ p: 4 }}>
                  <Typography variant="h6" gutterBottom fontWeight="bold">
                    Changing the Narrative on Observations:
                  </Typography>
                  <Stack spacing={2}>
                    {[
                      '"The child who couldn\'t sit still wasn\'t defiant — their Risk Seeking slider was maxed out. Once we gave them acceptable risks, they regulated."',
                      '"We kept trying to make an Independent child more Cohesive. Once we stopped, their anxiety dropped and their work improved."',
                      '"The student who absorbed everyone\'s emotions wasn\'t fragile — their Mirror Neurons were just set to Absorbent. They needed filters, not fixing."',
                    ].map((quote, i) => (
                      <Paper
                        key={i}
                        elevation={0}
                        sx={{
                          p: 2,
                          bgcolor: alpha(brandColor, 0.04),
                          borderLeft: `4px solid ${brandColor}`,
                          borderRadius: '0 8px 8px 0',
                        }}
                      >
                        <Typography variant="body2" sx={{ fontStyle: 'italic' }}>
                          "{quote}"
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

      {/* ── What ParentOS Is Not ── */}
      <Box sx={{ py: { xs: 6, md: 8 }, bgcolor: '#f5f7f9' }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            What ParentOS Is — And Isn't
          </Typography>
          
          <Grid container spacing={3} sx={{ mt: 2 }}>
            <Grid item xs={12} md={6}>
              <Paper elevation={2} sx={{ p: 3, height: '100%', bgcolor: 'white' }}>
                <Typography variant="h6" fontWeight="bold" gutterBottom sx={{ color: brandColor }}>
                  ✅ ParentOS Is:
                </Typography>
                <List dense>
                  {[
                    'A lens for observing your child’s patterns',
                    'A vocabulary for describing what you see',
                    'A bridge between home and school',
                    'A framework for understanding compatibility',
                    'Free foundational knowledge for every family',
                    'Built on 12 years of classroom evidence',
                  ].map((item) => (
                    <ListItem key={item} sx={{ px: 0 }}>
                      <ListItemIcon sx={{ minWidth: 36, color: brandColor }}>
                        <VerifiedIcon fontSize="small" />
                      </ListItemIcon>
                      <ListItemText primary={item} primaryTypographyProps={{ variant: 'body2' }} />
                    </ListItem>
                  ))}
                </List>
              </Paper>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Paper elevation={2} sx={{ p: 3, height: '100%', bgcolor: 'white' }}>
                <Typography variant="h6" fontWeight="bold" gutterBottom sx={{ color: '#c62828' }}>
                  ❌ ParentOS Is Not:
                </Typography>
                <List dense>
                  {[
                    'A diagnostic tool',
                    'A replacement for professional assessment',
                    'A checklist of "normal" development',
                    'A way to label or pathologize children',
                    'Prescriptive parenting advice',
                    'A scorecard for comparing children',
                  ].map((item) => (
                    <ListItem key={item} sx={{ px: 0 }}>
                      <ListItemIcon sx={{ minWidth: 36, color: '#c62828' }}>
                        <WarningIcon fontSize="small" />
                      </ListItemIcon>
                      <ListItemText primary={item} primaryTypographyProps={{ variant: 'body2' }} />
                    </ListItem>
                  ))}
                </List>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* ── CTA Section ── */}
      <Box
        sx={{
          py: { xs: 6, md: 8 },
          textAlign: 'center',
          background: `linear-gradient(135deg, #0f766e 0%, #14b8a6 100%)`,
          color: 'white',
        }}
      >
        <Container maxWidth="md">
          <FamilyIcon sx={{ fontSize: 56, mb: 2 }} />
          <Typography variant="h4" fontWeight="bold" gutterBottom>
            Start Observing, Stop Guessing
          </Typography>
          <Typography variant="body1" sx={{ mb: 4, maxWidth: 600, mx: 'auto', opacity: 0.95 }}>
            ParentOS gives you a framework for seeing what's already there. 
            No labels. No judgment. Just a clearer lens on the child you already know.
          </Typography>
          
          <Stack direction="row" spacing={2} justifyContent="center" flexWrap="wrap" useFlexGap>
            <Button
              variant="contained"
              size="large"
              component={Link}
              to="/parent-os/explorer"
              startIcon={<ExploreIcon />}
              sx={{
                px: 5,
                py: 1.5,
                fontWeight: 'bold',
                bgcolor: 'white',
                color: '#0f766e',
                '&:hover': { bgcolor: alpha('#ffffff', 0.9) },
              }}
            >
              Launch Spectrum Explorer
            </Button>
            <ContactReveal email="world@espthinking.co.uk" label="Get in Touch" />
            <Button
              variant="outlined"
              size="large"
              component="a"
              href="/parent-os/guide"
              sx={{
                px: 5,
                py: 1.5,
                fontWeight: 'bold',
                color: 'white',
                borderColor: 'white',
                '&:hover': { borderColor: 'white', bgcolor: alpha('#ffffff', 0.1) },
              }}
            >
              Read the Parent Guide
            </Button>
          </Stack>
          
          <Typography variant="caption" sx={{ display: 'block', mt: 4, opacity: 0.8 }}>
            Foundation access is always free. No credit card required. No data harvesting. 
            Just a better way to see your child.
          </Typography>
        </Container>
      </Box>

      {/* ── Footer ── */}
      <Box sx={{ py: 3, textAlign: 'center', bgcolor: '#0f766e', color: 'rgba(255,255,255,0.8)' }}>
        <Typography variant="body2">
          ParentOS — The Foundational Code Your Family Runs On — Part of the ESP Thinking Portfolio
        </Typography>
        <Typography variant="caption" sx={{ display: 'block', mt: 0.5, opacity: 0.7 }}>
          ParentOS is an observational framework, not a medical device. Always consult qualified professionals for developmental concerns.
        </Typography>
      </Box>
    </Box>
  );
};

export default ParentOSLanding;