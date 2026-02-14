# Production Dockerfile (Multi-stage) for the main 'assist' application.

# ---- Build Stage ----
FROM node:20-alpine3.19 AS builder

WORKDIR /app

COPY package*.json ./

# This preserves the complex install logic needed for this specific project.
RUN npm cache clean --force && \
    rm -rf node_modules package-lock.json && \
    npm config set legacy-peer-deps true && \
    npm config set audit-level none && \
    npm install --no-package-lock --force && \
    npm install react-scripts@5.0.1 --force

COPY . .

# Declare build arguments
ARG APP_VERSION=unknown
ARG REACT_APP_RECAPTCHA_SITE_KEY=""

# Set the React environment variables from the build arguments
ENV REACT_APP_VERSION=$APP_VERSION
ENV REACT_APP_RECAPTCHA_SITE_KEY=$REACT_APP_RECAPTCHA_SITE_KEY

ENV GENERATE_SOURCEMAP=false
ENV NODE_ENV=production
# Portfolio-first: serve at root. Override with --build-arg PUBLIC_URL=/esp-organizer for legacy.
ARG PUBLIC_URL=/
ENV PUBLIC_URL=$PUBLIC_URL

RUN npm run build

# ---- Serve Stage ----
FROM nginx:alpine

COPY --from=builder /app/build /usr/share/nginx/html

# Copy the local, self-contained Nginx config.
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
