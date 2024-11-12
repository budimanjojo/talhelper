package config

import (
	"strings"
	"testing"
)

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

func TestCheckNodeLabels(t *testing.T) {
	tests := []struct {
		name        string
		labels      map[string]string
		expectError bool
	}{
		{
			name: "single-no-template-value",
			labels: map[string]string{
				"valid-label": "valid-value",
			},
		},
		{
			name: "single-no-template-value-invalid",
			labels: map[string]string{
				"valid-label": strings.Repeat("a", 64),
			},
			expectError: true,
		},
		{
			name: "multiple-no-template-values",
			labels: map[string]string{
				"valid-label1": "valid-value1",
				"valid-label2": "valid-value2",
			},
		},
		{
			name: "single-template-values",
			labels: map[string]string{
				"valid-label": "{{ .Node.Hostname }}",
			},
		},
		{
			name: "single-template-value-invalid-for-non-template",
			labels: map[string]string{
				"valid-label": strings.Repeat("{{ 'a' }}", 64),
			},
		},
		{
			name: "mixed-values",
			labels: map[string]string{
				"valid-label1": "valid-value1",
				"valid-label2": "valid-value2",
				"valid-label":  "{{ .Node.Hostname }}",
			},
		},
	}

	for _, test := range tests {
		// Setup
		node := Node{
			NodeConfigs: NodeConfigs{
				NodeLabels: test.labels,
			},
		}

		// Test
		var providedErrors Errors
		resultErrors := checkNodeLabels(node, 0, &providedErrors)

		// Verify results
		if &providedErrors != resultErrors {
			t.Errorf("%s: provided errors and resultant errors don't point to the same set of values", test.name)
		}

		if len(*resultErrors) > 0 && !test.expectError {
			t.Errorf("%s: didn't expect an error but received %#v", test.name, *resultErrors)
		} else if len(*resultErrors) == 0 && test.expectError {
			t.Errorf("%s: exepected an error but didn't receive any", test.name)
		}
	}
}

func TestCheckNodeAnnotations(t *testing.T) {
	tests := []struct {
		name        string
		annotations map[string]string
		expectError bool
	}{
		{
			name: "single-no-template-value",
			annotations: map[string]string{
				"valid-annotation": "valid-value",
			},
		},
		{
			name: "multiple-no-template-values",
			annotations: map[string]string{
				"valid-annotation1": "valid-value1",
				"valid-annotation2": "valid-value2",
			},
		},
		{
			name: "single-template-values",
			annotations: map[string]string{
				"valid-annotation": "{{ .Node.Hostname }}",
			},
		},
		{
			name: "single-template-value-invalid-for-non-template",
			annotations: map[string]string{
				"valid-annotation": strings.Repeat("{{ 'a' }}", 64),
			},
		},
		{
			name: "mixed-values",
			annotations: map[string]string{
				"valid-annotation1": "valid-value1",
				"valid-annotation2": "valid-value2",
				"valid-annotation":  "{{ .Node.Hostname }}",
			},
		},
	}

	for _, test := range tests {
		// Setup
		node := Node{
			NodeConfigs: NodeConfigs{
				NodeAnnotations: test.annotations,
			},
		}

		// Test
		var providedErrors Errors
		resultErrors := checkNodeAnnotations(node, 0, &providedErrors)

		// Verify results
		if &providedErrors != resultErrors {
			t.Errorf("%s: provided errors and resultant errors don't point to the same set of values", test.name)
		}

		if len(*resultErrors) > 0 && !test.expectError {
			t.Errorf("%s: didn't expect an error but received %#v", test.name, *resultErrors)
		} else if len(*resultErrors) == 0 && test.expectError {
			t.Errorf("%s: exepected an error but didn't receive any", test.name)
		}
	}
}
