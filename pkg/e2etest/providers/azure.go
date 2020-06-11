package providers

import (
	"bytes"
	"text/template"
)

type InstallerConfigAzure struct {
	Name          string
	BaseDnsDomain string
	SSHKey        string
	BaseDomainRGN string
	Region        string
}

// function for filling out installconfig template, takes installConfig struct and returns string for secret creation
func GetInstallConfigAzure(instConfig InstallerConfigAzure) string {
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
    azure:
      osDisk:
      diskSizeGB: 128
    type:  Standard_D4s_v3
compute:
- hyperthreading: Enabled
  name: worker
  replicas: 3
  platform:
    azure:
      type:  Standard_D2s_v3
      osDisk:
      diskSizeGB: 128
      zones:
      - "1"
      - "2"
      - "3"
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  machineCIDR: 10.0.0.0/16
  networkType: OpenShiftSDN
  serviceNetwork:
  - 172.30.0.0/16
platform:
  azure:
    baseDomainResourceGroupName: {{.BaseDomainRGN}}
    region: {{.Region}}
pullSecret: "" # skip, hive will inject based on it's secrets
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
