package templating

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/Masterminds/sprig/v3"
)

// The split functions cannot be replaced with a single generic one, because there is
// no common type between maps and slices

// Given a map `src`, return two maps where the first contains key-value pairs that are not
// templates, and the latter containers only key-value pairs that are templates.
func SplitTemplatedMapItems(src map[string]string) (map[string]string, map[string]string) {
	nonTemplateKeyPairs := make(map[string]string, len(src))
	templateKeyPairs := make(map[string]string, len(src))

	for key, value := range src {
		if IsStringATemplate(value) {
			templateKeyPairs[key] = value
		} else {
			nonTemplateKeyPairs[key] = value
		}
	}

	return nonTemplateKeyPairs, templateKeyPairs
}

// Given a slice `src`, return two slices where the first contains items that are not
// templates, and the latter containers only items that are templates. Values of the
// opposing types will be set to an empty string, to preserve the position of all items.
func SplitTemplatedListItems(src []string) ([]string, []string) {
	nonTemplateItems := make([]string, len(src))
	templateKeyItems := make([]string, len(src))

	for index, value := range src {
		if IsStringATemplate(value) {
			templateKeyItems[index] = value
		} else {
			nonTemplateItems[index] = value
		}
	}

	return nonTemplateItems, templateKeyItems
}

// Returns true if the provided string contains a template, false otherwise
func IsStringATemplate(str string) bool {
	return strings.Contains(str, "{{") || strings.Contains(str, "}}")
}

// Render each value in the provided map.
func RenderMap[T any](templateKeyPairs map[string]string, data any) (map[string]T, error) {
	renderedKeyPairs := make(map[string]T, len(templateKeyPairs))
	for key, templateValue := range templateKeyPairs {
		renderedValue, err := RenderTemplate[T](templateValue, data)
		if err != nil {
			return nil, err
		}

		renderedKeyPairs[key] = renderedValue
	}

	return renderedKeyPairs, nil
}

// Render each template item in the provided slice.
func RenderList[T any](templateItems []string, data any) ([]T, error) {
	renderedItems := make([]T, len(templateItems))
	for i, templateValue := range templateItems {
		renderedValue, err := RenderTemplate[T](templateValue, data)
		if err != nil {
			return nil, err
		}

		renderedItems[i] = renderedValue
	}

	return renderedItems, nil
}

// Take a template string, render it, and convert it to the desired type.
// Only literal types (i.e. string, int, etc.) are supported.
func RenderTemplate[T any](templateText string, data any) (T, error) {
	var t T

	builtTemplate, err := template.New("template").Funcs(sprig.FuncMap()).Parse(templateText)
	if err != nil {
		return t, err
	}

	if _, ok := any(t).([]any); ok {
		// This almost definitely indicates a code bug, but unfortunately Go
		// does not support checking this via generics
		return t, fmt.Errorf("list types are not supported")
	}

	var b strings.Builder
	if err := builtTemplate.Execute(&b, data); err != nil {
		return t, err
	}

	if _, ok := any(t).(string); ok {
		return any(b.String()).(T), nil
	}

	if _, ok := any(t).([]byte); ok {
		return any([]byte(b.String())).(T), nil
	}

	// Convert the string to the target type
	if _, err := fmt.Sscan(b.String(), &t); err != nil {
		return t, err
	}
	return t, nil
}
