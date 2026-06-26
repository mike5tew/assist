import React from 'react';
import { Box, Container, Typography, Link as MuiLink } from '@mui/material';
import { Link } from 'react-router-dom';
import SEO from '../components/SEO';

const LAOTerms: React.FC = () => (
  <Box sx={{ bgcolor: '#fff', minHeight: '100vh' }}>
    <SEO
      title="Terms of Use — Little and Often (LAO)"
      description="Terms of use for the Little and Often GCSE Science revision app."
      path="/lao/terms"
    />

    {/* Header bar */}
    <Box sx={{ bgcolor: '#1a237e', color: 'white', py: 3 }}>
      <Container maxWidth="md">
        <MuiLink component={Link} to="/lao" sx={{ color: 'white', textDecoration: 'none', '&:hover': { textDecoration: 'underline' } }}>
          ← Back to Little and Often
        </MuiLink>
      </Container>
    </Box>

    <Container maxWidth="md" sx={{ py: 6 }}>
      <Typography variant="h3" fontWeight="bold" gutterBottom>Terms of Use</Typography>
      <Typography variant="body2" color="text.secondary" gutterBottom>
        Last updated: 18 April 2026
      </Typography>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 4, mb: 1 }}>1. Acceptance of Terms</Typography>
      <Typography paragraph>
        By downloading, installing, or using Little and Often ("LAO", "the App"), you agree to these Terms of Use. If you do not agree, please do not use the App.
      </Typography>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>2. Description of Service</Typography>
      <Typography paragraph>
        LAO is a GCSE Science revision app that provides flashcards, speed reading, audio summaries, and plain text content for AQA Biology, Chemistry, and Physics. The App is intended as a supplementary study aid and does not replace formal education or teaching.
      </Typography>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>3. In-App Purchases</Typography>
      <Typography paragraph>
        The App offers optional one-time, non-consumable in-app purchases to unlock additional course content. All purchases are processed through the Apple App Store and are subject to Apple's terms and conditions. Prices are displayed in the App before purchase.
      </Typography>
      <Box component="ul" sx={{ pl: 3 }}>
        <li><Typography><strong>No subscriptions:</strong> All purchases are one-time payments with lifetime access.</Typography></li>
        <li><Typography><strong>Refunds:</strong> Refund requests must be made through Apple. See{' '}
          <MuiLink href="https://support.apple.com/en-gb/HT204084" target="_blank" rel="noopener">Apple's refund policy</MuiLink>.
        </Typography></li>
        <li><Typography><strong>Restore purchases:</strong> You can restore previous purchases at any time via the Restore Purchases option in the App.</Typography></li>
      </Box>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>4. Intellectual Property</Typography>
      <Typography paragraph>
        All content in the App — including lesson text, summaries, keyword definitions, and audio — is the intellectual property of ESP Thinking (Michael Stewart) or used under licence. You may not reproduce, distribute, or create derivative works from the App's content without written permission.
      </Typography>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>5. Educational Content Disclaimer</Typography>
      <Typography paragraph>
        The content in LAO is designed to support GCSE Science revision and is aligned with AQA specifications. However:
      </Typography>
      <Box component="ul" sx={{ pl: 3 }}>
        <li><Typography>We do not guarantee exam results or academic outcomes.</Typography></li>
        <li><Typography>Specifications may change. We endeavour to keep content up to date but cannot guarantee it reflects the latest syllabus at all times.</Typography></li>
        <li><Typography>The App is a revision aid, not a substitute for classroom teaching or a textbook.</Typography></li>
      </Box>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>6. Acceptable Use</Typography>
      <Typography paragraph>
        You agree not to:
      </Typography>
      <Box component="ul" sx={{ pl: 3 }}>
        <li><Typography>Reverse-engineer, decompile, or disassemble the App.</Typography></li>
        <li><Typography>Distribute, share, or resell content from the App.</Typography></li>
        <li><Typography>Use the App for any unlawful purpose.</Typography></li>
      </Box>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>7. Limitation of Liability</Typography>
      <Typography paragraph>
        To the fullest extent permitted by law, ESP Thinking shall not be liable for any indirect, incidental, or consequential damages arising from your use of the App. The App is provided "as is" without warranties of any kind.
      </Typography>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>8. Changes to Terms</Typography>
      <Typography paragraph>
        We reserve the right to update these terms at any time. Continued use of the App after changes constitutes acceptance of the revised terms.
      </Typography>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>9. Governing Law</Typography>
      <Typography paragraph>
        These terms are governed by the laws of England and Wales.
      </Typography>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>10. Contact</Typography>
      <Typography paragraph>
        Questions about these terms? Contact us at:<br />
        <strong>Email:</strong> world@espthinking.co.uk<br />
        <strong>Website:</strong>{' '}
        <MuiLink href="https://espthinking.co.uk" target="_blank" rel="noopener">espthinking.co.uk</MuiLink>
      </Typography>

      {/* Footer */}
      <Box sx={{ mt: 6, pt: 3, borderTop: '1px solid #eee', textAlign: 'center' }}>
        <Typography variant="body2" color="text.secondary">
          © 2026 ESP Thinking (Michael Stewart). All rights reserved.
        </Typography>
      </Box>
    </Container>
  </Box>
);

export default LAOTerms;
