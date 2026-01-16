import React from 'react';
import { CssBaseline } from '@mui/material';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

import AppRoutes from './routes';
import DiagnosticViewer from './components/DiagnosticViewer';
import ImmunologyUpload from './components/ImmunologyUpload';
import Home from './components/Home';
import CoachMVPDemo from './components/CoachMVPDemo';
import SemanticLinkExtractor from './components/SemanticLinkExtractor';

const theme = createTheme({
  palette: {
    mode: 'light',
  },
});

export default function App() {
  // App is served under '/esp-organizer' inside the main proxy — use basename so direct URLs work
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router basename="/esp-organizer">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/upload/immunology" element={<ImmunologyUpload />} />
          <Route path="/diagnostics" element={<DiagnosticViewer />} />
          <Route path="/coach-mvp-demo" element={<CoachMVPDemo />} />
          <Route path="/semantic-links/extract" element={<SemanticLinkExtractor />} />
          <Route path="/*" element={<AppRoutes />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
}


