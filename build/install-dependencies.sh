#!/bin/bash -e

export GO111MODULE=off

# Go tools

if ! which patter > /dev/null; then      echo "Installing patter ..."; go get -u github.com/apg/patter; fi
if ! which gocovmerge > /dev/null; then  echo "Installing gocovmerge..."; go get -u github.com/wadey/gocovmerge; fi
if ! which golangci-lint > /dev/null; then
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.23.6
fi

# Install kubebuilder for unit test
echo "Install Kubebuilder components for test framework usage!"

_OS=$(go env GOOS)
_ARCH=$(go env GOARCH)
KubeBuilderVersion="2.3.1"
# download kubebuilder and extract it to tmp
curl -L https://go.kubebuilder.io/dl/"$KubeBuilderVersion"/"${_OS}"/"${_ARCH}" | tar -xz -C /tmp/

# move to a long-term location and put it on your path
# (you'll need to set the KUBEBUILDER_ASSETS env var if you put it somewhere else)
sudo mv /tmp/kubebuilder_"$KubeBuilderVersion"_"${_OS}"_"${_ARCH}" /usr/local/kubebuilder
export PATH=$PATH:/usr/local/kubebuilder/bin

# Build tools

# Image tools

# Check tools
