# Copyright Contributors to the Open Cluster Management project


set -e
#set -x

CURR_FOLDER_PATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
KIND_KUBECONFIG="${CURR_FOLDER_PATH}/../kind_kubeconfig.yaml"

export KIND_CLUSTER_NAME=library-go-functional-test
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
    GO111MODULE=off go get github.com/onsi/ginkgo/ginkgo
    GO111MODULE=off go get github.com/onsi/gomega/...
fi
if ! which gocovmerge > /dev/null; then
  echo "Installing gocovmerge..."
  GO111MODULE=off go get -u github.com/wadey/gocovmerge
fi

echo "setting up test tmp folder"
[ -d "$FUNCT_TEST_TMPDIR" ] && rm -r "$FUNCT_TEST_TMPDIR"
mkdir -p "$FUNCT_TEST_TMPDIR"
# mkdir -p "$FUNCT_TEST_TMPDIR/output"
mkdir -p "$FUNCT_TEST_TMPDIR/kind-config"

echo "creating cluster"
kind create cluster --name ${KIND_CLUSTER_NAME}

# setup kubeconfig
kind get kubeconfig --name ${KIND_CLUSTER_NAME} > ${KIND_KUBECONFIG}

# create namespace

echo "install cluster"
# setup cluster
make kind-cluster-setup

make functional-test

echo "delete cluster"
kind delete cluster --name ${KIND_CLUSTER_NAME}
