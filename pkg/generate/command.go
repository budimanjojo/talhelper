package generate

import (
	"fmt"
	"github.com/budimanjojo/talhelper/pkg/config"
)

func GenerateCommand(c *config.TalhelperConfig, outDir string, generateApply bool, generateUpgrade bool) error {
	if !generateApply && !generateUpgrade {
		fmt.Printf("Must select one\n")
		return nil
	}

	for _, node := range c.Nodes {
		if generateApply {
			fileName := outDir + "/" + c.ClusterName + "-" + node.Hostname + ".yaml"
			fmt.Printf("talosctl apply-config --talosconfig %s/talosconfig --nodes %s --file %s --insecure;\n", outDir, node.IPAddress, fileName)
		}
		
		if generateUpgrade {
			var image = "test123"
			fmt.Printf("talosctl upgrade --talosconfig %s/talosconfig --nodes %s --image %s;\n", outDir, node.IPAddress, image)
		}
	}
	return nil
}