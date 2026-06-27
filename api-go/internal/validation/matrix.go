package validation

import (
	"fmt"
	"math"
)

type ValidationError struct {
	Details []string
}

func (e ValidationError) Error() string {
	return "matrix validation failed"
}

func ValidateMatrix(matrix [][]float64) error {
	if len(matrix) == 0 {
		return ValidationError{Details: []string{"matrix must not be empty"}}
	}

	if len(matrix[0]) == 0 {
		return ValidationError{Details: []string{"matrix rows must not be empty"}}
	}

	expectedColumns := len(matrix[0])
	details := make([]string, 0)

	for rowIndex, row := range matrix {
		if len(row) == 0 {
			details = append(details, fmt.Sprintf("row %d must not be empty", rowIndex))
			continue
		}

		if len(row) != expectedColumns {
			details = append(details, fmt.Sprintf("row %d has length %d but expected %d", rowIndex, len(row), expectedColumns))
		}

		for columnIndex, value := range row {
			if math.IsNaN(value) || math.IsInf(value, 0) {
				details = append(details, fmt.Sprintf("matrix[%d][%d] must be a finite number", rowIndex, columnIndex))
			}
		}
	}

	if len(matrix) < expectedColumns {
		details = append(details, fmt.Sprintf("matrix must satisfy rows >= columns, got %d rows and %d columns", len(matrix), expectedColumns))
	}

	if len(details) > 0 {
		return ValidationError{Details: details}
	}

	return nil
}
