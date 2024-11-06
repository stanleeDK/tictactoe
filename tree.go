package main 

import (
		"fmt"
		"math"
		"strconv"
		// "encoding/json"
		"tictactoe/queue"
		"tictactoe/treantchart"
		// "io"
		// "math/rand"
    	// "time"
	)

// ---- NODE  DATA STRUCTURE START -----
type Node struct {
	PayLoad 	BoardInstance 
	Parent 		*Node //not used
	Traversed 	bool 
	Children 	[]*Node
	Level 		int 
	Index 		int 
}

//node constructor 
func newNode(pPayLoad BoardInstance, childrenCount int) *Node {
	var tempChildren 	= make([]*Node,childrenCount)
	var tempPayLoad 	= pPayLoad

	return &Node { PayLoad: tempPayLoad,Traversed: false, Children: tempChildren }
}

func (node *Node) AddChild(pPayLoad BoardInstance) *Node {
	newnode := newNode(pPayLoad,0)
	node.Children = append(node.Children, newnode)
	return newnode 		
}
// ---- NODE  DATA STRUCTURE END -----

// ----TREE DATA STRUCTURE START -----
type Tree struct {
	Root 				*Node 
	NodeCount 			int 
	CurrentLevel 		int 
	CurrentPlayerXorO 	byte 
}

func NewTree(pPayLoad BoardInstance, currentplayerxoro byte) *Tree {
	var rootNode *Node = newNode(pPayLoad,0)
	return &Tree { Root: rootNode, NodeCount: 1, CurrentLevel: 0, CurrentPlayerXorO: currentplayerxoro }
}

/*func (t *Tree) PrintNodesDepthFirst(traverseNode *Node) {}*/

func (t *Tree) PrintNodesBreadthFirst(/*root *Node*/) {
	
	var q = new(queue.Queue)
	q.Enqueue(/*root*/ t.Root)

	for q.IsEmpty() == false {
		traverseNode := q.PeekFront().(*Node)
		if len(traverseNode.Children) > 0 {
			for i:=0; i<len(traverseNode.Children); i++ {
				q.Enqueue(traverseNode.Children[i])
			}
		}
		fmt.Println(traverseNode.Index)
		traverseNode.Traversed = true 
		q.Dequeue()
	}
}
// ----TREE DATA STRUCTURE END -----


// ---- MINIMAX TREE CONSTRUCTION START -----
/*
	- Non-Recursive, queue based, breadth first minimax tree constuction 
	- This tree shows all the future moves given a current state of the tictactoe board 
	- Start by passing a Tree Root Node (see structs above), and the root node has a BoardInstance populated in payload then
	  it iwll construct all future tictactoe moves into the future. 
	- Control size of the tree (roughly) by controlling the size of the queue data structure 


*/
func (t *Tree) BuildMoveTreeBreadthFirst(root *Node) {
	
	var nodeCount int = 1 
	var nextMoveOorX  byte 

	var q = new(queue.Queue)
	q.Enqueue(root) //add the root node to the queue 
	
	for q.IsEmpty() == false { // look through the entire queue 
		traverseNode := q.PeekFront().(*Node)

		// make 2 copies of the current node being inspected/popped from queue; 
		// two are needed because you give the current state of a board:  
		// 		iterate through all the empty slots and place an x (or o depending on who's turn it is) to represent a possible future move
		// 		each iteration, return a board intance with the next move and add it to the tree 
		// you 
		var tempBoardIterator BoardInstance  = traverseNode.PayLoad // e.g. if there are 7 slots left, then iterate 7 times
		var parent BoardInstance 			 = traverseNode.PayLoad // e.g. for every iteration return a boardinstance with the move made, you return 7 of these in this example
		
		for isBoardFull(tempBoardIterator) == false {
			
			//alternate between cross and x depending on who's move it is
			if traverseNode.PayLoad.CurrentMove == 'x' {
				nextMoveOorX = 'o'				
			} else {
				nextMoveOorX = 'x'	
			}

			// increment the child node's tree level based off what the parent's level is in the tree
			var tempTreeLevel = traverseNode.PayLoad.CurrentTreeDepth + 1 

			nextMoveBoardInstance, isGameFinished := t.findNextMove(nextMoveOorX,&tempBoardIterator,parent,nodeCount,tempTreeLevel)		
			nextMoveBoardInstance.CurrentTreeDepth = tempTreeLevel 

			var childnode *Node  = traverseNode.AddChild(nextMoveBoardInstance)		
			
			childnode.Index = nodeCount
			nodeCount++

			
			if isGameFinished == false {
				q.Enqueue(childnode) 
			} 			
		}

		q.Dequeue()

		// use queue length to control the size of the future move tree you build 
		// if q.Length() > 300 { 
		// 	break 				
		// }
	}
	t.NodeCount = nodeCount
}

