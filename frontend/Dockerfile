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

# Set environment variables for build
ENV NODE_OPTIONS="--max-old-space-size=4096"
ENV GENERATE_SOURCEMAP=false
ENV NODE_ENV=production
ENV SKIP_PREFLIGHT_CHECK=true
# This variable is baked into the build. An empty string makes requests relative.
ENV REACT_APP_API_URL=""
ENV REACT_APP_ENVIRONMENT=production

# Build with error handling
RUN npm run build || (echo "Build failed, trying with different settings..." && \
    npm install --legacy-peer-deps --force && \
    npm run build)

# ---- Runtime Stage ----
FROM nginx:alpine

COPY --from=builder /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
