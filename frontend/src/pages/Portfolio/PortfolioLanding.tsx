import React from 'react';
import {
  Box,
  Container,
  Typography,
  Grid,
  Card,
  CardContent,
  CardActions,
  Button,
  Chip,
  Paper,
  Avatar,
  Stack,
  useTheme,
  alpha,
  Divider,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  LinearProgress,
  Menu,
  MenuItem,
} from '@mui/material';
import {
  School as SchoolIcon,
  WorkOff as WorkIcon,
  Balance as BalanceIcon,
  Psychology as PsychologyIcon,
  Analytics as AnalyticsIcon,
  AccountBalance as TrustIcon,
  Badge as BadgeIcon,
  PhoneIphone as MobileIcon,
  Storage as DataIcon,
  TrendingUp as GrowthIcon,
  AccountTree as SkillsIcon,
  Science as ScienceIcon,
  Groups as GroupsIcon,
  Verified as VerifiedIcon,
  SettingsAccessibility as PersonalizationIcon,
  Architecture as ArchitectureIcon,
  IntegrationInstructions as IntegrationIcon,
  Security as SecurityIcon,
  EmojiObjects as InsightIcon,
  ArrowDropDown as DropdownIcon,
Gavel as InterrogateIcon,
SettingsEthernet as DialecticIcon,
AdsClick as CoreIcon,
DataObject as CriteriaIcon,
} from '@mui/icons-material';
import { Link } from 'react-router-dom';
import DivideIcon from '../../components/icons/DivideIcon';
import ContactReveal from '../../components/ContactReveal';
import SEO from '../../components/SEO';

// Portfolio navigation bar
const PortfolioNav: React.FC = () => {
  const [toolsAnchor, setToolsAnchor] = React.useState<null | HTMLElement>(null);
  const openTools = (e: React.MouseEvent<HTMLElement>) => setToolsAnchor(e.currentTarget);
  const closeTools = () => setToolsAnchor(null);

  return (
    <Box
      sx={{
        position: 'relative',
        zIndex: 10,
        py: 2,
        px: 3,
        bgcolor: 'rgba(0,0,0,0.05)',
      }}
    >
      <Container maxWidth="lg">
        <Stack direction="row" justifyContent="space-between" alignItems="center">
          <Typography
            variant="h6"
            sx={{ color: 'white', fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: 1 }}
          >
            <Box component="img" src={`${process.env.PUBLIC_URL}/ESPLogoLong.png`} alt="ESP Thinking Portfolio" sx={{ height: { xs: 64, sm: 90, md: 120 }, width: 'auto' }} />
          </Typography>
          <Stack direction="row" spacing={1} alignItems="center">
            <Button
              component={Link}
              to="/etp-profile"
              size="small"
              sx={{ color: 'white', textTransform: 'none' }}
            >
              ETP Demo
            </Button>
            <Button
              component={Link}
              to="/chisg"
              size="small"
              sx={{ color: 'white', textTransform: 'none' }}
            >
              CHISG
            </Button>
            <Button
              component={Link}
              to="/skillstree"
              size="small"
              sx={{ color: 'white', textTransform: 'none' }}
            >
              Skills Map
            </Button>

            {/* Tools dropdown showing Active Projects */}
            <>
              <Button
                size="small"
                onClick={openTools}
                endIcon={<DropdownIcon />}
                sx={{ color: 'white', textTransform: 'none' }}
              >
                Tools
              </Button>
              <Menu anchorEl={toolsAnchor} open={Boolean(toolsAnchor)} onClose={closeTools}>
                <MenuItem component="a" href="/drb/" onClick={closeTools}>MAT Strategic Dashboard</MenuItem>
                <MenuItem component={Link} to="/lao" onClick={closeTools}>LAO Adaptive Revision</MenuItem>
                <MenuItem component={Link} to="/parent-os" onClick={closeTools}>ParentOS</MenuItem>
                <MenuItem component={Link} to="/primary-os" onClick={closeTools}>PrimaryOS (Neuron Navigators)</MenuItem>
                <MenuItem component={Link} to="/careeros" onClick={closeTools}>CareerOS</MenuItem>
                <MenuItem component={Link} to="/esp-world" onClick={closeTools}>ESP World</MenuItem>
              </Menu>
            </>

          </Stack>
        </Stack>
      </Container>
    </Box>
  );
};

