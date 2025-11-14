import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Card, 
  CardContent, 
  Grid,
  Alert,
  CircularProgress,
  Tabs,
  Tab,
  TextField,
  Button,
  AppBar,
  Toolbar,
  IconButton,
  Chip,
  Paper,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Divider,
  List,
  ListItem,
  ListItemText
} from '@mui/material';
import { Chapter, CaseStudy, MedicalTerm } from '../../types';
import { useNavigate } from 'react-router-dom';
import { 
  ArrowBack, 
  ExpandMore, 
  Science, 
  Search,
  Assignment,
  Book,
  Person
} from '@mui/icons-material';
import { apiClient } from '../../config/api';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;
  return (
    <div role="tabpanel" hidden={value !== index} {...other}>
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

// Update the StudyAreaConfig interface to include subject-specific data types
export interface StudyAreaConfig {
  title: string;
  domain: string;
  icon: React.ReactNode;
  color: string;
  searchPlaceholder: string;
  helpText?: string; // Optional help text for users
  apiEndpoints: {
    search: string;
    contentLibrary?: string;
    caseStudies?: string;
    terms?: string;
    chapters?: string;
    upload?: string;
    keywordsUpload?: string; // New endpoint for keywords upload
    extraction?: {
      status: (jobId: string) => string;
    };
  };
  uploadComponent?: React.ReactNode;
  dataTypes?: {
    supportsKeywords?: boolean;
  };
  additionalTabs?: string[];
}

// Define interfaces for our data types
interface SearchResult {
  ai_response: {
    type: 'success' | 'error' | 'not_found';
    content?: string;
  };
}

const GenericStudyArea: React.FC<{ config: StudyAreaConfig }> = ({ config }) => {
  const navigate = useNavigate();
  const [tabValue, setTabValue] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [query, setQuery] = useState('');
  const [searchResults, setSearchResults] = useState<any>(null);

  // Library content states
  const [chapters, setChapters] = useState<Chapter[]>([]);
  const [caseStudies, setCaseStudies] = useState<CaseStudy[]>([]);
  const [medicalTerms, setMedicalTerms] = useState<MedicalTerm[]>([]);
  const [contentLoading, setContentLoading] = useState(false);
  const [loadStatus, setLoadStatus] = useState<{
    chapters: boolean;
    caseStudies: boolean;
    terms: boolean;
  }>({
    chapters: false,
    caseStudies: false,
    terms: false
  });

  // Add state for keywords and definitions CSV
  const [csvContent, setCsvContent] = useState<string>('');
  const [csvUploadResult, setCsvUploadResult] = useState<any>(null);
  const [csvUploading, setCsvUploading] = useState(false);
  
  // Fetch content on component mount
  useEffect(() => {
    fetchLibraryContent();
  }, []);

  // Handle search
 const handleSearch = async () => {
  if (!query.trim()) {
    setError('Please enter a search term');
    return;
  }
  console.log("Searching with query:", query);
  setLoading(true);
  setError(null);

  const endpoint = config.apiEndpoints.search;
  if (!endpoint) {
    setError('Search endpoint is not configured.');
    setLoading(false);
    return;
  }

  // -----------------------------------------------------------------
  // 🎯 CORRECTION START: Constructing the Filter Object for Subject
  // -----------------------------------------------------------------
  
  // 1. Define the Operand for the subject restriction
  const subjectOperand = {
    // Path: ["subject"] is what the Go backend looks for
    Path: ["subject"], 
    Operator: "Equal", 
    // The subject value is the domain of the current study area (e.g., "immunology")
    ValueString: config.domain, 
    // Other value fields (ValueNumber, ValueBoolean) are omitted/default
    // NestedFilter is null/omitted
  };

  // 2. Define the main Filter object
  const subjectFilter = {
    // We only have one condition (the subject), so "And" is appropriate
    Operator: "And", 
    Operands: [subjectOperand],
    // NestedFilter is null/omitted
  };
  
    try {
      //This is the structure expected from the api
    //   var requestBody struct {
		// 	Query   string         `json:"query"`   // Changed to lowercase to match frontend
		// 	Message string         `json:"message"` // Changed to lowercase
		// 	Filters *models.Filter `json:"filters"`
		// }
//     type Filter struct {
// 	Operator     string    `json:"operator,omitempty"`     // "And" or "Or"
// 	Operands     []Operand `json:"operands,omitempty"`     // List of operands in this filter
// 	NestedFilter *Filter   `json:"nestedFilter,omitempty"` // Support for nested filters
// }
      //
      const result = await apiClient.post(endpoint, {
        query: query,
        message: "",
        filters: subjectFilter
      });
      setSearchResults(result);
    } catch (error) {
      console.error(`Error searching ${config.domain} content:`, error);
      setError(`Error performing ${config.domain} search. Make sure the API server is running.`);
    } finally {
      setLoading(false);
    }
  };

  // Helper function to fetch a specific content type
  const fetchContent = async (
    endpoint: string | undefined, 
    setData: (data: any[]) => void, 
    dataKey: string, 
    contentName: string,
    setLoadStatusKey: 'chapters' | 'caseStudies' | 'terms'
  ) => {
    if (!endpoint) {
      console.warn(`Endpoint for ${contentName} is not defined in config. Skipping fetch.`);
      setLoadStatus(prev => ({ ...prev, [setLoadStatusKey]: true }));
      return;
    }
    try {
      console.log(`Fetching ${config.domain} ${contentName} from ${endpoint}`);
      const response: any = await apiClient.get(endpoint);
      console.log(`${contentName} response:`, response);
      
      // More robust data extraction handling null values and providing better debug info
      let data: any[] = [];
      
      // Fix: Don't try to access the response object directly with dataKey
      // Instead, always go through response.data
      const responseData = response.data;
      
      if (responseData && typeof responseData === 'object') {
        if (dataKey in responseData) {
          // Data is directly in response.data[dataKey]
          const extractedData = responseData[dataKey];
          data = Array.isArray(extractedData) ? extractedData : (extractedData ? [extractedData] : []);
        } else if (responseData.data && typeof responseData.data === 'object' && dataKey in responseData.data) {
          // Data is nested in response.data.data[dataKey]
          const extractedData = responseData.data[dataKey];
          data = Array.isArray(extractedData) ? extractedData : (extractedData ? [extractedData] : []);
        }
      }
      
      // Check if data is still empty and log more details about the response structure
      if (data.length === 0) {
        console.warn(`No ${contentName} found in response. Response data:`, responseData);
      }
      
      setData(data);
    } catch (err) {
      console.error(`Error fetching ${contentName}:`, err);
      setData([]); // Set to empty array on error
    } finally {
      setLoadStatus(prev => ({ ...prev, [setLoadStatusKey]: true }));
    }
  };

  // Fetch library content
  const fetchLibraryContent = async () => {
    setContentLoading(true);
    setError(null);
    
    setLoadStatus({ chapters: false, caseStudies: false, terms: false });

    const fetchPromises = [
      fetchContent(config.apiEndpoints.chapters, setChapters, 'chapters', 'chapters', 'chapters'),
      fetchContent(config.apiEndpoints.caseStudies, setCaseStudies, 'case_studies', 'case studies', 'caseStudies'),
      fetchContent(config.apiEndpoints.terms, setMedicalTerms, 'medical_terms', 'terms', 'terms')
    ];

    try {
      await Promise.all(fetchPromises);
    } catch (err) {
      console.error(`Error fetching ${config.domain} content:`, err);
      setError(`Failed to load ${config.domain} content. Please check the API connection.`);
    } finally {
      setContentLoading(false);
    }
  };

  // Add function to handle CSV upload
  const handleCsvUpload = async () => {
    if (!csvContent.trim()) {
      setError('Please enter CSV content');
      return;
    }

    setCsvUploading(true);
    setError(null);
    
    try {
      // Use the keywords upload endpoint
      const endpoint = config.apiEndpoints.keywordsUpload || '/api/keywords/upload';
      const result = await apiClient.post(endpoint, {
        csvContent,
        domain: config.domain
      });
      setCsvUploadResult(result);
    } catch (error) {
      console.error(`Error uploading ${config.domain} keywords:`, error);
      setError(`Error uploading keywords and definitions. Please check the CSV format and try again.`);
    } finally {
      setCsvUploading(false);
    }
  };

  // Helper function to render content loading status
  const renderLoadingStatus = () => {
    return (
      <Box sx={{ mb: 3, p: 2, bgcolor: 'grey.100', borderRadius: 1 }}>
        <Typography variant="subtitle2" gutterBottom>Loading Status:</Typography>
        <List dense>
          {config.apiEndpoints.chapters && (
            <ListItem>
              <ListItemText 
                primary="Chapters" 
                secondary={
                  loadStatus.chapters ? 
                  `${chapters.length} loaded from ${config.apiEndpoints.chapters}` : 
                  `Loading from ${config.apiEndpoints.chapters}...`
                }
              />
              {loadStatus.chapters && (
                <Chip 
                  label={chapters.length > 0 ? "✓" : "No data"} 
                  color={chapters.length > 0 ? "success" : "warning"} 
                  size="small" 
                />
              )}
            </ListItem>
          )}
          
          {config.apiEndpoints.caseStudies && (
            <ListItem>
              <ListItemText 
                primary="Case Studies" 
                secondary={
                  loadStatus.caseStudies ? 
                  `${caseStudies.length} loaded from ${config.apiEndpoints.caseStudies}` : 
                  `Loading from ${config.apiEndpoints.caseStudies}...`
                }
              />
              {loadStatus.caseStudies && (
                <Chip 
                  label={caseStudies.length > 0 ? "✓" : "No data"} 
                  color={caseStudies.length > 0 ? "success" : "warning"} 
                  size="small" 
                />
              )}
            </ListItem>
          )}
          
          {config.apiEndpoints.terms && (
            <ListItem>
              <ListItemText 
                primary="Terms" 
                secondary={
                  loadStatus.terms ? 
                  `${medicalTerms.length} loaded from ${config.apiEndpoints.terms}` : 
                  `Loading from ${config.apiEndpoints.terms}...`
                }
              />
              {loadStatus.terms && (
                <Chip 
                  label={medicalTerms.length > 0 ? "✓" : "No data"} 
                  color={medicalTerms.length > 0 ? "success" : "warning"} 
                  size="small" 
                />
              )}
            </ListItem>
          )}
        </List>
      </Box>
    );
  };

  return (
    <Box sx={{ flexGrow: 1 }}>
      {/* Header */}
      <AppBar position="static" sx={{ bgcolor: config.color || 'primary.main' }}>
        <Toolbar>
          <IconButton
            edge="start"
            color="inherit"
            onClick={() => navigate('/')}
            sx={{ mr: 2 }}
          >
            <ArrowBack />
          </IconButton>
          {config.icon}
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            {config.title}
          </Typography>
          <Chip label={`Domain: ${config.domain}`} color="default" size="small" variant="outlined" />
        </Toolbar>
      </AppBar>

      <Box sx={{ p: 3 }}>
        {/* Navigation Tabs */}
        <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
          <Tabs value={tabValue} onChange={(e, newValue) => setTabValue(newValue)}>
            <Tab label="Search & Query" />
            <Tab label="Content Library" />
            <Tab label="Upload Content" />
            {/* Add Keywords & Definitions tab if supported */}
            {config.dataTypes?.supportsKeywords && (
              <Tab label="Keywords & Definitions" />
            )}
            {config.additionalTabs?.map((tabName, index) => (
              <Tab key={index} label={tabName} />
            ))}
          </Tabs>
        </Box>

        {/* Search Panel */}
        <TabPanel value={tabValue} index={0}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Search {config.domain.charAt(0).toUpperCase() + config.domain.slice(1)} Content
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                Search through {config.domain} content and get intelligent responses
              </Typography>
              
              <TextField
                fullWidth
                label={`Search ${config.domain} content...`}
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                sx={{ mb: 2 }}
                placeholder={config.searchPlaceholder}
              />
              
              <Button 
                variant="contained" 
                onClick={handleSearch}
                disabled={loading}
                sx={{ mb: 2 }}
              >
                Search Content
              </Button>

              {loading && <CircularProgress />}
              {error && <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>}
              
              {searchResults && (
                <Box sx={{ mt: 3 }}>
                  <Typography variant="h6" gutterBottom>Search Results</Typography>
                  
                  {searchResults?.ai_response?.type && (
                    <Alert 
                      severity={
                        searchResults.ai_response.type === "error" ? "error" : 
                        searchResults.ai_response.type === "not_found" ? "warning" : 
                        "info"
                      } 
                      sx={{ mb: 2 }}
                    >
                      <Typography variant="subtitle1">
                        {searchResults.ai_response.type.replace(/_/g, ' ').replace(/\b\w/g, (c: string) => c.toUpperCase())}
                      </Typography>
                      {searchResults.ai_response.content && (
                        <Box sx={{ mt: 1, mb: 1 }}>
                          <Typography variant="body2" style={{ whiteSpace: 'pre-line' }}>
                            {searchResults.ai_response.content}
                          </Typography>
                        </Box>
                      )}
                    </Alert>
                  )}
                  
                  <pre style={{ 
                    backgroundColor: '#f5f5f5', 
                    padding: '10px', 
                    borderRadius: '4px',
                    overflow: 'auto',
                    maxHeight: '400px'
                  }}>
                    {JSON.stringify(searchResults, null, 2)}
                  </pre>
                </Box>
              )}
            </CardContent>
          </Card>
        </TabPanel>

        {/* Content Library Tab */}
        <TabPanel value={tabValue} index={1}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                {config.title} Content Library
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                Browse uploaded content, materials, and terminology
              </Typography>
              
              {contentLoading ? (
                <Box sx={{ display: 'flex', justifyContent: 'center', my: 4 }}>
                  <CircularProgress />
                </Box>
              ) : (
                <>
                  {renderLoadingStatus()}
                  
                  {error && <Alert severity="error" sx={{ mb: 3 }}>{error}</Alert>}
                  
                  {/* Chapters Section */}
                  {config.apiEndpoints.chapters && (
                    <>
                      <Typography variant="h6" sx={{ mt: 3, mb: 2 }}>
                        Uploaded Chapters ({chapters.length})
                      </Typography>
                      
                      {chapters.length === 0 ? (
                        <Alert severity="info">No chapters have been uploaded yet.</Alert>
                      ) : (
                        <Box sx={{ mb: 4 }}>
                          {chapters.map((chapter, index) => (
                            <Card key={index} sx={{ mb: 2, bgcolor: 'background.paper' }}>
                              <CardContent>
                                <Typography variant="subtitle1">
                                  Chapter {chapter.chapter_number}: {chapter.chapter_title}
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                  Book: {chapter.book_title}
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                  Contains {chapter.case_studies_count || 0} case studies and {chapter.medical_terms_count || 0} terms
                                </Typography>
                              </CardContent>
                            </Card>
                          ))}
                        </Box>
                      )}
                    </>
                  )}
                  
                  {/* Case Studies Section */}
                  {config.apiEndpoints.caseStudies && (
                    <>
                      <Typography variant="h6" sx={{ mt: 4, mb: 2 }}>
                        Case Studies ({caseStudies.length})
                      </Typography>
                      
                      {caseStudies.length === 0 ? (
                        <Alert severity="info">No case studies have been extracted yet.</Alert>
                      ) : (
                        <Box sx={{ mb: 4 }}>
                          {caseStudies.map((caseStudy, index) => (
                            <Accordion key={index} sx={{ mb: 2 }}>
                              <AccordionSummary expandIcon={<ExpandMore />}>
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, width: '100%' }}>
                                  <Assignment color="primary" />
                                  <Typography variant="h6">
                                    Case {caseStudy.case_number}
                                  </Typography>
                                  {caseStudy.patient_info?.age && (
                                    <Chip 
                                      label={`${caseStudy.patient_info.age} years, ${caseStudy.patient_info.sex || caseStudy.patient_info.gender || 'Unknown'}`}
                                      size="small"
                                    />
                                  )}
                                </Box>
                              </AccordionSummary>
                              
                              <AccordionDetails>
                                <Grid container spacing={3}>
                                  <Grid item xs={12} md={6}>
                                    <Typography variant="subtitle2" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                      <Person />
                                      Patient Information
                                    </Typography>
                                    {caseStudy.patient_info?.presenting_symptoms?.length ? (
                                      <Box>
                                        <Typography variant="body2" gutterBottom>
                                          <strong>Presenting Symptoms:</strong>
                                        </Typography>
                                        <Box component="ul" sx={{ pl: 2, mt: 0 }}>
                                          {caseStudy.patient_info.presenting_symptoms.map((symptom: any, idx: any) => (
                                            <li key={idx}>
                                              <Typography variant="body2">{symptom}</Typography>
                                            </li>
                                          ))}
                                        </Box>
                                      </Box>
                                    ) : null}
                                  </Grid>
                                  
                                  <Grid item xs={12} md={6}>
                                    <Typography variant="subtitle2" gutterBottom>
                                      Clinical Findings
                                    </Typography>
                                    {caseStudy.clinical_findings?.length ? (
                                      <Box component="ul" sx={{ pl: 2, mt: 0 }}>
                                        {caseStudy.clinical_findings.map((finding: any, idx: any) => (
                                          <li key={idx}>
                                            <Typography variant="body2">{finding}</Typography>
                                          </li>
                                        ))}
                                      </Box>
                                    ) : (
                                      <Typography variant="body2" color="text.secondary">
                                        No clinical findings extracted
                                      </Typography>
                                    )}
                                  </Grid>
                                  
                                  <Grid item xs={12}>
                                    <Divider sx={{ my: 2 }} />
                                    <Typography variant="subtitle2" gutterBottom>
                                      Case Content
                                    </Typography>
                                    <Paper sx={{ p: 2, bgcolor: 'grey.50' }}>
                                      <Typography variant="body2" style={{ whiteSpace: 'pre-line' }}>
                                        {caseStudy.content}
                                      </Typography>
                                    </Paper>
                                  </Grid>
                                </Grid>
                              </AccordionDetails>
                            </Accordion>
                          ))}
                        </Box>
                      )}
                    </>
                  )}
                  
                  {/* Terms Section */}
                  {config.apiEndpoints.terms && (
                    <>
                      <Typography variant="h6" sx={{ mt: 4, mb: 2 }}>
                        Medical Terms ({medicalTerms.length})
                      </Typography>
                      
                      {medicalTerms.length === 0 ? (
                        <Alert severity="info">No terms have been extracted yet.</Alert>
                      ) : (
                        <Grid container spacing={2}>
                          {medicalTerms.map((term, index) => (
                            <Grid item xs={12} sm={6} md={4} key={index}>
                              <Card sx={{ height: '100%' }}>
                                <CardContent>
                                  <Typography variant="subtitle1">
                                    {term.term}
                                  </Typography>
                                  <Typography variant="caption" color="text.secondary">
                                    {term.category || 'General'}
                                  </Typography>
                                  <Typography variant="body2" sx={{ mt: 1 }}>
                                    {term.context && term.context[0]?.substring(0, 100)}
                                    {term.context && term.context[0]?.length > 100 ? '...' : ''}
                                  </Typography>
                                </CardContent>
                              </Card>
                            </Grid>
                          ))}
                        </Grid>
                      )}
                    </>
                  )}
                </>
              )}
              
              {/* Refresh button */}
              <Box sx={{ mt: 3, textAlign: 'center' }}>
                <Button 
                  variant="outlined" 
                  startIcon={<Science />} 
                  onClick={fetchLibraryContent}
                  disabled={contentLoading}
                >
                  Refresh Content Library
                </Button>
              </Box>
            </CardContent>
          </Card>
        </TabPanel>

        {/* Upload Content Tab */}
        <TabPanel value={tabValue} index={2}>
          <Card>
            <CardContent>
              {config.uploadComponent ? (
                config.uploadComponent
              ) : (
                <>
                  <Typography variant="h6" gutterBottom>
                    Upload {config.domain.charAt(0).toUpperCase() + config.domain.slice(1)} Content
                  </Typography>
                  <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                    To upload {config.domain} content, please use the Upload feature from the Home page 
                    or visit the dedicated Upload component.
                  </Typography>
                  
                  {config.apiEndpoints.upload && (
                    <Button 
                      variant="contained" 
                      onClick={() => navigate(`/${config.domain}-upload`)}
                      sx={{ mb: 2 }}
                    >
                      Go to {config.domain.charAt(0).toUpperCase() + config.domain.slice(1)} Upload
                    </Button>
                  )}
                </>
              )}
            </CardContent>
          </Card>
        </TabPanel>

        {/* Add CSV Upload Tab for Keywords & Definitions */}
        {config.dataTypes?.supportsKeywords && (
          <TabPanel value={tabValue} index={3}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  Upload Keywords & Definitions for {config.title}
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                  Paste CSV content with keywords and definitions. Format: Keyword,Definition
                </Typography>
                
                <TextField
                  fullWidth
                  multiline
                  rows={10}
                  label="Paste CSV Content"
                  value={csvContent}
                  onChange={(e) => setCsvContent(e.target.value)}
                  placeholder="keyword1,definition1&#10;keyword2,definition2&#10;..."
                  sx={{ mb: 2 }}
                />
                
                <Button 
                  variant="contained" 
                  onClick={handleCsvUpload}
                  disabled={csvUploading}
                  sx={{ mb: 2 }}
                >
                  Upload Keywords & Definitions
                </Button>

                {csvUploading && <CircularProgress />}
                {error && <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>}
                
                {csvUploadResult && (
                  <Box sx={{ mt: 3 }}>
                    <Typography variant="h6" gutterBottom>Upload Results</Typography>
                    <Alert severity="success" sx={{ mb: 2 }}>
                      Successfully uploaded {csvUploadResult.count || 0} keywords and definitions
                    </Alert>
                    <Paper sx={{ p: 2, bgcolor: 'grey.50' }}>
                      <pre style={{ overflow: 'auto' }}>
                        {JSON.stringify(csvUploadResult, null, 2)}
                      </pre>
                    </Paper>
                  </Box>
                )}

                <Box sx={{ mt: 3 }}>
                  <Typography variant="subtitle2" gutterBottom>CSV Format Example:</Typography>
                  <Paper sx={{ p: 2, bgcolor: 'grey.100' }}>
                    <Typography variant="body2" component="pre" sx={{ fontFamily: 'monospace' }}>
                      keyword,definition
                      {'\n'}antibody,A protective protein produced by the immune system
                      {'\n'}antigen,A substance that triggers an immune response
                    </Typography>
                  </Paper>
                </Box>
              </CardContent>
            </Card>
          </TabPanel>
        )}

        {/* Additional Tabs */}
        {config.additionalTabs?.map((_, index) => (
          <TabPanel key={index} value={tabValue} index={config.dataTypes?.supportsKeywords ? index + 4 : index + 3}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  {config.additionalTabs?.[index]}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  This tab is configured but not fully implemented yet.
                </Typography>
              </CardContent>
            </Card>
          </TabPanel>
        ))}
      </Box>
    </Box>
  );
};

export default GenericStudyArea;