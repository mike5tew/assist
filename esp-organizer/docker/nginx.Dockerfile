# /docker/nginx.Dockerfile

# Stage 1: Build the frontend application
# This uses the existing frontend/Dockerfile logic
FROM node:lts-alpine AS builder

WORKDIR /app

# Copy package.json files
COPY frontend/package*.json ./

# Install dependencies with necessary configurations for compatibility
RUN npm cache clean --force && \
    npm config set legacy-peer-deps true && \
    npm config set audit-level none && \
    npm install --no-package-lock --force && \
    npm install react-scripts@5.0.1 --force

# Copy frontend source code
COPY frontend/ ./

# Set environment variables for the build
# IMPORTANT: Remove the leading slash from API URL to prevent double slashes
ENV NODE_OPTIONS="--max-old-space-size=4096"
ENV GENERATE_SOURCEMAP=false
ENV NODE_ENV=production
ENV SKIP_PREFLIGHT_CHECK=true
ENV REACT_APP_API_URL=http://localhost:8080
ENV REACT_APP_API_BASE_PATH=""

# Build the frontend
RUN npm run build

# Stage 2: Setup Nginx with the built frontend
FROM nginx:alpine

# Copy the built static files from the frontend builder
COPY --from=frontend-builder /app/build /usr/share/nginx/html

# Copy our custom Nginx configuration
COPY nginx/local.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:80/ || exit 1

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
