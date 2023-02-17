package handlers

import (
	"log"
	"net/http"

	"github.com/KrzysztofSieczkiewicz/go-service-template/src/data"
)

type Mesh struct {
	l *log.Logger
}

func NewMesh(l *log.Logger) *Mesh {
	return &Mesh{l}
}

func (m *Mesh) ServeHTTP(rw http.ResponseWriter, h *http.Request) {

}

func (m *Mesh) getMeshes(rw http.ResponseWriter, h *http.Request) {
	lm := data.GetMeshes()
	err := lm.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
