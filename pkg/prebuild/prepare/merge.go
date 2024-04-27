// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"os"
	"path/filepath"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
)

type Merge struct {
	cfg.Base
}

func init() {
	RegisterTask(&Merge{
		Base: cfg.Base{
			Keyword: "merge",
			Msg:     "Merge all profiles into a unified apparmor.d directory",
		},
	})
}

func (p Merge) Apply() ([]string, error) {
	res := []string{}
	dirToMerge := []string{
		"groups/*/*", "groups",
		"profiles-*-*/*", "profiles-*",
	}

	idx := 0
	for idx < len(dirToMerge)-1 {
		dirMoved, dirRemoved := dirToMerge[idx], dirToMerge[idx+1]
		files, err := filepath.Glob(cfg.RootApparmord.Join(dirMoved).String())
		if err != nil {
			return res, err
		}
		for _, file := range files {
			err := os.Rename(file, cfg.RootApparmord.Join(filepath.Base(file)).String())
			if err != nil {
				return res, err
			}
		}

		files, err = filepath.Glob(cfg.RootApparmord.Join(dirRemoved).String())
		if err != nil {
			return []string{}, err
		}
		for _, file := range files {
			if err := paths.New(file).RemoveAll(); err != nil {
				return res, err
			}
		}
		idx = idx + 2
	}
	return res, nil
}
