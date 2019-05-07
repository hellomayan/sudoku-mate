package sudoku

import (
	"fmt"

	"github.com/scylladb/go-set/u8set"
)

func (s *Sudoku) isAllValid(pos int) bool {
	if !s.isValid3x3Square(pos) || !s.isValidRow(pos) || !s.isValidColumn(pos) {
		return false
	}
	return true
}

func (s *Sudoku) backwardNext(i int) int {
	s.solution[i] = 0
	for step := 1; step <= i; step++ {
		if !s.isPresetField(i - step) {
			return i - step
		}
	}
	return -1
}

func (s *Sudoku) forwardNext(i int) int {
	for step := 1; step < 81-i; step++ {
		if !s.isPresetField(i + step) {
			return i + step
		}
	}
	return 81
}

func (s *Sudoku) isPresetField(pos int) bool {
	if s.problem[pos] > 0 {
		return true
	}
	return false
}

func (s *Sudoku) isValidColumn(pos int) bool {
	colStart := pos % 9
	colPoses := [9]int{}
	for i := 0; i < 9; i++ {
		colPoses[i] = colStart + 9*i
	}
	return s.isValidSet(colPoses)
}

func (s *Sudoku) isValidRow(pos int) bool {
	rowStart := pos - pos%9
	rowPoses := [9]int{}
	for j := 0; j < 9; j++ {
		rowPoses[j] = rowStart + j
	}
	return s.isValidSet(rowPoses)
}

func (s *Sudoku) isValidSet(poses [9]int) bool {
	m := u8set.New()
	for _, pos := range poses {
		gold := s.solution[pos]
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

func (s *Sudoku) isValid3x3Square(idx int) bool {
	squarePoses := find3x3Square(idx)
	return s.isValidSet(squarePoses)
}

func printArray2D(arr [81]uint8) {
	for i, num := range arr {
		fmt.Printf("%d ", num)
		if i%9 == 8 {
			fmt.Print("\n")
		}
	}
	fmt.Print("\n\n")
}

//NewSudoku is used to create a Sudoku problem
func NewSudoku(input [9][9]uint8) *Sudoku {
	unrolled := [81]uint8{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			unrolled[i*9+j] = input[i][j]
		}
	}
	return &Sudoku{
		problem:  unrolled,
		solution: unrolled,
	}
}

//NewSudokuUnroll is used to create a Sudoku problem with a already unrolled input
func NewSudokuUnroll(input [81]uint8) *Sudoku {
	return &Sudoku{
		problem:  input,
		solution: input,
	}
}

//Sudoku is a unrolled form of the sudoku problem
//solution is a variable that keeps changing during the process of searching for a solution
type Sudoku struct {
	problem  [81]uint8
	solution [81]uint8
}

//FindSolutions returns all the solutions of a sudoku problem
func (s *Sudoku) FindSolutions() [][81]uint8 {

	solutions := make([][81]uint8, 0, 0)

	for i := 0; i <= 80; {
		if i == -1 {
			//fmt.Println("no more solutions")
			break
		}
		if s.isPresetField(i) {
			i = s.forwardNext(i)
			continue
		}
		foundAGold := false
		for maybeGold := s.solution[i] + 1; maybeGold <= uint8(9); maybeGold++ {
			s.solution[i] = maybeGold

			if !s.isAllValid(i) {
				continue
			}
			foundAGold = true
			break
		}

		if foundAGold == false {
			i = s.backwardNext(i)
			continue
		}

		temp := s.forwardNext(i)
		if temp == 81 {
			solutions = append(solutions, s.solution)

			//fmt.Printf("Found a new solution:\n")
			//solution.printArray2D()
			i = s.backwardNext(i)
			continue
		}
		i = temp
	}

	return solutions
}
