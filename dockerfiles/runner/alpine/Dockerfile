FROM alpine:3.12.0

ARG TARGETPLATFORM

RUN adduser -D -S -h /home/cicd-runner cicd-runner

# hadolint ignore=DL3018
RUN apk add --no-cache \
    bash \
    ca-certificates \
    git \
    tzdata \
    openssh-client

ARG DOCKER_MACHINE_VERSION
ARG DUMB_INIT_VERSION
ARG GIT_LFS_VERSION

COPY cicd-runner-linux-* /usr/bin/
COPY checksums-* install-deps install-cicd-runner /tmp/

# Install GNU wget for "-nv" flag support
# hadolint ignore=DL3018
RUN apk add --no-cache --virtual .fetch-deps wget && \
    /tmp/install-deps "${TARGETPLATFORM}" "${DOCKER_MACHINE_VERSION}" "${DUMB_INIT_VERSION}" "${GIT_LFS_VERSION}" && \
    apk del .fetch-deps

COPY entrypoint /
RUN chmod +x /entrypoint

STOPSIGNAL SIGQUIT
VOLUME ["/etc/cicd-runner", "/home/cicd-runner"]
ENTRYPOINT ["/usr/bin/dumb-init", "/entrypoint"]
CMD ["run", "--user=cicd-runner", "--working-directory=/home/cicd-runner"]
