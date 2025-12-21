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

# Declare a build argument for the version
ARG APP_VERSION=unknown
# Set the React environment variable from the build argument
ENV REACT_APP_VERSION=$APP_VERSION

ENV GENERATE_SOURCEMAP=false
ENV NODE_ENV=production
# Set the public URL so the app knows its assets are served from /esp-organizer/
ENV PUBLIC_URL=/esp-organizer

RUN npm run build

# ---- Serve Stage ----
FROM nginx:alpine

COPY --from=builder /app/build /usr/share/nginx/html

# Copy the local, self-contained Nginx config.
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
