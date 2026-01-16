import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
  LinearProgress,
  Chip,
  Grid,
  Paper,
  Divider,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions
} from '@mui/material';
import {
  CloudUpload,
  Book,
  Science,
  Info,
  Add,
  Delete
} from '@mui/icons-material';
import { apiClient, API_ENDPOINTS } from '../config/api'; // Use apiClient instead of apiHelpers

interface BookDetails {
  id?: string;
  _id?: string;
  title: string;
  isbn: string;
  authors: string[];
  publisher: string;
  year: string;
  type: string;
  targetExam: string;
}

interface ChapterDetails {
  number: string;
  title: string;
}

interface UploadResponse {
  status: string;
  message: string;
  batch_id?: string; // Make batch_id optional as it might not always be there
  file_info: {
    filename: string;
    size_mb: number;
  };
  book_source: BookDetails;
  processing: {
    job_id: string;
    status: string;
    estimated_time: string;
  };
  source_tracking: {
    citation_format: string;
  };
}

const emptyBook: BookDetails = {
  title: '',
  isbn: '',
  authors: [''],
  publisher: '',
  year: new Date().getFullYear().toString(),
  type: 'book',
  targetExam: 'FRC Path Part 1'
};

const ImmunologyUpload: React.FC = () => {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [books, setBooks] = useState<BookDetails[]>([]);
  const [selectedBook, setSelectedBook] = useState<BookDetails | null>(null);
  const [chapterDetails, setChapterDetails] = useState<ChapterDetails>({
    number: '1',
    title: 'Case 1: X-Linked Agammaglobulinemia'
  });
  const [uploading, setUploading] = useState(false);
  const [uploadResult, setUploadResult] = useState<UploadResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [showAddBookDialog, setShowAddBookDialog] = useState(false);
  const [newBook, setNewBook] = useState<BookDetails>(emptyBook);
  const [pollingJobId, setPollingJobId] = useState<string | null>(null);
  const [pollingStatus, setPollingStatus] = useState<string | null>(null);

  useEffect(() => {
    fetchBooks();
  }, []);

  // Polling effect for checking job status
  useEffect(() => {
    if (!pollingJobId) return;

    const interval = setInterval(async () => {
      try {
        // Use the correct API endpoint from the configuration
        const statusResult = await apiClient.get<any>(API_ENDPOINTS.extraction.status(pollingJobId));
        setPollingStatus(statusResult.status);

        if (statusResult.status === 'COMPLETE' || statusResult.status === 'FAILED') {
          clearInterval(interval);
          setPollingJobId(null);
          // Refresh the book list to show new content
          fetchBooks();
        }
      } catch (err) {
        console.error('Polling error:', err);
        clearInterval(interval);
        setPollingJobId(null);
      }
    }, 5000); // Poll every 5 seconds

    return () => clearInterval(interval);
  }, [pollingJobId]);

  const fetchBooks = async () => {
    try {
      // Use the apiClient.get method
      const data = await apiClient.get<any>('/api/sources?type=book');
      const fetchedBooks = data.sources || [];
      setBooks(fetchedBooks);
      if (fetchedBooks.length > 0) {
        setSelectedBook(fetchedBooks[0]);
      }
    } catch (err) {
      console.error('Error fetching books:', err);
      setError('Failed to fetch books. Is the API server running and accessible?');
    }
  };

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      if (file.type === 'application/pdf') {
        setSelectedFile(file);
        setError(null);
      } else {
        setError('Please select a PDF file');
        setSelectedFile(null);
      }
    }
  };

  const handleBookChange = (bookId: string) => {
    const book = books.find(b => (b.id || b._id) === bookId);
    setSelectedBook(book || null);
  };

  const handleAuthorChange = (index: number, value: string) => {
    const newAuthors = [...newBook.authors];
    newAuthors[index] = value;
    setNewBook({ ...newBook, authors: newAuthors });
  };

  const addAuthor = () => {
    setNewBook({
      ...newBook,
      authors: [...newBook.authors, '']
    });
  };

  const removeAuthor = (index: number) => {
    if (newBook.authors.length > 1) {
      const newAuthors = newBook.authors.filter((_, i) => i !== index);
      setNewBook({ ...newBook, authors: newAuthors });
    }
  };

  const handleAddNewBook = async () => {
    if (!newBook.title || !newBook.isbn) {
      alert('Please fill in at least Title and ISBN for the new book.');
      return;
    }
    try {
      const payload = {
        type: 'book',
        isbn: newBook.isbn,
        title: newBook.title,
        authors: (newBook.authors || []).map(a => a.trim()).filter(Boolean),
        publisher: newBook.publisher,
        year: (newBook.year || '').toString().trim() || new Date().getFullYear().toString(),
      };
      
      // Use apiClient.post
      const result = await apiClient.post<any>('/api/sources', payload);
      
      if (result.source) {
        setShowAddBookDialog(false);
        setNewBook(emptyBook);
        await fetchBooks();
        setSelectedBook(result.source);
      } else {
        alert(`Failed to add source: ${result.error || 'unknown error'}`);
      }
    } catch (err: any) {
      alert(`An error occurred while adding the source: ${err.message}`);
    }
  };

  const handleUpload = async () => {
    if (!selectedFile) {
      setError('Please select a file');
      return;
    }

    if (!selectedBook) {
      setError('Please select a book');
      return;
    }

    setUploading(true);
    setError(null);

    const formData = new FormData();
    formData.append('chapter', selectedFile, selectedFile.name);
    formData.append('chapter_number', chapterDetails.number);
    formData.append('chapter_title', chapterDetails.title);

    const sourceId = selectedBook?.id || selectedBook?._id || '';
    if (sourceId) {
      formData.append('source_id', sourceId);
    }
    
    formData.append('book_title', selectedBook.title);
    formData.append('book_isbn', selectedBook.isbn);
    formData.append('book_authors', (selectedBook.authors || []).filter(Boolean).join(','));
    formData.append('book_publisher', selectedBook.publisher || '');
    formData.append('book_year', (selectedBook.year || '').toString());
    formData.append('extraction_method', 'textract_immunology');
    formData.append('content_type', 'medical_chapter');
    formData.append('target_exam', selectedBook.targetExam || '');
    formData.append('domain', 'immunology');

    // Debug info
    console.log('Uploading to endpoint:', API_ENDPOINTS.studyAreas.immunology.upload); // FIXED: use correct path
    const formDataKeys: string[] = [];
    formData.forEach((value, key) => {
      formDataKeys.push(key);
    });
    console.log('Form data keys:', formDataKeys);

    try {
      // Use apiClient.upload for file uploads
      const result = await apiClient.upload<UploadResponse>(
        API_ENDPOINTS.studyAreas.immunology.upload,
        formData,
        (progressEvent) => {
          const progress = Math.round(
            (progressEvent.loaded * 100) / (progressEvent.total || 1)
          );
          console.log(`Upload progress: ${progress}%`);
        }
      );
      
      if (result) {
        setUploadResult(result);
        console.log('Upload successful:', result);
        // Start polling if we get a job ID
        const jobId = result.processing?.job_id || result.batch_id;
        if (jobId) {
          setPollingJobId(jobId);
          setPollingStatus('PENDING');
        }

        // Add quick action: open semantic extractor with this upload as default source
        // Use navigate with query params so extractor can prefill source
        const encodedSource = encodeURIComponent(JSON.stringify({
          source_id: result.book_source?.id || result.book_source?._id || '',
          batch_id: result.batch_id || result.processing?.job_id || '',
          title: result.book_source?.title || ''
        }));
        // Expose a convenience link via window.history state for manual opening, and also render button in UI below
        window.localStorage.setItem('last_upload_source', encodedSource);
      }
    } catch (err: any) {
      console.error('Upload error:', err);
      setError(err.message || 'Network error occurred');
    } finally {
      setUploading(false);
    }
  };

  return (
    <Box sx={{ maxWidth: 800, mx: 'auto', p: 3 }}>
      <Typography variant="h4" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <Science color="primary" />
        Immunology Chapter Upload
      </Typography>

      {/* Debug Info */}
      <Alert severity="info" sx={{ mb: 3 }}>
        <Typography variant="body2">
          <strong>Environment:</strong> {process.env.NODE_ENV || 'development'}
        </Typography>
        <Typography variant="body2">
          <strong>Upload Domain:</strong> immunology
        </Typography>
      </Alert>

      {/* File Upload Section */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <CloudUpload />
            Select Chapter PDF
          </Typography>
          
          <Button
            variant="outlined"
            component="label"
            startIcon={<CloudUpload />}
            sx={{ mb: 2 }}
          >
            Choose PDF File
            <input
              type="file"
              hidden
              accept=".pdf"
              onChange={handleFileSelect}
            />
          </Button>

          {selectedFile && (
            <Alert severity="success" sx={{ mt: 2 }}>
              Selected: {selectedFile.name} ({(selectedFile.size / (1024 * 1024)).toFixed(2)} MB)
            </Alert>
          )}

          {uploadResult && (
            <Box sx={{ mt: 2, display: 'flex', gap: 1, alignItems: 'center' }}>
              <Button
                variant="contained"
                onClick={() => {
                  // Navigate to semantic extractor with source info stored in localStorage
                  const encoded = window.localStorage.getItem('last_upload_source') || '';
                  const url = `/esp-organizer/semantic-links/extract${encoded ? `?source=${encodeURIComponent(encoded)}` : ''}`;
                  window.location.href = url;
                }}
              >
                Open Semantic Extractor (use this source)
              </Button>
              <Button
                variant="outlined"
                onClick={() => {
                  // open the upload job status page or scroll to polling info
                  if (uploadResult.processing?.job_id) {
                    setPollingJobId(uploadResult.processing.job_id);
                    setPollingStatus('PENDING');
                  }
                }}
              >
                View Processing Status
              </Button>
            </Box>
          )}
        </CardContent>
      </Card>

      {/* Book Selection Section */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Book />
            Book Information
          </Typography>

          <FormControl fullWidth sx={{ mb: 2 }}>
            <InputLabel>Select Book</InputLabel>
            <Select
              value={selectedBook?.id || selectedBook?._id || ''}
              onChange={(e) => handleBookChange(e.target.value as string)}
              label="Select Book"
            >
              {books.map((book) => (
                <MenuItem key={book.id || book._id} value={book.id || book._id}>
                  {book.title} ({book.year})
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          
          <Button startIcon={<Add />} onClick={() => setShowAddBookDialog(true)}>
            Add New Book
          </Button>

          {selectedBook && (
            <Paper sx={{ p: 2, mt: 2, bgcolor: 'grey.50' }}>
              <Typography variant="subtitle2" gutterBottom>
                Selected Book Details:
              </Typography>
              <Typography variant="body2">
                <strong>Title:</strong> {selectedBook.title}
              </Typography>
              {selectedBook.isbn && (
                <Typography variant="body2">
                  <strong>ISBN:</strong> {selectedBook.isbn}
                </Typography>
              )}
              <Typography variant="body2">
                <strong>Authors:</strong> {selectedBook.authors.filter(a => a).join(', ')}
              </Typography>
              {selectedBook.publisher && (
                <Typography variant="body2">
                  <strong>Publisher:</strong> {selectedBook.publisher} ({selectedBook.year})
                </Typography>
              )}
              {selectedBook.targetExam && (
                <Box sx={{ mt: 1 }}>
                  <Chip label={selectedBook.targetExam} color="primary" size="small" />
                </Box>
              )}
            </Paper>
          )}
        </CardContent>
      </Card>

      {/* Chapter Details Section */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Chapter Details
          </Typography>
          
          <Grid container spacing={2}>
            <Grid item xs={4}>
              <TextField
                fullWidth
                label="Chapter Number"
                value={chapterDetails.number}
                onChange={(e) => setChapterDetails({ ...chapterDetails, number: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={8}>
              <TextField
                fullWidth
                label="Chapter Title"
                value={chapterDetails.title}
                onChange={(e) => setChapterDetails({ ...chapterDetails, title: e.target.value })}
                required
              />
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Upload Button */}
      <Box sx={{ textAlign: 'center', mb: 3 }}>
        <Button
          variant="contained"
          size="large"
          onClick={handleUpload}
          disabled={!selectedFile || !selectedBook || uploading}
          startIcon={<CloudUpload />}
        >
          {uploading ? 'Processing...' : 'Upload & Process Chapter'}
        </Button>
      </Box>

      {/* Progress */}
      {uploading && (
        <Box sx={{ mb: 3 }}>
          <LinearProgress />
          <Typography variant="body2" sx={{ mt: 1, textAlign: 'center' }}>
            Uploading and processing chapter with AWS Textract...
          </Typography>
        </Box>
      )}

      {/* Polling Status */}
      {pollingJobId && (
        <Alert severity="info" sx={{ mb: 3 }}>
          Processing job <strong>{pollingJobId}</strong>. Current status: <strong>{pollingStatus || 'Initializing...'}</strong>
        </Alert>
      )}

      {/* Error Display */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Success Result */}
      {uploadResult && (
        <Card sx={{ mb: 3 }}>
          <CardContent>
            <Alert severity="success" sx={{ mb: 2 }}>
              {uploadResult.message}
            </Alert>
            
            <Typography variant="h6" gutterBottom>
              Processing Details
            </Typography>
            
            <Grid container spacing={2}>
              {uploadResult.processing && (
                <Grid item xs={6}>
                  <Typography variant="body2">
                    <strong>Job ID:</strong> {uploadResult.processing.job_id}
                  </Typography>
                  <Typography variant="body2">
                    <strong>Status:</strong> {uploadResult.processing.status}
                  </Typography>
                  <Typography variant="body2">
                    <strong>Estimated Time:</strong> {uploadResult.processing.estimated_time}
                  </Typography>
                </Grid>
              )}
              {uploadResult.file_info && (
                <Grid item xs={6}>
                  <Typography variant="body2">
                    <strong>File:</strong> {uploadResult.file_info.filename}
                  </Typography>
                  <Typography variant="body2">
                    <strong>Size:</strong> {uploadResult.file_info.size_mb.toFixed(2)} MB
                  </Typography>
                </Grid>
              )}
            </Grid>

            <Divider sx={{ my: 2 }} />
            
            {uploadResult.source_tracking && (
              <>
                <Typography variant="subtitle2" gutterBottom>
                  Source Citation Format:
                </Typography>
                <Paper sx={{ p: 1, bgcolor: 'grey.100' }}>
                  <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                    {uploadResult.source_tracking.citation_format}
                  </Typography>
                </Paper>
              </>
            )}
          </CardContent>
        </Card>
      )}

      {/* Add Book Dialog */}
      <Dialog open={showAddBookDialog} onClose={() => setShowAddBookDialog(false)} maxWidth="md" fullWidth>
        <DialogTitle>Add New Book Source</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ pt: 1 }}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Book Title"
                value={newBook.title}
                onChange={(e) => setNewBook({ ...newBook, title: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={6}>
              <TextField
                fullWidth
                label="ISBN"
                value={newBook.isbn}
                onChange={(e) => setNewBook({ ...newBook, isbn: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={6}>
              <TextField
                fullWidth
                label="Publisher"
                value={newBook.publisher}
                onChange={(e) => setNewBook({ ...newBook, publisher: e.target.value })}
              />
            </Grid>
            <Grid item xs={6}>
              <TextField
                fullWidth
                label="Publication Year"
                value={newBook.year}
                onChange={(e) => setNewBook({ ...newBook, year: e.target.value })}
              />
            </Grid>
            <Grid item xs={6}>
              <TextField
                fullWidth
                label="Target Exam"
                value={newBook.targetExam}
                onChange={(e) => setNewBook({ ...newBook, targetExam: e.target.value })}
                placeholder="e.g., FRC Path Part 1"
              />
            </Grid>
            <Grid item xs={12}>
              <Typography variant="subtitle2" gutterBottom>Authors</Typography>
              {newBook.authors.map((author, index) => (
                <Box key={index} sx={{ display: 'flex', gap: 1, mb: 1 }}>
                  <TextField
                    fullWidth
                    value={author}
                    onChange={(e) => handleAuthorChange(index, e.target.value)}
                    placeholder={`Author ${index + 1}`}
                  />
                  {newBook.authors.length > 1 && (
                    <IconButton onClick={() => removeAuthor(index)} color="error">
                      <Delete />
                    </IconButton>
                  )}
                </Box>
              ))}
              <Button startIcon={<Add />} onClick={addAuthor} size="small">
                Add Author
              </Button>
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowAddBookDialog(false)}>Cancel</Button>
          <Button onClick={handleAddNewBook} variant="contained">
            Save Book
          </Button>
        </DialogActions>
      </Dialog>

      {/* Help Information */}
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Info />
            Processing Information
          </Typography>
          
          <Typography variant="body2" paragraph>
            This system uses AWS Textract to extract content from medical textbooks with specialized processing for:
          </Typography>
          
          <Box component="ul" sx={{ pl: 2 }}>
            <li>Medical terminology and immunology concepts</li>
            <li>Clinical case studies and patient presentations</li>
            <li>Complex formatting (subscripts, superscripts)</li>
            <li>Laboratory results and diagnostic processes</li>
          </Box>
          
          <Typography variant="body2" sx={{ mt: 2 }}>
            All extracted content will include proper source attribution for academic integrity.
          </Typography>
        </CardContent>
      </Card>
    </Box>
  );
};

export default ImmunologyUpload;