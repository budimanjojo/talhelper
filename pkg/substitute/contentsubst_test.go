package substitute

import (
	"os"
	"strings"
	"testing"
)

func TestSubstituteFileContent_SopsEncryptedFile(t *testing.T) {
	os.Setenv("SOPS_AGE_KEY", "AGE-SECRET-KEY-172FENV3SDP8JSRRX2SWTA9JQMAW7MW3GSKJ2JZDNXS4GVFAS5STQUW8WN4")

	result, err := SubstituteFileContent("@./testdata/encrypted-manifest.sops.yaml", false)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(result, "password: supersecret") {
		t.Errorf("expected decrypted content to contain 'password: supersecret', got %q", result)
	}
	if !strings.Contains(result, "name: my-secret") {
		t.Errorf("expected decrypted content to contain 'name: my-secret', got %q", result)
	}
}

func TestSubstituteFileContent_SopsEncryptedWithMissingKey(t *testing.T) {
	os.Setenv("SOPS_AGE_KEY", "")

	_, err := SubstituteFileContent("@./testdata/encrypted-manifest.sops.yaml", false)
	if err == nil {
		t.Fatal("expected SOPS decryption error when key is missing, got nil")
	}
}

func TestSubstituteFileContent_SopsEncryptedWithEnvsubst(t *testing.T) {
	os.Setenv("SOPS_AGE_KEY", "AGE-SECRET-KEY-172FENV3SDP8JSRRX2SWTA9JQMAW7MW3GSKJ2JZDNXS4GVFAS5STQUW8WN4")
	os.Setenv("EXPECTED_NAME", "my-secret")

	result, err := SubstituteFileContent("@./testdata/encrypted-manifest.sops.yaml", true)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(result, "password: supersecret") {
		t.Errorf("expected decrypted content to contain 'password: supersecret', got %q", result)
	}
}

func TestSubstituteFileContent_NonSopsFile(t *testing.T) {
	result, err := SubstituteFileContent("@./testdata/content.yaml", false)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(result, "hubble-ui-nginx") {
		t.Errorf("expected content to contain 'hubble-ui-nginx', got %q", result)
	}
}

func TestSubstituteFileContent_NoAtPrefix(t *testing.T) {
	result, err := SubstituteFileContent("literal value", false)
	if err != nil {
		t.Fatal(err)
	}

	if result != "literal value" {
		t.Errorf("got %q, want %q", result, "literal value")
	}
}

func TestSubstituteFileContent(t *testing.T) {
	contents := []string{
		"@./testdata/content.yaml",
		"this is $host",
	}
	envs := map[string]string{
		"host":           "substhost",
		"remote_addr":    "substremote_addr",
		"request_method": "substrequest_method",
		"uri":            "substuri",
	}

	for env, value := range envs {
		os.Setenv(env, value)
	}

	expectedWithoutEnvsubst := []string{
		`---
# Source: cilium/templates/hubble-ui/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: hubble-ui-nginx
  namespace: kube-system
data:
  nginx.conf: "server {\n    listen       8081;\n    listen       [::]:8081;\n    server_name  localhost;\n    root /app;\n    index index.html;\n    client_max_body_size 1G;\n\n    location / {\n        proxy_set_header Host $host;\n        proxy_set_header X-Real-IP $remote_addr;\n\n        # CORS\n        add_header Access-Control-Allow-Methods \"GET, POST, PUT, HEAD, DELETE, OPTIONS\";\n        add_header Access-Control-Allow-Origin *;\n        add_header Access-Control-Max-Age 1728000;\n        add_header Access-Control-Expose-Headers content-length,grpc-status,grpc-message;\n        add_header Access-Control-Allow-Headers range,keep-alive,user-agent,cache-control,content-type,content-transfer-encoding,x-accept-content-transfer-encoding,x-accept-response-streaming,x-user-agent,x-grpc-web,grpc-timeout;\n        if ($request_method = OPTIONS) {\n            return 204;\n        }\n        # /CORS\n\n        location /api {\n            proxy_http_version 1.1;\n            proxy_pass_request_headers on;\n            proxy_hide_header Access-Control-Allow-Origin;\n            proxy_pass http://127.0.0.1:8090;\n        }\n        location / {\n            # double ` + "`" + "/index.html" + "`" + ` is required here \n            try_files $uri $uri/ /index.html /index.html;\n        }\n\n        # Liveness probe\n        location /healthz {\n            access_log off;\n            add_header Content-Type text/plain;\n            return 200 'ok';\n        }\n    }\n}"
---`,
		"this is $host",
	}

	expectedWithEnvsubst := []string{
		`apiVersion: v1
kind: ConfigMap
metadata:
  name: hubble-ui-nginx
  namespace: kube-system
data:
  nginx.conf: "server {\n    listen       8081;\n    listen       [::]:8081;\n    server_name  localhost;\n    root /app;\n    index index.html;\n    client_max_body_size 1G;\n\n    location / {\n        proxy_set_header Host substhost;\n        proxy_set_header X-Real-IP substremote_addr;\n\n        # CORS\n        add_header Access-Control-Allow-Methods \"GET, POST, PUT, HEAD, DELETE, OPTIONS\";\n        add_header Access-Control-Allow-Origin *;\n        add_header Access-Control-Max-Age 1728000;\n        add_header Access-Control-Expose-Headers content-length,grpc-status,grpc-message;\n        add_header Access-Control-Allow-Headers range,keep-alive,user-agent,cache-control,content-type,content-transfer-encoding,x-accept-content-transfer-encoding,x-accept-response-streaming,x-user-agent,x-grpc-web,grpc-timeout;\n        if (substrequest_method = OPTIONS) {\n            return 204;\n        }\n        # /CORS\n\n        location /api {\n            proxy_http_version 1.1;\n            proxy_pass_request_headers on;\n            proxy_hide_header Access-Control-Allow-Origin;\n            proxy_pass http://127.0.0.1:8090;\n        }\n        location / {\n            # double ` + "`" + "/index.html" + "`" + ` is required here \n            try_files substuri substuri/ /index.html /index.html;\n        }\n\n        # Liveness probe\n        location /healthz {\n            access_log off;\n            add_header Content-Type text/plain;\n            return 200 'ok';\n        }\n    }\n}"
---`, // header and comments are removed by SubstituteEnvFromByte()
		"this is $host", // ignored when it's not a file
	}

	for k, v := range contents {
		r1, err := SubstituteFileContent(v, false)
		if err != nil {
			t.Fatal(err)
		}

		if expectedWithoutEnvsubst[k] != strings.TrimSpace(r1) {
			t.Errorf("got\n%s,\bwant\n%s", strings.TrimSpace(r1), expectedWithoutEnvsubst[k])
		}

		r2, err := SubstituteFileContent(v, true)
		if err != nil {
			t.Fatal(err)
		}
		if expectedWithEnvsubst[k] != strings.TrimSpace(r2) {
			t.Errorf("got\n%s,\bwant\n%s", strings.TrimSpace(r2), expectedWithEnvsubst[k])
		}
	}
}
