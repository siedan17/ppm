package models

// Status constants for projects
const (
	ProjectActive    = "active"
	ProjectOnHold    = "on_hold"
	ProjectCompleted = "completed"
	ProjectArchived  = "archived"
)

// Status constants for tasks
const (
	TaskTodo       = "todo"
	TaskInProgress = "in_progress"
	TaskBlocked    = "blocked"
	TaskDone       = "done"
	TaskCancelled  = "cancelled"
)

// Category constants for tasks
const (
	CategoryProgramming     = "programming"
	CategoryDataEngineering = "data_engineering"
	CategorySpecification   = "specification"
	CategoryDesign          = "design"
	CategoryCommunication   = "communication"
	CategoryOther           = "other"
)

// Meeting type constants
const (
	MeetingInternal = "internal"
	MeetingExternal = "external"
)

var ProjectStatuses = []string{ProjectActive, ProjectOnHold, ProjectCompleted, ProjectArchived}
var TaskStatuses = []string{TaskTodo, TaskInProgress, TaskBlocked, TaskDone, TaskCancelled}
var TaskCategories = []string{CategoryProgramming, CategoryDataEngineering, CategorySpecification, CategoryDesign, CategoryCommunication, CategoryOther}
var MeetingTypes = []string{MeetingInternal, MeetingExternal}
var Priorities = []int{1, 2, 3, 4, 5}
