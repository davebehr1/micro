package user

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

func (model *User) BeforeCreate(scope *gorm.Scope) error {
	uuid, _ := uuid.NewV4()
	return scope.SetColumn("Id", uuid.String())
}
