import React, { useState } from 'react';
import {
  Box,
  Container,
  Typography,
  Button,
  Grid,
  Card,
  CardContent,
  CardActions,
  Stack,
  Paper,
  Chip,
  Avatar,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  TextareaAutosize,
  alpha,
  useTheme,
  IconButton,
  Menu,
  Divider,
  LinearProgress,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
} from '@mui/material';
import {
  Edit as EditIcon,
  Delete as DeleteIcon,
  Add as AddIcon,
  MoreVert as MoreIcon,
  ArrowBack as BackIcon,
  School as SchoolIcon,
  EventNote as EventIcon,
  Assignment as AssignmentIcon,
  CheckCircle as CheckIcon,
  AccessTime as ClockIcon,
  Group as GroupIcon,
  AdsClick as TargetIcon,
  Help as HelpIcon,
} from '@mui/icons-material';
import { Link } from 'react-router-dom';

// ─── Student Competency Interface ─────────────────────────
interface StudentCompetency {
  studentId: string;
  studentName: string;
  competencies: {
    [key: string]: number; // 1-10 scale
  };
}

// ─── Lesson Activity Interface ────────────────────────────
interface LessonActivity {
  id: string;
  name: string;
  description: string;
  skillsFoci: string[]; // which competencies this activity develops
  difficulty: 'beginner' | 'intermediate' | 'advanced';
}

// ─── Work Assignment Interface ────────────────────────────
interface WorkAssignment {
  id: string;
  studentId: string;
  activityId: string;
  assignedDate: string;
  dueDate: string;
  status: 'pending' | 'submitted' | 'graded';
  marks?: number;
  feedback?: string;
  targetCompetencies: string[]; // which competencies this assignment targets for THIS student
}

// ─── Dummy Lesson Plan Data ───────────────────────────────
interface LessonPlan {
  id: string;
  title: string;
  subject: string;
  year: string;
  date: string;
  duration: string; // minutes
  objectives: string[];
  starter: string;
  mainActivity: string;
  plenary: string;
  resources: string[];
  assessment: string;
  status: 'draft' | 'published' | 'completed';
  createdBy: string;
  lastModified: string;
  // ─── Skills Mapping ───
  targetCompetencies: {
    name: string;
    developmentLevel: number; // 1-10 what level students should reach
    focus: 'core' | 'supporting' | 'enrichment'; // how central to the lesson
  }[];
  studentAlignments?: {
    [studentId: string]: {
      readinessScore: number; // 0-100: how ready is student for this lesson
      developmentOpportunity: string; // tailored note
    };
  };
  // ─── Activities & Assignments ───
  activities: LessonActivity[];
  workAssignments: WorkAssignment[];
}

// ─── Competency Framework ─────────────────────────────────
const COMPETENCY_FRAMEWORK = {
  analyticalThinking: 'Analytical Thinking',
  scientificMethod: 'Scientific Method',
  criticalAnalysis: 'Critical Analysis',
  codingLogic: 'Coding & Logic',
  dataInterpretation: 'Data Interpretation',
  problemSolving: 'Problem Solving',
  collaboration: 'Collaboration',
  communication: 'Communication',
  creativity: 'Creativity',
  perseverance: 'Perseverance',
};

// ─── Dummy Student Data ────────────────────────────────────
const dummyStudents: StudentCompetency[] = [
  {
    studentId: 's-001',
    studentName: 'Alex Rivera',
    competencies: {
      analyticalThinking: 7,
      scientificMethod: 8,
      criticalAnalysis: 6,
      codingLogic: 4,
      dataInterpretation: 7,
      problemSolving: 8,
      collaboration: 7,
      communication: 6,
      creativity: 5,
      perseverance: 8,
    },
  },
  {
    studentId: 's-002',
    studentName: 'Jordan Chen',
    competencies: {
      analyticalThinking: 9,
      scientificMethod: 7,
      criticalAnalysis: 9,
      codingLogic: 9,
      dataInterpretation: 8,
      problemSolving: 9,
      collaboration: 5,
      communication: 6,
      creativity: 8,
      perseverance: 7,
    },
  },
  {
    studentId: 's-003',
    studentName: 'Sam Taylor',
    competencies: {
      analyticalThinking: 5,
      scientificMethod: 4,
      criticalAnalysis: 5,
      codingLogic: 3,
      dataInterpretation: 4,
      problemSolving: 6,
      collaboration: 9,
      communication: 8,
      creativity: 9,
      perseverance: 5,
    },
  },
  {
    studentId: 's-004',
    studentName: 'Morgan Lewis',
    competencies: {
      analyticalThinking: 6,
      scientificMethod: 5,
      criticalAnalysis: 6,
      codingLogic: 5,
      dataInterpretation: 6,
      problemSolving: 7,
      collaboration: 8,
      communication: 7,
      creativity: 6,
      perseverance: 8,
    },
  },
];

