FROM mcr.microsoft.com/devcontainers/go:1-1.22-bullseye

# Install godoc
RUN go install golang.org/x/tools/cmd/godoc

# Optionally, start godoc server for testing
# RUN godoc -http=:6060

# Create /data directory and change ownership to vscode user
#RUN mkdir -p /data/results && chown -R vscode:vscode /data && chmod -R 755 /data

# Create the docker group and add the 'vscode' user to it
RUN groupadd docker && usermod -aG docker vscode
