package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/scylladb/go-set/u8set"
)

var (
	silent bool
	wg     sync.WaitGroup
)

func main() {
	originalProblem := inputParam()
	solution := originalProblem
	su := Sudoku{
		Original: originalProblem,
		Solution: solution,
	}
	fmt.Println("The input sudoku is:")
	printArray2D(su.Original)
	//su.simpleIteration()
	//wg := sync.WaitGroup{}
	core := runtime.NumCPU()
	//wg.Add(core)
	su.parallelSearch(core)
	//wg.Wait()
}

func isAllValid(square81 [81]uint8, pos int) bool {
	if !isValid3x3Square(square81, pos) || !isValidRow(square81, pos) || !isValidColumn(square81, pos) {
		return false
	}
	return true
}

func (su *Sudoku) backwardNext(i int) int {
	su.Solution[i] = 0
	for step := 1; step <= i; step++ {
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
	reader := csv.NewReader(os.Stdin)
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
	return arr
}

func isValidColumn(square81 [81]uint8, pos int) bool {
	colStart := pos % 9
	colPoses := [9]int{}
	for i := 0; i < 9; i++ {
		colPoses[i] = colStart + 9*i
	}
	return isValidSet(square81, colPoses)
}

func isValidRow(square81 [81]uint8, pos int) bool {
	rowStart := pos - pos%9
	rowPoses := [9]int{}
	for j := 0; j < 9; j++ {
		rowPoses[j] = rowStart + j
	}
	return isValidSet(square81, rowPoses)
}

func isValidSet(square81 [81]uint8, poses [9]int) bool {
	m := u8set.New()
	for _, pos := range poses {
		gold := square81[pos]
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

func isValid3x3Square(square81 [81]uint8, idx int) bool {
	squarePoses := find3x3Square(idx)
	return isValidSet(square81, squarePoses)
}

//Sudoku is the problem
type Sudoku struct {
	Solution [81]uint8
	Original [81]uint8
}

func printArray2D(square81 [81]uint8) {
	for i, num := range square81 {
		fmt.Printf("%d ", num)
		if i%9 == 8 {
			fmt.Print("\n")
		}
	}
	fmt.Print("\n\n")
}

func (su *Sudoku) simpleIteration() {
	foundTotal := 0
	for i := 0; i <= 80; {
		if i == -1 {
			fmt.Println("no more solutions")
			break
		}
		if su.isPresetField(i) {
			i = su.forwardNext(i)
			continue
		}
		foundAGold := false
		for maybeGold := su.Solution[i] + 1; maybeGold <= uint8(9); maybeGold++ {
			su.Solution[i] = maybeGold
			//su.printSolution(su.Solution)
			//fmt.Println()
			if !isAllValid(su.Solution, i) {
				continue
			}
			foundAGold = true
			break
		}

		if foundAGold == false {
			i = su.backwardNext(i)
			continue
		}

		temp := su.forwardNext(i)
		if temp == 81 {
			foundTotal++
			if !silent {
				fmt.Printf("Found the %d th solution:\n", foundTotal)
				printArray2D(su.Solution)
			}
			i = su.backwardNext(i)
			continue
		}
		i = temp
		continue
	}
	fmt.Printf("Found a total of %d solutions\n", foundTotal)
}

func (su *Sudoku) parallelSearch(core int) {
	taskChan := make(chan *task, 10000)
	su.newWorkers(taskChan, core)
	tsk := &task{
		solution:   su.Solution,
		currentPos: forwardNextPara(su.Solution),
	}
	taskChan <- tsk
}

type task struct {
	solution   [81]uint8
	currentPos int
}

type workers struct {
	ws []worker
}

func (su *Sudoku) newWorkers(taskChan chan *task, core int) {
	workers := make([]worker, core)
	for i, w := range workers {
		w.id = i
		w.taskChan = taskChan
		w.doTask()
	}
}

type worker struct {
	id       int
	taskChan chan *task
}

func (w *worker) doTask() {
	for {
		select {
		case tk := <-w.taskChan:
			fmt.Printf("worker %d receiving new task, current pos: %d", w.id, tk.currentPos)
			for i := uint8(1); i < 10; i++ {
				tk.solution[tk.currentPos] = i
				if isAllValid(tk.solution, tk.currentPos) {
					if forwardNextPara(tk.solution) == 81 {
						fmt.Printf("Found a solution:\n")
						printArray2D(tk.solution)
						break
					}
					newTK := task{
						solution:   tk.solution,
						currentPos: forwardNextPara(tk.solution),
					}
					w.taskChan <- &newTK
				}
			}
		default:
			break
		}
	}
}

func forwardNextPara(s [81]uint8) int {
	for i := 0; i <= 80; i++ {
		if s[i] == 0 {
			return i
		}
	}
	return 81
}
