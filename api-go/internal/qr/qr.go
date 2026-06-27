package qr

import (
	"fmt"
	"math"
)

const Tolerance = 1e-9

func Factorize(matrix [][]float64) ([][]float64, [][]float64, error) {
	rows := len(matrix)
	if rows == 0 {
		return nil, nil, fmt.Errorf("matrix must not be empty")
	}

	columns := len(matrix[0])
	if rows < columns {
		return nil, nil, fmt.Errorf("matrix must satisfy rows >= columns")
	}

	q := makeMatrix(rows, columns)
	r := makeMatrix(columns, columns)
	vectors := cloneMatrix(matrix)

	for columnIndex := 0; columnIndex < columns; columnIndex++ {
		for previousColumnIndex := 0; previousColumnIndex < columnIndex; previousColumnIndex++ {
			projection := dotColumn(q, previousColumnIndex, vectors, columnIndex)
			r[previousColumnIndex][columnIndex] = projection

			for rowIndex := 0; rowIndex < rows; rowIndex++ {
				vectors[rowIndex][columnIndex] -= projection * q[rowIndex][previousColumnIndex]
			}
		}

		norm := columnNorm(vectors, columnIndex)
		if norm <= Tolerance {
			return nil, nil, fmt.Errorf("matrix columns are linearly dependent within tolerance %.1e", Tolerance)
		}

		r[columnIndex][columnIndex] = norm
		for rowIndex := 0; rowIndex < rows; rowIndex++ {
			q[rowIndex][columnIndex] = vectors[rowIndex][columnIndex] / norm
		}
	}

	return q, r, nil
}

func makeMatrix(rows, columns int) [][]float64 {
	result := make([][]float64, rows)
	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		result[rowIndex] = make([]float64, columns)
	}
	return result
}

func cloneMatrix(matrix [][]float64) [][]float64 {
	cloned := make([][]float64, len(matrix))
	for rowIndex := range matrix {
		cloned[rowIndex] = append([]float64(nil), matrix[rowIndex]...)
	}
	return cloned
}

func dotColumn(left [][]float64, leftColumn int, right [][]float64, rightColumn int) float64 {
	var result float64
	for rowIndex := range left {
		result += left[rowIndex][leftColumn] * right[rowIndex][rightColumn]
	}
	return result
}

func columnNorm(matrix [][]float64, column int) float64 {
	var sumSquares float64
	for rowIndex := range matrix {
		sumSquares += matrix[rowIndex][column] * matrix[rowIndex][column]
	}
	return math.Sqrt(sumSquares)
}
