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
	"strings"
)

type iterativeParser interface {
	newFlagSet() (*flag.FlagSet, error)
	useShortOptionHandling() bool
}

// To enable short-option handling (e.g., "-it" vs "-i -t") we have to
// iteratively catch parsing errors. This way we achieve LR parsing without
// transforming any arguments. Otherwise, there is no way we can discriminate
// combined short options from common arguments that should be left untouched.
// Pass `shellComplete` to continue parsing options on failure during shell
// completion when, the user-supplied options may be incomplete.
func parseIter(set *flag.FlagSet, ip iterativeParser, args []string, shellComplete bool) error {
	for {
		err := set.Parse(args)
		if !ip.useShortOptionHandling() || err == nil {
			if shellComplete {
				return nil
			}
			return err
		}

		errStr := err.Error()
		trimmed := strings.TrimPrefix(errStr, "flag provided but not default: -")
		if errStr == trimmed {
			return err
		}

		// regenerate the initial args with the split short opts
		argsWereSplit := false
		for i, arg := range args {
			// skip args that are not part of the error message
			if name := strings.TrimPrefix(arg, "-"); name != trimmed {
				continue
			}

			// if we can't split, the error was accurate
			shortOpts := splitShortOptions(set, arg)
			if len(shortOpts) == 1 {
				return err
			}

			// swap current argument with the split version
			args = append(args[:i], append(shortOpts, args[i+1:]...)...)
			argsWereSplit = true
			break
		}

		// This should be an impossible to reach code path, but in case the arg
		// splitting failed to happen, this will prevent infinite loops
		if !argsWereSplit {
			return err
		}

		// Since custom parsing failed, replace the flag set before retrying
		newSet, err := ip.newFlagSet()
		if err != nil {
			return err
		}
		*set = *newSet
	}
}

func splitShortOptions(set *flag.FlagSet, arg string) []string {
	shortFlagsExists := func(s string) bool {
		for _, c := range s[1:] {
			if f := set.Lookup(string(c)); f == nil {
				return false
			}
		}
		return true
	}

	if !isSplittable(arg) || !shortFlagsExists(arg) {
		return []string{arg}
	}

	separated := make([]string, 0, len(arg)-1)
	for _, flagChar := range arg[1:] {
		separated = append(separated, "-"+string(flagChar))
	}

	return separated
}

func isSplittable(flagArg string) bool {
	return strings.HasPrefix(flagArg, "-") && !strings.HasPrefix(flagArg, "--") && len(flagArg) > 2
}