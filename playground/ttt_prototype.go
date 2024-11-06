// package main 
package playground  

import (
		"fmt"
		"tictactoe/tree"
	)
import "math/rand"
import "time"

type Piece struct {
	pieceType 	string 
}
type Position struct {
	rank int 
	file int 
}

var board [3][3] Piece 
var rankFileToArrayIndicesMapping = make(map[string]Position)

func main(){

	var bigTree *tree.Tree = tree.NewTree("root")
	bigTree.Root.AddChild("root child 1")
	// tree.PrintNodesBreadthFirst(bigTree.Root)


	fmt.Println("Welcome to triple T")
	printBoard()
	
	var input 			 string 
	var gameOver         bool = false
	var player1sTurn	 bool = true 

	var aGoChannel chan bool = make(chan bool)
	

	initializeRAI()

	for gameOver != true {
		if player1sTurn {
			_, err := fmt.Scanln(&input) //get user input 
			if err == nil && input != "" {
				switch input {
				case "a1":
					makeUserMove("a1", &player1sTurn)			
				case "a2":
					makeUserMove("a2", &player1sTurn)
				case "a3":
					makeUserMove("a3", &player1sTurn)
				case "b1":
					makeUserMove("b1", &player1sTurn)
				case "b2":
					makeUserMove("b2", &player1sTurn)
				case "b3":
					makeUserMove("b3", &player1sTurn)
				case "c1":
					makeUserMove("c1", &player1sTurn)
				case "c2":
					makeUserMove("c2", &player1sTurn)
				case "c3":
					makeUserMove("c3", &player1sTurn)
				}
			}	
		} else {
			makeComputerMove()
			player1sTurn = true
		}	
		printBoard()
		go evaluateBoard(aGoChannel)
		gameOver = <- aGoChannel
	}
}

func makeUserMove(aMove string, pP1Turn *bool) {
	
	var pos = rankFileToArrayIndicesMapping[aMove]

	if board[pos.rank][pos.file].pieceType == ""  {
		board[pos.rank][pos.file] = Piece{pieceType:"cross"}
		*pP1Turn = false 
	} else {
		fmt.Println("square taken")
		*pP1Turn = true 
	}

	if *pP1Turn == false {
		fmt.Println("User Moved")
		
	}
}

func makeComputerMove(){

	// create two random number generators seeded with time to make it random
	var seed 	= rand.NewSource(time.Now().UnixNano())
	var randomRank 	= rand.New(seed)
	var randomFile 	= rand.New(seed)

	var computerRankPosition = randomRank.Intn(3)
	var computerFilePosition = randomFile.Intn(3)

	for board[computerRankPosition][computerFilePosition].pieceType != "" {
		fmt.Println("Computer trying a position")
		computerRankPosition = randomRank.Intn(3)
		computerFilePosition = randomFile.Intn(3)
	}
	fmt.Println("Computer Move", computerRankPosition, computerFilePosition)
	board[computerRankPosition][computerFilePosition] = Piece{pieceType:"circle"}
}

func evaluateBoard(pGameOver chan bool){

	threeInARow := false 

	horizontalMatchCountCircle 		:= 0 
	verticalMatchCountCircle 		:= 0

	horizontalMatchCountCross 		:= 0 
	verticalMatchCountCross 		:= 0  

	for i:=0;i<3; i++ {
		for j:=0;j<3; j++{		
			if board[i][j].pieceType == "circle" {
				horizontalMatchCountCircle++
			}
			if horizontalMatchCountCircle == 3 {
				threeInARow = true
			}
		}
		if threeInARow == true {
			break 
		} else {
			horizontalMatchCountCircle = 0 	
		}
	}

	if threeInARow == false {
		for i:=0;i<3;i++ {
			for j:=0;j<3;j++{		
				if board[i][j].pieceType == "cross" {
					horizontalMatchCountCross++
				}
				if horizontalMatchCountCross == 3 {
					threeInARow = true
				}
			}
			if threeInARow == true {
				break 
			} else {
				horizontalMatchCountCross = 0 	
			}
		}
	}
	//vertical 
	if threeInARow == false {
		for i:=0;i<3;i++ {
			for j:=0;j<3;j++ {
				if board[j][i].pieceType == "circle" {
					verticalMatchCountCircle++ 
				}
				if verticalMatchCountCircle == 3 {
					threeInARow = true
				}
			}
			if threeInARow == true {
				break 
			} else {
				verticalMatchCountCircle = 0 
			}
		}
	}
	if threeInARow == false {
		for i:=0;i<3;i++ {
			for j:=0;j<3;j++ {
				if board[j][i].pieceType == "cross" {
					verticalMatchCountCross++ 
				}
				if verticalMatchCountCross == 3 {
					threeInARow = true
				}
			}
			if threeInARow == true {
				break 
			} else {
				verticalMatchCountCross = 0 
			}
		}
	}
	//diagonal - bottom left to top right 
	if threeInARow == false {
		
		a1 := board[2][0].pieceType
		b2 := board[1][1].pieceType
		c3 := board[0][2].pieceType

		if a1 == "circle" && b2 == "circle" && c3 == "circle" {
			threeInARow = true 
		} else if a1 == "cross" && b2 == "cross" && c3 == "cross" {
			threeInARow = true
		}
	}


	pGameOver <- threeInARow
	fmt.Println("Evaluating Board... Game Over? ",threeInARow)
}






// --- View / Visual / FE Code ---- 
func printBoard() {
	
	// fmt.Println("\033[H\033[2J")
	fmt.Println("--------")
	for i:=0; i<3; i++ {
		var rankPrint string = ""
		
		for j:=0; j<3; j++ {		
			switch board[i][j].pieceType {
			case "circle":
				rankPrint = rankPrint + "[o]"
			case "cross":
				rankPrint = rankPrint + "[x]"
			case "":
				rankPrint = rankPrint + "[ ]"
			}
		}
		fmt.Println(rankPrint)
	}
}

func initializeRAI(){
	rankFileToArrayIndicesMapping["a1"] =  Position{2,0}
	rankFileToArrayIndicesMapping["a2"] =  Position{2,1}
	rankFileToArrayIndicesMapping["a3"] =  Position{2,2}

	rankFileToArrayIndicesMapping["b1"] =  Position{1,0}
	rankFileToArrayIndicesMapping["b2"] =  Position{1,1}
	rankFileToArrayIndicesMapping["b3"] =  Position{1,2}

	rankFileToArrayIndicesMapping["c1"] =  Position{0,0}
	rankFileToArrayIndicesMapping["c2"] =  Position{0,1}
	rankFileToArrayIndicesMapping["c3"] =  Position{0,2}
}



/*
	- two players can place pieces on the board 


*/