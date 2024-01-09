FROM mcr.microsoft.com/devcontainers/go:1-1.21-bullseye

# Install Qodana
RUN wget https://github.com/JetBrains/qodana-cli/releases/download/v2023.3.0/qodana_linux_x86_64 -O qodana && \
   chmod +x qodana && \
   mv qodana /usr/local/bin/

# Create /data directory and change ownership to vscode user
#RUN mkdir -p /data/results && chown -R vscode:vscode /data && chmod -R 755 /data

# Create the docker group and add the 'vscode' user to it
RUN groupadd docker && usermod -aG docker vscode
