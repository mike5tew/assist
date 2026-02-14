import React from 'react';
import { CssBaseline, Box } from '@mui/material';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { BrowserRouter as Router, Routes, Route, Outlet, useLocation } from 'react-router-dom';

import NavHeader from './components/NavHeader';
import PortfolioLanding from './pages/Portfolio/PortfolioLanding';
import ETPProfilePage from './pages/ETPProfilePage';
import ETPLanding from './pages/ETPLanding';
import LAOLanding from './pages/LAOLanding';
import ESPWorldLanding from './pages/ESPWorld/ESPWorldLanding';
import Home from './components/Home';
import CoachMVPDemo from './components/CoachMVPDemo';
import SemanticQuery from './components/SemanticQuery';
import AIChat from './components/AIChat';
import { ContentLoader } from './components/contentLoader';
import AdminPanel from './components/AdminPanel';
import DiagnosticViewer from './components/DiagnosticViewer';
import ImmunologyUpload from './components/ImmunologyUpload';
import SemanticLinkExtractor from './components/SemanticLinkExtractor';
import ErrorPage from './components/ErrorPage';
import ParentOSLanding from './pages/ESPWorld/ParentOSLanding';
import ParentOSExplorer from './pages/ESPWorld/ParentOSExplorer';
import ParentOSGuide from './pages/ESPWorld/ParentOSGuide';
import ESPWorldDemo from './pages/ESPWorld/ESPWorldDemo';
import LessonPlanningDemo from './pages/ESPWorld/LessonPlanningDemo';
import PrimaryOSLanding from './pages/PrimaryOS';
import CHISGLanding from './pages/CHISG/CHISGLanding';
import CareerOSLanding from './pages/careerOS';
import useAnalytics from './hooks/useAnalytics';
import AnalyticsDashboard from './pages/AnalyticsDashboard';

const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#1a365d', // Professional deep blue
      light: '#2c5282',
      dark: '#1a202c',
    },
    secondary: {
      main: '#38a169', // Trust-worthy green
      light: '#48bb78',
      dark: '#276749',
    },
  },
  typography: {
    fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
    h2: {
      fontWeight: 700,
    },
    h4: {
      fontWeight: 600,
    },
  },
});

// Layout with navigation header
const AppLayout: React.FC = () => (
  <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
    <NavHeader />
    <Box component="main" sx={{ flexGrow: 1 }}>
      <Outlet />
    </Box>
  </Box>
);

// Skills Tree redirect - served by nginx at /skillstree/
const SkillsTreeRedirect: React.FC = () => {
  React.useEffect(() => {
    // /skillstree/ is served by skills-frontend via the main proxy (port 80)
    // Always redirect to the root host (no dev port) to avoid loops when
    // running the portfolio app on http://localhost:3000.
    const { protocol, hostname } = window.location;
    const target = `${protocol}//${hostname}/skillstree/`;

    if (window.location.href !== target) {
      window.location.replace(target);
    }
  }, []);

  return <Box sx={{ p: 4, textAlign: 'center' }}>Redirecting to Skills Map...</Box>;
};

// Scroll to top on route change + page view tracking
const ScrollToTop: React.FC = () => {
  const { pathname } = useLocation();
  useAnalytics(); // auto-tracks page views on route change

  React.useEffect(() => {
    window.scrollTo(0, 0);
  }, [pathname]);

  return null;
};

export default function App() {
  // App is served under '/esp-organizer' inside the main proxy — use basename so direct URLs work.
  const routerBasename = process.env.REACT_APP_BASENAME || (window.location.pathname.startsWith('/esp-organizer') ? '/esp-organizer' : '/');
  
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router basename={routerBasename}>
        <ScrollToTop />
        <Routes>
          {/* Portfolio landing - no header (has its own hero) */}
          <Route path="/" element={<PortfolioLanding />} />
          {/* ETP landing - no header (has its own hero) */}
          <Route path="/etp-landing" element={<ETPLanding />} />
          {/* Primary OS landing - no header (has its own hero) */}
          <Route path="/primary-os" element={<PrimaryOSLanding />} />
          
          {/* LAO landing - no header (has its own hero) */}
          <Route path="/lao" element={<LAOLanding />} />
          {/* ESP World - consolidated skills + weaviate + assistants */}
          <Route path="/esp-world" element={<ESPWorldLanding />} />
          {/* CHISG - information quality landing */}
          <Route path="/chisg" element={<CHISGLanding />} />
          {/* Parent OS - ETP profile + skills + weaviate */}
          <Route path="/parent-os" element={<ParentOSLanding />} />

          {/* CareerOS - verified skills passport landing */}
          <Route path="/careeros/*" element={<CareerOSLanding />} />
          
          {/* Skills Tree redirect - no header wrapper */}
          <Route path="/skillstree" element={<SkillsTreeRedirect />} />
          
          {/* All other routes with navigation header */}
          <Route element={<AppLayout />}>
            {/* Demos */}
            <Route path="/etp-profile" element={<ETPProfilePage />} />
            <Route path="/coach-mvp-demo" element={<CoachMVPDemo />} />
            <Route path="/esp-world/demo" element={<ESPWorldDemo />} />
            <Route path="/esp-world/lesson-demo" element={<LessonPlanningDemo />} />
            <Route path="/parent-os/explorer" element={<ParentOSExplorer />} />
            <Route path="/parent-os/guide" element={<ParentOSGuide />} />
            
            {/* Tools */}
            <Route path="/tools" element={<Home />} />
            <Route path="/semantic-query" element={<SemanticQuery />} />
            <Route path="/ai-chat" element={<AIChat />} />
            <Route path="/contentLoader" element={<ContentLoader />} />
            <Route path="/admin/*" element={<AdminPanel />} />
            
            {/* Utilities */}
            <Route path="/diagnostics" element={<DiagnosticViewer />} />
            <Route path="/upload/immunology" element={<ImmunologyUpload />} />
            <Route path="/semantic-links/extract" element={<SemanticLinkExtractor />} />
            <Route path="/analytics" element={<AnalyticsDashboard />} />
            
            {/* Fallback */}
            <Route path="*" element={<ErrorPage />} />
          </Route>
        </Routes>
      </Router>
    </ThemeProvider>
  );
}