const dummyLessonPlans: LessonPlan[] = [
  {
    id: 'lp-001',
    title: 'Photosynthesis: The Light Reactions',
    subject: 'Biology',
    year: '10',
    date: '2026-02-14',
    duration: '60',
    objectives: [
      'Understand the location and function of light-dependent reactions',
      'Explain the role of chlorophyll and photosystems',
      'Describe electron transport chains in thylakoid membranes',
    ],
    starter: 'Show time-lapse video of plant growth in sunlight vs. darkness (5 min). Ask: "What\'s the difference and why?"',
    mainActivity: 'Practical: Use lamp and algae suspension to measure O2 production. Students predict, test, and analyse results.',
    plenary: 'Concept mapping: Students build a diagram linking light → chlorophyll → electrons → ATP + NADPH.',
    resources: ['Lamp apparatus', 'Algae suspension', 'Oxygen probe', 'Whiteboard markers', 'A3 paper'],
    assessment: 'Observation checklist during practical; concept map peer-review; exit question: "Explain why plants are green, not red."',
    status: 'published',
    createdBy: 'Dr. Sarah Chen',
    lastModified: '2026-02-12',
    // ─── Skills Mapping ───
    targetCompetencies: [
      { name: 'scientificMethod', developmentLevel: 8, focus: 'core' },
      { name: 'analyticalThinking', developmentLevel: 7, focus: 'core' },
      { name: 'dataInterpretation', developmentLevel: 7, focus: 'core' },
      { name: 'collaboration', developmentLevel: 6, focus: 'supporting' },
      { name: 'communication', developmentLevel: 6, focus: 'supporting' },
      { name: 'perseverance', developmentLevel: 5, focus: 'supporting' },
    ],
    studentAlignments: {
      's-001': {
        readinessScore: 85,
        developmentOpportunity: 'Strong in scientific method; can challenge with design modifications',
      },
      's-002': {
        readinessScore: 82,
        developmentOpportunity: 'Excellent analytical skills; pair with Sam for peer teaching',
      },
      's-003': {
        readinessScore: 55,
        developmentOpportunity: 'Needs scaffolding on interpretation; use collaborative practical',
      },
      's-004': {
        readinessScore: 72,
        developmentOpportunity: 'Solid foundations; consider extension question on photosystem structures',
      },
    },
    // ─── Activities & Work Assignments ───
    activities: [
      {
        id: 'act-001-a',
        name: 'Practical Lab: Measure O₂ Production',
        description: 'Students set up apparatus, predict outcomes, conduct experiment with algae suspension under different light conditions, record data',
        skillsFoci: ['scientificMethod', 'dataInterpretation', 'analyticalThinking'],
        difficulty: 'intermediate',
      },
      {
        id: 'act-001-b',
        name: 'Concept Mapping Exercise',
        description: 'Create visual diagram linking light energy → chlorophyll → electron transport → ATP and NADPH production',
        skillsFoci: ['analyticalThinking', 'communication', 'creativity'],
        difficulty: 'intermediate',
      },
      {
        id: 'act-001-c',
        name: 'Challenge: Design an Experiment',
        description: 'Design a novel experiment to test photosystem organization using given constraints and materials',
        skillsFoci: ['scientificMethod', 'problemSolving', 'creativity'],
        difficulty: 'advanced',
      },
      {
        id: 'act-001-d',
        name: 'Peer Teaching & Scaffolded Notes',
        description: 'Work with partner to complete guided worksheet explaining light reactions step-by-step with visual aids',
        skillsFoci: ['analyticalThinking', 'collaboration', 'communication'],
        difficulty: 'beginner',
      },
    ],
    workAssignments: [
      // Alex: Ready, so gets lab + mapping
      {
        id: 'wa-001',
        studentId: 's-001',
        activityId: 'act-001-a',
        assignedDate: '2026-02-14',
        dueDate: '2026-02-14',
        status: 'pending',
        targetCompetencies: ['scientificMethod', 'dataInterpretation'],
      },
      {
        id: 'wa-002',
        studentId: 's-001',
        activityId: 'act-001-b',
        assignedDate: '2026-02-14',
        dueDate: '2026-02-14',
        status: 'pending',
        targetCompetencies: ['analyticalThinking', 'communication'],
      },
      // Jordan: Advanced, so gets challenge + lab
      {
        id: 'wa-003',
        studentId: 's-002',
        activityId: 'act-001-c',
        assignedDate: '2026-02-14',
        dueDate: '2026-02-14',
        status: 'pending',
        targetCompetencies: ['problemSolving', 'creativity'],
      },
      {
        id: 'wa-004',
        studentId: 's-002',
        activityId: 'act-001-a',
        assignedDate: '2026-02-14',
        dueDate: '2026-02-14',
        status: 'pending',
        targetCompetencies: ['scientificMethod', 'analyticalThinking'],
      },
      // Sam: Needs support, gets scaffolded + peer teaching
      {
        id: 'wa-005',
        studentId: 's-003',
        activityId: 'act-001-d',
        assignedDate: '2026-02-14',
        dueDate: '2026-02-14',
        status: 'pending',
        targetCompetencies: ['analyticalThinking', 'collaboration'],
      },
      {
        id: 'wa-006',
        studentId: 's-003',
        activityId: 'act-001-b',
        assignedDate: '2026-02-14',
        dueDate: '2026-02-14',
        status: 'pending',
        targetCompetencies: ['communication', 'creativity'],
      },
      // Morgan: Solid, gets lab + mapping
      {
        id: 'wa-007',
        studentId: 's-004',
        activityId: 'act-001-a',
        assignedDate: '2026-02-14',
        dueDate: '2026-02-14',
        status: 'pending',
        targetCompetencies: ['scientificMethod', 'dataInterpretation'],
      },
      {
        id: 'wa-008',
        studentId: 's-004',
        activityId: 'act-001-c',
        assignedDate: '2026-02-14',
        dueDate: '2026-02-15',
        status: 'pending',
        targetCompetencies: ['problemSolving'],
      },
    ],
  },
  {
    id: 'lp-002',
    title: 'The English Civil War: Power & Persuasion',
    subject: 'History',
    year: '9',
    date: '2026-02-15',
    duration: '50',
    objectives: [
      'Evaluate competing historical narratives',
      'Analyse propaganda from both Royalist and Parliamentarian sides',
      'Construct evidence-based arguments about causation',
    ],
    starter: 'Display two contemporary propaganda posters (10 min). Jigsaw: Year 9s split into groups, each analyzes one poster.',
    mainActivity: 'Source analysis task: Evaluate reliability of personal accounts vs. official records. Groups present competing interpretations.',
    plenary: 'Whole-class debate: "Was the Civil War avoidable?" Students vote pre- and post-evidence.',
    resources: ['Propaganda poster reproductions', 'Primary source extracts (printed)', 'Debate cards', 'Whiteboard'],
    assessment: 'Formative: Note-taking during jigsaw; source evaluation checklist; debate contribution rubric.',
    status: 'published',
    createdBy: 'Mr. James Whitmore',
    lastModified: '2026-02-10',
    // ─── Skills Mapping ───
    targetCompetencies: [
      { name: 'criticalAnalysis', developmentLevel: 8, focus: 'core' },
      { name: 'analyticalThinking', developmentLevel: 7, focus: 'core' },
      { name: 'communication', developmentLevel: 7, focus: 'core' },
      { name: 'collaboration', developmentLevel: 7, focus: 'core' },
      { name: 'creativity', developmentLevel: 6, focus: 'supporting' },
      { name: 'perseverance', developmentLevel: 5, focus: 'supporting' },
    ],
    studentAlignments: {
      's-001': {
        readinessScore: 70,
        developmentOpportunity: 'Good analytical skills; needs push on communication during debate',
      },
      's-002': {
        readinessScore: 92,
        developmentOpportunity: 'Exceptional fit; consider debate leadership role',
      },
      's-003': {
        readinessScore: 88,
        developmentOpportunity: 'Strength in collaboration and communication; ideal group facilitator',
      },
      's-004': {
        readinessScore: 78,
        developmentOpportunity: 'Solid match; good opportunity to develop critical analysis',
      },
    },
    activities: [],
    workAssignments: [],
  },
  {
    id: 'lp-003',
    title: 'Introduction to Python: Variables & Data Types',
    subject: 'Computer Science',
    year: '8',
    date: '2026-02-16',
    duration: '60',
    objectives: [
      'Understand variables as named storage locations',
      'Distinguish between data types (int, float, string, bool)',
      'Write simple Python programs using input() and print()',
    ],
    starter: 'Live coding demo: Show how Python stores a calculation step-by-step. Ask students to predict each result.',
    mainActivity: 'Guided practice: Students code along on shared screen, then attempt independent "challenge" problems.',
    plenary: 'Code review: Students swap 2-line programs, peer-check for syntax and logic errors.',
    resources: ['Python IDLE or Replit', 'Handout: Data Types Cheat Sheet', 'Challenge problems (printed)'],
    assessment: 'Execution of peer-reviewed code; formative quizzes using Kahoot (data type naming).',
    status: 'draft',
    createdBy: 'Ms. Lisa Patel',
    lastModified: '2026-02-13',
    // ─── Skills Mapping ───
    targetCompetencies: [
      { name: 'codingLogic', developmentLevel: 7, focus: 'core' },
      { name: 'problemSolving', developmentLevel: 7, focus: 'core' },
      { name: 'analyticalThinking', developmentLevel: 6, focus: 'core' },
      { name: 'perseverance', developmentLevel: 6, focus: 'supporting' },
      { name: 'collaboration', developmentLevel: 5, focus: 'supporting' },
      { name: 'communication', developmentLevel: 5, focus: 'supporting' },
    ],
    studentAlignments: {
      's-001': {
        readinessScore: 45,
        developmentOpportunity: 'Weaker in coding; provide extra scaffolding and worked examples',
      },
      's-002': {
        readinessScore: 95,
        developmentOpportunity: 'Exceptional coder; extension: explore advanced data structures',
      },
      's-003': {
        readinessScore: 35,
        developmentOpportunity: 'Significant challenge for this student; 1:1 support or differentiated tasks recommended',
      },
      's-004': {
        readinessScore: 62,
        developmentOpportunity: 'Moderate readiness; pair with Jordan for peer programming sessions',
      },
    },
    activities: [],
    workAssignments: [],
  },
  {
    id: 'lp-004',
    title: 'Quadratic Equations: Factoring & Solving',
    subject: 'Mathematics',
    year: '10',
    date: '2026-02-17',
    duration: '55',
    objectives: [
      'Factor quadratic expressions (ax² + bx + c)',
      'Solve quadratic equations using factoring, completing the square, and the quadratic formula',
      'Interpret solutions in context',
    ],
    starter: 'Mental math: Factor simple binomials (x² + 5x + 6). Use think-pair-share.',
    mainActivity: 'Stations: Station 1 (factoring), Station 2 (completing the square), Station 3 (formula). Groups rotate every 12 min.',
    plenary: 'Gallery walk: Display student work showing different solution methods for the same equation.',
    resources: ['Algebra tiles', 'Factoring templates (printed)', 'Graphing calculators', 'Mini whiteboards'],
    assessment: 'Formative: Observation at stations; mini-quiz on formula application; homework set.',
    status: 'completed',
    createdBy: 'Mr. Prabhat Kumar',
    lastModified: '2026-02-02',
    // ─── Skills Mapping ───
    targetCompetencies: [
      { name: 'problemSolving', developmentLevel: 8, focus: 'core' },
      { name: 'analyticalThinking', developmentLevel: 8, focus: 'core' },
      { name: 'perseverance', developmentLevel: 7, focus: 'core' },
      { name: 'collaboration', developmentLevel: 6, focus: 'supporting' },
      { name: 'communication', developmentLevel: 5, focus: 'supporting' },
      { name: 'creativity', developmentLevel: 5, focus: 'enrichment' },
    ],
    studentAlignments: {
      's-001': {
        readinessScore: 88,
        developmentOpportunity: 'Excellent problem-solver; ready for complex multi-step equations',
      },
      's-002': {
        readinessScore: 95,
        developmentOpportunity: 'Advanced readiness; consider enrichment with quadratic applications',
      },
      's-003': {
        readinessScore: 52,
        developmentOpportunity: 'Needs patience-building; use visual algebra tiles and collaborative stations',
      },
      's-004': {
        readinessScore: 80,
        developmentOpportunity: 'Good fit; strong fundamentals with room to deepen formula flexibility',
      },
    },
    activities: [],
    workAssignments: [],
  },
];

