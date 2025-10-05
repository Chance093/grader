package types

type Assignment struct {
	Name  string
	Grade string
	Type  string
}

type Assignments = []Assignment

type (
	ClassAndWeightGradesMap = map[string]map[int][]float64
	ClassAndGradeMap        = map[string]string
)

type AssignmentWeight = struct {
	Weight  int
	Type_id int
}

type ClassAndGradeRaw = struct {
	ClassName string
	Grade     float64
	Weight    int
}

type ClassesAndGradesRaw = []ClassAndGradeRaw
