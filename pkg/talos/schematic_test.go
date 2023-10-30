package talos

import (
	"testing"

	"github.com/siderolabs/image-factory/pkg/schematic"
)

func TestGetInstallerURL(t *testing.T) {
	type testStruct struct {
		name        string
		cfg         *schematic.Schematic
		registryURL string
		version     string
		expectedURL string
	}

	for _, test := range []testStruct{
		{
			name:        "default",
			cfg:         &schematic.Schematic{},
			registryURL: "factory.talos.dev/installer",
			version:     "v1.5.4",
			expectedURL: "factory.talos.dev/installer/376567988ad370138ad8b2698212367b8edcb69b5fd68c80be1f2ec7d603b4ba:v1.5.4",
		},

		{
			name: "withExts",
			cfg: &schematic.Schematic{
				Customization: schematic.Customization{
					SystemExtensions: schematic.SystemExtensions{
						OfficialExtensions: []string{"siderolabs/drbd", "siderolabs/zfs"},
					},
				},
			},
			registryURL: "",
			version:     "v1.5.4",
			expectedURL: "/98442b5bb4e8d050f30978ce3e6ec22e7bf534d57cafcd51313235128057e612:v1.5.4",
		},

		{
			name: "withKernelArgs",
			cfg: &schematic.Schematic{
				Customization: schematic.Customization{
					ExtraKernelArgs: []string{"hihi", "hehe"},
				},
			},
			expectedURL: "/ff5083b14ccb03821ea738d712ac08a82b44d2693013622059edaae286665239:",
		},

		{
			name: "withBoth",
			cfg: &schematic.Schematic{
				Customization: schematic.Customization{
					SystemExtensions: schematic.SystemExtensions{
						OfficialExtensions: []string{"siderolabs/tailscale"},
					},
					ExtraKernelArgs: []string{"net.ifnames=0"},
				},
			},
			registryURL: "test.registry/",
			version:     "1.5.4",
			expectedURL: "test.registry/104c23dfe7c5bfeff6a4cc7e166d8b3bba0f371760592c7677c90c822bb1d109:1.5.4",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			url, err := GetInstallerURL(test.cfg, test.registryURL, test.version)
			if err != nil {
				t.Fatal(err)
			}
			if url != test.expectedURL {
				t.Errorf("got %s, want %s", url, test.expectedURL)
			}
		})
	}
}
