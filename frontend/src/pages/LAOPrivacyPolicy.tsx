import React from 'react';
import { Box, Container, Typography, Link as MuiLink } from '@mui/material';
import { Link } from 'react-router-dom';
import SEO from '../components/SEO';

const LAOPrivacyPolicy: React.FC = () => (
  <Box sx={{ bgcolor: '#fff', minHeight: '100vh' }}>
    <SEO
      title="Privacy Policy — Little and Often (LAO)"
      description="Privacy policy for the Little and Often GCSE Science revision app."
      path="/lao/privacy"
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
      <Typography variant="h3" fontWeight="bold" gutterBottom>Privacy Policy</Typography>
      <Typography variant="body2" color="text.secondary" gutterBottom>
        Last updated: 18 April 2026
      </Typography>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 4, mb: 1 }}>1. Introduction</Typography>
      <Typography paragraph>
        Little and Often ("LAO", "we", "us", or "our") is a GCSE Science revision app published by ESP Thinking (Michael Stewart, sole trader). This policy explains what data we collect, how we use it, and your rights.
      </Typography>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>2. Data We Collect</Typography>
      <Typography paragraph>
        <strong>We do not collect any personal data.</strong> LAO is designed to work entirely on your device:
      </Typography>
      <Box component="ul" sx={{ pl: 3 }}>
        <li><Typography>All lesson content, keywords, and summaries are stored locally in an embedded database on your device.</Typography></li>
        <li><Typography>Self-assessment confidence ratings are stored on-device only and are never transmitted to any server.</Typography></li>
        <li><Typography>We do not require user accounts, email addresses, or login credentials.</Typography></li>
        <li><Typography>We do not use analytics, tracking pixels, or advertising SDKs.</Typography></li>
      </Box>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>3. In-App Purchases</Typography>
      <Typography paragraph>
        LAO offers optional one-time in-app purchases to unlock additional course content. These transactions are processed entirely by Apple through the App Store. We receive confirmation of your purchase status via RevenueCat (our purchase management provider) but do not receive your name, email, or payment details. RevenueCat receives an anonymous app user ID only. See{' '}
        <MuiLink href="https://www.revenuecat.com/privacy" target="_blank" rel="noopener">RevenueCat's Privacy Policy</MuiLink>.
      </Typography>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>4. Third-Party Services</Typography>
      <Box component="ul" sx={{ pl: 3 }}>
        <li><Typography><strong>Apple App Store</strong> — processes payments and delivers the app. Apple's privacy policy applies to their services.</Typography></li>
        <li><Typography><strong>RevenueCat</strong> — manages purchase entitlements using an anonymous device identifier. No personal data is shared.</Typography></li>
        <li><Typography><strong>Expo / EAS</strong> — used to build and distribute the app. No user data is collected via Expo at runtime.</Typography></li>
      </Box>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>5. Children's Privacy</Typography>
      <Typography paragraph>
        LAO is designed for GCSE students (typically aged 14–16). Because we do not collect any personal data, we comply with applicable children's privacy regulations including COPPA and the UK Age Appropriate Design Code. No data about children is collected, stored, or shared.
      </Typography>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>6. Data Storage &amp; Security</Typography>
      <Typography paragraph>
        All user-generated data (confidence ratings, progress) is stored locally on your device using secure on-device storage. We have no server-side database of user information. If you delete the app, all local data is permanently removed.
      </Typography>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>7. Your Rights</Typography>
      <Typography paragraph>
        Under the UK GDPR and Data Protection Act 2018, you have the right to access, correct, or delete your personal data. Since we do not collect personal data, there is nothing to request. If you have questions, contact us at the address below.
      </Typography>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>8. Changes to This Policy</Typography>
      <Typography paragraph>
        We may update this privacy policy from time to time. Any changes will be posted on this page with an updated date.
      </Typography>

      <Typography variant="h5" fontWeight="bold" sx={{ mt: 3, mb: 1 }}>9. Contact</Typography>
      <Typography paragraph>
        If you have any questions about this privacy policy, contact us at:<br />
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

export default LAOPrivacyPolicy;
