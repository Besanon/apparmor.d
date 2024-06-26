// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"slices"
)

const SIGNAL Kind = "signal"

func init() {
	requirements[SIGNAL] = requirement{
		"access": {
			"r", "w", "rw", "read", "write", "send", "receive",
		},
		"set": {
			"hup", "int", "quit", "ill", "trap", "abrt", "bus", "fpe",
			"kill", "usr1", "segv", "usr2", "pipe", "alrm", "term", "stkflt",
			"chld", "cont", "stop", "stp", "ttin", "ttou", "urg", "xcpu",
			"xfsz", "vtalrm", "prof", "winch", "io", "pwr", "sys", "emt",
			"exists", "rtmin+0", "rtmin+1", "rtmin+2", "rtmin+3", "rtmin+4",
			"rtmin+5", "rtmin+6", "rtmin+7", "rtmin+8", "rtmin+9", "rtmin+10",
			"rtmin+11", "rtmin+12", "rtmin+13", "rtmin+14", "rtmin+15",
			"rtmin+16", "rtmin+17", "rtmin+18", "rtmin+19", "rtmin+20",
			"rtmin+21", "rtmin+22", "rtmin+23", "rtmin+24", "rtmin+25",
			"rtmin+26", "rtmin+27", "rtmin+28", "rtmin+29", "rtmin+30",
			"rtmin+31", "rtmin+32",
		},
	}
}

type Signal struct {
	RuleBase
	Qualifier
	Access []string
	Set    []string
	Peer   string
}

func newSignalFromLog(log map[string]string) Rule {
	return &Signal{
		RuleBase:  newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Access:    Must(toAccess(SIGNAL, log["requested_mask"])),
		Set:       []string{log["signal"]},
		Peer:      log["peer"],
	}
}

func (r *Signal) Validate() error {
	if err := validateValues(r.Kind(), "access", r.Access); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	if err := validateValues(r.Kind(), "set", r.Set); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *Signal) Less(other any) bool {
	o, _ := other.(*Signal)
	if len(r.Access) != len(o.Access) {
		return len(r.Access) < len(o.Access)
	}
	if len(r.Set) != len(o.Set) {
		return len(r.Set) < len(o.Set)
	}
	if r.Peer != o.Peer {
		return r.Peer < o.Peer
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Signal) Equals(other any) bool {
	o, _ := other.(*Signal)
	return slices.Equal(r.Access, o.Access) && slices.Equal(r.Set, o.Set) &&
		r.Peer == o.Peer && r.Qualifier.Equals(o.Qualifier)
}

func (r *Signal) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Signal) Constraint() constraint {
	return blockKind
}

func (r *Signal) Kind() Kind {
	return SIGNAL
}
