package handlers

import (
	"fmt"
	"net/http"

	"github.com/doha-ms/bnb-booking-app/cmd/pkg/config"
	"github.com/doha-ms/bnb-booking-app/cmd/pkg/models"
	"github.com/doha-ms/bnb-booking-app/cmd/pkg/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {

	return &Repository{
		App: a,
	}
}
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})

}
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SESSION DATA BEFORE RETRIEVAL:", m.App.Session.GetString(r.Context(), "remote_ip"))
	stringMap := make(map[string]string)
	stringMap["test"] = "oh hi i didnt see you there :)"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
