// Copyright 2020 Chaos Mesh Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"encoding/json"
	"github.com/pingcap/errors"
)

type FileCommand struct {
	CommonAttackConfig

	FileName string
	DirName  string
	DestDir  string
	Privilege uint32
}

var _ AttackConfig = &FileCommand{}

const (
	FileCreateAction = "create"
	FileModifyPrivilegeAction = "modify"
)

func (n *FileCommand) Validate() error {
	if err := n.CommonAttackConfig.Validate(); err != nil {
		return err
	}
	switch n.Action {
	case FileCreateAction:
		return n.validFileCreate()
	case FileModifyPrivilegeAction:
        return n.validFileModify()
	default:
		return errors.Errorf("network action %s not supported", n.Action)
	}
}

func (n *FileCommand) validFileCreate() error {
	return nil
}

func (n *FileCommand) validFileModify() error {
	if len(n.FileName) == 0 && len(n.DirName) == 0 {
		return errors.New("filename and dirname can not all null")
	}

	if n.Privilege == 0 {
		return errors.New("file privilege can not null")
	}

	return nil
}

func (n *FileCommand) CompleteDefaults() {
	switch n.Action {
	case FileCreateAction:
		n.setDefaultForFileCreate()
	}
}

func (n *FileCommand) setDefaultForFileCreate() {
	if len(n.FileName) == 0 && len(n.DirName) == 0 {
		n.FileName = "chaosd.file"
	}
	if len(n.DestDir) > 0 {
		n.DestDir = n.DestDir + "/"
	}
}

func (n FileCommand) RecoverData() string {
	data, _ := json.Marshal(n)

	return string(data)
}

func NewFileCommand() *FileCommand {
	return &FileCommand{
		CommonAttackConfig: CommonAttackConfig{
			Kind: FileAttack,
		},
	}
}
