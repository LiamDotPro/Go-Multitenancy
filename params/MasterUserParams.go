package params

type CreateNewTenantParams struct {
	SubDomainIdentifier string `form:"subDomainIdentifier" json:"subDomainIdentifier" binding:"required"`
}
