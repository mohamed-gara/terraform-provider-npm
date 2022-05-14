package client

type PackageMetadata struct {
	Name        string                             `json:"name"`
	Description string                             `json:"description"`
	Versions    map[string]*PackageVersionMetadata `json:"versions"`
}

type PackageVersionMetadata struct {
	Name    string `json:"name"`
	Version string `json:"version"`

	Dist struct {
		Integrity string `json:"integrity"`
		Shasum    string `json:"shasum"`
		Tarball   string `json:"tarball"`
	} `json:"dist"`
}
