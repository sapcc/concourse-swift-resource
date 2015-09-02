package resource

import "github.com/ncw/swift"

func NewClient(source Source) *swift.Connection {
	c := swift.Connection{
		UserName: source.Username,
		ApiKey:   source.ApiKey,
		AuthUrl:  source.AuthUrl,
		Domain:   source.Domain, // Name of the domain (v3 auth only)
		//Tenant:   "tenant", // Name of the tenant (v2 auth only)
	}

	err := c.Authenticate()

	if err != nil {
		Fatal("Authentication failed", err)
	}

	return &c

}
