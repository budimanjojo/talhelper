package secret

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/siderolabs/talos/pkg/machinery/config/generate/secrets"
	"gopkg.in/yaml.v3"
)

// PrintSecretBundle prints the generated `SecretsBundle` into the terminal.
// It returns an error, if any.
func PrintSecretBundle(secret *secrets.Bundle) error {
	buf := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)

	err := encoder.Encode(secret)
	if err != nil {
		return err
	}

	fmt.Print(buf.String())
	return nil
}

// PrintSortedSecrets takes a `SecretsBundle`, sorts them and prints them out.
func PrintSortedSecrets(secret *secrets.Bundle) {
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

// getSecrets takes a `SecretsBundle` and returns `map[string]string` of them.
func getSecrets(secret *secrets.Bundle) map[string]string {
	secrets := map[string]string{
		"etcdCert":             getEtcdCert(secret, "cert"),
		"etcdCertKey":          getEtcdCert(secret, "key"),
		"k8sServiceAccountKey": getK8sServiceAccountKey(secret),
		"k8sAggregatorCert":    getK8sAggregatorCert(secret, "cert"),
		"k8sAggregatorCertKey": getK8sAggregatorCert(secret, "key"),
		"clusterToken":         getClusterToken(secret),
		"aescbcEncryptionKey":  getAescbcEncryptionKey(secret),
		"clusterSecret":        getClusterSecret(secret),
		"machineToken":         getMachineToken(secret),
		"machineCert":          getMachineCert(secret, "cert"),
		"machineCertKey":       getMachineCert(secret, "key"),
		"clusterCert":          getClusterCert(secret, "cert"),
		"clusterCertKey":       getClusterCert(secret, "key"),
	}

	return secrets
}

// getEtcdCert takes a `SecretsBundle` and returns value of the specified
// etcd `key` or `cert` key.
func getEtcdCert(secret *secrets.Bundle, kind string) string {
	var etcdCert string
	switch kind {
	case "cert":
		etcdCert = base64.StdEncoding.EncodeToString(secret.Certs.Etcd.Crt)
	case "key":
		etcdCert = base64.StdEncoding.EncodeToString(secret.Certs.Etcd.Key)
	}
	return etcdCert
}

// getK8sServiceAccountKey takes a `SecretsBundle` and returns value of the
// service account key.
func getK8sServiceAccountKey(secret *secrets.Bundle) string {
	svcAccountKey := base64.StdEncoding.EncodeToString(secret.Certs.K8sServiceAccount.Key)
	return svcAccountKey
}

// getK8sAggregatorCert takes a `SecretsBundle` and returns value of the specified
// k8s aggregator `key` or `cert` key.
func getK8sAggregatorCert(secret *secrets.Bundle, kind string) string {
	var aggregatorCert string
	switch kind {
	case "cert":
		aggregatorCert = base64.StdEncoding.EncodeToString(secret.Certs.K8sAggregator.Crt)
	case "key":
		aggregatorCert = base64.StdEncoding.EncodeToString(secret.Certs.K8sAggregator.Key)
	}
	return aggregatorCert
}

// getClusterToken takes a `SecretsBundle` and returns value of the cluster token.
func getClusterToken(secret *secrets.Bundle) string {
	token := secret.Secrets.BootstrapToken
	return token
}

// getAescbcEncryptionKey takes a `SecretsBundle` and returns value of the Aescbc encryption key.
func getAescbcEncryptionKey(secret *secrets.Bundle) string {
	key := secret.Secrets.AESCBCEncryptionSecret
	return key
}

// getClusterSecret takes a `SecretsBundle` and returns value of the cluster secret key.
func getClusterSecret(secret *secrets.Bundle) string {
	key := secret.Cluster.Secret
	return key
}

// getMachineToken takes a `SecretsBundle` and returns value of the machine token key.
func getMachineToken(secret *secrets.Bundle) string {
	token := secret.TrustdInfo.Token
	return token
}

// getMachineCert takes a `SecretsBundle` and returns value of the specified
// machine `key` or `cert` key.
func getMachineCert(secret *secrets.Bundle, kind string) string {
	var machineCert string
	switch kind {
	case "cert":
		machineCert = base64.StdEncoding.EncodeToString(secret.Certs.OS.Crt)
	case "key":
		machineCert = base64.StdEncoding.EncodeToString(secret.Certs.OS.Key)
	}
	return machineCert
}

// getClusterCert takes a `SecretsBundle` and returns value of the specified
// cluster `key` or `cert` key.
func getClusterCert(secret *secrets.Bundle, kind string) string {
	var clusterCert string
	switch kind {
	case "cert":
		clusterCert = base64.StdEncoding.EncodeToString(secret.Certs.K8s.Crt)
	case "key":
		clusterCert = base64.StdEncoding.EncodeToString(secret.Certs.K8s.Key)
	}
	return clusterCert
}
