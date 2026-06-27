package rotation

import "fmt"

const (
	DirectionClockwise        = "clockwise"
	DirectionCounterClockwise = "counterclockwise"
)

func Rotate(matrix [][]float64, direction string) ([][]float64, string, error) {
	if direction == "" {
		direction = DirectionClockwise
	}

	switch direction {
	case DirectionClockwise:
		return rotateClockwise(matrix), direction, nil
	case DirectionCounterClockwise:
		return rotateCounterClockwise(matrix), direction, nil
	default:
		return nil, "", fmt.Errorf("direction must be one of: clockwise, counterclockwise")
	}
}

func rotateClockwise(matrix [][]float64) [][]float64 {
	rows := len(matrix)
	columns := len(matrix[0])
	result := makeRotatedMatrix(columns, rows)
	for columnIndex := 0; columnIndex < columns; columnIndex++ {
		for rowIndex := 0; rowIndex < rows; rowIndex++ {
			result[columnIndex][rows-1-rowIndex] = matrix[rowIndex][columnIndex]
		}
	}
	return result
}

func rotateCounterClockwise(matrix [][]float64) [][]float64 {
	rows := len(matrix)
	columns := len(matrix[0])
	result := makeRotatedMatrix(columns, rows)
	for columnIndex := 0; columnIndex < columns; columnIndex++ {
		for rowIndex := 0; rowIndex < rows; rowIndex++ {
			result[columns-1-columnIndex][rowIndex] = matrix[rowIndex][columnIndex]
		}
	}
	return result
}

func makeRotatedMatrix(rows, columns int) [][]float64 {
	result := make([][]float64, rows)
	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		result[rowIndex] = make([]float64, columns)
	}
	return result
}
