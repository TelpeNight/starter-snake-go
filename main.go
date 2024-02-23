package main

// Welcome to
// __________         __    __  .__                               __
// \______   \_____ _/  |__/  |_|  |   ____   ______ ____ _____  |  | __ ____
//  |    |  _/\__  \\   __\   __\  | _/ __ \ /  ___//    \\__  \ |  |/ // __ \
//  |    |   \ / __ \|  |  |  | |  |_\  ___/ \___ \|   |  \/ __ \|    <\  ___/
//  |________/(______/__|  |__| |____/\_____>______>___|__(______/__|__\\_____>
//
// This file can be a nice home for your Battlesnake logic and helper functions.
//
// To get you started we've included code to prevent your Battlesnake from moving backwards.
// For more info see docs.battlesnake.com

import (
	"log"
	"maps"
	"math/rand"
)

// info is called when you create your Battlesnake on play.battlesnake.com
// and controls your Battlesnake's appearance
// TIP: If you open your Battlesnake URL in a browser you should see this data
func info() BattlesnakeInfoResponse {
	log.Println("INFO")

	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "Team7",   // TODO: Your Battlesnake username
		Color:      "#02d9e8", // TODO: Choose color
		Head:       "dragon",  // TODO: Choose head
		Tail:       "default", // TODO: Choose tail
	}
}

// start is called when your Battlesnake begins a game
func start(state GameState) {
	log.Println("GAME START")
}

// end is called when your Battlesnake finishes a game
func end(state GameState) {
	log.Printf("GAME OVER\n\n")
}

// move is called on every turn and returns your next move
// Valid moves are "up", "down", "left", or "right"
// See https://docs.battlesnake.com/api/example-move for available data
func move(state GameState) BattlesnakeMoveResponse {

	isMoveSafe := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	// We've included code to prevent your Battlesnake from moving backwards
	myHead := state.You.Body[0] // Coordinates of your head
	myNeck := state.You.Body[1] // Coordinates of your "neck"

	if myNeck.X < myHead.X { // Neck is left of head, don't move left
		isMoveSafe["left"] = false

	} else if myNeck.X > myHead.X { // Neck is right of head, don't move right
		isMoveSafe["right"] = false

	} else if myNeck.Y < myHead.Y { // Neck is below head, don't move down
		isMoveSafe["down"] = false

	} else if myNeck.Y > myHead.Y { // Neck is above head, don't move up
		isMoveSafe["up"] = false
	}

	// TODO: Step 1 - Prevent your Battlesnake from moving out of bounds
	boardWidth := state.Board.Width
	boardHeight := state.Board.Height
	if myHead.X == 0 {
		isMoveSafe["left"] = false
	}
	if myHead.X == boardWidth-1 {
		isMoveSafe["right"] = false
	}
	if myHead.Y == 0 {
		isMoveSafe["down"] = false
	}
	if myHead.Y == boardHeight-1 {
		isMoveSafe["up"] = false
	}

	// TODO: Step 2 - Prevent your Battlesnake from colliding with itself
	mybody := state.You.Body
	for _, coord := range mybody {
		if step, blocks := coordBlocksMove(myHead, coord); blocks {
			isMoveSafe[step] = false
		}
	}

	// TODO: Step 3 - Prevent your Battlesnake from colliding with other Battlesnakes
	otherHeads := map[Coord]int{}
	opponents := state.Board.Snakes
	for _, op := range opponents {
		for _, coord := range op.Body {
			if step, blocks := coordBlocksMove(myHead, coord); blocks {
				isMoveSafe[step] = false
			}
		}

		// Collect opponents heads location
		otherHeads[op.Head] = op.Length
	}

	// TODO: Prevent head collisions
	movesWithoutHeadCollisions := maps.Clone(isMoveSafe)
	for move, isSafe := range movesWithoutHeadCollisions {
		if isSafe {
			if isHeadsCollisionMove(move, myHead, state.You.Length, otherHeads) {
				movesWithoutHeadCollisions[move] = false
			}
		}
	}
	if haveSafeMoves(movesWithoutHeadCollisions) {
		isMoveSafe = movesWithoutHeadCollisions
	}

	if len(isMoveSafe) > 0 {
		movesToClosedSpace := ClosedSpaceMoves(state, isMoveSafe)
		if len(movesToClosedSpace) > 0 {
			saveMovesCount := 0
			for _, safe := range isMoveSafe {
				if safe {
					saveMovesCount++
				}
			}
			if len(movesToClosedSpace) < saveMovesCount {
				//avoid any unsafe move
				for _, move := range movesToClosedSpace {
					isMoveSafe[move.Move] = false
				}
			} else {
				//keep only move to the widest space
				for _, move := range movesToClosedSpace[1:] {
					isMoveSafe[move.Move] = false
				}
			}
		}
	}

	// Are there any safe moves left?
	safeMoves := []string{}
	for move, isSafe := range isMoveSafe {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		log.Printf("MOVE %d: No safe moves detected! Moving down\n", state.Turn)
		return BattlesnakeMoveResponse{Move: "down"}
	}

	// Choose a random move from the safe ones
	nextMove := safeMoves[rand.Intn(len(safeMoves))]

	// TODO: Step 4 - Move towards food instead of random, to regain health and survive longer
	// food := state.Board.Food

	log.Printf("MOVE %d: %s\n", state.Turn, nextMove)
	return BattlesnakeMoveResponse{Move: nextMove}
}

func haveSafeMoves(moves map[string]bool) bool {
	for _, safe := range moves {
		if safe {
			return true
		}
	}
	return false
}

func coordBlocksMove(head, coord Coord) (string, bool) {
	if head.X == coord.X {
		if head.Y+1 == coord.Y {
			return "up", true
		}
		if head.Y-1 == coord.Y {
			return "down", true
		}
	}
	if head.Y == coord.Y {
		if head.X+1 == coord.X {
			return "right", true
		}
		if head.X-1 == coord.X {
			return "left", true
		}
	}
	return "", false
}

func isHeadsCollisionMove(move string, myHead Coord, myLen int, heads map[Coord]int) bool {
	var target Coord

	switch move {
	case "up":
		target = Coord{X: myHead.X, Y: myHead.Y + 1}
	case "down":
		target = Coord{X: myHead.X, Y: myHead.Y - 1}
	case "left":
		target = Coord{X: myHead.X - 1, Y: myHead.Y}
	case "right":
		target = Coord{X: myHead.X + 1, Y: myHead.Y}
	}

	for enemy, length := range heads {
		if enemy.Y == target.Y && (enemy.X-target.X == 1) || (enemy.X-target.X == -1) {
			if length >= myLen {
				return true
			}
		}

		if enemy.X == target.X && (enemy.Y-target.Y == 1) || (enemy.Y-target.Y == -1) {
			if length >= myLen {
				return true
			}
		}
	}

	return false
}

func main() {
	RunServer()
}
