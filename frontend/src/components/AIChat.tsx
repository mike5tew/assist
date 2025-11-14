import React, { useState, useRef, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Paper,
  Chip,
  IconButton,
  Divider,
  Alert,
  CircularProgress,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  List,
  ListItem,
  ListItemText,
  ListItemIcon
} from '@mui/material';
import {
  Send,
  Psychology,
  Clear,
  ExpandMore,
  School,
  Science,
  MenuBook,
  TrendingUp,
  Info,
  SmartToy
} from '@mui/icons-material';
import { apiClient } from '../config/api';

interface ChatMessage {
  id: string;
  type: 'user' | 'ai';
  content: string;
  timestamp: Date;
  queryData?: QueryResponse;
}

interface QueryResponse {
  query: string;
  timestamp: string;
  results_summary: {
    total_exact_matches: number;
    total_semantic_matches: number;
    total_related_skills: number;
  };
  exact_matches: Skill[];
  semantic_matches: SemanticMatch[];
  related_skills: Skill[];
  llm_synthesis: string;
  semantic_insights: SemanticInsights;
}

interface Skill {
  skill_name: string;
  description: string;
  skill_type: string;
  development_age: number;
  criteria_levels: number;
  criteria: string[];
  source_info?: SourceInfo;
}

interface SemanticMatch {
  skill_name: string;
  description: string;
  skill_type: string;
  semantic_relevance: number;
  source_info?: SourceInfo;
}

interface SourceInfo {
  title: string;
  isbn?: string;
  type: string;
  chapter_number?: string;
  chapter_title?: string;
  extraction_method: string;
}

interface SemanticInsights {
  query_complexity: string;
  confidence_levels: {
    high_confidence: number;
    medium_confidence: number;
    low_confidence: number;
  };
}

