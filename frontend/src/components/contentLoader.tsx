import React, { useState } from 'react';
import { 
  Button, Input, Box, Select, MenuItem, Snackbar, Alert, Tabs, Tab, TextField,
  Typography, Chip, FormControl, InputLabel, Checkbox, ListItemText 
} from '@mui/material';
import Grid from '@mui/material/Grid';
import PDFViewer from './SemanticLinkExtractor/PDFViewer';
import { useTranslation } from 'react-i18next';

type UploadType = 'keywords' | 'semantic' | 'rawText' | 'documentCreator';
type FileType = 'json' | 'csv' | 'txt' | 'pdf' | 'html' | 'markdown';
type SemanticCategory = 'technical' | 'legal' | 'medical' | 'financial' | 'educational';

export const ContentLoader: React.FC = () => {
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState<UploadType>('keywords');
  const [toast, setToast] = useState({
    open: false,
    message: '',
    severity: 'success' as 'success' | 'error'
  });

  // Common state
  const [loading, setLoading] = useState(false);
  const [file, setFile] = useState<File | null>(null);
  const [fileType, setFileType] = useState<FileType | null>(null);

  // Keyword upload state
  const [keyword, setKeyword] = useState('');
  const [definition, setDefinition] = useState('');

  // Semantic connections state - Updated
  const [sourceEntity, setSourceEntity] = useState('');
  const [targetEntity, setTargetEntity] = useState('');
  const [relationshipType, setRelationshipType] = useState('');
  const [contextualEvidence, setContextualEvidence] = useState('');
  const [sourceDocumentId, setSourceDocumentId] = useState('');
  const [confidenceScore, setConfidenceScore] = useState<number | ''>('');


  // Raw text upload state
  const [textContent, setTextContent] = useState('');
  const [metadata, setMetadata] = useState<Record<string, string>>({});
  const [url, setUrl] = useState('');

  // Document creator state
  const [documentTitle, setDocumentTitle] = useState('');
  const [documentContent, setDocumentContent] = useState('');
  const [selectedCategories, setSelectedCategories] = useState<SemanticCategory[]>([]);
  const [customCategories, setCustomCategories] = useState<string[]>([]);
  const [newCategory, setNewCategory] = useState('');

  const categories: SemanticCategory[] = ['technical', 'legal', 'medical', 'financial', 'educational'];

  const handleCloseToast = () => {
    setToast({ ...toast, open: false });
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: UploadType) => {
    setActiveTab(newValue);
  };

  // Metadata handlers
  const addMetadataField = () => setMetadata({ ...metadata, '': '' });

  const updateMetadata = (oldKey: string, newKey: string, newValue: string) => {
    const newMetadata = { ...metadata };
    delete newMetadata[oldKey];
    if (newKey) newMetadata[newKey] = newValue;
    setMetadata(newMetadata);
  };

  const removeMetadataField = (key: string) => {
    const newMetadata = { ...metadata };
    delete newMetadata[key];
    setMetadata(newMetadata);
  };

  // Category handlers
  const handleCategoryChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const value = event.target.value as SemanticCategory;
    setSelectedCategories(
      event.target.checked
        ? [...selectedCategories, value]
        : selectedCategories.filter(cat => cat !== value)
    );
  };

  const addCustomCategory = () => {
    if (newCategory && !customCategories.includes(newCategory)) {
      setCustomCategories([...customCategories, newCategory]);
      setNewCategory('');
    }
  };

  const removeCustomCategory = (category: string) => {
    setCustomCategories(customCategories.filter(c => c !== category));
  };

  // Upload handlers
  const uploadKeywords = async () => {
    try {
      setLoading(true);
      // API call would go here
      setToast({
        open: true,
        message: 'Keywords uploaded successfully',
        severity: 'success'
      });
    } catch (error) {
      setToast({
        open: true,
        message: 'Failed to upload keywords',
        severity: 'error'
      });
    } finally {
      setLoading(false);
    }
  };

  const uploadSemanticConnections = async () => {
    try {
      setLoading(true);
      // API call would go here, using:
      // sourceEntity, targetEntity, relationshipType, contextualEvidence, sourceDocumentId, confidenceScore
      setToast({
        open: true,
        message: 'Semantic connections uploaded successfully',
        severity: 'success'
      });
    } catch (error) {
      setToast({
        open: true,
        message: 'Failed to upload semantic connections',
        severity: 'error'
      });
    } finally {
      setLoading(false);
    }
  };

  const uploadRawContent = async () => {
    try {
      setLoading(true);
      // API call would go here
      setToast({
        open: true,
        message: 'Content uploaded successfully',
        severity: 'success'
      });
    } catch (error) {
      setToast({
        open: true,
        message: 'Failed to upload content',
        severity: 'error'
      });
    } finally {
      setLoading(false);
    }
  };

  const createDocument = async () => {
    try {
      setLoading(true);
      const allCategories = [...selectedCategories, ...customCategories];
      // API call would go here with documentContent and allCategories
      setToast({
        open: true,
        message: 'Document created successfully',
        severity: 'success'
      });
    } catch (error) {
      setToast({
        open: true,
        message: 'Failed to create document',
        severity: 'error'
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box sx={{ padding: 2 }}>
      <Tabs value={activeTab} onChange={handleTabChange} sx={{ mb: 3 }}>
        <Tab label="Keywords" value="keywords" />
        <Tab label="Semantic Links" value="semantic" />
        <Tab label="Upload Content" value="rawText" />
        <Tab label="Create Document" value="documentCreator" />
      </Tabs>

      {activeTab === 'keywords' && (
        <Grid container spacing={2}>
          <Grid item xs={12} md={6}>
            <TextField
              label="Keyword"
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              fullWidth
              sx={{ mb: 2 }}
            />
            <TextField
              label="Definition"
              value={definition}
              onChange={(e) => setDefinition(e.target.value)}
              multiline
              rows={4}
              fullWidth
              sx={{ mb: 2 }}
            />
            <Button
              variant="contained"
              onClick={uploadKeywords}
              disabled={loading || !keyword || !definition}
              fullWidth
            >
              {loading ? 'Uploading...' : 'Upload Keywords'}
            </Button>
          </Grid>
        </Grid>
      )}

      {activeTab === 'semantic' && (
        <Grid container spacing={2}>
          <Grid item xs={12} md={6}>
            <TextField
              label="Source Entity"
              value={sourceEntity}
              onChange={(e) => setSourceEntity(e.target.value)}
              fullWidth
              sx={{ mb: 2 }}
            />
            <TextField
              label="Target Entity"
              value={targetEntity}
              onChange={(e) => setTargetEntity(e.target.value)}
              fullWidth
              sx={{ mb: 2 }}
            />
            {/* Consider making Relationship Type a free text or a more comprehensive ontology-based select */}
            <TextField
              label="Relationship Type"
              value={relationshipType}
              onChange={(e) => setRelationshipType(e.target.value)}
              fullWidth
              sx={{ mb: 2 }}
              helperText="e.g., PRODUCES, ENABLES, CAUSES, IS_A"
            />
            <TextField
              label="Contextual Evidence / Original Phrasing"
              value={contextualEvidence}
              onChange={(e) => setContextualEvidence(e.target.value)}
              multiline
              rows={4}
              fullWidth
              sx={{ mb: 2 }}
            />
            <TextField
              label="Source Document ID"
              value={sourceDocumentId}
              onChange={(e) => setSourceDocumentId(e.target.value)}
              fullWidth
              sx={{ mb: 2 }}
            />
            <TextField
              label="Confidence Score (Optional)"
              type="number"
              value={confidenceScore}
              onChange={(e) => setConfidenceScore(e.target.value === '' ? '' : parseFloat(e.target.value))}
              fullWidth
              sx={{ mb: 2 }}
              inputProps={{ step: "0.01" }}
            />
            <Button
              variant="contained"
              onClick={uploadSemanticConnections}
              disabled={
                loading ||
                !sourceEntity ||
                !targetEntity ||
                !relationshipType ||
                !contextualEvidence ||
                !sourceDocumentId
              }
              fullWidth
            >
              {loading ? 'Uploading...' : 'Create Semantic Link'}
            </Button>
          </Grid>
        </Grid>
      )}

      {activeTab === 'rawText' && (
        <Grid container spacing={2}>
          <Grid item xs={12} md={6}>
            <FormControl fullWidth sx={{ mb: 2 }}>
              <InputLabel>Document Type</InputLabel>
              <Select
                value={fileType || ''}
                onChange={(e) => setFileType(e.target.value as FileType)}
                label="Document Type"
              >
                <MenuItem value="json">JSON</MenuItem>
                <MenuItem value="csv">CSV</MenuItem>
                <MenuItem value="txt">Text</MenuItem>
                <MenuItem value="pdf">PDF</MenuItem>
                <MenuItem value="html">HTML</MenuItem>
                <MenuItem value="markdown">Markdown</MenuItem>
              </Select>
            </FormControl>
            
            <input
              type="file"
              onChange={(e) => setFile(e.target.files?.[0] || null)}
              style={{ marginBottom: 16 }}
              accept=".json,.csv,.txt,.pdf,.html,.md"
            />

            {/* PDF Preview & Selection */}
            {fileType === 'pdf' && file && (
              <Box sx={{ mt: 2 }}>
                <Typography variant="subtitle1" sx={{ mb: 1 }}>PDF Preview (select text below)</Typography>
                <PDFViewer
                  file={file}
                  onTermSelected={(text) => {
                    // If no source selected yet, set it; otherwise set target
                    if (!sourceEntity) {
                      setSourceEntity(text);
                      setToast({ open: true, message: `Source term selected: ${text}`, severity: 'success' });
                    } else if (!targetEntity) {
                      setTargetEntity(text);
                      setToast({ open: true, message: `Target term selected: ${text}`, severity: 'success' });
                    } else {
                      // both already set; default to replacing target
                      setTargetEntity(text);
                      setToast({ open: true, message: `Target term replaced: ${text}`, severity: 'success' });
                    }
                  }}
                  isSelecting={true}
                />

                <Box sx={{ display: 'flex', gap: 1, mt: 1 }}>
                  <Chip label={`Source: ${sourceEntity || 'none'}`} onDelete={() => setSourceEntity('')} />
                  <Chip label={`Target: ${targetEntity || 'none'}`} onDelete={() => setTargetEntity('')} />
                </Box>

                <Box sx={{ display: 'flex', gap: 1, mt: 1 }}>
                  <Button
                    variant="outlined"
                    onClick={() => {
                      if (!sourceEntity || !targetEntity) {
                        setToast({ open: true, message: 'Select both source and target terms first', severity: 'error' });
                        return;
                      }
                      // Move to semantic tab and prefill fields
                      setActiveTab('semantic');
                    }}
                  >
                    Use selection to create semantic link
                  </Button>
                  <Button variant="text" onClick={() => { setSourceEntity(''); setTargetEntity(''); }}>
                    Clear selections
                  </Button>
                </Box>
              </Box>
            )}

            <TextField
              label="Or paste text content"
              value={textContent}
              onChange={(e) => setTextContent(e.target.value)}
              multiline
              rows={6}
              fullWidth
              sx={{ mb: 2 }}
            />

            <TextField
              label="Or enter URL"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              fullWidth
              sx={{ mb: 2 }}
            />

            <Box sx={{ mb: 2 }}>
              <Typography variant="h6" gutterBottom>Metadata</Typography>
              {Object.entries(metadata).map(([key, value]) => (
                <Box key={key} sx={{ display: 'flex', gap: 1, mb: 1 }}>
                  <TextField
                    placeholder="Metadata key"
                    value={key}
                    onChange={(e) => updateMetadata(key, e.target.value, value)}
                    sx={{ flex: 1 }}
                  />
                  <TextField
                    placeholder="Metadata value"
                    value={value}
                    onChange={(e) => updateMetadata(key, key, e.target.value)}
                    sx={{ flex: 1 }}
                  />
                  <Button
                    variant="outlined"
                    onClick={() => removeMetadataField(key)}
                  >
                    Remove
                  </Button>
                </Box>
              ))}
              <Button
                variant="outlined"
                onClick={addMetadataField}
                sx={{ mt: 1 }}
              >
                Add Metadata Field
              </Button>
            </Box>

            <Button
              variant="contained"
              onClick={uploadRawContent}
              disabled={loading || (!file && !textContent && !url)}
              fullWidth
            >
              {loading ? 'Uploading...' : 'Upload Content'}
            </Button>
          </Grid>
        </Grid>
      )}

      {activeTab === 'documentCreator' && (
        <Grid container spacing={2}>
          <Grid item xs={12} md={6}>
            <TextField
              label="Document Title"
              value={documentTitle}
              onChange={(e) => setDocumentTitle(e.target.value)}
              fullWidth
              sx={{ mb: 2 }}
            />
            
            <TextField
              label="Document Content"
              value={documentContent}
              onChange={(e) => setDocumentContent(e.target.value)}
              multiline
              rows={10}
              fullWidth
              sx={{ mb: 2 }}
            />

            <Typography variant="h6" gutterBottom>Semantic Categories</Typography>
            
            <FormControl fullWidth sx={{ mb: 2 }}>
              <Typography variant="subtitle1">Predefined Categories</Typography>
              {categories.map((category) => (
                <Box key={category} sx={{ display: 'flex', alignItems: 'center' }}>
                  <Checkbox
                    checked={selectedCategories.includes(category)}
                    onChange={handleCategoryChange}
                    value={category}
                  />
                  <Typography>{category.charAt(0).toUpperCase() + category.slice(1)}</Typography>
                </Box>
              ))}
            </FormControl>

            <Box sx={{ mb: 2 }}>
              <Typography variant="subtitle1">Custom Categories</Typography>
              <Box sx={{ display: 'flex', gap: 1, mb: 1 }}>
                <TextField
                  value={newCategory}
                  onChange={(e) => setNewCategory(e.target.value)}
                  placeholder="Add custom category"
                  sx={{ flex: 1 }}
                />
                <Button 
                  variant="outlined" 
                  onClick={addCustomCategory}
                  disabled={!newCategory}
                >
                  Add
                </Button>
              </Box>
              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                {customCategories.map((category) => (
                  <Chip
                    key={category}
                    label={category}
                    onDelete={() => removeCustomCategory(category)}
                  />
                ))}
              </Box>
            </Box>

            <Button
              variant="contained"
              onClick={createDocument}
              disabled={loading || !documentContent}
              fullWidth
            >
              {loading ? 'Creating...' : 'Create Document'}
            </Button>
          </Grid>
        </Grid>
      )}

      <Snackbar
        open={toast.open}
        autoHideDuration={3000}
        onClose={handleCloseToast}
      >
        <Alert onClose={handleCloseToast} severity={toast.severity}>
          {toast.message}
        </Alert>
      </Snackbar>
    </Box>
  );
};