package decrypt

import (
	"os"
	"strings"
	"testing"
)

func TestDecryptFileWithSops_EncryptedYaml(t *testing.T) {
	os.Setenv("SOPS_AGE_KEY", "AGE-SECRET-KEY-172FENV3SDP8JSRRX2SWTA9JQMAW7MW3GSKJ2JZDNXS4GVFAS5STQUW8WN4")

	result, err := DecryptFileWithSops("testdata/encrypted.yaml")
	if err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(string(result))
	expected := "hello: world\nsecret: mysecretvalue"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestDecryptFileWithSops_EncryptedJson(t *testing.T) {
	os.Setenv("SOPS_AGE_KEY", "AGE-SECRET-KEY-172FENV3SDP8JSRRX2SWTA9JQMAW7MW3GSKJ2JZDNXS4GVFAS5STQUW8WN4")

	result, err := DecryptFileWithSops("testdata/encrypted.json")
	if err != nil {
		t.Fatal(err)
	}

	got := string(result)
	if !strings.Contains(got, `"hello": "world"`) {
		t.Errorf("expected decrypted JSON to contain hello: world, got %q", got)
	}
	if !strings.Contains(got, `"secret": "mysecretvalue"`) {
		t.Errorf("expected decrypted JSON to contain secret: mysecretvalue, got %q", got)
	}
}

func TestDecryptFileWithSops_EncryptedDotenv(t *testing.T) {
	os.Setenv("SOPS_AGE_KEY", "AGE-SECRET-KEY-172FENV3SDP8JSRRX2SWTA9JQMAW7MW3GSKJ2JZDNXS4GVFAS5STQUW8WN4")

	result, err := DecryptFileWithSops("testdata/encrypted.env")
	if err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(string(result))
	if !strings.Contains(got, "HELLO=world") {
		t.Errorf("expected decrypted dotenv to contain HELLO=world, got %q", got)
	}
	if !strings.Contains(got, "SECRET=mysecretvalue") {
		t.Errorf("expected decrypted dotenv to contain SECRET=mysecretvalue, got %q", got)
	}
}

func TestDecryptFileWithSops_UnencryptedYaml(t *testing.T) {
	result, err := DecryptFileWithSops("testdata/unencrypted.yaml")
	if err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(string(result))
	expected := "hello: world\nsecret: mysecretvalue"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestDecryptFileWithSops_UnencryptedJson(t *testing.T) {
	result, err := DecryptFileWithSops("testdata/unencrypted.json")
	if err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(string(result))
	if !strings.Contains(got, `"hello": "world"`) {
		t.Errorf("expected JSON to contain hello: world, got %q", got)
	}
}

func TestDecryptFileWithSops_UnencryptedDotenv(t *testing.T) {
	result, err := DecryptFileWithSops("testdata/unencrypted.env")
	if err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(string(result))
	if !strings.Contains(got, "HELLO=world") {
		t.Errorf("expected dotenv to contain HELLO=world, got %q", got)
	}
}

func TestDecryptFileWithSops_PlainTextFile(t *testing.T) {
	result, err := DecryptFileWithSops("testdata/plain.txt")
	if err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(string(result))
	expected := "hello: world\nsecret: mysecretvalue"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestDecryptFileWithSops_NonExistentFile(t *testing.T) {
	_, err := DecryptFileWithSops("testdata/does-not-exist.yaml")
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}

func TestDecryptFileWithSops_EncryptedWithMissingKey(t *testing.T) {
	os.Setenv("SOPS_AGE_KEY", "")

	_, err := DecryptFileWithSops("testdata/encrypted.yaml")
	if err == nil {
		t.Fatal("expected SOPS decryption error when key is missing, got nil")
	}

	if !strings.Contains(err.Error(), "SOPS decryption failed") {
		t.Errorf("expected error to contain 'SOPS decryption failed', got %q", err.Error())
	}
}
