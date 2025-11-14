// this page allows the user to enter a query to the database and see the reposnse
import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Card, 
  CardContent, 
  Grid,
  Alert,
  CircularProgress,
  Tabs,
  Tab,
  TextField,
  Button
} from '@mui/material';
import { useTranslation } from 'react-i18next';
import { apiClient, API_ENDPOINTS } from '../config/api';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 3 }}>
          {children}
        </Box>
      )}
    </div>
  );
}

const Home: React.FC = () => {
  const { t } = useTranslation(); // Use proper hook instead of fallback

  const [tabValue, setTabValue] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [data, setData] = useState<any>(null);
  const [query, setQuery] = useState('');
  const [response, setResponse] = useState('');

  useEffect(() => {
    // Any initialization logic can go here
  }, []);

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const fetchData = async (endpoint: string, description: string) => {
    setLoading(true);
    setError(null);
    
    try {
      console.log(`Fetching data from: ${endpoint}`);
      const result = await apiClient.get(endpoint);
      console.log(`${description} data:`, result);
      setData(result);
    } catch (err) {
      console.error(`Error fetching ${description}:`, err);
      setError(err instanceof Error ? err.message : 'Unknown error occurred');
    } finally {
      setLoading(false);
    }
  };

  const handleTestConnection = () => {
    fetchData('/api/skills/semantic-query?q=test', 'API connection test');
  };

  const handleQuerySubmit = async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.post('/api/skills/semantic-query', { query });
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error occurred');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box sx={{ width: '100%', mt: 4 }}>
      <Typography variant="h4" gutterBottom>
        {t('home.title', 'Database Query')}
      </Typography>
      <Tabs value={tabValue} onChange={handleTabChange} aria-label="home tabs">
        <Tab label={t('home.queryTab', 'Query')} />
        <Tab label={t('home.testTab', 'Test Connection')} />
      </Tabs>
      <TabPanel value={tabValue} index={0}>
        <Card>
          <CardContent>
            <TextField
              label={t('home.queryLabel', 'Enter your query')}
              value={query}
              onChange={e => setQuery(e.target.value)}
              fullWidth
              margin="normal"
            />
            <Button
              variant="contained"
              color="primary"
              onClick={handleQuerySubmit}
              disabled={loading}
            >
              {t('home.submitButton', 'Submit')}
            </Button>
            {loading && <CircularProgress sx={{ mt: 2 }} />}
            {error && <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>}
            {data && (
              <Box sx={{ mt: 2 }}>
                <Typography variant="subtitle1">{t('home.response', 'Response')}:</Typography>
                <pre>{JSON.stringify(data, null, 2)}</pre>
              </Box>
            )}
          </CardContent>
        </Card>
      </TabPanel>
      <TabPanel value={tabValue} index={1}>
        <Button
          variant="contained"
          color="secondary"
          onClick={handleTestConnection}
          disabled={loading}
        >
          {t('home.testButton', 'Test API Connection')}
        </Button>
        {loading && <CircularProgress sx={{ mt: 2 }} />}
        {error && <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>}
        {data && (
          <Box sx={{ mt: 2 }}>
            <Typography variant="subtitle1">{t('home.response', 'Response')}:</Typography>
            <pre>{JSON.stringify(data, null, 2)}</pre>
          </Box>
        )}
      </TabPanel>
    </Box>
  );
};

export default Home;