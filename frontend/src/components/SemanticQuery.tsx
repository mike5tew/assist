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
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  Grid,
  Divider
} from '@mui/material';
import {
  Search,
  ExpandMore,
  Psychology,
  Analytics,
  School,
  Science
} from '@mui/icons-material';
import { apiClient } from '../config/api';

interface SemanticQueryResponse {
  query: string;
  timestamp: string;
  results_summary: {
    total_exact_matches: number;
    total_semantic_matches: number;
    total_related_skills: number;
  };
  exact_matches: Skill[];
  semantic_matches: SemanticMatch[];
  related_skills: Skill[];
  llm_synthesis: string;
  semantic_insights: {
    query_complexity: string;
    confidence_levels: {
      high_confidence: number;
      medium_confidence: number;
      low_confidence: number;
    };
  };
}

interface Skill {
  skill_name: string;
  description: string;
  skill_type: string;
  development_age: number;
  criteria_levels: number;
  criteria: string[];
  source_info?: SourceInfo;
}

interface SemanticMatch {
  skill_name: string;
  description: string;
  skill_type: string;
  semantic_relevance: number;
  source_info?: SourceInfo;
}

interface SourceInfo {
  title: string;
  isbn?: string;
  type: string;
  chapter_number?: string;
  chapter_title?: string;
  extraction_method: string;
  extracted_at?: string;
}

