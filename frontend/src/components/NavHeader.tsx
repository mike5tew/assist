import React, { useState } from 'react';
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  IconButton,
  Menu,
  MenuItem,
  Box,
  Container,
  useTheme,
  useMediaQuery,
  Divider,
} from '@mui/material';
import {
  Menu as MenuIcon,
  Home as HomeIcon,
  Psychology as PsychologyIcon,
  AccountTree as SkillsIcon,
  School as SchoolIcon,
  Dashboard as DashboardIcon,
  ArrowDropDown as DropdownIcon,
  PhoneIphone as MobileIcon,
  Analytics as AnalyticsIcon,
  Article as ArticleIcon,
} from '@mui/icons-material';
import { Link, useLocation } from 'react-router-dom';

interface NavItem {
  label: string;
  path?: string;
  icon?: React.ReactNode;
  children?: NavItem[];
}

const navItems: NavItem[] = [
  { label: 'Home', path: '/', icon: <HomeIcon /> },
  { label: 'Blog', path: '/blog', icon: <ArticleIcon /> },
  {
    label: 'Demos',
    icon: <DashboardIcon />,
    children: [
      { label: 'ETP Profile', path: '/etp-profile', icon: <PsychologyIcon /> },
      { label: 'Skills Map', path: '/skillstree', icon: <SkillsIcon /> },
      { label: 'Coach MVP', path: '/coach-mvp-demo', icon: <SchoolIcon /> },
    ],
  },
  {
    label: 'Tools',
    icon: <DashboardIcon />,
    children: [
      { label: 'MAT Strategic Dashboard', path: 'http://localhost/drb/', icon: <DashboardIcon /> },
      { label: 'LAO Adaptive Revision', path: '/lao', icon: <SchoolIcon /> },
      { label: 'ToddlerOS', path: '/toddler-os', icon: <SchoolIcon /> },
      { label: 'PrimaryOS (Neuron Navigators)', path: '/primary-os', icon: <MobileIcon /> },
      { label: 'CareerOS', path: '/careeros', icon: <DashboardIcon /> },
      { label: 'ESP World', path: '/esp-world', icon: <AnalyticsIcon /> },
    ],
  },
];

const NavHeader: React.FC = () => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const location = useLocation();

  // Mobile menu state
  const [mobileAnchor, setMobileAnchor] = useState<null | HTMLElement>(null);

  // Dropdown menu states
  const [dropdownAnchors, setDropdownAnchors] = useState<{ [key: string]: HTMLElement | null }>({});

  const handleMobileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setMobileAnchor(event.currentTarget);
  };

  const handleMobileMenuClose = () => {
    setMobileAnchor(null);
  };

  const handleDropdownOpen = (label: string, event: React.MouseEvent<HTMLElement>) => {
    setDropdownAnchors({ ...dropdownAnchors, [label]: event.currentTarget });
  };

  const handleDropdownClose = (label: string) => {
    setDropdownAnchors({ ...dropdownAnchors, [label]: null });
  };

  const isActive = (path?: string) => {
    if (!path) return false;
    if (path === '/') return location.pathname === '/';
    return location.pathname.startsWith(path);
  };

  // Desktop navigation
  const renderDesktopNav = () => (
    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
      {navItems.map((item) => {
        if (item.children) {
          return (
            <Box key={item.label}>
              <Button
                color="inherit"
                onClick={(e) => handleDropdownOpen(item.label, e)}
                endIcon={<DropdownIcon />}
                sx={{
                  textTransform: 'none',
                  fontWeight: item.children.some((c) => isActive(c.path)) ? 'bold' : 'normal',
                }}
              >
                {item.label}
              </Button>
              <Menu
                anchorEl={dropdownAnchors[item.label]}
                open={Boolean(dropdownAnchors[item.label])}
                onClose={() => handleDropdownClose(item.label)}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'left' }}
              >
                {item.children.map((child) => {
                  const isExternal = child.path?.startsWith('http');
                  return (
                    <MenuItem
                      key={child.path}
                      component={isExternal ? 'a' : Link}
                      to={isExternal ? undefined : child.path}
                      href={isExternal ? child.path : undefined}
                      target={isExternal ? '_blank' : undefined}
                      rel={isExternal ? 'noopener noreferrer' : undefined}
                      onClick={() => handleDropdownClose(item.label)}
                      selected={isActive(child.path)}
                    >
                      {child.icon && <Box sx={{ mr: 1, display: 'flex' }}>{child.icon}</Box>}
                      {child.label}
                    </MenuItem>
                  );
                })}
              </Menu>
            </Box>
          );
        }

        return (
          <Button
            key={item.path}
            component={Link}
            to={item.path!}
            color="inherit"
            sx={{
              textTransform: 'none',
              fontWeight: isActive(item.path) ? 'bold' : 'normal',
              borderBottom: isActive(item.path) ? '2px solid white' : 'none',
            }}
          >
            {item.label}
          </Button>
        );
      })}
    </Box>
  );

  // Mobile navigation
  const renderMobileNav = () => (
    <>
      <IconButton color="inherit" onClick={handleMobileMenuOpen} edge="end">
        <MenuIcon />
      </IconButton>
      <Menu
        anchorEl={mobileAnchor}
        open={Boolean(mobileAnchor)}
        onClose={handleMobileMenuClose}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
        transformOrigin={{ vertical: 'top', horizontal: 'right' }}
      >
        {navItems.map((item) => {
          if (item.children) {
            return [
              <MenuItem key={item.label} disabled sx={{ opacity: 0.7, fontWeight: 'bold' }}>
                {item.label}
              </MenuItem>,
              ...item.children.map((child) => {
                const isExternal = child.path?.startsWith('http');
                return (
                  <MenuItem
                    key={child.path}
                    component={isExternal ? 'a' : Link}
                    to={isExternal ? undefined : child.path}
                    href={isExternal ? child.path : undefined}
                    target={isExternal ? '_blank' : undefined}
                    rel={isExternal ? 'noopener noreferrer' : undefined}
                    onClick={handleMobileMenuClose}
                    selected={isActive(child.path)}
                    sx={{ pl: 4 }}
                  >
                    {child.label}
                  </MenuItem>
                );
              }),
              <Divider key={`${item.label}-divider`} />,
            ];
          }

          return (
            <MenuItem
              key={item.path}
              component={Link}
              to={item.path!}
              onClick={handleMobileMenuClose}
              selected={isActive(item.path)}
            >
              {item.label}
            </MenuItem>
          );
        })}
      </Menu>
    </>
  );

  return (
    <AppBar position="sticky" elevation={1}>
      <Container maxWidth="lg">
        <Toolbar disableGutters sx={{ justifyContent: 'space-between' }}>
          {/* Logo / Brand */}
          <Typography
            variant="h6"
            component={Link}
            to="/"
            sx={{
              textDecoration: 'none',
              color: 'inherit',
              fontWeight: 'bold',
              display: 'flex',
              alignItems: 'center',
              gap: 1,
            }}
          >
                      <Box component="img" src={`${process.env.PUBLIC_URL}/ESPLogoLong.png`} alt="ESP Thinking Portfolio" sx={{ height: 60, width: 'auto' }} />
        
          </Typography>

          {/* Navigation */}
          {isMobile ? renderMobileNav() : renderDesktopNav()}
        </Toolbar>
      </Container>
    </AppBar>
  );
};

export default NavHeader;
