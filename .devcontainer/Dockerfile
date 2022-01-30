FROM hadolint/hadolint@sha256:1b47f3192d2c69c622863b84d5cccbfa0c1fac9022bfc297aa931e6d44e51608 AS hadolint
FROM rancher/kubectl:v1.22.2@sha256:ee717db42f08236b52e57726b2684a1e560a838c20cab5adb81e214a470400ee AS kubectl

FROM registry.hub.docker.com/library/golang:1.17.6@sha256:301609ebecc0ec4cd3174294220a4d9c92aab9015b3a2958297d7663aac627a1

COPY --from=hadolint /bin/hadolint /bin/hadolint
COPY --from=kubectl /bin/kubectl /bin/kubectl

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
RUN apt-get update && apt-get install --assume-yes bash git sudo gpg apt-transport-https ca-certificates && \
    wget -O- -nv https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

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
        github.com/spf13/cobra/cobra@latest && \
    go install \
        github.com/go-delve/delve/cmd/dlv@latest && \
    go install \
        mvdan.cc/gofumpt@latest && \
    go install \
        github.com/godoctor/godoctor@latest
