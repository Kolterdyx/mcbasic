# Use the official community-maintained Go Alpine image
FROM golang:1.22-alpine

# Install essential build tools and security basics
RUN apk add --no-cache \
    git \
    build-base \
    postgresql-client \
    curl

ARG USER_UID=1001
ARG USER_GID=1001

# Create a dedicated 'agent' user so the AI never runs as root
RUN addgroup -g $USER_GID agent && \
    adduser -u $USER_UID -D -G agent agent
USER agent

# Set up the GOPATH and install tools as the 'agent' user
ENV GOPATH=/home/agent/go
ENV PATH=$PATH:/home/agent/go/bin

WORKDIR /workspace