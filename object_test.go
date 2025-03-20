// Copyright (c) The RPSL Go Authors.
// SPDX-License-Identifier: Apache-2.0

package rpsl

import (
	"bytes"
	"testing"
)

// Helper function to create string pointer.
func stringPtr(s string) *string {
	return &s
}

func TestObject(t *testing.T) {
	raw := []byte("organisation:      ORG-CEOf1-RIPE\n" +
		"description:       CERN")

	objects, err := parseObjects(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf(`parseObject => %v`, err)
	}

	if len(objects) != 1 {
		t.Fatalf(`parseObject => length of %v, want %v`, len(objects), 1)
	}

	obj := objects[0]
	if len(obj.Attributes) != 2 {
		t.Fatalf(`object.Attributes => length of %v, want %v`, len(obj.Attributes), 2)
	}
}

func TestObjectLen(t *testing.T) {
	raw := []byte("organisation:      ORG-CEOf1-RIPE\n" +
		"remarks:           This is a comment\n" +
		"description:       CERN\n" +
		"remarks:           This is another comment")

	objects, err := parseObjects(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf(`parseObject => %v`, err)
	}

	if len(objects) != 1 {
		t.Fatalf(`parseObject => length of %v, want %v`, len(objects), 1)
	}

	obj := objects[0]
	if obj.Len() != 4 {
		t.Fatalf(`object.Len() => %v, want %v`, obj.Len(), 4)
	}
}

func TestGetFirst(t *testing.T) {
	tests := []struct {
		name          string
		object        Object
		key           string
		expectedValue *string
		shouldBeNil   bool
	}{
		{
			name: "KeyExistsOnce",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "address", Value: "123 Main St"},
				},
			},
			key:           "person",
			expectedValue: stringPtr("John Doe"),
		},
		{
			name: "KeyExistsMultipleTimes",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "remarks", Value: "Remark 1"},
					{Name: "remarks", Value: "Remark 2"},
				},
			},
			key:           "remarks",
			expectedValue: stringPtr("Remark 1"),
		},
		{
			name: "CaseInsensitiveMatch",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "address", Value: "123 Main St"},
				},
			},
			key:           "ADDRESS",
			expectedValue: stringPtr("123 Main St"),
		},
		{
			name: "KeyDoesNotExist",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "address", Value: "123 Main St"},
				},
			},
			key:         "phone",
			shouldBeNil: true,
		},
		{
			name: "EmptyObject",
			object: Object{
				Attributes: []Attribute{},
			},
			key:         "person",
			shouldBeNil: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.object.GetFirst(tc.key)

			if tc.shouldBeNil {
				if result != nil {
					t.Errorf("Expected nil result, but got %v", *result)
				}
			} else {
				if result == nil {
					t.Errorf("Expected non-nil result, but got nil")
				} else if *result != *tc.expectedValue {
					t.Errorf("Result = %q, want %q", *result, *tc.expectedValue)
				}
			}
		})
	}
}

func TestObjectKeys(t *testing.T) {
	raw := []byte("organisation:      ORG-CEOf1-RIPE\n" +
		"remarks:           This is a comment\n" +
		"description:       CERN\n" +
		"remarks:           This is another comment")

	objects, err := parseObjects(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf(`parseObject => %v`, err)
	}

	if len(objects) != 1 {
		t.Fatalf(`parseObject => length of %v, want %v`, len(objects), 1)
	}

	obj := objects[0]
	keys := obj.Keys()
	if len(keys) != 3 {
		t.Fatalf(`object.Keys() => length of %v, want %v`, len(keys), 3)
	}

	if keys[0] != "organisation" {
		t.Fatalf(`object.Keys()[0] => %v, want %v`, keys[0], "organisation")
	}

	if keys[1] != "remarks" {
		t.Fatalf(`object.Keys()[1] => %v, want %v`, keys[1], "remarks")
	}

	if keys[2] != "description" {
		t.Fatalf(`object.Keys()[2] => %v, want %v`, keys[2], "description")
	}
}

