//go:build js && wasm

package main

import "github.com/2asm/daily_sudoku/sudoku"

func main() {
	sudoku.NewGame().Start()
}
