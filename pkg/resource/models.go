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
	From string `json:"from"`
	SegmentContainer string `json:"SegmentContainer"`
}

type OutRequest struct {
	Resource Source    `json:"source"`
	Params   OutParams `json:"params"`
}

type OutResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}
type Headers map[string]string

type LargeObjectOpts struct {
	Container        string  // Name of container to place object
	ObjectName       string  // Name of object
	Flags            int     // Creation flags
	CheckHash        bool    // If set Check the hash
	Hash             string  // If set use this hash to check
	ContentType      string  // Content-Type of the object
	Headers          Headers // Additional headers to upload the object with
	ChunkSize        int64   // Size of chunks of the object, defaults to 10MB if not set
	MinChunkSize     int64   // Minimum chunk size, automatically set for SLO's based on info
	SegmentContainer string  // Name of the container to place segments
	SegmentPrefix    string  // Prefix to use for the segments
	NoBuffer         bool    // Prevents using a bufio.Writer to write segments
}
