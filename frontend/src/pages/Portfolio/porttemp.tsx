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
  Stack,
  useTheme,
  alpha,
} from '@mui/material';
import {
  School as SchoolIcon,
  WorkOff as WorkIcon,
  Balance as BalanceIcon,
  Psychology as PsychologyIcon,
  Analytics as AnalyticsIcon,
  AccountBalance as TrustIcon,
  PhoneIphone as MobileIcon,
  Storage as DataIcon,
  TrendingUp as GrowthIcon,
  AccountTree as SkillsIcon,
} from '@mui/icons-material';
import { Link } from 'react-router-dom';
import DivideIcon from '../../components/icons/DivideIcon';

// Portfolio navigation bar
const PortfolioNav: React.FC = () => {
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
          <Typography
            variant="h6"
            sx={{ color: 'white', fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: 1 }}
          >
              <Box component="img" src={`${process.env.PUBLIC_URL}/ESPLogoLong.png`} alt="ESP Thinking Portfolio" sx={{ height: 90, width: 'auto' }} />
          </Typography>
          <Stack direction="row" spacing={1}>
            <Button
              component={Link}
              to="/tools"
              size="small"
              sx={{ color: 'white', textTransform: 'none' }}
            >
              ETP Demo
            </Button>
            <Button
              component="a"
              href="http://localhost/skillstree/"
              size="small"
              sx={{ color: 'white', textTransform: 'none' }}
            >
              Skills Map
            </Button>
            <Button
              component={Link}
              to="/tools"
              size="small"
              sx={{ color: 'white', textTransform: 'none' }}
            >
              Tools
            </Button>
          </Stack>
        </Stack>
      </Container>
    </Box>
  );
};

// Hero section with key value prop
const HeroSection: React.FC = () => {
  const theme = useTheme();
  
  return (
    <Box
      sx={{
        position: 'relative',
        background: `linear-gradient(135deg, ${theme.palette.primary.dark} 0%, ${theme.palette.primary.main} 50%, ${theme.palette.secondary.main} 100%)`,
        color: 'white',
        py: { xs: 6, md: 10 },
        pt: { xs: 10, md: 14 }, // Extra top padding for nav
        px: 2,
      }}
    >
      <PortfolioNav />
      <Container maxWidth="lg">
        <Grid container spacing={4} alignItems="center">
          <Grid item xs={12} md={7}>
            <Typography variant="h2" component="h1" fontWeight="bold" gutterBottom>
              Redefining the Boundaries of AI and Human Development.
            </Typography>
            <Typography variant="h5" sx={{ opacity: 0.95, mb: 3 }}>
              Human-centred AI for education — reducing AI hallucinations, explaining skills with CHISG, and personalising learning via ETPs
            </Typography>
            <Stack direction="row" spacing={2} flexWrap="wrap" useFlexGap> 
              <Chip label="Data Architecture" color="secondary" />
              <Chip label="EdTech" color="secondary" />
              <Chip label="Multi-Academy Trusts" color="secondary" />
              <Chip label="AI/ML Integration" color="secondary" />
            </Stack>
          </Grid>
          <Grid item xs={12} md={5}>
            <Paper
              elevation={6}
              sx={{
                p: 3,
                borderRadius: 3,
                background: alpha(theme.palette.background.paper, 0.95),
              }}
            >
              <Typography variant="h6" color="primary" gutterBottom>
                The Four "Horsemen" of Educational AI
              </Typography>
              <Stack spacing={1}>
                {[
                  { icon: <PsychologyIcon />, text: 'Cognitive Atrophy: Are we becoming "Short-circuit" thinkers?' },
                  { icon: <WorkIcon />, text: 'Job Obsolescence: Are we training children for a 1995 economy?' },
                  { icon: <BalanceIcon />, text: 'The Truth Decay: Is "AI Slop" and hallucination making reality a choice rather than a fact?' },
                  { icon: <DivideIcon />, text: 'Algorithmic Polarization: Are we coding outrage into the next generation?' },
                ].map((item, i) => (
                  <Box key={i} sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                    <Box sx={{ color: 'primary.main' }}>{item.icon}</Box>
                    <Typography variant="body1">{item.text}</Typography>
                  </Box>
                ))}
              </Stack>
            </Paper>
          </Grid>
        </Grid>
      </Container>
    </Box>
  );
};

