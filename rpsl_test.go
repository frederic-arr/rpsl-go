// Copyright (c) The RPSL Go Authors.
// SPDX-License-Identifier: Apache-2.0

package rpsl

import (
	"os"
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
			data, err := os.ReadFile("tests/data/" + dataset)
			if err != nil {
				t.Fatalf("unable to read file: %v", err)
			}

			objects, err := parseObjects(string(data))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(objects) == 0 {
				t.Fatalf(`parseObjects => length of %v, want %v`, len(objects), 0)
			}
		})
	}
}
