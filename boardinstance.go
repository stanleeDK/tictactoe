package main 

import (
	"strings"
	"strconv"
	// "fmt"
)

type BoardInstance struct {
	Board [3][3] 		byte  
	Html 				string 
	CellCount 			int 
	CurrentMove 		byte //signifies latest move just happened; so if x then the latest move was x and the next move is o's move
	CurrentTreeDepth 	int 
	Winner 				bool 
	BoardStateScore		float64 
	MiniMaxScore 	 	int 
	HumanPlayer			byte //capture which symbol the human is, x or o
	ComputerPlayer 		byte //capture which symbol the computer is, x or o
}

func NewBoardInstance() *BoardInstance{
	return &BoardInstance{
		Board: [3][3]byte {
			{'-','-','-'},
			{'-','-','-'},
			{'-','-','-'},
		},
		Html: "",
		CellCount: 0,
		Winner: false, 
		BoardStateScore: 0,
		MiniMaxScore: 0,
		CurrentMove: '-', //signifies latest move just happened	
		HumanPlayer: 'x',
		ComputerPlayer: 'o',
	}
}

func (b *BoardInstance) IsGameFinishedForOpposition(nextMoveOorX byte) bool {
	var opposition byte
	if nextMoveOorX == 'x' {
		opposition = 'o'
	} else {
		opposition = 'x'
	}
	return b.IsGameFinished(opposition)
}

func (b *BoardInstance) IsGameFinished(nextMoveOorX byte) bool {

	if b.DiagonalMatchEvaluation(nextMoveOorX) {
		return true 
	}
	if b.HorizontalMatchEvaluation(nextMoveOorX) {
		return true 
	}
	if b.VerticallMatchEvaluation(nextMoveOorX) {
		return true 
	}

	return false //no winner yet 
}

func (b *BoardInstance) HorizontalMatchEvaluation(nextMoveOorX byte) bool {
	var horizontalwin int = 0
	
	for i:=0;i<3;i++ {
		for j:=1;j<3;j++ {
			if b.Board[i][j-1] != b.Board[i][j] || b.Board[i][j] == '-' {
				break  
			} 
			if b.Board[i][j] == nextMoveOorX {
				horizontalwin++
			}
		}
		if horizontalwin == 2 {
			return true 
		} else {
			horizontalwin = 0 
		}
	}
	return false 
}

func (b *BoardInstance) VerticallMatchEvaluation(nextMoveOorX byte) bool {
	var verticalwin int = 0
	
	for i:=0;i<3;i++ {
		for j:=1;j<3;j++ {
			if b.Board[j-1][i] != b.Board[j][i] || b.Board[j][i] == '-' {
				break  
			} 
			if b.Board[j][i] == nextMoveOorX {
				verticalwin++
			}
		}
		if verticalwin == 2 {
			// fmt.Println(i)
			return true 
		} else {
			verticalwin = 0
		} 
	}
	return false 
}

func (b *BoardInstance) DiagonalMatchEvaluation(nextMoveOorX byte) bool {
	var diagonalwin int = 0

	var diagonal [3]byte 

	diagonal[0] = b.Board[0][0]
	diagonal[1] = b.Board[1][1]
	diagonal[2] = b.Board[2][2]

	for i:=1;i<3;i++ {
		if diagonal[i-1] != diagonal[i] || diagonal[i] == '-' {
			break
		}
		if diagonal[i] == nextMoveOorX {
			diagonalwin++
		}
	}
	
	if diagonalwin == 2 {
		return true 
	}

	// no wins in left right diagonal, so check right left 
	diagonalwin = 0 

	diagonal[0] = b.Board[0][2]
	diagonal[1] = b.Board[1][1]
	diagonal[2] = b.Board[2][0]

	for i:=1;i<3;i++ {
		if diagonal[i-1] != diagonal[i] || diagonal[i] == '-' {
			break
		}
		if diagonal[i] == nextMoveOorX {
			diagonalwin++
		}
	}
	
	if diagonalwin == 2 {
		return true 
	}
	return false 
}

// func (b *BoardInstance) IsBoardFull () bool {
// 	// var count int = 1 
// 	for i:=0;i<3;i++ {
// 		for j:=1;j<3;j++ {
// 			if b.Board[i][j] == '-' {
// 				return false 
// 			}
// 		}
// 	}
// 	return true 	
// }

func (b *BoardInstance) CreateHTMLforTableToFindNextMove(nodeIndex int, treeLevel int, isGameFinished bool,oppositionWon bool, boardstatescore int) {
	var writer strings.Builder 

	var nextMove string 

	if b.CurrentMove == 'x' {
		nextMove = "o"
	} else {
		nextMove = "x"
	}

	writer.WriteString("<div>nextMove:")
	writer.WriteString(nextMove)
	writer.WriteString("</div>")

	var MaxOrMin string 

	if treeLevel % 2 == 0 {
		MaxOrMin = "Maxi"
	} else {
		MaxOrMin = "Mini"
	}
	writer.WriteString("<div title=\"TESTING TOOL TIP\">MorM:")
	writer.WriteString(MaxOrMin)
	writer.WriteString("</div>")


	writer.WriteString("<div><b>i:")
	writer.WriteString(strconv.Itoa(nodeIndex))
	writer.WriteString("</b></div>")

	writer.WriteString("<span><b>l:")
	writer.WriteString(strconv.Itoa(treeLevel))
	writer.WriteString("</b></span>")

	writer.WriteString("</br><span><b>sc:")
	writer.WriteString(strconv.Itoa(boardstatescore))
	writer.WriteString("</b></span>")


	if isGameFinished == true && oppositionWon == false {
		writer.WriteString("<table class=\"winner\">")
	} else if oppositionWon == true {
		writer.WriteString("<table class=\"oppositionwon\">")	
	} else {
		writer.WriteString("<table>")
	}

	for i:=0;i<3;i++ {
		writer.WriteString("<tr>")		
		for j:=0;j<3;j++ {
			writer.WriteString("<td id=minimaxboard")
			writer.WriteString(strconv.Itoa(i))
			writer.WriteString(strconv.Itoa(j))
			writer.WriteString(">")
			switch b.Board[i][j] {
			case 'o':
				writer.WriteString("o")
			case 'x':
				writer.WriteString("x")
			case '-':
				writer.WriteString("-")
			}
			writer.WriteString("</td>")
		}
		writer.WriteString("</tr>")
	}

	writer.WriteString("</table>")
	b.Html = writer.String()
}

// func (b *BoardInstance) CreateGameBoardHtmlForActualGame () {
// 	var writer strings.Builder 

// 	for i:=0;i<3;i++ {
// 	writer.WriteString("<tr>")		
// 	for j:=0;j<3;j++ {
// 		writer.WriteString("<td>")
// 		switch b.Board[i][j] {
// 		case 'o':
// 			writer.WriteString("o")
// 		case 'x':
// 			writer.WriteString("x")
// 		case '-':
// 			writer.WriteString("-")
// 		}
// 		writer.WriteString("</td>")
// 	}
// 	writer.WriteString("</tr>")
// }

// 	writer.WriteString("</table>")
// 	b.Html = writer.String()
// }
