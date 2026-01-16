import React, { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
  Alert,
  Paper,
  CircularProgress,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
  IconButton,
  Tooltip,
} from '@mui/material';
import { CheckCircle, Error as ErrorIcon, Info as InfoIcon, Refresh, Download, CheckBox, CheckBoxOutlineBlank } from '@mui/icons-material';
import { ExtractionResult } from '../../types/semanticLinks';
import { apiClient } from '../../config/api';

interface SemanticLink {
  id: string;
  source_term: string;
  target_term: string;
  forward_relation: string;
  inverse_relation: string;
  statement: string;
  source_title?: string;
  chapter?: string;
  status: string;
  is_manual: boolean;
  domain: string;
  confidence: number;
  created_at: string;
}

interface ExtractionResultsProps {
  results: ExtractionResult | null;
  onReset: () => void;
}

/**
 * ExtractionResults - Displays results from batch semantic link extraction
 * and allows browsing/managing existing links
 */
export default function ExtractionResults({ results, onReset }: ExtractionResultsProps) {
  const [links, setLinks] = useState<SemanticLink[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // Filters
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [domainFilter, setDomainFilter] = useState<string>('');

  const fetchLinks = async () => {
    setLoading(true);
    setError(null);
    try {
      const params = new URLSearchParams();
      if (statusFilter) params.append('status', statusFilter);
      if (domainFilter) params.append('domain', domainFilter);
      params.append('limit', '100');

      const response = await apiClient.get<{ links: SemanticLink[], count: number }>(
        `/api/semantic-links/export?${params.toString()}`
      );
      setLinks(response.links || []);
    } catch (err: any) {
      setError(`Failed to load links: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLinks();
  }, [statusFilter, domainFilter]);

  const handleExport = () => {
    const params = new URLSearchParams();
    if (statusFilter) params.append('status', statusFilter);
    params.append('format', 'json');
    window.open(`/esp-organizer/api/semantic-links/export?${params.toString()}`, '_blank');
  };

  // Show last extraction results if available
  if (results) {
    const { inserted, errors } = results;
    const total = inserted + errors.length;
    const successRate = total > 0 ? ((inserted / total) * 100).toFixed(1) : '0';

    return (
      <Stack spacing={2}>
        {/* Summary Card */}
        <Card sx={{ backgroundColor: '#f5f5f5' }}>
          <CardContent>
            <Stack direction="row" spacing={3} sx={{ textAlign: 'center' }}>
              <Box>
                <Typography variant="h4" sx={{ fontWeight: 'bold', color: '#4caf50' }}>
                  {inserted}
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  Links Extracted
                </Typography>
              </Box>
              <Box>
                <Typography variant="h4" sx={{ fontWeight: 'bold', color: '#1976d2' }}>
                  {successRate}%
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  Success Rate
                </Typography>
              </Box>
              <Box>
                <Typography variant="h4" sx={{ fontWeight: 'bold', color: errors.length > 0 ? '#f44336' : '#4caf50' }}>
                  {errors.length}
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  Errors
                </Typography>
              </Box>
            </Stack>
          </CardContent>
        </Card>

        {inserted > 0 && (
          <Alert severity="success" icon={<CheckCircle />}>
            Successfully extracted {inserted} semantic link{inserted !== 1 ? 's' : ''}.
          </Alert>
        )}

        {errors.length > 0 && (
          <Alert severity="error" icon={<ErrorIcon />}>
            <Typography variant="body2" sx={{ fontWeight: 'bold', mb: 1 }}>
              {errors.length} extraction{errors.length !== 1 ? 's' : ''} failed:
            </Typography>
            <Box component="ul" sx={{ pl: 2, mb: 0 }}>
              {errors.map((error, idx) => (
                <li key={idx}>
                  <Typography variant="caption">{error}</Typography>
                </li>
              ))}
            </Box>
          </Alert>
        )}

        <Button variant="contained" onClick={onReset}>
          Continue Extracting
        </Button>
      </Stack>
    );
  }

  // Default view: browse existing links
  return (
    <Stack spacing={2}>
      {/* Header with filters */}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 2 }}>
        <Typography variant="h5">📋 Extracted Links</Typography>
        <Stack direction="row" spacing={1}>
          <FormControl size="small" sx={{ minWidth: 120 }}>
            <InputLabel>Status</InputLabel>
            <Select
              value={statusFilter}
              label="Status"
              onChange={(e) => setStatusFilter(e.target.value)}
            >
              <MenuItem value="">All</MenuItem>
              <MenuItem value="draft">Draft</MenuItem>
              <MenuItem value="validated">Validated</MenuItem>
              <MenuItem value="rejected">Rejected</MenuItem>
            </Select>
          </FormControl>
          <FormControl size="small" sx={{ minWidth: 120 }}>
            <InputLabel>Domain</InputLabel>
            <Select
              value={domainFilter}
              label="Domain"
              onChange={(e) => setDomainFilter(e.target.value)}
            >
              <MenuItem value="">All</MenuItem>
              <MenuItem value="immunology">Immunology</MenuItem>
              <MenuItem value="biology">Biology</MenuItem>
              <MenuItem value="manual">Manual</MenuItem>
            </Select>
          </FormControl>
          <Tooltip title="Refresh">
            <IconButton onClick={fetchLinks} disabled={loading}>
              <Refresh />
            </IconButton>
          </Tooltip>
          <Button
            variant="outlined"
            startIcon={<Download />}
            onClick={handleExport}
          >
            Export
          </Button>
        </Stack>
      </Box>

      {error && (
        <Alert severity="error" onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {loading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
          <CircularProgress />
        </Box>
      ) : links.length === 0 ? (
        <Alert severity="info">
          No semantic links found. Start creating links in the Link Builder tab!
        </Alert>
      ) : (
        <TableContainer component={Paper}>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Entity A</TableCell>
                <TableCell>Relation</TableCell>
                <TableCell>Entity B</TableCell>
                <TableCell>Source</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Type</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {links.map((link) => (
                <TableRow key={link.id} hover>
                  <TableCell>
                    <Typography variant="body2" fontWeight="medium">
                      {link.source_term}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Chip label={link.forward_relation} size="small" variant="outlined" />
                  </TableCell>
                  <TableCell>
                    <Typography variant="body2" fontWeight="medium">
                      {link.target_term}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="caption" color="text.secondary">
                      {link.source_title || link.chapter || '-'}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Chip 
                      label={link.status || 'draft'} 
                      size="small"
                      color={link.status === 'validated' ? 'success' : link.status === 'rejected' ? 'error' : 'default'}
                    />
                  </TableCell>
                  <TableCell>
                    <Chip 
                      label={link.is_manual ? 'Manual' : 'Auto'} 
                      size="small"
                      variant="outlined"
                      color={link.is_manual ? 'primary' : 'default'}
                    />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}

      <Typography variant="caption" color="text.secondary" sx={{ textAlign: 'center' }}>
        Showing {links.length} links • Filter by status to get validated training data
      </Typography>
    </Stack>
  );
}
