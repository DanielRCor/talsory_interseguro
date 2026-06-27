package validation

import "testing"

func TestValidateMatrix(t *testing.T) {
	tests := []struct {
		name    string
		matrix  [][]float64
		wantErr bool
	}{
		{
			name: "valid rectangular matrix",
			matrix: [][]float64{
				{1, 2},
				{3, 4},
				{5, 6},
			},
		},
		{
			name:    "empty matrix",
			matrix:  [][]float64{},
			wantErr: true,
		},
		{
			name: "non rectangular matrix",
			matrix: [][]float64{
				{1, 2},
				{3},
			},
			wantErr: true,
		},
		{
			name: "rows less than columns",
			matrix: [][]float64{
				{1, 2, 3},
				{4, 5, 6},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateMatrix(test.matrix)
			if test.wantErr && err == nil {
				t.Fatalf("expected error but got nil")
			}
			if !test.wantErr && err != nil {
				t.Fatalf("expected nil error but got %v", err)
			}
		})
	}
}
