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
	originalProblem := inputParam()
	fmt.Println("The input sudoku is:")
	originalProblem.printArray2D()

	solutions := originalProblem.FindSolutions()

	if len(solutions) == 0 {
		t.Fail()
	}

	fmt.Printf("Printing all the solutions, total=%d\n", len(solutions))
	for _, solution := range solutions {
		solution.printArray2D()
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
	arr := &Sudoku{}

	for i, record := range records {
		for j, s := range record {
			gold, err := strconv.ParseUint(s, 10, 8)
			if err != nil {
				log.Fatal(err)
			}
			arr[i*9+j] = uint8(gold)
		}
	}
	return arr
}