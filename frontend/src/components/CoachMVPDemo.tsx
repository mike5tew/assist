import React, { useState, useRef, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Paper,
  CircularProgress,
  Alert,
  Divider,
  Chip,
  Container
} from '@mui/material';
import {
  Send,
  SmartToy,
  Clear,
  Info
} from '@mui/icons-material';
import { apiClient } from '../config/api';

interface ChatMessage {
  id: string;
  type: 'user' | 'ai';
  content: string;
  timestamp: Date;
}

interface SemanticLink {
  source_term: string;
  target_term: string;
  relation_type: string;
  context: string;
}

interface CoachMVPDemoResponse {
  message: string;
  response_id: string;
  timestamp: string;
  knowledge_graph_routes: SemanticLink[];
}

const CoachMVPDemo: React.FC = () => {
  const [query, setQuery] = useState('');
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [response, setResponse] = useState<CoachMVPDemoResponse | null>(null);
  const responseEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    responseEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [response, messages]);

  const handleSendQuery = async () => {
    if (!query.trim() || loading) return;

    const userMessage: ChatMessage = {
      id: Date.now().toString(),
      type: 'user',
      content: query.trim(),
      timestamp: new Date()
    };

    setMessages((prev) => [...prev, userMessage]);
    setQuery('');
    setLoading(true);
    setError(null);

    try {
      // FIXED: Include /api prefix in endpoint
      const queryData = await apiClient.post<CoachMVPDemoResponse>(
        '/api/coach/mvp-demo',  // ← Changed from '/coach/mvp-demo'
        { message: userMessage.content }
      );

      if (!queryData) {
        throw new Error('No response from server');
      }

      setResponse(queryData);

      const aiMessage: ChatMessage = {
        id: (Date.now() + 1).toString(),
        type: 'ai',
        content: queryData.message,
        timestamp: new Date()
      };
      setMessages((prev) => [...prev, aiMessage]);

    } catch (err) {
      setError('Sorry, I encountered an error processing your question. Please try again.');
      console.error('Chat error:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleKeyPress = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      handleSendQuery();
    }
  };

  const handleClear = () => {
    setQuery('');
    setMessages([]);
    setResponse(null);
    setError(null);
  };

  const suggestedQueries = [
    'What is XLA?',
    'Tell me about X-Linked Agammaglobulinemia',
    'How does XLA affect the immune system?',
    'What are the symptoms of XLA?'
  ];

  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      {/* Header */}
      <Box sx={{ mb: 4, textAlign: 'center' }}>
        <Typography variant="h3" gutterBottom sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 1, mb: 2 }}>
          <SmartToy sx={{ fontSize: 40, color: 'primary.main' }} />
          MVP Coach Demo
        </Typography>
        <Typography variant="subtitle1" color="text.secondary">
          Explore the knowledge graph using natural language queries
        </Typography>
      </Box>

      {/* Info Alert */}
      <Alert icon={<Info />} severity="info" sx={{ mb: 3 }}>
        This demo shows how the Coach system retrieves and connects related concepts from a seeded knowledge graph. Try asking about XLA or related immunology concepts.
      </Alert>

      {/* Suggested Queries */}
      <Paper sx={{ p: 2, mb: 3, bgcolor: 'background.default' }}>
        <Typography variant="subtitle2" gutterBottom>
          Try these queries:
        </Typography>
        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
          {suggestedQueries.map((suggestedQuery, index) => (
            <Chip
              key={index}
              label={suggestedQuery}
              onClick={() => setQuery(suggestedQuery)}
              variant="outlined"
              sx={{ cursor: 'pointer' }}
            />
          ))}
        </Box>
      </Paper>

      {/* Query Input */}
      <Paper sx={{ p: 3, mb: 3 }}>
        <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
          <TextField
            fullWidth
            multiline
            maxRows={3}
            placeholder="Ask a question about medical concepts, immunology, or learning..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyPress={handleKeyPress}
            disabled={loading}
            variant="outlined"
          />
          <Button
            variant="contained"
            onClick={handleSendQuery}
            disabled={!query.trim() || loading}
            sx={{ minWidth: 120 }}
          >
            {loading ? <CircularProgress size={24} /> : <Send />}
          </Button>
          <Button
            variant="outlined"
            onClick={handleClear}
            disabled={loading}
            sx={{ minWidth: 100 }}
          >
            <Clear />
          </Button>
        </Box>
      </Paper>

      {/* Error Display */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Response Display */}
      {response && !loading && (
        <Card sx={{ mb: 3 }}>
          <CardContent>
            {/* Coach Message */}
            <Box sx={{ mb: 3, p: 2, bgcolor: 'primary.light', borderRadius: 1 }}>
              <Typography variant="body1" sx={{ color: 'primary.contrastText', whiteSpace: 'pre-wrap' }}>
                {response.message}
              </Typography>
            </Box>

            <Divider sx={{ my: 2 }} />

            {/* Knowledge Graph Routes */}
            <Typography variant="h6" gutterBottom sx={{ mt: 3, mb: 2 }}>
              Connected Concepts ({response.knowledge_graph_routes.length})
            </Typography>

            {response.knowledge_graph_routes.length > 0 ? (
              <Box>
                {response.knowledge_graph_routes.map((link, index) => (
                  <Card
                    key={index}
                    variant="outlined"
                    sx={{ mb: 2, p: 2, bgcolor: 'background.default' }}
                  >
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 1 }}>
                      <Chip
                        label={link.source_term}
                        color="primary"
                        variant="outlined"
                      />
                      <Typography variant="body2" color="text.secondary" sx={{ fontWeight: 'bold' }}>
                        {link.relation_type}
                      </Typography>
                      <Chip
                        label={link.target_term}
                        color="secondary"
                        variant="outlined"
                      />
                    </Box>
                    <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                      <strong>Context:</strong> {link.context}
                    </Typography>
                  </Card>
                ))}
              </Box>
            ) : (
              <Typography color="text.secondary">
                No connected concepts found. Try refining your query.
              </Typography>
            )}

            {/* Metadata */}
            <Divider sx={{ my: 2 }} />
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <Typography variant="caption" color="text.secondary">
                Response ID: {response.response_id}
              </Typography>
              <Typography variant="caption" color="text.secondary">
                {new Date(response.timestamp).toLocaleTimeString()}
              </Typography>
            </Box>
          </CardContent>
        </Card>
      )}

      {/* Empty State */}
      {!response && !loading && !error && (
        <Paper sx={{ p: 4, textAlign: 'center', bgcolor: 'background.default' }}>
          <SmartToy sx={{ fontSize: 48, color: 'text.disabled', mb: 2 }} />
          <Typography color="text.secondary">
            Enter a query above to explore the knowledge graph
          </Typography>
        </Paper>
      )}

      <div ref={responseEndRef} />
    </Container>
  );
};

export default CoachMVPDemo;
