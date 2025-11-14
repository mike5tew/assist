// Error page for the site. This will be displayed when a user tries to access a page that does not exist.

import React from 'react';
import { Typography, Box } from '@mui/material';

export default function ErrorPage() {
    return (
        <Box>
            <Typography variant="h1">Error 404</Typography>
            <Typography variant="h2">Page Not Found</Typography>
        </Box>
    );
}