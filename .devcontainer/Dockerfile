ARG GO_VERSION=1.17.6

FROM golang:${GO_VERSION}
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=1000

# Setup user
RUN groupadd -g 1000 mygroup && \
    adduser --shell /bin/bash --disabled-password --uid ${USER_UID} --gid 1000 ${USERNAME} && \
    mkdir -p /etc/sudoers.d && \
    echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

# Install useful things
RUN apt-get update && apt-get install --assume-yes bash git sudo gpg

USER $USERNAME
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install packages and Go language server
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh |\
        sh &&\
    go install \
        golang.org/x/tools/gopls@latest && \
    go install \
        github.com/fatih/gomodifytags@latest && \
    go install \
        github.com/josharian/impl@latest && \
    go install \
        github.com/lukehoban/go-outline@latest && \
    go install \
        github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest && \
    go install \
        github.com/cweill/gotests/...@latest && \
    go install \
        mvdan.cc/gofumpt/gofumports@v0.1.1