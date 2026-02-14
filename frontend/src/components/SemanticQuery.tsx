import React, { useEffect, useMemo, useState } from 'react';
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
  Divider,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  FormControlLabel,
  Switch,
  Stack,
  Tooltip
} from '@mui/material';
import {
  Search,
  ExpandMore,
  Psychology,
  Analytics,
  School,
  Science,
  Info
} from '@mui/icons-material';
import { Link, useSearchParams } from 'react-router-dom';
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

type ResultKind = 'exact' | 'semantic' | 'related';

type SelectedResult =
  | { kind: 'exact'; item: Skill }
  | { kind: 'related'; item: Skill }
  | { kind: 'semantic'; item: SemanticMatch };

const snippet = (text: string | undefined, max = 160) => {
  if (!text) return '';
  const cleaned = text.replace(/\s+/g, ' ').trim();
  if (cleaned.length <= max) return cleaned;
  return `${cleaned.slice(0, max).trim()}…`;
};

const relevanceLabel = (value: number) => `${(value * 100).toFixed(1)}%`;

const SemanticQuery: React.FC = () => {
  const [query, setQuery] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<SemanticQueryResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  const [selected, setSelected] = useState<SelectedResult | null>(null);
  const [showDebug, setShowDebug] = useState(false);
  const [autoRanForQuery, setAutoRanForQuery] = useState<string | null>(null);
  const [searchParams, setSearchParams] = useSearchParams();

  const exampleGroups = useMemo(
    () => [
      {
        label: 'Skills & Learning (safe defaults)',
        queries: [
          'working memory skills',
          'reading comprehension',
          'fine motor development',
          'problem solving abilities',
          'communication skills',
          'emotional intelligence',
        ],
      },
      {
        label: 'Quality of Connection (CHISG-style prompts)',
        queries: [
          'What skills are prerequisites for reading comprehension?',
          'What is often confused with working memory (but is different)?',
          'Show related skills for fine motor development, but avoid weak links',
          'List nearby skills for reading comprehension and explain why they belong',
        ],
      },
    ],
    []
  );

  useEffect(() => {
    // Support shareable links like /semantic-query?q=working%20memory
    const q = searchParams.get('q');
    if (q && q !== query) {
      setQuery(q);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchParams]);

  const qParam = searchParams.get('q') || '';
  const runParam = searchParams.get('run') === '1';

  useEffect(() => {
    // Optional auto-run for shareable demo links:
    // /semantic-query?q=working%20memory%20skills&run=1
    if (!runParam) return;
    if (!qParam) return;
    if (query !== qParam) return; // wait until state catches up
    if (loading) return;
    if (autoRanForQuery === qParam) return;

    setAutoRanForQuery(qParam);
    void handleSearch();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [runParam, qParam, query, loading, autoRanForQuery]);

  const handleSearch = async () => {
    if (!query.trim()) return;

    setLoading(true);
    setError(null);

    try {
      // Updated to use the correct apiClient pattern - it now returns data directly
      const data = await apiClient.post<SemanticQueryResponse>('/api/skills/semantic-query', { query });
      setResult(data);
      setSearchParams((prev) => {
        const next = new URLSearchParams(prev);
        next.set('q', query);
        return next;
      });
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

  const openDetails = (kind: ResultKind, item: Skill | SemanticMatch) => {
    if (kind === 'semantic') {
      setSelected({ kind: 'semantic', item: item as SemanticMatch });
      return;
    }
    if (kind === 'related') {
      setSelected({ kind: 'related', item: item as Skill });
      return;
    }
    setSelected({ kind: 'exact', item: item as Skill });
  };

  const closeDetails = () => setSelected(null);

  return (
    <Box sx={{ maxWidth: 1200, mx: 'auto', p: 3 }}>
      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={1} justifyContent="space-between" alignItems={{ xs: 'flex-start', sm: 'center' }} sx={{ mb: 1 }}>
        <Typography variant="h4" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Psychology color="primary" />
          Semantic Query Demo
        </Typography>
        <Stack direction="row" spacing={1}>
          <Button component={Link} to="/chisg" variant="outlined" size="small">
            Back to CHISG
          </Button>
          <Button component={Link} to="/ai-chat" variant="outlined" size="small">
            Grounded Chat
          </Button>
        </Stack>
      </Stack>

      <Typography variant="body1" paragraph color="text.secondary" sx={{ maxWidth: 980 }}>
        This demo compares <strong>exact matching</strong> (MongoDB) with <strong>semantic retrieval</strong> (Weaviate vector search).
        The intent is to show how stronger models can “join dots” too aggressively — and why the quality of connections matters.
      </Typography>

      <Alert severity="info" sx={{ mb: 3 }}>
        Current data note: this environment is primarily backed by skills content. If you load additional corpora later, we can add
        domain-specific example packs without changing the UI.
      </Alert>

      {/* Search Interface */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} alignItems={{ xs: 'stretch', md: 'flex-end' }}>
            <Box sx={{ flexGrow: 1 }}>
              <Typography variant="h6" gutterBottom>
                Search Query
              </Typography>
              <TextField
                fullWidth
                placeholder="Enter your query (e.g., prerequisites for reading comprehension, fine motor development…)"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                onKeyPress={handleKeyPress}
              />
            </Box>

            <Stack direction="row" spacing={1}>
              <Button
                variant="contained"
                onClick={handleSearch}
                disabled={loading || !query.trim()}
                startIcon={loading ? <CircularProgress size={20} /> : <Search />}
                sx={{ minWidth: 140 }}
              >
                {loading ? 'Searching…' : 'Search'}
              </Button>
              <Button
                variant="outlined"
                onClick={() => {
                  setQuery('');
                  setResult(null);
                  setError(null);
                  setSearchParams((prev) => {
                    const next = new URLSearchParams(prev);
                    next.delete('q');
                    next.delete('run');
                    return next;
                  });
                }}
                disabled={loading}
              >
                Clear
              </Button>
            </Stack>
          </Stack>

          <Box sx={{ mt: 1 }}>
            <Typography variant="caption" color="text.secondary">
              Tip: share a link like <strong>/semantic-query?q=working%20memory%20skills&amp;run=1</strong> to auto-run on open.
            </Typography>
          </Box>
          
          <Divider sx={{ my: 2 }} />

          <Typography variant="subtitle2" gutterBottom>
            Example Queries
          </Typography>
          <Grid container spacing={2}>
            {exampleGroups.map((group) => (
              <Grid item xs={12} md={4} key={group.label}>
                <Paper variant="outlined" sx={{ p: 1.5, height: '100%' }}>
                  <Typography variant="caption" color="text.secondary" display="block" sx={{ mb: 1 }}>
                    {group.label}
                  </Typography>
                  <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                    {group.queries.map((q) => (
                      <Chip
                        key={q}
                        label={q}
                        variant="outlined"
                        onClick={() => handleExampleClick(q)}
                        sx={{ cursor: 'pointer' }}
                      />
                    ))}
                  </Box>
                </Paper>
              </Grid>
            ))}
          </Grid>
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
          <Card sx={{ mb: 3 }}>
            <CardContent>
              <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} justifyContent="space-between" alignItems={{ xs: 'flex-start', md: 'center' }}>
                <Box>
                  <Typography variant="h6" gutterBottom>
                    Query
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    <strong>{result.query}</strong> • {new Date(result.timestamp).toLocaleString()}
                  </Typography>
                </Box>
                <Tooltip title="The synthesis is a convenience summary, not a source. Always check the items and sources below.">
                  <Chip icon={<Info />} label="Synthesis is not a source" variant="outlined" />
                </Tooltip>
              </Stack>
            </CardContent>
          </Card>

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
                        <TableRow key={index} hover sx={{ cursor: 'pointer' }} onClick={() => openDetails('exact', skill)}>
                          <TableCell>{skill.skill_name}</TableCell>
                          <TableCell>
                            <Chip label={skill.skill_type} size="small" />
                          </TableCell>
                          <TableCell>{skill.development_age}+</TableCell>
                          <TableCell sx={{ maxWidth: 300 }}>{snippet(skill.description)}</TableCell>
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
                        <TableRow key={index} hover sx={{ cursor: 'pointer' }} onClick={() => openDetails('semantic', match)}>
                          <TableCell>{match.skill_name}</TableCell>
                          <TableCell>
                            <Chip 
                              label={relevanceLabel(match.semantic_relevance)}
                              color={match.semantic_relevance > 0.8 ? 'success' : 
                                     match.semantic_relevance > 0.6 ? 'warning' : 'default'}
                              size="small"
                            />
                          </TableCell>
                          <TableCell>
                            <Chip label={match.skill_type} size="small" variant="outlined" />
                          </TableCell>
                          <TableCell sx={{ maxWidth: 300 }}>
                            {snippet(match.description)}
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

          {/* Related Skills */}
          {result.related_skills.length > 0 && (
            <Accordion sx={{ mb: 2 }}>
              <AccordionSummary expandIcon={<ExpandMore />}>
                <Typography variant="h6" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <Psychology color="warning" />
                  Related Skills ({result.related_skills.length})
                </Typography>
              </AccordionSummary>
              <AccordionDetails>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                  These are related concepts returned by the system’s relationship expansion. The goal is to surface helpful neighbors without “everything connects to everything”.
                </Typography>
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
                      {result.related_skills.map((skill, index) => (
                        <TableRow key={index} hover sx={{ cursor: 'pointer' }} onClick={() => openDetails('related', skill)}>
                          <TableCell>{skill.skill_name}</TableCell>
                          <TableCell>
                            <Chip label={skill.skill_type} size="small" variant="outlined" />
                          </TableCell>
                          <TableCell>{skill.development_age}+</TableCell>
                          <TableCell sx={{ maxWidth: 300 }}>{snippet(skill.description)}</TableCell>
                          <TableCell>
                            {skill.source_info ? (
                              <Box>
                                <Typography variant="caption" display="block">
                                  {skill.source_info.title}
                                </Typography>
                                {skill.source_info.chapter_title && (
                                  <Typography variant="caption" color="text.secondary">
                                    {skill.source_info.chapter_title}
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

          {/* LLM Synthesis */}
          {result.llm_synthesis && (
            <Card sx={{ mb: 2 }}>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  Grounded Synthesis
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                  This summary is generated for convenience. Treat it as a draft interpretation of the retrieved items, not as evidence.
                </Typography>
                <Typography variant="body1" sx={{ fontStyle: 'italic' }}>
                  {result.llm_synthesis}
                </Typography>
              </CardContent>
            </Card>
          )}

          {/* Technical Details */}
          <Card sx={{ mb: 2 }}>
            <CardContent>
              <FormControlLabel
                control={<Switch checked={showDebug} onChange={(e) => setShowDebug(e.target.checked)} />}
                label="Show debug / raw response (may expose implementation details)"
              />
              {showDebug && (
                <Box sx={{ bgcolor: 'grey.100', p: 2, borderRadius: 1, mt: 1 }}>
                  <Typography variant="body2" component="pre" sx={{ whiteSpace: 'pre-wrap', fontSize: '0.8rem' }}>
                    {JSON.stringify(result, null, 2)}
                  </Typography>
                </Box>
              )}
            </CardContent>
          </Card>
        </Box>
      )}

      {/* Details dialog */}
      <Dialog open={Boolean(selected)} onClose={closeDetails} maxWidth="md" fullWidth>
        {selected && (
          <>
            <DialogTitle>
              {selected.kind === 'semantic' ? 'Semantic Match Details' : selected.kind === 'related' ? 'Related Skill Details' : 'Exact Match Details'}
            </DialogTitle>
            <DialogContent dividers>
              <Stack spacing={1.25}>
                <Typography variant="h6">
                  {selected.item.skill_name}
                </Typography>

                {'semantic_relevance' in selected.item && (
                  <Box>
                    <Typography variant="body2" color="text.secondary">
                      Semantic relevance
                    </Typography>
                    <Chip
                      label={relevanceLabel(selected.item.semantic_relevance)}
                      color={selected.item.semantic_relevance > 0.8 ? 'success' : selected.item.semantic_relevance > 0.6 ? 'warning' : 'default'}
                      size="small"
                    />
                  </Box>
                )}

                <Box>
                  <Typography variant="body2" color="text.secondary">
                    Type
                  </Typography>
                  <Chip label={selected.item.skill_type} size="small" variant="outlined" />
                </Box>

                {'development_age' in selected.item && (
                  <Box>
                    <Typography variant="body2" color="text.secondary">
                      Development age
                    </Typography>
                    <Typography variant="body1">{selected.item.development_age}+</Typography>
                  </Box>
                )}

                <Box>
                  <Typography variant="body2" color="text.secondary">
                    Description
                  </Typography>
                  <Typography variant="body1">{selected.item.description}</Typography>
                </Box>

                {'criteria' in selected.item && selected.item.criteria?.length > 0 && (
                  <Box>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                      Criteria (summary)
                    </Typography>
                    <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                      {selected.item.criteria.slice(0, 12).map((c) => (
                        <Chip key={c} label={c} size="small" />
                      ))}
                    </Box>
                    {selected.item.criteria.length > 12 && (
                      <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
                        Showing first 12 criteria.
                      </Typography>
                    )}
                  </Box>
                )}

                <Divider />

                <Box>
                  <Typography variant="body2" color="text.secondary" gutterBottom>
                    Source
                  </Typography>
                  {selected.item.source_info ? (
                    <Stack spacing={0.5}>
                      <Typography variant="body2">
                        <strong>{selected.item.source_info.title}</strong>
                      </Typography>
                      {selected.item.source_info.type && (
                        <Typography variant="caption" color="text.secondary">
                          Type: {selected.item.source_info.type}
                        </Typography>
                      )}
                      {selected.item.source_info.isbn && (
                        <Typography variant="caption" color="text.secondary">
                          ISBN: {selected.item.source_info.isbn}
                        </Typography>
                      )}
                      {(selected.item.source_info.chapter_number || selected.item.source_info.chapter_title) && (
                        <Typography variant="caption" color="text.secondary">
                          {selected.item.source_info.chapter_number ? `Chapter ${selected.item.source_info.chapter_number}` : ''}
                          {selected.item.source_info.chapter_title ? ` • ${selected.item.source_info.chapter_title}` : ''}
                        </Typography>
                      )}
                      <Typography variant="caption" color="text.secondary">
                        Extraction: {selected.item.source_info.extraction_method}
                      </Typography>
                    </Stack>
                  ) : (
                    <Typography variant="body2" color="text.secondary">
                      ESP Skills Database
                    </Typography>
                  )}
                </Box>
              </Stack>
            </DialogContent>
            <DialogActions>
              <Button onClick={closeDetails}>Close</Button>
            </DialogActions>
          </>
        )}
      </Dialog>
    </Box>
  );
};

export default SemanticQuery;
   