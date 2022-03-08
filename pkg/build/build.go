// Copyright 2022 Chainguard, Inc.
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

package build

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
	apko_types "chainguard.dev/apko/pkg/build/types"
)

type Package struct {
	Name               string
	Version            string
	Epoch              uint64
	Description        string
	TargetArchitecture []string `yaml:"target-architecture"`
	Copyright          []Copyright
	Dependencies       Dependencies
}

type Copyright struct {
	Paths       []string
	Attestation string
	License     string
}

type Pipeline struct {
	Uses string
	With map[string]string
	Runs string
}

type Subpackage struct {
	Name     string
	Pipeline []Pipeline
}

type Configuration struct {
	Package     Package
	Environment apko_types.ImageConfiguration
	Pipeline    []Pipeline
	Subpackages []Subpackage
}

type Context struct {
	Configuration   Configuration
	ConfigFile      string
	SourceDateEpoch time.Time
}

type Dependencies struct {
	Runtime []string
}

func New(opts ...Option) (*Context, error) {
	ctx := Context{
		ConfigFile: ".melange.yaml",
	}

	for _, opt := range opts {
		if err := opt(&ctx); err != nil {
			return nil, err
		}
	}

	if err := ctx.Configuration.Load(ctx.ConfigFile); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// SOURCE_DATE_EPOCH will always overwrite the build flag
	if v, ok := os.LookupEnv("SOURCE_DATE_EPOCH"); ok {
		// The value MUST be an ASCII representation of an integer
		// with no fractional component, identical to the output
		// format of date +%s.
		sec, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			// If the value is malformed, the build process
			// SHOULD exit with a non-zero error code.
			return nil, fmt.Errorf("failed to parse SOURCE_DATE_EPOCH: %w", err)
		}

		ctx.SourceDateEpoch = time.Unix(sec, 0)
	}

	return &ctx, nil
}

type Option func(*Context) error

// WithConfig sets the configuration file used for the package build context.
func WithConfig(configFile string) Option {
	return func(ctx *Context) error {
		ctx.ConfigFile = configFile
		return nil
	}
}

// WithBuildDate sets the timestamps for the build context.
// The string is parsed according to RFC3339.
// An empty string is a special case and will default to
// the unix epoch.
func WithBuildDate(s string) Option {
	return func(bc *Context) error {
		// default to 0 for reproducibility
		if s == "" {
			bc.SourceDateEpoch = time.Unix(0, 0)
			return nil
		}

		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}

		bc.SourceDateEpoch = t
		return nil
	}
}

// Load the configuration data from the build context configuration file.
func (cfg *Configuration) Load(configFile string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("unable to load configuration file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("unable to parse configuration file: %w", err)
	}

	return nil
}

func (ctx *Context) BuildPackage() error {
	return nil
}