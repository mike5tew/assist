import React from 'react';
import { Routes, Route } from 'react-router-dom';
import MainLayout from './components/MainLayout';

// Import your page components
import Home from './components/Home';
import { ContentLoader } from './components/contentLoader';
import ErrorPage from './components/ErrorPage';
import ETPProfilePage from './pages/ETPProfilePage';
import ImmunologyPage from './pages/ImmunologyPage';
import ImmunologyUpload from './components/ImmunologyUpload';
import SemanticQuery from './components/SemanticQuery';
import SemanticLinkExtractor from './components/SemanticLinkExtractor';
import AIChat from './components/AIChat';
import StudyAreaGCSEPage from './components/StudyAreas/StudyAreaGCSE';
import StudyAreaCivilEngineeringPage from './components/StudyAreas/StudyAreaCivilEng';
import AdminPanel from './components/AdminPanel';
import CoachMVPDemo from './components/CoachMVPDemo';
import LAOLanding from './pages/LAOLanding';
import APhysicsRevision from './pages/APhysicsRevision';

const AppRoutes = () => {
  return (
    <Routes>
      {/* Routes with header layout */}
      <Route path="/" element={<MainLayout />}>
        <Route index element={<Home />} />
        <Route path="tools" element={<Home />} />
        <Route path="contentLoader" element={<ContentLoader />} />
        <Route path="error" element={<ErrorPage />} />
        <Route path="semantic-query" element={<SemanticQuery />} />
        <Route path="ai-chat" element={<AIChat />} />
        <Route path="coach-mvp-demo" element={<CoachMVPDemo />} />
        <Route path="etp-profile" element={<ETPProfilePage />} />
        
        {/* LAO Landing Page */}
        <Route path="lao" element={<LAOLanding />} />

        {/* A-Level Physics Revision */}
        <Route path="aphy" element={<APhysicsRevision />} />

        {/* Skills Map - external link handled in nav, but provide route for direct access */}
        <Route path="skillstree" element={<SkillsTreeRedirect />} />

        {/* Study areas */}
        <Route path="immunology" element={<ImmunologyPage />} />
        <Route path="immunology-upload" element={<ImmunologyUpload />} />
        <Route path="gcse/*" element={<StudyAreaGCSEPage />} />
        <Route path="GCSE-upload" element={<SemanticLinkExtractor />} />
        <Route path="civil-engineering/*" element={<StudyAreaCivilEngineeringPage />} />

        {/* Admin */}
        <Route path="admin/*" element={<AdminPanel />} />

        {/* Fallback */}
        <Route path="*" element={<ErrorPage />} />
      </Route>
    </Routes>
  );
};

// Redirect component for Skills Tree (hosted separately)
const SkillsTreeRedirect: React.FC = () => {
  React.useEffect(() => {
    // Skills Tree is served at /skillstree on the same server
    window.location.href = '/skillstree/';
  }, []);
  
  return (
    <div style={{ padding: '2rem', textAlign: 'center' }}>
      <p>Redirecting to Skills Map...</p>
    </div>
  );
};

export default AppRoutes;
