package validation

import (
	"strings"
	"testing"
)

func TestValidateWeights(t *testing.T) {
	tests := []struct {
		name     string
		weights  []string
		wantErr  bool
		errMatch string
	}{
		{
			name:    "Valid weights that sum to 100",
			weights: []string{"50", "30", "20"},
			wantErr: false,
		},
		{
			name:     "Invalid weights that sum to less than 100",
			weights:  []string{"40", "40", "10"},
			wantErr:  true,
			errMatch: "weights should add to 100",
		},
		{
			name:     "Invalid weights that sum to more than 100",
			weights:  []string{"60", "30", "20"},
			wantErr:  true,
			errMatch: "weights should add to 100",
		},
		{
			name:     "Invalid conversion from string to int",
			weights:  []string{"a", "50", "50"},
			wantErr:  true,
			errMatch: "String to int conversion failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWeights(tt.weights)
			if tt.wantErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMatch) {
				t.Errorf("expected error to contain %q but got %q", tt.errMatch, err.Error())
			}
		})
	}
}

func TestValidatePoints(t *testing.T) {
	tests := []struct {
		name           string
		totalPoints    string
		correctPoints  string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:          "Valid case",
			totalPoints:   "100",
			correctPoints: "80",
			wantErr:       false,
		},
		{
			name:           "Correct points greater than total",
			totalPoints:    "50",
			correctPoints:  "60",
			wantErr:        true,
			expectedErrMsg: "Total points must be greater than or equal to correct points",
		},
		{
			name:           "Invalid totalPoints conversion",
			totalPoints:    "abc",
			correctPoints:  "10",
			wantErr:        true,
			expectedErrMsg: "String to int conversion failed",
		},
		{
			name:           "Invalid correctPoints conversion",
			totalPoints:    "100",
			correctPoints:  "xyz",
			wantErr:        true,
			expectedErrMsg: "String to int conversion failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePoints(tt.totalPoints, tt.correctPoints)

			if tt.wantErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.expectedErrMsg) {
				t.Errorf("error message mismatch: expected %q in %q", tt.expectedErrMsg, err.Error())
			}
		})
	}
}
