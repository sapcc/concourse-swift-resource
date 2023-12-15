/*******************************************************************************
*
* Copyright 2015-2023 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package resource

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/ncw/swift/v2"
)

var tokenCacheFile = "/tmp/token.cache" //nolint:gosec // false positive
var cacheToken = false

func NewClient(ctx context.Context, source Source) *swift.Connection {
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConnsPerHost: 2048,
	}
	if source.DisableTLSVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // debugging option
	}
	c := swift.Connection{
		UserName:  source.Username,
		ApiKey:    source.APIKey,
		AuthUrl:   source.AuthURL,
		Domain:    source.Domain,   // Name of the domain (v3 auth only)
		Tenant:    source.Tenant,   // Name of the tenant
		TenantId:  source.TenantID, // Id of the tenant
		Retries:   1,
		UserAgent: fmt.Sprintf("%s (concourse swift resource; %s; container: %s)", swift.DefaultUserAgent, path.Base(os.Args[0]), source.Container),
		Transport: transport,
	}

	if err := c.Authenticate(ctx); err != nil {
		Fatal("Authentication failed", err)
	}

	if !cacheToken {
		return &c
	}

	if _, err := os.Stat(tokenCacheFile); err == nil {
		if cachedJSON, err := os.ReadFile(tokenCacheFile); err == nil {
			var cc swift.Connection
			if err := json.Unmarshal(cachedJSON, &cc); err == nil {
				if c.UserName == cc.UserName &&
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
	clientJSON, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		Fatal("Failed to marshal swift client", err)
	}
	if err := os.WriteFile(tokenCacheFile, clientJSON, 0600); err != nil {
		Sayf("Failed to cache token to %s: %s", tokenCacheFile, err)
		os.Exit(1)
	}
}
