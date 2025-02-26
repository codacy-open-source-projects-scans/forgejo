// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build !bindata

package migration

import (
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type SchemaLoader struct{}

func (*SchemaLoader) Load(s string) (any, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	basename := path.Base(u.Path)
	filename := basename
	//
	// Schema reference each other within the schemas directory but
	// the tests run in the parent directory.
	//
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		filename = filepath.Join("schemas", basename)
		//
		// Integration tests run from the git root directory, not the
		// directory in which the test source is located.
		//
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			filename = filepath.Join("modules/migration/schemas", basename)
		}
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return jsonschema.UnmarshalJSON(f)
}
