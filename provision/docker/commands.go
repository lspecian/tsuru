// Copyright 2014 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

import (
	"bytes"
	"fmt"
	"github.com/tsuru/config"
	"github.com/tsuru/tsuru/provision"
	"github.com/tsuru/tsuru/repository"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
)

// deployCmds returns the commands that is used when provisioner
// deploy an unit.
func deployCmds(app provision.App, version string) ([]string, error) {
	deployCmd, err := config.GetString("docker:deploy-cmd")
	if err != nil {
		return nil, err
	}
	appRepo := repository.ReadOnlyURL(app.GetName())
	var envs string
	for _, env := range app.Envs() {
		envs += fmt.Sprintf(`%s=%s `, env.Name, strings.Replace(env.Value, " ", "", -1))
	}
	cmds := []string{deployCmd, appRepo, version, envs}
	return cmds, nil
}

// runWithAgentCmds returns the commands that should be passed when the
// provisioner will run an unit using tsuru_unit_agent to start.
func runWithAgentCmds(app provision.App) ([]string, error) {
	runCmd, err := config.GetString("docker:run-cmd:bin")
	if err != nil {
		return nil, err
	}
	ssh, err := sshCmds()
	if err != nil {
		return nil, err
	}
	host := app.Envs()["TSURU_HOST"].Value
	token := app.Envs()["TSURU_APP_TOKEN"].Value
	unitAgentCmds := []string{"tsuru_unit_agent", host, token, app.GetName(), runCmd}
	unitAgentCmd := strings.Join(unitAgentCmds, " ")
	fallbackCmd := fmt.Sprintf("(%s || %s)", unitAgentCmd, runCmd)
	sshCmd := strings.Join(ssh, " && ")
	cmd := fmt.Sprintf("%s && %s", fallbackCmd, sshCmd)
	cmds := []string{"/bin/bash", "-c", cmd}
	return cmds, nil
}

// sshCmds returns the commands needed to start a ssh daemon.
func sshCmds() ([]string, error) {
	addKeyCommand, err := config.GetString("docker:ssh:add-key-cmd")
	if err != nil {
		return nil, err
	}
	keyFile, err := config.GetString("docker:ssh:public-key")
	if err != nil {
		if u, err := user.Current(); err == nil {
			keyFile = path.Join(u.HomeDir, ".ssh", "id_rsa.pub")
		} else {
			keyFile = os.ExpandEnv("${HOME}/.ssh/id_rsa.pub")
		}
	}
	f, err := filesystem().Open(keyFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	keyContent, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	sshdCommand, err := config.GetString("docker:ssh:sshd-path")
	if err != nil {
		sshdCommand = "sudo /usr/sbin/sshd"
	}
	return []string{
		fmt.Sprintf("%s %s", addKeyCommand, bytes.TrimSpace(keyContent)),
		sshdCommand + " -D",
	}, nil
}