// ─── Status Badge Component ───────────────────────────────
const StatusBadge: React.FC<{ status: LessonPlan['status'] }> = ({ status }) => {
  const statusConfig = {
    draft: { color: 'default', label: 'Draft' },
    published: { color: 'success', label: 'Published' },
    completed: { color: 'info', label: 'Completed' },
  };
  const config = statusConfig[status];
  return <Chip label={config.label} color={config.color as any} size="small" />;
};

// ─── Lesson Plan Card Component ───────────────────────────
interface LessonCardProps {
  plan: LessonPlan;
  onEdit: (plan: LessonPlan) => void;
  onDelete: (id: string) => void;
  onView: (plan: LessonPlan) => void;
}

const LessonCard: React.FC<LessonCardProps> = ({ plan, onEdit, onDelete, onView }) => {
  const theme = useTheme();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  return (
    <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column', transition: 'all 0.2s', '&:hover': { boxShadow: 4, transform: 'translateY(-4px)' } }}>
      <CardContent sx={{ flexGrow: 1 }}>
        <Stack direction="row" justifyContent="space-between" alignItems="start" sx={{ mb: 2 }}>
          <Box>
            <Typography variant="h6" fontWeight="bold" gutterBottom>
              {plan.title}
            </Typography>
            <Stack direction="row" spacing={1}>
              <Chip icon={<SchoolIcon />} label={plan.subject} size="small" variant="outlined" />
              <Chip label={`Year ${plan.year}`} size="small" variant="outlined" />
            </Stack>
          </Box>
          <IconButton size="small" onClick={handleMenuOpen}>
            <MoreIcon />
          </IconButton>
          <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={handleMenuClose}>
            <MenuItem onClick={() => { onEdit(plan); handleMenuClose(); }}>
              <EditIcon sx={{ mr: 1 }} /> Edit
            </MenuItem>
            <MenuItem onClick={() => { onDelete(plan.id); handleMenuClose(); }}>
              <DeleteIcon sx={{ mr: 1 }} /> Delete
            </MenuItem>
          </Menu>
        </Stack>

        <Stack spacing={2} sx={{ my: 2 }}>
          <Box>
            <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
              Date & Duration
            </Typography>
            <Typography variant="body2">
              {new Date(plan.date).toLocaleDateString()} • {plan.duration} min
            </Typography>
          </Box>

          <Box>
            <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
              Learning Objectives
            </Typography>
            <Stack spacing={0.5}>
              {plan.objectives.slice(0, 2).map((obj, idx) => (
                <Typography key={idx} variant="body2" sx={{ pl: 1 }}>
                  • {obj}
                </Typography>
              ))}
              {plan.objectives.length > 2 && (
                <Typography variant="caption" color="text.secondary" sx={{ pl: 1 }}>
                  +{plan.objectives.length - 2} more
                </Typography>
              )}
            </Stack>
          </Box>

          <Box>
            <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
              Skills Developed
            </Typography>
            <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
              {plan.targetCompetencies.slice(0, 3).map((comp) => (
                <Chip
                  key={comp.name}
                  label={(COMPETENCY_FRAMEWORK as any)[comp.name]}
                  size="small"
                  variant={comp.focus === 'core' ? 'filled' : 'outlined'}
                  color={comp.focus === 'core' ? 'primary' : 'default'}
                />
              ))}
              {plan.targetCompetencies.length > 3 && (
                <Chip label={`+${plan.targetCompetencies.length - 3}`} size="small" variant="outlined" />
              )}
            </Stack>
          </Box>

          <Divider />

          <Stack direction="row" spacing={2}>
            <Stack direction="row" spacing={0.5} alignItems="center">
              <ClockIcon sx={{ fontSize: 18, color: 'text.secondary' }} />
              <Typography variant="caption" color="text.secondary">
                By {plan.createdBy.split(' ')[0]}
              </Typography>
            </Stack>
            <StatusBadge status={plan.status} />
          </Stack>
        </Stack>
      </CardContent>

      <CardActions sx={{ pt: 0 }}>
        <Button size="small" variant="text" onClick={() => onView(plan)}>View Details</Button>
      </CardActions>
    </Card>
  );
};

