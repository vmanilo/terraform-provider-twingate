package resource

type Resource struct {
	ResourceName    string
	Name            string
	Address         string
	RemoteNetworkID string

	IsActive        *bool
	IsAuthoritative *bool
	Alias           *string
}
