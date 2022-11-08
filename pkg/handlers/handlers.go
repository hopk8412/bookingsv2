package handlers

import (
	"bookingsv2/pkg/config"
	"bookingsv2/pkg/models"
	"bookingsv2/pkg/render"
	"net/http"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type for the pattern
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository {
		App: a,
	}
}

// NewHandlers sets the Repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//handler function - in order to respond to a web request, must include a responsewriter and a pointer to a request

// Home is the home page handler
//giving these functions a receiver, this links it to the repository so all handlers have access to the repository
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remoteIp", remoteIp)
	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

// About is about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)

	stringMap["test"] = "Welcome back"

	remoteIp := m.App.Session.GetString(r.Context(), "remoteIp")
	stringMap["remoteIp"] = remoteIp
	// send data to template
	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

//using the receiver gives them access to everything in the repository