package data

import "testing"

func TestChecksValidation(t *testing.T) {
	m := &Mesh{
		Name:        "testname",
		Description: "testDesc",
		Address:     "something/else",
	}

	err := m.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