// Given current state of board, and a copy of it called "parent" which is actually a future parent after you predict moves
// function then returns a tree node that containing a future possible move in the game, 
// it's called over and over again until all the next possible moves, based off state of board, have been exhuasted  
// it also assigned the nodeindex (count) and the level/depth of the tree in which the returned future move sits 
func (t *Tree) findNextMove(nextMoveOorX byte, tempBoard *BoardInstance, parent BoardInstance, nodeIndex int, tempTreeLevel int)  (BoardInstance, bool) {
	
	var isGameFinished bool 
	var nextMove BoardInstance //create a node/possible future move 
	nextMove = parent //copy state of global gameboard 

	for i:=0; i<3; i++ {
		for j:=0; j<3; j++ {
			if tempBoard.Board[i][j] == '-' {

				tempBoard.Board[i][j] 	= nextMoveOorX //make the move on tempboad to trigger break to next level in tree
				nextMove.Board[i][j] 	= nextMoveOorX //make move on gamboard copy to create a node in the tree
				nextMove.CurrentMove 	= nextMoveOorX

				var oppositionWon bool = false 

				if nextMove.IsGameFinished(t.CurrentPlayerXorO) == false {
					
					if nextMove.IsGameFinishedForOpposition(t.CurrentPlayerXorO) == true {
						isGameFinished  = true 	
						oppositionWon = true 
						nextMove.BoardStateScore = -10
					} else {
						nextMove.BoardStateScore = -1
					}
					nextMove.CreateHTMLforTableToFindNextMove(nodeIndex,tempTreeLevel,isGameFinished,oppositionWon,int(nextMove.BoardStateScore))
					return nextMove, isGameFinished 
				} else {
					isGameFinished  = true 
					nextMove.Winner = isGameFinished
					nextMove.BoardStateScore = 10 //perfect score of 10 if game is won 
					nextMove.CreateHTMLforTableToFindNextMove(nodeIndex,tempTreeLevel,isGameFinished,oppositionWon,int(nextMove.BoardStateScore))
					return nextMove, isGameFinished 

				}
			}
		}
	}
	if nextMove.IsGameFinished(t.CurrentPlayerXorO) == false {
		isGameFinished = false 
		return nextMove, isGameFinished
	} else {
		isGameFinished = true 
		return nextMove, isGameFinished
	}
 }

func (t *Tree) Minimax(node *Node, maximizingPlayer bool) float64 {

	var x float64
	if maximizingPlayer == true {
		x = Maximizer(node)
	} else {
		x = Minimizer(node)
	}
	// fmt.Println("adf",x)
	return x
}

func Maximizer(node *Node) float64 {

	if len(node.Children) == 0 {
		node.PayLoad.MiniMaxScore = int(node.PayLoad.BoardStateScore)
		node.PayLoad.Html = "<div><b>mmS:"+  strconv.Itoa(node.PayLoad.MiniMaxScore) +"</b></div>" + node.PayLoad.Html
		return node.PayLoad.BoardStateScore 
	}
	
	maxEval  := -100.0 
	var tempEval float64
	
	for i:=0;i<len(node.Children);i++ {
		tempEval = Minimizer(node.Children[i])
		maxEval  = math.Max(maxEval,tempEval)
	}
	node.PayLoad.MiniMaxScore = int(maxEval)
	node.PayLoad.Html = "<div><b>mmS:"+  strconv.Itoa(node.PayLoad.MiniMaxScore) +"</b></div>" + node.PayLoad.Html
	return maxEval
}

