## ================================================================================================
# Utility versions
## ================================================================================================
ARG KUBECTL_VERSION=v1.32.3
ARG TALOSCTL_VERSION=v1.9.5
ARG GOLINT_VERSION=v2.0.2-alpine
ARG GORELEASER_VERSION=v2.8.2
ARG SOPS_VERSION=v3.10.1-alpine
ARG AGE_VERSION=v1.2.0
ARG AGE_KEYGEN_VERSION=V1.2.0


FROM registry.k8s.io/kubectl:${KUBECTL_VERSION} AS kubectl
FROM ghcr.io/siderolabs/talosctl:${TALOSCTL_VERSION} AS talosctl
FROM golangci/golangci-lint:${GOLINT_VERSION} AS golangci-lint
FROM goreleaser/goreleaser:${GORELEASER_VERSION} AS goreleaser
FROM ghcr.io/getsops/sops:${SOPS_VERSION} AS sops
FROM ghcr.io/mirceanton/age:${AGE_VERSION} AS age
FROM ghcr.io/mirceanton/age-keygen:${AGE_KEYGEN_VERSION} AS age-keygen


## ================================================================================================
## Main image
## ================================================================================================
FROM mcr.microsoft.com/devcontainers/go:1.24-bookworm@sha256:fb6bc4a5ae6b6c51b711a244c0031873bf9cf37ff80d09009ce301b663df1700 AS workspace

COPY --from=kubectl /bin/kubectl /usr/local/bin/kubectl
COPY --from=talosctl /talosctl /usr/local/bin/talosctl
COPY --from=golangci-lint /usr/bin/golangci-lint /usr/local/bin/golangci-lint
COPY --from=goreleaser /usr/bin/goreleaser /usr/local/bin/goreleaser
COPY --from=sops /usr/local/bin/sops /usr/local/bin/sops
COPY --from=age /age /usr/local/bin/age
COPY --from=age-keygen /age-keygen /usr/local/bin/age-keygen

RUN kubectl completion bash | sudo tee /etc/bash_completion.d/kubectl.bash > /dev/null
RUN talosctl completion bash | sudo tee /etc/bash_completion.d/talosctl.bash > /dev/null
RUN golangci-lint completion bash | sudo tee /etc/bash_completion.d/golangci-lint.bash > /dev/null
RUN goreleaser completion bash | sudo tee /etc/bash_completion.d/goreleaser.bash > /dev/null

USER vscode
ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
