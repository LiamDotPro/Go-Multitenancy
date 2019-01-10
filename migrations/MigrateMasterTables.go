package migrations

/**
This method uses the base tenant connection set out within init.
*/
func migrateMasterTenantDatabase() error {

	if err := Connection.AutoMigrate(&TenantConnectionInformation{}).Error; err != nil {
		return err
	}

	if err := Connection.AutoMigrate(&TenantSubscriptionInformation{}).Error; err != nil {
		return err
	}

	if err := Connection.AutoMigrate(&TenantSubscriptionType{}).Error; err != nil {
		return err
	}

	if err := Connection.AutoMigrate(&MasterUsers{}).Error; err != nil {
		return err
	}

	return nil

}
