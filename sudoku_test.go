package sudoku

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
)

func TestFinalSolutions(t *testing.T) {
	fileName := "testdata/problem_0_11_solutions.csv"
	sudoku := inputParam(fileName)
	fmt.Println("The input sudoku is:")
	printArray2D(sudoku.problem)

	solutions := sudoku.FindSolutions()
	fmt.Printf("Printing all the solutions, total=%d\n", len(solutions))
	for _, solution := range solutions {
		printArray2D(solution)
	}

	if len(solutions) == 0 {
		t.Fatalf("There is a solution to all the test problems, but it failed to find.")
	}

	for _, solution := range solutions {
		if !isValidFinalSolution(solution) {
			t.Fatalf("wrong solution")
		}
	}
	t.Log("All solutions are valid, congratulations!")
}

func inputParam(fileName string) *Sudoku {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("fail to open file %s", fileName)
	}
	reader := csv.NewReader(f)
	reader.Comma = ','
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	arr := [81]uint8{}

	for i, record := range records {
		for j, s := range record {
			gold, err := strconv.ParseUint(s, 10, 8)
			if err != nil {
				log.Fatal(err)
			}
			arr[i*9+j] = uint8(gold)
		}
	}
	return NewSudokuUnroll(arr)
}

func isValidFinalSolution(sol [81]uint8) bool {
	s := &Sudoku{
		solution: sol,
	}
	for i := range sol {
		if s.isAllValid(i) {
			continue
		}
		return false
	}
	return true

}

func TestValidity(t *testing.T) {
	s := inputParam("testdata/solution_wrong_at_27.csv")
	s.solution = s.problem
	wrongPos := 27
	row := 4
	col := 0
	if s.isValidColumn(wrongPos) {
		t.Fatalf("number at this pos hould be invilid, row=%d, column=%d", row, col)
	}
	if s.isValidRow(wrongPos) {
		t.Fatalf("number at this pos hould be invilid, row=%d, column=%d", row, col)
	}
	if s.isValid3x3Square(wrongPos) {
		t.Fatalf("number at this pos hould be invilid, row=%d, column=%d", row, col)
	}
}
