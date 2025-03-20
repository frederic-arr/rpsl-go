// Copyright (c) The RPSL Go Authors.
// SPDX-License-Identifier: Apache-2.0

package rpsl

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Object struct {
	Attributes []Attribute
}

// Keys returns a slice of unique keys present in the Object. If a key appears multiple times in the Object, it will
// only be included once in the returned slice.
func (o *Object) Keys() []string {
	// Most objects have < 16 unique keys.
	keyPresent := make(map[string]struct{}, 16)
	keyList := make([]string, 0, 16)

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

// GetFirst returns the first value for a given key in the Object. If the key is not present in the Object, an empty
// string will be returned. If a key appears multiple times in the Object, only the first value will be returned.
func (o *Object) GetFirst(key string) *string {
	key = strings.ToLower(key)
	for _, attr := range o.Attributes {
		if attr.Name == key {
			return &attr.Value
		}
	}

	return nil
}

// GetAll returns a slice of values for a given key in the Object. If the key is not present in the Object, an empty
// slice will be returned. If a key appears multiple times in the Object, all values will be included in the returned
// slice.
func (o *Object) GetAll(key string) []string {
	key = strings.ToLower(key)

	// Check how many attributes match this key to pre-allocate.
	count := 0
	for _, attr := range o.Attributes {
		if attr.Name == key {
			count++
		}
	}

	if count == 0 {
		return []string{}
	}

	attributes := make([]string, 0, count)
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
	// Pre-allocate a reasonably sized buffer.
	var str strings.Builder
	str.Grow(len(o.Attributes) * 64) // Assume average attribute length of ~64 chars.

	for i, attr := range o.Attributes {
		if i > 0 {
			str.WriteByte('\n')
		}
		str.WriteString(attr.Name)
		str.WriteByte(':')
		str.WriteString(attr.Value)
	}

	return str.String()
}

// EnsureClass ensures that the first attribute in the Object is of a given class.
func (o *Object) EnsureClass(class string) error {
	if len(o.Attributes) == 0 {
		return errors.New("object has no attributes")
	}

	first := o.Attributes[0].Name
	if first != strings.ToLower(class) {
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
	// Get count without allocating a slice.
	count := 0
	key = strings.ToLower(key)
	for _, attr := range o.Attributes {
		if attr.Name == key {
			count++
			if count > 1 {
				return fmt.Errorf("attribute '%s' is (optional, single) but found multiple", key)
			}
		}
	}

	return nil
}

// EnsureOne ensures that the Object has exactly one attribute with a given key.
func (o *Object) EnsureOne(key string) error {
	count := 0
	exists := false
	key = strings.ToLower(key)

	for _, attr := range o.Attributes {
		if attr.Name == key {
			exists = true
			count++
			if count > 1 {
				return fmt.Errorf("attribute '%s' is (mandatory, single) but found multiple", key)
			}
		}
	}

	if !exists {
		return fmt.Errorf("attribute '%s' is (mandatory, single) but found none", key)
	}

	return nil
}

func parseObjects(buf []byte) ([]Object, error) {
	// Most files contain a reasonable number of objects.
	objects := make([]Object, 0, 4)

	if len(buf) == 0 {
		return objects, nil
	}

	// Process comment lines first.
	// This is more efficient than processing them line by line.
	lines := bytes.Split(buf, []byte("\n"))

	// In-place filtering of comment lines.
	for i, line := range lines {
		if len(line) > 0 && (line[0] == '%' || line[0] == '#') {
			lines[i] = []byte{}
		}
	}

	buf = bytes.Join(lines, []byte("\n"))

	// Split by double newlines to get objects.
	parts := bytes.Split(buf, []byte("\n\n"))
	for _, part := range parts {
		part = bytes.TrimSpace(part)

		if len(part) == 0 {
			continue
		}

		attributes, err := parseAttributes(part)
		if err != nil {
			return nil, err
		}

		if len(attributes) > 0 {
			object := Object{Attributes: attributes}
			objects = append(objects, object)
		}
	}

	return objects, nil
}
