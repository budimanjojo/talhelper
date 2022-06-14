package decrypt

import (
	"testing"

	"sigs.k8s.io/yaml"
)

func TestIsSopsEncrypted(t *testing.T) {
	var m1 *sopsFile
	var m2 *sopsFile
	var data1 = `
this:
  is:
    not: encrypted`

	var data2 = `
this:
  is:
    super: encrypted
sops:
  kms: []
  gcp_kms: []
  azure_kv: []
  hc_vault: []
  age:
    - recipient: age123123123132133123454564646465
      enc: |
        -----BEGIN AGE ENCRYPTED FILE-----
        YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSBzaGFRdkVKa1l6Q1BnazZ2
        Y0FFV1RncExUL2VhUHBuVGVyUm5HdG8yMHlNCm5qR0JRMXBlNDdXa1pOTjR6bjN2
        TnlBQkVkMVNBcFZ5THJFTSsrdGNhdmMKLS0tIHYzVGd6cjgvSGhFdk5uMTg3UXVX
        R0NPUVF3M2oyMzR4NHNMRVcvQ1BGRjAKNXHWcLGCqb6XnMUZY57OTyAq5kaYW0wM
        fc97zk6rh/TanFRgBo296fSYtMiieNMqFUv/IayHdXJ8yFv/sDJjeA==
        -----END AGE ENCRYPTED FILE-----
  lastmodified: "2022-05-29T14:59:06Z"
  mac: ENC[AES256_GCM,data:ewqN6amkTExth5IZiQK+ReBd7OTyLZgXYkhVlGGud/Lm1xTyA+d9q2DfD3QpEDQ2AXdrrITxY9FlEjjUg1bNxvbTYEXB0t8WXcdZXskgk6yvlTn2mAmHlfHzNP9rT0mKRdsr7fny6Fkk7NmfeBzgGWUlZXl36+jHwUtrF0M5nKY=,iv:uddW4p5MBRl0KMIz6BdOBi  gPThFiahk4DMoehrHMRqI=,tag:72KKv3PUvzlav97KN6A5xg==,type:str]
  pgp: []
  encrypted_regex: encryptedPatches
  version: 3.7.3`

	err := yaml.Unmarshal([]byte(data1), &m1)
	if err != nil {
		t.Fatal(err)
	}

	err = yaml.Unmarshal([]byte(data2), &m2)
	if err != nil {
		t.Fatal(err)
	}

	ans1 := m1.isEncrypted()
	ans2 := m2.isEncrypted()
	if ans1 != false || ans2 != true {
		t.Errorf("got %t %t, want false true", ans1, ans2)
	}
}