// Hero section with clear positioning
const HeroSection: React.FC = () => {
  const theme = useTheme();
  
  return (
    <Box
      sx={{
        position: 'relative',
        background: `linear-gradient(135deg, ${theme.palette.primary.dark} 0%, ${theme.palette.primary.main} 50%, ${theme.palette.secondary.main} 100%)`,
        color: 'white',
        py: { xs: 6, md: 10 },
        px: 2,
      }}
    >
      <PortfolioNav />
      <Container maxWidth="lg">
        <Grid container spacing={4} alignItems="center">
          <Grid item xs={12} md={8}>
            <Typography variant="overline" sx={{ letterSpacing: 1.5, opacity: 0.9, mt: { xs: 2, md: 0 }, display: 'block' }}>
              PROJECTS PORTFOLIO
            </Typography>
            <Typography
              variant="h2"
              component="h1"
              fontWeight="bold"
              gutterBottom
              sx={{ fontSize: { xs: '1.6rem', sm: '2.25rem', md: '3rem' } }}
            >
              The Architecture of Emotional Readiness.
            </Typography>
            <Typography variant="h5" sx={{ opacity: 0.95, mb: 3 }}>
              A blocked mind cannot learn, so build systems that remove the barriers and unlock agency and potential.
            </Typography>
            
            <Stack direction="row" spacing={2} alignItems="center">
              <Box
                sx={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: 1,
                  px: 2,
                  py: 1,
                  bgcolor: alpha('#ffffff', 0.1),
                  borderRadius: 2,
                }}
              >
                <SchoolIcon />
                <Typography variant="body1">Classroom Reality</Typography>
              </Box>
              <Box sx={{ color: 'white' }}>→</Box>
              <Box
                sx={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: 1,
                  px: 2,
                  py: 1,
                  bgcolor: alpha('#ffffff', 0.1),
                  borderRadius: 2,
                }}
              >
                <ArchitectureIcon />
                <Typography variant="body1">Technical Architecture</Typography>
              </Box>
            </Stack>
          </Grid>
          
          <Grid item xs={12} md={4}>
            <Paper
              elevation={6}
              sx={{
                p: 3,
                borderRadius: 3,
                background: alpha(theme.palette.background.paper, 0.95),
              }}
            >
              <Typography variant="h6" color="primary" gutterBottom fontWeight="bold">
                The Unique Combination
              </Typography>
              <List dense>
                {[
                  { icon: <SchoolIcon />, text: '12+ years classroom experience' },
                  { icon: <ScienceIcon />, text: 'Scientific approach to learning biology' },
                  { icon: <ArchitectureIcon />, text: 'Full-stack technical architecture' },
                  { icon: <VerifiedIcon />, text: 'Focus on verifiable, grounded AI' },
                ].map((item, i) => (
                  <ListItem key={i} sx={{ px: 0 }}>
                    <ListItemIcon sx={{ minWidth: 40, color: 'primary.main' }}>
                      {item.icon}
                    </ListItemIcon>
                    <ListItemText primary={item.text} primaryTypographyProps={{ variant: 'body2' }} />
                  </ListItem>
                ))}
              </List>
            </Paper>
          </Grid>
        </Grid>
      </Container>
    </Box>
  );
};

