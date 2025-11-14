import React from 'react';
import { CssBaseline } from '@mui/material';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

import AppRoutes from './routes';
import DiagnosticViewer from './components/DiagnosticViewer';
import ImmunologyUpload from './components/ImmunologyUpload';
import Home from './components/Home';

const theme = createTheme({
  palette: {
    mode: 'light',
  },
});

export default function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <BrowserRouter>
        <div className="App">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/upload/immunology" element={<ImmunologyUpload />} />
            <Route path="/diagnostics" element={<DiagnosticViewer />} />
            <Route path="/*" element={<AppRoutes />} />
          </Routes>
        </div>
      </BrowserRouter>
    </ThemeProvider>
  );
}


