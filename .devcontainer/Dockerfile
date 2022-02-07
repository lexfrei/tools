FROM registry.hub.docker.com/hadolint/hadolint:v2.8.0 AS hadolint
FROM registry.hub.docker.com/rancher/kubectl:v1.23.3@sha256:e7c9f61123e2253cf037f00879cfe253e1d544543ad6026cceb6135c3a4d2e22 AS kubectl

FROM registry.hub.docker.com/library/golang:1.17.6

COPY --from=hadolint /bin/hadolint /bin/hadolint
COPY --from=kubectl /bin/kubectl /bin/kubectl

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=1000

ARG APT_TRANSPORT_HTTPS_VERSION="2.2.4"
ARG BASH_VERSION="5.1-2+b3"
ARG CA_CERTIFICATES_VERSION="20210119"
ARG GIT_VERSION="1:2.30.2-1"
ARG GPG_VERSION="2.2.27-2"
ARG SUDO_VERSION="1.9.5p2-3"

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Setup user
RUN groupadd -g 1000 mygroup && \
    adduser --shell /bin/bash --disabled-password --uid ${USER_UID} --gid 1000 ${USERNAME} && \
    mkdir -p /etc/sudoers.d && \
    echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

# Install useful things
RUN apt-get update && apt-get install --assume-yes --no-install-recommends \
        "apt-transport-https=${APT_TRANSPORT_HTTPS_VERSION}" \
        "bash=${BASH_VERSION}" \
        "ca-certificates=${CA_CERTIFICATES}" \
        "git=${GIT_VERSION}" \
        "gpg=${GPG_VERSION}" \
        "sudo=${GPG_VERSION}" \
    && wget -O- -nv https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

USER $USERNAME
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install packages and Go language server
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh &&\
    go install github.com/cweill/gotests/...@latest && \
    go install github.com/fatih/gomodifytags@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install github.com/godoctor/godoctor@latest && \
    go install github.com/josharian/impl@latest && \
    go install github.com/lukehoban/go-outline@latest && \
    go install github.com/spf13/cobra/cobra@latest && \
    go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest && \
    go install golang.org/x/tools/gopls@latest && \
    go install mvdan.cc/gofumpt@latest
