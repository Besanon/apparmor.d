// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
)

const (
	MOUNT   Kind = "mount"
	REMOUNT Kind = "remount"
	UMOUNT  Kind = "umount"
)

func init() {
	requirements[MOUNT] = requirement{
		"flags": {
			"acl", "async", "atime", "ro", "rw", "bind", "rbind", "dev",
			"diratime", "dirsync", "exec", "iversion", "loud", "mand", "move",
			"noacl", "noatime", "nodev", "nodiratime", "noexec", "noiversion",
			"nomand", "norelatime", "nosuid", "nouser", "private", "relatime",
			"remount", "rprivate", "rshared", "rslave", "runbindable", "shared",
			"silent", "slave", "strictatime", "suid", "sync", "unbindable",
			"user", "verbose",
		},
	}
}

type MountConditions struct {
	FsType  string
	Options []string
}

func newMountConditions(rule rule) (MountConditions, error) {
	options, err := toValues(MOUNT, "flags", rule.GetValuesAsString("options"))
	if err != nil {
		return MountConditions{}, err
	}
	return MountConditions{
		FsType:  rule.GetValuesAsString("fstype"),
		Options: options,
	}, nil
}

func newMountConditionsFromLog(log map[string]string) MountConditions {
	if _, present := log["flags"]; present {
		return MountConditions{
			FsType:  log["fstype"],
			Options: Must(toValues(MOUNT, "flags", log["flags"])),
		}
	}
	return MountConditions{FsType: log["fstype"]}
}

func (m MountConditions) Validate() error {
	return validateValues(MOUNT, "flags", m.Options)
}

func (m MountConditions) Compare(other MountConditions) int {
	if res := compare(m.FsType, other.FsType); res != 0 {
		return res
	}
	return compare(m.Options, other.Options)
}

func (m *MountConditions) Merge(other MountConditions) bool {
	if m.FsType == other.FsType {
		m.Options = merge(MOUNT, "flags", m.Options, other.Options)
		return true
	}
	return false
}

type Mount struct {
	Base
	Qualifier
	MountConditions
	Source     string
	MountPoint string
}

func newMount(q Qualifier, rule rule) (Rule, error) {
	mount, src := "", ""
	r := rule.GetSlice()
	if len(r) > 0 {
		if r[0] != tokARROW {
			src = r[0]
			r = r[1:]
		}
		if len(r) == 2 {
			if r[0] != tokARROW {
				return nil, fmt.Errorf("missing '%s' in rule: %s", tokARROW, rule)
			}
			mount = r[1]
		}
	}

	conditions, err := newMountConditions(rule)
	if err != nil {
		return nil, err
	}
	return &Mount{
		Base:            newBase(rule),
		Qualifier:       q,
		MountConditions: conditions,
		Source:          src,
		MountPoint:      mount,
	}, nil
}

func newMountFromLog(log map[string]string) Rule {
	return &Mount{
		Base:            newBaseFromLog(log),
		Qualifier:       newQualifierFromLog(log),
		MountConditions: newMountConditionsFromLog(log),
		Source:          log["srcname"],
		MountPoint:      log["name"],
	}
}

func (r *Mount) Validate() error {
	if err := r.MountConditions.Validate(); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *Mount) Compare(other Rule) int {
	o, _ := other.(*Mount)
	if res := compare(r.Source, o.Source); res != 0 {
		return res
	}
	if res := compare(r.MountPoint, o.MountPoint); res != 0 {
		return res
	}
	if res := r.MountConditions.Compare(o.MountConditions); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *Mount) Merge(other Rule) bool {
	o, _ := other.(*Mount)
	mc := &r.MountConditions

	if !r.Qualifier.Equal(o.Qualifier) {
		return false
	}
	if r.Source == o.Source && r.MountPoint == o.MountPoint &&
		mc.Merge(o.MountConditions) {
		b := &r.Base
		return b.merge(o.Base)
	}
	return false
}

func (r *Mount) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Mount) Constraint() constraint {
	return blockKind
}

func (r *Mount) Kind() Kind {
	return MOUNT
}

type Umount struct {
	Base
	Qualifier
	MountConditions
	MountPoint string
}

func newUmount(q Qualifier, rule rule) (Rule, error) {
	mount := ""
	r := rule.GetSlice()
	if len(r) > 0 {
		mount = r[0]
	}
	conditions, err := newMountConditions(rule)
	if err != nil {
		return nil, err
	}
	return &Umount{
		Base:            newBase(rule),
		Qualifier:       q,
		MountConditions: conditions,
		MountPoint:      mount,
	}, nil
}

func newUmountFromLog(log map[string]string) Rule {
	return &Umount{
		Base:            newBaseFromLog(log),
		Qualifier:       newQualifierFromLog(log),
		MountConditions: newMountConditionsFromLog(log),
		MountPoint:      log["name"],
	}
}

func (r *Umount) Validate() error {
	if err := r.MountConditions.Validate(); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *Umount) Compare(other Rule) int {
	o, _ := other.(*Umount)
	if res := compare(r.MountPoint, o.MountPoint); res != 0 {
		return res
	}
	if res := r.MountConditions.Compare(o.MountConditions); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *Umount) Merge(other Rule) bool {
	o, _ := other.(*Umount)
	mc := &r.MountConditions

	if !r.Qualifier.Equal(o.Qualifier) {
		return false
	}
	if r.MountPoint == o.MountPoint && mc.Merge(o.MountConditions) {
		b := &r.Base
		return b.merge(o.Base)
	}
	return false
}

func (r *Umount) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Umount) Constraint() constraint {
	return blockKind
}

func (r *Umount) Kind() Kind {
	return UMOUNT
}

type Remount struct {
	Base
	Qualifier
	MountConditions
	MountPoint string
}

func newRemount(q Qualifier, rule rule) (Rule, error) {
	mount := ""
	r := rule.GetSlice()
	if len(r) > 0 {
		mount = r[0]
	}

	conditions, err := newMountConditions(rule)
	if err != nil {
		return nil, err
	}
	return &Remount{
		Base:            newBase(rule),
		Qualifier:       q,
		MountConditions: conditions,
		MountPoint:      mount,
	}, nil
}

func newRemountFromLog(log map[string]string) Rule {
	return &Remount{
		Base:            newBaseFromLog(log),
		Qualifier:       newQualifierFromLog(log),
		MountConditions: newMountConditionsFromLog(log),
		MountPoint:      log["name"],
	}
}

func (r *Remount) Validate() error {
	if err := r.MountConditions.Validate(); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *Remount) Compare(other Rule) int {
	o, _ := other.(*Remount)
	if res := compare(r.MountPoint, o.MountPoint); res != 0 {
		return res
	}
	if res := r.MountConditions.Compare(o.MountConditions); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *Remount) Merge(other Rule) bool {
	o, _ := other.(*Remount)
	mc := &r.MountConditions

	if !r.Qualifier.Equal(o.Qualifier) {
		return false
	}
	if r.MountPoint == o.MountPoint && mc.Merge(o.MountConditions) {
		b := &r.Base
		return b.merge(o.Base)
	}
	return false
}

func (r *Remount) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Remount) Constraint() constraint {
	return blockKind
}

func (r *Remount) Kind() Kind {
	return REMOUNT
}