// ─── Competency Visualization Component ────────────────────
interface CompetencyVisualizationProps {
  plan: LessonPlan;
}

const CompetencyVisualization: React.FC<CompetencyVisualizationProps> = ({ plan }) => {
  const theme = useTheme();

  const getFocusColor = (focus: 'core' | 'supporting' | 'enrichment') => {
    if (focus === 'core') return theme.palette.primary.main;
    if (focus === 'supporting') return theme.palette.secondary.main;
    return theme.palette.warning.main;
  };

  const getReadinessColor = (score: number) => {
    if (score >= 80) return '#2e7d32'; // Green - Ready
    if (score >= 60) return '#f57c00'; // Orange - Moderate
    return '#c62828'; // Red - Needs support
  };

  return (
    <Paper sx={{ p: 3, bgcolor: alpha(theme.palette.primary.main, 0.02) }}>
      <Stack spacing={3}>
        {/* ─── Target Competencies ─── */}
        <Box>
          <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
            Skills Developed
          </Typography>
          <Stack spacing={1}>
            {plan.targetCompetencies.map((comp) => (
              <Box key={comp.name}>
                <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ mb: 0.5 }}>
                  <Typography variant="body2">
                    {(COMPETENCY_FRAMEWORK as any)[comp.name]}
                  </Typography>
                  <Chip
                    label={comp.focus === 'core' ? 'Core' : comp.focus === 'supporting' ? 'Supporting' : 'Enrichment'}
                    size="small"
                    sx={{
                      bgcolor: getFocusColor(comp.focus),
                      color: 'white',
                      height: 20,
                    }}
                  />
                </Stack>
                <Box sx={{ height: 6, bgcolor: '#e0e0e0', borderRadius: 1, overflow: 'hidden' }}>
                  <Box
                    sx={{
                      height: '100%',
                      width: `${(comp.developmentLevel / 10) * 100}%`,
                      bgcolor: getFocusColor(comp.focus),
                      transition: 'width 0.3s',
                    }}
                  />
                </Box>
              </Box>
            ))}
          </Stack>
        </Box>

        <Divider />

        {/* ─── Student Readiness ─── */}
        {plan.studentAlignments && (
          <Box>
            <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
              Student Readiness & Differentiation
            </Typography>
            <Stack spacing={2}>
              {dummyStudents.map((student) => {
                const alignment = plan.studentAlignments![student.studentId];
                if (!alignment) return null;

                return (
                  <Paper key={student.studentId} variant="outlined" sx={{ p: 2 }}>
                    <Stack spacing={1}>
                      <Stack direction="row" justifyContent="space-between" alignItems="center">
                        <Typography variant="body2" fontWeight="600">
                          {student.studentName}
                        </Typography>
                        <Box
                          sx={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: 1,
                          }}
                        >
                          <Box sx={{ width: 120, height: 6, bgcolor: '#e0e0e0', borderRadius: 1, overflow: 'hidden' }}>
                            <Box
                              sx={{
                                height: '100%',
                                width: `${alignment.readinessScore}%`,
                                bgcolor: getReadinessColor(alignment.readinessScore),
                                transition: 'width 0.3s',
                              }}
                            />
                          </Box>
                          <Typography variant="caption" fontWeight="bold" sx={{ minWidth: 45 }}>
                            {alignment.readinessScore}%
                          </Typography>
                        </Box>
                      </Stack>
                      <Typography variant="caption" color="text.secondary" sx={{ pl: 1 }}>
                        {alignment.developmentOpportunity}
                      </Typography>
                    </Stack>
                  </Paper>
                );
              })}
            </Stack>
          </Box>
        )}
      </Stack>
    </Paper>
  );
};

