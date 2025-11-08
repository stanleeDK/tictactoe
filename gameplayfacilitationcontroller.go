package main 

import (
	"net/http"
	"encoding/json"
	"fmt"
	"sync"
	// "time"
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
// this actually sends the computer move via a boardinstance struct that has html generated for the 3*3 array and sends the html
// to the browser reprsenting a move. the html is appended to a hidden div, which the javascript retrieves 
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


// create a new game by creating a new game session, which adds a game session to the gamesessions map
func (gp *gamePlayFacilitation)newGame() string {
	sessionID := gp.primaryGameSessions.CreateNewSession()
	return sessionID 
}



func (gp *gamePlayFacilitation)processHumanMove(writer http.ResponseWriter, request *http.Request){

    var wg sync.WaitGroup

	// gp.printAllGames() // debug statement to show all concurrent players' games to terminal 

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
	currentGameSession, exists  			 := gp.primaryGameSessions.GetSession(seshID)
	if !exists {//session not found, quit and return message 
		fmt.Fprintf(writer,"session not found") 
		return 
	}

	currentGameSession.GameBoard.MostRecentPlayer = value[0]


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

	currentGameSession.GameBoard.ReturnEvaluationOfGameBoard()
	
	switch currentGameSession.GameBoard.CurrentGameState {
	case Draw: 
		// fmt.Println("It's a draw!") 
		currentGameSession.GameBoard.CreateHTMLforTableToRenderUserFacingBoardState()
	case ComputerIsOAndWon:
		// fmt.Println("Computer Won!") 
		currentGameSession.GameBoard.CreateHTMLforTableToRenderUserFacingBoardState()
	case HumanisXandWon:
		// fmt.Println("Human Won!") 
		currentGameSession.GameBoard.CreateHTMLforTableToRenderUserFacingBoardState()
	case GameNotFinished:
		//then do computer move 
		currentGameSession.PredictionTree = nil
		currentGameSession.PredictionTree = NewTree(*currentGameSession.GameBoard)
		wg.Add(1)
		go func (){ //put this tree building in a goroutine so you can feed updates to a channel; wait for it to complete
			defer wg.Done()
			currentGameSession.PredictionTree.BuildMoveTreeBreadthFirst()
		 }()
		wg.Wait() 
		currentGameSession.PredictionTree.Minimax(currentGameSession.PredictionTree.Root,true)
		
		nodeIndex, bi := currentGameSession.PredictionTree.SearchMiniMaxedTreeForNextComputerMove() //with minimax tree + scores created; search top layer of children to find highest score which is the next move
		fmt.Println("next move should be nodeIndex:", nodeIndex)
		
		currentGameSession.GameBoard.Board = bi.Board  //overwrite the gameboard with the output from minimax search; thus indicating computer's move
		// currentGameSession.GameBoard.PrintBoardtoConsoleForDebugging()
		currentGameSession.GameBoard.ReturnEvaluationOfGameBoard()
		currentGameSession.GameBoard.CreateHTMLforTableToRenderUserFacingBoardState()//render the board in html so it can be accessed via the getCurrentBoardState() endpoit

	}
}


	// -- STREAMING PROGRESS CODE START ---
func (gp *gamePlayFacilitation)getMoveTreeBuildingProgress(writer http.ResponseWriter, request *http.Request){

    writer.Header().Set("Content-Type", "text/event-stream")
    writer.Header().Set("Cache-Control", "no-cache")
    writer.Header().Set("Connection", "keep-alive")
    writer.Header().Set("Access-Control-Allow-Origin", "*")
    flusher, _ 		:= writer.(http.Flusher)
    

	// GET parameters 
	params := request.URL.Query()
	seshID := params.Get("sessionid")


	progressInfo := map[string]interface{} {
		"PayLoadType": "thisisprogressdata",
		"data": "", 	
	}


	currentGameSession, exists  := gp.primaryGameSessions.GetSession(seshID)
	if !exists {
		progressInfo["data"] = "session not found"
		data, _ := json.Marshal(progressInfo)
		fmt.Fprintf(writer,"data: %s\n\n",data)
	}

	// get the channel of the current session 
	progressChannel := currentGameSession.PredictionTree.ProgressChannelForMoveTreeBuilding
	
	for {
		select{
		case progressInfo, ok := <-progressChannel:
		
			if !ok {
				fmt.Fprintf(writer, "data: {\"payloadtype\":\"complete\",\"progress\":\"complete\"}\n\n")
				// fmt.Println("channel closed")
				flusher.Flush()
				return 
			}
			// fmt.Println("tree size from channel: ", progressInfo)
			fmt.Fprintf(writer, "data: {\"payloadtype\":\"%d\",\"progress\":\"pending\"}\n\n", progressInfo)
			flusher.Flush()
			
		
		case <-request.Context().Done():
			// Client disconnected
			fmt.Printf("Client disconnected from session: %s", seshID)
			return
		}
	}
}

func (gp *gamePlayFacilitation)restartGame(writer http.ResponseWriter, request *http.Request){
	params := request.URL.Query()
	gp.primaryGameSessions.GameSessionsRunning[params.Get("sessionid")].GameBoard = NewBoardInstance()
	gp.primaryGameSessions.GameSessionsRunning[params.Get("sessionid")].PredictionTree = nil
	fmt.Println("helloggghgpoop[y] butt [line from henrik nov 2024...]")
}

func (gp *gamePlayFacilitation)showGameTree(writer http.ResponseWriter, request *http.Request){
	params 								:= request.URL.Query()
	var currentGameSession *GameSession = gp.primaryGameSessions.GameSessionsRunning[params.Get("sessionid")]

	// take the built prediction tree + minimax scores on the nodes, and change into a format the treant libary can render in browser
	treantTree := currentGameSession.PredictionTree.CreateTreantJSONTree()

	// send the tree, now in treant format for rendering, back to the browser 
	writer.Header().Set("Content-Type","application/json")
	json.NewEncoder(writer).Encode(treantTree)

}

func (gp *gamePlayFacilitation)printAllGames(){
	fmt.Println("all running game sessions")
	for key, value := range gp.primaryGameSessions.GameSessionsRunning{
		fmt.Println(key,value)
	}
}




