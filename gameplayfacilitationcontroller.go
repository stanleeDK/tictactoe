package main 

import (
	"net/http"
	"encoding/json"
	"fmt"
)

// this struct enables this controller to be instantiated in main.go; 
// it contains a GameSessions type, which is a map of multiple gameSession
type gamePlayFacilitation struct {
	primaryGameSessions *GameSessions

	gameBoard *BoardInstance 
	predictionTree *Tree 
}

// constructor which calls the constructor of GameSessions (which creates the game map)
// returns instance of this gamePlayFacilitator to main.go
func NewGamePlayFacilitation() *gamePlayFacilitation {
	gs := NewGameSessions()
	return &gamePlayFacilitation{
		primaryGameSessions: gs,
	}
}

// -- web server handlers 
func (gp *gamePlayFacilitation) getCurrentBoardState (writer http.ResponseWriter, request *http.Request) {
	params := request.URL.Query()
	writer.Header().Set("Content-Type", "text/html")
	fmt.Fprint(writer, gp.primaryGameSessions.GameSessionsRunning[params.Get("sessionid")].GameBoard.Html)
}

// handler for creating game session and returning sessionid to the client 
func (gp *gamePlayFacilitation) startGame (writer http.ResponseWriter, request *http.Request) {
	var response struct {
		SessionID string `json:"sessionid"`
	}
	response.SessionID = gp.newGame()
	writer.Header().Set("Content-Type","application/json")
	json.NewEncoder(writer).Encode(response)	
}

// func (gp *gamePlayFacilitation) startGame (writer http.ResponseWriter, request *http.Request) {
// 	gp.newGame()
// }

// create a new game by creating a new game session, which adds a game session to the gamesessions map
func (gp *gamePlayFacilitation)newGame() string {
	sessionID := gp.primaryGameSessions.CreateNewSession()
	return sessionID 
}

// func (gp *gamePlayFacilitation)newGame() string {
// 	gp.gameBoard = NewBoardInstance()
// 	gp.predictionTree = nil
// 	return ""
// }

func (gp *gamePlayFacilitation)processHumanMove(writer http.ResponseWriter, request *http.Request){

	gp.printAllGames()
/* use requestData as a map to extract the below 
            var humanMove = {
                cellIndex: cellIndex,
                value: value,
                sessionid: sessionStorage.getItem('sessionid') 
            }
*/

	var requestData map[string]interface{} 
	json.NewDecoder(request.Body).Decode(&requestData) // get the cellIndex and the value within that cell 

/*
	use string[0](value[0]) to access a single byte i.e the first byte which we know will contain 'x'; 
	and .CurrentMove is of type byte (your code could have been much simpler if you just used string...)
	tell gamboard what the current move is; so human is x, so currentmove is x. This is needed to calculate the tree
	In a previous iteration there was an idea to allow the user to be x or o, so the minimax tree
	need(ed) (does) use this to recognize what the current move is, so it can compute the opposite player(symbol's) moves

*/	
	var value string 						 = requestData["value"].(string) //cast "value" into a string
	var seshID 								 = requestData["sessionid"].(string)
	var currentGameSession *GameSession 	 = gp.primaryGameSessions.GameSessionsRunning[seshID]
	currentGameSession.GameBoard.CurrentMove = value[0]
	
	// update the in memory gameboard 
	switch requestData["cellIndex"] {
	case "00":
		currentGameSession.GameBoard.Board[0][0] = value[0]
	case "01":
		currentGameSession.GameBoard.Board[0][1] = value[0]
	case "02": 
		currentGameSession.GameBoard.Board[0][2] = value[0]
	case "10":
		currentGameSession.GameBoard.Board[1][0] = value[0]
	case "11":
		currentGameSession.GameBoard.Board[1][1] = value[0]
	case "12":
		currentGameSession.GameBoard.Board[1][2] = value[0]
	case "20":
		currentGameSession.GameBoard.Board[2][0] = value[0]
	case "21":
		currentGameSession.GameBoard.Board[2][1] = value[0]
	case "22":
		currentGameSession.GameBoard.Board[2][2] = value[0]
	}

	if currentGameSession.GameBoard.IsGameFinished(value[0]) == true {
		fmt.Println("you won!")
		currentGameSession.GameBoard.CreateHTMLforTableToFindNextMove(0,0,true,false,0)
	} else {
		// figure out computer's move by building prediction tree and computing minimax
		currentGameSession.PredictionTree = nil
		currentGameSession.PredictionTree = NewTree(*currentGameSession.GameBoard,currentGameSession.GameBoard.ComputerPlayer)
		currentGameSession.PredictionTree.BuildMoveTreeBreadthFirst(currentGameSession.PredictionTree.Root)
		currentGameSession.PredictionTree.Minimax(currentGameSession.PredictionTree.Root,true)

		//with minimax tree + scores created; search top layer of children to find highest score which is the next move
		nodeIndex, bi := currentGameSession.PredictionTree.SearchMiniMaxedTreeForNextComputerMove()
		fmt.Println("next move should be nodeIndex:", nodeIndex)

		//overwrite the gameboard with the output from minimax search; thus indicating computer's move
		currentGameSession.GameBoard.Board = bi.Board 

		currentGameSession.GameBoard.Html = ""

		// function takes a byte, but keeping "o" is more readable because computer is always "o"
		// so this ugly code casts string into byte slice and since it's one character return first entry in slice which is "o" in byte format
		if currentGameSession.GameBoard.IsGameFinished([]byte("o")[0]) == true {
			//computer won as "o" i.e. you lost
			fmt.Println("you lost", seshID)
			//render the board in html so it can be accessed via the getCurrentBoardState() endpoit
			currentGameSession.GameBoard.CreateHTMLforTableToFindNextMove(0,0,true,true,0)
		} else {
			//render the board in html so it can be accessed via the getCurrentBoardState() endpoit
			currentGameSession.GameBoard.CreateHTMLforTableToFindNextMove(0,0,false,false,0)
		}

		// take the built prediction tree + minimax scores on the nodes, and change into a format the treant libary can render in browser
		treantTree := currentGameSession.PredictionTree.CreateTreantJSONTree()
		// send the tree, now in treant format for rendering, back to the browser 
		writer.Header().Set("Content-Type","application/json")
		json.NewEncoder(writer).Encode(treantTree)
	}
}

