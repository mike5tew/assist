import React from 'react';
import GenericStudyArea, { StudyAreaConfig } from '../components/generic/GenericStudyArea';
import ImmunologyUpload from '../components/ImmunologyUpload';
import { Science } from '@mui/icons-material';

// Add additional examples to help users formulate better queries
const exampleQueries = [
  "What is X-Linked Agammaglobulinemia (XLA)?",
  "Explain the pathophysiology of SCID",
  "Compare and contrast B-cell and T-cell development",
  "What are the clinical features of Common Variable Immunodeficiency?",
  "How does the complement system work?"
];

// Get a random example query
const getRandomExample = () => {
  return exampleQueries[Math.floor(Math.random() * exampleQueries.length)];
};

const environment = process.env.REACT_APP_ENVIRONMENT || 'Local';

const immunologyConfig: StudyAreaConfig = {
  title: `Medical Immunology (${environment})`,
  domain: 'immunology',
  icon: <Science sx={{ mr: 1 }} />,
  color: '#6a1b9a', // A deep purple color
  searchPlaceholder: getRandomExample(),
  apiEndpoints: {
    search: '/api/immunology/search',
    chapters: '/api/immunology/chapters',
    caseStudies: '/api/immunology/case-studies',
    terms: '/api/immunology/medical-terms',
    upload: '/api/immunology/upload-chapter',
    extraction: {
      status: (jobId) => `/api/extraction/status/${jobId}`
    }
  },
  uploadComponent: <ImmunologyUpload />,
  // Add help text for users
  helpText: "Search for immunology topics, cases, or specific immune disorders. Try queries like 'What is XLA?' or 'Explain B-cell development'.",
  quickLinks: [
    {
      label: 'Open ACLI V(D)J Module',
      path: '/acli/vdj',
      description: 'FRCPath-style module with speedread, audio, flashcards, and exam mode.',
    },
  ],
};

const ImmunologyPage: React.FC = () => {
  return <GenericStudyArea config={immunologyConfig} />;
};

export default ImmunologyPage;
