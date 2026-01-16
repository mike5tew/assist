import React, { useState, useRef, useEffect, useCallback } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Stack,
  Alert,
  Chip,
  IconButton,
  Autocomplete,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Paper,
  Divider,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  Tooltip,
  CircularProgress,
  Slider,
  Tabs,
  Tab,
} from '@mui/material';
import {
  CloudUpload,
  Clear as ClearIcon,
  Add as AddIcon,
  Delete as DeleteIcon,
  Save as SaveIcon,
  CheckCircle,
  SwapHoriz,
  ContentPaste,
} from '@mui/icons-material';
import PDFViewer from './PDFViewer';
import { Source } from './SourceRegister';
import { apiClient } from '../../config/api';

interface LinkType {
  id: string;
  link_a_to_b: string;
  link_b_to_a: string;
  category: string;
  usage_count: number;
}

interface PendingLink {
  id: string;
  entityA: string;
  entityB: string;
  linkAToB: string;
  linkBToA: string;
  chapter?: string;
  excerptText?: string;
  pageNumber?: number;
  qualityScore?: number;
  qualityReason?: string;
  contextLinkIds?: string[];
}

interface ExistingLink {
  id: string;
  statement: string;
  source_term: string;
  target_term: string;
}

interface LinkBuilderProps {
  source: Source | null;
  onBack?: () => void;
}

type ActiveField = 'entityA' | 'entityB' | 'excerpt' | null;
type ViewTab = 'text' | 'pdf';

