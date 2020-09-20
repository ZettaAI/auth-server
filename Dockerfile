# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.140.1/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version: 1, 1.15, 1.14
ARG VARIANT="1"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# [Option] Install Node.js
ARG INSTALL_NODE="true"
ARG NODE_VERSION="lts/*"

RUN go get -v golang.org/x/oauth2 \
    && go get -v github.com/labstack/echo \
    && go get -v github.com/go-redis/redis/v8

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1