// ─── Lesson Detail Dialog ─────────────────────────────────
interface LessonDetailDialogProps {
  open: boolean;
  plan: LessonPlan | null;
  onClose: () => void;
}

const LessonDetailDialog: React.FC<LessonDetailDialogProps> = ({ open, plan, onClose }) => {
  if (!plan) return null;

  return (
    <Dialog open={open} onClose={onClose} maxWidth="lg" fullWidth>
      <DialogTitle sx={{ pb: 1 }}>
        <Typography variant="h6" fontWeight="bold">
          {plan.title}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          {plan.subject} • Year {plan.year} • {plan.duration} minutes
        </Typography>
      </DialogTitle>
      <DialogContent sx={{ pt: 2 }}>
        <Stack spacing={3}>
          {/* ─── Competencies First ─── */}
          <CompetencyVisualization plan={plan} />

          <Divider />

          {/* ─── Lesson Content ─── */}
          <Stack spacing={2}>
            <Box>
              <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
                Learning Objectives
              </Typography>
              <List dense>
                {plan.objectives.map((obj, idx) => (
                  <ListItem key={idx}>
                    <ListItemIcon>
                      <CheckIcon sx={{ fontSize: 20, color: 'success.main' }} />
                    </ListItemIcon>
                    <ListItemText primary={obj} />
                  </ListItem>
                ))}
              </List>
            </Box>

            <Box>
              <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
                Starter (5-10 min)
              </Typography>
              <Typography variant="body2" sx={{ pl: 1 }}>
                {plan.starter}
              </Typography>
            </Box>

            <Box>
              <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
                Main Activity (30-40 min)
              </Typography>
              <Typography variant="body2" sx={{ pl: 1 }}>
                {plan.mainActivity}
              </Typography>
            </Box>

            <Box>
              <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
                Plenary (8-10 min)
              </Typography>
              <Typography variant="body2" sx={{ pl: 1 }}>
                {plan.plenary}
              </Typography>
            </Box>

            <Box>
              <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
                Resources
              </Typography>
              <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                {plan.resources.map((resource, idx) => (
                  <Chip key={idx} label={resource} size="small" variant="outlined" />
                ))}
              </Stack>
            </Box>

            <Box>
              <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
                Assessment
              </Typography>
              <Typography variant="body2" sx={{ pl: 1 }}>
                {plan.assessment}
              </Typography>
            </Box>

            {/* ─── Work Assignments ─── */}
            {plan.workAssignments.length > 0 && (
              <Box>
                <Divider sx={{ my: 2 }} />
                <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
                  Work Assignments by Student
                </Typography>
                <Stack spacing={2}>
                  {dummyStudents.map((student) => {
                    const assignments = plan.workAssignments.filter((wa) => wa.studentId === student.studentId);
                    if (assignments.length === 0) return null;

                    return (
                      <Paper key={student.studentId} variant="outlined" sx={{ p: 2 }}>
                        <Typography variant="body2" fontWeight="bold" gutterBottom>
                          {student.studentName}
                        </Typography>
                        <Stack spacing={1}>
                          {assignments.map((assignment) => {
                            const activity = plan.activities.find((a) => a.id === assignment.activityId);
                            return (
                              <Box key={assignment.id} sx={{ pl: 1, py: 1, bgcolor: '#f5f5f5', borderRadius: 1, px: 1.5 }}>
                                <Typography variant="body2" fontWeight="600">
                                  {activity?.name}
                                </Typography>
                                <Typography variant="caption" color="text.secondary" display="block">
                                  {activity?.description}
                                </Typography>
                                <Stack direction="row" spacing={1} sx={{ mt: 1 }} flexWrap="wrap">
                                  {assignment.targetCompetencies.map((comp) => (
                                    <Chip
                                      key={comp}
                                      label={(COMPETENCY_FRAMEWORK as any)[comp]}
                                      size="small"
                                      variant="outlined"
                                      color={activity?.difficulty === 'advanced' ? 'warning' : 'default'}
                                    />
                                  ))}
                                </Stack>
                              </Box>
                            );
                          })}
                        </Stack>
                      </Paper>
                    );
                  })}
                </Stack>
              </Box>
            )}
          </Stack>
        </Stack>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
};

// ─── Help Dialog ───────────────────────────────────────────
interface HelpDialogProps {
  open: boolean;
  onClose: () => void;
}

const HelpDialog: React.FC<HelpDialogProps> = ({ open, onClose }) => {
  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle sx={{ pb: 1 }}>
        <Stack direction="row" spacing={1} alignItems="center">
          <HelpIcon />
          <Typography variant="h6" fontWeight="bold">How This Works</Typography>
        </Stack>
      </DialogTitle>
      <DialogContent sx={{ pt: 2 }}>
        <Stack spacing={2}>
          <Box sx={{ bgcolor: '#fff3e0', borderLeft: '4px solid #ff9800', p: 2, borderRadius: 1 }}>
            <Typography variant="body2" fontWeight="bold" color="#e65100">
              ℹ️ This is a demonstration of the lesson planning system
            </Typography>
          </Box>

          <Box>
            <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
              The Workflow
            </Typography>
            <Stack spacing={1}>
              <Box sx={{ display: 'flex', gap: 2 }}>
                <Box sx={{ color: 'primary.main', fontWeight: 'bold', minWidth: 24 }}>1</Box>
                <Box>
                  <Typography variant="body2" fontWeight="600">Lessons</Typography>
                  <Typography variant="caption" color="text.secondary">
                    Create structured lessons with learning objectives, starter/main/plenary activities, resources, and assessment strategies
                  </Typography>
                </Box>
              </Box>
              <Box sx={{ display: 'flex', gap: 2 }}>
                <Box sx={{ color: 'primary.main', fontWeight: 'bold', minWidth: 24 }}>2</Box>
                <Box>
                  <Typography variant="body2" fontWeight="600">Activities with Skills</Typography>
                  <Typography variant="caption" color="text.secondary">
                    Each lesson contains activities (beginner/intermediate/advanced), each tagged with the competencies it develops
                  </Typography>
                </Box>
              </Box>
              <Box sx={{ display: 'flex', gap: 2 }}>
                <Box sx={{ color: 'primary.main', fontWeight: 'bold', minWidth: 24 }}>3</Box>
                <Box>
                  <Typography variant="body2" fontWeight="600">Work Assignments</Typography>
                  <Typography variant="caption" color="text.secondary">
                    System automatically generates personalized work assignments for each student based on their competency profile and readiness
                  </Typography>
                </Box>
              </Box>
              <Box sx={{ display: 'flex', gap: 2 }}>
                <Box sx={{ color: 'primary.main', fontWeight: 'bold', minWidth: 24 }}>4</Box>
                <Box>
                  <Typography variant="body2" fontWeight="600">Marks & Tracking</Typography>
                  <Typography variant="caption" color="text.secondary">
                    Students complete their assignments, which are marked and tracked against their target competencies
                  </Typography>
                </Box>
              </Box>
            </Stack>
          </Box>

          <Box>
            <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
              Key Features to Try
            </Typography>
            <Stack spacing={1} sx={{ pl: 1 }}>
              <Typography variant="body2">
                • <strong>View Details</strong> – Click on any lesson to see how work assignments are automatically generated for each student
              </Typography>
              <Typography variant="body2">
                • <strong>Student Readiness</strong> – See how readiness scores (0-100%) are calculated from each student's competency profile
              </Typography>
              <Typography variant="body2">
                • <strong>Differentiation</strong> – Notice how struggling students get scaffolded activities while high achievers get challenges
              </Typography>
              <Typography variant="body2">
                • <strong>Skills Matching</strong> – Each assignment targets specific competencies for that individual student
              </Typography>
            </Stack>
          </Box>

          <Divider />

          <Box>
            <Typography variant="caption" color="text.secondary" display="block" sx={{ mb: 1 }}>
              <strong>Demo Data:</strong> The system includes 4 sample students with different competency profiles working through lessons in Biology, History, Computer Science, and Mathematics.
            </Typography>
            <Typography variant="caption" color="text.secondary">
              This demonstrates how lessons personalize to each student while covering the same curriculum material.
            </Typography>
          </Box>
        </Stack>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} variant="contained">Got It!</Button>
      </DialogActions>
    </Dialog>
  );
};

