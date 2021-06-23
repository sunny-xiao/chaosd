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
	"os"
	"os/exec"
	"strconv"

	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"go.uber.org/zap"

	"github.com/chaos-mesh/chaosd/pkg/core"
)

type fileAttack struct{}

var FileAttack AttackType = fileAttack{}

var FileMode int

var DirMode int

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

func (s *Server) modifyFilePrivilege(attack *core.FileCommand, uid string) error {

	//pri, _ := fmt.Printf("%04d", attack.Privilege)

	if len(attack.FileName) != 0 {
		/*var err error
        FileMode, err = getFileMode(attack.FileName)
        if err != nil {
        	return errors.WithStack(err)
		}*/

		t := os.FileMode(attack.Privilege)
		fmt.Println(t)

		if err := os.Chmod(attack.FileName, t); err != nil {
			return errors.WithStack(err)
		}
	}else if len(attack.DirName) != 0 {
		/*var err error
		DirMode, err = getFileMode(attack.DirName)
		if err != nil {
			return errors.WithStack(err)
		}*/

		if err := os.Chmod(attack.DirName, os.FileMode(attack.Privilege)); err != nil {
			return errors.WithStack(err)
		}
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
	case core.FileModifyPrivilegeAction:
		if err = env.Chaos.recoverModifyPrivilege(attack); err != nil {
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

func (s *Server) recoverModifyPrivilege(attack *core.FileCommand) error {

	if len(attack.FileName) != 0 {
		if err := os.Chmod(attack.FileName, os.FileMode(FileMode)); err != nil {
			return errors.WithStack(err)
		}
	}else if len(attack.DirName) != 0 {
		if err := os.Chmod(attack.FileName, os.FileMode(DirMode)); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}


func getFileMode(path string) (int, error) {
	cmd := fmt.Sprintf("stat -c %a %s", path)
	exeCmd := exec.Command("/bin/bash", "-c", cmd)
	stdout, err := exeCmd.CombinedOutput()
	if err != nil {
		log.Error(exeCmd.String()+string(stdout), zap.Error(err))
		return 0, errors.WithStack(err)
	}

	s := string(stdout)
	d, _ := strconv.Atoi(s)
	return d, nil
}