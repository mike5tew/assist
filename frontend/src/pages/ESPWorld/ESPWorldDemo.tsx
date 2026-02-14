import React, { useState } from 'react';
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
  Chip,
  Avatar,
  Slider,
  ToggleButton,
  ToggleButtonGroup,
  alpha,
  useTheme,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  LinearProgress,
  Tabs,
  Tab,
} from '@mui/material';
import {
  ArrowBack as BackIcon,
  School as SchoolIcon,
  Person as PersonIcon,
  Psychology as ETPIcon,
  Timeline as TimelineIcon,
  Notes as NotesIcon,
  TrendingUp as TrendingIcon,
  Group as GroupIcon,
  Verified as VerifiedIcon,
} from '@mui/icons-material';
import { Link } from 'react-router-dom';

// ─── Dummy Data: Student Profiles ───────────────────────────
const dummyStudents = [
  {
    id: 'student-001',
    name: 'Maya Chen',
    year: '6',
    form: 'Maple',
    etpProfile: {
      socialGravity: 28, // Slightly independent
      energyDirectionality: 65, // Outward (extroverted processor)
      voltageSensitivity: 72, // Conductive
      threatResponse: 45, // Balanced
      careResponse: 78, // Nurturing
      riskTolerance: 52, // Slightly risk-seeking
      integrityLogic: 58, // Slightly relativistic
      mirrorNeuronTuning: 68, // Absorbent
      orderliness: 42, // Flexible
    },
    recentNotes: [
      { date: '2026-02-13', content: 'Maya showed excellent collaboration during chemistry practical. High absorbency working well with peers.' },
      { date: '2026-02-12', content: 'Needs support shifting from group mode to independent focus work. Suggested "quiet corner" protocol.' },
    ],
    strengths: ['Communication', 'Collaboration', 'Empathy', 'Adaptability'],
    areasForSupport: ['Independent Focus', 'Emotional Boundaries'],
  },
  {
    id: 'student-002',
    name: 'Oliver Thompson',
    year: '6',
    form: 'Oak',
    etpProfile: {
      socialGravity: 78, // Very independent
      energyDirectionality: 28, // Inward (internal processor)
      voltageSensitivity: 25, // Insulated
      threatResponse: 32, // Passive (retreat style)
      careResponse: 55, // Balanced
      riskTolerance: 35, // Risk-averse
      integrityLogic: 72, // Absolutist
      mirrorNeuronTuning: 38, // Selective
      orderliness: 85, // Very ordered
    },
    recentNotes: [
      { date: '2026-02-13', content: 'Oliver thrives in structured environments. Pair programming with clear protocols gets best results.' },
      { date: '2026-02-11', content: 'Avoided group presentation. Consider offering written alternative or small group format.' },
    ],
    strengths: ['Focus', 'Precision', 'Independent Problem-Solving', 'Organisation'],
    areasForSupport: ['Group Communication', 'Risk-Taking in Learning', 'Emotional Processing'],
  },
  {
    id: 'student-003',
    name: 'Amara Okafor',
    year: '6',
    form: 'Maple',
    etpProfile: {
      socialGravity: 55, // Balanced
      energyDirectionality: 72, // Very outward
      voltageSensitivity: 58, // Moderately conductive
      threatResponse: 68, // Aggressive (assertive style)
      careResponse: 65, // Nurturing
      riskTolerance: 75, // Risk-seeking
      integrityLogic: 48, // Relativistic
      mirrorNeuronTuning: 62, // Moderately absorbent
      orderliness: 35, // Very flexible
    },
    recentNotes: [
      { date: '2026-02-13', content: 'Amara\'s assertiveness is channeling well. Leadership emerging in group projects.' },
      { date: '2026-02-10', content: 'High energy sometimes overwhelming peers. Set up "energetic ideas" capture system.' },
    ],
    strengths: ['Leadership', 'Initiative', 'Creative Problem-Solving', 'Resilience'],
    areasForSupport: ['Listening to Others', 'Structured Planning', 'Emotional Regulation'],
  },
];

