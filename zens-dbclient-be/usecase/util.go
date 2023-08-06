package usecase

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (u *Usecase) initDbConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
