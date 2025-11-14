import React, { useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Alert,
  Paper,
  Tabs,
  Tab,
  CircularProgress,
  Divider,
  List,
  ListItem,
  ListItemText
} from '@mui/material';
import { apiClient } from '../config/api';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;
  return (
    <div role="tabpanel" hidden={value !== index} {...other}>
      {value === index && <Box sx={{ p: 2 }}>{children}</Box>}
    </div>
  );
}

const DiagnosticViewer: React.FC = () => {
  const [tabValue, setTabValue] = useState(0);
  const [documentId, setDocumentId] = useState('');
  const [rawData, setRawData] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [collectionName, setCollectionName] = useState('immunology_chapters');

  const handleFetchRawData = async () => {
    if (!documentId.trim()) {
      setError('Please enter a document ID');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      // Fix: Manually construct the query string instead of using the params object
      const queryUrl = `/api/diagnostics/document?id=${encodeURIComponent(documentId)}&collection=${encodeURIComponent(collectionName)}`;
      const result: any = await apiClient.get(queryUrl);
      
      setRawData(result.data); // Access the data property of the response
    } catch (err: any) {
      console.error('Error fetching document:', err);
      setError(err.message || 'Failed to fetch document');
    } finally {
      setLoading(false);
    }
  };

  const handleCheckCollections = async () => {
    setLoading(true);
    setError(null);

    try {
      const result: any = await apiClient.get('/api/diagnostics/collections');
      setRawData(result.data); // Access the data property of the response
    } catch (err: any) {
      console.error('Error checking collections:', err);
      setError(err.message || 'Failed to check collections');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box sx={{ maxWidth: 900, mx: 'auto', p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Diagnostic Data Viewer
      </Typography>
      
      <Alert severity="info" sx={{ mb: 3 }}>
        This tool helps debug data storage issues by viewing raw MongoDB data and diagnosing pipeline errors.
      </Alert>

      <Tabs value={tabValue} onChange={(e, val) => setTabValue(val)} sx={{ mb: 2 }}>
        <Tab label="Document Lookup" />
        <Tab label="Collection Status" />
        <Tab label="Pipeline Diagnostics" />
      </Tabs>

      <TabPanel value={tabValue} index={0}>
        <Card>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Fetch Document by ID
            </Typography>

            <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
              <TextField
                label="MongoDB Document ID"
                value={documentId}
                onChange={(e) => setDocumentId(e.target.value)}
                fullWidth
                placeholder="e.g., 68cc6341895dfc76aa8a1d4b"
              />
              <TextField
                label="Collection"
                value={collectionName}
                onChange={(e) => setCollectionName(e.target.value)}
                sx={{ width: '40%' }}
              />
            </Box>

            <Button 
              variant="contained" 
              onClick={handleFetchRawData}
              disabled={loading}
            >
              Fetch Document
            </Button>

            {loading && <CircularProgress size={24} sx={{ ml: 2 }} />}
            
            {error && (
              <Alert severity="error" sx={{ mt: 2 }}>
                {error}
              </Alert>
            )}

            {rawData && (
              <Paper sx={{ mt: 3, p: 2, bgcolor: 'grey.50' }}>
                <Typography variant="subtitle2" gutterBottom>
                  Raw Document Data:
                </Typography>
                <Box sx={{ 
                  maxHeight: '400px', 
                  overflow: 'auto',
                  p: 2,
                  bgcolor: 'background.paper',
                  borderRadius: 1,
                  fontFamily: 'monospace',
                  fontSize: '0.85rem'
                }}>
                  <pre>{JSON.stringify(rawData, null, 2)}</pre>
                </Box>
              </Paper>
            )}
          </CardContent>
        </Card>
      </TabPanel>

      <TabPanel value={tabValue} index={1}>
        <Card>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Collection Status
            </Typography>
            
            <Button 
              variant="contained" 
              onClick={handleCheckCollections}
              disabled={loading}
            >
              Check Collections
            </Button>

            {loading && <CircularProgress size={24} sx={{ ml: 2 }} />}
            
            {error && (
              <Alert severity="error" sx={{ mt: 2 }}>
                {error}
              </Alert>
            )}

            {rawData && rawData.collections && (
              <Box sx={{ mt: 3 }}>
                <Typography variant="subtitle2" gutterBottom>
                  MongoDB Collections:
                </Typography>
                
                <List>
                  {Object.entries(rawData.collections).map(([collection, count]: [string, any]) => (
                    <ListItem key={collection} divider>
                      <ListItemText 
                        primary={`${collection}: ${count} documents`}
                        secondary={
                          count > 0 ? 
                          "Has data - use Document Lookup to view examples" : 
                          "Empty collection"
                        }
                      />
                    </ListItem>
                  ))}
                </List>
              </Box>
            )}
          </CardContent>
        </Card>
      </TabPanel>

      <TabPanel value={tabValue} index={2}>
        <Card>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Pipeline Status
            </Typography>
            
            <Alert severity="warning" sx={{ mb: 3 }}>
              <Typography variant="body2">
                <strong>AWS Bedrock Error:</strong> "The provided model identifier is invalid"
              </Typography>
              <Typography variant="body2" sx={{ mt: 1 }}>
                Check your AWS Bedrock configuration. The environment variable might not be updated properly.
              </Typography>
              <Button 
                variant="outlined" 
                size="small" 
                color="inherit"
                sx={{ mt: 1 }}
                onClick={async () => {
                  try {
                    setLoading(true);
                    setError(null);
                    console.log("Attempting to fetch config from: /api/diagnostics/config");
                    
                    // Use direct fetch to see raw response for debugging
                    const response = await fetch('/api/diagnostics/config');
                    console.log("Raw response status:", response.status);
                    
                    if (!response.ok) {
                      const errorText = await response.text();
                      console.error("Error response:", errorText);
                      throw new Error(`API returned ${response.status}: ${errorText}`);
                    }
                    
                    const config = await response.json();
                    console.log("Received config data:", config);
                    
                    // Force a tab switch to clearly show the results
                    setTabValue(2); // Stay on current tab
                    setRawData({
                      ...config,
                      _fetchTime: new Date().toISOString(),
                      _endpoint: '/api/diagnostics/config'
                    });
                    
                    // Use an informational message instead of error for success
                    setError("✅ Config values loaded successfully. Check results below.");
                  } catch (err: any) {
                    console.error('Error fetching config:', err);
                    // Show detailed error message
                    setError(`Failed to fetch configuration: ${err.message || 'Unknown error'}`);
                    
                    // Still provide something to display
                    setRawData({
                      error: err.message || 'Unknown error',
                      troubleshooting: {
                        message: "API endpoint may not be registered or running",
                        suggestions: [
                          "Check that the Go API server is running",
                          "Verify /api/diagnostics/config endpoint is registered",
                          "Check browser console for network errors"
                        ]
                      },
                      _fetchTime: new Date().toISOString()
                    });
                  } finally {
                    setLoading(false);
                  }
                }}
              >
                Check Actual Config Values
              </Button>

              {/* Add this button below the "Check Actual Config Values" button */}
              <Button 
                variant="outlined" 
                size="small" 
                color="primary"
                sx={{ mt: 1, ml: 1 }}
                onClick={() => {
                  // Show frontend environment variables only
                  setRawData({
                    frontend_env: {
                      NODE_ENV: process.env.NODE_ENV || 'undefined',
                      REACT_APP_API_URL: process.env.REACT_APP_API_URL || 'undefined',
                      REACT_APP_API_BASE_URL: process.env.REACT_APP_API_BASE_URL || 'undefined',
                      REACT_APP_ENVIRONMENT: process.env.REACT_APP_ENVIRONMENT || 'undefined'
                    },
                    note: "This data is from frontend environment variables only. Backend values may differ."
                  });
                  setError("✅ Frontend env values loaded (backend values not available)");
                }}
              >
                Show Frontend Env Only
              </Button>
            </Alert>
            
            {/* Always show the data section, even if empty */}
            <Paper sx={{ p: 2, mb: 3, bgcolor: '#f8f9fa' }}>
              <Typography variant="subtitle2" gutterBottom>
                Configuration Data:
              </Typography>
              
              {loading ? (
                <Box sx={{ display: 'flex', justifyContent: 'center', my: 3 }}>
                  <CircularProgress size={40} />
                </Box>
              ) : rawData ? (
                <Box sx={{ 
                  maxHeight: '400px', 
                  overflow: 'auto',
                  p: 2,
                  bgcolor: 'background.paper',
                  borderRadius: 1,
                  fontFamily: 'monospace',
                  fontSize: '0.85rem'
                }}>
                  <pre>{JSON.stringify(rawData, null, 2)}</pre>
                </Box>
              ) : (
                <Alert severity="info">
                  Click "Check Actual Config Values" or "Show Frontend Env Only" to view configuration data
                </Alert>
              )}
            </Paper>
            
            {/* Display AWS config values only if available */}
            {rawData && rawData.aws && (
              <Paper sx={{ p: 2, mb: 3, bgcolor: '#f8f9fa' }}>
                <Typography variant="subtitle2" gutterBottom>AWS Configuration Values:</Typography>
                <Box sx={{ fontFamily: 'monospace', fontSize: '0.85rem' }}>
                  <Typography variant="body2" color={rawData.aws.bedrock_embedding_model_id ? 'textPrimary' : 'error'}>
                    <strong>AWS_BEDROCK_EMBEDDING_MODEL_ID:</strong> {rawData.aws.bedrock_embedding_model_id || '(not set)'}
                  </Typography>
                  <Typography variant="body2" color={rawData.aws.bedrock_embed_model ? 'textPrimary' : 'error'}>
                    <strong>AWS_BEDROCK_EMBED_MODEL:</strong> {rawData.aws.bedrock_embed_model || '(not set)'}
                  </Typography>
                  <Typography variant="body2">
                    <strong>AWS_BEDROCK_REGION:</strong> {rawData.aws.bedrock_region || '(not set)'}
                  </Typography>
                </Box>
              </Paper>
            )}
            
            <Alert severity="info" sx={{ mb: 3 }}>
              <Typography variant="body2">
                <strong>Caching Issues:</strong> Environment variable changes may require restarting services.
              </Typography>
              <Typography variant="body2" sx={{ mt: 1 }}>
                1. Check for typos in <code>.env</code>: <code>AWS_BEDROCK_EMBED_MODEL</code> vs <code>AWS_BEDROCK_EMBEDDING_MODEL_ID</code>
              </Typography>
              <Typography variant="body2">
                2. Restart the Go API server to reload environment variables
              </Typography>
              <Typography variant="body2">
                3. Verify if <code>aamazon.titan-embed-text-v2:0</code> (typo with double 'a') appears anywhere
              </Typography>
            </Alert>
            
            <Alert severity="warning" sx={{ mb: 3 }}>
              <Typography variant="body2">
                <strong>Weaviate Error:</strong> "weaviate class 'SubjectAreaContent' does not exist"
              </Typography>
              <Typography variant="body2" sx={{ mt: 1 }}>
                <strong>Update:</strong> Architecture has changed. 'SubjectAreaContent' class has been removed. 
                Weaviate now only stores semantic links with subject as metadata, while MongoDB stores raw content.
              </Typography>
            </Alert>

            <Divider sx={{ my: 2 }} />
            
            <Typography variant="subtitle2" gutterBottom>
              Current Architecture:
            </Typography>
            
            <List>
              <ListItem>
                <ListItemText 
                  primary="MongoDB" 
                  secondary="Stores raw content (including subject/domain information)"
                />
              </ListItem>
              <ListItem>
                <ListItemText 
                  primary="Weaviate - SemanticLinks" 
                  secondary="Stores semantic connections between concepts with domain metadata"
                />
              </ListItem>
              <ListItem>
                <ListItemText 
                  primary="Weaviate - MedicalExcerpt" 
                  secondary="Specialized medical content vectors with references to MongoDB documents"
                />
              </ListItem>
            </List>
            
            <Typography variant="subtitle2" gutterBottom sx={{ mt: 2 }}>
              Common Solutions:
            </Typography>
            
            <List>
              <ListItem>
                <ListItemText 
                  primary="Clear Docker Cache" 
                  secondary="Run 'docker-compose down && docker-compose up -d' to ensure env vars are reloaded"
                />
              </ListItem>
              <ListItem>
                <ListItemText 
                  primary="Initialize Weaviate Schema" 
                  secondary="Run the schema initialization endpoint to create required classes"
                />
              </ListItem>
              <ListItem>
                <ListItemText 
                  primary="Check API Environment" 
                  secondary="Run 'docker-compose exec api env | grep AWS' to see actual environment variables"
                />
              </ListItem>
                <ListItem>
                <ListItemText
                    primary="Re-upload Content"
                    secondary="If collections are empty, re-upload chapters or case studies to populate data"
                  />
              </ListItem>
            </List>
          </CardContent>
        </Card>
      </TabPanel>
    </Box>
  );
}
export default DiagnosticViewer;

