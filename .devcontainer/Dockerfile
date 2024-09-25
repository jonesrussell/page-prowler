# Start from the Go image
FROM mcr.microsoft.com/devcontainers/go:1-1.22-bookworm

# Change the ownership of the Go directories
RUN chown -R vscode:vscode /go

USER vscode

# Install the necessary tools
RUN go install github.com/go-task/task/v3/cmd/task@latest && \
    go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install github.com/spf13/cobra-cli@latest
