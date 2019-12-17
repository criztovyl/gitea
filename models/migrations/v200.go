package migrations

import (
	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/structs"

	"xorm.io/xorm"
)

func addUserIdentityTable(x *xorm.Engine) error {
	type User struct {

		IdentityId   int64 `xorm:"NOT NULL DEFAULT 0"`
		Identity     models.Identity

		LowerName string `xorm:"-"`
		Name      string `xorm:"-"`
		FullName  string `xorm:"-"`

		Type        models.UserType      `xorm:"-"`

		// Avatar
		Avatar          string `xorm:"-"`
		AvatarEmail     string `xorm:"-"`
		UseCustomAvatar bool   `xorm:"-"`

		Visibility                structs.VisibleType `xorm:"-"`
	}

	err := x.Sync2(new(models.Identity))
	if err != nil {
		return err
	}

	err = x.Sync2(new(User))
	return err
}
