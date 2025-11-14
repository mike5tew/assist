import React from 'react';
import { School } from '@mui/icons-material';
import GenericStudyArea, { StudyAreaConfig } from '../generic/GenericStudyArea';

const config: StudyAreaConfig = {
  title: 'GCSE Study Area',
  domain: 'GCSE',
  icon: <School sx={{ mr: 2 }} />,
  color: 'primary.main',
  apiEndpoints: {
    search: '/api/gcse/search',
    contentLibrary: '/api/gcse',
    terms: '/api/gcse/keywords',
    chapters: '/api/gcse/chapters',
    upload: '/api/gcse/upload-chapter-real'
  },
  additionalTabs: ['Subjects', 'Keywords'],
  searchPlaceholder: 'Search GCSE content...'
};

const StudyAreaGCSEPage: React.FC = () => {
  return <GenericStudyArea config={config} />;
};

export default StudyAreaGCSEPage;
