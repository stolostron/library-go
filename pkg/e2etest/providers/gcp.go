package providers

import (
	"bytes"
	"text/template"
)

type InstallerConfigGCP struct {
	Name          string
	BaseDnsDomain string
	SSHKey        string
	ProjectID     string
	Region        string
}

// function for filling out installconfig template, takes installConfig struct and returns string for secret creation
func GetInstallConfigGCP(instConfig InstallerConfigGCP) string {
	const configTemplate = `
apiVersion: v1
metadata:
  name: {{.Name}}
baseDomain: {{.BaseDnsDomain}}
controlPlane:
  hyperthreading: Enabled
  name: master
  replicas: 3
  platform:
    gcp:
      type: n1-standard-4
compute:
- hyperthreading: Enabled
  name: worker
  replicas: 3
  platform:
      gcp:
        type: n1-standard-4
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  machineCIDR: 10.0.0.0/16
  networkType: OpenShiftSDN
  serviceNetwork:
  - 172.30.0.0/16
platform:
  gcp:
    projectID: {{.ProjectID}}
    region: {{.Region}}
pullSecret: ""
sshKey: {{.SSHKey}}
`
	stdoutBuffer := new(bytes.Buffer)
	t := template.Must(template.New("configTemplate").Parse(configTemplate))
	err := t.Execute(stdoutBuffer, instConfig)
	if err != nil {
		return err.Error()
	}
	return stdoutBuffer.String()
}
