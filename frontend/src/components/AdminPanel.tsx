import React from 'react';
import { 
  Box, 
  Typography, 
  AppBar,
  Toolbar,
  IconButton,
  Alert
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { ArrowBack, AdminPanelSettings } from '@mui/icons-material';

const AdminPanel: React.FC = () => {
  const navigate = useNavigate();

  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static" sx={{ bgcolor: 'error.main' }}>
        <Toolbar>
          <IconButton
            edge="start"
            color="inherit"
            onClick={() => navigate('/')}
            sx={{ mr: 2 }}
          >
            <ArrowBack />
          </IconButton>
          <AdminPanelSettings sx={{ mr: 2 }} />
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Admin Panel
          </Typography>
        </Toolbar>
      </AppBar>

      <Box sx={{ p: 3 }}>
        <Alert severity="warning">
          <Typography variant="h6" gutterBottom>
            Admin Panel - Access Restricted
          </Typography>
          <Typography>
            This area is for system administrators only. Administrative functions 
            will be available in future versions.
          </Typography>
        </Alert>
      </Box>
    </Box>
  );
};

export default AdminPanel;
