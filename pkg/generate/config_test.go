package generate

import (
	"os"
	"strings"
	"testing"
)

func TestGetFileContentByte_SopsEncryptedFile(t *testing.T) {
	os.Setenv("SOPS_AGE_KEY", "AGE-SECRET-KEY-172FENV3SDP8JSRRX2SWTA9JQMAW7MW3GSKJ2JZDNXS4GVFAS5STQUW8WN4")

	result, err := getFileContentByte("testdata/encrypted-manifest.sops.yaml")
	if err != nil {
		t.Fatal(err)
	}

	got := string(result)
	if !strings.Contains(got, "password: p4ssw0rd") {
		t.Errorf("expected decrypted content to contain 'password: p4ssw0rd', got %q", got)
	}
	if !strings.Contains(got, "name: db-credentials") {
		t.Errorf("expected decrypted content to contain 'name: db-credentials', got %q", got)
	}
}

func TestGetFileContentByte_UnencryptedFile(t *testing.T) {
	result, err := getFileContentByte("testdata/unencrypted-manifest.yaml")
	if err != nil {
		t.Fatal(err)
	}

	got := string(result)
	if !strings.Contains(got, "name: production") {
		t.Errorf("expected content to contain 'name: production', got %q", got)
	}
}

func TestGetFileContentByte_NonExistentFile(t *testing.T) {
	result, err := getFileContentByte("testdata/does-not-exist.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if result != nil {
		t.Errorf("expected nil for non-existent file, got %q", string(result))
	}
}

func TestGetFileContentByte_SopsEncryptedWithMissingKey(t *testing.T) {
	os.Setenv("SOPS_AGE_KEY", "")

	_, err := getFileContentByte("testdata/encrypted-manifest.sops.yaml")
	if err == nil {
		t.Fatal("expected SOPS decryption error when key is missing, got nil")
	}
}
