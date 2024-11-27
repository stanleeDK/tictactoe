

package main 

import(
	"fmt"
	"net/http"
	"os"
)

func main() {

// -- developpment -- 
	// tempGS := NewGameSession()
	// gss := NewGameSessions()
	// gss.CreateNewSession()
	// gss.CreateNewSession()
	// gss.CreateNewSession()
	// fmt.Println("hello", gss)

	// gpf := NewGamePlayFacilitation()


	// gamePlayFacilitator2 :=  NewGamePlayFacilitation()
	// http.HandleFunc("/startgame2/",gamePlayFacilitator2.startgame2)
	// http.HandleFunc("/processHumanMove2/",gamePlayFacilitator2.processHumanMove2)
// -- developpment -- 

// ----------------------------------


	fmt.Println("!!Welcome to the Triple T Tantilizer!!")

	gamePlayFacilitator := NewGamePlayFacilitation()

// -- web server stuff START -- 
	fmt.Println("Starting Server")

	http.HandleFunc("/",index)
	// http.HandleFunc("/computeGameTreeAndReturnInTreantFormatForRenderinginBrowser/",gamePlayFacilitator.computeGameTreeAndReturnInTreantFormatForRenderinginBrowser)
	http.HandleFunc("/startgame/",gamePlayFacilitator.startGame)
	http.HandleFunc("/restartGame/",gamePlayFacilitator.restartGame)
	http.HandleFunc("/processHumanMove/",gamePlayFacilitator.processHumanMove)
	// http.HandleFunc("/executeComputerMove/",gamePlayFacilitator.executeComputerMove)
	http.HandleFunc("/getCurrentBoardState/",gamePlayFacilitator.getCurrentBoardState)

	http.Handle("/staticfiles/", http.StripPrefix("/staticfiles/", http.FileServer(http.Dir("staticfiles"))))

	environment := os.Getenv("GO_ENV")	
	if (environment == "development") {
		http.ListenAndServe("localhost:5000", nil)
	} else {
		http.ListenAndServe("0.0.0.0:5000", nil)
	}


}


// -- Tutorial on how to create a tree and render it using the Treant library --
//steps to set up tree, compute minimax and from that find next move amongst the root's children 
// var moveTree *Tree = NewTree(gameBoard,nextMoveSymbol)
// moveTree.BuildMoveTreeBreadthFirst(moveTree.Root)
// moveTree.Minimax(moveTree.Root,true)
// fmt.Println(moveTree.SearchMiniMaxedTreeForNextComputerMove())

// -- UI Code 
// treantTree := moveTree.CreateTreantJSONTree()
// fmt.Println("Tree size:", moveTree.NodeCount)































