package templating

import "testing"

func TestSplitTemplatedMapItems(t *testing.T) {
	tests := []struct {
		name                string
		srcKeyPairs         map[string]string
		nonTemplateKeyPairs map[string]string
		templateKeyPairs    map[string]string
	}{
		{
			name: "single-non-template",
			srcKeyPairs: map[string]string{
				"nonTemplateKey1": "nonTemplateValue1",
			},
			nonTemplateKeyPairs: map[string]string{
				"nonTemplateKey1": "nonTemplateValue1",
			},
		},
		{
			name: "single-template",
			srcKeyPairs: map[string]string{
				"templateKey1": "{{ templateValue1 }}",
			},
			templateKeyPairs: map[string]string{
				"templateKey1": "{{ templateValue1 }}",
			},
		},
		{
			name: "multiple-mixed",
			srcKeyPairs: map[string]string{
				"nonTemplateKey1": "nonTemplateValue1",
				"nonTemplateKey2": "nonTemplateValue2",
				"templateKey1":    "{{ templateValue1 }}",
				"templateKey2":    "{{ templateValue2 }}",
			},
			nonTemplateKeyPairs: map[string]string{
				"nonTemplateKey1": "nonTemplateValue1",
				"nonTemplateKey2": "nonTemplateValue2",
			},
			templateKeyPairs: map[string]string{
				"templateKey1": "{{ templateValue1 }}",
				"templateKey2": "{{ templateValue2 }}",
			},
		},
	}

	for _, test := range tests {
		actualNonTemplateKeyPairs, actualTemplateKeyPairs := SplitTemplatedMapItems(test.srcKeyPairs)
		compareMaps(t, test.nonTemplateKeyPairs, actualNonTemplateKeyPairs)
		compareMaps(t, test.templateKeyPairs, actualTemplateKeyPairs)
	}
}

func TestSplitTemplatedListItems(t *testing.T) {
	tests := []struct {
		name             string
		srcItems         []string
		nonTemplateItems []string
		templateItems    []string
	}{
		{
			name: "single-non-template",
			srcItems: []string{
				"nonTemplateValue1",
			},
			nonTemplateItems: []string{
				"nonTemplateValue1",
			},
			templateItems: make([]string, 1),
		},
		{
			name: "single-template",
			srcItems: []string{
				"{{ templateValue1 }}",
			},
			nonTemplateItems: make([]string, 1),
			templateItems: []string{
				"{{ templateValue1 }}",
			},
		},
		{
			name: "multiple-mixed",
			srcItems: []string{
				"nonTemplateValue1",
				"{{ templateValue1 }}",
				"{{ templateValue2 }}",
				"nonTemplateValue2",
			},
			nonTemplateItems: []string{
				"nonTemplateValue1",
				"",
				"",
				"nonTemplateValue2",
			},
			templateItems: []string{
				"",
				"{{ templateValue1 }}",
				"{{ templateValue2 }}",
				"",
			},
		},
	}

	for _, test := range tests {
		actualNonTemplateItems, actualTemplateItems := SplitTemplatedListItems(test.srcItems)
		compareLists(t, test.nonTemplateItems, actualNonTemplateItems)
		compareLists(t, test.templateItems, actualTemplateItems)
	}
}

func TestIsStringATemplate(t *testing.T) {
	tests := []struct {
		name          string
		val           string
		isNonTemplate bool
	}{
		{
			name:          "empty-string",
			isNonTemplate: true,
		},
		{
			name:          "non-template",
			val:           "some test value",
			isNonTemplate: true,
		},
		{
			name: "basic-template",
			val:  "{{ .Basic.Template }}",
		},
		{
			name: "empty-braces",
			val:  "{{ }}",
		},
		{
			name: "braces-no-spaces",
			val:  "{{.No.Spaces}}",
		},
		{
			name:          "single-braces",
			val:           "{ .No.Spaces }",
			isNonTemplate: true,
		},
		{
			name: "only-opening",
			val:  "{{",
		},
		{
			name: "only-closing",
			val:  "}}",
		},
		{
			name: "out-of-order",
			val:  "}} {{ .Out.Of.Order }}",
		},
	}

	for _, test := range tests {
		result := IsStringATemplate(test.val)
		if result && test.isNonTemplate {
			t.Errorf("%s: result indicated template but test case was for non-template", test.name)
		} else if !result && !test.isNonTemplate {
			t.Errorf("%s: result indicated non-template but test case was for template", test.name)

		}
	}
}

