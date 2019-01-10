package migrations

import (
	"github.com/LiamDotPro/Go-Multitenancy/database"
	"github.com/LiamDotPro/Go-Multitenancy/tenants"
)

var Connection := database.Connection

/**
This method uses the base tenant connection set out within init.
*/
func migrateMasterTenantDatabase() error {

	if err := Connection.AutoMigrate(&tenants.TenantConnectionInformation{}).Error; err != nil {
		return err
	}

	if err := Connection.AutoMigrate(&tenants.TenantSubscriptionInformation{}).Error; err != nil {
		return err
	}

	if err := Connection.AutoMigrate(&tenants.TenantSubscriptionType{}).Error; err != nil {
		return err
	}

	if err := Connection.AutoMigrate(&tenants.MasterUsers{}).Error; err != nil {
		return err
	}

	return nil

}
