package migrations

import (
	"github.com/LiamDotPro/Go-Multitenancy/database"
	"github.com/LiamDotPro/Go-Multitenancy/master"
)

var Connection = database.Connection

/**
This method uses the base tenant connection set out within init.
*/
func migrateMasterTenantDatabase() error {

	if err := Connection.AutoMigrate(&master.TenantConnectionInformation{}).Error; err != nil {
		return err
	}

	if err := Connection.AutoMigrate(&master.TenantSubscriptionInformation{}).Error; err != nil {
		return err
	}

	if err := Connection.AutoMigrate(&master.TenantSubscriptionType{}).Error; err != nil {
		return err
	}

	if err := Connection.AutoMigrate(&MasterUsers{}).Error; err != nil {
		return err
	}

	return nil

}