const SemanticQuery: React.FC = () => {
  const [query, setQuery] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<SemanticQueryResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  const exampleQueries = [
    "working memory skills",
    "fine motor development",
    "reading comprehension",
    "emotional intelligence",
    "problem solving abilities",
    "communication skills",
    "immunology case studies",
    "T cell activation"
  ];

  const handleSearch = async () => {
    if (!query.trim()) return;

    setLoading(true);
    setError(null);

    try {
      // Updated to use the correct apiClient pattern - it now returns data directly
      const data = await apiClient.post<SemanticQueryResponse>('/api/skills/semantic-query', { query });
      setResult(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const handleKeyPress = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter') {
      handleSearch();
    }
  };

  const handleExampleClick = (exampleQuery: string) => {
    setQuery(exampleQuery);
  };

  return (
    <Box sx={{ maxWidth: 1200, mx: 'auto', p: 3 }}>
      <Typography variant="h4" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <Psychology color="primary" />
        Semantic Skills Query System
      </Typography>

      <Typography variant="body1" paragraph color="text.secondary">
        Test the semantic search system that powers the AI assistant. This interface directly queries 
        both MongoDB (exact matches) and Weaviate (vector similarity) to find related educational skills 
        and medical content.
      </Typography>

      {/* Search Interface */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Search Query
          </Typography>
          
          <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
            <TextField
              fullWidth
              placeholder="Enter your search query (e.g., 'working memory', 'immunology', 'fine motor skills')"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              onKeyPress={handleKeyPress}
            />
            <Button
              variant="contained"
              onClick={handleSearch}
              disabled={loading || !query.trim()}
              startIcon={loading ? <CircularProgress size={20} /> : <Search />}
              sx={{ minWidth: 120 }}
            >
              {loading ? 'Searching...' : 'Search'}
            </Button>
          </Box>

          <Typography variant="subtitle2" gutterBottom>
            Example Queries:
          </Typography>
          <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
            {exampleQueries.map((example, index) => (
              <Chip
                key={index}
                label={example}
                variant="outlined"
                onClick={() => handleExampleClick(example)}
                sx={{ cursor: 'pointer' }}
              />
            ))}
          </Box>
        </CardContent>
      </Card>

      {/* Error Display */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          Error: {error}
        </Alert>
      )}

      {/* Results Display */}
      {result && (
        <Box>
          {/* Results Summary */}
          <Card sx={{ mb: 3 }}>
            <CardContent>
              <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Analytics />
                Search Results Summary
              </Typography>
              
              <Grid container spacing={2}>
                <Grid item xs={12} md={4}>
                  <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'success.light', color: 'success.contrastText' }}>
                    <Typography variant="h4">{result.results_summary.total_exact_matches}</Typography>
                    <Typography variant="body2">Exact Matches (MongoDB)</Typography>
                  </Paper>
                </Grid>
                <Grid item xs={12} md={4}>
                  <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'info.light', color: 'info.contrastText' }}>
                    <Typography variant="h4">{result.results_summary.total_semantic_matches}</Typography>
                    <Typography variant="body2">Semantic Matches (Weaviate)</Typography>
                  </Paper>
                </Grid>
                <Grid item xs={12} md={4}>
                  <Paper sx={{ p: 2, textAlign: 'center', bgcolor: 'warning.light', color: 'warning.contrastText' }}>
                    <Typography variant="h4">{result.results_summary.total_related_skills}</Typography>
                    <Typography variant="body2">Related Skills</Typography>
                  </Paper>
                </Grid>
              </Grid>

              <Box sx={{ mt: 2, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                <Chip 
                  label={`${result.semantic_insights.confidence_levels.high_confidence} High Confidence`}
                  color="success"
                  size="small"
                />
                <Chip 
                  label={`${result.semantic_insights.confidence_levels.medium_confidence} Medium Confidence`}
                  color="warning"
                  size="small"
                />
                <Chip 
                  label={`${result.semantic_insights.confidence_levels.low_confidence} Low Confidence`}
                  color="default"
                  size="small"
                />
                <Chip 
                  label={`Complexity: ${result.semantic_insights.query_complexity}`}
                  variant="outlined"
                  size="small"
                />
              </Box>
            </CardContent>
          </Card>

          {/* Exact Matches */}
          {result.exact_matches.length > 0 && (
            <Accordion defaultExpanded sx={{ mb: 2 }}>
              <AccordionSummary expandIcon={<ExpandMore />}>
                <Typography variant="h6" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <School color="success" />
                  Exact Matches ({result.exact_matches.length}) - MongoDB
                </Typography>
              </AccordionSummary>
              <AccordionDetails>
                <TableContainer component={Paper}>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell><strong>Skill Name</strong></TableCell>
                        <TableCell><strong>Type</strong></TableCell>
                        <TableCell><strong>Age</strong></TableCell>
                        <TableCell><strong>Description</strong></TableCell>
                        <TableCell><strong>Source</strong></TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {result.exact_matches.map((skill, index) => (
                        <TableRow key={index}>
                          <TableCell>{skill.skill_name}</TableCell>
                          <TableCell>
                            <Chip label={skill.skill_type} size="small" />
                          </TableCell>
                          <TableCell>{skill.development_age}+</TableCell>
                          <TableCell sx={{ maxWidth: 300 }}>
                            {skill.description.substring(0, 150)}...
                          </TableCell>
                          <TableCell>
                            {skill.source_info ? (
                              <Box>
                                <Typography variant="caption" display="block">
                                  {skill.source_info.title}
                                </Typography>
                                {skill.source_info.isbn && (
                                  <Typography variant="caption" color="text.secondary">
                                    ISBN: {skill.source_info.isbn}
                                  </Typography>
                                )}
                              </Box>
                            ) : (
                              <Typography variant="caption" color="text.secondary">
                                ESP Skills Database
                              </Typography>
                            )}
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </TableContainer>
              </AccordionDetails>
            </Accordion>
          )}

          {/* Semantic Matches */}
          {result.semantic_matches.length > 0 && (
            <Accordion defaultExpanded sx={{ mb: 2 }}>
              <AccordionSummary expandIcon={<ExpandMore />}>
                <Typography variant="h6" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <Science color="info" />
                  Semantic Matches ({result.semantic_matches.length}) - Weaviate Vector Search
                </Typography>
              </AccordionSummary>
              <AccordionDetails>
                <TableContainer component={Paper}>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell><strong>Skill Name</strong></TableCell>
                        <TableCell><strong>Relevance</strong></TableCell>
                        <TableCell><strong>Type</strong></TableCell>
                        <TableCell><strong>Description</strong></TableCell>
                        <TableCell><strong>Source</strong></TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {result.semantic_matches.map((match, index) => (
                        <TableRow key={index}>
                          <TableCell>{match.skill_name}</TableCell>
                          <TableCell>
                            <Chip 
                              label={`${(match.semantic_relevance * 100).toFixed(1)}%`}
                              color={match.semantic_relevance > 0.8 ? 'success' : 
                                     match.semantic_relevance > 0.6 ? 'warning' : 'default'}
                              size="small"
                            />
                          </TableCell>
                          <TableCell>
                            <Chip label={match.skill_type} size="small" variant="outlined" />
                          </TableCell>
                          <TableCell sx={{ maxWidth: 300 }}>
                            {match.description.substring(0, 150)}...
                          </TableCell>
                          <TableCell>
                            {match.source_info ? (
                              <Box>
                                <Typography variant="caption" display="block">
                                  {match.source_info.title}
                                </Typography>
                                {match.source_info.chapter_number && (
                                  <Typography variant="caption" color="text.secondary">
                                    Chapter {match.source_info.chapter_number}
                                  </Typography>
                                )}
                                <Chip 
                                  label={match.source_info.extraction_method}
                                  size="small"
                                  variant="outlined"
                                  sx={{ mt: 0.5 }}
                                />
                              </Box>
                            ) : (
                              <Typography variant="caption" color="text.secondary">
                                ESP Skills Database
                              </Typography>
                            )}
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </TableContainer>
              </AccordionDetails>
            </Accordion>
          )}

          {/* LLM Synthesis */}
          {result.llm_synthesis && (
            <Card sx={{ mb: 2 }}>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  AI Synthesis
                </Typography>
                <Typography variant="body1" sx={{ fontStyle: 'italic' }}>
                  {result.llm_synthesis}
                </Typography>
              </CardContent>
            </Card>
          )}

          {/* Technical Details */}
          <Accordion>
            <AccordionSummary expandIcon={<ExpandMore />}>
              <Typography variant="h6">
                Technical Details & Raw Response
              </Typography>
            </AccordionSummary>
            <AccordionDetails>
              <Box sx={{ bgcolor: 'grey.100', p: 2, borderRadius: 1 }}>
                <Typography variant="body2" component="pre" sx={{ whiteSpace: 'pre-wrap', fontSize: '0.8rem' }}>
                  {JSON.stringify(result, null, 2)}
                </Typography>
              </Box>
            </AccordionDetails>
          </Accordion>
        </Box>
      )}
    </Box>
  );
};

export default SemanticQuery;
   