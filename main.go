package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
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
		if i == col || table[row][i] == 0 {
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
		if i == row || table[i][col] == 0 {
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
			if (i == row && col == j) || table[i][j] == 0 {
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
	var res []int
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

func findEmptyCells(table [9][9]int) ([]emptyCell, error) {
	emptyCells := make([]emptyCell, 0, 81)
	for i := 0; i < len(table); i++ {
		for j := 0; j < len(table[i]); j++ {
			if table[i][j] == 0 {
				values, err := findPossibleValues(table, i, j)
				if err != nil {
					return nil, err
				}
				emptyCells = append(emptyCells, emptyCell{
					row:            i,
					col:            j,
					possibleValues: values,
				})
			}
		}
	}
	return emptyCells, nil
}

func findValidTables(tables [][9][9]int, table [9][9]int, emptyCells []emptyCell, emptyCellsIndex int, tmp []regCell, ctx context.Context) [][9][9]int {
	select {
	case <-ctx.Done():
		return tables
	default:
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
			}), ctx)
		}
		return tables
	}
}

func solveSudoku(table [9][9]int) ([][9][9]int, error) {
	if !validateFullSudokuGrid(table) {
		return nil, fmt.Errorf("invalid Sudoku gird")
	}
	emptyCells, err := findEmptyCells(table)
	if err != nil {
		return nil, err
	}
	// using time out in case if the table had no answer and took too long for backtracking
	ctx, cancel := context.WithTimeout(context.Background(), 1300*time.Millisecond)
	defer cancel()
	tables := findValidTables([][9][9]int{}, table, emptyCells, 0, []regCell{}, ctx)
	if len(tables) > 0 {
		return tables, nil
	}
	return nil, fmt.Errorf("this Sudoku grid has no answer")
}

// ________________________________________

func printTable(table [9][9]int) {
	for i := 0; i < len(table); i++ {
		for j := 0; j < len(table[i]); j++ {
			if table[i][j] == 0 {
				whiteText := "\033[31m"
				greenBackground := "\033[47m"
				reset := "\033[0m"
				fmt.Printf("%v%v%v%v ", whiteText, greenBackground, table[i][j], reset)
			} else {
				fmt.Printf("%v ", table[i][j])
			}
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

// _______________________GeneratingSudoku_______________________

func fillTableRandomly(table [9][9]int) [9][9]int {
	numberOfFilledCells := rand.Intn(6) + 10
	for i := 0; i < numberOfFilledCells; i++ {
		value := rand.Intn(9) + 1
		row := rand.Intn(9)
		col := rand.Intn(9)
		for !validateOneCell(table, row, col, value) {
			value = rand.Intn(9) + 1
			row = rand.Intn(9)
			col = rand.Intn(9)
		}
		table[row][col] = value
	}
	return table
}

func generateCompleteTable() [9][9]int {
	table := [9][9]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	table = fillTableRandomly(table)
	tables, err := solveSudoku(table)
	for err != nil {
		table = [9][9]int{}
		table = fillTableRandomly(table)
		tables, err = solveSudoku(table)
	}
	return tables[0]
}

func findFilledCells(table [9][9]int) []regCell {
	var res []regCell
	for i := 0; i < len(table); i++ {
		for j := 0; j < len(table[i]); j++ {
			if table[i][j] != 0 {
				res = append(res, regCell{
					row:   i,
					col:   j,
					value: table[i][j],
				})
			}
		}
	}
	return res
}

func makeIncompleteTable(table [9][9]int, numberOfWantedEmptyCells int, ctx context.Context) ([9][9]int, bool) {
	select {
	case <-ctx.Done():
		if tables, err := solveSudoku(table); len(tables) != 1 || err != nil {
			return table, false
		}
		return table, true
	default:
		if tables, err := solveSudoku(table); len(tables) != 1 || err != nil {
			return table, false
		}
		if numberOfWantedEmptyCells == 0 {
			return table, true
		}

		filledCells := findFilledCells(table)
		rnd := rand.Intn(len(filledCells))
		for table[filledCells[rnd].row][filledCells[rnd].col] == 0 {
			rnd = rand.Intn(len(filledCells))
		}
		value := table[filledCells[rnd].row][filledCells[rnd].col]
		table[filledCells[rnd].row][filledCells[rnd].col] = 0
		table, isCorrect := makeIncompleteTable(table, numberOfWantedEmptyCells-1, ctx)
		if !isCorrect {
			table[filledCells[rnd].row][filledCells[rnd].col] = value
			table, _ = makeIncompleteTable(table, numberOfWantedEmptyCells, ctx)
		}

		return table, true
	}
}

func main() {
	table := generateCompleteTable()

	ctx, cancel := context.WithTimeout(context.Background(), 1100*time.Millisecond)
	defer cancel()
	emptyCellsNumber := rand.Intn(16) + 45
	table, _ = makeIncompleteTable(table, emptyCellsNumber, ctx)
	ctr := 0
	for i := 0; i < len(table); i++ {
		for j := 0; j < len(table[i]); j++ {
			if table[i][j] == 0 {
				ctr++
			}
		}
	}
	fmt.Printf("\nthe number of emptyCells: %v \n", ctr)
	fmt.Printf("question table:\n")
	printTable(table)
	fmt.Println()

	tables, err := solveSudoku(table)
	if err != nil {
		log.Fatalf("invalid grid Sudoku\n error: %v", err.Error())
	}
	fmt.Printf("\n\nAnswer table:\n")
	printTable(tables[0])
}
