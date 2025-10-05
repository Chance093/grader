package calculate

import (
	"reflect"
	"testing"

	"github.com/Chance093/grader/types"
)

func TestGetClassAndWeightGradesMap(t *testing.T) {
	raw := types.ClassesAndGradesRaw{
		{ClassName: "Math", Weight: 40, Grade: 90.0},
		{ClassName: "Math", Weight: 40, Grade: 80.0},
		{ClassName: "Math", Weight: 60, Grade: 70.0},
		{ClassName: "History", Weight: 100, Grade: 85.0},
	}

	got := getClassAndWeightGradesMap(raw)
	want := types.ClassAndWeightGradesMap{
		"Math": map[int][]float64{
			40: {90.0, 80.0},
			60: {70.0},
		},
		"History": map[int][]float64{
			100: {85.0},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("getClassAndWeightGradesMap() = %#v, want %#v", got, want)
	}
}

func TestGetClassAndGradeMap(t *testing.T) {
	classAndWeightGradesMap := types.ClassAndWeightGradesMap{
		"Math": map[int][]float64{
			40: {90.0, 80.0}, // avg 85
			60: {70.0},       // avg 70
		},
	}

	got := getClassAndGradeMap(classAndWeightGradesMap)
	want := types.ClassAndGradeMap{
		// Math: (85 * 0.4) + (70 * 0.6) = 34 + 42 = 76.0
		"Math": "76.0",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("getClassAndGradeMap() = %#v, want %#v", got, want)
	}
}

func TestCalculateGrades(t *testing.T) {
	raw := types.ClassesAndGradesRaw{
		{ClassName: "Math", Weight: 40, Grade: 100.0},
		{ClassName: "Math", Weight: 60, Grade: 80.0},

		{ClassName: "Science", Weight: 50, Grade: 70.0},
		{ClassName: "Science", Weight: 50, Grade: 90.0},
	}

	got := CalculateGrades(raw)

	// Math: (100 * .4) + (80 * .6) = 40 + 48 = 88
	// Science: avg(70) * .5 + avg(90) * .5 = (70*0.5 + 90*0.5) = 80
	want := types.ClassAndGradeMap{
		"Math":    "88.0",
		"Science": "80.0",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CalculateGrades() = %#v, want %#v", got, want)
	}
}
