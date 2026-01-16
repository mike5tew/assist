import React, { useEffect, useState } from 'react';
import {
  Box,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Slider,
  Stack,
  TextField,
  Typography,
  SelectChangeEvent,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Chip,
  IconButton,
} from '@mui/material';
import { 
  RelationshipType, 
  QualityFlag, 
  RELATIONSHIP_PAIRS 
} from '../../types/semanticLinks';

interface RelationshipSelectorProps {
  forwardRelation: string;
  inverseRelation: string;
  qualityFlag: QualityFlag | null;
  confidence: number;
  onForwardRelationChange: (forward: string, inverse: string) => void;
  onQualityFlagChange: (value: QualityFlag | null) => void;
  onConfidenceChange: (value: number) => void;
}

// Descriptions for common relationship types
const RELATIONSHIP_DESCRIPTIONS: Record<string, string> = {
  causes: 'A causes or produces B',
  causes_reduction_in: 'A causes a decrease in B',
  causes_increase_in: 'A causes an increase in B',
  treats: 'A is used to treat B',
  prevents: 'A prevents or stops B',
  is_a: 'A is a type or instance of B',
  part_of: 'A is a part or component of B',
  requires: 'A requires or needs B',
  results_in: 'A results in or leads to B',
  contrasts_with: 'A is opposite or contrasts with B',
  similar_to: 'A is similar to B',
  prerequisite_for: 'A must be learned before B',
  enables: 'A enables or allows B',
  associated_with: 'A is associated with B',
};

/**
 * RelationshipSelector - UI for selecting bidirectional relationship types
 * 
 * When a forward relation is selected, the inverse is auto-populated from RELATIONSHIP_PAIRS.
 * For custom relations, users can specify both directions.
 */
