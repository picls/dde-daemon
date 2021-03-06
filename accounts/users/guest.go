/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package users

import (
	"fmt"
	"math/rand"
	"time"
)

func CreateGuestUser() (string, error) {
	shell, _ := getDefaultShell(defaultConfigShell)
	if len(shell) == 0 {
		shell = "/bin/bash"
	}

	username := getGuestUserName()
	var cmd = fmt.Sprintf("%s -m -d /tmp/%s -s %s -l -p %s %s",
		userCmdAdd, username, shell, EncodePasswd(""), username)
	err := doAction(cmd)
	if err != nil {
		return "", err
	}

	return username, nil
}

func getGuestUserName() string {
	var (
		seedStr = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

		l    = len(seedStr)
		name = "guest-"
	)

	for i := 0; i < 6; i++ {
		rand.Seed(time.Now().UnixNano())
		index := rand.Intn(l)
		name += string(seedStr[index])
	}

	return name
}
