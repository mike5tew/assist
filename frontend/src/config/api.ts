import axios from 'axios';

// Get API base URL from environment or use empty string for relative paths
const API_BASE_URL = process.env.REACT_APP_API_URL || '';

// Centralized API endpoints configuration
export const API_ENDPOINTS = {
  coach: {
    respond: '/api/coach/respond',        // ← NEW
    status: '/api/coach/status/:jobId',   // ← NEW
  },
  chisg: {
    query: '/api/chisg/query',            // ← NEW
  },
  studyAreas: {
    immunology: {
      search: '/api/immunology/search',
      hsgSearch: '/api/immunology/hsg-search', // New HSG search endpoint
      chapters: '/api/immunology/chapters',
      caseStudies: '/api/immunology/case-studies',
      terms: '/api/immunology/medical-terms',
      upload: '/api/immunology/upload-chapter',
      keywordsUpload: '/api/immunology/keywords-upload',
    },
  },
  extraction: {
    status: (jobId: string) => `/api/extraction/status/${jobId}`
  },
  diagnostics: {
    config: '/api/diagnostics/config',
    document: (id: string, collection: string) => `/api/diagnostics/document?id=${encodeURIComponent(id)}&collection=${encodeURIComponent(collection)}`,
    collections: '/api/diagnostics/collections',
  }
};

// Helper functions for logging and error handling
const logRequest = (endpoint: string) => {
  console.log('Making API request to:', endpoint);
};

// API helper methods with consistent return patterns
export const apiHelpers = {
  get: async <T>(endpoint: string): Promise<T> => {
    const fullUrl = `${API_BASE_URL}${endpoint}`;
    logRequest(fullUrl);
    
    try {
      const response = await axios.get<T>(fullUrl);
      return response.data;
    } catch (error: any) {
      console.error('API Response Error:', error);
      console.error('Request failed for:', endpoint);
      throw error;
    }
  },

  post: async <T>(endpoint: string, data: any): Promise<T> => {
    const fullUrl = `${API_BASE_URL}${endpoint}`;
    logRequest(fullUrl);
    
    try {
      const response = await axios.post<T>(fullUrl, data);
      return response.data;
    } catch (error: any) {
      console.error('API Response Error:', error);
      console.error('POST request failed for:', endpoint);
      throw error;
    }
  },

  upload: async <T>(
    endpoint: string, 
    formData: FormData, 
    onProgress?: (progressEvent: any) => void
  ): Promise<T> => {
    const fullUrl = `${API_BASE_URL}${endpoint}`;
    logRequest(fullUrl);
    
    try {
      const response = await axios.post<T>(fullUrl, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        onUploadProgress: onProgress,
      });
      return response.data;
    } catch (error: any) {
      console.error('API Response Error:', error);
      console.error('Upload request failed for:', endpoint);
      throw error;
    }
  }
};

// Legacy API client for backward compatibility
// Make sure it's consistent with the new helpers
export const apiClient = {
  get: apiHelpers.get,
  post: apiHelpers.post,
  upload: apiHelpers.upload
};