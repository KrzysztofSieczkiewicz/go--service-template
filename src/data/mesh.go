package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type Mesh struct {
	ID          int    `json:"id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	CreatedOn   string `json:"-"`
	Address     string `json:"address" validate:"required,address"`
}

type Meshes []*Mesh

func (m *Mesh) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(m)
}

func (m *Meshes) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(m)
}

func (m *Mesh) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("address", validateAddress)

	return validate.Struct(m)
}

func validateAddress(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^(.*/)([^/]*)$`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}

func GetMeshes() Meshes {
	return meshList
}

func AddMesh(m *Mesh) {
	m.ID = GetNextID()
	meshList = append(meshList, m)
}

func UpdateMesh(id int, m *Mesh) error {
	_, pos, err := FindMesh(id)
	if err != nil {
		return err
	}

	m.ID = id
	meshList[pos] = m

	return nil
}

func GetNextID() int {
	m := meshList[len(meshList)-1]
	return m.ID + 1
}

var ErrMeshNotFound = fmt.Errorf("Mesh not found")

func FindMesh(id int) (*Mesh, int, error) {
	for pos, m := range meshList {
		if m.ID == id {
			return m, pos, nil
		}
	}
	return nil, -1, ErrMeshNotFound
}

var meshList = []*Mesh{
	{
		ID:          1,
		Name:        "Apple_01",
		Description: "A whole apple",
		CreatedOn:   time.Now().UTC().String(),
		Address:     "some/random/folder",
	},
	{
		ID:          2,
		Name:        "Apple_02",
		Description: "Half an apple",
		CreatedOn:   time.Now().UTC().String(),
		Address:     "some/other/random/folder",
	},
	{
		ID:          3,
		Name:        "Apple_03",
		Description: "Quarter an apple",
		CreatedOn:   time.Now().UTC().String(),
		Address:     "some/another/random/folder",
	},
}
