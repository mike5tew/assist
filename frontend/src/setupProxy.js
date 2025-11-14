const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  console.log('🚀 Setting up proxy to http://localhost:8080');
  
  const proxyOptions = {
    target: 'http://localhost:8080',
    changeOrigin: true,
    logLevel: 'debug',
    secure: false,
    pathRewrite: {
      '^/api': '/api' // Ensure path rewriting is correct
    },
    onProxyReq: (proxyReq, req, res) => {
      console.log(`🔀 Proxying ${req.method} ${req.originalUrl} -> http://localhost:8080${req.path}`);
    },
    onProxyRes: (proxyRes, req, res) => {
      console.log(`✅ Received ${proxyRes.statusCode} for ${req.method} ${req.originalUrl}`);
    },
    onError: (err, req, res) => {
      console.error('❌ Proxy error:', err.message);
    }
  };

  app.use('/api', createProxyMiddleware(proxyOptions));
};