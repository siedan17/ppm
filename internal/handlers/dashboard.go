package handlers

import (
	"net/http"

	"github.com/daniel/ppm/internal/render"
	"github.com/daniel/ppm/internal/repository"
)

type DashboardHandler struct {
	repo     *repository.DashboardRepo
	taskRepo *repository.TaskRepo
	render   *render.Renderer
}

func NewDashboardHandler(repo *repository.DashboardRepo, taskRepo *repository.TaskRepo, r *render.Renderer) *DashboardHandler {
	return &DashboardHandler{repo: repo, taskRepo: taskRepo, render: r}
}

func (h *DashboardHandler) Index(w http.ResponseWriter, r *http.Request) {
	data, err := h.repo.GetDashboard()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.render.Page(w, http.StatusOK, "dashboard.html", render.PageData{
		Title:   "Dashboard",
		Content: data,
	})
}

func (h *DashboardHandler) OverdueTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskRepo.ListOverdue()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.render.Partial(w, http.StatusOK, "dashboard-tasks", tasks)
}
