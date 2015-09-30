package resource

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/ncw/swift"
)

var tokenCacheFile = "/tmp/token.cache"

func NewClient(source Source) *swift.Connection {
	c := swift.Connection{
		UserName: source.Username,
		ApiKey:   source.ApiKey,
		AuthUrl:  source.AuthUrl,
		Domain:   source.Domain,   // Name of the domain (v3 auth only)
		Tenant:   source.Tenant,   // Name of the tenant
		TenantId: source.TenantId, // Id of the tenant
		Retries:  1,
	}

	if _, err := os.Stat(tokenCacheFile); err == nil {
		if cachedJson, err := ioutil.ReadFile(tokenCacheFile); err == nil {
			var cc swift.Connection
			if err := json.Unmarshal(cachedJson, &cc); err == nil {
				if cc.UserName == cc.UserName &&
					c.ApiKey == cc.ApiKey &&
					c.AuthUrl == cc.AuthUrl &&
					c.Domain == cc.Domain &&
					c.Tenant == cc.Tenant &&
					c.TenantId == cc.TenantId {
					c.AuthToken = cc.AuthToken
					c.StorageUrl = cc.StorageUrl
				}
			} else {
				Fatal("Failed to unmarshal cached token: ", err)
			}
		} else {
			Fatal("Failed to read cached token", err)
		}
	}

	return &c
}

func CacheClientToken(c *swift.Connection) {
	if !c.Authenticated() {
		return
	}
	clientJson, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		Fatal("Failed to marshal swift client", err)
	}
	if err := ioutil.WriteFile(tokenCacheFile, clientJson, 0600); err != nil {
		Sayf("Failed to cache token to %s: %s", tokenCacheFile, err)
		os.Exit(1)
	}
}
