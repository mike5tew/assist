import React from 'react';
import { Box, Chip, Typography, TextField, IconButton } from '@mui/material';
import { Close as CloseIcon } from '@mui/icons-material';

interface TermSelectionPanelProps {
  term: string;
  context: string;
  onClear: () => void;
  step: 'source' | 'target';
  disabled?: boolean;
}

/**
 * TermSelectionPanel - Displays selected term and its context
 * 
 * Shows the term selected in the PDF and allows users to clear it.
 */
export default function TermSelectionPanel({
  term,
  context,
  onClear,
  step,
  disabled = false,
}: TermSelectionPanelProps) {
  if (!term) {
    return (
      <Box
        sx={{
          p: 2,
          backgroundColor: disabled ? '#f5f5f5' : '#fafafa',
          border: '1px dashed #ddd',
          borderRadius: 1,
          textAlign: 'center',
          color: 'text.secondary',
          opacity: disabled ? 0.5 : 1,
        }}
      >
        <Typography variant="body2">
          {disabled
            ? 'Select source term first'
            : step === 'source'
              ? 'Click a term in the document'
              : 'Click a second term in the document'}
        </Typography>
      </Box>
    );
  }

  return (
    <Box>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1.5 }}>
        <Chip
          label={term}
          color="primary"
          variant="outlined"
          onDelete={onClear}
          sx={{ fontWeight: 'bold', fontSize: '0.95rem' }}
        />
        <IconButton size="small" onClick={onClear}>
          <CloseIcon fontSize="small" />
        </IconButton>
      </Box>

      {context && (
        <TextField
          fullWidth
          multiline
          rows={3}
          label="Context"
          value={context}
          disabled
          variant="filled"
          size="small"
          helperText="Context where term appears in document"
        />
      )}
    </Box>
  );
}
