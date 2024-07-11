//go:build js && wasm

package sudoku

func safe(input *[dim][dim]int, i, j int) bool {
	for x := range dim {
		if input[x][j] == input[i][j] && x != i {
			return false
		}
	}
	for y := range dim {
		if input[i][y] == input[i][j] && j != y {
			return false
		}
	}
	for x := range 3 {
		for y := range 3 {
			i1 := i - i%3 + x
			j1 := j - j%3 + y
			if input[i1][j1] == input[i][j] && i1 != i && j1 != j {
				return false
			}
		}
	}
	return true
}

func validateSudoku(input *[dim][dim]int) bool {
	for i := range dim {
		for j := range dim {
			if !safe(input, i, j) {
				return false
			}
		}
	}
	return true
}

func fillSudoku(i, j int, input *[dim][dim]int) bool {
	if i == dim-1 && j == dim {
		return true
	}

	if j == dim {
		i += 1
		j = 0
	}

	for _, x := range nums {
		input[i][j] = x
		if safe(input, i, j) {
			if fillSudoku(i, j+1, input) {
				return true
			}
		}
		input[i][j] = 0
	}
	return false
}

func generateRandomSudoku() [dim][dim]int {
	out := [dim][dim]int{}
	fillSudoku(0, 0, &out)
	return out
}
