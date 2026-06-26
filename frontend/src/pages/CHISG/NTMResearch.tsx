import React, { useCallback, useEffect, useRef, useState } from 'react';
import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  Container,
  Divider,
  Paper,
  Stack,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Tabs,
  TextField,
  Tooltip,
  Typography,
  alpha,
} from '@mui/material';
import {
  CloudUpload as UploadIcon,
  Psychology as QueryIcon,
  Science as ScienceIcon,
  CheckCircle as CheckIcon,
  Article as PaperIcon,
  Send as SendIcon,
  AccountTree as AccountTreeIcon,
  Edit as EditIcon,
  ListAlt as ListAltIcon,
  Lock as LockIcon,
  PlayArrow as RunIcon,
  HourglassEmpty as HourglassIcon,
} from '@mui/icons-material';
import { Link, useNavigate } from 'react-router-dom';
import SEO from '../../components/SEO';
import { apiClient, API_ENDPOINTS } from '../../config/api';
import { useAuth } from '../../context/AuthContext';

// ── Types ─────────────────────────────────────────────────────────────────────

interface SourceLink {
  entity_a: string;
  relation: string;
  entity_b: string;
  context?: string;
}

interface QueryResponse {
  answer: string;
  sources: SourceLink[];
  links_used: number;
}

interface UploadResponse {
  job_id: string;
  status: string;
  filename: string;
  message: string;
}

// ── Colour tokens ─────────────────────────────────────────────────────────────

const PRIMARY = '#0f172a';
const ACCENT = '#4f46e5';
const TEAL = '#0d9488';
const SURFACE = '#f8fafc';

// ── Nav bar ───────────────────────────────────────────────────────────────────

const NTMNav: React.FC = () => (
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
        <Box
          component="img"
          src={`${process.env.PUBLIC_URL}/ESPLogoLong.png`}
          alt="ESP Thinking"
          sx={{ height: 48, width: 'auto' }}
        />
        <Stack direction="row" spacing={1}>
          <Button
            component={Link}
            to="/chisg"
            variant="outlined"
            size="small"
            sx={{
              color: 'white',
              borderColor: 'rgba(255,255,255,0.5)',
              textTransform: 'none',
              bgcolor: 'rgba(15,23,42,0.3)',
              '&:hover': { borderColor: 'white', bgcolor: 'rgba(148,163,184,0.3)' },
            }}
          >
            CHISG
          </Button>
          <Button
            component={Link}
            to="/"
            variant="outlined"
            size="small"
            sx={{
              color: 'white',
              borderColor: 'rgba(255,255,255,0.5)',
              textTransform: 'none',
              bgcolor: 'rgba(15,23,42,0.3)',
              '&:hover': { borderColor: 'white', bgcolor: 'rgba(148,163,184,0.3)' },
            }}
          >
            Portfolio
          </Button>
        </Stack>
      </Stack>
    </Container>
  </Box>
);

// ── Upload tab ────────────────────────────────────────────────────────────────

