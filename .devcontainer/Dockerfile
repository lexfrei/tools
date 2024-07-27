FROM docker.io/rancher/kubectl:v1.30.3 AS kubectl

FROM docker.io/library/golang:1.22.5

COPY --from=kubectl /bin/kubectl /bin/kubectl

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=1000

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Setup user
RUN groupadd -g 1000 mygroup && \
    adduser --shell /bin/bash --disabled-password --uid ${USER_UID} --gid 1000 ${USERNAME} && \
    mkdir -p /etc/sudoers.d && \
    echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

# Install useful things
# hadolint ignore=DL3008
RUN apt-get update && apt-get install --assume-yes --no-install-recommends \
        "apt-transport-https" \
        "bash" \
        "ca-certificates" \
        "git" \
        "gpg" \
        "sudo" \
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
    go install github.com/spf13/cobra-cli@latest && \
    go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest && \
    go install golang.org/x/tools/gopls@latest && \
    go install mvdan.cc/gofumpt@latest && \
    go install github.com/davidrjenni/reftools/cmd/fillstruct@latest
