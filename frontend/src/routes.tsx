import React from 'react';
import { Routes, Route } from 'react-router-dom';
import MiniDrawer from './components/Drawer';

// Import your page components
import Home from './components/Home';
import { ContentLoader } from './components/contentLoader';
import ErrorPage from './components/ErrorPage';
import ImmunologyPage from './pages/ImmunologyPage';
import ImmunologyUpload from './components/ImmunologyUpload';
import SemanticQuery from './components/SemanticQuery';
import SemanticLinkExtractor from './components/SemanticLinkExtractor';
import AIChat from './components/AIChat';
import StudyAreaGCSEPage from './components/StudyAreas/StudyAreaGCSE';
import StudyAreaCivilEngineeringPage from './components/StudyAreas/StudyAreaCivilEng';
import AdminPanel from './components/AdminPanel';
import CoachMVPDemo from './components/CoachMVPDemo';

const AppRoutes = () => {
  return (
    <Routes>
      {/* Routes with drawer layout */}
      <Route path="/" element={<MiniDrawer />}>
        <Route index element={<Home />} />
        <Route path="home" element={<Home />} />
        <Route path="contentLoader" element={<ContentLoader />} />
        <Route path="error" element={<ErrorPage />} />
        <Route path="semantic-query" element={<SemanticQuery />} />
        <Route path="ai-chat" element={<AIChat />} />
        <Route path="coach-mvp-demo" element={<CoachMVPDemo />} />

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

export default AppRoutes;
