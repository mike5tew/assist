# ---- Build Stage ----
FROM node:lts-alpine AS builder

WORKDIR /app

# Copy package.json and package-lock.json (or yarn.lock)
COPY package*.json ./

# Force clean install with specific npm configurations to handle dependency conflicts
RUN npm cache clean --force && \
    rm -rf node_modules package-lock.json && \
    npm config set legacy-peer-deps true && \
    npm config set audit-level none && \
    npm install --no-package-lock --force && \
    npm install react-scripts@5.0.1 --force

# Copy the rest of the application source code
COPY . .

# Set environment variables to reduce memory usage and avoid conflicts
ENV NODE_OPTIONS="--max-old-space-size=4096"
ENV GENERATE_SOURCEMAP=false
ENV NODE_ENV=production
ENV SKIP_PREFLIGHT_CHECK=true

# Build with error handling
RUN npm run build || (echo "Build failed, trying with different settings..." && \
    npm install --legacy-peer-deps --force && \
    npm run build)

# The serve stage is removed. The main nginx-local service will now handle serving static files.
# The nginx-local service will handle serving the static files.
