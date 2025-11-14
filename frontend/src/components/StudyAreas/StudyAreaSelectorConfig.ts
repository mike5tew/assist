import { School, Engineering, LocalHospital } from '@mui/icons-material';
import React from 'react';

export interface StudyArea {
  id: string;
  title: string;
  description: string;
  icon: React.ElementType;
  path: string;
  status: 'active' | 'beta' | 'coming-soon';
  features: string[];
}

// This configuration can now be imported by StudyAreaSelector component
export const studyAreas: StudyArea[] = [
  {
    id: 'gcse',
    title: 'GCSE Studies',
    description: 'Comprehensive GCSE curriculum with interactive learning materials and assessment tools.',
    icon: School,
    path: '/gcse',
    status: 'coming-soon',
    features: ['Subject-specific content', 'Practice tests', 'Progress tracking', 'Revision notes']
  },
  {
    id: 'civil-engineering',
    title: 'Civil Engineering Degree',
    description: 'Advanced engineering concepts, calculations, and practical applications for degree-level studies.',
    icon: Engineering,
    path: '/civil-engineering',
    status: 'coming-soon',
    features: ['Structural analysis', 'Design principles', 'Project management', 'Technical standards']
  },
  {
    id: 'immunology',
    title: 'Medical Immunology',
    description: 'Specialized immunology content with case studies, medical terminology, and clinical scenarios.',
    icon: LocalHospital,
    path: '/immunology',
    status: 'active',
    features: ['Case studies', 'Medical terminology', 'Clinical scenarios', 'Textract processing']
  }
];
