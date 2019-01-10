package migrations

import (
	"fmt"
	"github.com/LiamDotPro/Go-Multitenancy/TenantUsers"
	"github.com/jinzhu/gorm"
)

// Attempts to migrate tables using database connection
func MigrateTenantTables(connection *gorm.DB) error {
	fmt.Println("Attempting to migrate tables to new database.")

	if err := connection.AutoMigrate(&TenantUsers.User{}).Error; err != nil {
		return err
	}

	return nil
}
