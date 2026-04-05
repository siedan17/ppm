package http

import (
	"net/http"
	"strconv"

	"github.com/daniel/ppm/internal/application"
	"github.com/daniel/ppm/internal/delivery/render"
	"github.com/daniel/ppm/internal/domain"
)

type PeopleHandler struct {
	svc    *application.PersonService
	render *render.Renderer
}

func NewPeopleHandler(svc *application.PersonService, r *render.Renderer) *PeopleHandler {
	return &PeopleHandler{svc: svc, render: r}
}

func (h *PeopleHandler) List(w http.ResponseWriter, r *http.Request) {
	people, err := h.svc.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.render.Page(w, http.StatusOK, "person_list.html", render.PageData{
		Title:   "People",
		Content: people,
	})
}

func (h *PeopleHandler) New(w http.ResponseWriter, r *http.Request) {
	h.render.Page(w, http.StatusOK, "person_form.html", render.PageData{
		Title: "New Person",
		Content: map[string]any{
			"Person":      &domain.Person{PersonType: domain.PersonExternal},
			"PersonTypes": domain.PersonTypes,
		},
	})
}

func (h *PeopleHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	p := &domain.Person{
		Name:       r.FormValue("name"),
		Company:    r.FormValue("company"),
		Role:       r.FormValue("role"),
		Email:      r.FormValue("email"),
		Phone:      r.FormValue("phone"),
		PersonType: r.FormValue("person_type"),
	}

	if p.Name == "" {
		h.render.Page(w, http.StatusUnprocessableEntity, "person_form.html", render.PageData{
			Title: "New Person",
			Content: map[string]any{
				"Person":      p,
				"PersonTypes": domain.PersonTypes,
			},
			Flash: "Name is required.",
		})
		return
	}

	if err := h.svc.Create(p); err != nil {
		h.render.Page(w, http.StatusUnprocessableEntity, "person_form.html", render.PageData{
			Title: "New Person",
			Content: map[string]any{
				"Person":      p,
				"PersonTypes": domain.PersonTypes,
			},
			Flash: "Error creating person: " + err.Error(),
		})
		return
	}

	http.Redirect(w, r, "/people", http.StatusSeeOther)
}

func (h *PeopleHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	person, err := h.svc.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.render.Page(w, http.StatusOK, "person_form.html", render.PageData{
		Title: "Edit Person",
		Content: map[string]any{
			"Person":      person,
			"PersonTypes": domain.PersonTypes,
		},
	})
}

func (h *PeopleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	p := &domain.Person{
		ID:         id,
		Name:       r.FormValue("name"),
		Company:    r.FormValue("company"),
		Role:       r.FormValue("role"),
		Email:      r.FormValue("email"),
		Phone:      r.FormValue("phone"),
		PersonType: r.FormValue("person_type"),
	}

	if p.Name == "" {
		h.render.Page(w, http.StatusUnprocessableEntity, "person_form.html", render.PageData{
			Title: "Edit Person",
			Content: map[string]any{
				"Person":      p,
				"PersonTypes": domain.PersonTypes,
			},
			Flash: "Name is required.",
		})
		return
	}

	if err := h.svc.Update(p); err != nil {
		h.render.Page(w, http.StatusUnprocessableEntity, "person_form.html", render.PageData{
			Title: "Edit Person",
			Content: map[string]any{
				"Person":      p,
				"PersonTypes": domain.PersonTypes,
			},
			Flash: "Error updating person: " + err.Error(),
		})
		return
	}

	http.Redirect(w, r, "/people", http.StatusSeeOther)
}

func (h *PeopleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.svc.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/people", http.StatusSeeOther)
}
