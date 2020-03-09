// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

type IdentityList []*Identity

func (identities IdentityList) getIdentityIDs() []int64 {
	idtyIDs := make([]int64, len(identities))
	for _, idty := range identities {
		idtyIDs = append(idtyIDs, idty.ID) //Considering that user id are unique in the list
	}
	return idtyIDs
}
