import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import { ThemeProvider } from '@mui/material/styles';
import { CssBaseline } from '@mui/material';
import { createTheme } from '@mui/material/styles';
const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

const theme = createTheme({
  palette: {
    mode: 'light',
  },
});

// Add debugging to help identify which server we're on
if (process.env.NODE_ENV === 'development') {
  console.log('=======================================');
  console.log('FRONTEND SERVER RUNNING ON PORT 3000');
  console.log('If you see "Welcome to ESP Organizer API", you\'ve been redirected to the API server (8080)');
  console.log('This is a routing issue - you should be on http://localhost:3000/');
  console.log('=======================================');
}

root.render(
  <React.StrictMode>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <App />
      </ThemeProvider>
  </React.StrictMode>
);
