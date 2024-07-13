//go:build js && wasm

package sudoku

import (
	"fmt"
	"math/rand"
	"strconv"
	"syscall/js"
	"time"
)

const dim = 9

var sand *rand.Rand // seeded rand
var nums = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
var seed int

func init() {
	y, m, d := time.Now().Date()
	seed = (y << 9) + int(m<<5) + d // daily seed
	sand = rand.New(rand.NewSource(int64(seed)))
	for i := dim - 1; i >= 1; i -= 1 { // random shuffle
		j := sand.Intn(i)
		nums[i], nums[j] = nums[j], nums[i]
	}
}

type coord struct {
	x, y int
}

type game struct {
	cell              [dim][dim]int
	possible_solution [dim][dim]int
}

func NewGame() *game {
	g := &game{}
	g.fillDailySudoku()
	return g
}

func (g *game) cellValid(x, y int) bool {
	return g.cell[x][y] > 0 && g.cell[x][y] <= dim
}

func (g *game) fillDailySudoku() {
	g.cell = generateRandomSudoku()
	g.possible_solution = g.cell
	removals := 0
	for removals > 60 || removals < 10 {
		removals = sand.Intn(dim*dim) + 1
	}
	for removals > 0 {
		i := sand.Intn(dim)
		j := sand.Intn(dim)
		if g.cell[i][j] != 0 {
			g.cell[i][j] = 0
			removals -= 1
		}
	}
}

func (g *game) checkSolution(input [dim][dim]int) bool {
	if input == g.possible_solution {
		return true
	}
	for i := range dim {
		for j := range dim {
			if input[i][j] < 1 || input[i][j] > 9 {
				return false
			}
			if g.cell[i][j] == 0 {
				continue
			}
			if g.cell[i][j] != input[i][j] {
				return false
			}
		}
	}
	return validateSudoku(&input)
}

func arrayToString(input [dim][dim]int) string {
	res := ""
	for i := range dim {
		for j := range dim {
			res += fmt.Sprint(input[i][j])
		}
	}
	return res
}

func stringToArray(input string) *[dim][dim]int {
	out := &[dim][dim]int{}
	for i, ch := range input {
		out[i/dim][i%dim] = int(ch - '0')
	}
	return out
}

func (g *game) Start() {
	check.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		res := getSudokuSolutionFromHTML()
		ok := g.checkSolution(res)
		g.updateResult(ok)
		if ok {
			js.Global().Get("localStorage").Call("setItem", seed, arrayToString(res))
		}
		return nil
	}))
	ans := js.Global().Get("localStorage").Call("getItem", seed)
	var out *[dim][dim]int = nil
	if !ans.IsNull() {
		out = stringToArray(ans.String())
	} else {
		js.Global().Get("localStorage").Call("clear")
	}
	for i := range dim {
		for j := range dim {
			grid[i][j].Set("value", g.cell[i][j])
			if g.cell[i][j] != 0 {
				grid[i][j].Set("innerText", g.cell[i][j])
				grid[i][j].Set("disabled", true)
				grid[i][j].Call("setAttribute", "style", "color:grey;")
			} else {
				if out != nil {
					grid[i][j].Set("innerText", out[i][j])
					grid[i][j].Set("value", out[i][j])
                    grid[i][j].Call("setAttribute", "style", "color:grey;background:#eee;")
				} else {
					grid[i][j].Set("innerText", "0")
                    grid[i][j].Call("setAttribute", "style", "color:transparent;background:#eee;")
				}
			}
		}
	}
	if out != nil {
		g.updateResult(true)
	}
	select {}
}

func (g *game) updateResult(ok bool) {
	if ok {
		result.Set("innerText", "Congratulations, you solved today's sudoku")
		result.Set("style", "color:green;")
	} else {
		result.Set("innerText", "Wrong")
		result.Set("style", "color:orange;")
		go func() {
			time.Sleep(time.Second * 3)
			result.Set("innerText", "")
		}()
	}
}

var (
	grid          [dim][dim]js.Value
	result, check js.Value
)

func coordToId(c coord) string {
	return fmt.Sprintf("i%v%v", c.x, c.y)
}

func idToCoord(id string) coord {
	x := int(id[1] - '0')
	y := int(id[2] - '0')
	if x <= 0 || x > dim {
		panic("id to coord")
	}
	if y <= 0 || y > dim {
		panic("id to coord")
	}
	out := coord{x, y}
	return out
}

func init() {
	for i := range dim {
		for j := range dim {
			grid[i][j] = js.Global().Get("document").Call("getElementById", coordToId(coord{i, j}))
		}
	}
	result = js.Global().Get("document").Call("getElementById", "result")
	check = js.Global().Get("document").Call("getElementById", "check")
}

func getSudokuSolutionFromHTML() [dim][dim]int {
	out := [dim][dim]int{}
	for i := range dim {
		for j := range dim {
			val, err := strconv.Atoi(grid[i][j].Get("value").String())
			if err != nil {
				panic("atoi")
			}
			out[i][j] = val
		}
	}
	return out
}
