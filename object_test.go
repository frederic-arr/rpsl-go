// Copyright (c) The RPSL Go Authors.
// SPDX-License-Identifier: Apache-2.0

package rpsl

import (
	"testing"
)

func TestObject(t *testing.T) {
	raw := "organisation:      ORG-CEOf1-RIPE\n" +
		"description:       CERN"

	objects, err := parseObjects(raw)
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

func TestObjectKeys(t *testing.T) {
	raw := "organisation:      ORG-CEOf1-RIPE\n" +
		"remarks:           This is a comment\n" +
		"description:       CERN\n" +
		"remarks:           This is another comment"

	objects, err := parseObjects(raw)
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

func TestObjectLen(t *testing.T) {
	raw := "organisation:      ORG-CEOf1-RIPE\n" +
		"remarks:           This is a comment\n" +
		"description:       CERN\n" +
		"remarks:           This is another comment"

	objects, err := parseObjects(raw)
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

func TestObjectGetFirst(t *testing.T) {
	// Helper function for string pointer.
	stringPtr := func(s string) *string {
		return &s
	}

	tests := []struct {
		name       string
		attributes []Attribute
		key        string
		want       *string
	}{
		{
			name: "MatchFirstAttribute",
			attributes: []Attribute{
				{Name: "source", Value: "RIPE"},
				{Name: "admin-c", Value: "ABC123"},
			},
			key:  "source",
			want: stringPtr("RIPE"),
		},
		{
			name: "MatchSecondAttribute",
			attributes: []Attribute{
				{Name: "origin", Value: "AS1234"},
				{Name: "source", Value: "RIPE"},
			},
			key:  "source",
			want: stringPtr("RIPE"),
		},
		{
			name: "CaseInsensitiveKeyMatch",
			attributes: []Attribute{
				{Name: "source", Value: "RIPE"},
			},
			key:  "SOURCE",
			want: stringPtr("RIPE"),
		},
		{
			name: "MultipleAttributesWithSameKey",
			attributes: []Attribute{
				{Name: "remarks", Value: "First remark"},
				{Name: "remarks", Value: "Second remark"},
			},
			key:  "remarks",
			want: stringPtr("First remark"),
		},
		{
			name:       "EmptyAttributes",
			attributes: []Attribute{},
			key:        "source",
			want:       nil,
		},
		{
			name: "KeyNotFound",
			attributes: []Attribute{
				{Name: "source", Value: "RIPE"},
				{Name: "admin-c", Value: "ABC123"},
			},
			key:  "origin",
			want: nil,
		},
		{
			name: "EmptyValue",
			attributes: []Attribute{
				{Name: "remarks", Value: ""},
			},
			key:  "remarks",
			want: stringPtr(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Object{Attributes: tt.attributes}

			got := o.GetFirst(tt.key)

			if (got == nil && tt.want != nil) || (got != nil && tt.want == nil) {
				t.Errorf("GetFirst(%q) = %v, want %v", tt.key, got, tt.want)
			} else if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("GetFirst(%q) = %q, want %q", tt.key, *got, *tt.want)
			}
		})
	}
}

func TestObjectGetAll(t *testing.T) {
	raw := "organisation:      ORG-CEOf1-RIPE\n" +
		"remarks:           This is a comment\n" +
		"description:       CERN\n" +
		"remarks:           This is another comment"

	objects, err := parseObjects(raw)
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
	raw := "" +
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
		"source:         RIPE\n"

	objects, err := parseObjects(raw)
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
