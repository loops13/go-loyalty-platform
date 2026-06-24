const target = process.env.API_PROXY_TARGET || 'http://backend:8080';

module.exports = {
  '/api': {
    target,
    secure: false,
    changeOrigin: true,
    pathRewrite: {
      '^/api': '',
    },
  },
};
