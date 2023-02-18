package handlers

import (
	"log"
	"net/http"

	"github.com/KrzysztofSieczkiewicz/go-service-template/src/data"
)

type Meshes struct {
	l *log.Logger
}

func NewMeshes(l *log.Logger) *Meshes {
	return &Meshes{l}
}

func (m *Meshes) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		m.getMeshes(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		m.addMesh(rw, r)
		return
	}

	//catch remaining
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (m *Meshes) getMeshes(rw http.ResponseWriter, r *http.Request) {
	m.l.Println("Handle GET Meshes")

	lm := data.GetMeshes()

	err := lm.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (m *Meshes) addMesh(rw http.ResponseWriter, r *http.Request) {
	m.l.Println("Handle POST Meshes")

	mesh := &data.Mesh{}
	err := mesh.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Undable to unmarshal json", http.StatusBadRequest)
	}

	m.l.Printf("Mesh %#v", mesh)
}
