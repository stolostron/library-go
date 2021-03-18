#!/bin/bash

# Copyright Contributors to the Open Cluster Management project

set -e
#set -x

CURR_FOLDER_PATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
KIND_KUBECONFIG="${CURR_FOLDER_PATH}/../kind_kubeconfig.yaml"

export CLUSTER_NAME=$PROJECT_NAME-functional-test
export KUBECONFIG=${KIND_KUBECONFIG}

export FUNCT_TEST_TMPDIR="${CURR_FOLDER_PATH}/../test/functional/tmp"

if ! which kind > /dev/null; then
    echo "installing kind"
    curl -Lo ./kind https://github.com/kubernetes-sigs/kind/releases/download/v0.7.0/kind-$(uname)-amd64
    chmod +x ./kind
    sudo mv ./kind /usr/local/bin/kind
fi
if ! which ginkgo > /dev/null; then
    echo "Installing ginkgo ..."
    pushd $(mktemp -d)
    GO111MODULE=off go get github.com/onsi/ginkgo/ginkgo
    GO111MODULE=off go get github.com/onsi/gomega/...
    popd
fi
if ! which gocovmerge > /dev/null; then
  echo "Installing gocovmerge..."
  pushd $(mktemp -d)
  GO111MODULE=off go get -u github.com/wadey/gocovmerge
  popd
fi

echo "setting up test tmp folder"
[ -d "$FUNCT_TEST_TMPDIR" ] && rm -r "$FUNCT_TEST_TMPDIR"
mkdir -p "$FUNCT_TEST_TMPDIR"
# mkdir -p "$FUNCT_TEST_TMPDIR/output"
mkdir -p "$FUNCT_TEST_TMPDIR/kind-config"

echo "creating cluster"
kind create cluster --name ${CLUSTER_NAME}

# setup kubeconfig
kind get kubeconfig --name ${CLUSTER_NAME} > ${KIND_KUBECONFIG}

# create namespace

echo "install cluster"
# setup cluster
make kind-cluster-setup

set +e
make functional-test
if [ $? != 0 ]; then
  ERR=$?
  POD_NAME=`kubectl get pods -n open-cluster-management -oname | grep clusterlifecycle-state-metrics`
  echo "$POD_NAME" | xargs -L 1 kubectl logs -n open-cluster-management
  exit $ERR
fi
set -e

echo "delete cluster"
kind delete cluster --name ${CLUSTER_NAME}
