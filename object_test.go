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
	data := "" +
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

	objects, err := parseObjects(data)
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
