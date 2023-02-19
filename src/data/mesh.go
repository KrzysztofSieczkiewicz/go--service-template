package data

import (
	"encoding/json"
	"io"
	"time"
)

type Mesh struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedOn   string `json:"-"`
	Address     string `json:"address"`
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

func GetMeshes() Meshes {
	return meshList
}

func AddMesh(m *Mesh) {
	m.ID = GetNextID()
	meshList = append(meshList, m)
}

func GetNextID() int {
	m := meshList[len(meshList)-1]
	return m.ID + 1
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
