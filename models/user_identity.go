// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"strings"

	"code.gitea.io/gitea/modules/structs"
)

type Identity struct {
	ID           int64  `xorm:"pk autoincr"`

						// combined unique index to allow empty IRI (= local
						// user) and keep the IRI still unique
	IRI           string `xorm:"UNIQUE(identity)"`
	UserName      string `xorm:"UNIQUE(identity) NOT NULL"`
	LowerUserName string `xorm:"UNIQUE(identity) NOT NULL"`
	DisplayName   string

	Type        UserType

	// Avatar
	Avatar          string `xorm:"VARCHAR(2048) NOT NULL"`
	AvatarEmail     string `xorm:"NOT NULL"`
	UseCustomAvatar bool

	Visibility                structs.VisibleType `xorm:"NOT NULL DEFAULT 0"`
}

func (i *Identity) BeforeUpdate() {
	i.LowerUserName = strings.ToLower(i.UserName)
}
