package qr

import "testing"

const testTolerance = 1e-6

func TestFactorize(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]float64
	}{
		{
			name: "identity matrix",
			matrix: [][]float64{
				{1, 0},
				{0, 1},
			},
		},
		{
			name: "rectangular matrix",
			matrix: [][]float64{
				{1, 2},
				{3, 4},
				{5, 6},
			},
		},
		{
			name: "negative values",
			matrix: [][]float64{
				{-1, 2},
				{3, -4},
				{5, 6},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			q, r, err := Factorize(test.matrix)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			qtq := multiply(transpose(q), q)
			identity := identityMatrix(len(qtq))
			assertMatricesClose(t, qtq, identity)

			reconstructed := multiply(q, r)
			assertMatricesClose(t, reconstructed, test.matrix)
			assertUpperTriangular(t, r)
		})
	}
}

func TestFactorizeRejectsRowsLessThanColumns(t *testing.T) {
	_, _, err := Factorize([][]float64{
		{1, 2, 3},
		{4, 5, 6},
	})
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
}

func multiply(left, right [][]float64) [][]float64 {
	rows := len(left)
	columns := len(right[0])
	inner := len(right)
	result := make([][]float64, rows)
	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		result[rowIndex] = make([]float64, columns)
		for columnIndex := 0; columnIndex < columns; columnIndex++ {
			for innerIndex := 0; innerIndex < inner; innerIndex++ {
				result[rowIndex][columnIndex] += left[rowIndex][innerIndex] * right[innerIndex][columnIndex]
			}
		}
	}
	return result
}

func transpose(matrix [][]float64) [][]float64 {
	result := make([][]float64, len(matrix[0]))
	for rowIndex := range result {
		result[rowIndex] = make([]float64, len(matrix))
		for columnIndex := range matrix {
			result[rowIndex][columnIndex] = matrix[columnIndex][rowIndex]
		}
	}
	return result
}

func identityMatrix(size int) [][]float64 {
	result := make([][]float64, size)
	for rowIndex := 0; rowIndex < size; rowIndex++ {
		result[rowIndex] = make([]float64, size)
		result[rowIndex][rowIndex] = 1
	}
	return result
}

func assertMatricesClose(t *testing.T, actual, expected [][]float64) {
	t.Helper()
	if len(actual) != len(expected) {
		t.Fatalf("row count mismatch: got %d want %d", len(actual), len(expected))
	}

	for rowIndex := range actual {
		if len(actual[rowIndex]) != len(expected[rowIndex]) {
			t.Fatalf("column count mismatch at row %d: got %d want %d", rowIndex, len(actual[rowIndex]), len(expected[rowIndex]))
		}
		for columnIndex := range actual[rowIndex] {
			diff := actual[rowIndex][columnIndex] - expected[rowIndex][columnIndex]
			if diff < -testTolerance || diff > testTolerance {
				t.Fatalf("value mismatch at [%d][%d]: got %.10f want %.10f", rowIndex, columnIndex, actual[rowIndex][columnIndex], expected[rowIndex][columnIndex])
			}
		}
	}
}

func assertUpperTriangular(t *testing.T, matrix [][]float64) {
	t.Helper()
	for rowIndex := range matrix {
		for columnIndex := 0; columnIndex < rowIndex && columnIndex < len(matrix[rowIndex]); columnIndex++ {
			value := matrix[rowIndex][columnIndex]
			if value < -testTolerance || value > testTolerance {
				t.Fatalf("expected upper triangular matrix, found %.10f at [%d][%d]", value, rowIndex, columnIndex)
			}
		}
	}
}