// Project card component
interface ProjectCardProps {
  title: string;
  subtitle: string;
  description: string;
  tech: string[];
  status: string;
  statusColor: 'success' | 'warning' | 'info';
  icon: React.ReactNode;
  demoLink?: string;
  isExternal?: boolean; // true = navigates via nginx (full page load), false = React Router
  relevance: string;
  impact?: string;
  brandColor?: string; // visual identity colour for this app
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
          <Chip label={status} color={statusColor} size="small" />
        </Box>
        
        <Typography variant="body2" paragraph>
          {description}
        </Typography>
        
        <Box sx={{ mb: 2 }}>
          <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
            Data Management Relevance:
          </Typography>
          <Typography variant="body2" sx={{ fontStyle: 'italic', color: 'primary.main' }}>
            {relevance}
          </Typography>
          </Box>
          <Box sx={{ mb: 2 }}>

          <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
            Impact:
          </Typography>
          <Typography variant="body2" sx={{ fontStyle: 'italic', color: 'primary.main' }}>
            {impact}
          </Typography>
        </Box>
        
        <Stack direction="row" spacing={0.5} flexWrap="wrap" useFlexGap>
          {tech.map((t) => (
            <Chip key={t} label={t} size="small" variant="outlined" sx={{ fontSize: '0.7rem' }} />
          ))}
        </Stack>
      </CardContent>
      
      {demoLink && (
        <CardActions sx={{ px: 2, pb: 2 }}>
          <Button
            {...(isExternal
              ? { component: 'a', href: demoLink }
              : { component: Link, to: demoLink })}
            variant="outlined"
            size="small"
          >
            View Demo
          </Button>
        </CardActions>
      )}
    </Card>
  );
};


// Main projects section
const projects: ProjectCardProps[] = [
  {
    title: 'MAT Strategic Dashboard',
    subtitle: 'Trust-Level Governance & Finance',
    description: 'A high-level oversight engine for Multi-Academy Trusts. Aggregates Open Banking data and utility spend to correlate financial health with physical and educational assets.',
    tech: ['Go', 'React', 'Vite', 'MySQL', 'MongoDB', 'Open Banking'],
    status: 'In Development',
    statusColor: 'warning',
    icon: <TrustIcon fontSize="large" />,
    demoLink: 'http://localhost/drb/app',
    isExternal: true,
    relevance: 'Secure financial data aggregation and multi-site governance reporting.',
    impact: 'Transforms "Structural Displacement" into fiscal and operational clarity for MAT leadership.',
    brandColor: '#374151',
  },
  {
    title: 'CHISG Skills Map',
    subtitle: 'Deterministic Knowledge Architecture',
    description: 'An interactive D3.js visualization engine for hierarchical skill structures. Uses a semantic knowledge graph to map competency relationships, enabling intelligent gap analysis and pathfinding.',
    tech: ['React', 'D3.js', 'Go', 'MySQL'],
    status: 'Active',
    statusColor: 'success',
    icon: <SkillsIcon fontSize="large" />,
    demoLink: 'http://localhost/skillstree/',
    isExternal: true,
    relevance: 'Solving "Truth Decay" through hierarchical data structures and relationship discovery.',
    impact: 'Maps 579+ skills and 1,120+ semantic links to eliminate AI hallucination in learning pathways.',
    brandColor: '#15803d',
  },
  {
    title: 'PrimaryOS and the Neuron Navigators Guide Book',
    subtitle: 'Real-Time ETP Sensor Array',
    description: 'A "point-of-observation" tool for teachers to log student competency and Emotional Trigger Points (ETP) in real-time. Captures the "biology of the classroom" as it happens.',
    tech: ['React Native', 'Expo', 'Weaviate', 'Go API'],
    status: 'In Development',
    statusColor: 'warning',
    icon: <MobileIcon fontSize="large" />,
    // demoLink: 'http://localhost/skillsmarkbook/', // Not yet deployed
    isExternal: true,
    relevance: 'High-concurrency data capture and real-time synchronization with vector databases.',
    impact: 'Reduces teacher administrative load by 30% through instant, tap-based competency logging.',
    brandColor: '#f59e0b',
  },
  {
    title: 'LAO Adaptive Revision',
    subtitle: 'Cognitive Friction Reduction',
    description: 'GCSE Science revision tool utilizing CHISG-driven gap analysis. Features speed-reading modes and audio summaries to match the specific "Voltage Sensitivity" of the learner.',
    tech: ['React Native', 'Go', 'MongoDB', 'AWS Polly'],
    status: 'TestFlight',
    statusColor: 'success',
    icon: <SchoolIcon fontSize="large" />,
    demoLink: 'http://localhost/lao',
    isExternal: true,
    relevance: 'Production-grade AI integration focusing on adaptive learning and cognitive load management.',
    impact: 'Drives student self-motivation by removing structural barriers to information retrieval.',
    brandColor: '#1a237e',
  },
  {
    title: 'ParentOS',
    subtitle: 'The Foundational Code Your Family Runs On.',
    description: 'A parent-facing portal that demystifies educational data. Translates complex ETP and CHISG data into human narratives, providing personalized guidance for home support.',
    tech: ['React', 'Weaviate', 'Go API', 'MySQL'],
    status: 'Planned',
    statusColor: 'info',
    icon: <SchoolIcon fontSize="large" />,
    demoLink: 'http://localhost/esp-organizer/skills-profile-parent',
    isExternal: true,
    relevance: 'Data democratization and multi-audience semantic explanation.',
    impact: 'Bridges the gap between school and home by making "Black Box" educational data transparent.',
    brandColor: '#0f766e',
  },
  {
    title: 'ESP World',
    subtitle: 'The Complete AI Enhanced Ecosystem',
    description: 'A comprehensive School Management Information System (MIS) serving as the integration layer for student profiles, group dynamics, and strategic visualization. It acts as the "Central Command" for grounded educational data.',
    tech: ['Go', 'React', 'MySQL', 'Docker', 'Microsoft Graph', 'Weaviate', 'AWS', 'Nginx'],
    status: 'In Development',
    statusColor: 'warning',
    icon: <AnalyticsIcon fontSize="large" />,
    demoLink: 'http://localhost/esp-organizer/etp-landing',
    isExternal: true,
    relevance: 'Architecting the full data lifecycle: moving from raw capture to high-fidelity strategic reporting.',
    impact: 'Consolidates fragmented school data into a unified, actionable knowledge base.',
    brandColor: '#4527a0',
  },
];

