import React from 'react';
import { Engineering } from '@mui/icons-material';
import GenericStudyArea, { StudyAreaConfig } from '../generic/GenericStudyArea';

const config: StudyAreaConfig = {
  title: 'Civil Engineering Study Area',
  domain: 'civil-engineering',
  icon: <Engineering sx={{ mr: 2 }} />,
  color: 'primary.main',
  searchPlaceholder: 'e.g., agammaglobulinemia, B cell deficiency, immune response',
  apiEndpoints: {
    search: '/api/civil-engineering/search',
    contentLibrary: '/api/civil-engineering',
    caseStudies: '/api/civil-engineering/case-studies',
    terms: '/api/civil-engineering/medical-terms',
    chapters: '/api/civil-engineering/chapters',
    upload: '/api/civil-engineering/upload-chapter-real'
  },
  additionalTabs: ['Case Studies', 'Medical Terms']
};

const StudyAreaCivilEngineeringPage: React.FC = () => {
  return <GenericStudyArea config={config} />;
};

export default StudyAreaCivilEngineeringPage;
