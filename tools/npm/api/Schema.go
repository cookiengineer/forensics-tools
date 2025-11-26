package api

type Schema struct {
	Identifier      string            `json:"_id"`
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Author          *Author           `json:"author,omitempty"`
	Maintainers     []Maintainer      `json:"maintainers"`
	Homepage        string            `json:"homepage,omitempty"`
	Bugs            *Bugs             `json:"bugs,omitempty"`
	Dist            *Dist             `json:"dist,omitempty"`
	Engines         map[string]string `json:"engines,omitempty"`
	Exports         map[string]Export `json:"exports,omitempty"`
	Scripts         map[string]string `json:"scripts,omitempty"`
	Description     string            `json:"description,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	DevDependencies map[string]string `json:"devDependencies,omitempty"`

	// XXX: Uninteresting properties
	// Directories     map[string]any    `json:"directories,omitempty"`
	// NPMVersion      string            `json:"_npmVersion,omitempty"`
	// NodeVersion     string            `json:"_nodeVersion,omitempty"`
	// HasShrinkWrap   string            `json:"_hasShrinkWrap,omitempty"`
	// NPMOperationalInternal map[string]any `json:"_npmOperationalInternal,omitempty"`
}

type Author struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type Maintainer struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type Bugs struct {
	URL string `json:"url,omitempty"`
}

type Dist struct {
	SHASum       string      `json:"shasum,omitempty"`
	Tarball      string      `json:"tarball,omitempty"`
	FileCount    int         `json:"fileCount,omitempty"`
	Integrity    string      `json:"integrity,omitempty"`
	Signatures   []Signature `json:"signatures,omitempty"`
	UnpackedSize int         `json:"unpackedSize,omitempty"`
}

type Signature struct {
	Sig   string `json:"sig,omitempty"`
	KeyID string `json:"keyid,omitempty"`
}

type Export struct {
	Types   string `json:"types,omitempty"`
	Import  string `json:"import,omitempty"`
	Require string `json:"require,omitempty"`
}
