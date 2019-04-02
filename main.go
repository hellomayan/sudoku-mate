package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/scylladb/go-set/iset"
	"github.com/scylladb/go-set/u8set"
)

func main() {
	originalProblem := inputParam()
	solution := originalProblem
	su := Sudoku{
		Original: originalProblem,
		Solution: solution,
	}

	//startNums is used to remember the next gold(number) to try
	startNums := [81]uint8{}
	for i := 0; i < 81; i++ {
		if su.Original[i] == 0 {
			startNums[i] = 1
		} else {
			startNums[i] = 10
		}
	}

	for i := 0; i >= 0 && i < 81; {
		if startNums[i] > 9 {
			i = su.forwardNext(i)
			continue
		}
		for maybeGold := startNums[i]; maybeGold <= uint8(9); maybeGold++ {
			su.Solution[i] = maybeGold
			if !su.isValid3x3Square(i) || !su.isValidRow(i) || !su.isValidColumn(i) {
				su.Solution[i] = 0
				continue
			}
			//found a valid maybeGold
			if su.Solution[i] > 0 {
				startNums[i] = startNums[i] + 1
				i = su.forwardNext(i)
				break
			}
		}
		i = su.backwardNext(i)
	}

}

func (su *Sudoku) backwardNext(i int) int {
	for step := 1; step < i; step++ {
		if su.Original[i-step] == 0 {
			return i - step
		}
	}
	return -1
}

func (su *Sudoku) forwardNext(i int) int {
	for step := 1; step < 81-i; step++ {
		if su.Original[i+step] == 0 {
			return i + step
		}
	}
	return 81
}

func inputParam() [81]uint8 {
	name := "problem.csv"
	f, _ := os.Open(name)
	defer f.Close()

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
	log.Print(arr)
	return arr
}

func (su *Sudoku) isValidColumn(pos int) bool {
	colStart := pos % 9
	colPoses := [9]int{}
	for i := 0; i < 9; i++ {
		colPoses[i] = colStart + 9*i
	}
	return su.isValidSet(colPoses)
}

func (su *Sudoku) isValidRow(pos int) bool {
	rowStart := pos - pos%9
	rowPoses := [9]int{}
	for j := 0; j < 9; j++ {
		rowPoses[j] = rowStart + j
	}
	return su.isValidSet(rowPoses)
}

func (su *Sudoku) isValidSet(poses [9]int) bool {
	m := u8set.New()
	for _, pos := range poses {
		gold := su.Solution[pos]
		if gold == 0 {
			continue
		}
		if m.Has(gold) {
			return false
		}
		m.Add(gold)
	}
	return true
}

func isValid9squares(rows [][]int) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if !isValidSquare(rows, 3*i, 3*j) {
				return false
			}
		}
	}
	return true
}

func isValidSquare(rows [][]int, row int, col int) bool {
	set := iset.New()
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			val := rows[row+i][col+j]
			if val == 0 {
				continue
			}
			if set.Has(val) {
				return false
			}
			set.Add(val)
		}
	}
	return true
}

func isValidColumns(rows [][]int) bool {
	for j := 0; j < 9; j++ {
		if !isValidOneColumn(rows, j) {
			return false
		}
	}
	return true
}

func isValidOneColumn(rows [][]int, col int) bool {
	set := iset.New()
	for i := 0; i < 9; i++ {
		if rows[i][col] == 0 {
			continue
		}
		if !set.Has(rows[i][col]) {
			set.Add(rows[i][col])
			continue
		}
		return false
	}
	return true
}

func find3x3Square(idx int) [9]int {
	a := idx / 9
	b := idx % 9
	c := a % 3
	d := b % 3
	start := idx - 9*c - d
	arr := [...]int{start, start + 1, start + 2, start + 9, start + 10, start + 11, start + 18, start + 19, start + 20}
	return arr
}

func (su *Sudoku) isValid3x3Square(idx int) bool {
	squarePoses := find3x3Square(idx)
	return su.isValidSet(squarePoses)
}

//Sudoku is the problem
type Sudoku struct {
	Solution [81]uint8
	Original [81]uint8
}
