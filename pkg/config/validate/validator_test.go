package validate

import (
	"testing"
)

func TestIsRFC6902List(t *testing.T) {
	c := Config{}
	data := make([]map[string]interface{}, 2)
	data1 := map[string]interface{}{
		"op":    "add",
		"path":  "/a/path",
		"value": "a value",
	}
	data2 := map[string]interface{}{
		"op":   "remove",
		"path": "/a/path",
	}

	data[0] = data1
	data[1] = data2

	expected := true

	if !c.IsRFC6902List(data) {
		t.Errorf("got %t, want %t", c.IsRFC6902List(data), expected)
	}
}

func TestIsSupportedK8sVersion(t *testing.T) {
	c := &Config{
		TalosVersion: "v1.3.7",
	}
	data := map[string]bool{
		"v1.23.0":  false,
		"v1.24.0":  true,
		"v1.25.1":  true,
		"v1.26.99": false,
	}

	for k, v := range data {
		if c.IsSupportedK8sVersion(k) != v {
			t.Errorf("%s: got %t, want %t", k, c.IsSupportedK8sVersion(k), v)
		}
	}
}

func TestIsSemVer(t *testing.T) {
	c := Config{}
	data := map[string]bool{
		"v1.2.3-4":  true,
		"1.2.3-444": true,
		"v12.3":     false,
		"v1.2.3.4":  false,
		"v1.2.3":    true,
	}

	for k, v := range data {
		if c.IsSemVer(k) != v {
			t.Errorf("%s: got %t, want %t", k, c.IsSemVer(k), v)
		}
	}
}

func TestIsCNIName(t *testing.T) {
	c := Config{}
	data := map[string]bool{
		"none":     true,
		"flannel":  true,
		"custom":   true,
		"nonea":    false,
		"aflannel": false,
	}

	for k, v := range data {
		if c.IsCNIName(k) != v {
			t.Errorf("%s: got %t, want %t", k, c.IsCNIName(k), v)
		}
	}
}

func TestIsCIDRList(t *testing.T) {
	c := Config{}
	data1 := []string{"0.0.0.0/0", "1.2.3.4/24"}
	data2 := []string{"0.0.0.0/0", "1.2.3.4"}

	if !c.IsCIDRList(data1) {
		t.Errorf("got %t, want true", c.IsCIDRList(data1))
	}

	if c.IsCIDRList(data2) {
		t.Errorf("got %t, want false", c.IsCIDRList(data2))
	}
}

func TestIsIPList(t *testing.T) {
	c := Config{}
	data1 := []string{"0.0.0.0", "1.2.3.4"}
	data2 := []string{"0.0.0.0.0", "1.2.3.4"}

	if !c.IsIPList(data1) {
		t.Errorf("got %t, want true", c.IsIPList(data1))
	}

	if c.IsIPList(data2) {
		t.Errorf("got %t, want false", c.IsIPList(data2))
	}
}

func TestIsURLList(t *testing.T) {
	c := Config{}
	data1 := []string{"https://www.www/path", "www.www/path", "www/path/file.ext"}
	data2 := []string{"htt_://www/path", "www.www/path"}

	if !c.IsURLList(data1) {
		t.Errorf("got %t, want true", c.IsURLList(data1))
	}

	if c.IsURLList(data2) {
		t.Errorf("got %t, want false", c.IsURLList(data2))
	}
}

func TestIsTalosEndpoint(t *testing.T) {
	c := Config{}
	data := map[string]bool{
		"1.1.1.1:443":          false,
		"http://hostname":      true,
		"http://hostname:6443": true,
		"https://1.1.1.1:6443": true,
	}

	for k, v := range data {
		if c.IsTalosEndpoint(k) != v {
			t.Errorf("%s: got %t, want %t", k, c.IsTalosEndpoint(k), v)
		}
	}
}

func TestIsDomain(t *testing.T) {
	c := Config{}
	data := map[string]bool{
		"adomain,test": false,
		"adomain.test": true,
		"adomain":      true,
		"adomain?":     false,
	}

	for k, v := range data {
		if c.IsDomain(k) != v {
			t.Errorf("%s: got %t, want %t", k, c.IsDomain(k), v)
		}
	}
}
