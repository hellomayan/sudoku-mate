package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/scylladb/go-set/u8set"
)

func main() {
	originalProblem := inputParam()
	solution := originalProblem
	su := Sudoku{
		Original: originalProblem,
		Solution: solution,
	}

	for i := 0; i >= 0 && i < 81; {
		if su.isPresetField(i) {
			i = su.forwardNext(i)
			continue
		}
		foundAGold := false
		for maybeGold := su.Solution[i] + 1; maybeGold <= uint8(9); maybeGold++ {
			su.Solution[i] = maybeGold
			su.printSolution(su.Solution)
			fmt.Println()
			if !su.isAllValid(i) {
				continue
			}
			foundAGold = true
			break
		}
		if foundAGold == false {
			i = su.backwardNext(i)
		} else {

			i = su.forwardNext(i)
		}
	}
	su.printSolution(su.Solution)
}

func (su *Sudoku) isAllValid(pos int) bool {
	if !su.isValid3x3Square(pos) || !su.isValidRow(pos) || !su.isValidColumn(pos) {
		return false
	}
	return true
}

func (su *Sudoku) backwardNext(i int) int {
	su.Solution[i] = 0
	for step := 1; step < i; step++ {
		if !su.isPresetField(i - step) {
			return i - step
		}
	}
	return -1
}

func (su *Sudoku) forwardNext(i int) int {
	for step := 1; step < 81-i; step++ {
		if !su.isPresetField(i + step) {
			return i + step
		}
	}
	return 81
}

func (su *Sudoku) isPresetField(pos int) bool {
	if su.Original[pos] > 0 {
		return true
	}
	return false
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

func (su *Sudoku) printSolution(solution [81]uint8) {
	for i, num := range solution {
		fmt.Print(num)
		if i%9 == 8 {
			fmt.Println()
		}
	}
}
