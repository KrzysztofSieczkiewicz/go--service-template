package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/KrzysztofSieczkiewicz/go-service-template/src/data"
	"github.com/gorilla/mux"
)

type Meshes struct {
	l *log.Logger
}

func NewMeshes(l *log.Logger) *Meshes {
	return &Meshes{l}
}

func (m *Meshes) GetMeshes(rw http.ResponseWriter, r *http.Request) {
	m.l.Println("Handle GET Meshes")

	lm := data.GetMeshes()

	err := lm.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (m *Meshes) AddMesh(rw http.ResponseWriter, r *http.Request) {
	m.l.Println("Handle POST Mesh")

	mesh := r.Context().Value(KeyMesh{}).(*data.Mesh)
	data.AddMesh(mesh)
}

func (m *Meshes) UpdateMeshes(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert ID to int", http.StatusBadRequest)
		return
	}

	m.l.Println("Handle PUT Meshes", id)
	mesh := r.Context().Value(KeyMesh{}).(*data.Mesh)

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

type KeyMesh struct{}

func (m Meshes) MiddlewareValidateMesh(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		mesh := &data.Mesh{}

		err := mesh.FromJSON(r.Body)
		if err != nil {
			m.l.Println("[ERROR] deserializing mesh JSON: ", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// validate the mesh
		err = mesh.Validate()
		if err != nil {
			m.l.Println("[ERROR] validating mesh: ", err)

			http.Error(
				rw,
				fmt.Sprintf("Error validating mesh: %s", err),
				http.StatusBadRequest)

			return
		}

		// add product to context
		ctx := context.WithValue(r.Context(), KeyMesh{}, mesh)
		req := r.WithContext(ctx)

		// call next handler
		next.ServeHTTP(rw, req)
	})
}
