FROM golang:1.22.4-alpine3.20

ENV PATH="/home/${USERNAME}/bin:${PATH}"

RUN apk add --no-cache \
        bash \
        bash-completion \
        git \
        openssh-client \
        build-base \
        git-bash-completion \
        sudo \
    && mkdir -p /etc/skel \
    && touch /etc/skel/.bashrc \
    && echo "source /usr/share/bash-completion/completions/git" >> /etc/skel/.bashrc

ARG USERNAME
ARG USER_ID
ARG GROUP_ID

RUN addgroup -g ${GROUP_ID} ${USERNAME} \
    && adduser -S -D -u ${USER_ID} -G ${USERNAME} ${USERNAME} \
    && echo "${USERNAME} ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/${USERNAME} \
    && chmod 0440 /etc/sudoers.d/${USERNAME} \
    && cd /home/${USERNAME} \
    && mkdir -p -m0755 src \
    && echo "export PATH=\${HOME}/bin:\${PATH}" >> .bashrc \
    && chown -R ${USER_ID}:${GROUP_ID} src .bashrc \
    && install -m 644 /dev/null .gitconfig

RUN sudo -u ${USERNAME} whoami \
    && go install github.com/ramya-rao-a/go-outline@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@latest \
    && go install golang.org/x/tools/gopls@latest \
    && go install github.com/spf13/cobra-cli@latest \
	&& go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
