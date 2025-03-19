// Copyright (c) The RPSL Go Authors.
// SPDX-License-Identifier: Apache-2.0

package rpsl

import (
	"errors"
	"fmt"
	"strings"
)

type Object struct {
	Attributes []Attribute
}

// Keys returns a slice of unique keys present in the Object.
// If a key appears multiple times in the Object, it will only be included once in the returned slice.
func (o *Object) Keys() []string {
	keyPresent := make(map[string]struct{})
	keyList := make([]string, 0)
	for _, attr := range o.Attributes {
		if _, exists := keyPresent[attr.Name]; !exists {
			keyPresent[attr.Name] = struct{}{}
			keyList = append(keyList, attr.Name)
		}
	}

	return keyList
}

// Len returns the number of attributes in the Object.
func (o *Object) Len() int {
	return len(o.Attributes)
}

// GetFirst returns the first value for a given key in the Object.
// If the key is not present in the Object, an empty string will be returned.
// If a key appears multiple times in the Object, only the first value will be returned.
func (o *Object) GetFirst(key string) *string {
	key = strings.ToLower(key)
	for _, attr := range o.Attributes {
		if attr.Name == key {
			return &attr.Value
		}
	}

	return nil
}

// GetAll returns a slice of values for a given key in the Object.
// If the key is not present in the Object, an empty slice will be returned.
// If a key appears multiple times in the Object, all values will be included in the returned slice.
func (o *Object) GetAll(key string) []string {
	attributes := make([]string, 0)
	for _, attr := range o.Attributes {
		if attr.Name == key {
			attributes = append(attributes, attr.Value)
		}
	}

	return attributes
}

// Exists returns true if the Object contains a given key.
func (o *Object) Exists(key string) bool {
	key = strings.ToLower(key)
	for _, attr := range o.Attributes {
		if attr.Name == key {
			return true
		}
	}

	return false
}

// String returns a string representation of the Object.
func (o *Object) String() string {
	var str strings.Builder

	var attributes []string
	for _, attr := range o.Attributes {
		attributes = append(attributes, attr.String())
	}

	str.WriteString(strings.Join(attributes, "\n"))
	return str.String()
}

// EnsureClass ensures that the first attribute in the Object is of a given class.
func (o *Object) EnsureClass(class string) error {
	if len(o.Attributes) == 0 {
		return errors.New("object has no attributes")
	}

	first := o.Attributes[0].Name
	if first != class {
		return fmt.Errorf("attribute '%s' should be the first, but found '%s' instead", class, first)
	}

	return nil
}

// EnsureAtLeastOne ensures that the Object has at least one attribute with a given key.
func (o *Object) EnsureAtLeastOne(key string) error {
	if !o.Exists(key) {
		return fmt.Errorf("attribute '%s' is (mandatory, multiple) but found none", key)
	}

	return nil
}

// EnsureAtMostOne ensures that the Object has at most one attribute with a given key.
func (o *Object) EnsureAtMostOne(key string) error {
	if len(o.GetAll(key)) > 1 {
		return fmt.Errorf("attribute '%s' is (optional, single) but found multiple", key)
	}

	return nil
}

// EnsureOne ensures that the Object has exactly one attribute with a given key.
func (o *Object) EnsureOne(key string) error {
	if err := o.EnsureAtLeastOne(key); err != nil {
		return fmt.Errorf("attribute '%s' is (mandatory, single) but found none", key)
	}

	if err := o.EnsureAtMostOne(key); err != nil {
		return fmt.Errorf("attribute '%s' is (mandatory, single) but found multiple", key)
	}

	return nil
}

func parseObjects(buf string) ([]Object, error) {
	var objects []Object

	if buf == "" {
		return objects, nil
	}

	lines := strings.Split(buf, "\n")

	// Process line by line, accumulating non-comment lines.
	var currentPart []string

	for i := 0; i <= len(lines); i++ {
		// Process a line or handle end of input.
		isEndOfFile := i == len(lines)
		isEmptyLine := !isEndOfFile && lines[i] == ""

		// When we hit an empty line or EOF, process the accumulated part.
		if isEmptyLine || isEndOfFile {
			if len(currentPart) > 0 {
				partText := strings.Join(currentPart, "\n")
				attributes, err := parseAttributes(partText)
				if err != nil {
					return nil, err
				}

				objects = append(objects, Object{Attributes: attributes})
				currentPart = currentPart[:0] // Clear slice without reallocating.
			}
			continue
		}

		// Skip comment lines.
		line := lines[i]
		if !strings.HasPrefix(line, "%") && !strings.HasPrefix(line, "#") {
			currentPart = append(currentPart, line)
		}
	}

	return objects, nil
}
