export interface ChapterInfo {
  chapter_number: string;
  chapter_title: string;
  // Add more fields as needed
}

export interface Source {
  _id?: string;
  title: string;
  isbn: string;
  authors?: string[];
  publisher?: string;
  year?: number;
  type?: string;
  // Add more fields as needed
}

/**
 * Type definitions for the ESP Organizer frontend
 * These correspond to structures in the backend's types.go
 */

// Chapter represents a chapter from a book or source
export interface Chapter {
  _id?: string;
  chapter_number: string;
  chapter_title: string;
  book_title: string;
  isbn: string;
  case_studies_count: number;
  medical_terms_count: number;
  extracted_at: string;
  job_id?: string;
}

// PatientInfo represents information about a patient in a case study
export interface PatientInfo {
  age?: string;
  gender?: string;
  sex?: string;
  medical_history?: string | string[];
  presenting_symptoms?: string[];
}

// CaseStudy represents a clinical case study
export interface CaseStudy {
  _id?: string;
  case_number: string;
  content: string;
  patient_info?: PatientInfo;
  clinical_findings?: string[];
  discussion_points?: string[];
  chapter_number?: string;
  book_isbn?: string;
  extracted_at?: string;
}

// MedicalTerm represents a medical term or definition
export interface MedicalTerm {
  _id?: string;
  term: string;
  category: string;
  frequency: number;
  context: string[];
  chapter_number?: string;
  book_isbn?: string;
  extracted_at?: string;
}

// Source represents a reference source like a book
export interface Source {
  _id?: string;
  title: string;
  isbn: string;
  type: string;
  authors?: string[];
  publisher?: string;
  year?: number;
  created_at?: string;
  updated_at?: string;
}

// SearchResponse represents a search API response
export interface SearchResponse {
  ai_response: {
    type: 'success' | 'error' | 'not_found' | 'medical_summary' | 'specific_condition';
    content?: string;
    confidence?: string;
    statistics?: {
      chapters_processed?: number;
      case_studies_found?: number;
      medical_terms_found?: number;
      relevant_cases?: number;
      relevant_terms?: number;
    };
  };
  domain?: string;
  query?: string;
  timestamp?: string;
}

// types/routes.ts
export type RoutePath = 
  | '/' 
  | '/diagnostics'
  | '/other-route';

export interface RouteConfig {
  path: RoutePath;
  element: React.ReactElement;
  label?: string;
}