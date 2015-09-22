package resource

type Source struct {
	Username string `json:"username"`
	ApiKey   string `json:"api_key"`
	AuthUrl  string `json:"auth_url"`
	Domain   string `json:"domain"`
	Tenant   string `json:"tenant"`
	TenantId string `json:"tenant_id"`

	Container string `json:"container"`
	Regex     string `json:"regex"`
}

type Version struct {
	Path string `json:"path,omitempty"`
}
