package providers

import (
	"bytes"
	"text/template"
)

type InstallerConfigAWS struct {
	Name          string
	BaseDnsDomain string
	SSHKey        string
	Region        string
}

// function for filling out installconfig template, takes installConfig struct and returns string for secret creation
func GetInstallConfigAWS(instConfig InstallerConfigAWS) string {
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
    aws:
      rootVolume:
        iops: 4000
        size: 500
        type: io1
      type: m4.xlarge
compute:
- hyperthreading: Enabled
  name: worker
  replicas: 3
  platform:
      aws:
      rootVolume:
        iops: 2000
        size: 500
        type: io1
      type: m4.large
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  machineCIDR: 10.0.0.0/16
  networkType: OpenShiftSDN
  serviceNetwork:
  - 172.30.0.0/16
platform:
  aws:
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
