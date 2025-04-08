# Use Alpine as the base image
FROM alpine:3.18

# Install required utilities and curl dependencies
RUN apk add --no-cache bash curl tar libstdc++ coreutils

# Install glibc 2.34 (required by Fabric CLI tools)
RUN curl -Lo /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
    curl -Lo glibc.apk https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.34-r0/glibc-2.34-r0.apk && \
    rm -f /etc/nsswitch.conf && \
    apk add --no-cache glibc.apk && \
    rm glibc.apk

# Optional: install extra tools (uncomment if needed)
# RUN apk add --no-cache jq openssl tree

# Copy Fabric binaries into container
# Make sure `bin` directory contains `peer`, `orderer`, `osnadmin`, etc.
COPY bin/* /usr/local/bin/

# Set working directory
WORKDIR /etc/hyperledger/fabric

# Default command
CMD ["sh"]
