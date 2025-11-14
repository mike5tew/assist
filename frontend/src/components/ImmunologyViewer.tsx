import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Tabs,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  TextField,
  InputAdornment,
  IconButton,
  Alert,
  CircularProgress,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Grid,
  Divider
} from '@mui/material';
import {
  Search,
  ExpandMore,
  Science,
  Book,
  Person,
  MedicalServices,
  Assignment
} from '@mui/icons-material';
import { apiClient } from '../config/api';

interface Chapter {
  id: string;
  chapter_number: string;
  chapter_title: string;
  book_title: string;
  isbn: string;
  case_studies_count: number;
  medical_terms_count: number;
  extracted_at: string;
}

interface CaseStudy {
  id: string;
  case_number: string;
  content: string;
  patient_info: {
    age?: string;
    sex?: string;
    presenting_symptoms: string[];
  };
  clinical_findings: string[];
  discussion_points: string[];
}

interface MedicalTerm {
  term: string;
  category: string;
  frequency: number;
  context: string[];
  definition?: string;
}

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

const TabPanel: React.FC<TabPanelProps> = ({ children, value, index }) => (
  <div hidden={value !== index}>
    {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
  </div>
);

const ImmunologyViewer: React.FC = () => {
  const [tabValue, setTabValue] = useState(0);
  const [chapters, setChapters] = useState<Chapter[]>([]);
  const [caseStudies, setCaseStudies] = useState<CaseStudy[]>([]);
  const [medicalTerms, setMedicalTerms] = useState<MedicalTerm[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [error, setError] = useState<string | null>(null);

  const categories = ['all', 'cell_marker', 'cytokine', 'immunoglobulin', 'mhc_molecule', 'general_immunology'];

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    setError(null);

    try {
      // Fetch chapters - fixed API call to use consistent pattern
      const chaptersData = await apiClient.get<{chapters: Chapter[]}>('/api/immunology/chapters');
      setChapters(chaptersData.chapters || []);

      // Fetch case studies - fixed API call
      const casesData = await apiClient.get<{case_studies: CaseStudy[]}>('/api/immunology/case-studies');
      setCaseStudies(casesData.case_studies || []);

      // Fetch medical terms - fixed API call
      const termsData = await apiClient.get<{medical_terms: MedicalTerm[]}>('/api/immunology/medical-terms');
      setMedicalTerms(termsData.medical_terms || []);

    } catch (err) {
      setError('Failed to load immunology data');
      console.error('Error fetching immunology data:', err);
    } finally {
      setLoading(false);
    }
  };

  const filteredTerms = medicalTerms.filter(term => {
    const matchesSearch = term.term.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         term.context.some(ctx => ctx.toLowerCase().includes(searchTerm.toLowerCase()));
    const matchesCategory = selectedCategory === 'all' || term.category === selectedCategory;
    return matchesSearch && matchesCategory;
  });

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 400 }}>
        <CircularProgress />
        <Typography sx={{ ml: 2 }}>Loading immunology data...</Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ maxWidth: 1200, mx: 'auto', p: 3 }}>
      <Typography variant="h4" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <Science color="primary" />
        Immunology Content Viewer
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      <Tabs value={tabValue} onChange={handleTabChange} sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tab icon={<Book />} label="Chapters" />
        <Tab icon={<Assignment />} label="Case Studies" />
        <Tab icon={<MedicalServices />} label="Medical Terms" />
      </Tabs>

      {/* Chapters Tab */}
      <TabPanel value={tabValue} index={0}>
        <Typography variant="h6" gutterBottom>
          Processed Chapters
        </Typography>
        
        {chapters.length === 0 ? (
          <Alert severity="info">
            No chapters have been processed yet. Use the upload component to add immunology content.
          </Alert>
        ) : (
          <Grid container spacing={2}>
            {chapters.map((chapter) => (
              <Grid item xs={12} md={6} key={chapter.id}>
                <Card>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Chapter {chapter.chapter_number}: {chapter.chapter_title}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                      {chapter.book_title}
                    </Typography>
                    <Typography variant="body2" gutterBottom>
                      ISBN: {chapter.isbn}
                    </Typography>
                    
                    <Box sx={{ mt: 2, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                      <Chip 
                        label={`${chapter.case_studies_count} Case Studies`} 
                        color="primary" 
                        size="small" 
                      />
                      <Chip 
                        label={`${chapter.medical_terms_count} Medical Terms`} 
                        color="secondary" 
                        size="small" 
                      />
                    </Box>
                    
                    <Typography variant="caption" display="block" sx={{ mt: 1 }}>
                      Processed: {new Date(chapter.extracted_at).toLocaleDateString()}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        )}
      </TabPanel>

      {/* Case Studies Tab */}
      <TabPanel value={tabValue} index={1}>
        <Typography variant="h6" gutterBottom>
          Clinical Case Studies
        </Typography>
        
        {caseStudies.length === 0 ? (
          <Alert severity="info">
            No case studies have been extracted yet.
          </Alert>
        ) : (
          <Box>
            {caseStudies.map((caseStudy) => (
              <Accordion key={caseStudy.id} sx={{ mb: 2 }}>
                <AccordionSummary expandIcon={<ExpandMore />}>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, width: '100%' }}>
                    <Assignment color="primary" />
                    <Typography variant="h6">
                      Case {caseStudy.case_number}
                    </Typography>
                    {caseStudy.patient_info.age && (
                      <Chip 
                        label={`${caseStudy.patient_info.age} years, ${caseStudy.patient_info.sex}`}
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
                      {caseStudy.patient_info.presenting_symptoms.length > 0 && (
                        <Box>
                          <Typography variant="body2" gutterBottom>
                            <strong>Presenting Symptoms:</strong>
                          </Typography>
                          <Box component="ul" sx={{ pl: 2, mt: 0 }}>
                            {caseStudy.patient_info.presenting_symptoms.map((symptom, index) => (
                              <li key={index}>
                                <Typography variant="body2">{symptom}</Typography>
                              </li>
                            ))}
                          </Box>
                        </Box>
                      )}
                    </Grid>
                    
                    <Grid item xs={12} md={6}>
                      <Typography variant="subtitle2" gutterBottom>
                        Clinical Findings
                      </Typography>
                      {caseStudy.clinical_findings.length > 0 ? (
                        <Box component="ul" sx={{ pl: 2, mt: 0 }}>
                          {caseStudy.clinical_findings.map((finding, index) => (
                            <li key={index}>
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
                        <Typography variant="body2">
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
      </TabPanel>

      {/* Medical Terms Tab */}
      <TabPanel value={tabValue} index={2}>
        <Box sx={{ mb: 3 }}>
          <Typography variant="h6" gutterBottom>
            Medical Terminology
          </Typography>
          
          <Grid container spacing={2} sx={{ mb: 2 }}>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                placeholder="Search medical terms..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <Search />
                    </InputAdornment>
                  ),
                }}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                select
                label="Category"
                value={selectedCategory}
                onChange={(e) => setSelectedCategory(e.target.value)}
                SelectProps={{ native: true }}
              >
                {categories.map((category) => (
                  <option key={category} value={category}>
                    {category === 'all' ? 'All Categories' : category.replace('_', ' ').toUpperCase()}
                  </option>
                ))}
              </TextField>
            </Grid>
          </Grid>
        </Box>

        {filteredTerms.length === 0 ? (
          <Alert severity="info">
            No medical terms found matching your criteria.
          </Alert>
        ) : (
          <TableContainer component={Paper}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell><strong>Term</strong></TableCell>
                  <TableCell><strong>Category</strong></TableCell>
                  <TableCell><strong>Frequency</strong></TableCell>
                  <TableCell><strong>Context</strong></TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {filteredTerms.map((term, index) => (
                  <TableRow key={index}>
                    <TableCell>
                      <Typography variant="body2" sx={{ fontWeight: 'medium' }}>
                        {term.term}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Chip 
                        label={term.category.replace('_', ' ')} 
                        size="small" 
                        color="primary"
                        variant="outlined"
                      />
                    </TableCell>
                    <TableCell>
                      <Chip label={term.frequency} size="small" />
                    </TableCell>
                    <TableCell>
                      <Box sx={{ maxWidth: 300 }}>
                        {term.context.slice(0, 2).map((ctx, ctxIndex) => (
                          <Typography 
                            key={ctxIndex} 
                            variant="caption" 
                            display="block"
                            sx={{ 
                              overflow: 'hidden',
                              textOverflow: 'ellipsis',
                              whiteSpace: 'nowrap'
                            }}
                          >
                            "{ctx}"
                          </Typography>
                        ))}
                        {term.context.length > 2 && (
                          <Typography variant="caption" color="text.secondary">
                            +{term.context.length - 2} more contexts
                          </Typography>
                        )}
                      </Box>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </TabPanel>
    </Box>
  );
};

export default ImmunologyViewer;
   