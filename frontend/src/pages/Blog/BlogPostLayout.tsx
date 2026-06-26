import React from 'react';
import {
  Box,
  Container,
  Typography,
  Button,
  Stack,
  Chip,
  Divider,
  alpha,
  useTheme,
} from '@mui/material';
import {
  Home as HomeIcon,
  ArrowBack as BackIcon,
  ArrowForward as ForwardIcon,
} from '@mui/icons-material';
import { Link } from 'react-router-dom';
import SEO from '../../components/SEO';

interface BlogPostLayoutProps {
  title: string;
  subtitle: string;
  seriesLabel: string;
  date: string;
  heroColor: string;
  heroGradient: string;
  icon: React.ReactNode;
  tags: string[];
  featuredImage?: string;
  prevPost?: { label: string; path: string };
  nextPost?: { label: string; path: string };
  children: React.ReactNode;
}

const BlogPostLayout: React.FC<BlogPostLayoutProps> = ({
  title,
  subtitle,
  seriesLabel,
  date,
  heroColor,
  heroGradient,
  icon,
  tags,
  featuredImage,
  prevPost,
  nextPost,
  children,
}) => {
  const theme = useTheme();

  return (
    <Box>
      <SEO
        title={`${title} — ETP Series`}
        description={subtitle}
        path={`/blog/${title.toLowerCase().replace(/\s+/g, '-')}`}
        image={featuredImage}
      />

      {/* Hero */}
      <Box
        sx={{
          background: heroGradient,
          borderTop: `6px solid ${heroColor}`,
          color: 'white',
          pt: { xs: 4, md: 8 },
          pb: { xs: 6, md: 10 },
          px: 2,
          position: 'relative',
          overflow: 'hidden',
          '&::after': {
            content: '""',
            position: 'absolute',
            bottom: 0,
            left: 0,
            right: 0,
            height: '60px',
            background: 'linear-gradient(to top right, #f5f7f9 50%, transparent 50%)',
          },
        }}
      >
        <Container maxWidth="md">
          <Stack spacing={2}>
            <Stack direction="row" spacing={2}>
              <Button
                component={Link}
                to="/"
                startIcon={<HomeIcon />}
                sx={{
                  color: 'rgba(255,255,255,0.8)',
                  textTransform: 'none',
                  '&:hover': { color: 'white', bgcolor: 'rgba(255,255,255,0.1)' },
                }}
              >
                Portfolio
              </Button>
              <Button
                component={Link}
                to="/blog"
                startIcon={<BackIcon />}
                sx={{
                  color: 'rgba(255,255,255,0.8)',
                  textTransform: 'none',
                  '&:hover': { color: 'white', bgcolor: 'rgba(255,255,255,0.1)' },
                }}
              >
                All Posts
              </Button>
            </Stack>

            <Chip
              label={seriesLabel}
              sx={{
                bgcolor: 'rgba(255,255,255,0.15)',
                color: 'white',
                alignSelf: 'flex-start',
                fontWeight: 600,
              }}
            />

            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
              <Box sx={{ fontSize: { xs: 40, md: 56 } }}>{icon}</Box>
              <Typography
                variant="h2"
                component="h1"
                fontWeight="bold"
                sx={{ lineHeight: 1.1, fontSize: { xs: '1.8rem', md: '2.8rem' } }}
              >
                {title}
              </Typography>
            </Box>

            <Typography variant="h5" sx={{ opacity: 0.92, maxWidth: 600 }}>
              {subtitle}
            </Typography>

            <Typography variant="body2" sx={{ opacity: 0.7 }}>
              {date} · Michael Stewart · ESP Thinking
            </Typography>
          </Stack>
        </Container>
      </Box>

      {/* Featured Image */}
      {featuredImage && (
        <Box sx={{ bgcolor: '#f5f7f9', pt: { xs: 3, md: 4 }, px: 2 }}>
          <Container maxWidth="md">
            <Box
              component="img"
              src={`${process.env.PUBLIC_URL}${featuredImage}`}
              alt={title}
              sx={{
                width: '100%',
                borderRadius: 3,
                boxShadow: '0 4px 20px rgba(0,0,0,0.12)',
              }}
            />
          </Container>
        </Box>
      )}

      {/* Article Body */}
      <Box sx={{ py: { xs: 4, md: 6 }, bgcolor: '#f5f7f9' }}>
        <Container maxWidth="md">
          <Box
            sx={{
              bgcolor: 'white',
              borderRadius: 3,
              p: { xs: 3, md: 5 },
              boxShadow: '0 2px 12px rgba(0,0,0,0.06)',
              '& h2': {
                color: heroColor,
                fontWeight: 700,
                fontSize: '1.5rem',
                mt: 5,
                mb: 2,
              },
              '& h3': {
                fontWeight: 600,
                fontSize: '1.2rem',
                mt: 4,
                mb: 1.5,
              },
              '& p': {
                lineHeight: 1.8,
                mb: 2,
                color: theme.palette.text.secondary,
              },
              '& blockquote': {
                borderLeft: `4px solid ${heroColor}`,
                pl: 3,
                py: 1,
                my: 3,
                bgcolor: alpha(heroColor, 0.04),
                borderRadius: '0 8px 8px 0',
                '& p': { fontStyle: 'italic', color: theme.palette.text.primary },
              },
            }}
          >
            {children}
          </Box>

          {/* Tags */}
          <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap sx={{ mt: 3 }}>
            {tags.map((tag) => (
              <Chip key={tag} label={tag} size="small" variant="outlined" />
            ))}
          </Stack>

          {/* Previous/Next navigation */}
          <Divider sx={{ my: 4 }} />
          <Stack
            direction={{ xs: 'column', sm: 'row' }}
            justifyContent="space-between"
            spacing={2}
          >
            {prevPost ? (
              <Button
                component={Link}
                to={prevPost.path}
                startIcon={<BackIcon />}
                sx={{ textTransform: 'none' }}
              >
                {prevPost.label}
              </Button>
            ) : (
              <Box />
            )}
            {nextPost ? (
              <Button
                component={Link}
                to={nextPost.path}
                endIcon={<ForwardIcon />}
                sx={{ textTransform: 'none' }}
              >
                {nextPost.label}
              </Button>
            ) : (
              <Box />
            )}
          </Stack>
        </Container>
      </Box>
    </Box>
  );
};

export default BlogPostLayout;
