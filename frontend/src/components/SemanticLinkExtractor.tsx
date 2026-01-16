import React, { useState } from 'react';
import {
  Box,
  Container,
  Typography,
  Tabs,
  Tab,
  Paper,
} from '@mui/material';
import {
  Book,
  Link as LinkIcon,
  List as ListIcon,
} from '@mui/icons-material';
import SourceRegister, { Source } from './SemanticLinkExtractor/SourceRegister';
import LinkBuilder from './SemanticLinkExtractor/LinkBuilder';
import ExtractionResults from './SemanticLinkExtractor/ExtractionResults';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel({ children, value, index, ...other }: TabPanelProps) {
  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`semantic-tabpanel-${index}`}
      aria-labelledby={`semantic-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ pt: 3 }}>{children}</Box>}
    </div>
  );
}

/**
 * SemanticLinkExtractor - Main component for the semantic link curation tool
 * 
 * Redesigned workflow:
 * 1. Source Register - Manage source documents (books, PDFs)
 * 2. Link Builder - Create semantic links from PDFs
 * 3. Results - View and manage extracted links
 */
export default function SemanticLinkExtractor() {
  const [currentTab, setCurrentTab] = useState(0);
  const [selectedSource, setSelectedSource] = useState<Source | null>(null);

  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    setCurrentTab(newValue);
  };

  const handleSourceSelect = (source: Source) => {
    setSelectedSource(source);
    // Optionally auto-switch to Link Builder
    setCurrentTab(1);
  };

  return (
    <Container maxWidth="xl" sx={{ py: 3 }}>
      {/* Header */}
      <Box sx={{ mb: 3 }}>
        <Typography variant="h4" gutterBottom sx={{ fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: 1 }}>
          🔗 Semantic Link Extractor
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Create training data for semantic extraction models. Register sources, upload PDFs, and create validated relationship links.
        </Typography>
      </Box>

      {/* Tabs */}
      <Paper sx={{ mb: 0 }}>
        <Tabs
          value={currentTab}
          onChange={handleTabChange}
          indicatorColor="primary"
          textColor="primary"
          sx={{ borderBottom: 1, borderColor: 'divider' }}
        >
          <Tab 
            icon={<Book />} 
            iconPosition="start" 
            label="Source Register" 
          />
          <Tab 
            icon={<LinkIcon />} 
            iconPosition="start" 
            label="Link Builder" 
          />
          <Tab 
            icon={<ListIcon />} 
            iconPosition="start" 
            label="Results" 
          />
        </Tabs>
      </Paper>

      {/* Tab Panels */}
      <TabPanel value={currentTab} index={0}>
        <SourceRegister 
          onSourceSelect={handleSourceSelect}
          selectedSourceId={selectedSource?.id || selectedSource?._id}
        />
      </TabPanel>

      <TabPanel value={currentTab} index={1}>
        <LinkBuilder 
          source={selectedSource}
          onBack={() => setCurrentTab(0)}
        />
      </TabPanel>

      <TabPanel value={currentTab} index={2}>
        <ExtractionResults 
          results={null}
          onReset={() => {}}
        />
      </TabPanel>
    </Container>
  );
}
