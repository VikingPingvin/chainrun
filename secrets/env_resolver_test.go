package secrets

import (
	"os"
	"testing"
)

func helper_SetupEnvs(t *testing.T) {
	t.Helper()

	envsToAdd := map[string]string{
		"TEST_ENV_1": "value1",
		"TEST_ENV_2": "value2",
		"TEST_ENV_3": "value3",
	}

	// add a set of envs for testing
	for k, v := range envsToAdd {
		if err := os.Setenv(k, v); err != nil {
			t.Fatalf("failed to set env %s: %v", k, err)
		}
	}
}

func TestResolveSecretsSingle(t *testing.T) {
	testCases := []struct {
		desc        string
		expectError bool
		keyName     string
	}{
		{desc: "existing env", expectError: false, keyName: "TEST_ENV_1"},
		{desc: "non-existing env", expectError: true, keyName: "NON_EXISTING_ENV"},
	}

	helper_SetupEnvs(t)
	resolver := NewEnvResolver()

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			value, found := resolver.Resolve(tC.keyName)
			if tC.expectError {
				if found {
					t.Errorf("expected error for key %s, but got value %s", tC.keyName, value)
				}
			} else {
				if !found {
					t.Errorf("expected value for key %s, but got error", tC.keyName)
				}
			}
		})
	}
}

func TestResolveSecretsAll(t *testing.T) {
	helper_SetupEnvs(t)

	resolver := NewEnvResolver()
	allSecrets := resolver.ResolveAll()

	expectedEnvs := []string{"TEST_ENV_1", "TEST_ENV_2", "TEST_ENV_3"}
	for _, env := range expectedEnvs {
		if _, found := allSecrets[env]; !found {
			t.Errorf("expected env %s to be present in resolved secrets, but it was not found", env)
		}
	}
}
