package main 

import (
	"net/http"
	"encoding/json"
	"fmt"
)

// this struct enables this controller to be instantiated in main.go; you can create a prediction tree and pass in 
// a gameboard instance which represents the game's initial state 
type gamePlayFacilitation struct {
	predictionTree *Tree 
	gameBoard BoardInstance 
}



// -- web server handlers 
func (gp *gamePlayFacilitation) getCurrentBoardState (writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html")
	fmt.Println(gp.gameBoard.Html)
	fmt.Fprint(writer, gp.gameBoard.Html)
}
func (gp *gamePlayFacilitation) startGame (writer http.ResponseWriter, request *http.Request) {
	gp.gameBoard = BoardInstance {
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
		HumanPlayer: '-',
		ComputerPlayer: '-',
	} 

	var requestData map[string]interface{} 
	json.NewDecoder(request.Body).Decode(&requestData)

	//simplify and have the user always be X
	if requestData["userIsX"] == "x" {
		gp.gameBoard.HumanPlayer 	= 'x'
		gp.gameBoard.ComputerPlayer = 'o'
	} 
}
func (gp *gamePlayFacilitation)processHumanMove(writer http.ResponseWriter, request *http.Request){
	var requestData map[string]interface{} 
	json.NewDecoder(request.Body).Decode(&requestData) // get the cellIndex and the value within that cell 

	fmt.Println(requestData["cellIndex"],requestData["value"])

	var value string = requestData["value"].(string)
	
	//use string[0] to access a single byte; and .CurrentMove is of type byte 
	//tell gamboard what the current move is; so human is x, so currentmove is x. This is needed to calculate the tree
	//and in a previous iteration there was an idea to allow the user to be x or o, so the minimax tree
	//need(ed) (does) use this to recognize what the current move is, so it can compute the opposite player(symbol's) 
	//moves
	gp.gameBoard.CurrentMove = value[0] 

	// update the in memory gameboard 
	switch requestData["cellIndex"] {
	case "00":
		gp.gameBoard.Board[0][0] = value[0]
	case "01":
		gp.gameBoard.Board[0][1] = value[0]
	case "02": 
		gp.gameBoard.Board[0][2] = value[0]
	case "10":
		gp.gameBoard.Board[1][0] = value[0]
	case "11":
		gp.gameBoard.Board[1][1] = value[0]
	case "12":
		gp.gameBoard.Board[1][2] = value[0]
	case "20":
		gp.gameBoard.Board[2][0] = value[0]
	case "21":
		gp.gameBoard.Board[2][1] = value[0]
	case "22":
		gp.gameBoard.Board[2][2] = value[0]
	}
	if gp.gameBoard.IsGameFinished(value[0]) == true {
		fmt.Println("you won!")
	} else {
		// figure out computer's move by building prediction and computing minimax
		gp.gameBoard.CreateHTMLforTableToFindNextMove(0,0,false,false,0)
		gp.predictionTree = nil
		gp.predictionTree = NewTree(gp.gameBoard,gp.gameBoard.ComputerPlayer)
		gp.predictionTree.BuildMoveTreeBreadthFirst(gp.predictionTree.Root)
		gp.predictionTree.Minimax(gp.predictionTree.Root,true)


		//with minimax tree + scores created; search top layer of children to find highest score which is the next move
		nodeIndex, bi := gp.predictionTree.SearchMiniMaxedTreeForNextComputerMove()
		fmt.Println("next move should be nodeIndex:", nodeIndex)

		//overwrite the gameboard with the output from minimax search; thus indicating computer's move
		gp.gameBoard.Board = bi.Board 
		
		//render the board in html so it can be accessed via the getCurrentBoardState() endpoit
		gp.gameBoard.Html = ""
		gp.gameBoard.CreateHTMLforTableToFindNextMove(0,0,false,false,0)


		// take the built prediction tree + minimax scores on the nodes, and change into a format the treant libary can render in browser
		treantTree := gp.predictionTree.CreateTreantJSONTree()
		// send the tree, now in treant format for rendering, back to the browser 
		writer.Header().Set("Content-Type","application/json")
		json.NewEncoder(writer).Encode(treantTree)

	}
}


/*
	- HOW TO RENDER THE GAME TREE AND RENDER IN BROWSER USING TREANT LIBRARY 
	- deprecatred function 
	- it takes a pre-configurd BoardInstance gameboard and loads up the prediciton tree, produces the TreantTree 
	- and send it to the browser. Deprecated because now the user will start the game
*/
/*func (gp *gamePlayFacilitation)computeGameTreeAndReturnInTreantFormatForRenderinginBrowser(writer http.ResponseWriter, request *http.Request){

	fmt.Println("hello")

	gp.predictionTree = NewTree(gp.gameBoard,'x')
	gp.predictionTree.BuildMoveTreeBreadthFirst(gp.predictionTree.Root)
	gp.predictionTree.Minimax(gp.predictionTree.Root,true)

	nodeIndex, _ := gp.predictionTree.SearchMiniMaxedTreeForNextComputerMove()
	fmt.Println("next move should be nodeIndex:", nodeIndex)

	// take the built prediction tree + minimax scores on the nodes, and change into a format the treant libary can render in browser
	treantTree := gp.predictionTree.CreateTreantJSONTree()

	// send the tree, now in treant format for rendering, back to the browser 
	writer.Header().Set("Content-Type","application/json")
	json.NewEncoder(writer).Encode(treantTree)
}*/


