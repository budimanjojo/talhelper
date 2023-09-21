package config

import "testing"

func TestIsRFC6902List(t *testing.T) {
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

	if !isRFC6902List(data) {
		t.Errorf("got %t, want %t", isRFC6902List(data), expected)
	}
}

func TestIsSemVer(t *testing.T) {
	data := map[string]bool{
		"v1.2.3-4":  true,
		"1.2.3-444": true,
		"v12.3":     false,
		"v1.2.3.4":  false,
		"v1.2.3":    true,
	}

	for k, v := range data {
		if isSemVer(k) != v {
			t.Errorf("%s: got %t, want %t", k, isSemVer(k), v)
		}
	}
}

func TestIsCIDRList(t *testing.T) {
	data1 := []string{"0.0.0.0/0", "1.2.3.4/24"}
	data2 := []string{"0.0.0.0/0", "1.2.3.4"}

	if !isCIDRList(data1) {
		t.Errorf("got %t, want true", isCIDRList(data1))
	}

	if isCIDRList(data2) {
		t.Errorf("got %t, want false", isCIDRList(data2))
	}
}

func TestIsDomain(t *testing.T) {
	data := map[string]bool{
		"adomain,test": false,
		"adomain.test": true,
		"adomain":      true,
		"adomain?":     false,
	}

	for k, v := range data {
		if isDomain(k) != v {
			t.Errorf("%s: got %t, want %t", k, isDomain(k), v)
		}
	}
}
