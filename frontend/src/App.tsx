import React from 'react';
import { CssBaseline, Box } from '@mui/material';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { BrowserRouter as Router, Routes, Route, Outlet, useLocation } from 'react-router-dom';
import { HelmetProvider } from 'react-helmet-async';

import NavHeader from './components/NavHeader';
import PortfolioLanding from './pages/Portfolio/PortfolioLanding';
import ETPProfilePage from './pages/ETPProfilePage';
import ETPLanding from './pages/ETPLanding';
import LAOLanding from './pages/LAOLanding';
import LAOPrivacyPolicy from './pages/LAOPrivacyPolicy';
import LAOTerms from './pages/LAOTerms';
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
import ToddlerOSLanding from './pages/ESPWorld/ToddlerOSLanding';
import ToddlerOSExplorer from './pages/ESPWorld/ToddlerOSExplorer';
import ToddlerOSGuide from './pages/ESPWorld/ToddlerOSGuide';
import ESPWorldDemo from './pages/ESPWorld/ESPWorldDemo';
import LessonPlanningDemo from './pages/ESPWorld/LessonPlanningDemo';
import PrimaryOSLanding from './pages/PrimaryOS';
import CHISGLanding from './pages/CHISG/CHISGLanding';
import ESPPilotLanding from './pages/ESPPilotLanding';
import GraphExplorer from './pages/CHISG/GraphExplorer';
import NTMResearch from './pages/CHISG/NTMResearch';
import CHISGLogin from './pages/CHISG/CHISGLogin';
import CHISGPapers from './pages/CHISG/CHISGPapers';
import CHISGLinks from './pages/CHISG/CHISGLinks';
import { AuthProvider } from './context/AuthContext';
import PrivateRoute from './components/PrivateRoute';
import ClassifierLanding from './pages/Classifier/ClassifierLanding';
import CareerOSLanding from './pages/careerOS';
import JoinLanding from './pages/JoinLanding';
import BlogIndex from './pages/Blog/BlogIndex';
import SocialGravity from './pages/Blog/SocialGravity';
import EnergyDirectionality from './pages/Blog/EnergyDirectionality';
import VoltageSensitivity from './pages/Blog/VoltageSensitivity';
import ThreatResponse from './pages/Blog/ThreatResponse';
import CareResponse from './pages/Blog/CareResponse';
import RiskTolerance from './pages/Blog/RiskTolerance';
import IntegrityLogic from './pages/Blog/IntegrityLogic';
import MirrorNeuronTuning from './pages/Blog/MirrorNeuronTuning';
import Orderliness from './pages/Blog/Orderliness';
import PilotStrength from './pages/Blog/PilotStrength';
import CurrentLoad from './pages/Blog/CurrentLoad';
import Libido from './pages/Blog/Libido';
import TheGovernor from './pages/Blog/TheGovernor';
import useAnalytics from './hooks/useAnalytics';
import AnalyticsDashboard from './pages/AnalyticsDashboard';
import APhysicsRevision from './pages/APhysicsRevision';

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
    <HelmetProvider>
    <AuthProvider>
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
          <Route path="/lao/privacy" element={<LAOPrivacyPolicy />} />
          <Route path="/lao/terms" element={<LAOTerms />} />
          {/* A-Level Physics Revision - no header (full screen) */}
          <Route path="/aphy" element={<APhysicsRevision />} />
          {/* ESP World - consolidated skills + weaviate + assistants */}
          <Route path="/esp-world" element={<ESPWorldLanding />} />
          {/* CHISG - information quality landing */}
          <Route path="/chisg" element={<CHISGLanding />} />
          {/* ESP Pilot - research pilot study landing */}
          <Route path="/esp-pilot" element={<ESPPilotLanding />} />
          <Route path="/chisg/graph" element={<GraphExplorer />} />
          {/* CHISG auth + protected admin pages */}
          <Route path="/chisg/login" element={<CHISGLogin />} />
          <Route path="/chisg/papers" element={<PrivateRoute><CHISGPapers /></PrivateRoute>} />
          <Route path="/chisg/links" element={<PrivateRoute><CHISGLinks /></PrivateRoute>} />
          {/* NTM Research - academic knowledge graph for mycobacteria */}
          <Route path="/ntm" element={<PrivateRoute><NTMResearch /></PrivateRoute>} />
          {/* Classifier - load balanced pipeline demo */}
          <Route path="/classifier-demo" element={<ClassifierLanding />} />
          {/* ToddlerOS - ETP for parents of young children */}
          <Route path="/toddler-os" element={<ToddlerOSLanding />} />
          {/* Legacy redirect */}
          <Route path="/parent-os" element={<ToddlerOSLanding />} />

          {/* CareerOS - verified skills passport landing */}
          <Route path="/careeros/*" element={<CareerOSLanding />} />

          {/* Join - email capture landing for LinkedIn posts */}
          <Route path="/join" element={<JoinLanding />} />
          
          {/* Blog - ETP deep-dive series */}
          <Route path="/blog" element={<BlogIndex />} />
          <Route path="/blog/social-gravity" element={<SocialGravity />} />
          <Route path="/blog/energy-directionality" element={<EnergyDirectionality />} />
          <Route path="/blog/voltage-sensitivity" element={<VoltageSensitivity />} />
          <Route path="/blog/threat-response" element={<ThreatResponse />} />
          <Route path="/blog/care-response" element={<CareResponse />} />
          <Route path="/blog/risk-tolerance" element={<RiskTolerance />} />
          <Route path="/blog/integrity-logic" element={<IntegrityLogic />} />
          <Route path="/blog/mirror-neuron-tuning" element={<MirrorNeuronTuning />} />
          <Route path="/blog/orderliness" element={<Orderliness />} />
          <Route path="/blog/pilot-strength" element={<PilotStrength />} />
          <Route path="/blog/current-load" element={<CurrentLoad />} />
          <Route path="/blog/libido" element={<Libido />} />
          <Route path="/blog/the-governor" element={<TheGovernor />} />

          {/* Skills Tree redirect - no header wrapper */}
          <Route path="/skillstree" element={<SkillsTreeRedirect />} />
          
          {/* All other routes with navigation header */}
          <Route element={<AppLayout />}>
            {/* Demos */}
            <Route path="/etp-profile" element={<ETPProfilePage />} />
            <Route path="/coach-mvp-demo" element={<CoachMVPDemo />} />
            <Route path="/esp-world/demo" element={<ESPWorldDemo />} />
            <Route path="/esp-world/lesson-demo" element={<LessonPlanningDemo />} />
            <Route path="/toddler-os/explorer" element={<ToddlerOSExplorer />} />
            <Route path="/toddler-os/guide" element={<ToddlerOSGuide />} />
            {/* Legacy redirects */}
            <Route path="/parent-os/explorer" element={<ToddlerOSExplorer />} />
            <Route path="/parent-os/guide" element={<ToddlerOSGuide />} />
            
            {/* Tools */}
            <Route path="/tools" element={<Home />} />
            <Route path="/semantic-query" element={<SemanticQuery />} />
            <Route path="/ai-chat" element={<AIChat />} />
            <Route path="/contentLoader" element={<ContentLoader />} />
            <Route path="/admin/*" element={<AdminPanel />} />
            
            {/* Utilities */}
            <Route path="/diagnostics" element={<DiagnosticViewer />} />
            <Route path="/upload/immunology" element={<ImmunologyUpload />} />
            <Route path="/semantic-links/extract" element={<PrivateRoute><SemanticLinkExtractor /></PrivateRoute>} />
            <Route path="/analytics" element={<AnalyticsDashboard />} />
            
            {/* Fallback */}
            <Route path="*" element={<ErrorPage />} />
          </Route>
        </Routes>
      </Router>
    </ThemeProvider>
    </AuthProvider>
    </HelmetProvider>
  );
}


