package entity

import (
	"net/http"

	"gorm.io/gorm"
)

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

type ConnectionNotFoundError struct{}

func (e ConnectionNotFoundError) Error() string {
	return "connection not found"
}

func (e ConnectionNotFoundError) ToHttpError() (int, HttpResponseError) {
	return http.StatusNotFound, HttpResponseError{
		Message: e.Error(),
	}
}
