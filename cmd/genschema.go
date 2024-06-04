package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/invopop/jsonschema"

	"github.com/spf13/cobra"
)

var genschemaFile string

var GenschemaCmd = &cobra.Command{
	Use:   "genschema",
	Short: "Generate `talconfig.yaml` JSON schema file",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.TalhelperConfig{}
		r := new(jsonschema.Reflector)
		r.FieldNameTag = "yaml"
		r.RequiredFromJSONSchemaTags = true

		// Doesn't work like I thought it should
		// if err := r.AddGoComments("github.com/budimanjojo/talhelper/v3/pkg/config", "./"); err != nil {
		// 	log.Fatalf("failed to add go comments: %v", err)
		// }
		// if err := r.AddGoComments("github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1", "./"); err != nil {
		// 	log.Fatalf("failed to add go comments: %v", err)
		// }
		schema := r.Reflect(&cfg)
		data, _ := json.MarshalIndent(schema, "", "  ")
		if err := os.WriteFile(genschemaFile, data, os.FileMode(0o644)); err != nil {
			log.Fatalf("failed to write file to %s: %v", genschemaFile, err)
		}
	},
}

func init() {
	RootCmd.AddCommand(GenschemaCmd)

	GenschemaCmd.Flags().StringVarP(&genschemaFile, "file", "f", "talconfig.json", "Where to dump the generated json-schema file")
}