interface LessonPlanFormProps {
  open: boolean;
  plan: LessonPlan | null;
  onClose: () => void;
  onSave: (plan: LessonPlan) => void;
}

const LessonPlanForm: React.FC<LessonPlanFormProps> = ({ open, plan, onClose, onSave }) => {
  const [formData, setFormData] = useState<LessonPlan>(
    plan || {
      id: `lp-${Date.now()}`,
      title: '',
      subject: '',
      year: '',
      date: new Date().toISOString().split('T')[0],
      duration: '50',
      objectives: [''],
      starter: '',
      mainActivity: '',
      plenary: '',
      resources: [''],
      assessment: '',
      status: 'draft',
      createdBy: 'Current User',
      lastModified: new Date().toISOString().split('T')[0],
      targetCompetencies: [],
      studentAlignments: {},
      activities: [],
      workAssignments: [],
    }
  );

  const handleInputChange = (field: keyof LessonPlan, value: any) => {
    setFormData({ ...formData, [field]: value });
  };

  const handleArrayItemChange = (field: string, index: number, value: string) => {
    const arr = [...(formData[field as keyof LessonPlan] as string[])];
    arr[index] = value;
    handleInputChange(field as keyof LessonPlan, arr);
  };

  const handleSave = () => {
    onSave(formData);
    onClose();
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>{plan ? 'Edit Lesson Plan' : 'Create New Lesson Plan'}</DialogTitle>
      <DialogContent sx={{ pt: 2 }}>
        <Stack spacing={2}>
          <TextField
            fullWidth
            label="Lesson Title"
            value={formData.title}
            onChange={(e) => handleInputChange('title', e.target.value)}
          />

          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Subject</InputLabel>
                <Select value={formData.subject} onChange={(e) => handleInputChange('subject', e.target.value)} label="Subject">
                  <MenuItem value="Biology">Biology</MenuItem>
                  <MenuItem value="Chemistry">Chemistry</MenuItem>
                  <MenuItem value="History">History</MenuItem>
                  <MenuItem value="English">English</MenuItem>
                  <MenuItem value="Mathematics">Mathematics</MenuItem>
                  <MenuItem value="Computer Science">Computer Science</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Year Group"
                value={formData.year}
                onChange={(e) => handleInputChange('year', e.target.value)}
              />
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                type="date"
                label="Date"
                InputLabelProps={{ shrink: true }}
                value={formData.date}
                onChange={(e) => handleInputChange('date', e.target.value)}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                type="number"
                label="Duration (minutes)"
                value={formData.duration}
                onChange={(e) => handleInputChange('duration', e.target.value)}
              />
            </Grid>
          </Grid>

          <Box>
            <Typography variant="subtitle2" gutterBottom>Learning Objectives</Typography>
            <Stack spacing={1}>
              {(formData.objectives as string[]).map((obj, idx) => (
                <TextField
                  key={idx}
                  fullWidth
                  size="small"
                  placeholder={`Objective ${idx + 1}`}
                  value={obj}
                  onChange={(e) => handleArrayItemChange('objectives', idx, e.target.value)}
                />
              ))}
            </Stack>
          </Box>

          <TextField
            fullWidth
            multiline
            rows={2}
            label="Starter Activity"
            value={formData.starter}
            onChange={(e) => handleInputChange('starter', e.target.value)}
          />

          <TextField
            fullWidth
            multiline
            rows={3}
            label="Main Activity"
            value={formData.mainActivity}
            onChange={(e) => handleInputChange('mainActivity', e.target.value)}
          />

          <TextField
            fullWidth
            multiline
            rows={2}
            label="Plenary"
            value={formData.plenary}
            onChange={(e) => handleInputChange('plenary', e.target.value)}
          />

          <TextField
            fullWidth
            multiline
            rows={2}
            label="Assessment"
            value={formData.assessment}
            onChange={(e) => handleInputChange('assessment', e.target.value)}
          />
        </Stack>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button onClick={handleSave} variant="contained">Save Lesson Plan</Button>
      </DialogActions>
    </Dialog>
  );
};

