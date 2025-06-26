// SPDX-FileCopyrightText: 2015-2023 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

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
