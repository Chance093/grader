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