func TestRenderMap(t *testing.T) {
	tests := []struct {
		name        string
		templates   map[string]string
		rendered    map[string]string
		expectError bool
	}{
		{
			name: "single-value",
			templates: map[string]string{
				"key1": "{{ .DummyField1 }}",
			},
			rendered: map[string]string{
				"key1": "val1",
			},
		},
		{
			name: "multiple-values",
			templates: map[string]string{
				"key1": "{{ .DummyField1 }}",
				"key2": "{{ .DummyField2 }}",
			},
			rendered: map[string]string{
				"key1": "val1",
				"key2": "val2",
			},
		},
		{
			name: "invalid-field",
			templates: map[string]string{
				"key1": "{{ .NonExistantField }}",
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		data := struct {
			DummyField1 string
			DummyField2 string
		}{
			DummyField1: "val1",
			DummyField2: "val2",
		}

		renderedValues, err := RenderMap[string](test.templates, data)

		if err == nil && test.expectError {
			t.Errorf("%s: expected error but didn't generate one", test.name)
		} else if err != nil && !test.expectError {
			t.Errorf("%s: didn't expect error but generated one", test.name)
		}
		compareMaps(t, test.rendered, renderedValues)
	}
}

func TestRenderList(t *testing.T) {
	tests := []struct {
		name        string
		templates   []string
		rendered    []string
		expectError bool
	}{
		{
			name: "single-value",
			templates: []string{
				"{{ .DummyField1 }}",
			},
			rendered: []string{
				"val1",
			},
		},
		{
			name: "multiple-values",
			templates: []string{
				"{{ .DummyField1 }}",
				"{{ .DummyField2 }}",
			},
			rendered: []string{
				"val1",
				"val2",
			},
		},
		{
			name: "invalid-field",
			templates: []string{
				"{{ .NonExistantField }}",
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		data := struct {
			DummyField1 string
			DummyField2 string
		}{
			DummyField1: "val1",
			DummyField2: "val2",
		}

		renderedValues, err := RenderList[string](test.templates, data)

		if err == nil && test.expectError {
			t.Errorf("%s: expected error but didn't generate one", test.name)
		} else if err != nil && !test.expectError {
			t.Errorf("%s: didn't expect error but generated one", test.name)
		}
		compareLists(t, test.rendered, renderedValues)
	}
}

func TestRenderTemplate(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		value       string
		expectError bool
	}{
		{
			name:     "simple-value",
			template: "{{ \"Simple Value\" }}",
			value:    "Simple Value",
		},
		{
			name:     "data-value",
			template: "{{ .DummyField1 }}",
			value:    "val1",
		},
		{
			name:     "multiple-templates",
			template: "f1:{{ .DummyField1 }} f2:{{ .DummyField2 }}",
			value:    "f1:val1 f2:val2",
		},
		{
			name:     "computed-value",
			template: "{{ .DummyField1 | upper }}",
			value:    "VAL1",
		},
		{
			name:        "invalid-template",
			template:    "{{ .NonExistantField }}",
			expectError: true,
		},
	}

	for _, test := range tests {
		data := struct {
			DummyField1 string
			DummyField2 string
		}{
			DummyField1: "val1",
			DummyField2: "val2",
		}

		actualValue, err := RenderTemplate[string](test.template, data)

		if err == nil && test.expectError {
			t.Errorf("%s: expected error but didn't generate one", test.name)
		} else if err != nil && !test.expectError {
			t.Errorf("%s: didn't expect error but generated one", test.name)
		}

		if actualValue != test.value {
			t.Errorf("%s: expected %q, got %q", test.name, test.value, actualValue)
		}
	}
}

func compareMaps[T comparable](t *testing.T, expected map[string]T, actual map[string]T) {
	for expectedKey, expectedValue := range expected {
		if actualValue, ok := actual[expectedKey]; ok {
			if actualValue != expectedValue {
				t.Errorf("actual value '%v' did not match expected value '%v'", expectedValue, actualValue)
			}
		} else {
			t.Errorf("missing key %q in actual values", expectedKey)
		}
	}
}

func compareLists[T comparable](t *testing.T, expected []T, actual []T) {
	if len(expected) != len(actual) {
		t.Errorf("Expected list of length %d, got list of length %d", len(expected), len(actual))
	}

	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("Values at index %d don't match. Expected %#v, got %#v", i, expected[i], actual[i])
		}
	}
}
