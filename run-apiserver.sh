KUBE_ADMISSION_CONTROL='' kube-apiserver \
    --authorization-mode=Node,RBAC \
    --advertise-address=127.0.0.1 \
    --allow-privileged=true \
    --enable-admission-plugins=NodeRestriction \
    --enable-bootstrap-token-auth=true \
    --etcd-servers=https://127.0.0.1:2378 \
    --insecure-port=0 \
    --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname \
    --requestheader-allowed-names=front-proxy-client \
    --client-ca-file=/Users/pourchet/p2k/milestone0/certs/ca.crt \
    --etcd-cafile=/Users/pourchet/p2k/milestone0/certs/etcd/ca.crt \
    --etcd-certfile=/Users/pourchet/p2k/milestone0/certs/apiserver-etcd-client.crt \
    --etcd-keyfile=/Users/pourchet/p2k/milestone0/certs/apiserver-etcd-client.key \
    --kubelet-client-certificate=/Users/pourchet/p2k/milestone0/certs/apiserver-kubelet-client.crt \
    --kubelet-client-key=/Users/pourchet/p2k/milestone0/certs/apiserver-kubelet-client.key \
    --proxy-client-cert-file=/Users/pourchet/p2k/milestone0/certs/front-proxy-client.crt \
    --proxy-client-key-file=/Users/pourchet/p2k/milestone0/certs/front-proxy-client.key \
    --requestheader-client-ca-file=/Users/pourchet/p2k/milestone0/certs/front-proxy-ca.crt \
    --service-account-key-file=/Users/pourchet/p2k/milestone0/certs/sa.pub \
    --tls-cert-file=/Users/pourchet/p2k/milestone0/certs/apiserver.crt \
    --tls-private-key-file=/Users/pourchet/p2k/milestone0/certs/apiserver.key \
    --requestheader-extra-headers-prefix=X-Remote-Extra- \
    --requestheader-group-headers=X-Remote-Group \
    --requestheader-username-headers=X-Remote-User \
    --secure-port=8443 \
    --service-cluster-ip-range=10.96.0.0/12
