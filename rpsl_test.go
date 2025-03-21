// Copyright (c) The RPSL Go Authors.
// SPDX-License-Identifier: Apache-2.0

package rpsl

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func TestIntegration(t *testing.T) {
	files, err := os.ReadDir("tests/data")
	if err != nil {
		t.Fatalf("unable to read directory: %v\nDid you run scripts/download-dumps.sh?", err)
	}

	datasets := make([]string, 0)
	for _, file := range files {
		datasets = append(datasets, file.Name())
	}

	for _, dataset := range datasets {
		t.Run(dataset, func(t *testing.T) {
			data, err := os.Open("tests/data/" + dataset)
			if err != nil {
				t.Fatalf("unable to read file: %v", err)
			}

			objects, err := parseObjects(data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(objects) == 0 {
				t.Fatalf(`parseObjectsFromBytes => length of %v, want > 0`, len(objects))
			}
		})
	}
}

func TestParseAPI(t *testing.T) {
	// Test the string-based Parse function (original API).
	rawString := "person: John Doe\naddress: 123 Example St\nnic-hdl: JD1-RIPE\nsource: RIPE"
	obj, err := Parse(rawString)
	if err != nil {
		t.Fatalf("Parse(string) error: %v", err)
	}

	if obj == nil {
		t.Fatalf("Parse(string) returned nil object")
	}

	if obj.Len() != 4 {
		t.Fatalf("Parse(string) object length: got %v, want 4", obj.Len())
	}

	// Test the byte-based ParseFromBytes function (new API).
	rawBytes := bytes.NewReader([]byte("person: Jane Smith\naddress: 456 Example Ave\nnic-hdl: JS1-RIPE\nsource: RIPE"))
	obj, err = ParseFromReader(rawBytes)
	if err != nil {
		t.Fatalf("ParseFromBytes error: %v", err)
	}

	if obj == nil {
		t.Fatalf("ParseFromBytes returned nil object")
	}

	if obj.Len() != 4 {
		t.Fatalf("ParseFromBytes object length: got %v, want 4", obj.Len())
	}
}

func TestParseManyAPI(t *testing.T) {
	// Test the string-based ParseMany function (original API).
	rawString := "person: John Doe\naddress: 123 Example St\nnic-hdl: JD1-RIPE\nsource: RIPE\n\nperson: Jane Smith\naddress: 456 Example Ave\nnic-hdl: JS1-RIPE\nsource: RIPE"
	objects, err := ParseMany(rawString)
	if err != nil {
		t.Fatalf("ParseMany(string) error: %v", err)
	}

	if len(objects) != 2 {
		t.Fatalf("ParseMany(string) objects count: got %v, want 2", len(objects))
	}

	// Test the byte-based ParseManyFromBytes function (new API).
	rawBytes := bytes.NewReader([]byte("person: Alice Brown\naddress: 789 Example Blvd\nnic-hdl: AB1-RIPE\nsource: RIPE\n\nperson: Bob Green\naddress: 101 Example Ct\nnic-hdl: BG1-RIPE\nsource: RIPE"))
	objects, err = ParseManyFromReader(rawBytes)
	if err != nil {
		t.Fatalf("ParseManyFromBytes error: %v", err)
	}

	if len(objects) != 2 {
		t.Fatalf("ParseManyFromBytes objects count: got %v, want 2", len(objects))
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectErr    bool
		errSubstring string
		validateObj  func(*Object) error
	}{
		{
			name:         "EmptyInput",
			input:        "",
			expectErr:    true,
			errSubstring: "no objects found",
		},
		{
			name:         "InvalidSyntax",
			input:        "this is not a valid RPSL object",
			expectErr:    true,
			errSubstring: "parseKey: illegal character",
		},
		{
			name: "SingleObject",
			input: "person:  John Doe\n" +
				"address: 123 Main St\n" +
				"phone:   +1-555-1234\n" +
				"source:  TEST",
			expectErr: false,
			validateObj: func(obj *Object) error {
				if obj == nil {
					return errors.New("object is nil")
				}
				if obj.Len() != 4 {
					return errors.New("expected 4 attributes")
				}

				person := obj.GetFirst("person")
				if person == nil || *person != "John Doe" {
					return errors.New("person attribute incorrect")
				}

				return nil
			},
		},
		{
			name: "MultipleObjects",
			input: "person:  John Doe\n" +
				"source:  TEST\n" +
				"\n" +
				"person:  Jane Smith\n" +
				"source:  TEST",
			expectErr:    true,
			errSubstring: "multiple objects found",
		},
		{
			name: "ObjectWithComments",
			input: "person:  John Doe # Inline comment\n" +
				"address: 123 Main St\n" +
				"source:  TEST",
			expectErr: false,
			validateObj: func(obj *Object) error {
				if obj == nil {
					return errors.New("object is nil")
				}
				if obj.Len() != 3 {
					return errors.New("expected 3 attributes")
				}

				address := obj.GetFirst("address")
				if address == nil || *address != "123 Main St" {
					return errors.New("address attribute incorrect")
				}

				return nil
			},
		},
		{
			name: "ObjectWithContinuationLines",
			input: "person:  John Doe\n" +
				"address: 123 Main St\n" +
				"remarks: First line\n" +
				" Continuation line 1\n" +
				" Continuation line 2\n" +
				"source:  TEST",
			expectErr: false,
			validateObj: func(obj *Object) error {
				if obj == nil {
					return errors.New("object is nil")
				}

				remarks := obj.GetFirst("remarks")
				if remarks == nil || *remarks != "First line Continuation line 1 Continuation line 2" {
					return errors.New("remarks with continuation lines incorrect")
				}

				return nil
			},
		},
		{
			name: "ObjectWithPlusNotation",
			input: "person:  John Doe\n" +
				"address: \n" +
				"+123 Main St\n" +
				"+Suite 100\n" +
				"source:  TEST",
			expectErr: false,
			validateObj: func(obj *Object) error {
				if obj == nil {
					return errors.New("object is nil")
				}

				address := obj.GetFirst("address")
				if address == nil || *address != "123 Main St Suite 100" {
					return errors.New("address with plus notation incorrect")
				}

				return nil
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			obj, err := Parse(tc.input)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if tc.errSubstring != "" && !strings.Contains(err.Error(), tc.errSubstring) {
					t.Errorf("Error %q does not contain expected substring %q", err.Error(), tc.errSubstring)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tc.validateObj != nil {
				if err = tc.validateObj(obj); err != nil {
					t.Errorf("Object validation failed: %v", err)
				}
			}
		})
	}
}

func TestParseMany(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectErr     bool
		errSubstring  string
		expectedCount int
		validateObjs  func([]Object) error
	}{
		{
			name:          "EmptyInput",
			input:         "",
			expectErr:     false,
			expectedCount: 0,
		},
		{
			name:         "InvalidSyntax",
			input:        "this is : not a valid RPSL object",
			expectErr:    true,
			errSubstring: "parseKey: illegal character",
		},
		{
			name: "SingleObject",
			input: "person:  John Doe\n" +
				"address: 123 Main St\n" +
				"source:  TEST",
			expectErr:     false,
			expectedCount: 1,
			validateObjs: func(objs []Object) error {
				if len(objs) != 1 {
					return errors.New("expected 1 object")
				}
				if objs[0].Len() != 3 {
					return errors.New("object has wrong number of attributes")
				}
				return nil
			},
		},
		{
			name: "MultipleObjects",
			input: "person:  John Doe\n" +
				"address: 123 Main St\n" +
				"source:  TEST\n" +
				"\n" +
				"person:  Jane Smith\n" +
				"address: 456 Oak Ave\n" +
				"source:  TEST",
			expectErr:     false,
			expectedCount: 2,
			validateObjs: func(objs []Object) error {
				if len(objs) != 2 {
					return errors.New("expected 2 objects")
				}

				person1 := objs[0].GetFirst("person")
				if person1 == nil || *person1 != "John Doe" {
					return errors.New("first person attribute incorrect")
				}

				person2 := objs[1].GetFirst("person")
				if person2 == nil || *person2 != "Jane Smith" {
					return errors.New("second person attribute incorrect")
				}

				return nil
			},
		},
		{
			name: "MultipleObjectsWithComments",
			input: "# Comment at start\n" +
				"person:  John Doe\n" +
				"source:  TEST\n" +
				"\n" +
				"# Comment before second object\n" +
				"person:  Jane Smith\n" +
				"source:  TEST\n" +
				"\n" +
				"# Comment before third object\n" +
				"person:  Bob Johnson\n" +
				"source:  TEST",
			expectErr:     false,
			expectedCount: 3,
			validateObjs: func(objs []Object) error {
				if len(objs) != 3 {
					return errors.New("expected 3 objects")
				}

				person3 := objs[2].GetFirst("person")
				if person3 == nil || *person3 != "Bob Johnson" {
					return errors.New("third person attribute incorrect")
				}

				return nil
			},
		},
		{
			name: "ObjectsWithMultipleEmptyLinesBetween",
			input: "person:  John Doe\n" +
				"source:  TEST\n" +
				"\n\n\n" +
				"person:  Jane Smith\n" +
				"source:  TEST",
			expectErr:     false,
			expectedCount: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			objs, err := ParseMany(tc.input)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if tc.errSubstring != "" && !strings.Contains(err.Error(), tc.errSubstring) {
					t.Errorf("Error %q does not contain expected substring %q", err.Error(), tc.errSubstring)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tc.expectedCount == 0 {
				if objs != nil {
					t.Errorf("Expected nil objects for empty input, got %v objects", len(objs))
				}
				return
			}

			if len(objs) != tc.expectedCount {
				t.Errorf("Expected %d objects, got %d", tc.expectedCount, len(objs))
				return
			}

			if tc.validateObjs != nil {
				if err = tc.validateObjs(objs); err != nil {
					t.Errorf("Objects validation failed: %v", err)
				}
			}
		})
	}
}

func TestParseManyFromReader(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		expectErr     bool
		errSubstring  string
		expectedCount int
		validateObjs  func([]Object) error
	}{
		{
			name:          "EmptyInput",
			input:         []byte{},
			expectErr:     false,
			expectedCount: 0,
		},
		{
			name:          "NilInput",
			input:         nil,
			expectErr:     false,
			expectedCount: 0,
		},
		{
			name: "MultipleObjectsWithDifferentAttributes",
			input: []byte("person:  John Doe\n" +
				"address: 123 Main St\n" +
				"phone:   +1-555-1234\n" +
				"source:  TEST\n" +
				"\n" +
				"organisation: ORG-EXAMPLE\n" +
				"org-name:    Example Organization\n" +
				"org-type:    OTHER\n" +
				"address:     789 Corporate Blvd\n" +
				"source:      TEST"),
			expectErr:     false,
			expectedCount: 2,
			validateObjs: func(objs []Object) error {
				if len(objs) != 2 {
					return errors.New("expected 2 objects")
				}

				// Check first object is a person
				if !objs[0].Exists("person") {
					return errors.New("first object should be a person")
				}

				// Check second object is an organisation
				orgName := objs[1].GetFirst("org-name")
				if orgName == nil || *orgName != "Example Organization" {
					return errors.New("org-name attribute incorrect")
				}

				// Check attribute counts
				if objs[0].Len() != 4 {
					return errors.New("first object should have 4 attributes")
				}

				if objs[1].Len() != 5 {
					return errors.New("second object should have 5 attributes")
				}

				return nil
			},
		},
		{
			name: "ObjectWithMultipleAttributesOfSameType",
			input: []byte("route:     192.168.0.0/16\n" +
				"origin:    AS12345\n" +
				"descr:     Primary Network\n" +
				"descr:     Corporate Use Only\n" +
				"descr:     No Public Access\n" +
				"source:    TEST"),
			expectErr:     false,
			expectedCount: 1,
			validateObjs: func(objs []Object) error {
				if len(objs) != 1 {
					return errors.New("expected 1 object")
				}

				// Check multiple descriptions
				descrs := objs[0].GetAll("descr")
				if len(descrs) != 3 {
					return errors.New("expected 3 description attributes")
				}

				if descrs[0] != "Primary Network" ||
					descrs[1] != "Corporate Use Only" ||
					descrs[2] != "No Public Access" {
					return errors.New("description attributes incorrect")
				}

				return nil
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			objs, err := parseObjects(bytes.NewReader(tc.input))
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if tc.errSubstring != "" && !strings.Contains(err.Error(), tc.errSubstring) {
					t.Errorf("Error %q does not contain expected substring %q", err.Error(), tc.errSubstring)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(objs) != tc.expectedCount {
				t.Errorf("Expected %d objects, got %d", tc.expectedCount, len(objs))
				return
			}

			if tc.validateObjs != nil {
				if err = tc.validateObjs(objs); err != nil {
					t.Errorf("Objects validation failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkParseFromReader(b *testing.B) {
	data := bytes.NewReader([]byte(`mntner:          DEV-MNT  # Comment \n" +
		"descr:           DEV maintainer\n" +
		"admin-c:         VM1-DEV\n" +
		"tech-c:          VM1-DEV\n" +
		"upd-to:          v.m@example.net\n" +
		"mnt-nfy:         auto@example.net\n" +
		"auth:            MD5-PW $1$q8Su3Hq/$rJt5M3TNLeRE4UoCh5bSH/\n" +
		"remarks:         password: secret\n" +
		"mnt-by:          DEV-MNT\n" +
		"source:          DEV\n

`))

	b.ResetTimer()
	for range b.N {
		if _, err := parseObjects(data); err != nil {
			b.Fatalf("ParseManyFromReader error: %v", err)
		}

		if _, err := data.Seek(0, io.SeekStart); err != nil {
			b.Fatalf("ParseManyFromReader error: %v", err)
		}
	}
}
