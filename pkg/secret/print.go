package secret

import (
	"encoding/base64"
	"fmt"
	"sort"

	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
)

func PrintSortedSecrets(input *generate.Input) {
	unsorted := getSecrets(input)

	sorted := make([]string, 0, len(unsorted))

	for k := range unsorted {
		sorted = append(sorted, k)
	}

	sort.Strings(sorted)

	for _, k := range sorted {
		fmt.Printf("%s: %s\n", k, unsorted[k])
	}
}
func getSecrets(input *generate.Input) map[string]string {
	secrets := map[string]string{
		"adminCert": getAdminCert(input, "cert"),
		"adminCertKey": getAdminCert(input, "key"),
		"etcdCert": getEtcdCert(input, "cert"),
		"etcdCertKey": getEtcdCert(input, "key"),
		"k8sServiceAccountKey": getK8sServiceAccountKey(input),
		"k8sAggregatorCert": getK8sAggregatorCert(input, "cert"),
		"k8sAggregatorCertKey": getK8sAggregatorCert(input, "key"),
		"clusterToken": getClusterToken(input),
		"aescbcEncryptionKey": getAescbcEncryptionKey(input),
		"clusterSecret": getClusterSecret(input),
		"machineToken": getMachineToken(input),
		"machineCert": getMachineCert(input, "cert"),
		"machineCertKey": getMachineCert(input, "key"),
		"clusterCert": getClusterCert(input, "cert"),
		"clusterCertKey": getClusterCert(input, "key"),
	}

	return secrets
}

func getAdminCert(input *generate.Input, kind string) string {
	var adminCert string
	switch kind {
	case "cert":
		adminCert = base64.StdEncoding.EncodeToString(input.Certs.Admin.Crt)
	case "key":
		adminCert = base64.StdEncoding.EncodeToString(input.Certs.Admin.Key)
	}
	return adminCert
}

func getMachineCert(input *generate.Input, kind string) string {
	var machineCert string
	switch kind {
	case "cert":
		machineCert = base64.StdEncoding.EncodeToString(input.Certs.OS.Crt)
	case "key":
		machineCert = base64.StdEncoding.EncodeToString(input.Certs.OS.Key)
	}
	return machineCert
}

func getK8sAggregatorCert(input *generate.Input, kind string) string {
	var aggregatorCert string
	switch kind {
	case "cert":
		aggregatorCert = base64.StdEncoding.EncodeToString(input.Certs.K8sAggregator.Crt)
	case "key":
		aggregatorCert = base64.StdEncoding.EncodeToString(input.Certs.K8sAggregator.Key)
	}
	return aggregatorCert
}

func getEtcdCert(input *generate.Input, kind string) string {
	var etcdCert string
	switch kind {
	case "cert":
		etcdCert = base64.StdEncoding.EncodeToString(input.Certs.Etcd.Crt)
	case "key":
		etcdCert = base64.StdEncoding.EncodeToString(input.Certs.Etcd.Key)
	}
	return etcdCert
}

func getClusterCert(input *generate.Input, kind string) string {
	var clusterCert string
	switch kind {
	case "cert":
		clusterCert = base64.StdEncoding.EncodeToString(input.Certs.K8s.Crt)
	case "key":
		clusterCert = base64.StdEncoding.EncodeToString(input.Certs.K8s.Key)
	}
	return clusterCert
}

func getK8sServiceAccountKey(input *generate.Input) string {
	etcdCert := base64.StdEncoding.EncodeToString(input.Certs.K8sServiceAccount.Key)
	return etcdCert
}

func getClusterToken(input *generate.Input) string {
	token := input.Secrets.BootstrapToken
	return token
}

func getAescbcEncryptionKey(input *generate.Input) string {
	key := input.Secrets.AESCBCEncryptionSecret
	return key
}

func getClusterSecret(input *generate.Input) string {
	key := input.ClusterSecret
	return key
}

func getMachineToken(input *generate.Input) string {
	token := input.TrustdInfo.Token
	return token
}
