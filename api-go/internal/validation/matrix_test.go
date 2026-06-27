package validation

import "testing"

func TestValidateRectangularMatrix(t *testing.T) {
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
			name: "rows less than columns is still rectangular",
			matrix: [][]float64{
				{1, 2, 3},
				{4, 5, 6},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateRectangularMatrix(test.matrix)
			if test.wantErr && err == nil {
				t.Fatalf("expected error but got nil")
			}
			if !test.wantErr && err != nil {
				t.Fatalf("expected nil error but got %v", err)
			}
		})
	}
}

func TestValidateQRMatrix(t *testing.T) {
	err := ValidateQRMatrix([][]float64{
		{1, 2, 3},
		{4, 5, 6},
	})

	if err == nil {
		t.Fatalf("expected error but got nil")
	}
}
