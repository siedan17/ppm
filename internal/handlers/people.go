package handlers

import (
	"net/http"
	"strconv"

	"github.com/daniel/ppm/internal/models"
	"github.com/daniel/ppm/internal/render"
	"github.com/daniel/ppm/internal/repository"
)

type PeopleHandler struct {
	repo   *repository.PersonRepo
	render *render.Renderer
}

func NewPeopleHandler(repo *repository.PersonRepo, r *render.Renderer) *PeopleHandler {
	return &PeopleHandler{repo: repo, render: r}
}

func (h *PeopleHandler) List(w http.ResponseWriter, r *http.Request) {
	people, err := h.repo.List()
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
		Title:   "New Person",
		Content: &models.Person{},
	})
}

func (h *PeopleHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	p := &models.Person{
		Name:    r.FormValue("name"),
		Company: r.FormValue("company"),
		Role:    r.FormValue("role"),
		Email:   r.FormValue("email"),
		Phone:   r.FormValue("phone"),
	}

	if p.Name == "" {
		h.render.Page(w, http.StatusUnprocessableEntity, "person_form.html", render.PageData{
			Title:   "New Person",
			Content: p,
			Flash:   "Name is required.",
		})
		return
	}

	if err := h.repo.Create(p); err != nil {
		h.render.Page(w, http.StatusUnprocessableEntity, "person_form.html", render.PageData{
			Title:   "New Person",
			Content: p,
			Flash:   "Error creating person: " + err.Error(),
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

	person, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.render.Page(w, http.StatusOK, "person_form.html", render.PageData{
		Title:   "Edit Person",
		Content: person,
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

	p := &models.Person{
		ID:      id,
		Name:    r.FormValue("name"),
		Company: r.FormValue("company"),
		Role:    r.FormValue("role"),
		Email:   r.FormValue("email"),
		Phone:   r.FormValue("phone"),
	}

	if p.Name == "" {
		h.render.Page(w, http.StatusUnprocessableEntity, "person_form.html", render.PageData{
			Title:   "Edit Person",
			Content: p,
			Flash:   "Name is required.",
		})
		return
	}

	if err := h.repo.Update(p); err != nil {
		h.render.Page(w, http.StatusUnprocessableEntity, "person_form.html", render.PageData{
			Title:   "Edit Person",
			Content: p,
			Flash:   "Error updating person: " + err.Error(),
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

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// For htmx requests, return empty (row removed)
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/people", http.StatusSeeOther)
}
