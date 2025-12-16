## ================================================================================================
# Utility versions
## ================================================================================================
ARG KUBECTL_VERSION=v1.34.3
ARG TALOSCTL_VERSION=v1.11.6
ARG GOLINT_VERSION=v2.7.2-alpine
ARG GORELEASER_VERSION=v2.13.1
ARG SOPS_VERSION=v3.11.0-alpine
ARG AGE_VERSION=1.2.1
ARG AGE_KEYGEN_VERSION=V1.2.0


FROM registry.k8s.io/kubectl:${KUBECTL_VERSION} AS kubectl
FROM ghcr.io/siderolabs/talosctl:${TALOSCTL_VERSION} AS talosctl
FROM golangci/golangci-lint:${GOLINT_VERSION} AS golangci-lint
FROM goreleaser/goreleaser:${GORELEASER_VERSION} AS goreleaser
FROM ghcr.io/getsops/sops:${SOPS_VERSION} AS sops
FROM ghcr.io/mirceanton/age:${AGE_VERSION} AS age


## ================================================================================================
## Main image
## ================================================================================================
FROM mcr.microsoft.com/devcontainers/go:2.0-bookworm@sha256:1d9c1070bd9b7ad75f47afae4d0bc9e3a39f93490923a284bbbfa7e1513696ca AS workspace

COPY --from=kubectl /bin/kubectl /usr/local/bin/kubectl
COPY --from=talosctl /talosctl /usr/local/bin/talosctl
COPY --from=golangci-lint /usr/bin/golangci-lint /usr/local/bin/golangci-lint
COPY --from=goreleaser /usr/bin/goreleaser /usr/local/bin/goreleaser
COPY --from=sops /usr/local/bin/sops /usr/local/bin/sops
COPY --from=age /usr/local/bin/age /usr/local/bin/age
COPY --from=age /usr/local/bin/age-keygen /usr/local/bin/age-keygen

RUN kubectl completion bash | sudo tee /etc/bash_completion.d/kubectl.bash > /dev/null
RUN talosctl completion bash | sudo tee /etc/bash_completion.d/talosctl.bash > /dev/null
RUN golangci-lint completion bash | sudo tee /etc/bash_completion.d/golangci-lint.bash > /dev/null
RUN goreleaser completion bash | sudo tee /etc/bash_completion.d/goreleaser.bash > /dev/null

USER vscode
ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
