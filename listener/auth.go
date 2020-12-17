package listener

type AuthHeader struct {
	Identity     Identity     `json:"identity"`
	Entitlements Entitlements `json:"entitlements"`
}

type Identity struct {
	AccountNumber string `json:"account_number"`
	Type          string `json:"type"`
	User          User   `json:"user"`
}

type User struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	IsOrgAdmin bool   `json:"is_org_admin"`
}

type Entitlements struct {
	CostManagement CostMgmt `json:"cost_management"`
}

type CostMgmt struct {
	IsEntitled bool `json:"is_entitled"`
}
