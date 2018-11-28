package main

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"strings"
)

type TenantConnectionInformation struct {
	gorm.Model
	TenantId                  uint `gorm:"AUTO_INCREMENT"`
	TenantSubDomainIdentifier string
	ConnectionString          string
}

// Helper method that create's and returns the database connection.
func (t TenantConnectionInformation) getConnection() (connection *gorm.DB, err error) {

	if len(strings.TrimSpace(t.ConnectionString)) == 0 {
		return nil, errors.New("Connection string was not found or was empty..")
	}

	db, err := gorm.Open("postgres", t.ConnectionString)

	if err == nil {
		return nil, err
	}

	return db, nil
}