export default function RelationshipSelector({
  forwardRelation,
  inverseRelation,
  qualityFlag,
  confidence,
  onForwardRelationChange,
  onQualityFlagChange,
  onConfidenceChange,
}: RelationshipSelectorProps) {
  const QUALITY_FLAGS: QualityFlag[] = ['high_confidence', 'needs_verification', 'problematic'];

  const [customRelations, setCustomRelations] = useState<{ key: string; inverse: string; desc: string }[]>([]);
  const [showAddDialog, setShowAddDialog] = useState(false);
  const [newRelationKey, setNewRelationKey] = useState('');
  const [newRelationInverse, setNewRelationInverse] = useState('');
  const [newRelationDesc, setNewRelationDesc] = useState('');

  useEffect(() => {
    try {
      const stored = localStorage.getItem('custom_relations_v2');
      if (stored) setCustomRelations(JSON.parse(stored));
    } catch (err) {
      console.error('Failed to load custom relations', err);
    }
  }, []);

  useEffect(() => {
    localStorage.setItem('custom_relations_v2', JSON.stringify(customRelations));
  }, [customRelations]);

  // Get inverse for a forward relation
  const getInverse = (forward: string): string => {
    // Check built-in pairs
    if (RELATIONSHIP_PAIRS[forward]) {
      return RELATIONSHIP_PAIRS[forward];
    }
    // Check custom relations
    const custom = customRelations.find(c => c.key === forward);
    if (custom) return custom.inverse;
    // Default: same relation (symmetric)
    return forward;
  };

  const handleForwardChange = (forward: string) => {
    const inverse = getInverse(forward);
    onForwardRelationChange(forward, inverse);
  };

  const addCustomRelation = () => {
    const key = newRelationKey.trim().replace(/\s+/g, '_');
    const inverse = newRelationInverse.trim().replace(/\s+/g, '_') || key;
    if (!key) return;
    setCustomRelations([...customRelations, { key, inverse, desc: newRelationDesc }]);
    setNewRelationKey('');
    setNewRelationInverse('');
    setNewRelationDesc('');
    setShowAddDialog(false);
    onForwardRelationChange(key, inverse);
  };

  const handleQualityChange = (event: SelectChangeEvent<string>) => {
    const value = event.target.value;
    onQualityFlagChange((value as QualityFlag) || null);
  };

  // Combine built-in and custom relations
  const allRelations = [
    ...Object.keys(RELATIONSHIP_PAIRS),
    ...customRelations.map(c => c.key)
  ];

  return (
    <Stack spacing={2.5}>
      {/* Quick-select chips */}
      <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
        {allRelations.slice(0, 8).map((key) => (
          <Chip
            key={key}
            label={key.replace(/_/g, ' ')}
            clickable
            color={key === forwardRelation ? 'primary' : 'default'}
            onClick={() => handleForwardChange(key)}
          />
        ))}
        <Chip label="Add..." variant="outlined" onClick={() => setShowAddDialog(true)} />
      </Box>

      {/* Forward Relation Selector */}
      <FormControl fullWidth>
        <InputLabel>Forward Relation (A → B)</InputLabel>
        <Select
          value={forwardRelation}
          onChange={(e) => handleForwardChange(e.target.value)}
          label="Forward Relation (A → B)"
        >
          {Object.entries(RELATIONSHIP_DESCRIPTIONS).map(([key, desc]) => (
            <MenuItem key={key} value={key}>
              {key.replace(/_/g, ' ')} — {desc}
            </MenuItem>
          ))}
          {customRelations.map((c) => (
            <MenuItem key={c.key} value={c.key}>
              {c.key.replace(/_/g, ' ')} — {c.desc}
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      {/* Bidirectional relation display */}
      <Box
        sx={{
          p: 1.5,
          backgroundColor: '#e3f2fd',
          borderLeft: '4px solid #1976d2',
          borderRadius: 1,
        }}
      >
        <Typography variant="caption" color="primary" sx={{ display: 'block', mb: 0.5 }}>
          <strong>Forward:</strong> Source <em>{forwardRelation.replace(/_/g, ' ')}</em> Target
        </Typography>
        <Typography variant="caption" color="secondary">
          <strong>Inverse:</strong> Target <em>{inverseRelation.replace(/_/g, ' ')}</em> Source
        </Typography>
      </Box>

      {/* Quality Flag */}
      <FormControl fullWidth>
        <InputLabel>Quality Assessment</InputLabel>
        <Select
          value={qualityFlag || ''}
          onChange={handleQualityChange}
          label="Quality Assessment"
        >
          <MenuItem value="">
            <em>None (default)</em>
          </MenuItem>
          {QUALITY_FLAGS.map((flag) => (
            <MenuItem key={flag} value={flag}>
              {flag.replace(/_/g, ' ').charAt(0).toUpperCase() + flag.replace(/_/g, ' ').slice(1)}
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      {/* Confidence Slider */}
      <Box>
        <Typography variant="body2" gutterBottom>
          Confidence: <strong>{(confidence * 100).toFixed(0)}%</strong>
        </Typography>
        <Slider
          min={0}
          max={1}
          step={0.05}
          value={confidence}
          onChange={(_, value) => onConfidenceChange(value as number)}
          marks={[
            { value: 0, label: '0%' },
            { value: 0.5, label: '50%' },
            { value: 1, label: '100%' },
          ]}
          valueLabelDisplay="auto"
          valueLabelFormat={(value) => `${(value * 100).toFixed(0)}%`}
        />
      </Box>

      {/* Quality Flag Description */}
      {qualityFlag && (
        <Box
          sx={{
            p: 1.5,
            backgroundColor:
              qualityFlag === 'high_confidence'
                ? '#e8f5e9'
                : qualityFlag === 'needs_verification'
                  ? '#fff3e0'
                  : '#ffebee',
            border: `1px solid ${
              qualityFlag === 'high_confidence'
                ? '#4caf50'
                : qualityFlag === 'needs_verification'
                  ? '#ff9800'
                  : '#f44336'
            }`,
            borderRadius: 1,
          }}
        >
          <Typography variant="caption">
            <strong>Flag: {qualityFlag.replace(/_/g, ' ')}</strong>
            {qualityFlag === 'high_confidence' && ' — Curator confident in this link'}
            {qualityFlag === 'needs_verification' && ' — Link may need manual review'}
            {qualityFlag === 'problematic' && ' — Link has issues or inconsistencies'}
          </Typography>
        </Box>
      )}

      {/* Add relation dialog */}
      <Dialog open={showAddDialog} onClose={() => setShowAddDialog(false)} fullWidth maxWidth="sm">
        <DialogTitle>Add Bidirectional Relationship Type</DialogTitle>
        <DialogContent>
          <TextField
            label="Forward relation (e.g., 'inhibits')"
            value={newRelationKey}
            onChange={(e) => setNewRelationKey(e.target.value)}
            fullWidth
            sx={{ mb: 2, mt: 1 }}
            helperText="The relation from source to target"
          />
          <TextField
            label="Inverse relation (e.g., 'is_inhibited_by')"
            value={newRelationInverse}
            onChange={(e) => setNewRelationInverse(e.target.value)}
            fullWidth
            sx={{ mb: 2 }}
            helperText="The relation from target back to source (leave empty for symmetric)"
          />
          <TextField
            label="Description"
            value={newRelationDesc}
            onChange={(e) => setNewRelationDesc(e.target.value)}
            fullWidth
            multiline
            rows={2}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowAddDialog(false)}>Cancel</Button>
          <Button onClick={addCustomRelation} variant="contained">Add</Button>
        </DialogActions>
      </Dialog>
    </Stack>
  );
}
