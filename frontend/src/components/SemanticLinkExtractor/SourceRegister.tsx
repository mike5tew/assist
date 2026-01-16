import React, { useState, useEffect, useCallback } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  Chip,
  Alert,
  CircularProgress,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Stack,
  Tooltip,
} from '@mui/material';
import {
  Add as AddIcon,
  Delete as DeleteIcon,
  Edit as EditIcon,
  CloudUpload,
  Book,
  Description,
  Refresh,
} from '@mui/icons-material';
import { apiClient } from '../../config/api';

export interface Source {
  id?: string;
  _id?: string;
  title: string;
  author?: string;
  authors?: string[];
  publisher?: string;
  isbn?: string;
  year?: number;
  domain: string;
  subject?: string;
  description?: string;
  pdf_path?: string;
  pdf_size?: number;
  total_pages?: number;
  processing_status?: string;
  text_extracted?: boolean;
  link_count?: number;
  chapters?: Chapter[];
  created_at?: string;
  updated_at?: string;
}

export interface Chapter {
  id?: string;
  _id?: string;
  title: string;
  number?: number;
  start_page?: number;
  end_page?: number;
  text?: string;
  link_count?: number;
}

interface SourceRegisterProps {
  onSourceSelect?: (source: Source) => void;
  selectedSourceId?: string;
}

const emptySource: Partial<Source> = {
  title: '',
  author: '',
  publisher: '',
  isbn: '',
  year: new Date().getFullYear(),
  domain: 'immunology',
  subject: 'GCSE Biology',
  description: '',
};

