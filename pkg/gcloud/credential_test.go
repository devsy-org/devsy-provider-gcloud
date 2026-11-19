package gcloud

import "testing"

func TestCredentialValue(t *testing.T) {
	t.Setenv("PROVIDER_ID", "")
	t.Setenv("KEY_FILE", "")
	t.Setenv("DEVSY_PROVIDER_GCLOUD_KEY_FILE", "")

	if got := credentialValue("KEY_FILE"); got != "" {
		t.Fatalf("empty environment returned %q", got)
	}

	t.Setenv("DEVSY_PROVIDER_GCLOUD_KEY_FILE", "from-env")
	if got := credentialValue("KEY_FILE"); got != "from-env" {
		t.Fatalf("scoped env = %q", got)
	}

	t.Setenv("KEY_FILE", "from-option")
	if got := credentialValue("KEY_FILE"); got != "from-option" {
		t.Fatalf("injected option = %q, want the resolved option to win", got)
	}

	t.Setenv("KEY_FILE", "")
	t.Setenv("PROVIDER_ID", "mygcloud")
	t.Setenv("DEVSY_PROVIDER_MYGCLOUD_KEY_FILE", "renamed")
	if got := credentialValue("KEY_FILE"); got != "renamed" {
		t.Fatalf("renamed provider = %q", got)
	}

	t.Setenv("DEVSY_PROVIDER_MYGCLOUD_KEY_FILE", "")
	if got := credentialValue("KEY_FILE"); got != "" {
		t.Fatalf("renamed provider must not reuse the gcloud key, got %q", got)
	}
}
