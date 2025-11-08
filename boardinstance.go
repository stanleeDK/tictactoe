package main 

import (
	"strings"
	"strconv"
	"fmt"
)

type GameState int 

const (
	Draw 				GameState = 0
	ComputerIsOAndWon 	GameState = 1 
	HumanisXandWon 		GameState = 2
	GameNotFinished 	GameState = 3 
)

type BoardInstance struct {
	Board [3][3] 		byte  
	Html 				string 
	CellCount 			int 
	MostRecentPlayer 	byte //signifies latest move just happened; so if player played most recently then this is x and the next move is o's move
	CurrentTreeDepth 	int 
	Winner 				bool 
	BoardStateScore		int 
	MiniMaxScore 	 	int 
	HumanPlayer			byte //capture which symbol the human is, x or o
	ComputerPlayer 		byte //capture which symbol the computer is, x or o
	CurrentGameState 	GameState
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
		MostRecentPlayer: '-', //signifies latest move just happened	
		HumanPlayer: 'x',
		ComputerPlayer: 'o',
		CurrentGameState: GameNotFinished,
	}
}

func (b *BoardInstance) ReturnEvaluationOfGameBoard () {

	if b.HorizontalMatchEvaluation('o') == true || b.VerticallMatchEvaluation('o') == true || b.DiagonalMatchEvaluation('o') == true { 
		b.CurrentGameState = ComputerIsOAndWon
		return
	}
	if b.HorizontalMatchEvaluation('x') == true || b.VerticallMatchEvaluation('x') == true || b.DiagonalMatchEvaluation('x') == true { 
		b.CurrentGameState = HumanisXandWon
		return
	}
	if b.IsMatchADraw() == true {
		b.CurrentGameState = Draw
		return
	}

	b.CurrentGameState = GameNotFinished
}



//finish draw function
func (b *BoardInstance) IsMatchADraw() bool {

	if b.HorizontalMatchEvaluation('o') == false || b.VerticallMatchEvaluation('o') == false || b.DiagonalMatchEvaluation('o') == false {
		// fmt.Println("o didn;t win")
		if b.HorizontalMatchEvaluation('x') == false || b.VerticallMatchEvaluation('x') == false || b.DiagonalMatchEvaluation('x') == false {
			// fmt.Println("x didn;t win")
			if b.IsBoardFull() == true {
				// fmt.Println("board is full")
				return true 
			}
		}
	}

	return false 

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

func (b *BoardInstance) IsBoardFull () bool {
	// var count int = 1 
	for i:=0;i<3;i++ {
		for j:=0;j<3;j++ {
			// fmt.Println(string(b.Board[i][j]))
			if b.Board[i][j] == '-' {
				return false 
			}
		}
	}
	return true 	
}
func (b *BoardInstance) CreateHTMLforTableToRenderUserFacingBoardState(){
	var writer strings.Builder 
	// fmt.Println("game state", b.CurrentGameState)
	
	switch b.CurrentGameState {
	case Draw: 
		writer.WriteString("<table class=\"draw\">")
	case ComputerIsOAndWon:
		writer.WriteString("<table class=\"computerisIandWon\">")
	case HumanisXandWon:
		writer.WriteString("<table class=\"humanisXandWon\">")	
	case GameNotFinished:
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
	//fmt.Println(writer.String())
	b.Html = ""
	b.Html = writer.String()

}
func (b *BoardInstance) CreateHTMLforTableToFindNextMove(nodeIndex int, treeLevel int) {
	var writer strings.Builder 

	var nextMove string 

	if b.MostRecentPlayer == 'x' {
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
	writer.WriteString(strconv.Itoa(b.BoardStateScore))
	writer.WriteString("</b></span>")

	switch b.CurrentGameState {
	case Draw: 
		writer.WriteString("<table class=\"draw\">")
	case ComputerIsOAndWon:
		writer.WriteString("<table class=\"computerisIandWon\">")
	case HumanisXandWon:
		writer.WriteString("<table class=\"humanisXandWon\">")	
	case GameNotFinished:
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

func (b *BoardInstance) PrintBoardtoConsoleForDebugging() {

    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            // Convert byte to string, using a dash if the byte is 0
            cell := "-"
            if b.Board[i][j] != 0 {
                cell = string(b.Board[i][j])
            }
            
            // Print the cell with vertical separators
            fmt.Printf(" %s ", cell)
            
            // Add vertical line between columns, except after the last column
            if j < 2 {
                fmt.Print("|")
            }
        }
        
        // Print a new line after each row
        fmt.Println()
        
        // Print horizontal separators between rows, except after the last row
        if i < 2 {
            fmt.Println("-----------")
        }
    }
}
