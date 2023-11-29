package talos

import (
	"errors"
	"testing"

	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/siderolabs/image-factory/pkg/schematic"
)

func TestOnlineMode(t *testing.T) {
	data := &schematic.Schematic{
		Customization: schematic.Customization{
			ExtraKernelArgs: []string{"net.ifnames=0"},
			SystemExtensions: schematic.SystemExtensions{
				OfficialExtensions: []string{"siderolabs/intel-ucode"},
			},
		},
	}
	jsonByte, err := data.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	onlineOut := map[string]string{}
	if err := doHTTPPOSTRequest(jsonByte, "https://factory.talos.dev/schematics", &onlineOut); err != nil {
		if errors.Is(err, errNotStatusCreated) {
			t.Skipf("%v. Skipping this test...", err)
		}
		t.Fatal(err)
	}
	expectedID, err := data.ID()
	if err != nil {
		t.Fatal(err)
	}

	if onlineOut["id"] != expectedID {
		t.Errorf("got %s, want %s", onlineOut["id"], expectedID)
	}
}

func TestGetInstallerURL(t *testing.T) {
	type testStruct struct {
		name        string
		cfg         *schematic.Schematic
		iFactory    *config.ImageFactory
		version     string
		expectedURL string
	}

	for _, test := range []testStruct{
		{
			name: "default",
			cfg:  &schematic.Schematic{},
			iFactory: &config.ImageFactory{
				RegistryURL: "factory.talos.dev",
			},
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
			iFactory: &config.ImageFactory{
				RegistryURL: "",
			},
			version:     "v1.5.4",
			expectedURL: "factory.talos.dev/installer/98442b5bb4e8d050f30978ce3e6ec22e7bf534d57cafcd51313235128057e612:v1.5.4",
		},

		{
			name: "withKernelArgs",
			cfg: &schematic.Schematic{
				Customization: schematic.Customization{
					ExtraKernelArgs: []string{"hihi", "hehe"},
				},
			},
			iFactory:    &config.ImageFactory{},
			expectedURL: "factory.talos.dev/installer/ff5083b14ccb03821ea738d712ac08a82b44d2693013622059edaae286665239:",
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
			iFactory: &config.ImageFactory{
				RegistryURL: "test.registry/",
			},
			version:     "1.5.4",
			expectedURL: "test.registry//installer/104c23dfe7c5bfeff6a4cc7e166d8b3bba0f371760592c7677c90c822bb1d109:1.5.4",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := &config.TalhelperConfig{}
			cfg.ImageFactory = *test.iFactory
			url, err := GetInstallerURL(test.cfg, cfg.GetImageFactory(), test.version, true)
			if err != nil {
				t.Fatal(err)
			}
			if url != test.expectedURL {
				t.Errorf("got %s, want %s", url, test.expectedURL)
			}
		})
	}
}