export default function LinkBuilder({ source, onBack }: LinkBuilderProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);
  
  // PDF and text state
  const [pdfFile, setPdfFile] = useState<File | null>(null);
  const [extractedText, setExtractedText] = useState<string>('');
  const [isExtractingText, setIsExtractingText] = useState(false);
  const [viewTab, setViewTab] = useState<ViewTab>('text');
  const [selectedText, setSelectedText] = useState<string>('');
  
  // Link entry fields
  const [entityA, setEntityA] = useState('');
  const [entityB, setEntityB] = useState('');
  const [linkAToB, setLinkAToB] = useState('');
  const [linkBToA, setLinkBToA] = useState('');
  const [chapter, setChapter] = useState('');
  const [excerptText, setExcerptText] = useState('');
  const [pageNumber, setPageNumber] = useState<number | undefined>();
  const [qualityScore, setQualityScore] = useState<number>(0.8);
  const [qualityReason, setQualityReason] = useState('');
  const [contextLinkIds, setContextLinkIds] = useState<string[]>([]);
  
  // Active field for click-to-populate
  const [activeField, setActiveField] = useState<ActiveField>(null);
  
  // Link types for autocomplete
  const [linkTypes, setLinkTypes] = useState<LinkType[]>([]);
  const [showAddLinkTypeDialog, setShowAddLinkTypeDialog] = useState(false);
  const [newLinkAToB, setNewLinkAToB] = useState('');
  const [newLinkBToA, setNewLinkBToA] = useState('');
  const [newCategory, setNewCategory] = useState('causal');
  
  // Existing links for context linking
  const [existingLinks, setExistingLinks] = useState<ExistingLink[]>([]);
  
  // Pending links queue
  const [pendingLinks, setPendingLinks] = useState<PendingLink[]>([]);
  
  // UI state
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  // Fetch link types on mount
  useEffect(() => {
    fetchLinkTypes();
    fetchExistingLinks();
  }, []);

  // Global Enter key handler for quick submission
  useEffect(() => {
    const handleGlobalKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Enter' && !e.shiftKey && source && entityA && entityB && linkAToB && linkBToA) {
        // Don't trigger if typing in a text field that uses Enter for newlines
        const target = e.target as HTMLElement;
        if (target.tagName === 'TEXTAREA') return;
        
        e.preventDefault();
        handleAddToPending();
      }
    };

    window.addEventListener('keydown', handleGlobalKeyDown);
    return () => window.removeEventListener('keydown', handleGlobalKeyDown);
  }, [source, entityA, entityB, linkAToB, linkBToA, chapter, excerptText, pageNumber, qualityScore, qualityReason, contextLinkIds]);

  const fetchLinkTypes = async () => {
    try {
      const response = await apiClient.get<{ link_types: LinkType[] }>('/api/link-types');
      setLinkTypes(response.link_types || []);
    } catch (err) {
      console.error('Failed to fetch link types:', err);
    }
  };

  const fetchExistingLinks = async () => {
    try {
      const sourceId = source?.id || source?._id;
      const url = sourceId 
        ? `/api/semantic-links?source_id=${sourceId}&limit=100`
        : '/api/semantic-links?limit=100';
      const response = await apiClient.get<{ links: ExistingLink[] }>(url);
      setExistingLinks(response.links || []);
    } catch (err) {
      console.error('Failed to fetch existing links:', err);
    }
  };

  // Handle PDF file selection
  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file && file.type === 'application/pdf') {
      setPdfFile(file);
      setError(null);
      
      // Extract text from PDF
      extractTextFromPdf(file);
    } else {
      setError('Please select a PDF file');
    }
    // Reset input so same file can be selected again
    event.target.value = '';
  };

  // Extract text from PDF using pdfjs
  const extractTextFromPdf = async (file: File) => {
    setIsExtractingText(true);
    setExtractedText('');
    
    try {
      const pdfjsLib = await import('pdfjs-dist');
      // Use CDN worker - must match installed package version
      const pdfjsVersion = '4.10.38';
      pdfjsLib.GlobalWorkerOptions.workerSrc = `https://unpkg.com/pdfjs-dist@${pdfjsVersion}/build/pdf.worker.min.mjs`;
      
      const arrayBuffer = await file.arrayBuffer();
      const pdf = await pdfjsLib.getDocument({ data: arrayBuffer }).promise;
      
      let fullText = '';
      for (let i = 1; i <= pdf.numPages; i++) {
        const page = await pdf.getPage(i);
        const textContent = await page.getTextContent();
        const pageText = textContent.items
          .map((item: any) => item.str)
          .join(' ');
        fullText += `\n--- Page ${i} ---\n${pageText}\n`;
      }
      
      setExtractedText(fullText);
    } catch (err) {
      console.error('Failed to extract text:', err);
      setError('Failed to extract text from PDF. Check console for details.');
    } finally {
      setIsExtractingText(false);
    }
  };

  // Handle text selection from extracted text
  const handleTextSelection = () => {
    const selection = window.getSelection();
    if (!selection || selection.isCollapsed) return;
    
    const selectedStr = selection.toString().trim();
    if (!selectedStr) return;
    
    // Get surrounding context
    const range = selection.getRangeAt(0);
    const container = range.commonAncestorContainer;
    const fullText = container.textContent || '';
    const start = Math.max(0, range.startOffset - 50);
    const end = Math.min(fullText.length, range.endOffset + 50);
    const context = fullText.slice(start, end);
    
    handleTermSelected(selectedStr, context);
  };

  // Handle text selection from PDF - this now populates the active field
  const handleTermSelected = useCallback((term: string, context: string, position?: { start: number; end: number }) => {
    setSelectedText(term);
    
    // If an entity field is active, populate it
    if (activeField === 'entityA') {
      setEntityA(term);
      setActiveField('entityB'); // Auto-advance to next field
    } else if (activeField === 'entityB') {
      setEntityB(term);
      setActiveField(null); // Done with entities
    } else if (activeField === 'excerpt') {
      setExcerptText(context || term);
      setActiveField(null);
    }
  }, [activeField]);

  // Handle keyboard events for field activation
  const handleKeyDown = useCallback((e: React.KeyboardEvent, field: ActiveField) => {
    if (e.key === 'Enter' && selectedText && activeField === field) {
      // Confirm the selection
      e.preventDefault();
      if (field === 'entityA') {
        setEntityA(selectedText);
        setActiveField('entityB');
      } else if (field === 'entityB') {
        setEntityB(selectedText);
        setActiveField(null);
      } else if (field === 'excerpt') {
        setExcerptText(selectedText);
        setActiveField(null);
      }
    }
  }, [selectedText, activeField]);

  // Handle link type selection
  const handleLinkTypeChange = (_: any, value: LinkType | string | null) => {
    if (typeof value === 'string') {
      // Custom input - open dialog to create new
      setNewLinkAToB(value);
      setShowAddLinkTypeDialog(true);
    } else if (value) {
      setLinkAToB(value.link_a_to_b);
      setLinkBToA(value.link_b_to_a);
    }
  };

  // Create new link type
  const handleAddLinkType = async () => {
    if (!newLinkAToB) return;
    
    try {
      const response = await apiClient.post<{ link_type: LinkType }>('/api/link-types', {
        link_a_to_b: newLinkAToB,
        link_b_to_a: newLinkBToA || undefined,
        category: newCategory,
      });
      
      if (response.link_type) {
        setLinkTypes([...linkTypes, response.link_type]);
        setLinkAToB(response.link_type.link_a_to_b);
        setLinkBToA(response.link_type.link_b_to_a);
      }
      
      setShowAddLinkTypeDialog(false);
      setNewLinkAToB('');
      setNewLinkBToA('');
    } catch (err: any) {
      setError(`Failed to create link type: ${err.message}`);
    }
  };

  // Add link to pending queue
  const handleAddToPending = () => {
    if (!source) {
      setError('Please select a source first from the Source Register tab');
      return;
    }
    if (!entityA || !entityB) {
      setError('Please fill in both Entity A and Entity B');
      return;
    }
    if (!linkAToB || !linkBToA) {
      setError('Please fill in both forward (A→B) and reverse (B←A) link types');
      return;
    }

    const newLink: PendingLink = {
      id: `pending_${Date.now()}`,
      entityA,
      entityB,
      linkAToB,
      linkBToA,
      chapter: chapter || undefined,
      excerptText: excerptText || undefined,
      pageNumber: pageNumber,
      qualityScore: qualityScore,
      qualityReason: qualityReason || undefined,
      contextLinkIds: contextLinkIds.length > 0 ? contextLinkIds : undefined,
    };

    setPendingLinks([...pendingLinks, newLink]);
    
    // Reset form for next link (keep some fields for convenience)
    setEntityA('');
    setEntityB('');
    // Keep link type, chapter, quality score for convenience
    setExcerptText('');
    setContextLinkIds([]);
    setActiveField('entityA');
    setSuccess('Link added to queue (Press Enter for quick add)');
    setTimeout(() => setSuccess(null), 2000);
  };

  // Remove link from pending queue
  const handleRemovePending = (id: string) => {
    setPendingLinks(pendingLinks.filter(l => l.id !== id));
  };

  // Submit all pending links
  const handleSubmitAll = async () => {
    if (pendingLinks.length === 0) {
      setError('No links to submit');
      return;
    }

    setSaving(true);
    setError(null);

    try {
      const payload = {
        document_id: source?.id || source?._id || 'manual',
        source_id: source?.id || source?._id,
        source_title: source?.title,
        domain: source?.domain || 'manual',
        subject: source?.subject,
        selections: pendingLinks.map(link => ({
          statement: `${link.entityA} ${link.linkAToB} ${link.entityB}`,
          source_term: link.entityA,
          target_term: link.entityB,
          forward_relation: link.linkAToB,
          inverse_relation: link.linkBToA,
          chapter: link.chapter,
          excerpt_text: link.excerptText,
          page_number: link.pageNumber,
          quality_score: link.qualityScore,
          quality_reason: link.qualityReason,
          context_link_ids: link.contextLinkIds,
        })),
      };

      await apiClient.post('/api/semantic-links/extract', payload);
      
      setPendingLinks([]);
      setSuccess(`Successfully saved ${payload.selections.length} links`);
    } catch (err: any) {
      setError(`Failed to save links: ${err.message}`);
    } finally {
      setSaving(false);
    }
  };

  const getSourceId = () => source?.id || source?._id || '';

  return (
    <Box sx={{ display: 'flex', height: 'calc(100vh - 200px)', gap: 2 }}>
      {/* Left Panel: Text/PDF Viewer with Tabs */}
      <Paper sx={{ flex: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
        {/* Header with source info and upload button */}
        <Box sx={{ p: 2, borderBottom: 1, borderColor: 'divider', display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Box>
            <Typography variant="h6" sx={{ mb: 0.5 }}>
              📄 Document Content
            </Typography>
            {source && (
              <Chip 
                label={`Source: ${source.title}`} 
                size="small" 
                color="primary" 
                variant="outlined"
              />
            )}
          </Box>
          <Button
            size="small"
            variant="outlined"
            startIcon={pdfFile ? <ClearIcon /> : <AddIcon />}
            onClick={() => {
              if (pdfFile) {
                setPdfFile(null);
                setExtractedText('');
              } else {
                fileInputRef.current?.click();
              }
            }}
          >
            {pdfFile ? 'Remove PDF' : 'Upload PDF'}
          </Button>
        </Box>

        {/* Text/PDF Tabs - only show when PDF is loaded */}
        {pdfFile && (
          <Tabs
            value={viewTab}
            onChange={(_, newValue) => setViewTab(newValue)}
            sx={{ borderBottom: 1, borderColor: 'divider', minHeight: 40 }}
          >
            <Tab 
              value="text" 
              label="📝 Extracted Text" 
              sx={{ textTransform: 'none' }}
            />
            <Tab 
              value="pdf" 
              label="📄 PDF View" 
              sx={{ textTransform: 'none' }}
            />
          </Tabs>
        )}
        
        <Box sx={{ flex: 1, overflow: 'auto', p: 2 }}>
          {!pdfFile ? (
            <Box
              sx={{
                height: '100%',
                border: '2px dashed #ccc',
                borderRadius: 2,
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                cursor: 'pointer',
                transition: 'all 0.3s',
                '&:hover': { borderColor: 'primary.main', backgroundColor: 'action.hover' },
              }}
              onClick={() => fileInputRef.current?.click()}
            >
              <CloudUpload sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
              <Typography variant="h6" color="text.secondary">
                Click to upload PDF
              </Typography>
              <Typography variant="caption" color="text.secondary">
                Text will be extracted for easy selection
              </Typography>
            </Box>
          ) : viewTab === 'text' ? (
            <Box sx={{ height: '100%' }}>
              {isExtractingText ? (
                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100%', flexDirection: 'column', gap: 2 }}>
                  <CircularProgress />
                  <Typography color="text.secondary">Extracting text from PDF...</Typography>
                </Box>
              ) : (
                <Box
                  onMouseUp={handleTextSelection}
                  sx={{
                    whiteSpace: 'pre-wrap',
                    fontFamily: 'monospace',
                    fontSize: '0.9rem',
                    lineHeight: 1.6,
                    p: 2,
                    bgcolor: activeField ? 'warning.light' : 'grey.50',
                    borderRadius: 1,
                    cursor: activeField ? 'text' : 'default',
                    userSelect: 'text',
                    minHeight: '100%',
                    transition: 'background-color 0.2s',
                  }}
                >
                  {activeField && (
                    <Alert severity="info" sx={{ mb: 2 }}>
                      Highlight text to populate <strong>{activeField === 'entityA' ? 'Entity A' : activeField === 'entityB' ? 'Entity B' : 'Excerpt'}</strong>
                    </Alert>
                  )}
                  {extractedText || 'No text extracted from PDF.'}
                </Box>
              )}
            </Box>
          ) : (
            <Box sx={{ height: '100%' }}>
              <PDFViewer
                file={pdfFile}
                onTermSelected={handleTermSelected}
                isSelecting={activeField !== null}
              />
            </Box>
          )}
        </Box>

        <input
          ref={fileInputRef}
          type="file"
          hidden
          accept=".pdf"
          onChange={handleFileSelect}
        />
      </Paper>

      {/* Right Panel: Link Entry */}
      <Paper sx={{ width: 400, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
        <Box sx={{ p: 2, borderBottom: 1, borderColor: 'divider' }}>
          <Typography variant="h6">🔗 Create Link</Typography>
          <Typography variant="caption" color="text.secondary">
            Click a field, then highlight text to populate
          </Typography>
        </Box>

        <Box sx={{ flex: 1, overflow: 'auto', p: 2 }}>
          {/* Alerts */}
          {error && (
            <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
              {error}
            </Alert>
          )}
          {success && (
            <Alert severity="success" sx={{ mb: 2 }} onClose={() => setSuccess(null)}>
              {success}
            </Alert>
          )}

          <Stack spacing={2}>
            {/* Entity A */}
            <TextField
              label="Entity A"
              value={entityA}
              onChange={(e) => setEntityA(e.target.value)}
              onFocus={() => setActiveField('entityA')}
              onKeyDown={(e) => handleKeyDown(e, 'entityA')}
              fullWidth
              size="small"
              placeholder="Click here, then highlight in PDF"
              sx={{
                '& .MuiOutlinedInput-root': {
                  backgroundColor: activeField === 'entityA' ? 'warning.light' : undefined,
                },
              }}
              InputProps={{
                endAdornment: activeField === 'entityA' && (
                  <Chip label="ACTIVE" size="small" color="warning" />
                ),
              }}
            />

            {/* Link Type with A→B */}
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <Autocomplete
                freeSolo
                options={linkTypes}
                getOptionLabel={(option) => 
                  typeof option === 'string' ? option : option.link_a_to_b
                }
                value={linkTypes.find(lt => lt.link_a_to_b === linkAToB) || null}
                onChange={handleLinkTypeChange}
                onInputChange={(_, value) => {
                  if (value && !linkTypes.find(lt => lt.link_a_to_b === value)) {
                    setLinkAToB(value);
                  }
                }}
                renderInput={(params) => (
                  <TextField
                    {...params}
                    label="A → B (link type)"
                    size="small"
                    placeholder="e.g., causes, activates"
                  />
                )}
                renderOption={(props, option) => (
                  <li {...props}>
                    <Box>
                      <Typography variant="body2">{option.link_a_to_b}</Typography>
                      <Typography variant="caption" color="text.secondary">
                        ↔ {option.link_b_to_a} ({option.category})
                      </Typography>
                    </Box>
                  </li>
                )}
                sx={{ flex: 1 }}
              />
              <Tooltip title="Add new link type">
                <IconButton 
                  size="small" 
                  onClick={() => setShowAddLinkTypeDialog(true)}
                >
                  <AddIcon />
                </IconButton>
              </Tooltip>
            </Box>

            {/* Entity B */}
            <TextField
              label="Entity B"
              value={entityB}
              onChange={(e) => setEntityB(e.target.value)}
              onFocus={() => setActiveField('entityB')}
              onKeyDown={(e) => handleKeyDown(e, 'entityB')}
              fullWidth
              size="small"
              placeholder="Click here, then highlight in PDF"
              sx={{
                '& .MuiOutlinedInput-root': {
                  backgroundColor: activeField === 'entityB' ? 'warning.light' : undefined,
                },
              }}
              InputProps={{
                endAdornment: activeField === 'entityB' && (
                  <Chip label="ACTIVE" size="small" color="warning" />
                ),
              }}
            />

            {/* B→A (inferred, editable) */}
            <TextField
              label="B ← A (reverse)"
              value={linkBToA}
              onChange={(e) => setLinkBToA(e.target.value)}
              fullWidth
              size="small"
              placeholder="Auto-inferred from link type"
            />

            <Divider />

            {/* Context Links - Display selected */}
            {contextLinkIds.length > 0 && (
              <Box>
                <Typography variant="caption" color="text.secondary">
                  Context Links ({contextLinkIds.length})
                </Typography>
                <List dense sx={{ py: 0 }}>
                  {existingLinks.filter(l => contextLinkIds.includes(l.id)).map(link => (
                    <ListItem key={link.id} sx={{ py: 0, px: 1 }}>
                      <ListItemText 
                        primary={
                          <Typography variant="body2" noWrap>
                            {link.statement || `${link.source_term} → ${link.target_term}`}
                          </Typography>
                        }
                      />
                      <ListItemSecondaryAction>
                        <IconButton 
                          size="small" 
                          onClick={() => setContextLinkIds(contextLinkIds.filter(id => id !== link.id))}
                        >
                          <ClearIcon fontSize="small" />
                        </IconButton>
                      </ListItemSecondaryAction>
                    </ListItem>
                  ))}
                </List>
              </Box>
            )}

            {/* Add Context Link */}
            <Autocomplete
              options={existingLinks.filter(l => !contextLinkIds.includes(l.id))}
              getOptionLabel={(option) => option.statement || `${option.source_term} → ${option.target_term}`}
              value={null}
              onChange={(_, value) => {
                if (value) {
                  setContextLinkIds([...contextLinkIds, value.id]);
                }
              }}
              renderInput={(params) => (
                <TextField
                  {...params}
                  label="Add Context Link"
                  size="small"
                  placeholder="Search for related links..."
                />
              )}
              renderOption={(props, option) => (
                <li {...props}>
                  <Typography variant="body2" noWrap>
                    {option.statement || `${option.source_term} → ${option.target_term}`}
                  </Typography>
                </li>
              )}
            />

            {/* Quality Score and Reason - side by side */}
            <Stack direction="row" spacing={1} alignItems="flex-start">
              <Box sx={{ flex: 1 }}>
                <Typography variant="caption" color="text.secondary">
                  Quality: {(qualityScore * 100).toFixed(0)}%
                </Typography>
                <Slider
                  value={qualityScore}
                  onChange={(_, value) => setQualityScore(value as number)}
                  min={0}
                  max={1}
                  step={0.1}
                  size="small"
                  valueLabelDisplay="auto"
                  valueLabelFormat={(v) => `${(v * 100).toFixed(0)}%`}
                />
              </Box>
              <TextField
                label="Reason"
                value={qualityReason}
                onChange={(e) => setQualityReason(e.target.value)}
                size="small"
                placeholder="Why?"
                sx={{ flex: 1 }}
              />
            </Stack>

            {/* Save Link Button */}
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={handleAddToPending}
              disabled={!source || !entityA || !entityB || !linkAToB || !linkBToA}
              fullWidth
            >
              Save Link (or press Enter)
            </Button>

            <Divider />

            {/* Section/Chapter */}
            <Stack direction="row" spacing={1}>
              <TextField
                label="Section / Chapter"
                value={chapter}
                onChange={(e) => setChapter(e.target.value)}
                size="small"
                placeholder="e.g., Chapter 5 - Immunity"
                sx={{ flex: 2 }}
              />
              <TextField
                label="Page"
                type="number"
                value={pageNumber || ''}
                onChange={(e) => setPageNumber(parseInt(e.target.value) || undefined)}
                size="small"
                sx={{ flex: 1 }}
              />
            </Stack>

            {/* Source Excerpt */}
            <TextField
              label="Source Excerpt"
              value={excerptText}
              onChange={(e) => setExcerptText(e.target.value)}
              onFocus={() => setActiveField('excerpt')}
              multiline
              rows={2}
              fullWidth
              size="small"
              placeholder="Highlight text in PDF to capture excerpt"
              sx={{
                '& .MuiOutlinedInput-root': {
                  backgroundColor: activeField === 'excerpt' ? 'warning.light' : undefined,
                },
              }}
            />
          </Stack>

          {/* Pending Links Queue for this section */}
          <Box sx={{ mt: 3 }}>
            <Typography variant="subtitle2" gutterBottom>
              Links for this Section ({pendingLinks.length})
            </Typography>
            {pendingLinks.length === 0 ? (
              <Typography variant="body2" color="text.secondary" sx={{ fontStyle: 'italic' }}>
                No links added yet. Create links above and they will appear here.
              </Typography>
            ) : (
              <List dense>
                {pendingLinks.map((link) => (
                  <ListItem key={link.id} sx={{ bgcolor: 'action.hover', borderRadius: 1, mb: 0.5 }}>
                    <ListItemText
                      primary={
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, flexWrap: 'wrap' }}>
                          <Typography variant="body2">
                            <strong>{link.entityA}</strong> → <em>{link.linkAToB}</em> → <strong>{link.entityB}</strong>
                          </Typography>
                          {link.qualityScore !== undefined && (
                            <Chip 
                              label={`${(link.qualityScore * 100).toFixed(0)}%`} 
                              size="small" 
                              color={link.qualityScore >= 0.7 ? 'success' : link.qualityScore >= 0.4 ? 'warning' : 'error'}
                            />
                          )}
                        </Box>
                      }
                      secondary={
                        <Typography variant="body2" color="text.secondary">
                          <strong>{link.entityB}</strong> → <em>{link.linkBToA}</em> → <strong>{link.entityA}</strong>
                          {link.contextLinkIds?.length ? ` • ${link.contextLinkIds.length} context` : ''}
                        </Typography>
                      }
                    />
                    <ListItemSecondaryAction>
                      <IconButton 
                        size="small" 
                        onClick={() => handleRemovePending(link.id)}
                      >
                        <DeleteIcon fontSize="small" />
                      </IconButton>
                    </ListItemSecondaryAction>
                  </ListItem>
                ))}
              </List>
            )}
          </Box>
        </Box>

        {/* Save All Links Button */}
        <Box sx={{ p: 2, borderTop: 1, borderColor: 'divider' }}>
          <Button
            variant="contained"
            color="success"
            startIcon={saving ? <CircularProgress size={20} /> : <SaveIcon />}
            onClick={handleSubmitAll}
            disabled={pendingLinks.length === 0 || saving}
            fullWidth
            size="large"
          >
            Save All Links ({pendingLinks.length})
          </Button>
        </Box>
      </Paper>

      {/* Add Link Type Dialog */}
      <Dialog open={showAddLinkTypeDialog} onClose={() => setShowAddLinkTypeDialog(false)}>
        <DialogTitle>Add New Link Type</DialogTitle>
        <DialogContent>
          <Stack spacing={2} sx={{ mt: 1, minWidth: 300 }}>
            <TextField
              label="A → B (forward)"
              value={newLinkAToB}
              onChange={(e) => setNewLinkAToB(e.target.value)}
              fullWidth
              placeholder="e.g., inhibits"
            />
            <TextField
              label="B → A (reverse)"
              value={newLinkBToA}
              onChange={(e) => setNewLinkBToA(e.target.value)}
              fullWidth
              placeholder="e.g., is inhibited by (auto-inferred if blank)"
              helperText="Leave blank to auto-infer"
            />
            <TextField
              label="Category"
              value={newCategory}
              onChange={(e) => setNewCategory(e.target.value)}
              fullWidth
              placeholder="e.g., causal, structural"
            />
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowAddLinkTypeDialog(false)}>Cancel</Button>
          <Button variant="contained" onClick={handleAddLinkType}>
            Add
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
