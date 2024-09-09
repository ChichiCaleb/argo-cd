#!/bin/bash

# This runs at Codespace creation - not part of pre-build

echo "post-create start"
echo "$(date)    post-create start" >> "$HOME/status"


# Update package list and install necessary packages
sudo apt-get update
sudo apt-get install -y curl make 

# Function to install Go
install_go() {
  local go_version="1.22.2"
  local go_tar="go${go_version}.linux-amd64.tar.gz"
  local go_url="https://golang.org/dl/${go_tar}"

  # Remove any existing Go installation
  sudo rm -rf /usr/local/go

  # Download and install Go
  curl -OL "${go_url}"
  sudo tar -C /usr/local -xzf "${go_tar}"

   # Add Go to PATH for the current user
  echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile

  # Source the updated profile for the current user
  source ~/.profile

  # Clean up downloaded tar file
  rm "${go_tar}"

}

# Run installation functions
install_go

# Install and configure nginx ingress for kind cluster
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

echo "post-create complete"
echo "$(date +'%Y-%m-%d %H:%M:%S')    post-create complete" >> "$HOME/status"