// --- COMPONENT: FIRST PRINCIPLES / ARCHITECTURE SECTION ---
const UnderlyingArchitectureSection: React.FC = () => {
  const theme = useTheme();

  const methods = [
    {
      title: 'Recursive Interrogation',
      subtitle: 'The "Why" Loop',
      icon: <InterrogateIcon fontSize="large" />,
      color: '#2563eb', // Blue
      description:
        'I employ a Recursive Root-Cause Analysis to separate "Functional Requirements" (what is needed) from "Legacy Artifacts" (how it used to be done). I decouple the human intent from the legacy implementation.',
      chip: 'Root Cause Analysis',
    },
    {
      title: 'Dialectical Engineering',
      subtitle: 'The Hegelian Engine',
      icon: <DialecticIcon fontSize="large" />,
      color: '#7c3aed', // Purple
      description:
        'I resolve diametrically opposed requirements (e.g., Safety vs. Agency) not by compromise, but by architectural synthesis. I use orthogonal engineering to satisfy both constraints on different layers of the stack.',
      chip: 'Synthesis over Compromise',
    },
    {
      title: 'Biological Abstraction',
      subtitle: 'The System Core',
      icon: <CoreIcon fontSize="large" />,
      color: '#059669', // Green
      description:
        'I view educational problems as Information Theory problems wrapped in Biological Constraints. "Emotional Trigger Points" are not fuzzy feelings; they are system interrupts that require specific error-handling protocols.',
      chip: 'Bio-Informatics',
    },
  ];

  return (
    <Box sx={{ py: 10, bgcolor: alpha('#64748b', 0.03) }}>
      <Container maxWidth="lg">
        {/* Section Header */}
        <Box sx={{ mb: 8, textAlign: 'center' }}>
          <Typography variant="overline" sx={{ letterSpacing: 2, color: 'text.secondary', fontWeight: 'bold' }}>
            Operational Philosophy
          </Typography>
          <Typography 
            variant="h3" 
            component="h2" 
            fontWeight="800" 
            gutterBottom
            sx={{ 
              background: `linear-gradient(45deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
            }}
          >
            Recursive Deconstruction
          </Typography>
          <Typography variant="h5" color="text.secondary" sx={{ maxWidth: 800, mx: 'auto', lineHeight: 1.6 }}>
            I don't just build features; I interrogate systems until their fundamental truths are exposed. 
            My process is defined by a refusal to accept "legacy" as a requirement.
          </Typography>
        </Box>

        {/* The 3 Pillars */}
        <Grid container spacing={4}>
          {methods.map((method) => (
            <Grid item xs={12} md={4} key={method.title}>
              <Card 
                elevation={0}
                sx={{ 
                  height: '100%', 
                  bgcolor: 'transparent', // Transparent to blend with background
                  border: `1px solid ${alpha(method.color, 0.2)}`,
                  background: `linear-gradient(180deg, ${alpha(method.color, 0.05)} 0%, ${alpha('#ffffff', 0.4)} 100%)`,
                  transition: 'transform 0.2s',
                  '&:hover': {
                    transform: 'translateY(-5px)',
                    boxShadow: `0 10px 30px -10px ${alpha(method.color, 0.3)}`
                  }
                }}
              >
                <CardContent sx={{ p: 4 }}>
                  <Box 
                    sx={{ 
                      display: 'inline-flex', 
                      p: 2, 
                      borderRadius: 3, 
                      bgcolor: alpha(method.color, 0.1),
                      color: method.color,
                      mb: 3
                    }}
                  >
                    {method.icon}
                  </Box>
                  
                  <Typography variant="h5" fontWeight="bold" gutterBottom>
                    {method.title}
                  </Typography>
                  <Typography variant="subtitle2" sx={{ color: method.color, mb: 2, fontWeight: 'bold' }}>
                    {method.subtitle}
                  </Typography>
                  
                  <Divider sx={{ my: 2, borderColor: alpha(method.color, 0.2) }} />
                  
                  <Typography variant="body1" color="text.secondary" paragraph sx={{ minHeight: 80 }}>
                    {method.description}
                  </Typography>
                  
                  <Chip 
                    label={method.chip} 
                    size="small" 
                    sx={{ 
                      bgcolor: alpha(method.color, 0.1), 
                      color: method.color,
                      fontWeight: 'bold',
                      mt: 1
                    }} 
                  />
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>

        {/* The "Refusal to Let Go" Block */}
        <Paper 
          elevation={3}
          sx={{
            mt: 8,
            p: { xs: 3, md: 5 },
            borderRadius: 4,
            bgcolor: '#1e293b', // Dark Slate
            color: 'white',
            position: 'relative',
            overflow: 'hidden'
          }}
        >
          {/* Decorative background element */}
          <Box sx={{
            position: 'absolute',
            top: -50,
            right: -50,
            width: 200,
            height: 200,
            borderRadius: '50%',
            bgcolor: alpha(theme.palette.primary.main, 0.2),
            zIndex: 0
          }} />

          <Grid container spacing={4} alignItems="center" sx={{ position: 'relative', zIndex: 1 }}>
            <Grid item xs={12} md={8}>
              <Stack direction="row" spacing={2} alignItems="center" sx={{ mb: 2 }}>
                <CriteriaIcon sx={{ color: theme.palette.secondary.light }} />
                <Typography variant="overline" sx={{ color: theme.palette.secondary.light, letterSpacing: 1 }}>
                  Rigorous Termination Criteria
                </Typography>
              </Stack>
              <Typography variant="h4" fontWeight="bold" gutterBottom>
                The Refusal to Let Go
              </Typography>
              <Typography variant="body1" sx={{ opacity: 0.8, fontSize: '1.1rem', lineHeight: 1.8 }}>
                I operate with high-fidelity "Definitions of Done." I do not accept a solution that works "most of the time." 
                If the deep requirement is not met, the architecture is not finished. My process is iterative, aggressive, 
                and grounded in the belief that if a system is broken, it is because we haven't asked the right question yet.
              </Typography>
            </Grid>
            <Grid item xs={12} md={4} sx={{ display: 'flex', justifyContent: 'center' }}>
               {/* Contextual Visual Placeholder -  */}
               <Box 
                sx={{ 
                  p: 3, 
                  border: `2px dashed ${alpha('#ffffff', 0.2)}`, 
                  borderRadius: 4,
                  textAlign: 'center' 
                }}
              >
                <Typography variant="h2" fontWeight="900" sx={{ color: alpha('#ffffff', 0.1) }}>
                  100%
                </Typography>
                <Typography variant="caption" sx={{ color: alpha('#ffffff', 0.6) }}>
                  Fidelity to Core Requirement
                </Typography>
              </Box>
            </Grid>
          </Grid>
        </Paper>
      </Container>
    </Box>
  );
};


// Add this component after HeroSection, before ProblemFrameworkSection

const EndorsementSection: React.FC = () => {
  const theme = useTheme();
  
  return (
    <Box sx={{ 
      py: 4, 
      bgcolor: alpha('#232f3e', 0.02), // AWS brand-adjacent
      borderBottom: `1px solid ${alpha('#232f3e', 0.1)}`,
      borderTop: `1px solid ${alpha('#232f3e', 0.1)}`,
    }}>
      <Container maxWidth="lg">
        <Grid container spacing={3} alignItems="center">
          <Grid item xs={12} md={3} sx={{ textAlign: { xs: 'center', md: 'left' } }}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
  <Avatar sx={{ bgcolor: '#232f3e', width: 40, height: 40 }}>TF</Avatar>
  <Box>
    <Typography variant="subtitle2" fontWeight="bold">Toby Fotherby</Typography>
    <Typography variant="caption" color="text.secondary">
      Senior AI/ML Specialist Solutions Architect
    </Typography>
    <Typography variant="caption" color="text.secondary" display="block">
      AWS, Strategic Accounts
    </Typography>
  </Box>
</Box>  
          </Grid>
          
          <Grid item xs={12} md={9}>
            <Box sx={{ position: 'relative' }}>
              <Typography 
                variant="h4" 
                sx={{ 
                  position: 'absolute', 
                  left: -20, 
                  top: -20, 
                  color: alpha('#232f3e', 0.1),
                  fontSize: '6rem',
                  fontFamily: 'Georgia, serif',
                  userSelect: 'none',
                }}
              >
                "
              </Typography>
              <Typography 
                variant="body1" 
                sx={{ 
                  position: 'relative', 
                  zIndex: 1,
                  fontStyle: 'italic',
                  fontSize: '1.1rem',
                  lineHeight: 1.7,
                  color: 'text.primary',
                  fontWeight: 'medium',
                }}
              >
                Michael quickly stood out as one of the most thoughtful and capable practitioners that I've worked with. 
                He demonstrated mastery of all key concepts—from model evaluation and orchestration to scalability and 
                responsible AI principles.
              </Typography>
              <Typography 
                variant="body1" 
                sx={{ 
                  mt: 2,
                  fontStyle: 'italic',
                  fontSize: '1.1rem',
                  lineHeight: 1.7,
                  color: 'text.primary',
                }}
              >
                I am deeply impressed with the innovative AI-powered teaching strategy and framework he developed. 
                It's an outstanding example of how AI can be applied not only to solve technical problems, but also 
                to transform the way people learn.
              </Typography>
              <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2 }}>
                <Typography variant="subtitle2" sx={{ fontWeight: 'bold' }}>
                  — Toby Fotherby
                </Typography>
                <Typography variant="caption" sx={{ ml: 1, color: 'text.secondary' }}>
                  AWS Senior AI/ML Specialist Solutions Architect
                </Typography>
              </Box>
            </Box>
          </Grid>
        </Grid>
      </Container>
    </Box>
  );
};

// Problem Framework Section
const ProblemFrameworkSection: React.FC = () => {
  const theme = useTheme();
  
  const problems = [
    {
      title: 'Cognitive Atrophy',
      icon: <PsychologyIcon />,
      symptom: '"Short-circuit" thinkers',
      rootCause: 'Information retrieval replacing deep learning',
      solution: 'CHISG - Grounded semantic knowledge structures',
      color: '#4f46e5',
    },
    {
      title: 'Job Obsolescence',
      icon: <WorkIcon />,
      symptom: '1996 curriculum in 2026 economy',
      rootCause: 'Static skill frameworks',
      solution: 'Dynamic competency mapping & XLA relationships',
      color: '#059669',
    },
    {
      title: 'Truth Decay',
      icon: <BalanceIcon />,
      symptom: 'AI hallucinations & "slop"',
      rootCause: 'Ungrounded vector embeddings',
      solution: 'Deterministic validation layers',
      color: '#dc2626',
    },
    {
      title: 'Algorithmic Polarization',
      icon: <GroupsIcon />,
      symptom: 'Coded outrage',
      rootCause: 'One-size-fits-all content delivery',
      solution: 'ETP-informed personalization',
      color: '#d97706',
    },
  ];

  return (
    <Box sx={{ py: 8, bgcolor: 'grey.50' }}>
      <Container maxWidth="lg">
        <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
          The Four Failures of Current Educational AI
        </Typography>
        <Typography variant="body1" textAlign="center" sx={{ mb: 4, maxWidth: 800, mx: 'auto' }}>
          Most educational tech addresses symptoms. I architect solutions for the root causes.
        </Typography>

        <Grid container spacing={3}>
          {problems.map((problem) => (
            <Grid item xs={12} sm={6} key={problem.title}>
              <Card
                elevation={2}
                sx={{
                  height: '100%',
                  borderTop: `4px solid ${problem.color}`,
                }}
              >
                <CardContent>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                    <Box
                      sx={{
                        p: 1.5,
                        borderRadius: 2,
                        bgcolor: alpha(problem.color, 0.1),
                        color: problem.color,
                      }}
                    >
                      {problem.icon}
                    </Box>
                    <Typography variant="h6" fontWeight="bold">
                      {problem.title}
                    </Typography>
                  </Box>

                  <Box sx={{ mb: 2 }}>
                    <Typography variant="caption" color="text.secondary" display="block">
                      Symptom
                    </Typography>
                    <Typography variant="body2" sx={{ fontStyle: 'italic' }}>
                      {problem.symptom}
                    </Typography>
                  </Box>

                  <Box sx={{ mb: 2 }}>
                    <Typography variant="caption" color="text.secondary" display="block">
                      Root Cause
                    </Typography>
                    <Typography variant="body2">
                      {problem.rootCause}
                    </Typography>
                  </Box>

                  <Box>
                    <Typography variant="caption" color="text.secondary" display="block">
                      My Solution
                    </Typography>
                    <Typography variant="body2" fontWeight="medium" color="primary">
                      {problem.solution}
                    </Typography>
                  </Box>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      </Container>
    </Box>
  );
};

// Innovation Stack Visualization
const InnovationStackSection: React.FC = () => {
  const theme = useTheme();
  
  return (
    <Box sx={{ py: 8 }}>
      <Container maxWidth="lg">
        <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
          The Innovation Stack
        </Typography>
        <Typography variant="body1" textAlign="center" sx={{ mb: 4, maxWidth: 700, mx: 'auto' }}>
          Three core innovations that combine to solve educational AI failures
        </Typography>

        <Grid container spacing={4}>
          <Grid item xs={12} md={4}>
            <Paper
              elevation={3}
              sx={{
                p: 3,
                height: '100%',
                border: `2px solid ${alpha('#4f46e5', 0.2)}`,
              }}
            >
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                <VerifiedIcon sx={{ fontSize: 40, color: '#4f46e5' }} />
                <Typography variant="h5" fontWeight="bold">
                  CHISG
                </Typography>
              </Box>
              <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
                Contextualised Hierarchical Iterative Semantic Groupings
              </Typography>
              <Typography variant="body2" paragraph>
                A method for parsing and validating information for vector database ingestion. Creates grounded knowledge structures that eliminate AI hallucinations.
              </Typography>
              <Chip label="Solves: Truth Decay" size="small" sx={{ bgcolor: alpha('#4f46e5', 0.1), color: '#4f46e5' }} />

              <CardActions sx={{ px: 0, pb: 0, pt: 2 }}>
                <Button component={Link} to="/chisg" variant="outlined" size="small">
                  VIEW Demo
                </Button>
              </CardActions>
            </Paper>
          </Grid>

          <Grid item xs={12} md={4}>
            <Paper
              elevation={3}
              sx={{
                p: 3,
                height: '100%',
                border: `2px solid ${alpha('#059669', 0.2)}`,
              }}
            >
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                <PersonalizationIcon sx={{ fontSize: 40, color: '#059669' }} />
                <Typography variant="h5" fontWeight="bold">
                  ETP's and the HumanOS
                </Typography>
              </Box>
              <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
                Emotional Trigger Points
              </Typography>
              <Typography variant="body2" paragraph>
                9 biological spectrums that describe learner tendencies. Enables personalization based on individual cognitive and emotional patterns.
              </Typography>
              <Chip label="Solves: Algorithmic Polarization" size="small" sx={{ bgcolor: alpha('#059669', 0.1), color: '#059669' }} />
            <CardActions sx={{ px: 2, pb: 2 }}>
            <Button component={Link} to={'/etp-landing'} variant="outlined" size="small">
              View Demo
            </Button>
        </CardActions>
         </Paper>
          </Grid>

          <Grid item xs={12} md={4}>
            <Paper
              elevation={3}
              sx={{
                p: 3,
                height: '100%',
                border: `2px solid ${alpha('#dc2626', 0.2)}`,
              }}
            >
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                <SkillsIcon sx={{ fontSize: 40, color: '#dc2626' }} />
                <Typography variant="h5" fontWeight="bold">
                  Skills Map
                </Typography>
              </Box>
              <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
                Developmental Vector Database
              </Typography>
              <Typography variant="body2" paragraph>
                5,000+ indexed objects including skills, semantic links, and documentation. Powers intelligent search and curriculum analysis across the ecosystem.
              </Typography>
              <Chip label="Solves: Job Obsolescence" size="small" sx={{ bgcolor: alpha('#dc2626', 0.1), color: '#dc2626' }} />
             <CardActions sx={{ px: 2, pb: 2 }}>
            <Button component={Link} to={'/skillstree'} variant="outlined" size="small">
              View Demo
            </Button>
        </CardActions>
            </Paper>
          </Grid>
        </Grid>
      </Container>
    </Box>
  );
};

// Projects section (rich project cards with brand colours and demo links)
interface ProjectCardProps {
  title: string;
  subtitle: string;
  description: string;
  tech: string[];
  status: 'Active' | 'In Development' | 'TestFlight' | 'Planned' | string;
  statusColor: 'success' | 'warning' | 'info' | 'default';
  icon: React.ReactNode;
  demoLink?: string;
  isExternal?: boolean;
  relevance?: string;
  impact?: string;
  brandColor?: string;
}

const ProjectCard: React.FC<ProjectCardProps> = ({
  title,
  subtitle,
  description,
  tech,
  status,
  statusColor,
  icon,
  demoLink,
  isExternal,
  relevance,
  impact,
  brandColor,
}) => {
  const theme = useTheme();

  return (
    <Card
      elevation={2}
      sx={{
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        transition: 'transform 0.2s, box-shadow 0.2s',
        borderTop: brandColor ? `4px solid ${brandColor}` : undefined,
        '&:hover': {
          transform: 'translateY(-4px)',
          boxShadow: theme.shadows[8],
        },
      }}
    >
      <CardContent sx={{ flexGrow: 1 }}>
        <Box sx={{ display: 'flex', alignItems: 'flex-start', gap: 2, mb: 2 }}>
          <Box
            sx={{
              p: 1.5,
              borderRadius: 2,
              bgcolor: alpha(brandColor || theme.palette.primary.main, 0.1),
              color: brandColor || 'primary.main',
            }}
          >
            {icon}
          </Box>
          <Box sx={{ flexGrow: 1 }}>
            <Typography variant="h6" component="h3">
              {title}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {subtitle}
            </Typography>
          </Box>
          <Chip label={status} color={statusColor as any} size="small" />
        </Box>

        <Typography variant="body2" paragraph>
          {description}
        </Typography>

        {relevance && (
          <Box sx={{ mb: 2 }}>
            <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
              Data Management Relevance:
            </Typography>
            <Typography variant="body2" sx={{ fontStyle: 'italic', color: 'primary.main' }}>
              {relevance}
            </Typography>
          </Box>
        )}

        {impact && (
          <Box sx={{ mb: 2 }}>
            <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
              Impact:
            </Typography>
            <Typography variant="body2" sx={{ fontStyle: 'italic', color: 'primary.main' }}>
              {impact}
            </Typography>
          </Box>
        )}

        <Stack direction="row" spacing={0.5} flexWrap="wrap" useFlexGap>
          {tech.map((t) => (
            <Chip key={t} label={t} size="small" variant="outlined" sx={{ fontSize: '0.7rem' }} />
          ))}
        </Stack>
      </CardContent>

      {demoLink && (
        <CardActions sx={{ px: 2, pb: 2 }}>
          {isExternal ? (
            <Button component="a" href={demoLink} target="_blank" rel="noreferrer" variant="outlined" size="small">
              View Demo
            </Button>
          ) : (
            <Button component={Link} to={demoLink} variant="outlined" size="small">
              View Demo
            </Button>
          )}
        </CardActions>
      )}
    </Card>
  );
};

const projects: ProjectCardProps[] = [
  {
    title: 'MAT Strategic Dashboard',
    subtitle: 'Trust-Level Governance & Finance',
    description:
      'A high-level oversight engine for Multi-Academy Trusts. Aggregates Open Banking data and utility spend to correlate financial health with physical and educational assets.',
    tech: ['Go', 'React', 'Vite', 'MySQL', 'MongoDB', 'Open Banking'],
    status: 'Concept Demo',
    statusColor: 'info',
    icon: <TrustIcon fontSize="large" />,
    demoLink: '/drb/',
    isExternal: true,
    relevance: 'Secure financial data aggregation and multi-site governance reporting.',
    impact: 'Transforms "Structural Displacement" into fiscal and operational clarity for MAT leadership.',
    brandColor: '#374151',
  },
  {
    title: 'LAO Adaptive Revision',
    subtitle: 'Cognitive Friction Reduction',
    description:
      'GCSE Science revision tool utilizing CHISG-driven gap analysis. Features speed-reading modes and audio summaries to match the specific "Voltage Sensitivity" of the learner.',
    tech: ['React Native', 'Go', 'MongoDB', 'AWS Polly'],
    status: 'TestFlight',
    statusColor: 'success',
    icon: <SchoolIcon fontSize="large" />,
    demoLink: '/lao',
    isExternal: false,
    relevance: 'Production-grade AI integration focusing on adaptive learning and cognitive load management.',
    impact: 'Drives student self-motivation by removing structural barriers to information retrieval.',
    brandColor: '#1a237e',
  },
  {
    title: 'ParentOS',
    subtitle: 'The Foundational Code Your Family Runs On.',
    description:
      'A parent-facing portal that demystifies educational data. Translates complex ETP and CHISG data into human narratives, providing personalized guidance for home support.',
    tech: ['React', 'Weaviate', 'Go API', 'MySQL'],
    status: 'Planned',
    statusColor: 'info',
    icon: <SchoolIcon fontSize="large" />,
    demoLink: '/parent-os',
    isExternal: false,
    relevance: 'Data democratization and multi-audience semantic explanation.',
    impact: 'Bridges the gap between school and home by making "Black Box" educational data transparent.',
    brandColor: '#0f766e',
  },
  {
    title: 'PrimaryOS and the Neuron Navigators Guide Book',
    subtitle: 'Real-Time ETP Sensor Array',
    description:
      'A "point-of-observation" tool for teachers to log student competency and Emotional Trigger Points (ETP) in real-time. Captures the "biology of the classroom" as it happens.',
    tech: ['React Native', 'Expo', 'Weaviate', 'Go API'],
    status: 'Concept',
    statusColor: 'info',
    icon: <MobileIcon fontSize="large" />,
    demoLink: '/primary-os',
    isExternal: false,
    relevance: 'High-concurrency data capture and real-time synchronization with vector databases.',
    impact: 'Reduces teacher administrative load through instant, tap-based competency logging.',
    brandColor: '#f59e0b',
  },
  {
    title: 'CareerOS',
    subtitle: 'Verified Skills Passport + Intelligent Matching',
    description:
      'A post-CV skills infrastructure: role profiles, portable verified skills passports, and semantic matching that measures distance-to-competency instead of keyword overlap.',
    tech: ['React/TypeScript', 'CHISG', 'Weaviate', 'Go API'],
    status: 'Planned',
    statusColor: 'info',
    icon: <BadgeIcon fontSize="large" />,
    demoLink: '/careeros',
    isExternal: false,
    relevance: 'Semantic gap analysis and evidence-backed capability verification across education → employment.',
    impact: 'Replaces keyword matching with defensible, auditable competency distance scoring.',
    brandColor: '#2563eb',
  },
  {
    title: 'ESP World',
    subtitle: 'The Complete AI Enhanced Ecosystem',
    description:
      'A comprehensive School Management Information System (MIS) serving as the integration layer for student profiles, group dynamics, and strategic visualization. It acts as the "Central Command" for grounded educational data.',
    tech: ['Go', 'React', 'MySQL', 'Docker', 'Microsoft Graph', 'Weaviate', 'AWS', 'Nginx'],
    status: 'In Development',
    statusColor: 'warning',
    icon: <AnalyticsIcon fontSize="large" />,
    demoLink: '/esp-world',
    isExternal: false,
    relevance: 'Architecting the full data lifecycle: moving from raw capture to high-fidelity strategic reporting.',
    impact: 'Consolidates fragmented school data into a unified, actionable knowledge base.',
    brandColor: '#4527a0',
  },
];

const ProjectsSection: React.FC = () => {
  return (
    <Box sx={{ py: 8, bgcolor: 'grey.50' }}>
      <Container maxWidth="lg">
        <Typography variant="h4" component="h2" gutterBottom textAlign="center">
          Active Projects
        </Typography>
        <Typography variant="body1" color="text.secondary" textAlign="center" sx={{ mb: 4, maxWidth: 700, mx: 'auto' }}>
          An interconnected ecosystem of educational technology tools, demonstrating data architecture, 
          system integration, and domain expertise in the education sector.
        </Typography>
        
        <Grid container spacing={3}>
          {projects.map((project) => (
            <Grid item xs={12} md={6} key={project.title}>
              <ProjectCard {...project} />
            </Grid>
          ))}
        </Grid>
      </Container>
    </Box>
  );
};

// Technical Capabilities Section (Condensed)
const TechnicalCapabilitiesSection: React.FC = () => {
  const capabilities = {
    'Data Architecture': ['PostgreSQL', 'MySQL', 'MongoDB', 'Weaviate (Vector DB)', 'CHISG Methodology'],
    'Backend & APIs': ['Go', 'RESTful/GraphQL APIs', 'Microservices', 'Docker'],
    'Frontend & Mobile': ['React/TypeScript', 'React Native/Expo', 'Material UI', 'D3.js'],
    'Education Domain': ['12+ Years Classroom', 'Multi-Academy Trusts', 'Curriculum Design', 'Assessment Systems'],
  };

  return (
    <Box sx={{ py: 8 }}>
      <Container maxWidth="lg">
        <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
          Technical & Domain Capabilities
        </Typography>

        <Grid container spacing={3} sx={{ mt: 2 }}>
          {Object.entries(capabilities).map(([category, items]) => (
            <Grid item xs={12} sm={6} md={3} key={category}>
              <Typography variant="subtitle1" fontWeight="bold" gutterBottom color="primary">
                {category}
              </Typography>
              <Stack spacing={0.5}>
                {items.map((item) => (
                  <Box
                    key={item}
                    sx={{
                      px: 1.5,
                      py: 0.75,
                      bgcolor: 'grey.50',
                      borderRadius: 1,
                      border: '1px solid',
                      borderColor: 'divider',
                    }}
                  >
                    <Typography variant="body2">{item}</Typography>
                  </Box>
                ))}
              </Stack>
            </Grid>
          ))}
        </Grid>
      </Container>
    </Box>
  );
};

// Design Principles Section
const DesignPrinciplesSection: React.FC = () => {
  const principles = [
    {
      title: 'Ground Truth',
      description: 'No claim without provenance. Every skill, link, and recommendation must be traceable to evidence or explicitly marked as hypothesis.',
      icon: <VerifiedIcon />,
      color: '#2563eb',
    },
    {
      title: 'Biological First',
      description: 'Behaviour is the output of biological systems. Design for the nervous system, not the curriculum.',
      icon: <PsychologyIcon />,
      color: '#7c3aed',
    },
    {
      title: 'Replace Shame, Not Settings',
      description: 'No child is "bad." Mismatch and stuck sliders are engineering problems, not moral ones.',
      icon: <BalanceIcon />,
      color: '#059669',
    },
    {
      title: 'Honest Uncertainty',
      description: '"I don\'t know" is always a valid answer. Confidence scores accompany every inference.',
      icon: <ScienceIcon />,
      color: '#dc2626',
    },
    {
      title: 'Portable Identity',
      description: 'Skills should travel with the learner — from classroom to career, institution to institution.',
      icon: <BadgeIcon />,
      color: '#f59e0b',
    },
    {
      title: 'Compatibility, Not Correction',
      description: 'When two systems conflict, design a protocol — don\'t declare one system broken.',
      icon: <GroupsIcon />,
      color: '#0f766e',
    },
  ];

  return (
    <Box sx={{ py: 8, bgcolor: '#ffffff' }}>
      <Container maxWidth="lg">
        <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
          Design Principles
        </Typography>
        <Typography variant="body1" textAlign="center" sx={{ mb: 5, maxWidth: 700, mx: 'auto' }}>
          These aren't features. They're commitments — to learners, to teachers, and to honest engineering.
        </Typography>

        <Grid container spacing={3}>
          {principles.map((principle) => (
            <Grid item xs={12} md={6} key={principle.title}>
              <Paper
                elevation={0}
                sx={{
                  p: 3,
                  height: '100%',
                  border: `1px solid ${alpha(principle.color, 0.2)}`,
                  borderRadius: 2,
                  transition: 'border-color 0.2s',
                  '&:hover': {
                    borderColor: principle.color,
                  },
                }}
              >
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 1.5 }}>
                  <Box sx={{ color: principle.color }}>{principle.icon}</Box>
                  <Typography variant="h6" fontWeight="bold">
                    {principle.title}
                  </Typography>
                </Box>
                <Typography variant="body2" color="text.secondary">
                  {principle.description}
                </Typography>
              </Paper>
            </Grid>
          ))}
        </Grid>
      </Container>
    </Box>
  );
};

// Current Focus Section
const CurrentFocusSection: React.FC = () => {
  const focusAreas = [
    {
      area: 'CHISG Validation Pipeline',
      description: 'Moving from internal benchmarks to structured external evaluation',
      milestone: 'Q2 2024',
      progress: 65,
      color: '#2563eb',
    },
    {
      area: 'ETP Classroom Pilots',
      description: '3 primary schools, 12 teachers, 240 students',
      milestone: 'Data collection complete',
      progress: 90,
      color: '#7c3aed',
    },
    {
      area: 'Skills Map Expansion',
      description: 'Adding cross-curricular semantic links (target: +400)',
      milestone: 'Q3 2024',
      progress: 45,
      color: '#059669',
    },
    {
      area: 'ParentOS Beta',
      description: 'Recruiting 50 parent testers from pilot schools',
      milestone: 'June 2024',
      progress: 30,
      color: '#0f766e',
    },
  ];

  return (
    <Box sx={{ py: 8, bgcolor: '#f8fafc' }}>
      <Container maxWidth="lg">
        <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
          Current Focus
        </Typography>
        <Typography variant="body1" textAlign="center" sx={{ mb: 4, maxWidth: 700, mx: 'auto' }}>
          What I'm building right now. No smoke, no mirrors.
        </Typography>

        <Grid container spacing={3}>
          {focusAreas.map((area) => (
            <Grid item xs={12} md={6} key={area.area}>
              <Paper elevation={2} sx={{ p: 3 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 1.5 }}>
                  <Typography variant="subtitle1" fontWeight="bold">
                    {area.area}
                  </Typography>
                  <Chip label={area.milestone} size="small" sx={{ bgcolor: alpha(area.color, 0.1), color: area.color }} />
                </Box>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                  {area.description}
                </Typography>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <LinearProgress
                    variant="determinate"
                    value={area.progress}
                    sx={{
                      flexGrow: 1,
                      height: 8,
                      borderRadius: 4,
                      bgcolor: alpha(area.color, 0.1),
                      '& .MuiLinearProgress-bar': {
                        bgcolor: area.color,
                      },
                    }}
                  />
                  <Typography variant="caption" fontWeight="bold">
                    {area.progress}%
                  </Typography>
                </Box>
              </Paper>
            </Grid>
          ))}
        </Grid>
      </Container>
    </Box>
  );
};

// Impact Metrics Section
const ImpactSection: React.FC = () => {
  const metrics = [
    { value: '70%', label: 'Reduction in AI hallucinations (internal tests)', icon: <SecurityIcon /> },
    { value: '579', label: 'CHISG-mapped skills with semantic validation', icon: <VerifiedIcon /> },
    { value: '1,120+', label: 'Semantic relationships preventing "slop"', icon: <SkillsIcon /> },
    { value: '30%', label: 'Reduction in teacher admin load (projected)', icon: <SchoolIcon /> },
    { value: '9', label: 'Biological learning dimensions (ETP spectra)', icon: <ScienceIcon /> },
    { value: '12+', label: 'Years of classroom experience informing design', icon: <InsightIcon /> },
  ];

  return (
    <Box sx={{ py: 8, bgcolor: 'primary.light', background: 'linear-gradient(135deg, #e0f2fe 0%, #f0f9ff 100%)' }}>
      <Container maxWidth="lg">
        <Typography variant="h4" component="h2" textAlign="center" fontWeight="bold" gutterBottom>
          Measurable Impact
        </Typography>
        <Typography variant="body1" textAlign="center" sx={{ mb: 4, maxWidth: 700, mx: 'auto' }}>
          Combining innovation with practicality to deliver tangible results
        </Typography>

        <Grid container spacing={3}>
          {metrics.map((metric) => (
            <Grid item xs={12} sm={6} md={4} key={metric.label}>
              <Paper
                elevation={1}
                sx={{
                  p: 3,
                  height: '100%',
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  textAlign: 'center',
                }}
              >
                <Box sx={{ color: 'primary.main', mb: 1 }}>
                  {metric.icon}
                </Box>
                <Typography variant="h3" fontWeight="bold" gutterBottom>
                  {metric.value}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {metric.label}
                </Typography>
              </Paper>
            </Grid>
          ))}
        </Grid>
      </Container>
    </Box>
  );
};

// CTA Section
const CTASection: React.FC = () => {
  const theme = useTheme();
  
  return (
    <Box
      sx={{
        py: 8,
        background: `linear-gradient(135deg, ${theme.palette.primary.dark} 0%, ${theme.palette.primary.main} 100%)`,
        color: 'white',
      }}
    >
      <Container maxWidth="md" sx={{ textAlign: 'center' }}>
        <Typography variant="h4" gutterBottom fontWeight="bold">
          Ready to Build Trustworthy Educational AI?
        </Typography>
        <Typography variant="body1" sx={{ opacity: 0.9, mb: 4, maxWidth: 600, mx: 'auto' }}>
          Let's discuss how grounded semantic architectures and biologically-informed personalization can transform your educational technology.
        </Typography>
        
        <Stack direction="row" spacing={2} justifyContent="center" flexWrap="wrap" useFlexGap>
          <ContactReveal email="world@espthinking.co.uk" label="Start a Conversation" />
          <Button
            variant="outlined"
            size="large"
            sx={{ color: 'white', borderColor: 'white', fontWeight: 'bold' }}
            component={Link}
            to="/skillstree"
          >
            Explore the Skills Map
          </Button>
          <Button
            variant="outlined"
            size="large"
            sx={{ color: 'white', borderColor: 'white', fontWeight: 'bold' }}
            component={Link}
            to="/etp-profile"
          >
            Experience ETP Demo
          </Button>
        </Stack>
        
        <Typography variant="body2" sx={{ mt: 4, opacity: 0.8, fontStyle: 'italic' }}>
          "Bridging classroom reality with technical architecture to make AI trustworthy in education."
        </Typography>
      </Container>
    </Box>
  );
};

// Main portfolio landing page
const PortfolioLanding: React.FC = () => {
  return (
    <Box>
      <SEO
        title="ESP Thinking — Educational Technology That Understands Learners"
        description="Portfolio of educational technology tools built on CHISG knowledge graphs and ETP biological spectra. Skills mapping, adaptive revision, AI coaching, and strategic oversight for schools and trusts."
        path="/"
      />
      <HeroSection />
      <EndorsementSection />
      <ProblemFrameworkSection />
      <InnovationStackSection />
      <DesignPrinciplesSection />
      <UnderlyingArchitectureSection/>
      <ProjectsSection />
      <CurrentFocusSection />
      <TechnicalCapabilitiesSection />
      {/* <ImpactSection /> */}
      <CTASection />
    </Box>
  );
};

export default PortfolioLanding;