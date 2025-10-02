package types

type AssignmentsRaw struct {
	Assignment     string
	Grade          string
	AssignmentType string
}

type (
	ClassAndWeightGradesMap = map[string]map[int][]float64
	ClassAndGradeMap        = map[string]string
)
