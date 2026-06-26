import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
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
} from '@mui/material';
import LogoutIcon from '@mui/icons-material/Logout';
import AccountTreeIcon from '@mui/icons-material/AccountTree';
import EditIcon from '@mui/icons-material/Edit';
import { useAuth } from '../../context/AuthContext';

const API_BASE = process.env.REACT_APP_API_URL || '/api';

interface Paper {
  paper_id: string;
  title: string;
  filename: string;
  pages: number;
  total_chars: number;
  link_count: number;
  entity_count: number;
}

const CHISGPapers: React.FC = () => {
  const { token, logout } = useAuth();
  const navigate = useNavigate();
  const [papers, setPapers] = useState<Paper[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const load = async () => {
      try {
        const res = await fetch(`${API_BASE}/chisg/papers`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        if (res.status === 401) {
          logout();
          navigate('/chisg/login');
          return;
        }
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const data = await res.json();
        setPapers(data.papers || []);
      } catch (e: any) {
        setError(e.message || 'Failed to load papers.');
      } finally {
        setLoading(false);
      }
    };
    load();
  }, [token, logout, navigate]);

  const handleLogout = () => {
    logout();
    navigate('/chisg/login');
  };

  return (
    <Box sx={{ minHeight: '100vh', bgcolor: '#f5f7fa' }}>
      <AppBar position="static" elevation={1} sx={{ bgcolor: '#1a365d' }}>
        <Toolbar>
          <Typography variant="h6" fontWeight={700} sx={{ flexGrow: 1 }}>
            CHISG Knowledge Base
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
          Processed Papers
        </Typography>
        <Typography variant="body2" color="text.secondary" mb={3}>
          Academic papers that have been processed and stored in the knowledge graph.
        </Typography>

        {loading && (
          <Box display="flex" justifyContent="center" py={8}>
            <CircularProgress />
          </Box>
        )}

        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

        {!loading && !error && papers.length === 0 && (
          <Alert severity="info">No papers have been processed yet.</Alert>
        )}

        {!loading && papers.length > 0 && (
          <TableContainer component={Paper} elevation={2}>
            <Table>
              <TableHead>
                <TableRow sx={{ bgcolor: '#f0f4f8' }}>
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
                        <Typography variant="caption" color="text.secondary">
                          {p.filename}
                        </Typography>
                      )}
                    </TableCell>
                    <TableCell>
                      <Chip label={p.paper_id} size="small" variant="outlined" />
                    </TableCell>
                    <TableCell align="center">
                      <Chip
                        label={p.link_count}
                        size="small"
                        color="primary"
                        variant="filled"
                      />
                    </TableCell>
                    <TableCell align="center">{p.entity_count}</TableCell>
                    <TableCell align="center">{p.pages || '—'}</TableCell>
                    <TableCell align="center">
                      <Box display="flex" gap={1} justifyContent="center">
                        <Tooltip title="View in graph">
                          <Button
                            size="small"
                            variant="outlined"
                            startIcon={<AccountTreeIcon />}
                            onClick={() => navigate(`/chisg/graph?paper_id=${p.paper_id}`)}
                          >
                            Graph
                          </Button>
                        </Tooltip>
                        <Tooltip title="Review / edit links">
                          <Button
                            size="small"
                            variant="contained"
                            startIcon={<EditIcon />}
                            onClick={() =>
                              navigate(`/chisg/links?paper_id=${p.paper_id}`, {
                                state: { paper: p },
                              })
                            }
                          >
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
        )}
      </Box>
    </Box>
  );
};

export default CHISGPapers;
