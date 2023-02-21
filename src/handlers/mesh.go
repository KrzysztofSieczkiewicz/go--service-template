package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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
	if r.Method == http.MethodPut {
		m.l.Println("PUT", r.URL.Path)

		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			m.l.Println("Invalid URI - more than one ID")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			m.l.Println("Invalid URI - more than one capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)

		if err != nil {
			m.l.Println("Invalid URI - unable to convert to number")
			http.Error(rw, "Can't retrieve ID as int", http.StatusBadRequest)
			return
		}

		m.updateMeshes(id, rw, r)
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
	m.l.Println("Handle POST Mesh")

	mesh := &data.Mesh{}

	err := mesh.FromJSON(r.Body)
	if err != nil {
		m.l.Printf("Error unmarshaling JSON: %s\n", err)
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)

	}

	data.AddMesh(mesh)
}

func (m *Meshes) updateMeshes(id int, rw http.ResponseWriter, r *http.Request) {
	m.l.Println("Handle PUT Meshes")

	mesh := &data.Mesh{}

	err := mesh.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err = data.UpdateMesh(id, mesh)
	if err == data.ErrMeshNotFound {
		http.Error(rw, "Mesh not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Mesh update error", http.StatusInternalServerError)
		return
	}

}
