
/*
- When serving a page using html/template you have to update the path of the .js and style sheets in the html files 

1. index.html is loaded 
2. the user is defaulted to being x 
3. gamePlayFacilitation.startGame creates the BoardInstance and initializes it 

*/
package main 

import(
	"fmt"
	"net/http"
)


func main() {
	fmt.Println("!!Welcome to the Triple T Tantilizer!!")

	/* this was used to set up the game board for the first time 
 	   and mainly to test the prediction tree 

	*/
	// var theMainGameBoard BoardInstance = BoardInstance  {
	// 	Board: [3][3]byte {
	// 		{'-','-','-'},
	// 		{'-','x','-'},
	// 		{'o','-','-'},
	// 	},
	// 	Html: "",
	// 	CellCount: 0,
	// 	Winner: false, 
	// 	BoardStateScore: 0,
	// 	MiniMaxScore: 0,
	// 	CurrentMove: 'o', //signifies latest move just happened
	// }

	// theMainGameBoard.CreateHTMLforTableToFindNextMove(0,0,false,false,0)
	

	var gamePlayFacilitator gamePlayFacilitation 

	// var nextMoveSymbol byte 
	// if gameBoard.CurrentMove == 'x' {
	// 	nextMoveSymbol = 'o'
	// } else {
	// 	nextMoveSymbol = 'x'
	// }

	// -- Tutorial on how to create a tree and render it using the Treant library
	//steps to set up tree, compute minimax and from that find next move amongst the root's children 
	// var moveTree *Tree = NewTree(gameBoard,nextMoveSymbol)
	// moveTree.BuildMoveTreeBreadthFirst(moveTree.Root)
	// moveTree.Minimax(moveTree.Root,true)
	// fmt.Println(moveTree.SearchMiniMaxedTreeForNextComputerMove())

	// -- UI Code 
	// treantTree := moveTree.CreateTreantJSONTree()
	// fmt.Println("Tree size:", moveTree.NodeCount)
	

// -- web server stuff START -- 
	fmt.Println("Starting Server")

	http.HandleFunc("/",index)
	// http.HandleFunc("/computeGameTreeAndReturnInTreantFormatForRenderinginBrowser/",gamePlayFacilitator.computeGameTreeAndReturnInTreantFormatForRenderinginBrowser)
	http.HandleFunc("/startgame/",gamePlayFacilitator.startGame)
	http.HandleFunc("/processHumanMove/",gamePlayFacilitator.processHumanMove)
	// http.HandleFunc("/executeComputerMove/",gamePlayFacilitator.executeComputerMove)
	http.HandleFunc("/getCurrentBoardState/",gamePlayFacilitator.getCurrentBoardState)

	

	http.Handle("/staticfiles/", http.StripPrefix("/staticfiles/", http.FileServer(http.Dir("staticfiles"))))

	
	http.ListenAndServe(":5000", nil)


}

































