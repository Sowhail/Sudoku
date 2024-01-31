package main

import (
	"fmt"
	"log"
)

type emptyCell struct {
	row            int
	col            int
	possibleValues []int
}

func validateByRow(table [9][9]int, row, col, value int) bool {
	for i := 0; i < len(table[row]); i++ {
		if i == col {
			continue
		}
		if table[row][i] == value {
			return false
		}
	}
	return true
}

func validateByCol(table [9][9]int, row, col, value int) bool {
	for i := 0; i < len(table[col]); i++ {
		if i == row {
			continue
		}
		if table[i][col] == value {
			return false
		}
	}
	return true
}

func setBox(first, last *int, pos int) {
	switch {
	case pos >= 0 && pos <= 2:
		*first = 0
		*last = 2
	case pos >= 3 && pos <= 5:
		*first = 3
		*last = 5
	case pos >= 6 && pos <= 8:
		*first = 6
		*last = 8
	}
}

func validateByBox(table [9][9]int, row, col, value int) bool {
	var rowStart, rowEnd, colStart, colEnd int
	setBox(&rowStart, &rowEnd, row)
	setBox(&colStart, &colEnd, col)
	for i := rowStart; i <= rowEnd; i++ {
		for j := colStart; j <= colEnd; j++ {
			if i == row && col == j {
				continue
			}
			if table[i][j] == value {
				return false
			}
		}
	}
	return true
}

func findPossibleValues(table [9][9]int, row, col int) ([]int, error) {
	res := []int{}
	for i := 1; i <= 9; i++ {
		isValid := validateByBox(table, row, col, i) && validateByCol(table, row, col, i) && validateByRow(table, row, col, i)
		if isValid {
			res = append(res, i)
		}
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("No valid value was found for position row=%v col=%v", row, col)
	}
	return res, nil
}

func main() {
	table := [9][9]int{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	}
	emptyCells := make([]emptyCell, 0, 81)
	// cellsPermutation := make([]emptyCell, 0)
	for i := 0; i < len(table); i++ {
		for j := 0; j < len(table[i]); j++ {
			if table[i][j] == 0 {
				values, err := findPossibleValues(table, i, j)
				if err != nil {
					log.Fatalf("invalid grid Sudoku beacause:\n%v", err.Error())
				}
				emptyCells = append(emptyCells, emptyCell{
					row:            i,
					col:            j,
					possibleValues: values,
				})
			}
		}
	}

	fmt.Println(emptyCells)

}
