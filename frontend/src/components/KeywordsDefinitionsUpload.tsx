import React, { useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Alert,
  CircularProgress,
  Grid,
  Paper,
  Divider,
  Chip,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow
} from '@mui/material';
import { 
  CloudUpload, 
  ImportContacts, 
  Info,
  CheckCircle
} from '@mui/icons-material';
import Papa from 'papaparse';
import { apiClient } from '../config/api';

interface KeywordDefinition {
  keyword: string;
  definition: string;
  subject?: string;
  category?: string;
}

interface UploadResponse {
  success: boolean;
  message: string;
  count: number;
  failed?: {
    count: number;
    entries: KeywordDefinition[];
  };
  uploaded?: KeywordDefinition[];
}

const KeywordsDefinitionsUpload: React.FC<{ subject: string }> = ({ subject }) => {
  const [csvContent, setCsvContent] = useState<string>('');
  const [parsedData, setParsedData] = useState<KeywordDefinition[]>([]);
  const [uploading, setUploading] = useState(false);
  const [uploadResult, setUploadResult] = useState<UploadResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [previewMode, setPreviewMode] = useState(false);

  const handleCsvParse = () => {
    if (!csvContent.trim()) {
      setError('Please paste CSV content');
      return;
    }

    try {
      // Parse CSV content
      const result = Papa.parse<any>(csvContent, {
        header: true,
        skipEmptyLines: true
      });

      if (result.errors.length > 0) {
        setError(`CSV parsing error: ${result.errors[0].message}`);
        return;
      }

      // Map to our format and validate
      const keywords: KeywordDefinition[] = result.data.map((row: any) => ({
        keyword: row.keyword || row.Keyword || row.term || row.Term || '',
        definition: row.definition || row.Definition || row.description || row.Description || '',
        category: row.category || row.Category || '',
        subject: subject
      }));

      // Basic validation
      const invalidEntries = keywords.filter(k => !k.keyword || !k.definition);
      if (invalidEntries.length > 0) {
        setError(`Found ${invalidEntries.length} entries missing keyword or definition`);
      } else {
        setError(null);
      }

      setParsedData(keywords);
      setPreviewMode(true);
    } catch (err) {
      setError('Failed to parse CSV. Please check the format.');
    }
  };

  const handleUpload = async () => {
    if (parsedData.length === 0) {
      setError('No valid data to upload');
      return;
    }

    setUploading(true);
    setError(null);

    try {
      // Use apiClient.post which now directly returns the data, not a response object
      const response = await apiClient.post<UploadResponse>(`/api/${subject.toLowerCase()}/keywords-upload`, {
        keywords: parsedData,
        subject: subject
      });

      setUploadResult(response); // No need for .data anymore
    } catch (err) {
      console.error('Error uploading keywords:', err);
      setError('Failed to upload keywords and definitions');
    } finally {
      setUploading(false);
    }
  };

  return (
    <Box sx={{ maxWidth: 800, mx: 'auto', p: 3 }}>
      <Typography variant="h4" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <ImportContacts color="primary" />
        {subject} Keywords & Definitions Upload
      </Typography>

      {/* CSV Input Section */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Paste CSV Content
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            Format should have columns for keyword/term and definition/description.
          </Typography>

          <TextField
            fullWidth
            multiline
            rows={10}
            label="CSV Content"
            value={csvContent}
            onChange={(e) => setCsvContent(e.target.value)}
            placeholder="keyword,definition&#10;antibody,A protective protein produced by the immune system"
            sx={{ mb: 2 }}
          />

          <Button
            variant="contained"
            onClick={handleCsvParse}
            startIcon={<CloudUpload />}
            sx={{ mr: 2 }}
          >
            Parse CSV
          </Button>

          <Button
            variant="outlined"
            onClick={() => {
              setCsvContent('');
              setParsedData([]);
              setPreviewMode(false);
              setUploadResult(null);
              setError(null);
            }}
          >
            Clear
          </Button>
        </CardContent>
      </Card>

      {/* Preview Section */}
      {previewMode && parsedData.length > 0 && (
        <Card sx={{ mb: 3 }}>
          <CardContent>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6">
                Preview ({parsedData.length} entries)
              </Typography>
              <Button
                variant="contained"
                color="primary"
                onClick={handleUpload}
                disabled={uploading}
                startIcon={uploading ? <CircularProgress size={20} /> : <CloudUpload />}
              >
                {uploading ? 'Uploading...' : 'Upload Keywords'}
              </Button>
            </Box>

            <TableContainer component={Paper} sx={{ maxHeight: 400 }}>
              <Table stickyHeader size="small">
                <TableHead>
                  <TableRow>
                    <TableCell>#</TableCell>
                    <TableCell>Keyword</TableCell>
                    <TableCell>Definition</TableCell>
                    <TableCell>Category</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {parsedData.slice(0, 50).map((row, index) => (
                    <TableRow key={index}>
                      <TableCell>{index + 1}</TableCell>
                      <TableCell>{row.keyword}</TableCell>
                      <TableCell>{row.definition.substring(0, 100)}
                        {row.definition.length > 100 ? '...' : ''}
                      </TableCell>
                      <TableCell>{row.category || 'General'}</TableCell>
                    </TableRow>
                  ))}
                  {parsedData.length > 50 && (
                    <TableRow>
                      <TableCell colSpan={4} align="center">
                        <Typography variant="body2" color="text.secondary">
                          {parsedData.length - 50} more entries not shown in preview
                        </Typography>
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </TableContainer>
          </CardContent>
        </Card>
      )}

      {/* Error Display */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Upload Results */}
      {uploadResult && (
        <Card sx={{ mb: 3 }}>
          <CardContent>
            <Alert 
              severity={uploadResult.success ? "success" : "warning"}
              icon={uploadResult.success ? <CheckCircle /> : undefined}
              sx={{ mb: 2 }}
            >
              <Typography variant="subtitle1">{uploadResult.message}</Typography>
            </Alert>

            <Grid container spacing={2}>
              <Grid item xs={12} md={6}>
                <Paper sx={{ p: 2, bgcolor: 'background.default' }}>
                  <Typography variant="subtitle2" gutterBottom>Upload Summary</Typography>
                  <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                    <Chip 
                      label={`${uploadResult.count || 0} Uploaded`} 
                      color="success" 
                      variant="outlined" 
                    />
                    {uploadResult.failed && (
                      <Chip 
                        label={`${uploadResult.failed.count || 0} Failed`} 
                        color="error" 
                        variant="outlined" 
                      />
                    )}
                  </Box>
                </Paper>
              </Grid>
              
              <Grid item xs={12}>
                <Divider sx={{ my: 2 }} />
                
                {uploadResult.uploaded && uploadResult.uploaded.length > 0 && (
                  <Box sx={{ mt: 2 }}>
                    <Typography variant="subtitle2" gutterBottom>
                      Sample of Uploaded Keywords:
                    </Typography>
                    <TableContainer component={Paper} sx={{ maxHeight: 200 }}>
                      <Table size="small">
                        <TableHead>
                          <TableRow>
                            <TableCell>Keyword</TableCell>
                            <TableCell>Definition</TableCell>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          {uploadResult.uploaded.slice(0, 5).map((item, idx) => (
                            <TableRow key={idx}>
                              <TableCell>{item.keyword}</TableCell>
                              <TableCell>{item.definition}</TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </TableContainer>
                  </Box>
                )}
              </Grid>
            </Grid>
          </CardContent>
        </Card>
      )}

      {/* Help Section */}
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Info />
            CSV Format Information
          </Typography>
          
          <Typography variant="body2" paragraph>
            The CSV should contain at least these columns:
          </Typography>
          
          <Box component="ul" sx={{ pl: 2 }}>
            <li>
              <Typography variant="body2">
                <strong>keyword</strong> or <strong>term</strong> (required) - The keyword or term
              </Typography>
            </li>
            <li>
              <Typography variant="body2">
                <strong>definition</strong> or <strong>description</strong> (required) - The definition
              </Typography>
            </li>
            <li>
              <Typography variant="body2">
                <strong>category</strong> (optional) - Category or classification
              </Typography>
            </li>
          </Box>
          
          <Box sx={{ mt: 2 }}>
            <Typography variant="subtitle2" gutterBottom>Example CSV:</Typography>
            <Paper sx={{ p: 2, bgcolor: 'grey.100' }}>
              <Typography variant="body2" component="pre" sx={{ fontFamily: 'monospace' }}>
                keyword,definition,category
                {'\n'}antibody,A protective protein produced by the immune system,immunology
                {'\n'}antigen,A substance that triggers an immune response,immunology
                {'\n'}cytokine,Cell signaling molecules that aid cell to cell communication,cell biology
              </Typography>
            </Paper>
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
};

export default KeywordsDefinitionsUpload;