export default function SourceRegister({ onSourceSelect, selectedSourceId }: SourceRegisterProps) {
  const [sources, setSources] = useState<Source[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // Dialog state
  const [showAddDialog, setShowAddDialog] = useState(false);
  const [showUploadDialog, setShowUploadDialog] = useState(false);
  const [editingSource, setEditingSource] = useState<Partial<Source>>(emptySource);
  const [isEditing, setIsEditing] = useState(false);
  
  // Upload state
  const [uploadFile, setUploadFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);

  const fetchSources = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.get<{ sources: Source[] }>('/api/sources?type=book');
      setSources(response.sources || []);
    } catch (err) {
      console.error('Failed to fetch sources:', err);
      setError('Failed to load sources. Please check if the API is running.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchSources();
  }, [fetchSources]);

  const handleAddSource = async () => {
    if (!editingSource.title) {
      setError('Title is required');
      return;
    }

    try {
      const payload = {
        type: 'book',
        title: editingSource.title,
        author: editingSource.author,
        authors: editingSource.author ? [editingSource.author] : [],
        publisher: editingSource.publisher,
        isbn: editingSource.isbn,
        year: editingSource.year?.toString(),
        domain: editingSource.domain || 'immunology',
        subject: editingSource.subject,
        description: editingSource.description,
      };

      if (isEditing && (editingSource.id || editingSource._id)) {
        // Update existing
        await apiClient.put(`/api/sources/${editingSource.id || editingSource._id}`, payload);
      } else {
        // Create new
        await apiClient.post('/api/sources', payload);
      }

      setShowAddDialog(false);
      setEditingSource(emptySource);
      setIsEditing(false);
      fetchSources();
    } catch (err: any) {
      setError(`Failed to save source: ${err.message}`);
    }
  };

  const handleDeleteSource = async (sourceId: string) => {
    if (!window.confirm('Are you sure you want to delete this source?')) return;
    
    try {
      await apiClient.delete(`/api/sources/${sourceId}`);
      fetchSources();
    } catch (err: any) {
      setError(`Failed to delete source: ${err.message}`);
    }
  };

  const handleEditSource = (source: Source) => {
    setEditingSource({
      ...source,
      author: source.authors?.[0] || source.author || '',
    });
    setIsEditing(true);
    setShowAddDialog(true);
  };

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file && file.type === 'application/pdf') {
      setUploadFile(file);
    } else {
      setError('Please select a PDF file');
    }
  };

  const handleUploadPDF = async (sourceId: string) => {
    if (!uploadFile) {
      setError('Please select a file first');
      return;
    }

    setUploading(true);
    try {
      const formData = new FormData();
      formData.append('pdf', uploadFile);
      formData.append('source_id', sourceId);

      await apiClient.upload('/api/sources/upload-pdf', formData);
      
      setShowUploadDialog(false);
      setUploadFile(null);
      fetchSources();
    } catch (err: any) {
      setError(`Failed to upload PDF: ${err.message}`);
    } finally {
      setUploading(false);
    }
  };

  const getStatusColor = (status?: string) => {
    switch (status) {
      case 'complete': return 'success';
      case 'processing': return 'warning';
      case 'error': return 'error';
      default: return 'default';
    }
  };

  const getSourceId = (source: Source) => source.id || source._id || '';

  return (
    <Box>
      {/* Header */}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h5" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Book color="primary" />
          Source Register
        </Typography>
        <Stack direction="row" spacing={1}>
          <Button
            variant="outlined"
            startIcon={<Refresh />}
            onClick={fetchSources}
            disabled={loading}
          >
            Refresh
          </Button>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => {
              setEditingSource(emptySource);
              setIsEditing(false);
              setShowAddDialog(true);
            }}
          >
            Add Source
          </Button>
        </Stack>
      </Box>

      {/* Error display */}
      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {/* Loading state */}
      {loading && (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
          <CircularProgress />
        </Box>
      )}

      {/* Sources table */}
      {!loading && (
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Title</TableCell>
                <TableCell>Author</TableCell>
                <TableCell>Domain</TableCell>
                <TableCell>Status</TableCell>
                <TableCell align="center">Links</TableCell>
                <TableCell align="right">Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {sources.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} align="center">
                    <Typography color="text.secondary" sx={{ py: 4 }}>
                      No sources found. Add your first source to get started.
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                sources.map((source) => {
                  const sourceId = getSourceId(source);
                  const isSelected = selectedSourceId === sourceId;
                  
                  return (
                    <TableRow 
                      key={sourceId}
                      sx={{ 
                        cursor: onSourceSelect ? 'pointer' : 'default',
                        backgroundColor: isSelected ? 'action.selected' : undefined,
                        '&:hover': { backgroundColor: 'action.hover' },
                      }}
                      onClick={() => onSourceSelect?.(source)}
                    >
                      <TableCell>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          {source.pdf_path ? (
                            <Description color="success" fontSize="small" />
                          ) : (
                            <Description color="disabled" fontSize="small" />
                          )}
                          <Box>
                            <Typography variant="body2" fontWeight="medium">
                              {source.title}
                            </Typography>
                            {source.isbn && (
                              <Typography variant="caption" color="text.secondary">
                                ISBN: {source.isbn}
                              </Typography>
                            )}
                          </Box>
                        </Box>
                      </TableCell>
                      <TableCell>
                        {source.authors?.join(', ') || source.author || '-'}
                      </TableCell>
                      <TableCell>
                        <Chip 
                          label={source.domain || 'general'} 
                          size="small" 
                          variant="outlined"
                        />
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={source.processing_status || 'pending'}
                          size="small"
                          color={getStatusColor(source.processing_status)}
                        />
                      </TableCell>
                      <TableCell align="center">
                        <Chip
                          label={source.link_count || 0}
                          size="small"
                          color={source.link_count ? 'primary' : 'default'}
                        />
                      </TableCell>
                      <TableCell align="right">
                        <Stack direction="row" spacing={0} justifyContent="flex-end">
                          <Tooltip title="Upload PDF">
                            <IconButton
                              size="small"
                              onClick={(e) => {
                                e.stopPropagation();
                                setEditingSource(source);
                                setShowUploadDialog(true);
                              }}
                            >
                              <CloudUpload fontSize="small" />
                            </IconButton>
                          </Tooltip>
                          <Tooltip title="Edit">
                            <IconButton
                              size="small"
                              onClick={(e) => {
                                e.stopPropagation();
                                handleEditSource(source);
                              }}
                            >
                              <EditIcon fontSize="small" />
                            </IconButton>
                          </Tooltip>
                          <Tooltip title="Delete">
                            <IconButton
                              size="small"
                              color="error"
                              onClick={(e) => {
                                e.stopPropagation();
                                handleDeleteSource(sourceId);
                              }}
                            >
                              <DeleteIcon fontSize="small" />
                            </IconButton>
                          </Tooltip>
                        </Stack>
                      </TableCell>
                    </TableRow>
                  );
                })
              )}
            </TableBody>
          </Table>
        </TableContainer>
      )}

      {/* Add/Edit Source Dialog */}
      <Dialog open={showAddDialog} onClose={() => setShowAddDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>{isEditing ? 'Edit Source' : 'Add New Source'}</DialogTitle>
        <DialogContent>
          <Stack spacing={2} sx={{ mt: 1 }}>
            <TextField
              label="Title"
              value={editingSource.title || ''}
              onChange={(e) => setEditingSource({ ...editingSource, title: e.target.value })}
              fullWidth
              required
            />
            <TextField
              label="Author(s)"
              value={editingSource.author || ''}
              onChange={(e) => setEditingSource({ ...editingSource, author: e.target.value })}
              fullWidth
              helperText="For multiple authors, separate with commas"
            />
            <Stack direction="row" spacing={2}>
              <TextField
                label="ISBN"
                value={editingSource.isbn || ''}
                onChange={(e) => setEditingSource({ ...editingSource, isbn: e.target.value })}
                fullWidth
              />
              <TextField
                label="Year"
                type="number"
                value={editingSource.year || ''}
                onChange={(e) => setEditingSource({ ...editingSource, year: parseInt(e.target.value) })}
                sx={{ width: 120 }}
              />
            </Stack>
            <TextField
              label="Publisher"
              value={editingSource.publisher || ''}
              onChange={(e) => setEditingSource({ ...editingSource, publisher: e.target.value })}
              fullWidth
            />
            <Stack direction="row" spacing={2}>
              <FormControl fullWidth>
                <InputLabel>Domain</InputLabel>
                <Select
                  value={editingSource.domain || 'immunology'}
                  label="Domain"
                  onChange={(e) => setEditingSource({ ...editingSource, domain: e.target.value })}
                >
                  <MenuItem value="immunology">Immunology</MenuItem>
                  <MenuItem value="biology">Biology</MenuItem>
                  <MenuItem value="medicine">Medicine</MenuItem>
                  <MenuItem value="education">Education</MenuItem>
                  <MenuItem value="other">Other</MenuItem>
                </Select>
              </FormControl>
              <TextField
                label="Subject"
                value={editingSource.subject || ''}
                onChange={(e) => setEditingSource({ ...editingSource, subject: e.target.value })}
                fullWidth
                placeholder="e.g., GCSE Biology"
              />
            </Stack>
            <TextField
              label="Description"
              value={editingSource.description || ''}
              onChange={(e) => setEditingSource({ ...editingSource, description: e.target.value })}
              fullWidth
              multiline
              rows={2}
            />
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowAddDialog(false)}>Cancel</Button>
          <Button variant="contained" onClick={handleAddSource}>
            {isEditing ? 'Save Changes' : 'Add Source'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Upload PDF Dialog */}
      <Dialog open={showUploadDialog} onClose={() => setShowUploadDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Upload PDF for: {editingSource.title}</DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 2 }}>
            <Button
              variant="outlined"
              component="label"
              startIcon={<CloudUpload />}
              fullWidth
            >
              Select PDF File
              <input
                type="file"
                hidden
                accept=".pdf"
                onChange={handleFileSelect}
              />
            </Button>
            {uploadFile && (
              <Alert severity="info" sx={{ mt: 2 }}>
                Selected: {uploadFile.name} ({(uploadFile.size / (1024 * 1024)).toFixed(2)} MB)
              </Alert>
            )}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowUploadDialog(false)}>Cancel</Button>
          <Button 
            variant="contained" 
            onClick={() => handleUploadPDF(getSourceId(editingSource as Source))}
            disabled={!uploadFile || uploading}
          >
            {uploading ? <CircularProgress size={20} /> : 'Upload'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
