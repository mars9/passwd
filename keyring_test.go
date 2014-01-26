package passwd

import "testing"

func TestKeyringFactotum(t *testing.T) {
	f := new(factotum)
	err := f.Set("test.service", "test.user", "test.password")
	if err != nil {
		t.Fatalf("set: %v", err)
	}

	password, err := f.Get("test.service", "test.user")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if password != "test.password" {
		t.Fatalf("get: expected test.password, got %s", password)
	}

	err = f.Delete("test.service", "test.user")
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
}
