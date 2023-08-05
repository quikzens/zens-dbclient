package entity

import "gorm.io/gorm"

type Connection struct {
	Id         int
	Client     *gorm.DB
	Credential Credential
}

type Credential struct {
	Host         string
	Port         string
	DatabaseName string
	User         string
	Password     string
}
