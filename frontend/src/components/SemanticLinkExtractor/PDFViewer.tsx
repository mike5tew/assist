import React, { useRef, useEffect, useState, useCallback } from 'react';
import { Box, Typography, CircularProgress, Alert, Pagination, Chip, Stack } from '@mui/material';
// @ts-ignore - pdfjs-dist has type issues
import * as pdfjsLib from 'pdfjs-dist';

// Set up PDF.js worker - must match installed package version
const pdfjsVersion = '4.10.38';
if (pdfjsLib.GlobalWorkerOptions) {
  pdfjsLib.GlobalWorkerOptions.workerSrc = `https://unpkg.com/pdfjs-dist@${pdfjsVersion}/build/pdf.worker.min.mjs`;
}

interface PDFViewerProps {
  file: File;
  onTermSelected: (term: string, context: string, position?: { start: number; end: number }) => void;
  isSelecting: boolean;
}

/**
 * PDFViewer - Renders PDF pages with selectable text overlay
 * 
 * Features:
 * - Multi-page navigation with pagination
 * - Dual display: canvas (visual) + text layer (selectable)
 * - Text selection fires callback for semantic link workflow
 * - Visual feedback for selection state
 */
export default function PDFViewer({ file, onTermSelected, isSelecting }: PDFViewerProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedText, setSelectedText] = useState<string>('');
  const [pageText, setPageText] = useState<string>('');
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [pdfDoc, setPdfDoc] = useState<any>(null);
  const [renderKey, setRenderKey] = useState(0);

  // Handle text selection from the text overlay
  const handleTextSelection = useCallback(() => {
    const selection = window.getSelection();
    if (!selection || selection.toString().trim().length === 0) return;

    const text = selection.toString().trim();
    if (text.length > 0) {
      setSelectedText(text);
      
      // Get surrounding context (up to 200 chars around selection)
      const fullText = pageText;
      const selectionStart = fullText.indexOf(text);
      const contextStart = Math.max(0, selectionStart - 100);
      const contextEnd = Math.min(fullText.length, selectionStart + text.length + 100);
      const context = fullText.slice(contextStart, contextEnd);

      onTermSelected(text, context, 
        selectionStart >= 0 ? { start: selectionStart, end: selectionStart + text.length } : undefined
      );
      
      // Clear browser selection
      selection.removeAllRanges();
    }
  }, [pageText, onTermSelected]);

  // Load PDF document
  useEffect(() => {
    let cancelled = false;

    const loadPDF = async () => {
      if (!file) return;

      try {
        setIsLoading(true);
        setError(null);
        setPageText('');
        setPdfDoc(null);

        const fileReader = new FileReader();
        fileReader.onload = async (e) => {
          if (cancelled) return;
          try {
            const pdfData = e.target?.result as ArrayBuffer;
            const pdf = await pdfjsLib.getDocument({ data: pdfData }).promise;
            if (cancelled) return;
            setPdfDoc(pdf);
            setTotalPages(pdf.numPages);
            setCurrentPage(1);
            setRenderKey(k => k + 1);
          } catch (err) {
            if (!cancelled) {
              setError(`Failed to load PDF: ${err instanceof Error ? err.message : 'Unknown error'}`);
              setIsLoading(false);
            }
          }
        };
        fileReader.readAsArrayBuffer(file);
      } catch (err) {
        if (!cancelled) {
          setError(`Failed to read file: ${err instanceof Error ? err.message : 'Unknown error'}`);
          setIsLoading(false);
        }
      }
    };

    loadPDF();

    return () => {
      cancelled = true;
    };
  }, [file]);

  // Render current page
  useEffect(() => {
    let cancelled = false;

    const renderPage = async () => {
      if (!pdfDoc || !canvasRef.current) return;

      try {
        setIsLoading(true);
        const page = await pdfDoc.getPage(currentPage);
        if (cancelled) return;

        const scale = 1.5;
        const viewport = page.getViewport({ scale });

        const canvas = canvasRef.current;
        const context = canvas.getContext('2d');
        if (!context) throw new Error('Could not get canvas context');

        canvas.width = viewport.width;
        canvas.height = viewport.height;

        // Render PDF page to canvas
        await page.render({
          canvasContext: context,
          viewport: viewport,
        }).promise;

        if (cancelled) return;

        // Extract text content for selection
        const txtContent = await page.getTextContent();
        const extractedText = txtContent.items.map((it: any) => (it.str || '')).join(' ');
        setPageText(extractedText);
        setIsLoading(false);
      } catch (err) {
        if (!cancelled) {
          setError(`Failed to render page: ${err instanceof Error ? err.message : 'Unknown error'}`);
          setIsLoading(false);
        }
      }
    };

    renderPage();

    return () => {
      cancelled = true;
    };
  }, [pdfDoc, currentPage, renderKey]);

  const handlePageChange = (_event: React.ChangeEvent<unknown>, page: number) => {
    setCurrentPage(page);
    setSelectedText('');
  };

  return (
    <Box sx={{ width: '100%' }}>
      {/* Selection feedback */}
      {selectedText && (
        <Box sx={{ mb: 1, p: 1, backgroundColor: '#e8f5e9', borderRadius: 1 }}>
          <Typography variant="caption" color="success.main">
            ✓ Selected: <strong>"{selectedText}"</strong>
          </Typography>
        </Box>
      )}

      {/* Error display */}
      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

      {/* PDF Canvas Container */}
      <Box
        sx={{
          width: '100%',
          minHeight: 300,
          maxHeight: 400,
          backgroundColor: '#f5f5f5',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: isLoading && !pageText ? 'center' : 'flex-start',
          borderRadius: 1,
          overflow: 'auto',
          border: '1px solid #ddd',
          mb: 2,
        }}
      >
        {isLoading && !pdfDoc && (
          <Box sx={{ textAlign: 'center', py: 4 }}>
            <CircularProgress size={40} sx={{ mb: 2 }} />
            <Typography variant="body2" color="textSecondary">
              Loading PDF...
            </Typography>
          </Box>
        )}
        
        <canvas 
          ref={canvasRef} 
          style={{ 
            display: pdfDoc ? 'block' : 'none',
            maxWidth: '100%', 
            height: 'auto' 
          }} 
        />
      </Box>

      {/* Selectable Text Layer */}
      {pageText && (
        <Box
          sx={{
            p: 2,
            backgroundColor: '#fffef5',
            border: isSelecting ? '2px solid #1976d2' : '1px solid #ddd',
            borderRadius: 1,
            maxHeight: 250,
            overflow: 'auto',
            transition: 'border 0.2s',
          }}
        >
          <Typography 
            variant="caption" 
            sx={{ 
              display: 'block', 
              mb: 1, 
              color: isSelecting ? '#1565c0' : '#666',
              fontWeight: 500,
              backgroundColor: isSelecting ? '#e3f2fd' : '#f5f5f5',
              p: 1,
              borderRadius: 0.5,
            }}
          >
            {isSelecting 
              ? '👆 Highlight text below to select terms or statements' 
              : '📄 Text preview (enable selection above)'}
          </Typography>
          <Typography
            component="div"
            onMouseUp={isSelecting ? handleTextSelection : undefined}
            sx={{
              whiteSpace: 'pre-wrap',
              fontSize: '14px',
              lineHeight: 1.6,
              color: '#333',
              cursor: isSelecting ? 'text' : 'default',
              userSelect: isSelecting ? 'text' : 'none',
            }}
          >
            {pageText}
          </Typography>
        </Box>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <Stack direction="row" spacing={2} alignItems="center" justifyContent="center" sx={{ mt: 2 }}>
          <Pagination 
            count={totalPages} 
            page={currentPage} 
            onChange={handlePageChange}
            color="primary"
            size="small"
          />
          <Chip label={`Page ${currentPage} of ${totalPages}`} size="small" variant="outlined" />
        </Stack>
      )}
    </Box>
  );
}
