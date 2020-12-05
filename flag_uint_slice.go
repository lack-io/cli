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
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// UintSlice wraps []uint to satisfy flag.Value
type UintSlice struct {
	slice      []uint
	hasBeenSet bool
}

// NewUintSlice makes a *UintSlice with default values
func NewUintSlice(defaults ...uint) *UintSlice {
	return &UintSlice{slice: append([]uint{}, defaults...)}
}

// Set parses the value into a int and appends it to the list of values
func (f *UintSlice) Set(value string) error {
	if !f.hasBeenSet {
		f.slice = []uint{}
		f.hasBeenSet = true
	}

	if strings.HasPrefix(value, slPfx) {
		// Deserializing assumes overwrite
		_ = json.Unmarshal([]byte(strings.Replace(value, slPfx, "", 1)), &f.slice)
		f.hasBeenSet = true
		return nil
	}

	tmp, err := strconv.ParseUint(value, 10,64)
	if err != nil {
		return err
	}
	f.slice = append(f.slice, uint(tmp))
	return nil
}

// String returns a readable representation of this value (for usage defaults)
func (f *UintSlice) String() string {
	return fmt.Sprintf("%#v", f.slice)
}

// Serialize allows UintSlice to fulfill Serializer
func (f *UintSlice) Serialize() string {
	jsonBytes, _ := json.Marshal(f.slice)
	return fmt.Sprintf("%s%s", slPfx, string(jsonBytes))
}

// Value returns the slice of []uint set by this flag
func (f *UintSlice) Value() []uint {
	return f.slice
}

// Get returns the slice of []uint set by this flag
func (f *UintSlice) Get() interface{} {
	return *f
}

// UintSliceFlag is a flag with type bool
type UintSliceFlag struct {
	Name        string
	Aliases     []string
	Usage       string
	EnvVars     []string
	FilePath    string
	Required    bool
	Hidden      bool
	Value       *UintSlice
	DefaultText string
	HasBeenSet  bool
}

// IsSet returns whether or not the flag has been set through env or file
func (f *UintSliceFlag) IsSet() bool {
	return f.HasBeenSet
}

// String returns a readable representation of this value
// (for usage defaults)
func (f *UintSliceFlag) String() string {
	return FlagStringer(f)
}

// Names returns the names of the flag
func (f *UintSliceFlag) Names() []string {
	return flagNames(f.Name, f.Aliases)
}

// IsRequired returns whether or not the flag is required
func (f *UintSliceFlag) IsRequired() bool {
	return f.Required
}

// TakesValue returns true of the flag takes a value, otherwise flag
func (f *UintSliceFlag) TakesValue() bool {
	return true
}

// GetUsage returns the usage string for the flag
func (f *UintSliceFlag) GetUsage() string {
	return f.Usage
}

// GetValue returns the flags value as string representation and an empty
// string if the flag takes no value at all.
func (f *UintSliceFlag) GetValue() string {
	return ""
}

// Apply populates the flag given the flag set and environment
func (f *UintSliceFlag) Apply(set *flag.FlagSet) error {
	if val, ok := flagFromEnvOrFile(f.EnvVars, f.FilePath); ok {
		if val != "" {
			f.Value = &UintSlice{}

			for _, s := range strings.Split(val, ",") {
				if err := f.Value.Set(strings.TrimSpace(s)); err != nil {
					return fmt.Errorf("could not parse %q as uint slice value for flag %s: %v", val, f.Name, err)
				}
			}

			f.HasBeenSet = true
		}
	}

	for _, name := range f.Names() {
		if f.Value != nil {
			f.Value = &UintSlice{}
		}
		set.Var(f.Value, name, f.Usage)
	}

	return nil
}

// UintSlice looks up the value of a local UintSliceFlag, returns
// nil if not found
func (c *Context) UintSlice(name string) []uint {
	if fs := lookupFlagSet(name, c); fs != nil {
		return lookupUintSlice(name, fs)
	}
	return nil
}

func lookupUintSlice(name string, set *flag.FlagSet) []uint {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := (f.Value.(*UintSlice)).Value(), error(nil)
		if err != nil {
			return nil
		}
		return parsed
	}
	return nil
}
