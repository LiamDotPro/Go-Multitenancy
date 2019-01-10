package main

import "github.com/jinzhu/gorm"

type MasterUsers struct {
	gorm.Model
	Email         string
	Password      string `json:",omitempty"`
	AccountType   int
	FirstName     string
	LastName      string
	PhoneNumber   string
	RecoveryEmail string
}

