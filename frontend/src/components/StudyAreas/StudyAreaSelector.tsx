import React from 'react';
import { 
  Box, 
  Typography, 
  Card, 
  CardContent, 
  Grid,
  Button,
  CardActions,
  Chip
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { AdminPanelSettings } from '@mui/icons-material';
import { studyAreas, StudyArea } from './StudyAreaSelectorConfig';

const StudyAreaSelector: React.FC = () => {
  const navigate = useNavigate();

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'success';
      case 'beta': return 'warning';
      case 'coming-soon': return 'default';
      default: return 'default';
    }
  };

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'active': return 'Available';
      case 'beta': return 'Beta';
      case 'coming-soon': return 'Coming Soon';
      default: return status;
    }
  };

  return (
    <Box sx={{ flexGrow: 1, p: 4, minHeight: '100vh', bgcolor: 'grey.50' }}>
      <Box sx={{ textAlign: 'center', mb: 6 }}>
        <Typography variant="h3" component="h1" gutterBottom color="primary">
          ESP Study Areas
        </Typography>
        <Typography variant="h6" color="text.secondary" sx={{ maxWidth: 800, mx: 'auto' }}>
          Choose your study area to access specialized content, tools, and resources tailored to your educational needs.
        </Typography>
      </Box>

      <Grid container spacing={4} sx={{ maxWidth: 1200, mx: 'auto' }}>
        {studyAreas.map((area) => (
          <Grid item xs={12} md={4} key={area.id}>
            <Card 
              sx={{ 
                height: '100%', 
                display: 'flex', 
                flexDirection: 'column',
                transition: 'transform 0.2s, box-shadow 0.2s',
                '&:hover': {
                  transform: 'translateY(-4px)',
                  boxShadow: 4
                }
              }}
            >
              <CardContent sx={{ flexGrow: 1, textAlign: 'center', p: 3 }}>
                <Box sx={{ color: 'primary.main', mb: 2 }}>
                  {area.icon && React.createElement(area.icon as React.ElementType, { fontSize: "large" })}
                </Box>
                
                <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: 1, mb: 2 }}>
                  <Typography variant="h5" component="h2">
                    {area.title}
                  </Typography>
                  <Chip 
                    label={getStatusLabel(area.status)} 
                    color={getStatusColor(area.status) as any}
                    size="small"
                  />
                </Box>
                
                <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
                  {area.description}
                </Typography>
                
                <Box sx={{ textAlign: 'left' }}>
                  <Typography variant="subtitle2" gutterBottom>
                    Features:
                  </Typography>
                  {area.features.map((feature, index) => (
                    <Typography key={index} variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>
                      • {feature}
                    </Typography>
                  ))}
                </Box>
              </CardContent>
              
              <CardActions sx={{ p: 3, pt: 0 }}>
                <Button 
                  fullWidth
                  variant={area.status === 'active' ? 'contained' : 'outlined'}
                  size="large"
                  disabled={area.status === 'coming-soon'}
                  onClick={() => navigate(area.path)}
                >
                  {area.status === 'active' ? 'Enter Study Area' : 
                   area.status === 'beta' ? 'Try Beta Version' : 
                   'Coming Soon'}
                </Button>
              </CardActions>
            </Card>
          </Grid>
        ))}
      </Grid>

      {/* Admin access */}
      <Box sx={{ textAlign: 'center', mt: 6, pt: 4, borderTop: 1, borderColor: 'divider' }}>
        <Button
          variant="text"
          startIcon={<AdminPanelSettings />}
          onClick={() => navigate('/admin')}
          sx={{ color: 'text.secondary' }}
        >
          Admin Panel
        </Button>
      </Box>
    </Box>
  );
};

export default StudyAreaSelector;