func TestObjectGetAll(t *testing.T) {
	raw := []byte("organisation:      ORG-CEOf1-RIPE\n" +
		"remarks:           This is a comment\n" +
		"description:       CERN\n" +
		"remarks:           This is another comment")

	objects, err := parseObjects(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf(`parseObject => %v`, err)
	}

	if len(objects) != 1 {
		t.Fatalf(`parseObject => length of %v, want %v`, len(objects), 1)
	}

	obj := objects[0]
	attrs := obj.GetAll("organisation")
	if len(attrs) != 1 {
		t.Fatalf(`object.GetAll("organisation") => length of %v, want %v`, len(attrs), 1)
	}

	if attrs[0] != "ORG-CEOf1-RIPE" {
		t.Fatalf(`object.GetAll("organisation") => %v, want %v`, attrs[0], "ORG-CEOf1-RIPE")
	}

	attrs = obj.GetAll("description")
	if len(attrs) != 1 {
		t.Fatalf(`object.GetAll("description") => length of %v, want %v`, len(attrs), 1)
	}

	if attrs[0] != "CERN" {
		t.Fatalf(`object.GetAll("description") => %v, want %v`, attrs[0], "CERN")
	}

	attrs = obj.GetAll("remarks")
	if len(attrs) != 2 {
		t.Fatalf(`object.GetAll("remarks") => length of %v, want %v`, len(attrs), 2)
	}

	if attrs[0] != "This is a comment" {
		t.Fatalf(`object.GetAll("remarks")[0] => %v, want %v`, attrs[0], "This is a comment")
	}

	if attrs[1] != "This is another comment" {
		t.Fatalf(`object.GetAll("remarks")[1] => %v, want %v`, attrs[1], "This is another comment")
	}
}

func TestMultipleObjects(t *testing.T) {
	data := []byte("" +
		"poem:           POEM-LIR\n" +
		"form:           FORM-HAIKU\n" +
		"text:           hello ripe please\n" +
		"text:           consider this offer, make lir\n" +
		"text:           just for free\n" +
		"descr:          Does RIPE still allow creation of these objects?\n" +
		"created:        2024-04-30T18:06:01Z\n" +
		"last-modified:  2024-04-30T18:06:01Z\n" +
		"source:         RIPE\n" +
		"mnt-by:         DUMMY-MNT\n" +
		"\n" +
		"poem:           poem-ipv6-adoption\n" +
		"form:           FORM-HAIKU\n" +
		"text:           Bound by old NAT's chains,\n" +
		"text:           Joy of routing slips away,\n" +
		"text:           IPv6 scorned.\n" +
		"author:         DUMY-RIPE\n" +
		"notify:         dummy@example.com\n" +
		"mnt-by:         dummy-mnt\n" +
		"created:        2024-06-01T23:28:08Z\n" +
		"last-modified:  2024-06-01T23:28:08Z\n" +
		"source:         RIPE\n")

	objects, err := parseObjects(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("(error): %v", err)
	}

	if len(objects) != 2 {
		t.Fatalf("(length): got %v, want %v", len(objects), 2)
	}

	if objects[0].Len() != 10 {
		t.Fatalf("(0.length): got %v, want %v", objects[0].Len(), 10)
	}

	if objects[1].Len() != 11 {
		t.Fatalf("(1.length): got %v, want %v", objects[1].Len(), 11)
	}
}

