# This Dockerfile builds the main Vultr reverse proxy.
# It should be built from the root of the 'assist' monorepo.

FROM nginx:1.27.0-alpine

# Remove the default Nginx configuration
RUN rm /etc/nginx/conf.d/default.conf

# Copy the portfolio-first proxy configuration
# (serves portfolio at root /, tools at subpaths, DRB at /drb/)
COPY main-proxy/nginx-vultr-portfolio.conf /etc/nginx/conf.d/default.conf

# Keep legacy config available for rollback
COPY main-proxy/nginx.conf /etc/nginx/conf.d/legacy.conf.bak

EXPOSE 80
EXPOSE 443

CMD ["nginx", "-g", "daemon off;"]
