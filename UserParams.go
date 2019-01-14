package main

/**
Model binding types.
*/

type CreateUserParams struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Type     int    `form:"type" json:"type"`
}

type UpdateUserParams struct {
	Id            uint   `form:"id" json:"id" binding:"required"`
	Email         string `form:"email" json:"email"`
	AccountType   int    `form:"accountType" json:"accountType"`
	FirstName     string `form:"firstName" json:"firstName"`
	LastName      string `form:"lastName" json:"lastName"`
	PhoneNumber   string `form:"phoneNumber" json:"phoneNumber"`
	RecoveryEmail string `form:"recoveryEmail" json:"recoveryEmail"`
}

type LoginParams struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type DeleteUserParams struct {
	Id uint `form:"id" json:"id" binding:"required"`
}