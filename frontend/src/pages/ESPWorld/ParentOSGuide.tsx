import React from 'react';
import { Box, Container, Typography, Button } from '@mui/material';
import { Link } from 'react-router-dom';

const ParentOSGuide: React.FC = () => {
  return (
    <Box sx={{ py: 6 }}>
      <Container maxWidth="md">
        <Button component={Link} to="/parent-os" variant="text" sx={{ mb: 2 }}>← Back to ParentOS</Button>
        <Typography variant="h3" gutterBottom>ParentOS — Parent Guide</Typography>
        <Typography variant="body1" color="text.secondary" paragraph>
          The Parent Guide gives practical advice for observing and supporting your child across the 9 spectra. It is a companion to the interactive Explorer and the classroom PrimaryOS tools.
        </Typography>
        <Typography variant="body2" color="text.secondary">
          (Placeholder content) Add guide articles, downloads, and references here when ready.
        </Typography>
      </Container>
    </Box>
  );
};

export default ParentOSGuide;
