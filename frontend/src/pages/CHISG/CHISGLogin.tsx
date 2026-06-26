import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  Box,
  Card,
  CardContent,
  TextField,
  Button,
  Typography,
  Alert,
  CircularProgress,
} from '@mui/material';
import { useAuth } from '../../context/AuthContext';

const API_BASE = process.env.REACT_APP_API_URL || '/api';

const CHISGLogin: React.FC = () => {
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const from = (location.state as any)?.from?.pathname || '/ntm';

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [demoLoading, setDemoLoading] = useState(false);

  const handleDemoLogin = async () => {
    setError('');
    setDemoLoading(true);
    try {
      const res = await fetch(`${API_BASE}/auth/demo-login`);
      if (!res.ok) {
        setError('Demo login unavailable — please try again.');
        return;
      }
      const data = await res.json();
      login(data.token);
      navigate(from, { replace: true });
    } catch {
      setError('Network error — please try again.');
    } finally {
      setDemoLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const res = await fetch(`${API_BASE}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      });
      if (!res.ok) {
        setError('Invalid username or password.');
        return;
      }
      const data = await res.json();
      login(data.token);
      navigate(from, { replace: true });
    } catch {
      setError('Network error — please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        bgcolor: '#f5f7fa',
      }}
    >
      <Card sx={{ width: 380, boxShadow: 3 }}>
        <CardContent sx={{ p: 4 }}>
          <Typography variant="h5" fontWeight={700} mb={0.5}>
            CHISG Admin
          </Typography>
          <Typography variant="body2" color="text.secondary" mb={3}>
            Sign in to access the knowledge graph tools.
          </Typography>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          <Box component="form" onSubmit={handleSubmit} noValidate>
            <TextField
              label="Username"
              fullWidth
              required
              autoFocus
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              sx={{ mb: 2 }}
            />
            <TextField
              label="Password"
              type="password"
              fullWidth
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              sx={{ mb: 3 }}
            />
            <Button
              type="submit"
              variant="contained"
              fullWidth
              disabled={loading || !username || !password}
              sx={{ py: 1.5 }}
            >
              {loading ? <CircularProgress size={22} color="inherit" /> : 'Sign In'}
            </Button>
          </Box>

          <Box sx={{ mt: 2, textAlign: 'center' }}>
            <Typography variant="caption" color="text.secondary" display="block" sx={{ mb: 1 }}>
              — or —
            </Typography>
            <Button
              variant="outlined"
              fullWidth
              disabled={demoLoading}
              onClick={handleDemoLogin}
              sx={{ py: 1 }}
            >
              {demoLoading ? <CircularProgress size={20} color="inherit" /> : 'Try Demo Access'}
            </Button>
            <Typography variant="caption" color="text.secondary" sx={{ mt: 0.5, display: 'block' }}>
              Read-only access to the CHISG knowledge base
            </Typography>
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
};

export default CHISGLogin;