// Projects section
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

// Core innovations section
const InnovationsSection: React.FC = () => {
  const theme = useTheme();
  
  const innovations = [
    {
      title: 'CHISG',
      subtitle: 'Contextualised Hierarchical Iterative Semantic Groupings',
      description: 'A method for parsing and validating information for vector database ingestion. Enables intelligent gap analysis, learning pathway optimisation, and reduces AI hallucinations through grounded knowledge structures.',
      icon: <PsychologyIcon fontSize="large" />,
    },
    {
      title: 'ETP',
      subtitle: 'Emergent Tendency Profiles',
      description: '9 biological spectrums that describe learner tendencies (social gravity, energy directionality, voltage sensitivity, etc.). Enables personalised learning approaches based on individual cognitive and emotional patterns.',
      icon: <GrowthIcon fontSize="large" />,
    },
    {
      title: 'Skills Map',
      subtitle: 'Developmental Vector Database',
      description: 'Over 5,000 indexed objects including skills, semantic links, and documentation. Powers intelligent search, relationship discovery, and curriculum analysis across the entire ecosystem.',
      icon: <SkillsIcon fontSize="large" />,
    },
  ];

  return (
    <Box sx={{ py: 8 }}>
      <Container maxWidth="lg">
        <Typography variant="h4" component="h2" gutterBottom textAlign="center">
          Core Innovations
        </Typography>
        <Typography variant="body1" color="text.secondary" textAlign="center" sx={{ mb: 4 }}>
          Novel approaches to educational data that differentiate this ecosystem
        </Typography>
        
        <Grid container spacing={4}>
          {innovations.map((item) => (
            <Grid item xs={12} md={4} key={item.title}>
              <Paper
                elevation={0}
                sx={{
                  p: 3,
                  height: '100%',
                  border: `1px solid ${theme.palette.divider}`,
                  borderRadius: 2,
                }}
              >
                <Box sx={{ color: 'primary.main', mb: 2 }}>{item.icon}</Box>
                <Typography variant="h6" gutterBottom>
                  {item.title}
                </Typography>
                <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
                  {item.subtitle}
                </Typography>
                <Typography variant="body2">
                  {item.description}
                </Typography>
              </Paper>
            </Grid>
          ))}
        </Grid>
      </Container>
    </Box>
  );
};