// ─── Tab Navigation ─────────────────────────────────────────
interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;
  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`esp-tabpanel-${index}`}
      aria-labelledby={`esp-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ pt: 3 }}>{children}</Box>}
    </div>
  );
}

// ─── ETP Spectrum Slider Visualizer ───────────────────────
interface SpectrumDataPoint {
  label: string;
  leftLabel: string;
  rightLabel: string;
  value: number;
  color: string;
  description: string;
}

const spectrumDefinitions: Record<string, SpectrumDataPoint> = {
  socialGravity: {
    label: 'Social Gravity',
    leftLabel: 'Independent',
    rightLabel: 'Cohesive',
    value: 0,
    color: '#1565c0',
    description: 'Does this student recharge in solitude or with others?',
  },
  energyDirectionality: {
    label: 'Energy Directionality',
    leftLabel: 'Inward',
    rightLabel: 'Outward',
    value: 0,
    color: '#6a1b9a',
    description: 'Does this student process internally or externally (by talking)?',
  },
  voltageSensitivity: {
    label: 'Voltage Sensitivity',
    leftLabel: 'Insulated',
    rightLabel: 'Conductive',
    value: 0,
    color: '#f57c00',
    description: 'How much emotional current does this student absorb?',
  },
  threatResponse: {
    label: 'Threat Response',
    leftLabel: 'Passive (Retreat)',
    rightLabel: 'Aggressive (Assert)',
    value: 0,
    color: '#c62828',
    description: 'How does this student react when threatened?',
  },
  careResponse: {
    label: 'Care Response',
    leftLabel: 'Detached',
    rightLabel: 'Nurturing',
    value: 0,
    color: '#ad1457',
    description: 'How does this student respond to others\' vulnerability?',
  },
  riskTolerance: {
    label: 'Risk Tolerance',
    leftLabel: 'Averse',
    rightLabel: 'Seeking',
    value: 0,
    color: '#2e7d32',
    description: 'How does risk affect this student\'s emotional state?',
  },
  integrityLogic: {
    label: 'Integrity Logic',
    leftLabel: 'Relativistic',
    rightLabel: 'Absolutist',
    value: 0,
    color: '#4527a0',
    description: 'How does this student process moral decisions?',
  },
  mirrorNeuronTuning: {
    label: 'Mirror Neuron Tuning',
    leftLabel: 'Selective',
    rightLabel: 'Absorbent',
    value: 0,
    color: '#00838f',
    description: 'How much does this student absorb others\' emotional states?',
  },
  orderliness: {
    label: 'Orderliness',
    leftLabel: 'Flexible',
    rightLabel: 'Ordered',
    value: 0,
    color: '#5d4037',
    description: 'What\'s this student\'s preference for structure?',
  },
};

// ─── Spectrum Visualization Component ──────────────────────
const SpectrumVisualization: React.FC<{ etpProfile: Record<string, number> }> = ({ etpProfile }) => {
  return (
    <Stack spacing={3}>
      {Object.entries(etpProfile).map(([key, value]) => {
        const spec = spectrumDefinitions[key];
        if (!spec) return null;
        return (
          <Paper key={key} elevation={0} sx={{ p: 3, borderLeft: `4px solid ${spec.color}`, bgcolor: alpha(spec.color, 0.05) }}>
            <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
              {spec.label}
            </Typography>
            <Typography variant="caption" color="text.secondary" paragraph>
              {spec.description}
            </Typography>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 1 }}>
              <Typography variant="caption" sx={{ minWidth: 100 }}>
                {spec.leftLabel}
              </Typography>
              <LinearProgress
                variant="determinate"
                value={value}
                sx={{
                  flex: 1,
                  height: 8,
                  borderRadius: 4,
                  bgcolor: alpha(spec.color, 0.2),
                  '& .MuiLinearProgress-bar': {
                    bgcolor: spec.color,
                  },
                }}
              />
              <Typography variant="caption" sx={{ minWidth: 100, textAlign: 'right' }}>
                {spec.rightLabel}
              </Typography>
            </Box>
            <Typography variant="body2" fontWeight="bold" sx={{ color: spec.color }}>
              {value.toFixed(0)}/100
            </Typography>
          </Paper>
        );
      })}
    </Stack>
  );
};

// ─── Main Demo Component ────────────────────────────────────
const ESPWorldDemo: React.FC = () => {
  const theme = useTheme();
  const [selectedStudent, setSelectedStudent] = useState(dummyStudents[0]);
  const [tabValue, setTabValue] = useState(0);

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  return (
    <Box sx={{ minHeight: '100vh', bgcolor: '#f5f7f9' }}>
      {/* ── Header ── */}
      <Box sx={{ bgcolor: 'white', borderBottom: '1px solid', borderColor: 'divider', mb: 4 }}>
        <Container maxWidth="lg" sx={{ py: 3 }}>
          <Stack direction="row" spacing={2} alignItems="center">
            <Button
              component={Link}
              to="/esp-world"
              startIcon={<BackIcon />}
              variant="text"
              sx={{ color: 'text.secondary' }}
            >
              Back to ESP World
            </Button>
            <Box>
              <Typography variant="h5" gutterBottom fontWeight="bold">
                ESP World Demo
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Interactive demonstration of student profiles, ETP positioning, and teacher insights
              </Typography>
            </Box>
          </Stack>
        </Container>
      </Box>

      <Container maxWidth="lg">
        <Grid container spacing={4}>
          {/* ── Left Panel: Student Selector ── */}
          <Grid item xs={12} md={3}>
            <Typography variant="subtitle1" fontWeight="bold" gutterBottom sx={{ px: 2 }}>
              Class Roster
            </Typography>
            <Stack spacing={1.5}>
              {dummyStudents.map((student) => (
                <Card
                  key={student.id}
                  onClick={() => {
                    setSelectedStudent(student);
                    setTabValue(0); // Reset to first tab
                  }}
                  sx={{
                    cursor: 'pointer',
                    bgcolor: selectedStudent.id === student.id ? alpha(theme.palette.primary.main, 0.1) : 'white',
                    borderLeft: selectedStudent.id === student.id ? `4px solid ${theme.palette.primary.main}` : 'none',
                    transition: 'all 0.2s',
                    '&:hover': {
                      boxShadow: 2,
                      transform: 'translateX(4px)',
                    },
                  }}
                >
                  <CardContent sx={{ p: 2 }}>
                    <Stack direction="row" spacing={2} alignItems="center">
                      <Avatar sx={{ bgcolor: theme.palette.primary.main, width: 40, height: 40 }}>
                        {student.name.split(' ')[0][0]}
                        {student.name.split(' ')[1][0]}
                      </Avatar>
                      <Box>
                        <Typography variant="subtitle2" fontWeight="bold">
                          {student.name}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                          Year {student.year}, {student.form}
                        </Typography>
                      </Box>
                    </Stack>
                  </CardContent>
                </Card>
              ))}
            </Stack>
          </Grid>

          {/* ── Right Panel: Student Details ── */}
          <Grid item xs={12} md={9}>
            <Paper elevation={2} sx={{ p: 4, borderRadius: 3 }}>
              {/* ── Student Header ── */}
              <Stack direction="row" spacing={3} alignItems="center" sx={{ mb: 4 }}>
                <Avatar sx={{ bgcolor: theme.palette.primary.main, width: 80, height: 80, fontSize: '2rem' }}>
                  {selectedStudent.name.split(' ')[0][0]}
                  {selectedStudent.name.split(' ')[1][0]}
                </Avatar>
                <Box>
                  <Typography variant="h4" fontWeight="bold" gutterBottom>
                    {selectedStudent.name}
                  </Typography>
                  <Stack direction="row" spacing={2}>
                    <Chip icon={<SchoolIcon />} label={`Year ${selectedStudent.year}`} variant="outlined" />
                    <Chip label={`Form ${selectedStudent.form}`} variant="outlined" />
                    <Chip icon={<VerifiedIcon />} label="ETP Profile Complete" variant="outlined" color="success" />
                  </Stack>
                </Box>
              </Stack>

              {/* ── Tabs ── */}
              <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
                <Tabs
                  value={tabValue}
                  onChange={handleTabChange}
                  aria-label="Student Profile Tabs"
                  sx={{ '& .MuiTab-root': { textTransform: 'none' } }}
                >
                  <Tab label="ETP Spectra" icon={<ETPIcon />} iconPosition="start" id="esp-tab-0" />
                  <Tab label="Recent Notes" icon={<NotesIcon />} iconPosition="start" id="esp-tab-1" />
                  <Tab label="Strengths & Support" icon={<TrendingIcon />} iconPosition="start" id="esp-tab-2" />
                </Tabs>
              </Box>

              {/* ── Tab 1: ETP Spectra ── */}
              <TabPanel value={tabValue} index={0}>
                <Typography variant="body2" color="text.secondary" paragraph>
                  This student's positioning across the 9 biological spectra. Use this to understand how they're likely to react,
                  process information, and interact with peers.
                </Typography>
                <SpectrumVisualization etpProfile={selectedStudent.etpProfile} />
              </TabPanel>

              {/* ── Tab 2: Recent Notes ── */}
              <TabPanel value={tabValue} index={1}>
                <Typography variant="body2" color="text.secondary" paragraph>
                  Teacher observations logged using real-time ETP sensor tapping (via PrimaryOS).
                </Typography>
                <Stack spacing={2}>
                  {selectedStudent.recentNotes.map((note, idx) => (
                    <Paper key={idx} elevation={0} sx={{ p: 3, bgcolor: alpha(theme.palette.info.main, 0.05), borderLeft: `4px solid ${theme.palette.info.main}` }}>
                      <Typography variant="caption" color="text.secondary" sx={{ fontWeight: 'bold' }}>
                        {new Date(note.date).toLocaleDateString()}
                      </Typography>
                      <Typography variant="body2" sx={{ mt: 1 }}>
                        {note.content}
                      </Typography>
                    </Paper>
                  ))}
                </Stack>
              </TabPanel>

              {/* ── Tab 3: Strengths & Support ── */}
              <TabPanel value={tabValue} index={2}>
                <Grid container spacing={3}>
                  <Grid item xs={12} sm={6}>
                    <Paper elevation={0} sx={{ p: 3, bgcolor: alpha('#2e7d32', 0.05), borderLeft: '4px solid #2e7d32' }}>
                      <Typography variant="subtitle2" fontWeight="bold" gutterBottom sx={{ color: '#2e7d32' }}>
                        Strengths
                      </Typography>
                      <Stack spacing={1}>
                        {selectedStudent.strengths.map((strength, idx) => (
                          <Stack key={idx} direction="row" spacing={1} alignItems="center">
                            <Box sx={{ width: 8, height: 8, borderRadius: '50%', bgcolor: '#2e7d32' }} />
                            <Typography variant="body2">{strength}</Typography>
                          </Stack>
                        ))}
                      </Stack>
                    </Paper>
                  </Grid>
                  <Grid item xs={12} sm={6}>
                    <Paper elevation={0} sx={{ p: 3, bgcolor: alpha('#f57c00', 0.05), borderLeft: '4px solid #f57c00' }}>
                      <Typography variant="subtitle2" fontWeight="bold" gutterBottom sx={{ color: '#f57c00' }}>
                        Areas for Support
                      </Typography>
                      <Stack spacing={1}>
                        {selectedStudent.areasForSupport.map((area, idx) => (
                          <Stack key={idx} direction="row" spacing={1} alignItems="center">
                            <Box sx={{ width: 8, height: 8, borderRadius: '50%', bgcolor: '#f57c00' }} />
                            <Typography variant="body2">{area}</Typography>
                          </Stack>
                        ))}
                      </Stack>
                    </Paper>
                  </Grid>
                </Grid>
              </TabPanel>
            </Paper>
          </Grid>
        </Grid>

        {/* ── Footer Info ── */}
        <Paper elevation={0} sx={{ mt: 6, p: 4, bgcolor: alpha(theme.palette.primary.main, 0.05), borderRadius: 2 }}>
          <Stack spacing={2}>
            <Typography variant="h6" fontWeight="bold">
              About This Demo
            </Typography>
            <Typography variant="body2" color="text.secondary">
              This demonstration showcases the core features of ESP World:
            </Typography>
            <Stack spacing={1.5} sx={{ ml: 2 }}>
              <Typography variant="body2">
                <strong>ETP Profiles:</strong> Each student's positioning across 9 biological spectra, capturing how they naturally operate.
              </Typography>
              <Typography variant="body2">
                <strong>Real-Time Observations:</strong> Teachers use PrimaryOS to log ETP-informed notes during lessons, building a rich longitudinal picture.
              </Typography>
              <Typography variant="body2">
                <strong>Actionable Insights:</strong> Spectra positioning informs matched teaching strategies, not labels or deficits.
              </Typography>
              <Typography variant="body2">
                <strong>Data Privacy:</strong> All student data is school-owned and accessible only to authorized staff. No external aggregation or algorithmic ranking.
              </Typography>
            </Stack>
          </Stack>
        </Paper>
      </Container>
    </Box>
  );
};

export default ESPWorldDemo;
