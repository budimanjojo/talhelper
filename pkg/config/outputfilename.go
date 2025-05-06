package config

import (
	"bytes"
	"text/template"
)

type filenameTmpl struct {
	ClusterName string
	Hostname    string
	IPAddress   string
	Role        string
}

func (n *Node) GetOutputFileName(c *TalhelperConfig) (string, error) {
	role := "controlplane"
	if !n.ControlPlane {
		role = "worker"
	}

	tmplData := filenameTmpl{
		ClusterName: c.ClusterName,
		Hostname:    n.Hostname,
		IPAddress:   n.IPAddress,
		Role:        role,
	}

	t, err := template.New("filename").Parse(n.GetFilenameTmpl())
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, tmplData); err != nil {
		return "", err
	}

	return buf.String(), nil
}