// Impact & Outcomes section
const ImpactSection: React.FC = () => {
  const theme = useTheme();

  return (
    <Box sx={{ py: 8, background: `linear-gradient(90deg, ${theme.palette.secondary.light}, ${theme.palette.primary.light})` }}>
      <Container maxWidth="lg">
        <Typography variant="h4" component="h2" gutterBottom textAlign="center">
          Impact & Outcomes
        </Typography>
        <Box sx={{ display: 'flex', justifyContent: 'center', mb: 2 }}>
          <DivideIcon width={36} height={36} color={theme.palette.primary.main} />
        </Box>
        <Typography variant="body1" color="text.secondary" textAlign="center" sx={{ mb: 4, maxWidth: 800, mx: 'auto' }}>
          Driving progress in human–AI interaction by combining CHISG's grounded semantic knowledge with ETP-informed personalisation to make AI safer and more useful in education.
        </Typography>

        <Grid container spacing={3}>
          {[
            { title: '579', label: 'CHISG Skills' },
            { title: '1,120', label: 'Semantic Links' },
            { title: '300+', label: 'XLA Relationships' },
            { title: 'Internal Evaluations', label: 'Hallucination Reduction' },
          ].map((item) => (
            <Grid item xs={12} sm={6} md={3} key={item.label}>
              <Paper elevation={3} sx={{ p: 3, textAlign: 'center' }}>
                <Typography variant="h5" fontWeight="bold">{item.title}</Typography>
                <Typography variant="caption" color="text.secondary">{item.label}</Typography>
              </Paper>
            </Grid>
          ))}
        </Grid>

        <Box sx={{ mt: 3 }}>
          <Typography variant="body2" color="text.secondary" textAlign="center">
            <strong>Outcomes:</strong> Improved interpretability, actionable guidance for parents & teachers, and reduced hallucination risk in model-in-the-loop scenarios.
          </Typography>
        </Box>
      </Container>
    </Box>
  );
};

// Technical capabilities section
const CapabilitiesSection: React.FC = () => {
  const capabilities = {
    'Data & Databases': ['PostgreSQL', 'MySQL', 'MongoDB', 'Weaviate (Vector DB)', 'SQLite'],
    'Backend & APIs': ['Go', 'RESTful APIs', 'GraphQL', 'Docker', 'Microservices'],
    'Frontend & Mobile': ['React', 'TypeScript', 'React Native', 'Material UI', 'Expo'],
    'Cloud & DevOps': ['AWS (S3, Polly)', 'Vultr VPS', 'Nginx', 'CI/CD', 'Linux'],
    'Domain Expertise': ['Education Sector', 'Multi-Academy Trusts', 'Open Banking', 'Knowledge Graphs'],
  };

  return (
    <Box sx={{ py: 8, bgcolor: 'grey.50' }}>
      <Container maxWidth="lg">
        <Typography variant="h4" component="h2" gutterBottom textAlign="center">
          Technical Capabilities
        </Typography>
        
        <Grid container spacing={3} sx={{ mt: 2 }}>
          {Object.entries(capabilities).map(([category, items]) => (
            <Grid item xs={12} sm={6} md={4} key={category}>
              <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
                {category}
              </Typography>
              <Stack direction="row" spacing={0.5} flexWrap="wrap" useFlexGap>
                {items.map((item) => (
                  <Chip key={item} label={item} size="small" sx={{ mb: 0.5 }} />
                ))}
              </Stack>
            </Grid>
          ))}
        </Grid>
      </Container>
    </Box>
  );
};

// Contact/CTA section
const CTASection: React.FC = () => {
  const theme = useTheme();
  
  return (
    <Box
      sx={{
        py: 6,
        background: `linear-gradient(135deg, ${theme.palette.primary.main} 0%, ${theme.palette.primary.dark} 100%)`,
        color: 'white',
      }}
    >
      <Container maxWidth="md" sx={{ textAlign: 'center' }}>
        <Typography variant="h5" gutterBottom>
          Interested in discussing data management solutions?
        </Typography>
        <Typography variant="body1" sx={{ opacity: 0.9, mb: 3 }}>
          I bring a unique combination of scientific rigour, classroom experience, and technical capability to education data challenges.
        </Typography>
        <Stack direction="row" spacing={2} justifyContent="center">
          <Button
            variant="contained"
            color="secondary"
            size="large"
            href="mailto:contact@espthinking.co.uk"
          >
            Get in Touch
          </Button>
          <Button
            variant="outlined"
            size="large"
            sx={{ color: 'white', borderColor: 'white' }}
            component={Link}
            to="/etp-profile"
          >
            Explore Demos
          </Button>
        </Stack>
      </Container>
    </Box>
  );
};

// Main portfolio landing page
const PortfolioLanding: React.FC = () => {
  return (
    <Box>
      <HeroSection />
      <ProjectsSection />
      <InnovationsSection />
      <ImpactSection />
      <CapabilitiesSection />
      <CTASection />
    </Box>
  );
};

export default PortfolioLanding;

