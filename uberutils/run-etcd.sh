KUBE_CERTS=${KUBE_CERTS:-$PWD/uberutils/certs}
rm -rf /tmp/etcd && mkdir -p /tmp/etcd
etcd \
    --client-cert-auth=true \
    --data-dir=/tmp/etcd \
    --trusted-ca-file=$KUBE_CERTS/etcd/ca.crt \
    --peer-key-file=$KUBE_CERTS/etcd/peer.key \
    --peer-trusted-ca-file=$KUBE_CERTS/etcd/ca.crt \
    --listen-peer-urls=https://127.0.0.1:2380 \
    --listen-client-urls=https://127.0.0.1:2379 \
    --peer-client-cert-auth=true \
    --cert-file=$KUBE_CERTS/etcd/server.crt \
    --key-file=$KUBE_CERTS/etcd/server.key \
    --peer-cert-file=$KUBE_CERTS/etcd/peer.crt \
    --advertise-client-urls=https://127.0.0.1:2379

    # k8s.gcr.io/etcd-amd64:3.1.12 \
#     --advertise-client-urls=https://127.0.0.1:2379 \
#     --cert-file=/etc/kubernetes/pki/etcd/server.crt \
#     --client-cert-auth=true \
#     --data-dir=/var/lib/etcd \
#     --initial-advertise-peer-urls=https://127.0.0.1:2380 \
#     --initial-cluster=kubeadm=https://127.0.0.1:2380 \
#     --key-file=/etc/kubernetes/pki/etcd/server.key \
#     --listen-client-urls=https://127.0.0.1:2379 \
#     --listen-peer-urls=https://127.0.0.1:2380 \
#     --name=kubeadm \
#     --peer-cert-file=/etc/kubernetes/pki/etcd/peer.crt \
#     --peer-client-cert-auth=true \
#     --peer-key-file=/etc/kubernetes/pki/etcd/peer.key \
#     --peer-trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt \
#     --snapshot-count=10000 \
#     --trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt
