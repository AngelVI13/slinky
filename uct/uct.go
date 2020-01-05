package uct

import (
	"fmt"
	"math/rand"
	"slinky/board"
	"sort"
	"time"
)

type rankedMove struct {
	move  int
	score float64
	node  *Node
}

type uctResult struct {
	node             *Node
	move             int
	wins             float64
	visits           float64
	totalSimulations int
}

func uct(rootstate *board.ChessBoard, node *Node, originMove int, timeData timeInfo) uctResult {
	rootstate.MakeMove(originMove)
	/* Check for immediate result
	It is possible the game is already over by this point
	in which the value of the move should be immediately computed and
	put in the result from the view point of the enemy
	since here moves are evaluated from that viewpoint */
	isResult, result := isImmediateResult(rootstate, originMove)
	if isResult == true {
		return result
	}

	var rootnode *Node
	if node != nil {
		if node.move <= 0 {
			rootnode = CreateRootNode(rootstate)
		} else {
			fmt.Printf("Reusing node tree: move %s visits %f wins %f\n", board.PrintMove(node.move), node.visits, node.wins)
			rootnode = node
		}
	} else {
		rootnode = CreateRootNode(rootstate)
	}

	state := rootstate

	simulations := 0

	elapsedTime := time.Since(timeData.startTime).Seconds() * 1000 // get elapsed time in ms
	for timeData.isTimeSet == true && elapsedTime < float64(timeData.stopTime) {
		simulations++ // count number of simulations done
		node := rootnode
		movesToRoot := 0

		// Select stage
		// node is fully expanded and non-terminal
		for len(node.untriedMoves) == 0 && len(node.childNodes) > 0 {
			node = node.SelectChild()
			state.MakeMove(node.move)
			movesToRoot++
		}

		// Expand
		// if we can expand (i.e. state/node is non-terminal)
		if len(node.untriedMoves) > 0 {
			move := node.untriedMoves[rand.Intn(len(node.untriedMoves))]
			state.MakeMove(move)
			movesToRoot++
			// add child and descend tree
			node = node.AddChild(move, state)
		}

		// Rollout
		//  - this can often be made orders of magnitude quicker
		//    using a state.GetRandomMove() function

		// while state is non-terminal
		for state.GetResult(state.GetPlayerJustMoved()) == board.NoWinner {
			moves := state.GetMoves()
			m := moves[rand.Intn(len(moves))]
			state.MakeMove(m)
			movesToRoot++
		}

		// Backpropagate
		// backpropagate from the expanded node and work back to the root node
		for node != nil {
			// state is terminal.
			// Update node with result from POV of node.playerJustMoved
			gameResult := state.GetResult(node.playerJustMoved)
			node.Update(float64(gameResult))
			node = node.parent
		}

		// Revert all the made moves
		for j := 0; j < movesToRoot; j++ {
			state.TakeMove()
		}

		elapsedTime = time.Since(timeData.startTime).Seconds() * 1000 // get elapsed time in ms
	}

	// Todo try returning move with highest average score
	sort.Slice(rootnode.childNodes, func(i, j int) bool {
		return rootnode.childNodes[i].visits > rootnode.childNodes[j].visits
	})
	// above we sort by descending order -> move with most visits is the first element
	bestMove := rootnode.childNodes[0]
	// todo not needed to send the rest of the items back since they are all inside bestMove
	return uctResult{node: bestMove, move: originMove, wins: bestMove.wins, visits: bestMove.visits, totalSimulations: simulations}
}

type timeInfo struct {
	startTime time.Time
	stopTime  int
	isTimeSet bool
}

type uctArg struct {
	state    *board.ChessBoard
	node     *Node
	move     int
	timeData timeInfo
}

func worker(jobs <-chan uctArg, results chan<- uctResult) {
	for uctArguments := range jobs {
		results <- uct(uctArguments.state, uctArguments.node, uctArguments.move, uctArguments.timeData)
	}
}

// GetEngineMoveFast returns the best move found by the UCT (computed in parallel)
// todo return Struct instead of a lot of individual vars
func GetEngineMoveFast(state *board.ChessBoard, info *board.SearchInfo, node *Node) (move int, score float64, totalSim int, moveNode *Node) {
	// todo this could use the values from Node instead of asking for moves again
	availableMoves := state.GetMoves()
	numMoves := len(availableMoves)

	if numMoves == 0 {
		panic("Game is already over, can't get engine move for a finished game!")
	} else if numMoves == 1 {
		// no clue what the score is here since we haven't actually searched the move
		return availableMoves[0], 0.5, 0, &Node{}
	}

	// todo return whole node from mcts and then send it down on the jobs for the next run
	// todo send a cut-off of depth -> even if game is not over - evaluate position and return evaluation as if it's a proper result
	// todo make sure to take care of quescent position i.e. do not stop immediatelly but on a quiet position
	// todo add move ordering i.e. more promissing moves might get more iterations ??
	// create channels to share data between goroutines
	jobs := make(chan uctArg, numMoves)
	results := make(chan uctResult, numMoves)

	// spawn workers ready to process data
	for _i := 0; _i < numMoves; _i++ {
		go worker(jobs, results)
	}

	timeData := timeInfo{startTime: info.StartTime, stopTime: info.StopTime, isTimeSet: info.TimeSet}
	bestMove := rankedMove{move: -1, score: 1.1}
	// if we have a valid node tree provided -> enter last user move and make that node root
	if node != nil {
		if node.move != 0 {
			for _, child := range node.childNodes {
				if child.move == state.GetLastMove() {
					node = child
					break
				}
			}
		}
	}

	for _, move := range availableMoves {
		// create a copy of the board in order to be sent to the goroutine
		b := *state

		var childNode *Node
		if node != nil {
			for _, child := range node.childNodes {
				if child.move == move {
					childNode = child
					childNode.parent = nil
					break
				}
			}
		}

		jobs <- uctArg{
			state:    &b,
			node:     childNode,
			move:     move,
			timeData: timeData,
		}
	}

	close(jobs) // close jobs channel

	var result uctResult
	var scoreValue float64
	var totalSimulations int
	for _i := 0; _i < numMoves; _i++ {
		result = <-results
		scoreValue = result.wins / result.visits

		fmt.Printf("Move: %s: %.3f -> %.1f / %.0f (%d)\n",
			board.PrintMove(result.move), scoreValue, result.wins, result.visits, result.totalSimulations)
		// here the move_score refers to the best enemy reply
		// therefore we want to minimize that i.e. chose the move
		// which leads to the lowest scored best enemy reply
		if scoreValue < bestMove.score {
			bestMove.score = scoreValue
			bestMove.move = result.move
			bestMove.node = result.node
		}
		totalSimulations += result.totalSimulations
	}
	fmt.Printf("Total simulations done: %d\n", totalSimulations)

	return bestMove.move, bestMove.score, totalSimulations, bestMove.node
}

func isImmediateResult(state *board.ChessBoard, move int) (isResult bool, result uctResult) {
	enemy := state.GetEnemy(state.GetPlayerJustMoved())
	gameResult := state.GetResult(enemy)

	if gameResult != board.NoWinner {
		result = uctResult{node: &Node{}, move: move, wins: float64(gameResult), visits: 1.0}
		isResult = true
	}

	return
}

// GetEngineMove returns the best move found by the UCT
// func GetEngineMove(state board.ChessBoard, simulations int) int {
// 	availableMoves := state.GetMoves()

// 	if len(availableMoves) == 0 {
// 		panic("Game is already over, can't get engine move for a finished game!")
// 	} else if len(availableMoves) == 1 {
// 		return availableMoves[0]
// 	}

// 	simPerMove := simulations / len(availableMoves)

// 	bestMove := rankedMove{move: -1, score: 1.1}
// 	for _, move := range availableMoves {
// 		b := state // todo does this copy ? or points

// 		mScore := uct(b, move, simPerMove)
// 		scoreValue := mScore.wins / mScore.visits

// 		fmt.Printf("Move: %d: %.3f -> %f / %f\n", move, scoreValue, mScore.wins, mScore.visits)
// 		// here the move_score refers to the best enemy reply
// 		// therefore we want to minimize that i.e. chose the move
// 		// which leads to the lowest scored best enemy reply
// 		if scoreValue < bestMove.score {
// 			bestMove.score = scoreValue
// 			bestMove.move = move
// 		}
// 		// take move since b here is not a copy but points to the same position as state
// 		b.TakeMove()
// 	}
// 	return bestMove.move
// }
