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
} from '@mui/material';
import {
  Home as HomeIcon,
  Psychology as PsychologyIcon,
  Tune as TuneIcon,
  Bolt as BoltIcon,
  Groups as GroupsIcon,
  Visibility as VisibilityIcon,
  Security as SecurityIcon,
  Balance as BalanceIcon,
  SelfImprovement as SelfImprovementIcon,
  Favorite as FavoriteIcon,
  Speed as SpeedIcon,
  ElectricalServices as ElectricalIcon,
  Shield as ShieldIcon,
  School as SchoolIcon,
  Hub as HubIcon,
  MenuBook as BookIcon,
  LocalFireDepartment as FireIcon,
  CrisisAlert as AlertIcon,
  Lock as LockIcon,
  LockOpen as LockOpenIcon,
  SwapHoriz as SwapIcon,
  Sensors as SensorsIcon,
  Straighten as StraightenIcon,
} from '@mui/icons-material';
import { Link } from 'react-router-dom';
import ContactReveal from '../components/ContactReveal';

// ─── Spectra Data ───────────────────────────────────────────
interface Spectrum {
  id: string;
  label: string;
  leftLabel: string;
  rightLabel: string;
  description: string;
  icon: React.ReactNode;
  color: string;
}

const spectra: Spectrum[] = [
  {
    id: 'social_gravity',
    label: 'Social Gravity',
    leftLabel: 'Independent',
    rightLabel: 'Cohesive',
    description: 'How social interaction affects your energy. Neither solitude nor company is "wrong" — they\'re different charging modes.',
    icon: <GroupsIcon />,
    color: '#1565c0',
  },
  {
    id: 'energy_directionality',
    label: 'Energy Directionality',
    leftLabel: 'Inward',
    rightLabel: 'Outward',
    description: 'Where your energy naturally flows. Internal processors need time to think; external processors think by talking.',
    icon: <SwapIcon />,
    color: '#6a1b9a',
  },
  {
    id: 'voltage_sensitivity',
    label: 'Voltage Sensitivity',
    leftLabel: 'Insulated',
    rightLabel: 'Conductive',
    description: 'How much emotional current you absorb. Insulated people need thick boundaries; conductive people thrive on emotional exchange.',
    icon: <BoltIcon />,
    color: '#f57c00',
  },
  {
    id: 'threat_response',
    label: 'Threat Response',
    leftLabel: 'Passive',
    rightLabel: 'Aggressive',
    description: 'How you react when threats are detected. Neither freeze-flee nor fight is inherently wrong — dysfunction is being stuck.',
    icon: <ShieldIcon />,
    color: '#c62828',
  },
  {
    id: 'care_response',
    label: 'Care Response',
    leftLabel: 'Detached',
    rightLabel: 'Nurturing',
    description: 'How you respond to vulnerability. Detachment enables objectivity; nurturing enables connection. Both are needed.',
    icon: <FavoriteIcon />,
    color: '#ad1457',
  },
  {
    id: 'risk_tolerance',
    label: 'Risk Tolerance',
    leftLabel: 'Averse',
    rightLabel: 'Seeking',
    description: 'How risk affects your emotional state. Caution prevents catastrophe; risk-seeking drives discovery. Context determines value.',
    icon: <SpeedIcon />,
    color: '#2e7d32',
  },
  {
    id: 'integrity_logic',
    label: 'Integrity Logic',
    leftLabel: 'Relativistic',
    rightLabel: 'Absolutist',
    description: 'How you process moral decisions. Relativism enables empathy; absolutism enables conviction. Dysfunction is being unable to switch.',
    icon: <BalanceIcon />,
    color: '#4527a0',
  },
  {
    id: 'mirror_neuron_tuning',
    label: 'Mirror Neuron Tuning',
    leftLabel: 'Selective',
    rightLabel: 'Absorbent',
    description: 'How much you absorb others\' emotional states. Selective filtering enables focus; absorption enables deep empathy.',
    icon: <SensorsIcon />,
    color: '#00838f',
  },
  {
    id: 'orderliness',
    label: 'Orderliness',
    leftLabel: 'Flexible',
    rightLabel: 'Ordered',
    description: 'Your preference for structure vs spontaneity. Flexibility enables adaptation; order enables efficiency. Both are tools.',
    icon: <StraightenIcon />,
    color: '#5d4037',
  },
];

