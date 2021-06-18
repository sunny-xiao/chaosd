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

package chaosd

import (
	"os/exec"

	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"go.uber.org/zap"

	"github.com/chaos-mesh/chaosd/pkg/core"
)

type fileAttack struct{}

var FileAttack AttackType = fileAttack{}

func (fileAttack) Attack(options core.AttackConfig, env Environment) (err error) {
	attack := options.(*core.FileCommand)

	switch attack.Action {
	case core.FileCreateAction:
		if err = env.Chaos.createFile(attack, env.AttackUid); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (s *Server) createFile(attack *core.FileCommand, uid string) error {

	createCmd := ""
	if len(attack.FileName) > 0 {
		createCmd = "touch " + attack.DestDir + attack.FileName
	} else if len(attack.DirName) > 0 {
		createCmd = "mkdir " + attack.DestDir + attack.DirName
	}

	cmd := exec.Command("/bin/bash", "-c", createCmd) // #nosec
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(cmd.String()+string(stdout), zap.Error(err))
		return errors.WithStack(err)
	}

	return nil
}

func (fileAttack) Recover(exp core.Experiment, env Environment) error {
	config, err := exp.GetRequestCommand()
	if err != nil {
		return err
	}
	attack := config.(*core.FileCommand)

	switch attack.Action {
	case core.FileCreateAction:
		if err = env.Chaos.recoverCreateFile(attack); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (s *Server) recoverCreateFile(attack *core.FileCommand) error {

	deleteCmd := ""
	if len(attack.FileName) > 0 {
		deleteCmd = "rm -rf " + attack.DestDir + attack.FileName
	} else if len(attack.DirName) > 0 {
		deleteCmd = "rm -rf " + attack.DestDir + attack.DirName
	}

	cmd := exec.Command("/bin/bash", "-c", deleteCmd) // #nosec
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(cmd.String()+string(stdout), zap.Error(err))
		return errors.WithStack(err)
	}
	return nil
}
