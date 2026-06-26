import React from 'react';
import { Box, Container, Typography, Button, Paper, Stack } from '@mui/material';
import { Link } from 'react-router-dom';
import MenuBookIcon from '@mui/icons-material/MenuBook';
import ContactReveal from '../../components/ContactReveal';
import SEO from '../../components/SEO';

const ToddlerOSGuide: React.FC = () => {
  return (
    <Box sx={{ py: 6, minHeight: '80vh', bgcolor: '#fafafa' }}>
      <SEO
        title="ToddlerOS Parent Guide"
        description="Practical articles, observation frameworks, and resources to help parents support their child's development across the 9 biological spectra."
        path="/toddler-os/guide"
      />
      <Container maxWidth="md">
        <Button component={Link} to="/toddler-os" variant="text" sx={{ mb: 2 }}>← Back to ToddlerOS</Button>
        <Paper elevation={0} sx={{ p: 5, textAlign: 'center', borderRadius: 3, border: '1px solid #e0e0e0' }}>
          <MenuBookIcon sx={{ fontSize: 64, color: '#7c4dff', mb: 2 }} />
          <Typography variant="h3" gutterBottom>Parent Guide</Typography>
          <Typography variant="h6" color="text.secondary" paragraph>
            Coming Soon
          </Typography>
          <Typography variant="body1" color="text.secondary" paragraph sx={{ maxWidth: 500, mx: 'auto' }}>
            The Parent Guide will offer practical articles, observation frameworks, and downloadable resources to help you support your child's development across the 9 spectra at home.
          </Typography>
          <Stack direction="row" spacing={2} justifyContent="center" sx={{ mt: 3 }}>
            <ContactReveal email="world@espthinking.co.uk" label="Register Interest" />
          </Stack>
        </Paper>
      </Container>
    </Box>
  );
};

export default ToddlerOSGuide;
