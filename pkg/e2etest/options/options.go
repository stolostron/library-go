package options

type TestOptionsContainer struct {
	Options TestOptions `yaml:"options"`
}

// Define options available for Tests to consume
type TestOptions struct {
	HubConfigDir             string          `yaml:"hubConfigDir"`
	ManagedClustersConfigDir string          `yaml:"managedClustersConfigDir"`
	ImageRegistry            Registry        `yaml:"imageRegistry,omitempty"`
	IdentityProvider         int             `yaml:"identityProvider,omitempty"`
	Connection               CloudConnection `yaml:"cloudConnection,omitempty"`
	Headless                 string          `yaml:"headless,omitempty"`
	OwnerPrefix              string          `yaml:"ownerPrefix,omitempty"`
}

// Define the image registry
type Registry struct {
	// example: quay.io/open-cluster-management
	Server   string `yaml:"server"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// CloudConnection struct for bits having to do with Connections
type CloudConnection struct {
	PullSecret    string  `yaml:"pullSecret"`
	SSHPrivateKey string  `yaml:"sshPrivatekey"`
	SSHPublicKey  string  `yaml:"sshPublickey"`
	Keys          APIKeys `yaml:"apiKeys,omitempty"`
	OCPRelease    string  `yaml:"ocpRelease,omitempty"`
}

type APIKeys struct {
	AWS   AWSAPIKey   `yaml:"aws,omitempty"`
	GCP   GCPAPIKey   `yaml:"gcp,omitempty"`
	Azure AzureAPIKey `yaml:"azure,omitempty"`
}

type AWSAPIKey struct {
	AWSAccessID     string `yaml:"awsAccessKeyID"`
	AWSAccessSecret string `yaml:"awsSecretAccessKeyID"`
	BaseDnsDomain   string `yaml:"baseDnsDomain"`
	Region          string `yaml:"region"`
}

type GCPAPIKey struct {
	ProjectID             string `yaml:"gcpProjectID"`
	ServiceAccountJsonKey string `yaml:"gcpServiceAccountJsonKey"`
	BaseDnsDomain         string `yaml:"baseDnsDomain"`
	Region                string `yaml:"region"`
}

type AzureAPIKey struct {
	BaseDnsDomain        string `yaml:"baseDnsDomain"`
	BaseDomainRGN        string `yaml:"azureBaseDomainRGN"`
	ServicePrincipalJson string `yaml:"azureServicePrincipalJson"`
	Region               string `yaml:"region"`
}
