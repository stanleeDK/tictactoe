

package main 

import(
	"fmt"
	"net/http"
	"os"
	"time"
	// "encoding/json"
)

func main() {



	// board := NewBoardInstance() 
	// board.Board[0][0] = 'x'
	// board.Board[0][1] = 'o'
	// board.Board[0][2] = 'o'

	// board.Board[1][0] = '-'
	// board.Board[1][1] = 'o'
	// board.Board[1][2] = 'x'
	
	// board.Board[2][0] = '-'
	// board.Board[2][1] = 'x'
	// board.Board[2][2] = 'x'
	// tree := NewTree(*board)
	// fmt.Println(board.Board)

	// tree.BuildMoveTreeBreadthFirst(/*tree.Root*/)
	// tree.Minimax(tree.Root, true)
	// // board.PrintBoardtoConsoleForDebugging()
	// board.ReturnEvaluationOfGameBoard()
	// fmt.Println(board.CurrentGameState)

	// // http.HandleFunc("/",index)
	// http.HandleFunc("/showdebuggametree", func(writer http.ResponseWriter, request *http.Request) {
	// 	fmt.Println("hello")
	// 	treantTree := tree.CreateTreantJSONTree()
	// 	writer.Header().Set("Content-Type","application/json")
	// 	json.NewEncoder(writer).Encode(treantTree)
	// })



// -- developpment -- 

// ----------------------------------


	fmt.Println("!!Welcome to the Triple T Tantilizer!!")
	fmt.Println("Starting Server")

	gamePlayFacilitator := NewGamePlayFacilitation()

// -- web server stuff START -- 
	http.HandleFunc("/",index)
	http.HandleFunc("/startgame/",gamePlayFacilitator.startGame)
	http.HandleFunc("/restartGame/",gamePlayFacilitator.restartGame)
	http.HandleFunc("/processHumanMove/",gamePlayFacilitator.processHumanMove)
	http.HandleFunc("/getCurrentBoardState/",gamePlayFacilitator.getCurrentBoardState)
	http.HandleFunc("/showgametree/", gamePlayFacilitator.showGameTree)
	http.HandleFunc("/getprogress/", gamePlayFacilitator.getMoveTreeBuildingProgress)
	http.Handle("/staticfiles/", http.StripPrefix("/staticfiles/", http.FileServer(http.Dir("staticfiles"))))
	http.Handle("/tmnt/", http.StripPrefix("/tmnt/", http.FileServer(http.Dir("tmnthangman"))))



	environment := os.Getenv("GO_ENV")	
	fmt.Println (environment)
	if (environment == "development") {
		http.ListenAndServe("localhost:8888", nil)
	} else {
		// err := http.ListenAndServe("0.0.0.0:80", nil)

		srv := &http.Server{
		    Addr:              "0.0.0.0:80",
		    Handler:           nil,  // nil = use http.DefaultServeMux (your existing handlers)
		    // ReadTimeout:       5 * time.Minute,  // For large uploads
		    WriteTimeout:      60 * time.Second,
		    // ReadHeaderTimeout: 10 * time.Second,
		    IdleTimeout:       120 * time.Second,
		}
		err := srv.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}



}




























