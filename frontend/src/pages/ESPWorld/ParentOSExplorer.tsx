import React from 'react';
import { Box, Container, Typography, Button } from '@mui/material';
import { Link } from 'react-router-dom';

const ParentOSExplorer: React.FC = () => {
  return (
    <Box sx={{ py: 6 }}>
      <Container maxWidth="md">
        <Button component={Link} to="/parent-os" variant="text" sx={{ mb: 2 }}>← Back to ParentOS</Button>
        <Typography variant="h3" gutterBottom>ParentOS — Spectrum Explorer</Typography>
        <Typography variant="body1" color="text.secondary" paragraph>
          This Explorer provides interactive access to the 9 Biological Spectra for parents. Use it to explore spectrum positions, slider states, and age-stage guidance.
        </Typography>
        <Typography variant="body2" color="text.secondary">
          (Placeholder content) Replace with interactive explorer component or embedded tool as required.
        </Typography>
      </Container>
    </Box>
  );
};

export default ParentOSExplorer;
