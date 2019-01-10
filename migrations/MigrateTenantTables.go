package migrations

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

// Attempts to migrate tables using database connection
func migrateTenantTables(connection *gorm.DB) error {
	fmt.Println("Attempting to migrate tables to new database.")

	if err := connection.AutoMigrate(&User{}).Error; err != nil {
		return err
	}

	return nil
}
