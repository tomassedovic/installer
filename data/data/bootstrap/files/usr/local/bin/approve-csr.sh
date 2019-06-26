#!/usr/bin/env bash

KUBECONFIG="${1}"

if [ -f /usr/local/bin/haproxy-watcher.sh ]; then
    echo "NOTE: We're on the service VM. Skipping approve-csr"
    exit 0
fi

echo "Approving all CSR requests until bootstrapping is complete..."
while [ ! -f /opt/openshift/.bootkube.done ]
do
    oc --config="$KUBECONFIG" get csr --no-headers | grep Pending | \
        awk '{print $1}' | \
        xargs --no-run-if-empty oc --config="$KUBECONFIG" adm certificate approve
	sleep 20
done
