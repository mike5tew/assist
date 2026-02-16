import React from 'react';
import { Box, Container, Typography, Button, Paper, Stack } from '@mui/material';
import { Link } from 'react-router-dom';
import ExploreIcon from '@mui/icons-material/Explore';
import ContactReveal from '../../components/ContactReveal';
import SEO from '../../components/SEO';

const ParentOSExplorer: React.FC = () => {
  return (
    <Box sx={{ py: 6, minHeight: '80vh', bgcolor: '#fafafa' }}>
      <SEO
        title="ParentOS Spectrum Explorer"
        description="Interactive spectrum explorer for parents to visualise and understand their child's position across the 9 ETP biological spectra."
        path="/parent-os/explorer"
      />
      <Container maxWidth="md">
        <Button component={Link} to="/parent-os" variant="text" sx={{ mb: 2 }}>← Back to ParentOS</Button>
        <Paper elevation={0} sx={{ p: 5, textAlign: 'center', borderRadius: 3, border: '1px solid #e0e0e0' }}>
          <ExploreIcon sx={{ fontSize: 64, color: '#7c4dff', mb: 2 }} />
          <Typography variant="h3" gutterBottom>Spectrum Explorer</Typography>
          <Typography variant="h6" color="text.secondary" paragraph>
            Coming Soon
          </Typography>
          <Typography variant="body1" color="text.secondary" paragraph sx={{ maxWidth: 500, mx: 'auto' }}>
            The interactive Spectrum Explorer will let parents visualise and understand their child's position across the 9 biological spectra — with age-appropriate guidance and practical observation tips.
          </Typography>
          <Stack direction="row" spacing={2} justifyContent="center" sx={{ mt: 3 }}>
            <ContactReveal email="world@espthinking.co.uk" label="Register Interest" />
          </Stack>
        </Paper>
      </Container>
    </Box>
  );
};

export default ParentOSExplorer;
