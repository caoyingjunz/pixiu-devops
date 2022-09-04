/*
Copyright 2021 The Pixiu Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package db

import (
	"github.com/caoyingjunz/gopixiu/pkg/db/sys"
	"github.com/casbin/casbin/v2"

	"gorm.io/gorm"

	"github.com/caoyingjunz/gopixiu/pkg/db/cloud"
	"github.com/caoyingjunz/gopixiu/pkg/db/demo"
	"github.com/caoyingjunz/gopixiu/pkg/db/user"
)

type ShareDaoFactory interface {
	User() user.UserInterface
	Demo() demo.DemoInterface
	Cloud() cloud.CloudInterface
	Role() sys.RoleInterface
	Menu() sys.MenuInterface
	Casbin() sys.CasbinInterface
}

type shareDaoFactory struct {
	db       *gorm.DB
	enforcer *casbin.Enforcer
}

func (f *shareDaoFactory) Demo() demo.DemoInterface {
	return demo.NewDemo(f.db)
}

func (f *shareDaoFactory) Cloud() cloud.CloudInterface {
	return cloud.NewCloud(f.db)
}

func (f *shareDaoFactory) User() user.UserInterface {
	return user.NewUser(f.db)
}
func (f *shareDaoFactory) Role() sys.RoleInterface {
	return sys.NewRole(f.db)
}
func (f *shareDaoFactory) Menu() sys.MenuInterface {
	return sys.NewMenu(f.db)
}

func (f *shareDaoFactory) Casbin() sys.CasbinInterface {
	return sys.NewCasbin(f.db, f.enforcer)
}

func NewDaoFactory(db *gorm.DB, enforcer *casbin.Enforcer) ShareDaoFactory {
	return &shareDaoFactory{
		db:       db,
		enforcer: enforcer,
	}
}
