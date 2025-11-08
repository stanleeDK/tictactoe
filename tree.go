 package main 

import (
		"fmt"
		"math"
		// "log"
		"strconv"
		"tictactoe/queue"
		"tictactoe/treantchart"
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

/* 

	-- to use this Tree model, you run the following functions (maybe you should wrap these in the constructor?)
	gp.predictionTree = NewTree(gp.gameBoard,gp.gameBoard.ComputerPlayer)
	gp.predictionTree.BuildMoveTreeBreadthFirst(gp.predictionTree.Root)
	gp.predictionTree.Minimax(gp.predictionTree.Root,true)
	//with minimax tree + scores created; search top layer of children to find highest score which is the next move
	nodeIndex, bi := gp.predictionTree.SearchMiniMaxedTreeForNextComputerMove()

*/ 

type Tree struct {
	Root 				*Node 
	NodeCount 			int 
	CurrentLevel 		int 
	ProgressChannelForMoveTreeBuilding chan int
}

func NewTree(pPayLoad BoardInstance) *Tree {
	var rootNode *Node = newNode(pPayLoad,0)
	progressChannel    := make(chan int,1000)
	return &Tree { 
		Root: 								rootNode, 
		NodeCount: 							1, 
		CurrentLevel: 						0, 
		ProgressChannelForMoveTreeBuilding: progressChannel,
	}
}

/*func (t *Tree) PrintNodesDepthFirst(traverseNode *Node) {}*/

func (t *Tree) PrintNodesBreadthFirst(/*root *Node*/) {
	
	var q = new(queue.Queue)
	q.Enqueue(t.Root)

	for q.IsEmpty() == false {
		traverseNode := q.PeekFront().(*Node)
		if len(traverseNode.Children) > 0 {
			for i:=0; i<len(traverseNode.Children); i++ {
				q.Enqueue(traverseNode.Children[i])
			}
		}
		// fmt.Println(traverseNode.Index)
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
func (t *Tree) BuildMoveTreeBreadthFirst(/*root *Node*/) {
	
	defer close(t.ProgressChannelForMoveTreeBuilding)
	var nodeCount 		int = 1 
	var nextMoveOorX 	byte 
	
	t.Root.PayLoad.MostRecentPlayer = 'x'

	var q 				= new(queue.Queue)
	q.Enqueue(t.Root) //add the root node to the queue 
	
	for q.IsEmpty() == false { // look through the entire queue 
		traverseNode := q.PeekFront().(*Node)

		// make 2 copies of the current node being inspected/popped from queue; 
		// two are needed because you give the current state of a board:  
		// 		iterate through all the empty slots and place an x (or o depending on who's turn it is) to represent a possible future move, see below how this is just a tool to create the child nodes which represent the actual game nodes. 
		// 		each iteration, return a board intance with the next move and add it to the tree 
		var tempBoardIterator BoardInstance  = traverseNode.PayLoad // a temporary vehcile simmply to drive the loop through the remaining slots and not used for persistence or core game logic. e.g. if there are 8 slots left becaue the human has moved, then the next level in the tree needs 8 children. This copy of the node/boardinstance is just used to loop through 8 times to ensure the creation of 8 children. It's passed by reference because the function that finds the next move can edit it at the same time as creating the child node. 
		var originalBoardState BoardInstance = traverseNode.PayLoad // this is a mechanism to be able to capture the original board state, which is needed because each child node is a separate move so you can create and pass in this copy to the findmove function, which will apply the move indpendently. If you didn't use this then each move made (and thus each child) would be applied to the original flooding the children with all of the incremental moves. 
		
		for isBoardFull(tempBoardIterator) == false {
			
			//alternate between cross and x depending on who's move it is because here you are building the game tree and each level alternates players
			if traverseNode.PayLoad.MostRecentPlayer == 'x' {
				nextMoveOorX = 'o'				
			} else {
				nextMoveOorX = 'x'	
			}

			
			var tempTreeLevel 						= traverseNode.PayLoad.CurrentTreeDepth + 1  // increment the child node's tree level based off what the parent's level is in the tree
			nextMoveBoardInstance 					:= t.findNextMove(nextMoveOorX,&tempBoardIterator,originalBoardState,nodeCount,tempTreeLevel)
			nextMoveBoardInstance.CurrentTreeDepth 	= tempTreeLevel 
			
			var childnode *Node  					= traverseNode.AddChild(nextMoveBoardInstance)		
			childnode.Index 						= nodeCount
			nodeCount++
			
			// put the count of nodes nto channel which will be streamed to the browser as progress updates
			select {
			case t.ProgressChannelForMoveTreeBuilding <- nodeCount: //send queue length into channel	
				// fmt.Println("Tree Size: ", nodeCount, len(t.ProgressChannelForMoveTreeBuilding))
				// time.Sleep(1 * time.Second)
			default:// Channel full or no listener, drop the message
        		// log.Println("Progress update dropped. Tree Size", nodeCount)
    			// fmt.Println("channel input blocked so skipped",len(t.ProgressChannelForMoveTreeBuilding))

        	}
			
			if nextMoveBoardInstance.CurrentGameState == GameNotFinished {
				q.Enqueue(childnode)  
			} 			
		}
		q.Dequeue()
		// use queue length to control the size of the future move tree you build 
		// if q.Length() > 10 { 
		// 	break 				
		// }
	}
	t.NodeCount = nodeCount
}


/*
	SYNPOSIS
	- Invoked in the context of building the game tree for the "AI". Specifically BuildMoveTreeBreadthFirst()
	- The AI uses this tree to map out all game states and find the next move, and this function creates a single game state boardinstance and is called 
	  repeatedly until all game outcome permutations (and the paths to them) are mapped out and scored. 
	- The goal is to find the next move given a boardinstance state, could be root but could be a node at any stage of the game. Create a child for this 
	  node so that represents the next move captured in the tree. Apply a score to that child node, as well as notate level and index number. 
	  This function is called until the entire tree is filled out i.e. all game permutations are stored in the tree data structure. 

	- To ensure you create the correct number of child nodes and in the right permutation ie. you don't want to create child nodes with duplicate game states; 
	  you have a vehicle 'currentTraverseNodeBoardInstance'. It helps control the iteration (see above explanation) by storing the next move temporarily.
	  Since the child node created by this function has the same gamestate as this currentTraverseNodeBoardInstance, it makes sense to pass it by reference
	  and simply update it within this function. 

	- nextMoveOorX - Even though the computer is always 'o', BuildMoveTreeBreadthFirst() can pass in 'o' or 'x' because it builds the entire move tree layer by layer, which each layer alternating simularing 'o' or 'x'
	- currentTraverseNodeBoardInstance
	- nextMove - this is the boardinstance that has the next move applied and will be returned, and the caller will add it to the tree as an official 
				 possible game move. It should be evaluated to see if one of the 'GameState' enums was achieved. This is in contract to currentTraverseNodeBoardInstance
				 which was just a vehicle for iteration, to ensure all the possible single moves can be iterated through and each one of those produces a nextMove. 


*/
func (t *Tree) findNextMove(nextMoveOorX byte, currentTraverseNodeBoardInstance *BoardInstance,originalBoardState BoardInstance, nodeIndex int, nodeTreeLevel int)(BoardInstance) {
	
	nextMove := originalBoardState

	for i:=0; i<3; i++ {
		for j:=0; j<3; j++ {
			if currentTraverseNodeBoardInstance.Board[i][j] == '-' {

				currentTraverseNodeBoardInstance.Board[i][j] = nextMoveOorX //make the move on tempboad by ref to trigger next iteration in the invoker function 
				nextMove.Board[i][j] 						 = nextMoveOorX 
				nextMove.MostRecentPlayer 					 = nextMoveOorX //update this boardinstance which will be a child node returned, with what player just made a move, x or o
				nextMove.ReturnEvaluationOfGameBoard() 
				
				switch nextMove.CurrentGameState {
				case Draw: 
					nextMove.BoardStateScore = 0 
				case ComputerIsOAndWon:
					nextMove.BoardStateScore = (10 - nodeTreeLevel) //computer wins at lower tree depths should be prioritied with higher scores 
				case HumanisXandWon:
					nextMove.BoardStateScore = (-10 + nodeTreeLevel) //shallow depth human wins should be prioritized. Because also here would be minimizing (selecting the minimum)
				case GameNotFinished:
					nextMove.BoardStateScore = 0 
				}
				nextMove.CreateHTMLforTableToFindNextMove(nodeIndex,nodeTreeLevel)
				return nextMove
			}
		}
	}

	return nextMove
}

func (t *Tree) Minimax(node *Node, maximizingPlayer bool) int {

	var x int
	if maximizingPlayer == true {
		x = Maximizer(node)
	} else {
		x = Minimizer(node)
	}
	return x
}

func Maximizer(node *Node) int {

	if len(node.Children) == 0 {
		node.PayLoad.MiniMaxScore = int(node.PayLoad.BoardStateScore)
		node.PayLoad.Html = "<div><b>mmS:"+  strconv.Itoa(node.PayLoad.MiniMaxScore) +"</b></div>" + node.PayLoad.Html
		return node.PayLoad.BoardStateScore 
	}
	
	maxEval  := -100.0 
	var tempEval int
	
	for i:=0;i<len(node.Children);i++ {
		tempEval = Minimizer(node.Children[i])
		maxEval  = math.Max(float64(maxEval),float64(tempEval))
	}
	node.PayLoad.MiniMaxScore = int(maxEval)
	node.PayLoad.Html = "<div><b>mmS:"+  strconv.Itoa(node.PayLoad.MiniMaxScore) +"</b></div>" + node.PayLoad.Html
	return int(maxEval)
}

func Minimizer(node *Node) int {

	if len(node.Children) == 0 {
		node.PayLoad.MiniMaxScore = int(node.PayLoad.BoardStateScore)
		node.PayLoad.Html = "<div><b>mmS:"+  strconv.Itoa(node.PayLoad.MiniMaxScore) +"</b></div>" + node.PayLoad.Html
		return node.PayLoad.BoardStateScore
	}
	minEval  := 100.0 
	var tempEval int

	for i:=0;i<len(node.Children);i++ {
		tempEval = Maximizer(node.Children[i])
		minEval  = math.Min(float64(minEval),float64(tempEval))
	}
	node.PayLoad.MiniMaxScore = int(minEval)
	node.PayLoad.Html = "<div><b>mmS:"+  strconv.Itoa(node.PayLoad.MiniMaxScore) +"</b></div>" + node.PayLoad.Html
	return int(minEval)
}

func (t *Tree) SearchMiniMaxedTreeForNextComputerMove() (int, BoardInstance) {
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
// ---- MINIMAX TREE CONSTRUCTION END -----





// ---- BROWSER JAVASCRIPT TREE RENDERING LOGIC START -----

/*
	1. to create a debug tree you have this div <div class="chart" id="OrganiseChart-simple">tree lives here</div>
	2. In the code that builds the move tree, processHumanMove() and BuildMoveTreeBreadthFirst() and CreateTreantJSONTree()
		the code below will generate the html of the tree and send it to the browser 

		treantTree := currentGameSession.PredictionTree.CreateTreantJSONTree()
		// send the tree, now in treant format for rendering, back to the browser 
		writer.Header().Set("Content-Type","application/json")
		json.NewEncoder(writer).Encode(treantTree)

	3. you then send the data by encoding it see last line of code above
	4. then simply receive it and initialize the tree 
		globalUITree = new Treant(data);
	5. TreantTree library finds your div from #1 and renders the tree within 


*/

// -- produces the json format needed by the UI library treant that renders the tree in the browser for debugging 
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
		PayLoadType: "thisisachart",
	}

	// t.Root.PayLoad.PrintBoardtoConsoleForDebugging()
	// fmt.Print(t.Root.PayLoad.Html)
	t.CreateTreantJSONPrintNodesDepthFirst(t.Root, &TChart.NodeStructure)	
	return TChart
}

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

// ---- BROWSER JAVASCRIPT TREE RENDERING LOGIC END -----





