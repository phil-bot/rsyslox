FROM ubuntu:24.04

# Prevent services from starting during apt-get install
RUN echo '#!/bin/sh\nexit 101' > /usr/sbin/policy-rc.d && \
    chmod +x /usr/sbin/policy-rc.d

# Install dependencies
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        rsyslog \
        rsyslog-mysql \
        mariadb-server \
        curl \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Remove policy-rc.d after installation
RUN rm -f /usr/sbin/policy-rc.d

# Create installation directory
RUN mkdir -p /opt/rsyslox

# Copy entrypoint
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose API port
EXPOSE 8000

ENTRYPOINT ["/entrypoint.sh"]
CMD ["tail", "-f", "/dev/null"]