// ─── Main Demo Component ──────────────────────────────────
const LessonPlanningDemo: React.FC = () => {
  const theme = useTheme();
  const [plans, setPlans] = useState<LessonPlan[]>(dummyLessonPlans);
  const [formOpen, setFormOpen] = useState(false);
  const [selectedPlan, setSelectedPlan] = useState<LessonPlan | null>(null);
  const [detailOpen, setDetailOpen] = useState(false);
  const [detailedPlan, setDetailedPlan] = useState<LessonPlan | null>(null);
  const [filterSubject, setFilterSubject] = useState<string>('');
  const [helpOpen, setHelpOpen] = useState(false);

  const handleEdit = (plan: LessonPlan) => {
    setSelectedPlan(plan);
    setFormOpen(true);
  };

  const handleDelete = (id: string) => {
    setPlans(plans.filter((p) => p.id !== id));
  };

  const handleSave = (plan: LessonPlan) => {
    if (selectedPlan) {
      setPlans(plans.map((p) => (p.id === plan.id ? plan : p)));
    } else {
      setPlans([...plans, plan]);
    }
    setSelectedPlan(null);
  };

  const filteredPlans = filterSubject ? plans.filter((p) => p.subject === filterSubject) : plans;

  const subjects = Array.from(new Set(dummyLessonPlans.map((p) => p.subject)));

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
            <Box sx={{ flexGrow: 1 }}>
              <Typography variant="h5" gutterBottom fontWeight="bold">
                Lesson Planning Demo
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Create, manage, and deliver structured lesson plans with clear objectives and assessment criteria
              </Typography>
            </Box>
            <IconButton
              onClick={() => setHelpOpen(true)}
              sx={{ color: 'primary.main', '&:hover': { bgcolor: alpha(theme.palette.primary.main, 0.1) } }}
              title="Learn how this works"
            >
              <HelpIcon />
            </IconButton>
          </Stack>
        </Container>
      </Box>

      <Container maxWidth="lg">
        {/* ── Overview Stats ── */}
        <Grid container spacing={3} sx={{ mb: 4 }}>
          <Grid item xs={12} sm={6} md={3}>
            <Paper sx={{ p: 3, textAlign: 'center', bgcolor: alpha(theme.palette.primary.main, 0.05) }}>
              <Typography variant="h4" fontWeight="bold" color="primary">
                {plans.length}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Lesson Plans
              </Typography>
            </Paper>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Paper sx={{ p: 3, textAlign: 'center', bgcolor: alpha('#2e7d32', 0.05) }}>
              <Typography variant="h4" fontWeight="bold" sx={{ color: '#2e7d32' }}>
                {plans.filter((p) => p.status === 'published').length}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Published
              </Typography>
            </Paper>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Paper sx={{ p: 3, textAlign: 'center', bgcolor: alpha('#f57c00', 0.05) }}>
              <Typography variant="h4" fontWeight="bold" sx={{ color: '#f57c00' }}>
                {plans.filter((p) => p.status === 'draft').length}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Drafts
              </Typography>
            </Paper>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Paper sx={{ p: 3, textAlign: 'center', bgcolor: alpha('#0288d1', 0.05) }}>
              <Typography variant="h4" fontWeight="bold" sx={{ color: '#0288d1' }}>
                {plans.filter((p) => p.status === 'completed').length}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Delivered
              </Typography>
            </Paper>
          </Grid>
        </Grid>

        {/* ── Filtering & Actions ── */}
        <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} sx={{ mb: 4, alignItems: 'center' }}>
          <FormControl sx={{ minWidth: 200 }}>
            <InputLabel>Filter by Subject</InputLabel>
            <Select value={filterSubject} onChange={(e) => setFilterSubject(e.target.value)} label="Filter by Subject">
              <MenuItem value="">All Subjects</MenuItem>
              {subjects.map((subject) => (
                <MenuItem key={subject} value={subject}>
                  {subject}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <Box sx={{ flexGrow: 1 }} />
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => {
              setSelectedPlan(null);
              setFormOpen(true);
            }}
          >
            New Lesson Plan
          </Button>
        </Stack>

        {/* ── Lesson Plan Grid ── */}
        <Grid container spacing={3} sx={{ mb: 6 }}>
          {filteredPlans.map((plan) => (
            <Grid item xs={12} md={6} lg={4} key={plan.id}>
              <LessonCard plan={plan} onEdit={handleEdit} onDelete={handleDelete} onView={(p) => { setDetailedPlan(p); setDetailOpen(true); }} />
            </Grid>
          ))}
        </Grid>

        {/* ── Form Dialog ── */}
        <LessonPlanForm open={formOpen} plan={selectedPlan} onClose={() => { setFormOpen(false); setSelectedPlan(null); }} onSave={handleSave} />

        {/* ── Detail Dialog ── */}
        <LessonDetailDialog open={detailOpen} plan={detailedPlan} onClose={() => { setDetailOpen(false); setDetailedPlan(null); }} />

        {/* ── Help Dialog ── */}
        <HelpDialog open={helpOpen} onClose={() => setHelpOpen(false)} />

        {/* ── Footer Information ── */}
        <Paper elevation={0} sx={{ p: 4, bgcolor: alpha(theme.palette.primary.main, 0.05), borderRadius: 2 }}>
          <Stack spacing={2}>
            <Typography variant="h6" fontWeight="bold">
              About the Lesson Planning System
            </Typography>
            <Typography variant="body2" color="text.secondary">
              This demonstration showcases a modern lesson planning interface that combines the original <strong>LessonPlan.php</strong> API logic with intelligent skills-to-competency matching for personalized learning.
            </Typography>
            <Stack spacing={1}>
              <Typography variant="body2">
                <strong>Core Features:</strong>
              </Typography>
              <List dense>
                <ListItem>
                  <ListItemIcon>
                    <CheckIcon sx={{ fontSize: 20 }} />
                  </ListItemIcon>
                  <ListItemText primary="Structured lesson planning with clear objectives, activities, and assessments" />
                </ListItem>
                <ListItem>
                  <ListItemIcon>
                    <CheckIcon sx={{ fontSize: 20 }} />
                  </ListItemIcon>
                  <ListItemText primary="Draft, publish, and completion workflow for lesson delivery" />
                </ListItem>
                <ListItem>
                  <ListItemIcon>
                    <CheckIcon sx={{ fontSize: 20 }} />
                  </ListItemIcon>
                  <ListItemText primary="Reusable lesson library searchable by subject and year group" />
                </ListItem>
                <ListItem>
                  <ListItemIcon>
                    <CheckIcon sx={{ fontSize: 20 }} />
                  </ListItemIcon>
                  <ListItemText primary="Skills-to-competency mapping: Align lesson objectives with student competencies to personalize learning pathways" />
                </ListItem>
                <ListItem>
                  <ListItemIcon>
                    <CheckIcon sx={{ fontSize: 20 }} />
                  </ListItemIcon>
                  <ListItemText primary="Readiness scoring: Automatic calculation of student readiness for each lesson based on current competency levels" />
                </ListItem>
                <ListItem>
                  <ListItemIcon>
                    <CheckIcon sx={{ fontSize: 20 }} />
                  </ListItemIcon>
                  <ListItemText primary="Differentiation guidance: Actionable tailored recommendations for each student's learning needs and development opportunities" />
                </ListItem>
                <ListItem>
                  <ListItemIcon>
                    <CheckIcon sx={{ fontSize: 20 }} />
                  </ListItemIcon>
                  <ListItemText primary="Work assignment generation: Create personalized activity assignments from lesson activities, matching each student's competency profile and readiness level" />
                </ListItem>
              </List>
            </Stack>
          </Stack>
        </Paper>
      </Container>
    </Box>
  );
};

export default LessonPlanningDemo;
