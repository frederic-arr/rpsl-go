// Copyright (c) The RPSL Go Authors.
// SPDX-License-Identifier: Apache-2.0

package rpsl

import (
	"errors"
)

// Parse parses a string containing a single RPSL object
// and returns a representation of the parsed data.
// If the string contains multiple objects, an error will be returned.
// If the string is empty, an error will be returned.
//
// Example:
//
//	raw := "person:	John Doe\n"+
//	    "address:	1234 Elm Street\n"+
//	    "phone:		+1 555 123456\n"+
//	    "nic-hdl:	JD1234-RIPE\n"+
//	    "mnt-by:	EXAMPLE-MNT\n"+
//	    "source:	RIPE\n"
//	obj, err := Parse(raw)
//
//	if err != nil {
//	    log.Fatalf("Failed to parse RPSL object: %v", err)
//	}
//
//	fmt.Printf("Parsed Object: %+v\n", obj)
func Parse(raw string) (*Object, error) {
	objects, err := parseObjects(raw)
	if err != nil {
		return nil, err
	}

	if len(objects) == 0 {
		return nil, errors.New("no objects found in input")
	}

	if len(objects) > 1 {
		return nil, errors.New("multiple objects found in input")
	}

	return &objects[0], nil
}

// ParseMany parses a string containing a multiple RPSL object
// and returns a representation of the parsed data.
// If the string does not contain any objects, nil will be returned.
//
// Example:
//
//	raw := "person:	John Doe\n"+
//		"address:	1234 Elm Street\n"+
//		"phone:		+1 555 123456\n"+
//		"nic-hdl:	JD1234-RIPE\n"+
//		"mnt-by:	EXAMPLE-MNT\n"+
//		"source:	RIPE\n"+
//		"\n"+ // Objects are delimited by one or more empty lines
//		"person:	Jane Smith\n"+
//		"address:	5678 Oak Street\n"+
//		"phone:		+1 555 654321\n"+
//		"nic-hdl:	JS5678-RIPE\n"+
//		"mnt-by:	EXAMPLE-MNT\n"+
//		"source:	RIPE"
//	obj, err := ParseMany(raw)
//
//	if err != nil {
//		log.Fatalf("Failed to parse RPSL object: %v", err)
//	}
//
//	for _, obj := range *objs {
//		fmt.Printf("Parsed Object: %+v\n", obj)
//	}
func ParseMany(raw string) ([]Object, error) {
	objects, err := parseObjects(raw)
	if err != nil {
		return nil, err
	}

	if len(objects) == 0 {
		return nil, nil
	}

	return objects, nil
}
