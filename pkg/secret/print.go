package secret

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
	"gopkg.in/yaml.v3"
)

func PrintSecretBundle(secret *generate.SecretsBundle) error {
	buf := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)

	err := encoder.Encode(secret)
	if err != nil {
		return err
	}

	fmt.Printf(buf.String())
	return nil
}

func PrintSortedSecrets(secret *generate.SecretsBundle) {
	unsorted := getSecrets(secret)

	sorted := make([]string, 0, len(unsorted))

	for k := range unsorted {
		sorted = append(sorted, k)
	}

	sort.Strings(sorted)

	for _, k := range sorted {
		fmt.Printf("%s%s %s\n", color.BlueString(k), color.BlueString(":"), unsorted[k])
	}
}

func getSecrets(secret *generate.SecretsBundle) map[string]string {
	secrets := map[string]string{
		"etcdCert": getEtcdCert(secret, "cert"),
		"etcdCertKey": getEtcdCert(secret, "key"),
		"k8sServiceAccountKey": getK8sServiceAccountKey(secret),
		"k8sAggregatorCert": getK8sAggregatorCert(secret, "cert"),
		"k8sAggregatorCertKey": getK8sAggregatorCert(secret, "key"),
		"clusterToken": getClusterToken(secret),
		"aescbcEncryptionKey": getAescbcEncryptionKey(secret),
		"clusterSecret": getClusterSecret(secret),
		"machineToken": getMachineToken(secret),
		"machineCert": getMachineCert(secret, "cert"),
		"machineCertKey": getMachineCert(secret, "key"),
		"clusterCert": getClusterCert(secret, "cert"),
		"clusterCertKey": getClusterCert(secret, "key"),
	}

	return secrets
}

func getEtcdCert(secret *generate.SecretsBundle, kind string) string {
	var etcdCert string
	switch kind {
	case "cert":
		etcdCert = base64.StdEncoding.EncodeToString(secret.Certs.Etcd.Crt)
	case "key":
		etcdCert = base64.StdEncoding.EncodeToString(secret.Certs.Etcd.Key)
	}
	return etcdCert
}

func getK8sServiceAccountKey(secret *generate.SecretsBundle) string {
	etcdCert := base64.StdEncoding.EncodeToString(secret.Certs.K8sServiceAccount.Key)
	return etcdCert
}

func getK8sAggregatorCert(secret *generate.SecretsBundle, kind string) string {
	var aggregatorCert string
	switch kind {
	case "cert":
		aggregatorCert = base64.StdEncoding.EncodeToString(secret.Certs.K8sAggregator.Crt)
	case "key":
		aggregatorCert = base64.StdEncoding.EncodeToString(secret.Certs.K8sAggregator.Key)
	}
	return aggregatorCert
}

func getClusterToken(secret *generate.SecretsBundle) string {
	token := secret.Secrets.BootstrapToken
	return token
}

func getAescbcEncryptionKey(secret *generate.SecretsBundle) string {
	key := secret.Secrets.AESCBCEncryptionSecret
	return key
}

func getClusterSecret(secret *generate.SecretsBundle) string {
	key := secret.Cluster.Secret
	return key
}

func getMachineToken(secret *generate.SecretsBundle) string {
	token := secret.TrustdInfo.Token
	return token
}

func getMachineCert(secret *generate.SecretsBundle, kind string) string {
	var machineCert string
	switch kind {
	case "cert":
		machineCert = base64.StdEncoding.EncodeToString(secret.Certs.OS.Crt)
	case "key":
		machineCert = base64.StdEncoding.EncodeToString(secret.Certs.OS.Key)
	}
	return machineCert
}

func getClusterCert(secret *generate.SecretsBundle, kind string) string {
	var clusterCert string
	switch kind {
	case "cert":
		clusterCert = base64.StdEncoding.EncodeToString(secret.Certs.K8s.Crt)
	case "key":
		clusterCert = base64.StdEncoding.EncodeToString(secret.Certs.K8s.Key)
	}
	return clusterCert
}
