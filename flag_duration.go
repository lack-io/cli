// Copyright 2020 The vine Authors
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

package cli

import (
	"flag"
	"fmt"
	"time"
)

// DurationFlag is a flag with type bool
type DurationFlag struct {
	Name        string
	Aliases     []string
	Usage       string
	EnvVars     []string
	FilePath    string
	Required    bool
	Hidden      bool
	Value       time.Duration
	DefaultText string
	Destination *time.Duration
	HasBeenSet  bool
}

// IsSet returns whether or not the flag has been set through env or file
func (f *DurationFlag) IsSet() bool {
	return f.HasBeenSet
}

// String returns a readable representation of this value
// (for usage defaults)
func (f *DurationFlag) String() string {
	return FlagStringer(f)
}

// Names returns the names of the flag
func (f *DurationFlag) Names() []string {
	return flagNames(f.Name, f.Aliases)
}

// IsRequired returns whether or not the flag is required
func (f *DurationFlag) IsRequired() bool {
	return f.Required
}

// TakesValue returns true of the flag takes a value, otherwise flag
func (f *DurationFlag) TakesValue() bool {
	return true
}

// GetUsage returns the usage string for the flag
func (f *DurationFlag) GetUsage() string {
	return f.Usage
}

// GetValue returns the flags value as string representation and an empty
// string if the flag takes no value at all.
func (f *DurationFlag) GetValue() string {
	return ""
}

// Apply populates the flag given the flag set and environment
func (f *DurationFlag) Apply(set *flag.FlagSet) error {
	if val, ok := flagFromEnvOrFile(f.EnvVars, f.FilePath); ok {
		if val != "" {
			valDuration, err := time.ParseDuration(val)

			if err != nil {
				return fmt.Errorf("could not parse %q as duration value for flag %s: %v", val, f.Name, err)
			}

			f.Value = valDuration
			f.HasBeenSet = true
		}
	}

	for _, name := range f.Names() {
		if f.Destination != nil {
			set.DurationVar(f.Destination, name, f.Value, f.Usage)
			continue
		}
		set.Duration(name, f.Value, f.Usage)
	}

	return nil
}

// Duration looks up the value of a local DurationFlag, returns
// 0 if not found
func (c *Context) Duration(name string) time.Duration {
	if fs := lookupFlagSet(name, c); fs != nil {
		return lookupDuration(name, fs)
	}
	return 0
}

func lookupDuration(name string, set *flag.FlagSet) time.Duration {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := time.ParseDuration(f.Value.String())
		if err != nil {
			return 0
		}
		return parsed
	}
	return 0
}