package data

import "testing"

func TestChecksValidation(t *testing.T) {
	m := &Mesh{
		Name:        "testname",
		Description: "testDesc",
		Address:     "test/mesh.log",
	}

	err := m.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
