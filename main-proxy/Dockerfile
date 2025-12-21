# This Dockerfile builds the main Vultr reverse proxy.
# It should be built from the root of the 'assist' monorepo.

FROM nginx:1.27.0-alpine

# Remove the default Nginx configuration
RUN rm /etc/nginx/conf.d/default.conf

# Copy our custom main proxy configuration
COPY main-proxy/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
EXPOSE 443

CMD ["nginx", "-g", "daemon off;"]
