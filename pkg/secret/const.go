package secret

const SecretPatch = `controlPlane:
  inlinePatch:
    cluster:
      aescbcEncryptionSecret: ${aescbcEncryptionKey}
      aggregatorCA:
        crt: ${k8sAggregatorCert}
        key: ${k8sAggregatorCertKey}
      ca:
        crt: ${clusterCert}
        key: ${clusterCertKey}
      etcd:
        ca:
          crt: ${etcdCert}
          key: ${etcdCertKey}
      secret: ${clusterSecret}
      serviceAccount:
        key: ${k8sServiceAccountKey}
      token: ${clusterToken}
    machine:
      ca:
        crt: ${machineCert}
        key: ${machineCertKey}
      token: ${machineToken}
worker:
  inlinePatch:
    cluster:
      aescbcEncryptionSecret: ${aescbcEncryptionKey}
      ca:
        crt: ${clusterCert}
        key: ${clusterCertKey}
      secret: ${clusterSecret}
      token: ${clusterToken}
    machine:
      ca:
        crt: ${machineCert}
        key: ${machineCertKey}
      token: ${machineToken}
`