func Minimizer(node *Node) float64 {

	if len(node.Children) == 0 {
		node.PayLoad.MiniMaxScore = int(node.PayLoad.BoardStateScore)
		node.PayLoad.Html = "<div><b>mmS:"+  strconv.Itoa(node.PayLoad.MiniMaxScore) +"</b></div>" + node.PayLoad.Html
		return node.PayLoad.BoardStateScore
	}
	minEval  := 100.0 
	var tempEval float64

	for i:=0;i<len(node.Children);i++ {
		tempEval = Maximizer(node.Children[i])
		minEval  = math.Min(minEval,tempEval)
	}
	node.PayLoad.MiniMaxScore = int(minEval)
	node.PayLoad.Html = "<div><b>mmS:"+  strconv.Itoa(node.PayLoad.MiniMaxScore) +"</b></div>" + node.PayLoad.Html
	return minEval
}

func (t *Tree) SearchMiniMaxedTreeForNextComputerMove() (int, BoardInstance) {

	fmt.Println("looking at root's children for next move")

	var minimaxScoreOfCurrentNode int = -100 
	var nodeIndexOfNextMove int = 0
	var boardinstanceOfNextMove BoardInstance  

	for i:=0; i<len(t.Root.Children); i++ {		
		fmt.Println("nodeid:", t.Root.Children[i].Index,"minimax score:", t.Root.Children[i].PayLoad.MiniMaxScore, "-, current prevailing minimaxscore", minimaxScoreOfCurrentNode, "nodeid of the next move: ")
		if t.Root.Children[i].PayLoad.MiniMaxScore > minimaxScoreOfCurrentNode {	
			nodeIndexOfNextMove       = t.Root.Children[i].Index
			minimaxScoreOfCurrentNode = t.Root.Children[i].PayLoad.MiniMaxScore
			boardinstanceOfNextMove   = t.Root.Children[i].PayLoad 
		}
	}
	return nodeIndexOfNextMove, boardinstanceOfNextMove

}

			// minEval := 100.0 
			// for i:=0;i<len(node.Children);i++ {
			// 	fmt.Println("   invoking maximizer on  node: ", node.Children[i].Index," ; minEVal currently is:", minEval/*," tempEval: ", tempEval)*/ )
			// 	tempEval 								:= t.Minimax(node.Children[i],true)
			// 	fmt.Println(node.Children[i].Index,":",minEval,tempEval)
			// 	minEval 								= math.Min(minEval,tempEval)
			// 	fmt.Println("min ouput:",minEval)
			// 	minimaxOutput 							= minEval
			// 	node.Children[i].PayLoad.MiniMaxScore 	= int(minimaxOutput)
				
			// 	node.Children[i].PayLoad.Html = "<div><b>mmS:"+  strconv.Itoa(node.Children[i].PayLoad.MiniMaxScore) +"</b></div>" + node.Children[i].PayLoad.Html
			// 	fmt.Println("   post invokation on node: ", node.Children[i].Index," its MMs is: ", node.Children[i].PayLoad.MiniMaxScore," child returned tempEval: ", tempEval)


// ---- MINIMAX TREE CONSTRUCTION END -----





// ---- BROWSER JAVASCRIPT TREE RENDERING LOGIC START -----
func (t *Tree) CreateTreantJSONPrintNodesDepthFirst(traverseNode *Node, treantNode *treantchart.TreantNodeStructure) {
	
	treantNode.InnerHTML = traverseNode.PayLoad.Html
	var nodeCount int = len(traverseNode.Children)
	if nodeCount > 0 {
		for i:=0; i<nodeCount; i++ {
			treantNode.Children = append(treantNode.Children,treantchart.TreantNodeStructure{})
			t.CreateTreantJSONPrintNodesDepthFirst(traverseNode.Children[i], &treantNode.Children[i])	
		}
	} else {
		return 
	}
}

func (t *Tree) CreateTreantJSONTree() treantchart.TreantChart {
	var TChart = treantchart.TreantChart {
		Chart: treantchart.TreantChartContainer {
			Container: "#OrganiseChart-simple",
			Scrollbar: "native",
		},
		NodeStructure: treantchart.TreantNodeStructure {
			InnerHTML: "",
			// Children: []TreantNodeStructure,		
		},
	}
	t.CreateTreantJSONPrintNodesDepthFirst(t.Root, &TChart.NodeStructure)
	
	// var writer strings.Builder 
	// encoder := json.NewEncoder(&writer)
	// encoder.Encode(&TChart)
	// return writer
	// fmt.Print(writer.String())
	
	return TChart
}
// ---- BROWSER JAVASCRIPT TREE RENDERING LOGIC END -----





