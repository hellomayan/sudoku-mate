package sudoku

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
)

func TestSudoku(t *testing.T) {
	sudoku := inputParam()
	fmt.Println("The input sudoku is:")
	printArray2D(sudoku.problem)

	solutions := sudoku.FindSolutions()

	if len(solutions) == 0 {
		t.Fail()
	}

	fmt.Printf("Printing all the solutions, total=%d\n", len(solutions))
	for _, solution := range solutions {
		printArray2D(solution)
	}

}

func inputParam() *Sudoku {
	fileName := "testdata/problem_0_11_solutions.csv"
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
