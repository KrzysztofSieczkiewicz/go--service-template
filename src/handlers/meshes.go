// Package classification of Product API
//
// # Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
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

// swagger:response meshesResponse
type meshesResponseWrapper struct {
	// All meshes in the system
	// in: body
	Body []data.Mesh
}

// swagger:response noContent
type meshNoContent struct{}

// swagger:parameters deleteMesh
type meshIDParameterWrapper struct {
	// the id of the mesh to be deleted from the database
	// in: path
	// required: true
	ID int `json:"id"`
}

type Meshes struct {
	l *log.Logger
}

func NewMeshes(l *log.Logger) *Meshes {
	return &Meshes{l}
}

// swagger:route GET /meshes meshes listMeshes
// Returns a list of meshes from the data store
// responses:
// 200: meshesResponse

// GetMeshes returns all meshes from the data store
func (m *Meshes) GetMeshes(rw http.ResponseWriter, r *http.Request) {
	m.l.Println("Handle GET Meshes")

	lm := data.GetMeshes()

	err := lm.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// swagger:route POST /meshes meshes postMeshes
// Returns a 200 message
// responses:
// 200: noContent

// PutMesh creates new mesh in database
func (m *Meshes) AddMesh(rw http.ResponseWriter, r *http.Request) {
	m.l.Println("Handle POST Mesh")

	mesh := r.Context().Value(KeyMesh{}).(*data.Mesh)
	data.AddMesh(mesh)
}

// swagger:route PUT /meshes/{id} meshes putMesh
// Returns a 200 message
// responses:
// 200: noContent

// PutMesh updates mesh in database or creates new
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

// swagger:route DELETE /meshes/{id} meshes deleteMesh
// Returns a 201 message
// responses:
// 201: noContent

// DeleteMesh deletes mesh from database
func (m *Meshes) DeleteMesh(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	m.l.Println("Handle DELETE mesh")

	err := data.DeleteMesh(id)

	if err == data.ErrMeshNotFound {
		http.Error(rw, "Mesh not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Error during mesh deletion", http.StatusInternalServerError)
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
