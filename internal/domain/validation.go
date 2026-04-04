package domain

func ValidateProject(p *Project) map[string]string {
	errors := make(map[string]string)
	if p.Name == "" {
		errors["name"] = "Name is required"
	}
	if p.StartDate == "" {
		errors["start_date"] = "Start date is required"
	}
	if p.Priority < 1 || p.Priority > 5 {
		errors["priority"] = "Priority must be between 1 and 5"
	}
	return errors
}

func ValidateTask(t *Task) map[string]string {
	errors := make(map[string]string)
	if t.Title == "" {
		errors["title"] = "Title is required"
	}
	if t.ProjectID == 0 {
		errors["project_id"] = "Project is required"
	}
	if t.StartDate == "" {
		errors["start_date"] = "Start date is required"
	}
	if t.Deadline == "" {
		errors["deadline"] = "Deadline is required"
	}
	if t.EstimatedHours <= 0 {
		errors["estimated_hours"] = "Estimated hours must be greater than 0"
	}
	return errors
}

func ValidateMeeting(m *Meeting) map[string]string {
	errors := make(map[string]string)
	if m.Title == "" {
		errors["title"] = "Title is required"
	}
	if m.Date == "" {
		errors["date"] = "Date is required"
	}
	if m.ProjectID == 0 {
		errors["project_id"] = "Project is required"
	}
	return errors
}
