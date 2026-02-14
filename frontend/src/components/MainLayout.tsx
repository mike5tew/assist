import React from 'react';
import { Box } from '@mui/material';
import { Outlet } from 'react-router-dom';
import NavHeader from './NavHeader';

/**
 * MainLayout - Wraps pages with the navigation header
 * Used for multipage navigation instead of drawer-based navigation
 */
const MainLayout: React.FC = () => {
  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <NavHeader />
      <Box component="main" sx={{ flexGrow: 1 }}>
        <Outlet />
      </Box>
    </Box>
  );
};

export default MainLayout;
