import React, { useEffect, useState } from 'react';
import { useNavigate, useSearchParams, useLocation } from 'react-router-dom';
import {
  Box,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Button,
  Chip,
  CircularProgress,
  Alert,
  AppBar,
  Toolbar,
  IconButton,
  Tooltip,
  TextField,
} from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import LogoutIcon from '@mui/icons-material/Logout';
import { useAuth } from '../../context/AuthContext';

const API_BASE = process.env.REACT_APP_API_URL || '/api';

interface AcademicLink {
  entity_a: string;
  relation: string;
  entity_b: string;
  backward_relation: string;
  statement: string;
  context: string;
  paper_id: string;
  chunk_id: string;
  page_num: number;
}

const CHISGLinks: React.FC = () => {
  const { token, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [searchParams] = useSearchParams();
  const paperID = searchParams.get('paper_id') || '';

  // Paper metadata passed from CHISGPapers via router state
  const paperFromState = (location.state as any)?.paper as { paper_id: string; title?: string; filename?: string } | undefined;
  const paperTitle = paperFromState?.title || paperFromState?.filename || paperID;

  const [links, setLinks] = useState<AcademicLink[]>([]);
  const [filtered, setFiltered] = useState<AcademicLink[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');

  // Redirect to papers list if no paper_id — prevents showing mixed links
  useEffect(() => {
    if (!paperID) {
      navigate('/chisg/papers', { replace: true });
    }
  }, [paperID, navigate]);

  useEffect(() => {
    if (!paperID) return;
    const load = async () => {
      try {
        const url = `${API_BASE}/chisg/links?paper_id=${encodeURIComponent(paperID)}&limit=1000`;
        const res = await fetch(url, {
          headers: { Authorization: `Bearer ${token}` },
        });
        if (res.status === 401) {
          logout();
          navigate('/chisg/login');
          return;
        }
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const data = await res.json();
        setLinks(data.links || []);
        setFiltered(data.links || []);
      } catch (e: any) {
        setError(e.message || 'Failed to load links.');
      } finally {
        setLoading(false);
      }
    };
    load();
  }, [token, paperID, logout, navigate]);

  useEffect(() => {
    if (!search.trim()) {
      setFiltered(links);
      return;
    }
    const q = search.toLowerCase();
    setFiltered(
      links.filter(
        (l) =>
          l.entity_a.toLowerCase().includes(q) ||
          l.entity_b.toLowerCase().includes(q) ||
          l.relation.toLowerCase().includes(q) ||
          l.statement.toLowerCase().includes(q)
      )
    );
  }, [search, links]);

  const handleLogout = () => {
    logout();
    navigate('/chisg/login');
  };

  return (
    <Box sx={{ minHeight: '100vh', bgcolor: '#f5f7fa' }}>
      <AppBar position="static" elevation={1} sx={{ bgcolor: '#1a365d' }}>
        <Toolbar>
          <Tooltip title="Back to papers">
            <IconButton color="inherit" onClick={() => navigate('/chisg/papers')} sx={{ mr: 1 }}>
              <ArrowBackIcon />
            </IconButton>
          </Tooltip>
          <Typography variant="h6" fontWeight={700} sx={{ flexGrow: 1 }}>
            CHISG Knowledge Base — Links
          </Typography>
          <Tooltip title="Logout">
            <IconButton color="inherit" onClick={handleLogout}>
              <LogoutIcon />
            </IconButton>
          </Tooltip>
        </Toolbar>
      </AppBar>

      <Box sx={{ p: 4 }}>
        <Typography variant="h5" fontWeight={700} mb={0.5}>
          Academic Links
          {paperID && (
            <Chip label={paperID} size="small" sx={{ ml: 1.5, verticalAlign: 'middle' }} />
          )}
        </Typography>
        <Typography variant="body2" color="text.secondary" mb={3}>
          Approved semantic triples extracted from the paper. {filtered.length} of {links.length} shown.
        </Typography>

        <Box display="flex" gap={2} mb={3} alignItems="center">
          <TextField
            size="small"
            placeholder="Search entities or relations…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            sx={{ width: 340 }}
          />
          <Button
            variant="outlined"
            onClick={() =>
              navigate(`/semantic-links/extract`, {
                state: {
                  chisgPaper: {
                    id: paperID,
                    _id: paperID,
                    title: paperTitle,
                    domain: 'academic',
                  },
                },
              })
            }
          >
            Open Link Builder
          </Button>
        </Box>

        {loading && (
          <Box display="flex" justifyContent="center" py={8}>
            <CircularProgress />
          </Box>
        )}

        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

        {!loading && !error && filtered.length === 0 && (
          <Alert severity="info">
            {links.length === 0
              ? 'No links found for this paper.'
              : 'No links match your search.'}
          </Alert>
        )}

        {!loading && filtered.length > 0 && (
          <TableContainer component={Paper} elevation={2}>
            <Table size="small">
              <TableHead>
                <TableRow sx={{ bgcolor: '#f0f4f8' }}>
                  <TableCell><strong>Entity A</strong></TableCell>
                  <TableCell><strong>Relation</strong></TableCell>
                  <TableCell><strong>Entity B</strong></TableCell>
                  <TableCell><strong>Statement</strong></TableCell>
                  <TableCell align="center"><strong>Page</strong></TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {filtered.map((link, i) => (
                  <TableRow key={i} hover>
                    <TableCell>
                      <Typography variant="body2" fontWeight={500}>
                        {link.entity_a}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={link.relation}
                        size="small"
                        color="primary"
                        variant="outlined"
                        sx={{ fontSize: '0.7rem' }}
                      />
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" fontWeight={500}>
                        {link.entity_b}
                      </Typography>
                    </TableCell>
                    <TableCell sx={{ maxWidth: 400 }}>
                      <Typography variant="caption" color="text.secondary" sx={{ display: 'block' }}>
                        {link.statement}
                      </Typography>
                    </TableCell>
                    <TableCell align="center">
                      {link.page_num > 0 ? link.page_num : '—'}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </Box>
    </Box>
  );
};

export default CHISGLinks;
