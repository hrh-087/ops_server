package global

import (
	"gorm.io/gorm"
	"time"
)

type OpsModel struct {
	ID        uint `gorm:"primary_key;auto_increment;not_null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
