/**
 * services/semanticLinkService.ts
 * 
 * Service layer for semantic link extraction operations
 * Handles all communication with the backend API
 */

import axios from 'axios';
import { SemanticLinkSelection, ExtractionResult, SearchResult } from '../types/semanticLinks';

const API_BASE = process.env.REACT_APP_API_URL || '/api';

/**
 * SemanticLinkExtractionService - Singleton service for semantic link operations
 * 
 * Provides methods to:
 * - Extract semantic links from user selections
 * - Search the knowledge graph
 * - Validate extracted links
 */
export class SemanticLinkExtractionService {
  private static instance: SemanticLinkExtractionService;
  private apiClient = axios.create({
    baseURL: API_BASE,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  private constructor() {}

  /**
   * Get singleton instance
   */
  static getInstance(): SemanticLinkExtractionService {
    if (!this.instance) {
      this.instance = new SemanticLinkExtractionService();
    }
    return this.instance;
  }

  /**
   * Extract semantic links from document
   * 
   * POST /api/semantic-links/extract
   * 
   * @param documentId - ID of the document being processed
   * @param selections - Array of semantic link selections from user
   * @returns Extraction result with count and errors
   */
  async extractSemanticLinks(
    documentId: string,
    selections: SemanticLinkSelection[]
  ): Promise<ExtractionResult> {
    try {
      const response = await this.apiClient.post('/semantic-links/extract', {
        document_id: documentId,
        selections: selections.map((sel) => ({
          // New statement-centric fields
          statement: sel.statement,
          source_term: sel.sourceTerm,
          target_term: sel.targetTerm,
          forward_relation: sel.forwardRelation,
          inverse_relation: sel.inverseRelation,
          conditions: sel.conditions?.map((cond) => ({
            term: cond.term,
            forward_relation: cond.forwardRelation,
            inverse_relation: cond.inverseRelation,
            required: cond.required,
            description: cond.description,
          })),
          // Metadata
          confidence: sel.confidence,
          quality_flag: sel.qualityFlag,
          position: sel.position,
        })),
      });

      return {
        inserted: response.data.inserted || 0,
        errors: response.data.errors || [],
      };
    } catch (error) {
      const errorMsg = axios.isAxiosError(error)
        ? error.response?.data?.message || error.message
        : 'Unknown error occurred';
      throw new Error(`Failed to extract semantic links: ${errorMsg}`);
    }
  }

  /**
   * Search semantic links by term
   * 
   * GET /api/semantic-links/search?term=...&domain=...&max=10
   * 
   * @param query - Search term
   * @param domain - Optional domain filter (e.g., 'medical', 'educational')
   * @param maxResults - Maximum number of results to return
   * @returns Search results with semantic context and answers
   */
  async searchSemanticLinks(
    query: string,
    domain?: string,
    maxResults: number = 10
  ): Promise<SearchResult> {
    try {
      const params = new URLSearchParams({
        term: query,
        max: maxResults.toString(),
      });

      if (domain) {
        params.append('domain', domain);
      }

      const response = await this.apiClient.get('/semantic-links/search', {
        params: Object.fromEntries(params),
      });

      return {
        semanticContext: response.data.semantic_context || [],
        summaryContext: response.data.summary_context || [],
        documents: response.data.documents || [],
        answer: response.data.answer || '',
      };
    } catch (error) {
      const errorMsg = axios.isAxiosError(error)
        ? error.response?.data?.message || error.message
        : 'Unknown error occurred';
      throw new Error(`Failed to search semantic links: ${errorMsg}`);
    }
  }

  /**
   * Validate extracted semantic links using LLM quality scoring
   * 
   * POST /api/semantic-links/validate
   * 
   * @param linkId - ID of the link to validate
   * @returns Validation result with quality score
   */
  async validateSemanticLink(linkId: string): Promise<{
    linkId: string;
    qualityScore: number;
    feedback: string;
    isValid: boolean;
  }> {
    try {
      const response = await this.apiClient.post('/semantic-links/validate', {
        link_id: linkId,
      });

      return {
        linkId: response.data.link_id,
        qualityScore: response.data.quality_score || 0,
        feedback: response.data.feedback || '',
        isValid: response.data.is_valid || false,
      };
    } catch (error) {
      const errorMsg = axios.isAxiosError(error)
        ? error.response?.data?.message || error.message
        : 'Unknown error occurred';
      throw new Error(`Failed to validate semantic link: ${errorMsg}`);
    }
  }

  /**
   * Batch validate multiple semantic links
   * 
   * @param linkIds - Array of link IDs to validate
   * @returns Array of validation results
   */
  async batchValidateLinks(
    linkIds: string[]
  ): Promise<Array<{
    linkId: string;
    qualityScore: number;
    feedback: string;
    isValid: boolean;
  }>> {
    try {
      const responses = await Promise.all(
        linkIds.map((id) => this.validateSemanticLink(id))
      );
      return responses;
    } catch (error) {
      throw new Error(`Failed to batch validate links: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  }

  /**
   * Health check endpoint
   */
  async healthCheck(): Promise<boolean> {
    try {
      const response = await this.apiClient.get('/health');
      return response.status === 200;
    } catch {
      return false;
    }
  }
}
