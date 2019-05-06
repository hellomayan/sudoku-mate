package sudoku
import(
	"fmt"
	"github.com/scylladb/go-set/u8set"
)


func (su *Sudoku) isAllValid(pos int) bool {
	if !su.isValid3x3Square(pos) || !su.isValidRow(pos) || !su.isValidColumn(pos) {
		return false
	}
	return true
}

func (su *Sudoku) backwardNext(i int) int {
	su[i] = 0
	for step := 1; step <= i; step++ {
		if !su.IsPresetField(i - step) {
			return i - step
		}
	}
	return -1
}

func (su *Sudoku) forwardNext(i int) int {
	for step := 1; step < 81-i; step++ {
		if !su.IsPresetField(i + step) {
			return i + step
		}
	}
	return 81
}

//IsPresetField is used to check whether a field in the original problem is preset
func (su *Sudoku) IsPresetField(pos int) bool {
	if su[pos] > 0 {
		return true
	}
	return false
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
		gold := su[pos]
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

func (su *Sudoku) printArray2D() {
	for i, num := range su {
		fmt.Printf("%d ", num)
		if i%9 == 8 {
			fmt.Print("\n")
		}
	}
	fmt.Print("\n\n")
}

//NewSudoku is used to create a Sudoku problem
func NewSudoku(input [9][9]uint8)Sudoku{
	unrolled :=[81]uint8{}
	for i:=0;i<9;i++{
		for j:=0;j<9;j++{
			unrolled[i*9+j]=input[i][j]
		}
	}
	return unrolled
}

//Sudoku is a unrolled form of the sudoku problem
type Sudoku [81]uint8

//FindSolutions finds all the solutions of a sudoku problem
func (su *Sudoku)FindSolutions() []*Sudoku{
	
	solutions :=make([]*Sudoku,0,0)
	solution :=&Sudoku{}
	*solution=*su

	solution.printArray2D()

	for i := 0; i <= 80; {
		if i == -1 {
			fmt.Println("no more solutions")
			break
		}
		if su.IsPresetField(i) {
			i = solution.forwardNext(i)
			continue
		}
		foundAGold := false
		for maybeGold := solution[i] + 1; maybeGold <= uint8(9); maybeGold++ {
			solution[i] = maybeGold

			if !solution.isAllValid(i) {
				continue
			}
			foundAGold = true
			break
		}

		if foundAGold == false {
			i = solution.backwardNext(i)
			continue
		}

		temp := solution.forwardNext(i)
		if temp == 81 {
			solutions=append(solutions,solution)

			fmt.Printf("Found a new solution:\n")
			
			i = solution.backwardNext(i)
			continue
		}
		i = temp
		continue
	}
	
	return solutions
}