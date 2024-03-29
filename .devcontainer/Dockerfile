FROM mcr.microsoft.com/devcontainers/go:1-1.22-bullseye

# Create /data directory and change ownership to vscode user
#RUN mkdir -p /data/results && chown -R vscode:vscode /data && chmod -R 755 /data

# Create the docker group and add the 'vscode' user to it
RUN groupadd docker && usermod -aG docker vscode
