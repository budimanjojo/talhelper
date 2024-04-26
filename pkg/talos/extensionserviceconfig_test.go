package talos

import (
	"testing"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/siderolabs/talos/pkg/machinery/config/types/runtime/extensions"
	"gopkg.in/yaml.v3"
)

func TestGenerateNodeExtensionServiceConfig(t *testing.T) {
	data := []byte(`nodes:
  - hostname: node1
    extensionServices:
      - name: nut-client
        configFiles:
          - content: MONITOR ${upsmonHost} 1 remote pass password
            mountPath: /usr/local/etc/nut/upsmon.conf
        environment:
          - UPS_NAME=ups
      - name: nut-client2
        configFiles:
          - content: hello
            mountPath: /etc/hello
          - content: hello2
            mountPath: /etc/hello2`)

	var m config.TalhelperConfig
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	expectedExtension1Name := "nut-client"
	expectedExtension1ConfigFiles := []extensions.ConfigFile{
		{
			ConfigFileContent:   "MONITOR ${upsmonHost} 1 remote pass password",
			ConfigFileMountPath: "/usr/local/etc/nut/upsmon.conf",
		},
	}
	expectedExtension1Environment := []string{"UPS_NAME=ups"}
	expectedExtension2Name := "nut-client2"
	expectedExtension2ConfigFiles := []extensions.ConfigFile{
		{
			ConfigFileContent:   "hello",
			ConfigFileMountPath: "/etc/hello",
		},
		{
			ConfigFileContent:   "hello2",
			ConfigFileMountPath: "/etc/hello2",
		},
	}

	result, err := GenerateNodeExtensionServiceConfig(m.Nodes[0].ExtensionServices)
	if err != nil {
		t.Fatal(err)
	}

	compare(result[0].Name(), expectedExtension1Name, t)
	compare(result[0].ServiceConfigFiles, expectedExtension1ConfigFiles, t)
	compare(result[0].Environment(), expectedExtension1Environment, t)
	compare(result[1].Name(), expectedExtension2Name, t)
	compare(result[1].ServiceConfigFiles, expectedExtension2ConfigFiles, t)
}
