package tidb

import (
	"gorm.io/gorm"
)

type connectionStrategy func() (*gorm.DB, error)
