package main

import "github.com/jinzhu/gorm"

type TenantSubscriptionType struct {
	gorm.Model
	SubscriptionName    string
	SubscriptionPrice   uint
	SubscriptionPeriod  uint // renewal period denoted as 1-24
	SubscriptionRenewal bool
}