func TestObjectString(t *testing.T) {
	tests := []struct {
		name     string
		object   Object
		expected string
	}{
		{
			name: "EmptyObject",
			object: Object{
				Attributes: []Attribute{},
			},
			expected: "",
		},
		{
			name: "SingleAttribute",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
				},
			},
			expected: "person:John Doe",
		},
		{
			name: "MultipleAttributes",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "address", Value: "123 Main St"},
					{Name: "phone", Value: "+1-555-1234"},
				},
			},
			expected: "person:John Doe\naddress:123 Main St\nphone:+1-555-1234",
		},
		{
			name: "DuplicateAttributes",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "remarks", Value: "Remark 1"},
					{Name: "remarks", Value: "Remark 2"},
				},
			},
			expected: "person:John Doe\nremarks:Remark 1\nremarks:Remark 2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.object.String()
			if result != tc.expected {
				t.Errorf("Object.String() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestEnsureClass(t *testing.T) {
	tests := []struct {
		name      string
		object    Object
		class     string
		expectErr bool
		errMsg    string
	}{
		{
			name: "CorrectClass",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "address", Value: "123 Main St"},
				},
			},
			class:     "person",
			expectErr: false,
		},
		{
			name: "CaseInsensitiveMatch",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "address", Value: "123 Main St"},
				},
			},
			class:     "PERSON",
			expectErr: false,
		},
		{
			name: "IncorrectClass",
			object: Object{
				Attributes: []Attribute{
					{Name: "route", Value: "192.168.0.0/16"},
					{Name: "origin", Value: "AS12345"},
				},
			},
			class:     "person",
			expectErr: true,
			errMsg:    "attribute 'person' should be the first, but found 'route' instead",
		},
		{
			name: "EmptyObject",
			object: Object{
				Attributes: []Attribute{},
			},
			class:     "person",
			expectErr: true,
			errMsg:    "object has no attributes",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.object.EnsureClass(tc.class)

			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if err.Error() != tc.errMsg {
					t.Errorf("Error message = %q, want %q", err.Error(), tc.errMsg)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestEnsureAtLeastOne(t *testing.T) {
	tests := []struct {
		name      string
		object    Object
		key       string
		expectErr bool
		errMsg    string
	}{
		{
			name: "KeyExistsOnce",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "address", Value: "123 Main St"},
				},
			},
			key:       "address",
			expectErr: false,
		},
		{
			name: "KeyExistsMultipleTimes",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "remarks", Value: "Remark 1"},
					{Name: "remarks", Value: "Remark 2"},
				},
			},
			key:       "remarks",
			expectErr: false,
		},
		{
			name: "CaseInsensitiveMatch",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "address", Value: "123 Main St"},
				},
			},
			key:       "ADDRESS",
			expectErr: false,
		},
		{
			name: "KeyDoesNotExist",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "address", Value: "123 Main St"},
				},
			},
			key:       "phone",
			expectErr: true,
			errMsg:    "attribute 'phone' is (mandatory, multiple) but found none",
		},
		{
			name: "EmptyObject",
			object: Object{
				Attributes: []Attribute{},
			},
			key:       "person",
			expectErr: true,
			errMsg:    "attribute 'person' is (mandatory, multiple) but found none",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.object.EnsureAtLeastOne(tc.key)

			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if err.Error() != tc.errMsg {
					t.Errorf("Error message = %q, want %q", err.Error(), tc.errMsg)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestEnsureAtMostOne(t *testing.T) {
	tests := []struct {
		name      string
		object    Object
		key       string
		expectErr bool
		errMsg    string
	}{
		{
			name: "KeyDoesNotExist",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "address", Value: "123 Main St"},
				},
			},
			key:       "phone",
			expectErr: false,
		},
		{
			name: "KeyExistsOnce",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "address", Value: "123 Main St"},
				},
			},
			key:       "address",
			expectErr: false,
		},
		{
			name: "CaseInsensitiveMatch",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "address", Value: "123 Main St"},
				},
			},
			key:       "ADDRESS",
			expectErr: false,
		},
		{
			name: "KeyExistsMultipleTimes",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "remarks", Value: "Remark 1"},
					{Name: "remarks", Value: "Remark 2"},
				},
			},
			key:       "remarks",
			expectErr: true,
			errMsg:    "attribute 'remarks' is (optional, single) but found multiple",
		},
		{
			name: "EmptyObject",
			object: Object{
				Attributes: []Attribute{},
			},
			key:       "person",
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.object.EnsureAtMostOne(tc.key)

			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if err.Error() != tc.errMsg {
					t.Errorf("Error message = %q, want %q", err.Error(), tc.errMsg)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestEnsureOne(t *testing.T) {
	tests := []struct {
		name      string
		object    Object
		key       string
		expectErr bool
		errMsg    string
	}{
		{
			name: "KeyExistsOnce",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "address", Value: "123 Main St"},
				},
			},
			key:       "address",
			expectErr: false,
		},
		{
			name: "CaseInsensitiveMatch",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "address", Value: "123 Main St"},
				},
			},
			key:       "ADDRESS",
			expectErr: false,
		},
		{
			name: "KeyDoesNotExist",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "address", Value: "123 Main St"},
				},
			},
			key:       "phone",
			expectErr: true,
			errMsg:    "attribute 'phone' is (mandatory, single) but found none",
		},
		{
			name: "KeyExistsMultipleTimes",
			object: Object{
				Attributes: []Attribute{
					{Name: "person", Value: "John Doe"},
					{Name: "remarks", Value: "Remark 1"},
					{Name: "remarks", Value: "Remark 2"},
				},
			},
			key:       "remarks",
			expectErr: true,
			errMsg:    "attribute 'remarks' is (mandatory, single) but found multiple",
		},
		{
			name: "EmptyObject",
			object: Object{
				Attributes: []Attribute{},
			},
			key:       "person",
			expectErr: true,
			errMsg:    "attribute 'person' is (mandatory, single) but found none",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.object.EnsureOne(tc.key)

			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if err.Error() != tc.errMsg {
					t.Errorf("Error message = %q, want %q", err.Error(), tc.errMsg)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
