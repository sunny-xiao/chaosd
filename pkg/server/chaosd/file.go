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
	"fmt"
	"github.com/chaos-mesh/chaosd/pkg/core"
	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"go.uber.org/zap"
	"os"
	"os/exec"
)

type fileAttack struct{}

var FileAttack AttackType = fileAttack{}

var FileMode string

func (fileAttack) Attack(options core.AttackConfig, env Environment) (err error) {
	attack := options.(*core.FileCommand)

	switch attack.Action {
	case core.FileCreateAction:
		if err = env.Chaos.createFile(attack, env.AttackUid); err != nil {
			return errors.WithStack(err)
		}
	case core.FileModifyPrivilegeAction:
		if err = env.Chaos.modifyFilePrivilege(attack, env.AttackUid); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (s *Server) createFile(attack *core.FileCommand, uid string) error {

	var err error
	if len(attack.FileName) > 0 {
		_, err = os.Create(attack.DestDir + attack.FileName)
	} else if len(attack.DirName) > 0 {
		err = os.Mkdir(attack.DestDir+attack.DirName, os.ModePerm)
	}

	if err != nil {
		log.Error("create file/dir faild", zap.Error(err))
		return errors.WithStack(err)
	}

	return nil
}

func (s *Server) modifyFilePrivilege(attack *core.FileCommand, uid string) error {

	cmdStr := "stat -c %a" + " "+ attack.FileName
	cmd := exec.Command("bash", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(output), zap.Error(err))
		return errors.WithStack(err)
	}
	FileMode = string(output)

	cmdStr = fmt.Sprintf("chmod %d %s", attack.Privilege, attack.FileName)

	cmd = exec.Command("bash", "-c", cmdStr)
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Error(string(output), zap.Error(err))
		return errors.WithStack(err)
	}
	log.Info(string(output))

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
	case core.FileModifyPrivilegeAction:
		if err = env.Chaos.recoverModifyPrivilege(attack); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (s *Server) recoverCreateFile(attack *core.FileCommand) error {

	var err error
	if len(attack.FileName) > 0 {
		err = os.Remove(attack.DestDir + attack.FileName)
	} else if len(attack.DirName) > 0 {
		err = os.RemoveAll(attack.DestDir + attack.DirName)
	}

	if err != nil {
		log.Error("delete file/dir faild", zap.Error(err))
		return errors.WithStack(err)
	}
	return nil
}

func (s *Server) recoverModifyPrivilege(attack *core.FileCommand) error {

	cmdStr := fmt.Sprintf("chmod %s %s", FileMode, attack.FileName)
	cmd := exec.Command("bash", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(output), zap.Error(err))
		return errors.WithStack(err)
	}
	return nil
}