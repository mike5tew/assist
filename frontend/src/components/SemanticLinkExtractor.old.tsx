import React, { useState, useRef } from 'react';
import {
  Box,
  Button,
  Card,
  CardContent,
  Container,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Divider,
  FormControl,
  Grid,
  InputLabel,
  LinearProgress,
  MenuItem,
  Paper,
  Select,
  Stack,
  TextField,
  Typography,
  Alert,
  Chip,
  SelectChangeEvent,
  IconButton,
} from '@mui/material';
import {
  CloudUpload,
  CheckCircle,
  Error as ErrorIcon,
  Clear as ClearIcon,
  Add as AddIcon,
  Delete as DeleteIcon,
} from '@mui/icons-material';
import PDFViewer from './SemanticLinkExtractor/PDFViewer';
import TermSelectionPanel from './SemanticLinkExtractor/TermSelectionPanel';
import RelationshipSelector from './SemanticLinkExtractor/RelationshipSelector';
import ExtractionResults from './SemanticLinkExtractor/ExtractionResults';
import { SemanticLinkExtractionService } from '../services/semanticLinkService';
import {
  SemanticLinkSelection,
  RelationshipType,
  ExtractionResult,
  QualityFlag,
  LinkCondition,
  RELATIONSHIP_PAIRS,
  CONDITION_RELATIONS,
} from '../types/semanticLinks';

// Local type for accumulated links before batch submission
interface ExtractedLink {
  statement: string;
  sourceTerm: string;
  targetTerm: string;
  forwardRelation: string;
  inverseRelation: string;
  conditions: LinkCondition[];
  confidence?: number;
  qualityFlag?: QualityFlag;
  position?: { start: number; end: number };
}

/**
 * SemanticLinkExtractor - Main component for the semantic link curation tool
 * 
 * Workflow for creating training data (Golden Set):
 * 1. Highlight statement in PDF (the full vectorizable fact)
 * 2. Tag source and target terms within the statement
 * 3. Select forward relation (inverse auto-populates)
 * 4. Optionally add conditions with their own relations
 */
