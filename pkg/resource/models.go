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

type Source struct {
	Username         string `json:"username"`
	APIKey           string `json:"api_key"`
	AuthURL          string `json:"auth_url"`
	Domain           string `json:"domain"`
	Tenant           string `json:"tenant"`
	TenantID         string `json:"tenant_id"`
	Container        string `json:"container"`
	Regex            string `json:"regex"`
	DisableTLSVerify bool   `json:"disable_tls_verify"`
}

type Version struct {
	Path string `json:"path,omitempty"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CheckRequest struct {
	Resource Source  `json:"source"`
	Version  Version `json:"version"`
}

type InRequest struct {
	Resource Source  `json:"source"`
	Version  Version `json:"version"`
}

type InResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

type OutParams struct {
	From             string `json:"from"`
	SegmentContainer string `json:"segment_container"`
	SegmentSize      int64  `json:"segment_size"`
	DeleteAfter      int64  `json:"delete_after"`
}

type OutRequest struct {
	Resource Source    `json:"source"`
	Params   OutParams `json:"params"`
}

type OutResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}
