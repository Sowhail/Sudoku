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

type regCell struct {
	row   int
	col   int
	value int
}

// ____________________validation____________________

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

func validateOneCell(table [9][9]int, row, col, value int) bool {
	return validateByBox(table, row, col, value) && validateByCol(table, row, col, value) && validateByRow(table, row, col, value)
}

func validateFullSudokuGrid(table [9][9]int) bool {
	for i := 0; i < len(table); i++ {
		for j := 0; j < len(table[i]); j++ {
			if table[i][j] == 0 {
				continue
			}
			isValidCell := validateOneCell(table, i, j, table[i][j])
			if !isValidCell {
				return false
			}
		}
	}
	return true
}

//________________________________________

// ____________________solvingSudoku____________________

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

func fillSudokuGrid(table [9][9]int, emptyCells []regCell) [9][9]int {
	for i := 0; i < len(emptyCells); i++ {
		table[emptyCells[i].row][emptyCells[i].col] = emptyCells[i].value
	}
	return table
}

func findValidTables(tables [][9][9]int, table [9][9]int, emptyCells []emptyCell, emptyCellsIndex int, tmp []regCell) [][9][9]int {
	if len(tmp) > 0 && !validateOneCell(table, tmp[len(tmp)-1].row, tmp[len(tmp)-1].col, tmp[len(tmp)-1].value) {
		return tables
	}
	if len(tmp) == len(emptyCells) {
		tempTable := fillSudokuGrid(table, tmp)
		tables = append(tables, tempTable)
		return tables
	}
	table = fillSudokuGrid(table, tmp)
	for j := 0; j < len(emptyCells[emptyCellsIndex].possibleValues); j++ {

		// if there were multiple answers only return the first and second answers and stop the progress
		if len(tables) > 1 {
			break
		}
		tables = findValidTables(tables, table, emptyCells, emptyCellsIndex+1, append(tmp, regCell{
			row:   emptyCells[emptyCellsIndex].row,
			col:   emptyCells[emptyCellsIndex].col,
			value: emptyCells[emptyCellsIndex].possibleValues[j],
		}))
	}
	return tables
}

func solveSudoku(table [9][9]int, emptyCells []emptyCell) ([][9][9]int, error) {
	tables := findValidTables([][9][9]int{}, table, emptyCells, 0, []regCell{})
	if len(tables) > 0 {
		return tables, nil
	}
	return [][9][9]int{}, fmt.Errorf("this Sudoku grid has no answer")
}

// ________________________________________

func printTable(table [9][9]int) {
	for i := 0; i < len(table); i++ {
		for j := 0; j < len(table[i]); j++ {
			fmt.Printf("%v ", table[i][j])
			if j == 2 || j == 5 {
				fmt.Printf("| ")
			}
		}
		fmt.Println()
		if i == 2 || i == 5 {
			fmt.Printf("----------------------\n")
		}
	}
}

func main() {
	table := [9][9]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 1},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	emptyCells := make([]emptyCell, 0, 81)
	// cellsPermutation := make([]emptyCell, 0)
	if !validateFullSudokuGrid(table) {
		log.Fatalf("invalid Sudoku gird")
	}
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
	tables, err := solveSudoku(table, emptyCells)
	if err != nil {
		log.Fatalf("invalid grid Sudoku because of:\n%v", err.Error())
	}

	//fmt.Println(ttmp)
	printTable(tables[0])
}