const AIChat: React.FC = () => {
  const [messages, setMessages] = useState<ChatMessage[]>([
    {
      id: '1',
      type: 'ai',
      content: 'Hello! I can help you explore educational skills and medical content. Ask me about:\n\n• Skill relationships and development\n• Learning objectives and criteria\n• Medical terminology and case studies\n• Educational pathways and assessments\n\nWhat would you like to know?',
      timestamp: new Date()
    }
  ]);
  const [inputValue, setInputValue] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const suggestedQuestions = [
    "What skills are related to working memory?",
    "Tell me about fine motor skills development",
    "What immunology concepts are covered in the textbook?",
    "How do reading skills develop with age?",
    "What are the assessment criteria for problem solving?",
    "Show me case studies about immune system disorders",
    "What skills should a 7-year-old have developed?",
    "Compare semantic search results for emotional intelligence"
  ];

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleSendMessage = async () => {
    if (!inputValue.trim() || loading) return;

    const userMessage: ChatMessage = {
      id: Date.now().toString(),
      type: 'user',
      content: inputValue.trim(),
      timestamp: new Date()
    };

    setMessages(prev => [...prev, userMessage]);
    setInputValue('');
    setLoading(true);
    setError(null);

    try {
      // Query the semantic API - updated to use new API client format
      const queryData = await apiClient.post<QueryResponse>('/api/skills/semantic-query', { 
        query: userMessage.content 
      });
      
      // No need to check response.status or access response.data anymore
      // since our updated apiClient directly returns the data
      
      if (!queryData) {
        throw new Error('No data received from AI system');
      }
      
      // Generate AI response based on the query results
      const aiResponse = generateAIResponse(userMessage.content, queryData);
      
      const aiMessage: ChatMessage = {
        id: (Date.now() + 1).toString(),
        type: 'ai',
        content: aiResponse,
        timestamp: new Date(),
        queryData
      };

      setMessages(prev => [...prev, aiMessage]);
    } catch (err) {
      setError('Sorry, I encountered an error processing your question. Please try again.');
      console.error('Chat error:', err);
    } finally {
      setLoading(false);
    }
  };

  const generateAIResponse = (question: string, data: QueryResponse): string => {
    const { results_summary, exact_matches, semantic_matches, semantic_insights } = data;
    
    let response = `Based on your question "${question}", here's what I found:\n\n`;

    // Summary
    const totalResults = results_summary.total_exact_matches + results_summary.total_semantic_matches;
    if (totalResults === 0) {
      return `I couldn't find any direct matches for "${question}". Try rephrasing your question or asking about specific educational skills, learning objectives, or medical concepts.`;
    }

    response += `📊 **Summary**: Found ${totalResults} relevant results (${results_summary.total_exact_matches} exact matches, ${results_summary.total_semantic_matches} semantic matches)\n\n`;

    // Exact matches
    if (exact_matches.length > 0) {
      response += `🎯 **Direct Matches**:\n`;
      exact_matches.slice(0, 3).forEach((skill, index) => {
        response += `${index + 1}. **${skill.skill_name}** (Age ${skill.development_age}+)\n   ${skill.description.substring(0, 150)}...\n\n`;
      });
    }

    // Semantic matches with high relevance
    const highRelevanceMatches = semantic_matches.filter(match => match.semantic_relevance > 0.8);
    if (highRelevanceMatches.length > 0) {
      response += `🔍 **Highly Relevant Skills**:\n`;
      highRelevanceMatches.slice(0, 3).forEach((match, index) => {
        response += `${index + 1}. **${match.skill_name}** (${(match.semantic_relevance * 100).toFixed(1)}% match)\n   ${match.description.substring(0, 120)}...\n\n`;
      });
    }

    // Source information
    const sources = new Set<string>();
    [...exact_matches, ...semantic_matches].forEach(item => {
      if (item.source_info?.title) {
        sources.add(item.source_info.title);
      }
    });

    if (sources.size > 0) {
      response += `📚 **Sources**: ${Array.from(sources).join(', ')}\n\n`;
    }

    // Insights
    if (semantic_insights.confidence_levels.high_confidence > 0) {
      response += `✨ **Quality**: ${semantic_insights.confidence_levels.high_confidence} high-confidence matches found\n\n`;
    }

    response += `💡 **Tip**: Click on the detailed results below to explore assessment criteria, related skills, and source information.`;

    return response;
  };

  const handleSuggestedQuestion = (question: string) => {
    setInputValue(question);
  };

  const clearChat = () => {
    setMessages([messages[0]]); // Keep the welcome message
    setError(null);
  };

  const handleKeyPress = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      handleSendMessage();
    }
  };

  return (
    <Box sx={{ maxWidth: 1000, mx: 'auto', p: 3, height: '80vh', display: 'flex', flexDirection: 'column' }}>
      <Typography variant="h4" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <SmartToy color="primary" />
        AI Educational Assistant
      </Typography>

      {/* Suggested Questions */}
      <Paper sx={{ p: 2, mb: 2, bgcolor: 'background.default' }}>
        <Typography variant="subtitle2" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Psychology fontSize="small" />
          Suggested Questions:
        </Typography>
        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
          {suggestedQuestions.slice(0, 4).map((question, index) => (
            <Chip
              key={index}
              label={question}
              variant="outlined"
              size="small"
              onClick={() => handleSuggestedQuestion(question)}
              sx={{ cursor: 'pointer' }}
            />
          ))}
        </Box>
      </Paper>

      {/* Chat Messages */}
      <Paper sx={{ flex: 1, p: 2, mb: 2, overflow: 'auto', bgcolor: 'grey.50' }}>
        {messages.map((message) => (
          <Box key={message.id} sx={{ mb: 3 }}>
            <Box
              sx={{
                display: 'flex',
                justifyContent: message.type === 'user' ? 'flex-end' : 'flex-start',
                mb: 1
              }}
            >
              <Paper
                sx={{
                  p: 2,
                  maxWidth: '80%',
                  bgcolor: message.type === 'user' ? 'primary.main' : 'white',
                  color: message.type === 'user' ? 'white' : 'text.primary'
                }}
              >
                <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>
                  {message.content}
                </Typography>
                <Typography variant="caption" sx={{ opacity: 0.7, display: 'block', mt: 1 }}>
                  {message.timestamp.toLocaleTimeString()}
                </Typography>
              </Paper>
            </Box>

            {/* Detailed Results */}
            {message.type === 'ai' && message.queryData && (
              <Box sx={{ ml: 2 }}>
                <Accordion variant="outlined">
                  <AccordionSummary expandIcon={<ExpandMore />}>
                    <Typography variant="subtitle2" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Info fontSize="small" />
                      Detailed Results & Sources
                    </Typography>
                  </AccordionSummary>
                  <AccordionDetails>
                    <QueryResultsDisplay data={message.queryData} />
                  </AccordionDetails>
                </Accordion>
              </Box>
            )}
          </Box>
        ))}

        {loading && (
          <Box sx={{ display: 'flex', justifyContent: 'flex-start', mb: 2 }}>
            <Paper sx={{ p: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
              <CircularProgress size={20} />
              <Typography variant="body2">Thinking...</Typography>
            </Paper>
          </Box>
        )}

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <div ref={messagesEndRef} />
      </Paper>

      {/* Input Area */}
      <Box sx={{ display: 'flex', gap: 1, alignItems: 'flex-end' }}>
        <TextField
          fullWidth
          multiline
          maxRows={3}
          placeholder="Ask about educational skills, learning objectives, or medical content..."
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyPress={handleKeyPress}
          disabled={loading}
        />
        <Button
          variant="contained"
          onClick={handleSendMessage}
          disabled={!inputValue.trim() || loading}
          sx={{ minWidth: 'auto', px: 2 }}
        >
          <Send />
        </Button>
        <IconButton onClick={clearChat} disabled={loading}>
          <Clear />
        </IconButton>
      </Box>
    </Box>
  );
};

// Component to display detailed query results
const QueryResultsDisplay: React.FC<{ data: QueryResponse }> = ({ data }) => {
  return (
    <Box>
      {/* Exact Matches */}
      {data.exact_matches.length > 0 && (
        <Box sx={{ mb: 3 }}>
          <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <School color="primary" />
            Exact Matches ({data.exact_matches.length})
          </Typography>
          <List dense>
            {data.exact_matches.map((skill, index) => (
              <ListItem key={index} sx={{ bgcolor: 'background.paper', mb: 1, borderRadius: 1 }}>
                <ListItemIcon>
                  <TrendingUp color="success" />
                </ListItemIcon>
                <ListItemText
                  primary={`${skill.skill_name} (Age ${skill.development_age}+)`}
                  secondary={
                    <Box>
                      <Typography variant="body2">{skill.description}</Typography>
                      {skill.criteria.length > 0 && (
                        <Typography variant="caption" color="text.secondary">
                          Assessment Levels: {skill.criteria_levels}
                        </Typography>
                      )}
                      {skill.source_info && (
                        <Chip 
                          label={skill.source_info.title} 
                          size="small" 
                          variant="outlined" 
                          sx={{ mt: 1 }} 
                        />
                      )}
                    </Box>
                  }
                />
              </ListItem>
            ))}
          </List>
        </Box>
      )}

      {/* Semantic Matches */}
      {data.semantic_matches.length > 0 && (
        <Box sx={{ mb: 3 }}>
          <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Science color="secondary" />
            Semantic Matches ({data.semantic_matches.length})
          </Typography>
          <List dense>
            {data.semantic_matches.slice(0, 5).map((match, index) => (
              <ListItem key={index} sx={{ bgcolor: 'background.paper', mb: 1, borderRadius: 1 }}>
                <ListItemIcon>
                  <MenuBook color="secondary" />
                </ListItemIcon>
                <ListItemText
                  primary={
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      {match.skill_name}
                      <Chip 
                        label={`${(match.semantic_relevance * 100).toFixed(1)}% match`}
                        size="small"
                        color={match.semantic_relevance > 0.8 ? 'success' : 'default'}
                      />
                    </Box>
                  }
                  secondary={
                    <Box>
                      <Typography variant="body2">{match.description}</Typography>
                      {match.source_info && (
                        <Box sx={{ mt: 1 }}>
                          <Chip 
                            label={match.source_info.title}
                            size="small" 
                            variant="outlined"
                          />
                          {match.source_info.chapter_number && (
                            <Chip 
                              label={`Chapter ${match.source_info.chapter_number}`}
                              size="small"
                              variant="outlined"
                              sx={{ ml: 1 }}
                            />
                          )}
                        </Box>
                      )}
                    </Box>
                  }
                />
              </ListItem>
            ))}
          </List>
        </Box>
      )}

      {/* Query Insights */}
      <Box sx={{ bgcolor: 'background.default', p: 2, borderRadius: 1 }}>
        <Typography variant="subtitle2" gutterBottom>
          Search Quality Insights
        </Typography>
        <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
          <Chip 
            label={`${data.semantic_insights.confidence_levels.high_confidence} High Confidence`}
            color="success"
            size="small"
          />
          <Chip 
            label={`${data.semantic_insights.confidence_levels.medium_confidence} Medium Confidence`}
            color="warning"
            size="small"
          />
          <Chip 
            label={`${data.semantic_insights.confidence_levels.low_confidence} Low Confidence`}
            color="default"
            size="small"
          />
          <Chip 
            label={`Query: ${data.semantic_insights.query_complexity}`}
            variant="outlined"
            size="small"
          />
        </Box>
      </Box>
    </Box>
  );
};

export default AIChat;
