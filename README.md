## sudokumore/solver
sudoku-mate is a library to find all solutions to a sudoku problem, although usually a sudoku only has one solution.
## How to use it
#### Input
- Use NewSudoku(input [9][9]uint8) or NewSudokuUnroll(input [81]uint8) to input your problem.
- Empty fields in the original problem should be filled with "0".

#### output
- FindSolutions() finds all the solutions to the sudoku problem