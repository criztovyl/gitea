// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"strconv"
	"strings"

	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/modules/timeutil"
)

type Identity struct {
	ID           int64  `xorm:"pk autoincr"`

						// combined unique index to allow empty IRI (= local
						// user) and keep the IRI still unique
	IRI           string `xorm:"UNIQUE(identity)"`
	UserName      string `xorm:"UNIQUE(identity) NOT NULL"`
	LowerUserName string `xorm:"UNIQUE(identity) NOT NULL"`
	DisplayName   string

	Description   string
	Location      string
	Website       string

	Type        UserType

	// Avatar
	Avatar          string `xorm:"VARCHAR(2048) NOT NULL"`
	AvatarEmail     string `xorm:"NOT NULL"`
	UseCustomAvatar bool

	Visibility                structs.VisibleType `xorm:"NOT NULL DEFAULT 0"`

	CreatedUnix   timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix   timeutil.TimeStamp `xorm:"INDEX updated"`

	// Counters
	NumFollowers int
	NumFollowing int `xorm:"NOT NULL DEFAULT 0"`
	NumStars     int
	NumRepos     int

	// For organization
	NumTeams                  int
	NumMembers                int
	Teams                     []*Team             `xorm:"-"`
	Members                   IdentityList        `xorm:"-"`
	MembersIsPublic           map[int64]bool      `xorm:"-"`
	RepoAdminChangeTeamAccess bool                `xorm:"NOT NULL DEFAULT false"`


}

/*
Todo
- CustomAvatar
- Followers (Local+Remote)
- IsOrganization
- IsPublicMember
- GetOrganizationCount (<-> User@org.Members)
- Get{,Org,Mirror}Repositories (<-> Repository.Owner)
- Get{,Owned}Organizations
- GitName (?)
- ShortName (?)
*/

func (i *Identity) BeforeUpdate() {
	i.LowerUserName = strings.ToLower(i.UserName)
}

// IsUserPartOfOrg returns true if user with userID is part of the u organisation.
func (u *Identity) IsUserPartOfOrg(userID int64) bool {
	return u.isUserPartOfOrg(x, userID)
}

func (u *Identity) isUserPartOfOrg(e Engine, userID int64) bool {
	isMember, err := isOrganizationMember(e, u.ID, userID)
	if err != nil {
		log.Error("IsOrganizationMember: %v", err)
		return false
	}
	return isMember
}

// HomeLink returns the identity or organization home page link.
func (i *Identity) HomeLink() string {
	if len(i.IRI) == 0 { // IRI empty means local identity
		return setting.AppSubURL + "/" + i.UserName
	} else {
		return i.IRI
	}
}

// HTMLURL returns the user or organization's full link.
func (i *Identity) HTMLURL() string {
	if len(i.IRI) == 0 { // IRI empty means local identity
		return setting.AppURL + i.UserName
	} else {
		return i.IRI
	}
}

// SizedRelAvatarLink returns a link to the identity's avatar via
// the local explore page. Function returns immediately.
// When applicable, the link is for an avatar of the indicated size (in pixels).
func (i *Identity) SizedRelAvatarLink(size int) string {
	// i.UserName won't work for federated identities.
	return strings.TrimRight(setting.AppSubURL, "/") + "/user/avatar/" + i.UserName + "/" + strconv.Itoa(size)
}

// RelAvatarLink returns a relative link to the identity's avatar. The link
// may either be a sub-URL to this site, or a full URL to an external avatar
// service.
func (i *Identity) RelAvatarLink() string {
	return i.SizedRelAvatarLink(base.DefaultAvatarSize)
}

// AvatarLink returns identity avatar absolute link.
func (i *Identity) AvatarLink() string {
	link := i.RelAvatarLink()
	if link[0] == '/' && link[1] != '/' {
		return setting.AppURL + strings.TrimPrefix(link, setting.AppSubURL)[1:]
	}
	return link
}

/*

// GetUsersByIDs returns all resolved users from a list of Ids.
func GetUsersByIDs(ids []int64) ([]*User, error) {
	ous := make([]*User, 0, len(ids))
	if len(ids) == 0 {
		return ous, nil
	}
	err := x.In("id", ids).
		Asc("name").
		Find(&ous)
	return ous, err
}
*/

func GetIdentitiesByIDs(ids []int64) ([]*Identity, error) {
	idts := make([]*Identity, 0, len(ids))
	if len(ids) == 0 {
		return idts, nil
	}
	err := x.In("id", ids).
		Asc("user_name").
		Find(&idts)
	return idts, err
}

/*
func (repo *Repository) getOwner(e Engine) (err error) {
	if repo.Owner != nil {
		return nil
	}

	repo.Owner, err = getUserByID(e, repo.OwnerID)
	return err
}

// GetOwner returns the repository owner
func (repo *Repository) GetOwner() error {
	return repo.getOwner(x)
}

func (repo *Repository) mustOwner(e Engine) *User {
	if err := repo.getOwner(e); err != nil {
		return &User{
			Name:     "error",
			FullName: err.Error(),
		}
	}

	return repo.Owner
}

*/
