import React, { useEffect, useState, useCallback } from 'react';
import {
  Box, Typography, Paper, Grid, CircularProgress, ToggleButtonGroup, ToggleButton,
  Table, TableBody, TableCell, TableHead, TableRow, TextField, Alert,
} from '@mui/material';

interface PageStat   { path: string; views: number }
interface RefStat    { referrer: string; count: number }
interface EventStat  { name: string; count: number }
interface DayStat    { date: string; views: number }

interface Stats {
  period: string;
  total_views: number;
  unique_visitors: number;
  top_pages: PageStat[];
  top_referrers: RefStat[];
  top_events: EventStat[];
  views_by_day: DayStat[];
}

const AnalyticsDashboard: React.FC = () => {
  const [period, setPeriod] = useState('7d');
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [token, setToken] = useState(() => localStorage.getItem('analytics_token') || '');

  const fetchStats = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const headers: Record<string, string> = { 'Content-Type': 'application/json' };
      if (token) headers['Authorization'] = `Bearer ${token}`;

      const resp = await fetch(`/api/analytics/stats?period=${period}`, { headers });
      if (resp.status === 401) {
        setError('Unauthorized — enter your analytics token below.');
        setLoading(false);
        return;
      }
      if (!resp.ok) throw new Error(`HTTP ${resp.status}`);
      const data = await resp.json();
      setStats(data);
    } catch (e: any) {
      setError(e.message || 'Failed to load stats');
    } finally {
      setLoading(false);
    }
  }, [period, token]);

  useEffect(() => { fetchStats(); }, [fetchStats]);

  const saveToken = (t: string) => {
    setToken(t);
    localStorage.setItem('analytics_token', t);
  };

  // Simple bar chart using CSS
  const maxDay = stats?.views_by_day?.reduce((m, d) => Math.max(m, d.views), 0) || 1;

  return (
    <Box sx={{ p: 3, maxWidth: 1000, mx: 'auto' }}>
      <Typography variant="h4" gutterBottom>Visitor Analytics</Typography>

      <Box sx={{ display: 'flex', gap: 2, mb: 3, alignItems: 'center', flexWrap: 'wrap' }}>
        <ToggleButtonGroup
          value={period}
          exclusive
          onChange={(_, v) => v && setPeriod(v)}
          size="small"
        >
          <ToggleButton value="24h">24 h</ToggleButton>
          <ToggleButton value="7d">7 days</ToggleButton>
          <ToggleButton value="30d">30 days</ToggleButton>
          <ToggleButton value="90d">90 days</ToggleButton>
        </ToggleButtonGroup>

        <TextField
          size="small"
          label="Token"
          type="password"
          value={token}
          onChange={(e) => saveToken(e.target.value)}
          sx={{ width: 220 }}
        />
      </Box>

      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
      {loading && <CircularProgress />}

      {stats && !loading && (
        <>
          {/* Summary cards */}
          <Grid container spacing={2} sx={{ mb: 3 }}>
            <Grid item xs={6} md={3}>
              <Paper sx={{ p: 2, textAlign: 'center' }}>
                <Typography variant="h3">{stats.total_views.toLocaleString()}</Typography>
                <Typography color="text.secondary">Page Views</Typography>
              </Paper>
            </Grid>
            <Grid item xs={6} md={3}>
              <Paper sx={{ p: 2, textAlign: 'center' }}>
                <Typography variant="h3">{stats.unique_visitors.toLocaleString()}</Typography>
                <Typography color="text.secondary">Unique Visitors</Typography>
              </Paper>
            </Grid>
            <Grid item xs={6} md={3}>
              <Paper sx={{ p: 2, textAlign: 'center' }}>
                <Typography variant="h3">{stats.top_pages?.length || 0}</Typography>
                <Typography color="text.secondary">Pages Visited</Typography>
              </Paper>
            </Grid>
            <Grid item xs={6} md={3}>
              <Paper sx={{ p: 2, textAlign: 'center' }}>
                <Typography variant="h3">{stats.top_events?.length || 0}</Typography>
                <Typography color="text.secondary">Event Types</Typography>
              </Paper>
            </Grid>
          </Grid>

          {/* Views by day chart */}
          {stats.views_by_day && stats.views_by_day.length > 0 && (
            <Paper sx={{ p: 2, mb: 3 }}>
              <Typography variant="h6" gutterBottom>Views by Day</Typography>
              <Box sx={{ display: 'flex', alignItems: 'flex-end', gap: '2px', height: 120 }}>
                {stats.views_by_day.map((d) => (
                  <Box
                    key={d.date}
                    title={`${d.date}: ${d.views}`}
                    sx={{
                      flex: 1,
                      minWidth: 6,
                      height: `${(d.views / maxDay) * 100}%`,
                      minHeight: 4,
                      bgcolor: 'primary.main',
                      borderRadius: '2px 2px 0 0',
                      '&:hover': { bgcolor: 'primary.light' },
                    }}
                  />
                ))}
              </Box>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 0.5 }}>
                <Typography variant="caption">{stats.views_by_day[0]?.date}</Typography>
                <Typography variant="caption">{stats.views_by_day[stats.views_by_day.length - 1]?.date}</Typography>
              </Box>
            </Paper>
          )}

          {/* Tables */}
          <Grid container spacing={2}>
            <Grid item xs={12} md={4}>
              <Paper sx={{ p: 2 }}>
                <Typography variant="h6" gutterBottom>Top Pages</Typography>
                <Table size="small">
                  <TableHead>
                    <TableRow>
                      <TableCell>Page</TableCell>
                      <TableCell align="right">Views</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {(stats.top_pages || []).map((p) => (
                      <TableRow key={p.path}>
                        <TableCell sx={{ maxWidth: 180, overflow: 'hidden', textOverflow: 'ellipsis' }}>{p.path}</TableCell>
                        <TableCell align="right">{p.views}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </Paper>
            </Grid>

            <Grid item xs={12} md={4}>
              <Paper sx={{ p: 2 }}>
                <Typography variant="h6" gutterBottom>Referrers</Typography>
                <Table size="small">
                  <TableHead>
                    <TableRow>
                      <TableCell>Source</TableCell>
                      <TableCell align="right">Hits</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {(stats.top_referrers || []).map((r) => (
                      <TableRow key={r.referrer}>
                        <TableCell sx={{ maxWidth: 180, overflow: 'hidden', textOverflow: 'ellipsis' }}>{r.referrer}</TableCell>
                        <TableCell align="right">{r.count}</TableCell>
                      </TableRow>
                    ))}
                    {(!stats.top_referrers || stats.top_referrers.length === 0) && (
                      <TableRow><TableCell colSpan={2}>No referrer data yet</TableCell></TableRow>
                    )}
                  </TableBody>
                </Table>
              </Paper>
            </Grid>

            <Grid item xs={12} md={4}>
              <Paper sx={{ p: 2 }}>
                <Typography variant="h6" gutterBottom>Events (CTAs)</Typography>
                <Table size="small">
                  <TableHead>
                    <TableRow>
                      <TableCell>Event</TableCell>
                      <TableCell align="right">Count</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {(stats.top_events || []).map((e) => (
                      <TableRow key={e.name}>
                        <TableCell>{e.name}</TableCell>
                        <TableCell align="right">{e.count}</TableCell>
                      </TableRow>
                    ))}
                    {(!stats.top_events || stats.top_events.length === 0) && (
                      <TableRow><TableCell colSpan={2}>No events yet</TableCell></TableRow>
                    )}
                  </TableBody>
                </Table>
              </Paper>
            </Grid>
          </Grid>
        </>
      )}
    </Box>
  );
};

export default AnalyticsDashboard;
