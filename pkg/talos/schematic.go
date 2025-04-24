package talos

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/budimanjojo/talhelper/v3/pkg/config"
	"github.com/siderolabs/image-factory/pkg/schematic"
)

var (
	errNotStatusCreated = errors.New("server not replying StatusCreated")
	errFailedtoPost     = errors.New("server not wanting to reply")
)

type factoryPOSTResult struct {
	ID string `json:"id"`
}

type installerTmpl struct {
	RegistryURL string
	ID          string
	Version     string
	Secureboot  bool
}

type imageTmpl struct {
	BootMethod  string
	Suffix      string
	Protocol    string
	RegistryURL string
	ID          string
	Version     string
	Mode        string
	Arch        string
	Secureboot  bool
	UseUKI      bool
}

func GetInstallerURL(cfg *schematic.Schematic, factory *config.ImageFactory, spec *config.MachineSpec, version string, offlineMode bool) (string, error) {
	tmplData := installerTmpl{
		RegistryURL: factory.RegistryURL,
		Version:     version,
		Secureboot:  spec.Secureboot,
	}

	id, err := getSchematicID(cfg, factory, offlineMode)
	if err != nil {
		return "", err
	}
	tmplData.ID = id

	t, err := template.New("installer").Parse(factory.InstallerURLTmpl)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, tmplData); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetImageURL(cfg *schematic.Schematic, factory *config.ImageFactory, spec *config.MachineSpec, version string, offlineMode bool) (string, error) {
	data := imageTmpl{
		BootMethod:  spec.BootMethod,
		Suffix:      spec.ImageSuffix,
		Protocol:    factory.Protocol,
		RegistryURL: factory.RegistryURL,
		Version:     version,
		Mode:        spec.Mode,
		Arch:        spec.Arch,
		Secureboot:  spec.Secureboot,
		UseUKI:      spec.UseUKI,
	}

	if spec.ImageSuffix != "" {
		data.Suffix = "." + spec.ImageSuffix
	} else {
		data.Suffix = parseImageSuffix(data.BootMethod, data.Secureboot, data.UseUKI)
	}

	id, err := getSchematicID(cfg, factory, offlineMode)
	if err != nil {
		return "", err
	}
	data.ID = id

	return genImageURL(&data, factory)
}

func genImageURL(data *imageTmpl, factory *config.ImageFactory) (string, error) {
	t, err := template.New("image").Parse(factory.ImageURLTmpl)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func parseImageSuffix(bm string, sb, uki bool) string {
	var ext string
	switch bm {
	case "":
		ext = ".iso"
	case "iso":
		if sb && uki {
			ext = ".efi"
		} else {
			ext = ".iso"
		}
	case "disk-image":
		ext = ".raw.zxt"
	case "pxe":
		ext = ""
	}
	return ext
}

func getSchematicID(cfg *schematic.Schematic, iFactory *config.ImageFactory, offlineMode bool) (string, error) {
	body, err := cfg.Marshal()
	if err != nil {
		return "", err
	}

	slog.Debug(fmt.Sprintf("defined schematic:\n%s", body))

	if offlineMode {
		slog.Debug("generating schematic ID in offline mode")
		id, err := cfg.ID()
		if err != nil {
			return "", err
		}
		return id, nil
	}
	var resp factoryPOSTResult
	schematicURL := iFactory.Protocol + "://" + iFactory.RegistryURL + iFactory.SchematicEndpoint
	slog.Debug(fmt.Sprintf("generating schematic ID from %s", schematicURL))
	if err := doHTTPPOSTRequest(body, schematicURL, &resp); err != nil {
		return "", err
	}
	return resp.ID, nil
}

func doHTTPPOSTRequest(body []byte, url string, out interface{}) error {
	resp, err := http.Post(url, "", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("%w: %v", errFailedtoPost, err)
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("%w (%v): %v", errNotStatusCreated, url, resp.Status)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return err
	}
	return nil
}