const UploadTab: React.FC = () => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [file, setFile] = useState<File | null>(null);
  const [title, setTitle] = useState('');
  const [uploading, setUploading] = useState(false);
  const [result, setResult] = useState<UploadResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  const [extracting, setExtracting] = useState(false);
  const [extractDone, setExtractDone] = useState(false);
  const [extractLinks, setExtractLinks] = useState(0);
  const [extractError, setExtractError] = useState<string | null>(null);
  const pollRef = useRef<ReturnType<typeof setInterval> | null>(null);

  // Cleanup polling on unmount
  useEffect(() => () => { if (pollRef.current) clearInterval(pollRef.current); }, []);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const f = e.target.files?.[0] ?? null;
    setFile(f);
    setResult(null);
    setError(null);
    if (f && !title) setTitle(f.name.replace(/\.pdf$/i, ''));
  };

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    const f = e.dataTransfer.files[0];
    if (f?.type === 'application/pdf') {
      setFile(f);
      setResult(null);
      setError(null);
      if (!title) setTitle(f.name.replace(/\.pdf$/i, ''));
    }
  }, [title]);

  const handleSubmit = async () => {
    if (!file) return;
    setUploading(true);
    setError(null);
    setExtractDone(false);
    setExtractError(null);
    try {
      const fd = new FormData();
      fd.append('paper', file);
      fd.append('title', title || file.name);
      const res = await apiClient.upload<UploadResponse>(API_ENDPOINTS.ntm.upload, fd);
      setResult(res);
    } catch (err: any) {
      setError(err?.response?.data?.error ?? 'Upload failed — is the API running?');
    } finally {
      setUploading(false);
    }
  };

  const handleRunExtraction = async () => {
    if (!result) return;
    setExtracting(true);
    setExtractError(null);
    setExtractDone(false);
    try {
      await apiClient.post(API_ENDPOINTS.ntm.extract, { job_id: result.job_id });
      // Poll for status every 4s until complete or error
      pollRef.current = setInterval(async () => {
        try {
          const st = await apiClient.get<{ status: string; link_count: number }>(
            `${API_ENDPOINTS.ntm.status}/${result.job_id}`
          );
          if (st.status === 'complete') {
            clearInterval(pollRef.current!);
            pollRef.current = null;
            setExtracting(false);
            setExtractDone(true);
            setExtractLinks(st.link_count ?? 0);
          } else if (st.status.startsWith('error')) {
            clearInterval(pollRef.current!);
            pollRef.current = null;
            setExtracting(false);
            setExtractError(`Extraction failed: ${st.status.replace('error: ', '')}`);
          }
        } catch {
          // non-fatal polling error — keep trying
        }
      }, 4000);
    } catch (err: any) {
      setExtracting(false);
      setExtractError(err?.response?.data?.error ?? 'Failed to start extraction');
    }
  };

  return (
    <Stack spacing={3}>
      <Typography variant="body1" color="text.secondary">
        Upload a PDF academic paper to queue it for CHISG extraction. The extraction
        pipeline reads the paper, identifies semantic relationships between biological
        entities, and ingests them into the knowledge graph.
      </Typography>

      {/* Drop zone */}
      <Paper
        onDragOver={(e) => e.preventDefault()}
        onDrop={handleDrop}
        onClick={() => fileInputRef.current?.click()}
        sx={{
          border: `2px dashed ${file ? TEAL : alpha(ACCENT, 0.4)}`,
          borderRadius: 2,
          p: 5,
          textAlign: 'center',
          cursor: 'pointer',
          bgcolor: file ? alpha(TEAL, 0.04) : 'transparent',
          transition: 'all 0.2s',
          '&:hover': { bgcolor: alpha(ACCENT, 0.04), borderColor: ACCENT },
        }}
      >
        <input
          ref={fileInputRef}
          type="file"
          accept=".pdf"
          style={{ display: 'none' }}
          onChange={handleFileChange}
        />
        {file ? (
          <Stack alignItems="center" spacing={1}>
            <PaperIcon sx={{ fontSize: 40, color: TEAL }} />
            <Typography fontWeight={600}>{file.name}</Typography>
            <Typography variant="body2" color="text.secondary">
              {(file.size / 1024 / 1024).toFixed(2)} MB — click to change
            </Typography>
          </Stack>
        ) : (
          <Stack alignItems="center" spacing={1}>
            <UploadIcon sx={{ fontSize: 40, color: alpha(ACCENT, 0.5) }} />
            <Typography fontWeight={500}>Drop a PDF here or click to browse</Typography>
            <Typography variant="body2" color="text.secondary">
              Academic papers only · max 50 MB
            </Typography>
          </Stack>
        )}
      </Paper>

      {file && (
        <TextField
          label="Paper title (optional)"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          fullWidth
          size="small"
          placeholder="Auto-filled from filename"
        />
      )}

      <Button
        variant="contained"
        disabled={!file || uploading}
        onClick={handleSubmit}
        startIcon={uploading ? <CircularProgress size={16} color="inherit" /> : <UploadIcon />}
        sx={{ bgcolor: ACCENT, '&:hover': { bgcolor: '#4338ca' }, alignSelf: 'flex-start' }}
      >
        {uploading ? 'Uploading…' : 'Upload Paper'}
      </Button>

      {error && <Alert severity="error">{error}</Alert>}

      {result && (
        <Alert
          severity="success"
          icon={<CheckIcon />}
          sx={{ bgcolor: alpha(TEAL, 0.08), border: `1px solid ${alpha(TEAL, 0.3)}` }}
        >
          <Typography fontWeight={600}>{result.filename} queued</Typography>
          <Typography variant="body2">{result.message}</Typography>
          <Typography variant="caption" color="text.secondary" sx={{ mt: 0.5, display: 'block' }}>
            Job ID: {result.job_id}
          </Typography>
        </Alert>
      )}

      {/* Extraction controls — shown after a successful upload */}
      {result && !extractDone && (
        <Stack direction="row" alignItems="center" spacing={2}>
          <Button
            variant="contained"
            disabled={extracting}
            onClick={handleRunExtraction}
            startIcon={extracting ? <CircularProgress size={16} color="inherit" /> : <RunIcon />}
            sx={{ bgcolor: TEAL, '&:hover': { bgcolor: '#0b7a72' } }}
          >
            {extracting ? 'Extracting links…' : 'Run CHISG Extraction'}
          </Button>
          {extracting && (
            <Stack direction="row" alignItems="center" spacing={1}>
              <HourglassIcon sx={{ fontSize: 14, color: 'text.secondary' }} />
              <Typography variant="caption" color="text.secondary">
                Processing chunks — typically 1–3 minutes
              </Typography>
            </Stack>
          )}
        </Stack>
      )}

      {extractError && <Alert severity="error">{extractError}</Alert>}

      {extractDone && (
        <Alert
          severity="success"
          icon={<CheckIcon />}
          sx={{ bgcolor: alpha(TEAL, 0.08), border: `1px solid ${alpha(TEAL, 0.3)}` }}
        >
          <Typography fontWeight={600}>
            Extraction complete — {extractLinks} semantic link{extractLinks !== 1 ? 's' : ''} added to knowledge graph
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Links are now searchable via the Query tab.
          </Typography>
        </Alert>
      )}

      <Divider />

      <Box>
        <Typography variant="subtitle2" fontWeight={700} gutterBottom>
          Extraction pipeline
        </Typography>
        <Typography variant="body2" color="text.secondary">
          After uploading, run the extraction pipeline from the server:
        </Typography>
        <Paper sx={{ p: 1.5, mt: 1, bgcolor: '#0f172a', borderRadius: 1 }}>
          <Typography
            component="code"
            sx={{ fontFamily: 'monospace', fontSize: 12, color: '#86efac', display: 'block' }}
          >
            python3 scripts/extract_chisg_papers.py --input /data/ntm-papers/&lt;filename&gt;
          </Typography>
        </Paper>
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
          Then ingest the output JSON into Weaviate:
        </Typography>
        <Paper sx={{ p: 1.5, mt: 1, bgcolor: '#0f172a', borderRadius: 1 }}>
          <Typography
            component="code"
            sx={{ fontFamily: 'monospace', fontSize: 12, color: '#86efac', display: 'block' }}
          >
            python3 scripts/ingest_chisg_academic.py --ingest --input &lt;output.json&gt;
          </Typography>
        </Paper>
      </Box>
    </Stack>
  );
};

