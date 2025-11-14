import React from 'react';
import { Routes, Route } from 'react-router-dom';
import MiniDrawer from './components/Drawer';


// Import your page components here
import Home from './components/Home';
import { ContentLoader } from './components/contentLoader';
import ErrorPage from './components/ErrorPage';
import ImmunologyPage from './pages/ImmunologyPage'; // Corrected import
import ImmunologyUpload from './components/ImmunologyUpload';
import SemanticQuery from './components/SemanticQuery';
import AIChat from './components/AIChat';
import StudyAreaSelector from './components/StudyAreas/StudyAreaSelector';
import StudyAreaGCSEPage from './components/StudyAreas/StudyAreaGCSE';
import StudyAreaCivilEngineeringPage from './components/StudyAreas/StudyAreaCivilEng';
import AdminPanel from './components/AdminPanel';

// This component defines a placeholder for routes that aren't implemented yet
const PlaceholderPage: React.FC<{ name: string }> = ({ name }) => (
  <div style={{ padding: '20px' }}>
    <h2>{name} Page</h2>
    <p>This page is not implemented yet.</p>
  </div>
);

const AppRoutes = () => {
  return (
    <Routes>
      {/* Routes with the drawer layout */}
      <Route path="/" element={<MiniDrawer />}>
        {/* Default route */}
        <Route
          index
          element={<Home />} // Corrected: Render Home component for the index route
        />
        
        <Route path="home" element={<Home />} />
        <Route path="contentLoader" element={<ContentLoader />} />
        <Route path="error" element={<ErrorPage />} />
        <Route path="semantic-query" element={<SemanticQuery />} />
        <Route path="ai-chat" element={<AIChat />} />

        {/* Study area routes */}
        <Route path="immunology" element={<ImmunologyPage />} />
        <Route path="immunology-upload" element={<ImmunologyUpload />} />

        <Route path="linking/:linkParams" element={<div>Linking page placeholder</div>} />
        
          {/* Main landing page */}
          <Route path="/" element={<StudyAreaSelector />} />
          <Route path="/home" element={<Home />} />
          
          {/* Study area specific routes */}
          <Route path="/gcse/*" element={<StudyAreaGCSEPage />} />
          <Route path="/civil-engineering/*" element={<StudyAreaCivilEngineeringPage />} />
          <Route path="/immunology/*" element={<ImmunologyPage />} />
          
          {/* Admin routes */}
          <Route path="/admin/*" element={<AdminPanel />} />
          
          {/* Legacy routes for backward compatibility */}
          <Route path="/esp-organizer" element={<Home />} />
          <Route path="/esp-organizer/*" element={<Home />} />
          
          {/* Fallback */}
          <Route path="*" element={<StudyAreaSelector />} />
        
        {/* Catch-all for undefined routes */}
        <Route path="*" element={
          <div style={{ padding: '20px' }}>
            <h2>404 - Page Not Found</h2>
            <p>The page you're looking for doesn't exist.</p>
          </div>
        } />
      </Route>
    </Routes>
  );
};

// Ensure no premature usage of initGraphAuth here

export default AppRoutes;
