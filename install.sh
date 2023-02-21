#!/bin/bash

cd /tmp

sudo wget -q https://github.com/surajincloud/kubectl-eks/releases/download/v0.1.2/kubectl-eks_0.1.2_checksums.txt
sudo wget https://github.com/surajincloud/kubectl-eks/releases/download/v0.1.2/kubectl-eks_0.1.2_linux_arm64.tar.gz


sudo tar -xvf kubectl-eks_0.1.2_linux_arm64.tar.gz

sudo mv -v kubectl-eks /usr/local/bin/

echo "kubectl-eks installation completed"