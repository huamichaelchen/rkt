// Copyright 2016 The rkt Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/coreos/rkt/pkg/log"
)

func TestCalculateDataDir(t *testing.T) {
	// Used in calculateDataDir.
	// TODO(tmrts): Restructure this pkg, specifically by using dependency
	// injection, to eliminate these work-arounds.
	stderr = log.New(os.Stderr, "TestCalculateDataDir", globalFlags.Debug)

	_, err := getConfig()
	if err != nil {
		panic(fmt.Errorf("getConfig() got error %q", err))
	}

	if cachedConfig == nil {
		panic(fmt.Errorf("getConfig() should've set `cachedConfig`"))
	}

	dirFlag := cmdRkt.PersistentFlags().Lookup("dir")
	dirFlag.Changed = false
	defDirFlagVal := dirFlag.Value.String()
	defCfgDataDir := cachedConfig.Paths.DataDir

	resetConfigState := func() {
		cmdRkt.PersistentFlags().Set("dir", defDirFlagVal)
		dirFlag.Changed = false

		cachedConfig.Paths.DataDir = defCfgDataDir
	}

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(fmt.Errorf("ioutil.TempDir(%q, %q) got error %q", "", "", err))
	}
	defer os.Remove(tmpDir)

	// TODO(tmrts): Write a utility function that generates unused random paths.
	// Example signature would be fileutils.GenerateUniquePath(prefix string) (string, error).
	nonExistantDir, err := ioutil.TempDir("", "non-existant-")
	if err != nil {
		panic(fmt.Errorf("ioutil.TempDir(%q, %q) got error %q", "", "", err))
	}
	if err := os.Remove(nonExistantDir); err != nil {
		panic(fmt.Errorf("os.Remove(%q) got error %q", nonExistantDir, err))
	}

	testCases := []struct {
		flagDataDir   string
		configDataDir string
		out           string
	}{
		{"", "", defaultDataDir},
		{"", tmpDir, tmpDir},
		{tmpDir, "", tmpDir},
		{nonExistantDir, "", nonExistantDir},
		{"", nonExistantDir, nonExistantDir},
	}

	for _, tc := range testCases {
		cmdRkt.PersistentFlags().Set("dir", tc.flagDataDir)

		cachedConfig.Paths.DataDir = tc.configDataDir

		if dataDir := calculateDataDir(); dataDir != tc.out {
			t.Errorf("main.calculateDataDir() with setup %q, expected %q, got %q", tc, tc.out, dataDir)
		}

		resetConfigState()
	}
}
