package http

import (
	"net/http"

	"github.com/daniel/ppm/internal/application"
	"github.com/daniel/ppm/internal/delivery/render"
)

type DashboardHandler struct {
	svc    *application.DashboardService
	render *render.Renderer
}

func NewDashboardHandler(svc *application.DashboardService, r *render.Renderer) *DashboardHandler {
	return &DashboardHandler{svc: svc, render: r}
}

func (h *DashboardHandler) Index(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GetDashboard()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.render.Page(w, http.StatusOK, "dashboard.html", render.PageData{
		Title:   "Dashboard",
		Content: data,
	})
}

func (h *DashboardHandler) ActiveTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.svc.ListActiveTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.render.Partial(w, http.StatusOK, "dashboard-tasks", tasks)
}