// ── Query tab ─────────────────────────────────────────────────────────────────

const QueryTab: React.FC = () => {
  const [question, setQuestion] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<QueryResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  const EXAMPLE_QUESTIONS = [
    'How does CarD regulate mycobacterial starvation response?',
    'What is the role of Clp protease in NTM survival?',
    'How does (p)ppGpp affect gene regulation under nutrient stress?',
    'What mechanisms does Mycobacterium use to survive in low-nutrient conditions?',
  ];

  const handleSubmit = async () => {
    if (!question.trim()) return;
    setLoading(true);
    setError(null);
    setResult(null);
    try {
      const res = await apiClient.post<QueryResponse>(API_ENDPOINTS.ntm.query, { question });
      setResult(res);
    } catch (err: any) {
      setError(err?.response?.data?.error ?? 'Query failed — is the API running?');
    } finally {
      setLoading(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) handleSubmit();
  };

  return (
    <Stack spacing={3}>
      <Typography variant="body1" color="text.secondary">
        Ask a research question about Non-Tuberculous Mycobacteria starvation survival.
        The system retrieves relevant semantic links from the knowledge graph and uses
        Claude (AWS Bedrock) to synthesise an evidence-based answer.
      </Typography>

      {/* Example chips */}
      <Box>
        <Typography variant="caption" color="text.secondary" sx={{ mb: 1, display: 'block' }}>
          Example questions:
        </Typography>
        <Stack direction="row" flexWrap="wrap" gap={1}>
          {EXAMPLE_QUESTIONS.map((q) => (
            <Chip
              key={q}
              label={q}
              size="small"
              clickable
              onClick={() => setQuestion(q)}
              sx={{
                bgcolor: alpha(ACCENT, 0.08),
                '&:hover': { bgcolor: alpha(ACCENT, 0.16) },
                fontSize: 11,
              }}
            />
          ))}
        </Stack>
      </Box>

      <TextField
        multiline
        minRows={3}
        maxRows={6}
        label="Research question"
        value={question}
        onChange={(e) => setQuestion(e.target.value)}
        onKeyDown={handleKeyDown}
        fullWidth
        placeholder="e.g. What is the role of CarD in mycobacterial starvation response?"
        helperText="⌘+Enter to submit"
      />

      <Button
        variant="contained"
        disabled={!question.trim() || loading}
        onClick={handleSubmit}
        startIcon={loading ? <CircularProgress size={16} color="inherit" /> : <SendIcon />}
        sx={{ bgcolor: ACCENT, '&:hover': { bgcolor: '#4338ca' }, alignSelf: 'flex-start' }}
      >
        {loading ? 'Searching knowledge graph…' : 'Ask'}
      </Button>

      {error && <Alert severity="error">{error}</Alert>}

      {result && (
        <Stack spacing={2}>
          {/* Answer */}
          <Card
            elevation={0}
            sx={{
              border: `1px solid ${alpha(TEAL, 0.3)}`,
              bgcolor: alpha(TEAL, 0.04),
              borderRadius: 2,
            }}
          >
            <CardContent>
              <Stack direction="row" alignItems="center" spacing={1} sx={{ mb: 1.5 }}>
                <QueryIcon sx={{ color: TEAL, fontSize: 18 }} />
                <Typography variant="subtitle2" fontWeight={700} color={TEAL}>
                  Answer
                </Typography>
                <Chip
                  label={`${result.links_used} links used`}
                  size="small"
                  sx={{ ml: 'auto', bgcolor: alpha(TEAL, 0.1) }}
                />
              </Stack>
              <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap', lineHeight: 1.8 }}>
                {result.answer}
              </Typography>
            </CardContent>
          </Card>

          {/* Sources */}
          {result.sources && result.sources.length > 0 && (
            <Box>
              <Typography variant="subtitle2" fontWeight={700} sx={{ mb: 1.5 }}>
                Knowledge graph sources ({result.sources.length})
              </Typography>
              <Stack spacing={1}>
                {result.sources.slice(0, 15).map((s, i) => (
                  <Paper
                    key={i}
                    elevation={0}
                    sx={{
                      p: 1.5,
                      border: `1px solid ${alpha('#94a3b8', 0.2)}`,
                      borderRadius: 1,
                      bgcolor: SURFACE,
                    }}
                  >
                    <Stack direction="row" flexWrap="wrap" alignItems="center" gap={0.75}>
                      <Chip
                        label={s.entity_a}
                        size="small"
                        sx={{ bgcolor: alpha(ACCENT, 0.1), fontWeight: 600, fontSize: 11 }}
                      />
                      <Typography variant="caption" color="text.secondary" sx={{ fontStyle: 'italic' }}>
                        {s.relation}
                      </Typography>
                      <Chip
                        label={s.entity_b}
                        size="small"
                        sx={{ bgcolor: alpha(TEAL, 0.1), fontWeight: 600, fontSize: 11 }}
                      />
                    </Stack>
                    {s.context && (
                      <Typography
                        variant="caption"
                        color="text.secondary"
                        sx={{ mt: 0.5, display: 'block', pl: 0.5, borderLeft: `2px solid ${alpha(ACCENT, 0.3)}` }}
                      >
                        {s.context}
                      </Typography>
                    )}
                  </Paper>
                ))}
                {result.sources.length > 15 && (
                  <Typography variant="caption" color="text.secondary">
                    + {result.sources.length - 15} more links
                  </Typography>
                )}
              </Stack>
            </Box>
          )}
        </Stack>
      )}
    </Stack>
  );
};

// ── Papers tab ────────────────────────────────────────────────────────────────

interface PaperItem {
  paper_id: string;
  title: string;
  filename: string;
  pages: number;
  total_chars: number;
  link_count: number;
  entity_count: number;
}

const PapersTab: React.FC = () => {
  const { token, isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [papers, setPapers] = useState<PaperItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isAuthenticated) return;
    setLoading(true);
    fetch('/api/chisg/papers', { headers: { Authorization: `Bearer ${token}` } })
      .then((r) => r.json())
      .then((d) => setPapers(d.papers || []))
      .catch(() => setError('Failed to load papers.'))
      .finally(() => setLoading(false));
  }, [isAuthenticated, token]);

  if (!isAuthenticated) {
    return (
      <Stack alignItems="center" spacing={2} sx={{ py: 6 }}>
        <LockIcon sx={{ fontSize: 48, color: alpha(ACCENT, 0.4) }} />
        <Typography variant="h6" fontWeight={600}>Sign in to view processed papers</Typography>
        <Typography variant="body2" color="text.secondary" textAlign="center" maxWidth={400}>
          The papers list is restricted to authorised researchers. Sign in with your CHISG admin credentials.
        </Typography>
        <Button
          variant="contained"
          onClick={() => navigate('/chisg/login')}
          sx={{ bgcolor: ACCENT, '&:hover': { bgcolor: '#4338ca' } }}
        >
          Sign in
        </Button>
      </Stack>
    );
  }

  if (loading) return <Box sx={{ py: 6, textAlign: 'center' }}><CircularProgress /></Box>;
  if (error) return <Alert severity="error">{error}</Alert>;

  if (papers.length === 0) {
    return (
      <Alert severity="info" action={
        <Button size="small" onClick={() => navigate('/ntm')}>Upload one</Button>
      }>
        No papers have been processed yet.
      </Alert>
    );
  }

  return (
    <TableContainer component={Paper} elevation={0} sx={{ border: `1px solid ${alpha('#94a3b8', 0.15)}`, borderRadius: 1 }}>
      <Table size="small">
        <TableHead>
          <TableRow sx={{ bgcolor: alpha(PRIMARY, 0.04) }}>
            <TableCell><strong>Title</strong></TableCell>
            <TableCell><strong>Paper ID</strong></TableCell>
            <TableCell align="center"><strong>Links</strong></TableCell>
            <TableCell align="center"><strong>Entities</strong></TableCell>
            <TableCell align="center"><strong>Pages</strong></TableCell>
            <TableCell align="center"><strong>Actions</strong></TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {papers.map((p) => (
            <TableRow key={p.paper_id} hover>
              <TableCell>
                <Typography variant="body2" fontWeight={500}>
                  {p.title || p.filename || p.paper_id}
                </Typography>
                {p.filename && p.title && (
                  <Typography variant="caption" color="text.secondary">{p.filename}</Typography>
                )}
              </TableCell>
              <TableCell>
                <Chip label={p.paper_id} size="small" variant="outlined" sx={{ fontSize: 11 }} />
              </TableCell>
              <TableCell align="center">
                <Chip label={p.link_count} size="small" color="primary" variant="filled" />
              </TableCell>
              <TableCell align="center">{p.entity_count}</TableCell>
              <TableCell align="center">{p.pages || '—'}</TableCell>
              <TableCell align="center">
                <Box display="flex" gap={1} justifyContent="center">
                  <Tooltip title="View in graph">
                    <Button size="small" variant="outlined" startIcon={<AccountTreeIcon />}
                      onClick={() => navigate(`/chisg/graph?paper_id=${p.paper_id}`)}>
                      Graph
                    </Button>
                  </Tooltip>
                  <Tooltip title="Review / edit links">
                    <Button size="small" variant="contained" startIcon={<EditIcon />}
                      sx={{ bgcolor: ACCENT, '&:hover': { bgcolor: '#4338ca' } }}
                      onClick={() => navigate(`/chisg/links?paper_id=${p.paper_id}`)}>
                      Links
                    </Button>
                  </Tooltip>
                </Box>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

// ── Main page ─────────────────────────────────────────────────────────────────

const NTMResearch: React.FC = () => {
  const [tab, setTab] = useState(0);

  return (
    <Box sx={{ minHeight: '100vh', bgcolor: SURFACE }}>
      <SEO
        title="NTM Research — CHISG Knowledge Graph | ESP Thinking"
        description="Academic knowledge graph for Non-Tuberculous Mycobacteria starvation survival research. Upload papers, query the CHISG semantic network using AI."
        path="/ntm"
      />

      {/* Hero */}
      <Box
        sx={{
          position: 'relative',
          background: `linear-gradient(135deg, ${PRIMARY} 0%, #1e293b 50%, ${TEAL} 100%)`,
          color: 'white',
          pt: { xs: 10, md: 13 },
          pb: { xs: 6, md: 8 },
          overflow: 'hidden',
        }}
      >
        <NTMNav />
        <Container maxWidth="lg">
          <Stack spacing={2}>
            <Stack direction="row" spacing={1} alignItems="center">
              <ScienceIcon sx={{ fontSize: 20, opacity: 0.8 }} />
              <Typography variant="overline" sx={{ opacity: 0.8, letterSpacing: 2 }}>
                CHISG Academic Repository
              </Typography>
            </Stack>
            <Typography variant="h3" fontWeight={800} sx={{ lineHeight: 1.2 }}>
              NTM Starvation Survival
            </Typography>
            <Typography variant="h6" sx={{ opacity: 0.8, fontWeight: 400, maxWidth: 600 }}>
              Knowledge graph for{' '}
              <em>
                "Understanding the Global Control of Mechanisms for Starvation Survival in
                Non-Tuberculous Mycobacteria"
              </em>
            </Typography>
            <Stack direction="row" spacing={1} flexWrap="wrap" sx={{ mt: 1 }}>
              {[
                '454 semantic links',
                '647 biological entities',
                '30 relation types',
                'Claude via Bedrock',
              ].map((tag) => (
                <Chip
                  key={tag}
                  label={tag}
                  size="small"
                  sx={{
                    bgcolor: 'rgba(255,255,255,0.12)',
                    color: 'white',
                    border: '1px solid rgba(255,255,255,0.2)',
                    fontSize: 11,
                  }}
                />
              ))}
            </Stack>
          </Stack>
        </Container>
      </Box>

      {/* Tabs */}
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Paper elevation={0} sx={{ border: `1px solid ${alpha('#94a3b8', 0.2)}`, borderRadius: 2 }}>
          <Tabs
            value={tab}
            onChange={(_, v) => setTab(v)}
            sx={{
              px: 2,
              borderBottom: `1px solid ${alpha('#94a3b8', 0.2)}`,
              '& .MuiTab-root': { textTransform: 'none', fontWeight: 600 },
            }}
          >
            <Tab
              icon={<UploadIcon sx={{ fontSize: 18 }} />}
              iconPosition="start"
              label="Upload Paper"
            />
            <Tab
              icon={<QueryIcon sx={{ fontSize: 18 }} />}
              iconPosition="start"
              label="Ask a Question"
            />
            <Tab
              icon={<ListAltIcon sx={{ fontSize: 18 }} />}
              iconPosition="start"
              label="Processed Papers"
            />
          </Tabs>

          <Box sx={{ p: 3 }}>
            {tab === 0 && <UploadTab />}
            {tab === 1 && <QueryTab />}
            {tab === 2 && <PapersTab />}
          </Box>
        </Paper>

        {/* Footer note */}
        <Typography
          variant="caption"
          color="text.disabled"
          sx={{ display: 'block', textAlign: 'center', mt: 3 }}
        >
          Research collaboration with Michael McGrath, University of Brighton ·{' '}
          <Link to="/chisg" style={{ color: 'inherit' }}>
            CHISG methodology
          </Link>
        </Typography>
      </Container>
    </Box>
  );
};

export default NTMResearch;
