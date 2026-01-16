import React from 'react';
import { styled, useTheme, Theme, CSSObject } from '@mui/material/styles';
import Box from '@mui/material/Box';
import MuiDrawer from '@mui/material/Drawer';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import List from '@mui/material/List';
import CssBaseline from '@mui/material/CssBaseline';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import UploadIcon from '@mui/icons-material/Upload';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton'; // Import ListItemButton
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import HomeIcon from '@mui/icons-material/Home';
import { Outlet, useNavigate } from 'react-router-dom';
import Grid from '@mui/material/Grid';



const MiniDrawer = () => {
  const drawerWidth = 240;


  const openedMixin = (theme: Theme): CSSObject => ({
    width: drawerWidth,
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
    overflowX: 'hidden',
  });

  const closedMixin = (theme: Theme): CSSObject => ({
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    overflowX: 'hidden',
    width: `calc(${theme.spacing(7)} + 1px)`,
    [theme.breakpoints.down('sm')]: {
      width: `calc(${theme.spacing(6)} + 1px)`,
    },
  });

  const DrawerHeader = styled('div')(({ theme }) => ({
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
    padding: theme.spacing(0, 1),
    ...theme.mixins.toolbar,
  }));

  const Drawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== 'open' })(
    ({ theme, open }) => ({
      width: drawerWidth,
      flexShrink: 0,
      whiteSpace: 'nowrap',
      boxSizing: 'border-box',
      ...(open && {
        ...openedMixin(theme),
        '& .MuiDrawer-paper': openedMixin(theme),
      }),
      ...(!open && {
        ...closedMixin(theme),
        '& .MuiDrawer-paper': closedMixin(theme),
      }),
    }),
  );

  interface AppBarProps extends MuiAppBarProps {
    open?: boolean;
  }

  const AppBar = styled(MuiAppBar, {
    shouldForwardProp: (prop) => prop !== 'open',
  })<AppBarProps>(({ theme, open }) => ({
    zIndex: theme.zIndex.drawer + 1,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    ...(open && {
      marginLeft: drawerWidth,
      width: `calc(100% - ${drawerWidth}px)`,
      transition: theme.transitions.create(['width', 'margin'], {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.enteringScreen,
      }),
    }),
  }));

  const [menuName, setMenuName] = React.useState([
    { EntName: 'Home', type: 'page', icon: <HomeIcon />, link: 'home' },
    { EntName: 'ContentLoader', type: 'page', icon: <UploadIcon />, link: 'contentLoader' },
    { EntName: 'Semantic Query', type: 'page', icon: <UploadIcon />, link: 'semantic-query' },
    { EntName: 'AI Chat', type: 'page', icon: <UploadIcon />, link: 'ai-chat' },
    { EntName: 'Coach MVP Demo', type: 'page', icon: <UploadIcon />, link: 'coach-mvp-demo' },
    { EntName: 'Immunology', type: 'page', icon: <UploadIcon />, link: 'immunology' },
    { EntName: 'Immunology Upload', type: 'page', icon: <UploadIcon />, link: 'immunology-upload' },
    { EntName: 'GCSE Study Area', type: 'page', icon: <UploadIcon />, link: 'gcse' },
    { EntName: 'Civil Engineering Study Area', type: 'page', icon: <UploadIcon />, link: 'civil-engineering' },
    { EntName: 'Admin Panel', type: 'page', icon: <UploadIcon />, link: 'admin' },
    { EntName: 'Back', type: 'navigation', icon: <ChevronLeftIcon />, link: 'back' },
   
  ]);

  const theme = useTheme();
  const [open, setOpen] = React.useState(false);
  const handleDrawerOpen = () => {
    setOpen(true);
  };
  //set the background image to an empty file object
  const handleDrawerClose = () => {
    setOpen(false);
  };
  // A menuitem is a single entry in the menu, e.g. Home, Planning, Students, etc.
  interface menuItem { EntName: string, type: string, icon: any, link: string }
  
  // This is the function that will be called when a menu item is clicked.
  const navigate = useNavigate();

  // Using the name that is passed it will either set the menuName state or call the page using the router.
  const clickListener = (mItem: menuItem) => {
    if (mItem.type === 'navigation' && mItem.link === 'back') {
      navigate(-1);
    } else if (mItem.type === 'menu') {
      for (let i = 0; i < menuName.length; i++) {
        if (menuName[i].link === mItem.link) {
          break;
        }
      }
    } else {
      navigate(mItem.link);
    }
  };


  return (
    <Box sx={{
      display: 'flex',
      backgroundImage: `url(${localStorage.getItem('BackgroundImage')})`,
      backgroundSize: 'cover',
      backgroundPosition: 'center',
      backgroundAttachment: 'fixed',
      backgroundRepeat: 'no-repeat',
      minHeight: '100vh',
      width: '100%',
      maxWidth: '100vw',
      overflow: 'hidden',
      backgroundColor: localStorage.getItem('BackgroundColor') || 'inherit',
      // Add z-index to ensure proper stacking
      zIndex: 1,
      position: 'relative',
      // Prevent inheritance issues
      '& *': {
        boxSizing: 'border-box'
      }
    }}>
      <CssBaseline />
      <AppBar 
        position="fixed" 
        open={open}
        sx={{
          backgroundColor: localStorage.getItem('AppBarColor') || 'primary',
          width: '100%',
          [theme.breakpoints.down('sm')]: {
            width: '100%',
            marginLeft: 0,
          }
        }}
      >
        <Toolbar sx={{ minHeight: { xs: '48px', sm: '64px' } }}>
          <Grid container spacing={0} alignItems="center">
            <Grid>
              <IconButton
                color="inherit"
                aria-label="open drawer"
                onClick={handleDrawerOpen}
                edge="start"
                sx={{
                  marginRight: { xs: 0.5, sm: 2 },
                  ...(open && { display: 'none' }),
                }}
              >
                <MenuIcon />
              </IconButton>
            </Grid>
            <Grid item xs>
              <Typography 
                variant="h6" 
                noWrap 
                component="div"
                sx={{
                  fontSize: { 
                    xs: '0.8rem', 
                    sm: '1rem', 
                    md: '1.25rem' 
                  },
                  textAlign: { xs: 'left', sm: 'left' },
                  paddingLeft: { xs: 1, sm: 2 },
                  whiteSpace: { xs: 'normal', sm: 'nowrap' },
                  lineHeight: { xs: 1.2, sm: 'normal' }
                }}
              >
                ESP Organizer
              </Typography>
            </Grid>
          </Grid>
        </Toolbar>
      </AppBar>
      <Drawer
        variant="permanent"
        open={open}
        onKeyDown={handleDrawerClose}
        onClick={handleDrawerClose}
        sx={{
          '& .MuiDrawer-paper': {
            backgroundColor: localStorage.getItem('MenuColour') || 'inherit',
            color: localStorage.getItem('MenuTextColour') || 'inherit',
            [theme.breakpoints.down('sm')]: {
              top: '48px',
              position: 'fixed',
              height: 'calc(100% - 48px)',
              width: open ? '200px' : '50px'
            },
            '& .MuiListItemText-root': {
              whiteSpace: 'normal',
              '& span': {
                fontSize: { xs: '0.8rem', sm: '0.9rem' }
              }
            },
            '& .MuiListItemIcon-root': {
              color: 'inherit',
              minWidth: { xs: '35px', sm: '40px' }
            },
            '& .MuiDivider-root': {
              borderColor: 'inherit'
            }
          },
        }}
      >
        <DrawerHeader>
          <IconButton 
            onClick={handleDrawerClose}
            aria-label="Close menu"
            tabIndex={0}
          >
            {theme.direction === 'rtl' ? <ChevronRightIcon /> : <ChevronLeftIcon />}
          </IconButton>
        </DrawerHeader>
        <Divider />
        <List>
          {/* Add the items from the array that has been selected */}
          {menuName ? menuName.map((text) => (
            // pass the menuitem object to the click listener
            <ListItemButton // Changed from ListItem to ListItemButton
              key={text.EntName} 
              onClick={() => clickListener(text)}
              // The 'button' prop is removed as it's inherent to ListItemButton
              tabIndex={0}
              aria-label={text.EntName}
              sx={{
                '&:focus': {
                  outline: '2px solid',
                  outlineOffset: '-2px'
                },
                '&:hover': {
                  backgroundColor: 'rgba(255, 255, 255, 0.08)'
                }
              }}
            >
              <ListItemIcon aria-hidden="true">
                {text.icon}
              </ListItemIcon>
              <ListItemText 
                primary={text.EntName} 
                sx={{ 
                  opacity: open ? 1 : 0,
                  visibility: open ? 'visible' : 'hidden' 
                }} 
              />
            </ListItemButton> // Changed from ListItem to ListItemButton
          )) : null}
        </List>
      </Drawer>
      <Box 
        component="main" // Changed from span to main for semantic HTML
        sx={{ 
          p: { xs: 1, sm: 2, md: 3 }, 
          pt: { xs: 7, sm: 8, md: 9 },
          flexGrow: 1,
          width: '100%',
          maxWidth: '100%',
          overflowX: 'hidden',
          backgroundColor: localStorage.getItem('ContentBackgroundColor') || 'transparent',
          // Prevent content from being hidden under drawer
          position: 'relative',
          zIndex: 0,
          [theme.breakpoints.down('sm')]: {
            marginLeft: open ? '50px' : '0',
          }
        }}
      >
        <Outlet />
      </Box>
    </Box>
  );
};


export default MiniDrawer;

