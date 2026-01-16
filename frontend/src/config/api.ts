import axios from 'axios';

// Use RELATIVE paths - Nginx will handle proxying
// Don't specify a base URL at all
const API_BASE = '';

console.log('🔗 API Routes: Using relative paths (proxied by Nginx)');

// Centralized API endpoints configuration
export const API_ENDPOINTS = {
  coach: {
    respond: '/api/coach/respond',
    status: '/api/coach/status/:jobId',
    mvcDemo: '/api/coach/mvp-demo',
  },
  chisg: {
    query: '/api/chisg/query',
  },
  studyAreas: {
    immunology: {
      search: '/api/immunology/search',
      hsgSearch: '/api/immunology/hsg-search',
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
export const apiClient = {
  get: async <T>(endpoint: string): Promise<T> => {
    // This becomes just: /api/coach/mvp-demo (relative)
    const fullUrl = `${API_BASE}${endpoint}`;
    logRequest(fullUrl);
    
    try {
      const response = await axios.get<T>(fullUrl);
      return response.data;
    } catch (error: any) {
      console.error('API Response Error:', error);
      throw error;
    }
  },

  post: async <T>(endpoint: string, data: any): Promise<T> => {
    const fullUrl = `${API_BASE}${endpoint}`;
    logRequest(fullUrl);
    
    try {
      const response = await axios.post<T>(fullUrl, data);
      return response.data;
    } catch (error: any) {
      console.error('API Response Error:', error);
      throw error;
    }
  },

  upload: async <T>(
    endpoint: string,
    formData: FormData,
    onProgress?: (progressEvent: any) => void
  ): Promise<T> => {
    const fullUrl = `${API_BASE}${endpoint}`;
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
      throw error;
    }
  },

  put: async <T>(endpoint: string, data: any): Promise<T> => {
    const fullUrl = `${API_BASE}${endpoint}`;
    logRequest(fullUrl);
    
    try {
      const response = await axios.put<T>(fullUrl, data);
      return response.data;
    } catch (error: any) {
      console.error('API Response Error:', error);
      throw error;
    }
  },

  delete: async <T>(endpoint: string): Promise<T> => {
    const fullUrl = `${API_BASE}${endpoint}`;
    logRequest(fullUrl);
    
    try {
      const response = await axios.delete<T>(fullUrl);
      return response.data;
    } catch (error: any) {
      console.error('API Response Error:', error);
      throw error;
    }
  }
};