package main

import "github.com/jinzhu/gorm"

type TenantSubscriptionInformation struct {
	gorm.Model
	TenantId uint
	SubscriptionType uint // This is linked to the TenantSubscriptionType Table
}
