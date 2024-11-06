package main 

import (
	"fmt"
	// "tictactoe/boardinstance"

)

func printBoard(b BoardInstance) {
	fmt.Print("\n\n")
	for i:=0;i<3;i++ {
		fmt.Print("\t\t")
		for j:=0;j	<3;j++ {
			fmt.Print(b.Board[i][j]," ")
		}
		fmt.Print("\n")
	}
	fmt.Print("\n\n")
}

func isBoardFull(b BoardInstance) bool {
	var cellCount int = 0 
	for i:=0;i<3;i++ {
		for j:=0;j<3;j++ {
			if b.Board[i][j] == 'x' || b.Board[i][j] == 'o' {
				cellCount++
			}
		}
	}
	// fmt.Println(cellCount)
	if cellCount == 9 {
		return true 
	} else {
		return false 
	}
}