// ─── Slider State Data ──────────────────────────────────────
const sliderStates = [
  {
    state: 'Fluid',
    icon: <LockOpenIcon />,
    description: 'Moves freely — healthy. Can adapt to context.',
    color: '#2e7d32',
  },
  {
    state: 'Sticky',
    icon: <TuneIcon />,
    description: 'Moves with difficulty — developing. Needs practice to shift.',
    color: '#f57c00',
  },
  {
    state: 'Stuck',
    icon: <LockIcon />,
    description: 'Cannot move — dysfunction. Frozen at one position regardless of context.',
    color: '#c62828',
  },
  {
    state: 'Hijacked',
    icon: <AlertIcon />,
    description: 'Moves involuntarily — triggered. External stimuli move the slider without consent.',
    color: '#6a1b9a',
  },
];

// ─── Component ──────────────────────────────────────────────
const ETPLanding: React.FC = () => {
  const theme = useTheme();

  return (
    <Box>
      {/* ── Hero Section ── */}
      <Box
        sx={{
          background: `linear-gradient(135deg, #4527a0 0%, #6a4bd6 45%, #b69bff 100%)`,
          borderTop: '6px solid #4527a0',
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
                src={`${process.env.PUBLIC_URL}/ESPLogoShort.png`}
                alt="ESP Thinking Portfolio"
                sx={{ height: { xs: 48, md: 64 }, width: 'auto', mb: 2 }}
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
                  <PsychologyIcon sx={{ fontSize: { xs: 48, md: 64 } }} />
                  <Typography variant="h2" component="h1" fontWeight="bold" sx={{ lineHeight: 1.1 }}>
                    HumanOS
                  </Typography>
                </Box>

                <Typography variant="h5" sx={{ opacity: 0.92, maxWidth: 600 }}>
                  A neurophysiologically-grounded emotional intelligence framework that models human behaviour as positioning on biological spectra — enabling morally-neutral intervention design and self-calibration training.
                </Typography>

                <Stack direction="row" spacing={1.5} flexWrap="wrap" useFlexGap>
                  <Chip label="9 Core Spectra" sx={{ bgcolor: 'rgba(255,255,255,0.15)', color: 'white' }} />
                  <Chip label="Morally Neutral" sx={{ bgcolor: 'rgba(255,255,255,0.15)', color: 'white' }} />
                  <Chip label="Trainable Skills" sx={{ bgcolor: 'rgba(255,255,255,0.15)', color: 'white' }} />
                  <Chip label="Classroom-Tested" sx={{ bgcolor: 'rgba(255,255,255,0.15)', color: 'white' }} />
                </Stack>

                <Stack direction="row" spacing={2} sx={{ pt: 1 }}>
                  <Button
                    variant="contained"
                    size="large"
                    component={Link}
                    to="/etp-profile"
                    sx={{
                      bgcolor: 'white',
                      color: '#1a237e',
                      fontWeight: 'bold',
                      px: 4,
                      '&:hover': { bgcolor: alpha('#ffffff', 0.9) },
                    }}
                  >
                    Launch ETP Demo
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
                <Typography variant="h6" color="primary" gutterBottom fontWeight="bold">
                  The Animal–Pilot Model
                </Typography>
                <Typography variant="body2" color="text.secondary" paragraph>
                  Every person has two systems: the <strong>Animal</strong> (the biological hardware — respresented by 9 Emotional Trigger Points (ETP's) sliders
                  that respond automatically to stimuli) and the <strong>Pilot</strong> (executive function — the capacity
                  to <em>choose</em> where to position each slider).
                </Typography>
                <Divider sx={{ my: 1.5 }} />
                <Stack spacing={1.5}>
                  {[
                    { icon: <ElectricalIcon fontSize="small" />, text: 'ETPs are not personality flaws — they are positioning on biological spectra.' },
                    { icon: <BoltIcon fontSize="small" />, text: 'Dysfunction is not having a slider at an extreme — it\'s having it stuck there.' },
                    { icon: <TuneIcon fontSize="small" />, text: 'The goal: agency to move sliders yourself, rather than be trapped by your default settings.' },
                  ].map((item, i) => (
                    <Box key={i} sx={{ display: 'flex', alignItems: 'flex-start', gap: 1.5 }}>
                      <Box sx={{ color: 'primary.main', mt: 0.3 }}>{item.icon}</Box>
                      <Typography variant="body2">{item.text}</Typography>
                    </Box>
                  ))}
                </Stack>
                <Typography variant="body2" sx={{ mt: 2, fontStyle: 'italic', color: 'text.secondary' }}>
                  "You don't get angry at Hydrogen for being explosive — you just learn how to handle it so it doesn't blow up the lab."
                </Typography>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* ── The Neurocomputational Model ── */}
<Box sx={{ py: { xs: 6, md: 8 }, bgcolor: '#f8f9fa' }}>
  <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
    <Grid container spacing={4} alignItems="center">
      <Grid item xs={12} md={6}>
        <Typography variant="h4" component="h2" fontWeight="bold" gutterBottom>
          The Science Behind the Sliders
        </Typography>
        <Typography variant="body1" paragraph sx={{ lineHeight: 1.7 }}>
          HumanOS emerged from a simple observation: every behavioral response follows a predictable computational path:
        </Typography>
        
        <Box sx={{ pl: 2, borderLeft: `3px solid ${theme.palette.primary.main}` }}>
          {[
            "1. **Event occurs** → Sensory input enters the limbic system",
            "2. **Emotional valuation** → Evolutionary instincts, societal norms, and personal history compete to assign 'pain' or 'pleasure' values",
            "3. **Spectra positioning** → We have distilled these inputs into 9 emotional trigger points (ETPs).  These sliders represent the comfort zome of an individual.",
            "4. **Behavior emerges** → The collective slider positions determine what action 'feels right' in that moment",
            "5. **Conscious modulation** → The Pilot (executive function) can override, but requires significant energy"
          ].map((step, i) => (
            <Typography key={i} variant="body1" sx={{ mb: 1.5, '& strong': { color: theme.palette.primary.main } }}>
              <strong>Step {i+1}:</strong> {step.replace(/^\d+\.\s+\*\*(.*?)\*\*\s+/, '')}
            </Typography>
          ))}
        </Box>

        <Typography variant="body1" sx={{ mt: 3, fontStyle: 'italic', color: 'text.secondary' }}>
          "The sliders aren't arbitrary categories — they're the control knobs of your emotional valuation system. 
          Understanding where yours are stuck is understanding why certain situations trigger you."
        </Typography>
      </Grid>
      
      <Grid item xs={12} md={6}>
        <Paper elevation={3} sx={{ p: 3, borderRadius: 2 }}>
          <Typography variant="h6" color="primary" gutterBottom fontWeight="bold">
            From Theory to AI Empathy
          </Typography>
          <Typography variant="body2" paragraph>
            This computational model enables something unprecedented: <strong>AI systems that can genuinely model human emotional logic.</strong>
          </Typography>
          <Stack spacing={2} sx={{ mt: 2 }}>
            {[
              "AI can simulate how different slider positions would value the same event",
              "Predict why two people with different spectra react oppositely",
              "Design compatibility protocols that work with, not against, biological wiring",
              "Move beyond scripted empathy to genuine emotional prediction"
            ].map((item, i) => (
              <Box key={i} sx={{ display: 'flex', gap: 1.5 }}>
                <BoltIcon fontSize="small" sx={{ color: 'primary.main', mt: 0.2 }} />
                <Typography variant="body2">{item}</Typography>
              </Box>
            ))}
          </Stack>
        </Paper>
      </Grid>
    </Grid>
  </Container>
</Box>

      {/* ── Core Philosophy Section ── */}
      <Box sx={{ py: { xs: 6, md: 8 }, bgcolor: '#f5f7f9' }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>

<Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
  The Philosophy: Repurposing Shame, Expanding Agency
</Typography>
<Typography
  variant="body1"
  color="text.secondary"
  textAlign="center"
  sx={{ mb: 5, maxWidth: 780, mx: 'auto' }}
>
  Shame evolved to regulate social behavior — it's the pain of exclusion from the tribe. 
  But in modern contexts, it's often misapplied: we shame people for their <em>biological wiring</em> rather than their <em>choices</em>.
</Typography>
<Typography
  variant="body1"
  textAlign="center"
  sx={{ mb: 5, maxWidth: 700, mx: 'auto', fontStyle: 'italic', color: 'text.secondary' }}
>
  HumanOS redirects shame from "who you are" to "how you're stuck." 
  The goal isn't to eliminate shame, but to harness its motivational power for growth rather than self-hatred.
</Typography>

          <Grid container spacing={3}>
            {[
              {
                icon: <VisibilityIcon fontSize="large" />,
                title: 'Moral Neutrality',
                description:
                  'No spectrum position is "good" or "bad." Being independent is not worse than being cohesive. Being risk-averse is not worse than being risk-seeking. They are different operating settings for different contexts.',
              },
              {
                icon: <TuneIcon fontSize="large" />,
                title: 'Range of Motion, Not "Correct" Position',
                description:
                  'The goal is never to push someone to a "healthy" end of a spectrum. It\'s to expand their range — so they can choose the right position for the right situation, rather than being stuck.',
              },
              {
                icon: <SelfImprovementIcon fontSize="large" />,
                title: 'The Electric Fence Metaphor',
                description:
                  'Your comfort zone has an electric fence around it. Crossing it hurts. But the voltage can be reduced through practice, and seeing what\'s on the other side (the "Freedom Vision") motivates the crossing.',
              },
              {
                icon: <SchoolIcon fontSize="large" />,
                title: 'Options as Education\'s Core Goal',
                description:
                  'Education\'s purpose is maximising available life paths, not creating specific outcomes. Provide keys (capabilities), not destinations. Success = breadth of accessible options.',
              },
              {
                icon: <HubIcon fontSize="large" />,
                title: 'Nature + Nurture + Operating System',
                description:
                  'Nature = hardware (genetics). Nurture = software (learned patterns). HumanOS = the operating system that runs on the hardware and is programmed by experience.',
              },
              {
                icon: <SecurityIcon fontSize="large" />,
                title: 'Neurodiversity as Different Hardware',
                description:
                  'Autism, ADHD, and other neurotypes are different neurological hardware — not deficiencies. Different defaults, different keyrings of capabilities. Diverse paths, not hierarchy of worth.',
              },
            ].map((card) => (
              <Grid item xs={12} md={6} key={card.title}>
                <Card
                  elevation={0}
                  sx={{
                    height: '100%',
                    border: `1px solid ${alpha(theme.palette.primary.main, 0.12)}`,
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
                          bgcolor: alpha(theme.palette.primary.main, 0.08),
                          color: 'primary.main',
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

      {/* ── The 9 Spectra Section ── */}
      <Box sx={{ py: { xs: 6, md: 8 } }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            The 9 Core Spectra
          </Typography>
          <Typography
            variant="body1"
            color="text.secondary"
            textAlign="center"
            sx={{ mb: 2, maxWidth: 740, mx: 'auto' }}
          >
            Each spectrum is a biological slider — innate, morally neutral, and orthogonal to the others.
            Every person sits somewhere on each spectrum, and the position itself is never the problem.
          </Typography>
          <Typography
            variant="body2"
            color="text.secondary"
            textAlign="center"
            sx={{ mb: 5, maxWidth: 600, mx: 'auto', fontStyle: 'italic' }}
          >
            Like a mixing desk in a recording studio — every slider contributes to the overall sound. 
            The art is in the mix.
          </Typography>

          <Grid container spacing={2.5}>
            {spectra.map((s) => (
              <Grid item xs={12} sm={6} md={4} key={s.id}>
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
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 1.5 }}>
                    <Box sx={{ color: s.color }}>{s.icon}</Box>
                    <Typography variant="subtitle1" fontWeight="bold">
                      {s.label}
                    </Typography>
                  </Box>

                  {/* Spectrum visual bar */}
                  <Box sx={{ mb: 1.5 }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
                      <Chip label={s.leftLabel} size="small" variant="outlined" sx={{ fontSize: '0.7rem' }} />
                      <Chip label={s.rightLabel} size="small" variant="outlined" sx={{ fontSize: '0.7rem' }} />
                    </Box>
                    <Box
                      sx={{
                        height: 6,
                        borderRadius: 3,
                        background: `linear-gradient(90deg, ${alpha(s.color, 0.3)}, ${s.color})`,
                      }}
                    />
                  </Box>

                  <Typography variant="body2" color="text.secondary" sx={{ lineHeight: 1.6, flexGrow: 1 }}>
                    {s.description}
                  </Typography>
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
              startIcon={<TuneIcon />}
              sx={{ px: 5, py: 1.5, fontWeight: 'bold', bgcolor: '#1a237e' }}
            >
              Try the Interactive ETP Demo →
            </Button>
          </Box>
        </Container>
      </Box>

      {/* ── Slider States Section ── */}
      <Box sx={{ py: { xs: 6, md: 8 }, bgcolor: '#f5f7f9' }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            Slider State Model
          </Typography>
          <Typography
            variant="body1"
            color="text.secondary"
            textAlign="center"
            sx={{ mb: 5, maxWidth: 700, mx: 'auto' }}
          >
            Dysfunction isn't about <em>where</em> a slider sits — it's about whether you can <em>move</em> it. 
            Each slider can be in one of four states:
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
                    borderLeft: `4px solid ${s.color}`,
                  }}
                >
                  <Box sx={{ color: s.color, mb: 1 }}>{s.icon}</Box>
                  <Typography variant="h6" fontWeight="bold" gutterBottom>
                    {s.state}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {s.description}
                  </Typography>
                </Paper>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* ── Global Moderators Section ── */}
      <Box sx={{ py: { xs: 6, md: 8 } }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            Global Moderators
          </Typography>
          <Typography
            variant="body1"
            color="text.secondary"
            textAlign="center"
            sx={{ mb: 5, maxWidth: 700, mx: 'auto' }}
          >
            Two meta-controls that affect <em>all</em> 9 spectra simultaneously. They are the "master volume"
            and "power supply" of the entire system.
          </Typography>

          <Grid container spacing={4} justifyContent="center">
            <Grid item xs={12} md={5}>
              <Paper elevation={3} sx={{ p: 4, borderRadius: 2, height: '100%' }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                  <PsychologyIcon sx={{ fontSize: 40, color: '#1565c0' }} />
                  <Typography variant="h5" fontWeight="bold">Pilot Strength</Typography>
                </Box>
                <Typography variant="body2" color="text.secondary" paragraph sx={{ lineHeight: 1.7 }}>
                  Executive function capacity — the hand on ALL sliders. When pilot strength is <strong>high</strong>,
                  you move sliders deliberately through conscious choice. When <strong>low</strong>, triggers
                  control your positioning and the sliders move reactively.
                </Typography>
                <Divider sx={{ my: 2 }} />
                <Typography variant="body2" sx={{ fontStyle: 'italic', color: 'text.secondary' }}>
                  Think of it as the difference between driving a car with your hands on the wheel vs. 
                  being in the back seat while something else steers.
                </Typography>
              </Paper>
            </Grid>

            <Grid item xs={12} md={5}>
              <Paper elevation={3} sx={{ p: 4, borderRadius: 2, height: '100%' }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                  <FireIcon sx={{ fontSize: 40, color: '#c62828' }} />
                  <Typography variant="h5" fontWeight="bold">Current Load</Typography>
                </Box>
                <Typography variant="body2" color="text.secondary" paragraph sx={{ lineHeight: 1.7 }}>
                  Stress and depletion level — narrows range of motion on <em>all</em> spectra. When load is 
                  <strong> low</strong>, the full range of motion is available. When <strong>high</strong>,
                  range contracts toward your comfort zone and flexibility drops.
                </Typography>
                <Divider sx={{ my: 2 }} />
                <Typography variant="body2" sx={{ fontStyle: 'italic', color: 'text.secondary' }}>
                  A well-rested person can handle more social interaction, take more risks, and tolerate 
                  more ambiguity. A depleted person retreats to default settings.
                </Typography>
              </Paper>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* ── Voltage Compatibility Section ── */}
      <Box
        sx={{
          py: { xs: 6, md: 8 },
          bgcolor: alpha('#1a237e', 0.03),
          borderTop: `1px solid ${alpha('#1a237e', 0.08)}`,
        }}
      >
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            Compatibility, Not Correction
          </Typography>
          <Typography
            variant="body1"
            color="text.secondary"
            textAlign="center"
            sx={{ mb: 5, maxWidth: 780, mx: 'auto' }}
          >
            When two people have different settings, the traditional approach is to "fix" one of them.
            HumanOS designs <strong>compatibility protocols</strong> — systems that let different settings 
            work together productively.
          </Typography>

          <Grid container spacing={2}>
            {[
              { spectrum: 'Social Gravity', solution: 'Social Voltage Zones', bridge: 'Paired work with clear recharge breaks', color: '#1565c0' },
              { spectrum: 'Energy Direction', solution: 'Energy Circuit Design', bridge: '"You talk, I listen, then switch"', color: '#6a1b9a' },
              { spectrum: 'Voltage Sensitivity', solution: 'Step-Down Transformers', bridge: 'Insulated observes, Conductive participates, then discuss', color: '#f57c00' },
              { spectrum: 'Threat Response', solution: 'Threat Calibration', bridge: '"When to stand ground, when to withdraw — both are tactics"', color: '#c62828' },
              { spectrum: 'Care Response', solution: 'Care Boundaries', bridge: '"Put your own oxygen mask on first"', color: '#ad1457' },
              { spectrum: 'Risk Tolerance', solution: 'Risk Gradients', bridge: 'Seeker scouts, Averse maps safety', color: '#2e7d32' },
              { spectrum: 'Integrity Logic', solution: 'Integrity Framing', bridge: '"These 3 things are absolute, the rest is relative"', color: '#4527a0' },
              { spectrum: 'Mirror Neurons', solution: 'Empathy Filters', bridge: '"I\'m resonating with your X" / "I see you\'re feeling X"', color: '#00838f' },
              { spectrum: 'Orderliness', solution: 'Structure Gradients', bridge: '"Here\'s the framework, you choose the order within it"', color: '#5d4037' },
            ].map((item) => (
              <Grid item xs={12} sm={6} md={4} key={item.spectrum}>
                <Box
                  sx={{
                    p: 2,
                    borderLeft: `3px solid ${item.color}`,
                    bgcolor: alpha(item.color, 0.04),
                    borderRadius: '0 8px 8px 0',
                    height: '100%',
                  }}
                >
                  <Typography variant="subtitle2" fontWeight="bold" sx={{ color: item.color }}>
                    {item.spectrum}
                  </Typography>
                  <Typography variant="body2" fontWeight="medium" gutterBottom>
                    {item.solution}
                  </Typography>
                  <Typography variant="caption" color="text.secondary" sx={{ fontStyle: 'italic' }}>
                    {item.bridge}
                  </Typography>
                </Box>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* ── Classroom Evidence Section ── */}
      <Box sx={{ py: { xs: 6, md: 8 } }}>
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            12+ Years of Classroom Evidence
          </Typography>
          <Typography
            variant="body1"
            color="text.secondary"
            textAlign="center"
            sx={{ mb: 5, maxWidth: 700, mx: 'auto' }}
          >
            HumanOS isn't academic theory — it was forged in real classrooms, observing real students, 
            solving real problems that existing frameworks couldn't explain.
          </Typography>

          <Grid container spacing={3}>
            {[
              {
                title: 'Voltage Reduction',
                description: 'Start with familiar ground. Recap previous content. Connect to known material. Lower the "electric fence voltage" before introducing new territory.',
                icon: <BoltIcon />,
              },
              {
                title: 'The Behavioral Pipeline',
                description: 'Stimuli → Processing Mechanisms → Emotional Response → ETP Threshold → Behaviour. Understanding the pipeline lets you intervene at the right stage.',
                icon: <HubIcon />,
              },
              {
                title: 'Neuro-Chemical Circuitry',
                description: 'Reward (Dopamine), Praise (Dopamine + Oxytocin), Risk (Adrenaline + Cortisol), Humility (GABA + Serotonin). Every classroom moment has a chemical signature.',
                icon: <ElectricalIcon />,
              },
              {
                title: 'Group Harmonics',
                description: 'Constructive vs destructive interference. High-voltage "alpha nodes" as transmitters. The teacher as "heat sink" — the ground wire in the classroom circuit.',
                icon: <GroupsIcon />,
              },
              {
                title: 'The "Sticker Book" Curriculum',
                description: 'Year 3: "I have different settings." Year 4: "I can change my setting." Year 5: "Different situations need different settings." Year 6: "Can I access all parts?"',
                icon: <BookIcon />,
              },
              {
                title: 'AI Mentorship Scaffold',
                description: 'Phase 1: High guidance → Phase 2: Gradual transfer → Phase 3: Student leads → Phase 4: Full independence. The fading scaffold approach.',
                icon: <SchoolIcon />,
              },
            ].map((item) => (
              <Grid item xs={12} sm={6} md={4} key={item.title}>
                <Card elevation={1} sx={{ height: '100%' }}>
                  <CardContent sx={{ p: 3 }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 1.5 }}>
                      <Box sx={{ color: 'primary.main' }}>{item.icon}</Box>
                      <Typography variant="subtitle1" fontWeight="bold">{item.title}</Typography>
                    </Box>
                    <Typography variant="body2" color="text.secondary" sx={{ lineHeight: 1.7 }}>
                      {item.description}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* ── What Makes This Different ── */}
      <Box
        sx={{
          py: { xs: 6, md: 8 },
          background: `linear-gradient(135deg, ${alpha('#1a237e', 0.06)} 0%, ${alpha('#4527a0', 0.06)} 100%)`,
        }}
      >
        <Container maxWidth={false} sx={{ maxWidth: '1400px' }}>
          <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
            What Makes ETP Different
          </Typography>

          <Grid container spacing={2} sx={{ mt: 2 }}>
            {[
              'Moral neutrality — no position is "wrong"',
              'Dynamic sliders, not static personality labels',
              'Concrete, trainable skills per spectrum',
              'Compatibility protocols, not moral correction',
              'Explicit circuit breakers for operating modes',
              'Identity hijacking detection built in',
              'Grounded in 12+ years of classroom evidence',
              'Federated, privacy-first learning architecture',
              'The Periodic Table analogy — tangible, no judgment',
            ].map((point) => (
              <Grid item xs={12} sm={6} md={4} key={point}>
                <Box sx={{ display: 'flex', gap: 1.5, p: 1.5 }}>
                  <Box sx={{ color: '#4527a0', mt: 0.25 }}>✓</Box>
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
        }}
      >
        <Container maxWidth="md">
          <PsychologyIcon sx={{ fontSize: 56, color: '#1a237e', mb: 2 }} />
          <Typography variant="h4" fontWeight="bold" gutterBottom>
            Explore the ETP Framework
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 4, maxWidth: 600, mx: 'auto' }}>
            The interactive demo lets you position yourself on all 9 spectra and 2 global moderators. 
            See how the sliders interact to create your unique behavioural profile.
          </Typography>
          <Button
            variant="contained"
            size="large"
            component={Link}
            to="/etp-profile"
            startIcon={<TuneIcon />}
            sx={{
              px: 5,
              py: 1.5,
              fontWeight: 'bold',
              bgcolor: '#1a237e',
              '&:hover': { bgcolor: '#283593' },
            }}
          >
            Launch ETP Demo →
          </Button>
          <ContactReveal email="world@espthinking.co.uk" label="Request Demo" />
        </Container>
      </Box>

      {/* ── Footer ── */}
      <Box sx={{ py: 3, textAlign: 'center', bgcolor: '#4527a0', color: 'rgba(255,255,255,0.85)' }}>
        <Typography variant="body2">
          HumanOS — Emergent Tendency Profiles — Part of the ESP Thinking Portfolio
        </Typography>
      </Box>
    </Box>
  );
};

export default ETPLanding;
