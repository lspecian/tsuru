// Copyright 2013 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/globocom/tsuru/cmd"
)

type tokenCmd struct{}

func (c *tokenCmd) Run(context *cmd.Context, client *cmd.Client) error {
	return nil
}

func (tokenCmd) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "token",
		Usage:   "token",
		Desc:    "Generates a tsuru token.",
		MinArgs: 0,
	}
}