// func (gp *gamePlayFacilitation)processHumanMove(writer http.ResponseWriter, request *http.Request){
// 	var requestData map[string]interface{} 
// 	json.NewDecoder(request.Body).Decode(&requestData) // get the cellIndex and the value within that cell 

// 	fmt.Println(requestData["cellIndex"],requestData["value"],requestData["sessionid"])

// 	var value string = requestData["value"].(string)
	
// 	//use string[0] to access a single byte; and .CurrentMove is of type byte 
// 	//tell gamboard what the current move is; so human is x, so currentmove is x. This is needed to calculate the tree
// 	//and in a previous iteration there was an idea to allow the user to be x or o, so the minimax tree
// 	//need(ed) (does) use this to recognize what the current move is, so it can compute the opposite player(symbol's) 
// 	//moves
// 	gp.gameBoard.CurrentMove = value[0] 

// 	// update the in memory gameboard 
// 	switch requestData["cellIndex"] {
// 	case "00":
// 		gp.gameBoard.Board[0][0] = value[0]
// 	case "01":
// 		gp.gameBoard.Board[0][1] = value[0]
// 	case "02": 
// 		gp.gameBoard.Board[0][2] = value[0]
// 	case "10":
// 		gp.gameBoard.Board[1][0] = value[0]
// 	case "11":
// 		gp.gameBoard.Board[1][1] = value[0]
// 	case "12":
// 		gp.gameBoard.Board[1][2] = value[0]
// 	case "20":
// 		gp.gameBoard.Board[2][0] = value[0]
// 	case "21":
// 		gp.gameBoard.Board[2][1] = value[0]
// 	case "22":
// 		gp.gameBoard.Board[2][2] = value[0]
// 	}

// 	if gp.gameBoard.IsGameFinished(value[0]) == true {
// 		fmt.Println("you won!")
// 		gp.gameBoard.CreateHTMLforTableToFindNextMove(0,0,true,false,0)
// 	} else {
// 		// figure out computer's move by building prediction tree and computing minimax
// 		gp.predictionTree = nil
// 		gp.predictionTree = NewTree(*gp.gameBoard,gp.gameBoard.ComputerPlayer)
// 		gp.predictionTree.BuildMoveTreeBreadthFirst(gp.predictionTree.Root)
// 		gp.predictionTree.Minimax(gp.predictionTree.Root,true)

// 		//with minimax tree + scores created; search top layer of children to find highest score which is the next move
// 		nodeIndex, bi := gp.predictionTree.SearchMiniMaxedTreeForNextComputerMove()
// 		fmt.Println("next move should be nodeIndex:", nodeIndex)

// 		//overwrite the gameboard with the output from minimax search; thus indicating computer's move
// 		gp.gameBoard.Board = bi.Board 

// 		gp.gameBoard.Html = ""

// 		// function takes a byte, but keeping "o" is more readable because computer is always "o"
// 		// so this ugly code casts string into byte slice and since it's one character return first entry in slice which is "o" in byte format
// 		if gp.gameBoard.IsGameFinished([]byte("o")[0]) == true {
// 			//computer won as "o" i.e. you lost
// 			fmt.Println("you lost")
// 			//render the board in html so it can be accessed via the getCurrentBoardState() endpoit
// 			gp.gameBoard.CreateHTMLforTableToFindNextMove(0,0,true,true,0)
// 		} else {
// 			//render the board in html so it can be accessed via the getCurrentBoardState() endpoit
// 			gp.gameBoard.CreateHTMLforTableToFindNextMove(0,0,false,false,0)
// 		}

// 		// take the built prediction tree + minimax scores on the nodes, and change into a format the treant libary can render in browser
// 		treantTree := gp.predictionTree.CreateTreantJSONTree()
// 		// send the tree, now in treant format for rendering, back to the browser 
// 		writer.Header().Set("Content-Type","application/json")
// 		json.NewEncoder(writer).Encode(treantTree)
// 	}
// }
func (gp *gamePlayFacilitation)restartGame(writer http.ResponseWriter, request *http.Request){
	params := request.URL.Query()
	gp.primaryGameSessions.GameSessionsRunning[params.Get("sessionid")].GameBoard = NewBoardInstance()
	gp.primaryGameSessions.GameSessionsRunning[params.Get("sessionid")].PredictionTree = nil
	fmt.Println("helloggghgpoop[y] butt [line from henrik nov 2024...]")
}

func (gp *gamePlayFacilitation)printAllGames(){
	fmt.Println("all running game sessions")
	for key, value := range gp.primaryGameSessions.GameSessionsRunning{
		fmt.Println(key,value)
	}
}



