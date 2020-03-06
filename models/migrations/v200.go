package migrations

import (
	"xorm.io/xorm"
)

func addUserIdentityTable(x *xorm.Engine) error {

	type Identity struct {
		ID				int64	`xorm:"pk autoincr"`

		// combined unique index to allow empty IRI (= local
		// user) and keep the IRI still unique
		IRI				string	`xorm:"UNIQUE(identity)"`
		UserName		string	`xorm:"UNIQUE(identity) NOT NULL"`
		LowerUserName	string	`xorm:"UNIQUE(identity) NOT NULL"`
		DisplayName		string

		Type			int // originally UserType

		// Avatar
		Avatar			string	`xorm:"VARCHAR(2048) NOT NULL"`
		AvatarEmail		string	`xorm:"NOT NULL"`
		UseCustomAvatar	bool

		Visibility		int		`xorm:"NOT NULL DEFAULT 0"` // originally VisibleType
	}

	type User struct {

		IdentityId		int64 `xorm:"NOT NULL DEFAULT 0"`
		Identity		Identity `xorm:"-"`

		LowerName		string
		Name			string	`xorm:"-"`
		FullName		string	`xorm:"-"`

		Type			int		`xorm:"-"`

		// Avatar
		Avatar			string	`xorm:"-"`
		AvatarEmail		string	`xorm:"-"`
		UseCustomAvatar	bool	`xorm:"-"`

		Visibility		int		`xorm:"-"`
	}

	err := x.Sync2(new(Identity))
	if err != nil {
		return err
	}

	sess := x.NewSession()
	defer sess.Close()

	return dropTableColumns(sess, "user", "name", "full_name", "type", "avatar", "avatar_email", "use_custom_avatar", "visibility")

}
