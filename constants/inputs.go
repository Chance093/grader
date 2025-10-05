package constants

const (
	VIEW_OVERALL_GRADES = "View Overall Grades"
	SELECT_CLASS        = "Select a Class"
	ADD_CLASS           = "Add a Class"
	EDIT_CLASS          = "Edit a Class"
	DELETE_CLASS        = "Delete a Class"
	VIEW_ASSIGNMENTS    = "View Assignments"
	ADD_ASSIGNMENT      = "Add Assignment"
	EDIT_ASSIGNMENT     = "Edit Assignment"
	DELETE_ASSIGNMENT   = "Delete Assignment"
	GO_BACK             = "Go Back"
	MAIN_MENU           = "Main Menu"
	QUIT                = "Quit"
)

var (
	EDIT_CLASS_OPTS      = []string{"Name", "Grade Weights"}
	EDIT_ASSIGNMENT_OPTS = []string{"Name", "Grade", "Type"}
	ASSIGNMENT_TYPES     = []string{"Test", "Quiz", "Homework"}
)