export default function SemanticLinkExtractor() {
  // Read optional default source from URL query param `source`
  const [searchParams] = React.useState(() => new URLSearchParams(window.location.search));
  const initialSourceEncoded = searchParams.get('source') || '';
  const initialSource = initialSourceEncoded ? JSON.parse(decodeURIComponent(initialSourceEncoded)) : null;
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [pdfFile, setPdfFile] = useState<File | null>(null);
  const [documentId, setDocumentId] = useState<string>(initialSource?.batch_id || initialSource?.source_id || '');
  const [documentTitle, setDocumentTitle] = useState<string>(initialSource?.title || '');
  const [selectedText, setSelectedText] = useState<string>('');
  const [showSourceDialog, setShowSourceDialog] = useState<boolean>(false);
  const [availableSources, setAvailableSources] = useState<any[]>([]);
  const [sourceLoading, setSourceLoading] = useState(false);
  
  // Statement-centric workflow state
  const [statement, setStatement] = useState<string>(''); // The full vectorizable fact
  const [sourceTerm, setSourceTerm] = useState<string>('');
  const [targetTerm, setTargetTerm] = useState<string>('');
  const [forwardRelation, setForwardRelation] = useState<string>('causes');
  const [inverseRelation, setInverseRelation] = useState<string>('is_caused_by');
  const [conditions, setConditions] = useState<LinkCondition[]>([]);
  const [qualityFlag, setQualityFlag] = useState<QualityFlag | null>(null);
  const [confidence, setConfidence] = useState<number>(0.95);

  // UI state
  const [isLoading, setIsLoading] = useState(false);
  const [extractedLinks, setExtractedLinks] = useState<ExtractedLink[]>([]);
  const [results, setResults] = useState<ExtractionResult | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [showConfirmDialog, setShowConfirmDialog] = useState(false);
  const [currentStep, setCurrentStep] = useState<'select_statement' | 'tag_terms' | 'select_relation'>('select_statement');

  // Condition dialog state
  const [showConditionDialog, setShowConditionDialog] = useState(false);
  const [newConditionTerm, setNewConditionTerm] = useState('');
  const [newConditionForward, setNewConditionForward] = useState('occurs_when');
  const [newConditionInverse, setNewConditionInverse] = useState('occurs_when');

  const service = SemanticLinkExtractionService.getInstance();

  // Handle PDF file selection
  const handleFileSelect = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    if (!file.type.includes('pdf') && !file.type.includes('image')) {
      setError('Please select a valid PDF or image file');
      return;
    }

    setPdfFile(file);
    setDocumentId(`doc_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`);
    setError(null);
    setExtractedLinks([]);
  };

  // Handle text selection from PDF viewer - captures the statement
  const handleTermSelected = (term: string, context: string, position?: { start: number; end: number }) => {
    if (currentStep === 'select_statement') {
      // First selection captures the full statement
      setStatement(context || term);
      setCurrentStep('tag_terms');
    } else if (currentStep === 'tag_terms') {
      // Tag source and target terms
      if (!sourceTerm) {
        setSourceTerm(term);
      } else if (!targetTerm) {
        setTargetTerm(term);
        setCurrentStep('select_relation');
      }
    }
  };

  // Handle relation change with auto-population of inverse
  const handleRelationChange = (forward: string, inverse: string) => {
    setForwardRelation(forward);
    setInverseRelation(inverse);
  };

  // Add a condition to the current link
  const handleAddCondition = () => {
    if (!newConditionTerm.trim()) return;
    setConditions([...conditions, {
      term: newConditionTerm.trim(),
      forwardRelation: newConditionForward,
      inverseRelation: newConditionInverse,
    }]);
    setNewConditionTerm('');
    setNewConditionForward('occurs_when');
    setNewConditionInverse('occurs_when');
    setShowConditionDialog(false);
  };

  // Remove a condition
  const handleRemoveCondition = (index: number) => {
    setConditions(conditions.filter((_, i) => i !== index));
  };

  // Reset term selection
  const resetTermSelection = () => {
    setStatement('');
    setSourceTerm('');
    setTargetTerm('');
    setForwardRelation('causes');
    setInverseRelation('is_caused_by');
    setConditions([]);
    setQualityFlag(null);
    setCurrentStep('select_statement');
  };

  // Submit semantic link extraction
  const handleSubmitLink = async () => {
    if (!statement || !sourceTerm || !targetTerm || !forwardRelation) {
      setError('Please capture statement, tag both terms, and select a relationship');
      return;
    }

    const newLink: ExtractedLink = {
      statement,
      sourceTerm,
      targetTerm,
      forwardRelation,
      inverseRelation,
      conditions: [...conditions],
      confidence,
      qualityFlag: qualityFlag || undefined,
    };

    setExtractedLinks([...extractedLinks, newLink]);
    resetTermSelection();
  };

  // Batch submit all extracted links
  const handleBatchSubmit = async () => {
    if (extractedLinks.length === 0) {
      setError('No links to submit');
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const selections: SemanticLinkSelection[] = extractedLinks.map((link) => ({
        statement: link.statement,
        sourceTerm: link.sourceTerm,
        targetTerm: link.targetTerm,
        forwardRelation: link.forwardRelation,
        inverseRelation: link.inverseRelation,
        conditions: link.conditions,
        confidence: link.confidence,
        qualityFlag: link.qualityFlag,
      }));

      const result = await service.extractSemanticLinks(documentId, selections);
      setResults(result);
      setExtractedLinks([]);
    } catch (err) {
      setError(`Failed to extract semantic links: ${err instanceof Error ? err.message : 'Unknown error'}`);
    } finally {
      setIsLoading(false);
    }
  };

  // Clear all state
  const handleReset = () => {
    setPdfFile(null);
    setDocumentId('');
    setExtractedLinks([]);
    setResults(null);
    setError(null);
    resetTermSelection();
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const fetchSources = async () => {
    setSourceLoading(true);
    try {
      const { sources } = await (await fetch('/esp-organizer/api/sources?type=book')).json();
      setAvailableSources(sources || []);
    } catch (err) {
      console.error('Failed to fetch sources', err);
    } finally {
      setSourceLoading(false);
    }
  };

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h3" gutterBottom sx={{ mb: 3, fontWeight: 'bold' }}>
        🔗 Semantic Link Extraction Tool
      </Typography>
      <Typography variant="body1" sx={{ mb: 4, color: 'text.secondary' }}>
        Extract meaningful relationships from educational texts. Select terms in the document and define
        their relationships to build the knowledge graph.
      </Typography>

      <Grid container spacing={3}>
        {/* Left Panel: PDF Viewer */}
        <Grid item xs={12} md={6}>
          <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            <CardContent sx={{ flex: 1 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                <Typography variant="h6" gutterBottom>
                  📄 Document Viewer
                </Typography>
                <Chip label={documentTitle ? `Source: ${documentTitle}` : (documentId ? `Source ID: ${documentId}` : 'No source selected')} />
                <Button size="small" onClick={() => { setShowSourceDialog(true); }}>Change source</Button>
              </Box>

              {!pdfFile ? (
                <Box
                  sx={{
                    border: '2px dashed #ccc',
                    borderRadius: 2,
                    p: 4,
                    textAlign: 'center',
                    cursor: 'pointer',
                    transition: 'all 0.3s',
                    '&:hover': { borderColor: '#1976d2', backgroundColor: '#f5f5f5' },
                  }}
                  onClick={() => fileInputRef.current?.click()}
                >
                  <CloudUpload sx={{ fontSize: 48, color: '#999', mb: 2 }} />
                  <Typography variant="body1">Click to upload PDF or image</Typography>
                  <Typography variant="caption" sx={{ color: 'text.secondary' }}>
                    Supported: PDF, JPG, PNG
                  </Typography>
                </Box>
              ) : (
                <Box>
                  <PDFViewer
                    file={pdfFile}
                    onTermSelected={handleTermSelected}
                    isSelecting={currentStep === 'select_statement' || currentStep === 'tag_terms'}
                  />
                  <Button
                    fullWidth
                    variant="text"
                    startIcon={<ClearIcon />}
                    onClick={() => {
                      setPdfFile(null);
                      setDocumentId('');
                    }}
                    sx={{ mt: 2 }}
                  >
                    Remove Document
                  </Button>
                </Box>
              )}

              <input
                ref={fileInputRef}
                type="file"
                hidden
                accept=".pdf,image/*"
                onChange={handleFileSelect}
              />
            </CardContent>
          </Card>
        </Grid>

        {/* Right Panel: Statement & Term Selection */}
        <Grid item xs={12} md={6}>
          <Stack spacing={2} sx={{ height: '100%' }}>
            {/* Statement Capture Card */}
            <Card sx={{ backgroundColor: currentStep === 'select_statement' ? '#fff3e0' : undefined }}>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  1️⃣ Capture Statement
                </Typography>
                <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mb: 1 }}>
                  Highlight the complete fact from the document (this will be vectorized)
                </Typography>
                {statement ? (
                  <Box sx={{ p: 2, backgroundColor: '#e8f5e9', borderRadius: 1, position: 'relative' }}>
                    <Typography variant="body2" sx={{ fontStyle: 'italic' }}>"{statement}"</Typography>
                    <IconButton 
                      size="small" 
                      onClick={() => { setStatement(''); setCurrentStep('select_statement'); }}
                      sx={{ position: 'absolute', top: 4, right: 4 }}
                    >
                      <ClearIcon fontSize="small" />
                    </IconButton>
                  </Box>
                ) : (
                  <Typography variant="body2" color="text.secondary">
                    Select text in the PDF to capture the statement...
                  </Typography>
                )}
              </CardContent>
            </Card>

            {/* Source Term Card */}
            <Card sx={{ backgroundColor: currentStep === 'tag_terms' && !sourceTerm ? '#fff3e0' : undefined }}>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  2️⃣ Tag Source Term
                </Typography>
                <TermSelectionPanel
                  term={sourceTerm}
                  context={statement}
                  onClear={() => setSourceTerm('')}
                  step="source"
                  disabled={!statement}
                />
              </CardContent>
            </Card>

            {/* Target Term Card */}
            <Card sx={{ backgroundColor: currentStep === 'tag_terms' && sourceTerm && !targetTerm ? '#fff3e0' : undefined }}>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  3️⃣ Tag Target Term
                </Typography>
                <TermSelectionPanel
                  term={targetTerm}
                  context={statement}
                  onClear={() => setTargetTerm('')}
                  step="target"
                  disabled={!sourceTerm}
                />
              </CardContent>
            </Card>

            {/* Relationship Type Card */}
            {sourceTerm && targetTerm && (
              <Card sx={{ backgroundColor: '#f0f7ff' }}>
                <CardContent>
                  <Typography variant="h6" gutterBottom>
                    4️⃣ Select Relationship
                  </Typography>
                  <RelationshipSelector
                    forwardRelation={forwardRelation}
                    inverseRelation={inverseRelation}
                    qualityFlag={qualityFlag}
                    confidence={confidence}
                    onForwardRelationChange={handleRelationChange}
                    onQualityFlagChange={(value: QualityFlag | null) => setQualityFlag(value)}
                    onConfidenceChange={(value: number) => setConfidence(value)}
                  />
                </CardContent>
              </Card>
            )}

            {/* Conditions Card */}
            {sourceTerm && targetTerm && (
              <Card>
                <CardContent>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
                    <Typography variant="h6">
                      5️⃣ Conditions (Optional)
                    </Typography>
                    <Button 
                      size="small" 
                      startIcon={<AddIcon />}
                      onClick={() => setShowConditionDialog(true)}
                    >
                      Add Condition
                    </Button>
                  </Box>
                  <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mb: 1 }}>
                    When does this relationship hold? (e.g., "occurs when maternal antibodies wane")
                  </Typography>
                  {conditions.length === 0 ? (
                    <Typography variant="body2" color="text.secondary">No conditions added</Typography>
                  ) : (
                    <Stack spacing={1}>
                      {conditions.map((cond, idx) => (
                        <Box key={idx} sx={{ display: 'flex', alignItems: 'center', gap: 1, p: 1, backgroundColor: '#f5f5f5', borderRadius: 1 }}>
                          <Chip label={cond.forwardRelation.replace(/_/g, ' ')} size="small" color="primary" variant="outlined" />
                          <Typography variant="body2" sx={{ flex: 1 }}>{cond.term}</Typography>
                          <IconButton size="small" onClick={() => handleRemoveCondition(idx)}>
                            <DeleteIcon fontSize="small" />
                          </IconButton>
                        </Box>
                      ))}
                    </Stack>
                  )}
                </CardContent>
              </Card>
            )}

            {/* Quality Flag - Optional */}
            {sourceTerm && targetTerm && (
              <Card>
                <CardContent>
                  <Typography variant="body2" sx={{ fontWeight: 'bold', mb: 2 }}>
                    6️⃣ Quality Assessment (Optional)
                  </Typography>
                  <FormControl fullWidth size="small">
                    <InputLabel>Quality Flag</InputLabel>
                    <Select
                      value={qualityFlag || ''}
                      onChange={(e) => setQualityFlag((e.target.value as QualityFlag) || null)}
                      label="Quality Flag"
                    >
                      <MenuItem value="">No flag</MenuItem>
                      <MenuItem value="high_confidence">High confidence</MenuItem>
                      <MenuItem value="needs_verification">Needs verification</MenuItem>
                      <MenuItem value="problematic">Problematic</MenuItem>
                    </Select>
                  </FormControl>
                  <TextField
                    fullWidth
                    type="number"
                    label="Confidence (0-1)"
                    value={confidence}
                    onChange={(e) => setConfidence(parseFloat(e.target.value))}
                    inputProps={{ step: 0.05, min: 0, max: 1 }}
                    size="small"
                    sx={{ mt: 2 }}
                  />
                </CardContent>
              </Card>
            )}

            {/* Action Buttons */}
            <Stack direction="row" spacing={1}>
              <Button
                variant="contained"
                fullWidth
                disabled={!statement || !sourceTerm || !targetTerm}
                onClick={handleSubmitLink}
              >
                Add Link
              </Button>
              <Button variant="outlined" fullWidth onClick={resetTermSelection}>
                Clear Selection
              </Button>
            </Stack>
          </Stack>
        </Grid>

        {/* Extracted Links Summary */}
        {extractedLinks.length > 0 && (
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  📋 Extracted Links ({extractedLinks.length})
                </Typography>
                <Box sx={{ maxHeight: 300, overflowY: 'auto', mb: 2 }}>
                  {extractedLinks.map((link, idx) => (
                    <Box key={idx} sx={{ p: 1.5, mb: 1, backgroundColor: '#f5f5f5', borderRadius: 1, borderLeft: '4px solid #1976d2' }}>
                      <Typography variant="body2" sx={{ fontStyle: 'italic', mb: 1 }}>
                        "{link.statement}"
                      </Typography>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, flexWrap: 'wrap' }}>
                        <Chip label={link.sourceTerm} size="small" color="primary" />
                        <Typography variant="caption">→ {link.forwardRelation.replace(/_/g, ' ')} →</Typography>
                        <Chip label={link.targetTerm} size="small" color="secondary" />
                        {link.conditions.length > 0 && (
                          <Chip 
                            label={`${link.conditions.length} condition(s)`} 
                            size="small" 
                            variant="outlined" 
                          />
                        )}
                        <Button
                          size="small"
                          variant="text"
                          color="error"
                          onClick={() => setExtractedLinks(extractedLinks.filter((_, i) => i !== idx))}
                        >
                          Remove
                        </Button>
                      </Box>
                    </Box>
                  ))}
                </Box>
                <Stack direction="row" spacing={1}>
                  <Button
                    variant="contained"
                    fullWidth
                    onClick={handleBatchSubmit}
                    disabled={isLoading}
                  >
                    {isLoading ? 'Submitting...' : 'Submit All Links'}
                  </Button>
                  <Button variant="outlined" fullWidth onClick={() => setExtractedLinks([])}>
                    Clear All
                  </Button>
                </Stack>
              </CardContent>
            </Card>
          </Grid>
        )}

        {/* Results */}
        {results && (
          <Grid item xs={12}>
            <ExtractionResults results={results} onReset={handleReset} />
          </Grid>
        )}

        {/* Error Alert */}
        {error && (
          <Grid item xs={12}>
            <Alert severity="error" onClose={() => setError(null)}>
              {error}
            </Alert>
          </Grid>
        )}

        {/* Loading Progress */}
        {isLoading && <LinearProgress sx={{ mt: 2 }} />}
      </Grid>

      {/* Source selection dialog */}
      <Dialog open={showSourceDialog} onClose={() => setShowSourceDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Select Source</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
            <Button onClick={fetchSources} disabled={sourceLoading}>Load Books</Button>
            <Button onClick={() => {
              // populate with last upload if available
              const enc = window.localStorage.getItem('last_upload_source');
              if (enc) {
                try {
                  const obj = JSON.parse(decodeURIComponent(enc));
                  setDocumentId(obj.batch_id || obj.source_id || '');
                  setDocumentTitle(obj.title || '');
                  setShowSourceDialog(false);
                } catch (err) {
                  console.error('Failed to parse last upload source', err);
                }
              }
            }}>Use last upload</Button>
          </Box>
          <Box sx={{ maxHeight: 320, overflow: 'auto' }}>
            {availableSources.length === 0 ? (
              <Typography variant="body2" sx={{ color: 'text.secondary' }}>No sources loaded — click "Load Books"</Typography>
            ) : (
              availableSources.map((s: any) => (
                <Box key={s.id || s._id} sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', py: 1, borderBottom: '1px solid #eee' }}>
                  <Box>
                    <Typography variant="subtitle2">{s.title}</Typography>
                    <Typography variant="caption" sx={{ color: 'text.secondary' }}>{(s.authors || []).join(', ')}</Typography>
                  </Box>
                  <Box>
                    <Button onClick={() => { setDocumentId(s.id || s._id || ''); setDocumentTitle(s.title || ''); setShowSourceDialog(false); }} size="small">Select</Button>
                  </Box>
                </Box>
              ))
            )}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowSourceDialog(false)}>Close</Button>
        </DialogActions>
      </Dialog>

      {/* Confirm Overwrite Dialog */}
      <Dialog open={showConfirmDialog} onClose={() => setShowConfirmDialog(false)}>
        <DialogTitle>Clear Existing Selection?</DialogTitle>
        <DialogContent>
          <Typography>
            You already have both terms selected. Would you like to clear and start a new selection?
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowConfirmDialog(false)}>Keep</Button>
          <Button
            onClick={() => {
              resetTermSelection();
              setShowConfirmDialog(false);
            }}
            variant="contained"
          >
            Clear & Start Over
          </Button>
        </DialogActions>
      </Dialog>

      {/* Add Condition Dialog */}
      <Dialog open={showConditionDialog} onClose={() => setShowConditionDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Add Condition</DialogTitle>
        <DialogContent>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            Under what condition does this relationship hold?
          </Typography>
          <TextField
            label="Condition Term"
            value={newConditionTerm}
            onChange={(e) => setNewConditionTerm(e.target.value)}
            fullWidth
            sx={{ mb: 2 }}
            placeholder="e.g., maternal antibody waning"
          />
          <FormControl fullWidth sx={{ mb: 2 }}>
            <InputLabel>Relation to Link</InputLabel>
            <Select
              value={newConditionForward}
              onChange={(e) => {
                const forward = e.target.value;
                setNewConditionForward(forward);
                // Auto-populate inverse
                const inverse = CONDITION_RELATIONS[forward] || forward;
                setNewConditionInverse(inverse);
              }}
              label="Relation to Link"
            >
              {Object.entries(CONDITION_RELATIONS).map(([key, inverse]) => (
                <MenuItem key={key} value={key}>
                  {key.replace(/_/g, ' ')} {key !== inverse ? `(inverse: ${inverse.replace(/_/g, ' ')})` : '(symmetric)'}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <Box sx={{ p: 1.5, backgroundColor: '#e3f2fd', borderRadius: 1 }}>
            <Typography variant="caption">
              This relationship <strong>{newConditionForward.replace(/_/g, ' ')}</strong> "{newConditionTerm || '...'}"
            </Typography>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowConditionDialog(false)}>Cancel</Button>
          <Button onClick={handleAddCondition} variant="contained" disabled={!newConditionTerm.trim()}>
            Add Condition
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
}
