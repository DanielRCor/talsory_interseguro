package rotation

import "testing"

func TestRotateClockwise(t *testing.T) {
	result, direction, err := Rotate([][]float64{
		{1, 2, 3},
		{4, 5, 6},
	}, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := [][]float64{
		{4, 1},
		{5, 2},
		{6, 3},
	}

	if direction != DirectionClockwise {
		t.Fatalf("expected clockwise direction, got %s", direction)
	}

	assertRotationEqual(t, result, expected)
}

func TestRotateCounterClockwise(t *testing.T) {
	result, _, err := Rotate([][]float64{
		{1, 2, 3},
		{4, 5, 6},
	}, DirectionCounterClockwise)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := [][]float64{
		{3, 6},
		{2, 5},
		{1, 4},
	}

	assertRotationEqual(t, result, expected)
}

func TestRotateRejectsInvalidDirection(t *testing.T) {
	_, _, err := Rotate([][]float64{{1}}, "diagonal")
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
}

func assertRotationEqual(t *testing.T, actual, expected [][]float64) {
	t.Helper()
	if len(actual) != len(expected) {
		t.Fatalf("row count mismatch: got %d want %d", len(actual), len(expected))
	}
	for rowIndex := range actual {
		if len(actual[rowIndex]) != len(expected[rowIndex]) {
			t.Fatalf("column count mismatch at row %d", rowIndex)
		}
		for columnIndex := range actual[rowIndex] {
			if actual[rowIndex][columnIndex] != expected[rowIndex][columnIndex] {
				t.Fatalf("value mismatch at [%d][%d]: got %.2f want %.2f", rowIndex, columnIndex, actual[rowIndex][columnIndex], expected[rowIndex][columnIndex])
			}
		}
	}
}
