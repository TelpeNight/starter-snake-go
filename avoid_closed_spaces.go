package main

import (
	"cmp"
	"log"
	"slices"
)

type MoveToClosedSpace struct {
	Move   string
	Volume int
}

// ClosedSpaceMoves sorted by volume from high to low
func ClosedSpaceMoves(state GameState, testMoves map[string]bool) []MoveToClosedSpace {
	obstSet := make(map[Coord]struct{})
	for _, coord := range state.You.Body {
		obstSet[coord] = struct{}{}
	}
	for _, snake := range state.Board.Snakes {
		for _, coord := range snake.Body {
			obstSet[coord] = struct{}{}
		}
	}

	cellsCount := state.Board.Width * state.Board.Height
	freeCellsCount := cellsCount - len(obstSet)

	var result []MoveToClosedSpace
	for step, test := range testMoves {
		if !test {
			continue
		}
		nextCoord := getNextCooed(state.You.Head, step)
		reachableCellsCount := getReachableCellsCount(nextCoord, state.Board, obstSet)
		if reachableCellsCount == freeCellsCount {
			//not closed
			continue
		}
		result = append(result, MoveToClosedSpace{
			Move:   step,
			Volume: reachableCellsCount,
		})
	}
	slices.SortFunc(result, func(a, b MoveToClosedSpace) int {
		return -cmp.Compare(a.Volume, b.Volume)
	})
	return result
}

func getNextCooed(head Coord, step string) Coord {
	switch step {
	case "up":
		return Coord{
			X: head.X,
			Y: head.Y + 1,
		}
	case "down":
		return Coord{
			X: head.X,
			Y: head.Y - 1,
		}
	case "left":
		return Coord{
			X: head.X - 1,
			Y: head.Y,
		}
	case "right":
		return Coord{
			X: head.X + 1,
			Y: head.Y,
		}
	}
	log.Println("INVALID STEP", step)
	return head
}

func getReachableCellsCount(head Coord, board Board, obst map[Coord]struct{}) int {
	componentNum := 1
	visited := make(map[Coord]int)
	reachableNodes := bfs(componentNum, head, visited, board, obst)
	return len(reachableNodes)
}

func bfs(componentNum int, src Coord, visited map[Coord]int, board Board, obst map[Coord]struct{}) []Coord {
	queue := []Coord{src}
	visited[src] = componentNum
	var reachableNodes []Coord
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		reachableNodes = append(reachableNodes, u)
		adjs := getAdjCells(u, board, obst)
		for _, adj := range adjs {
			if visited[adj] == 0 {
				visited[adj] = componentNum
				queue = append(queue, adj)
			}
		}
	}
	return reachableNodes
}

func getAdjCells(u Coord, board Board, obst map[Coord]struct{}) []Coord {
	adj := []Coord{
		{u.X, u.Y + 1}, //up
		{u.X + 1, u.Y}, //right
		{u.X, u.Y - 1}, //down
		{u.X - 1, u.Y}, //left
	}
	adj = slices.DeleteFunc(adj, func(coord Coord) bool {
		if coord.X < 0 || coord.Y < 0 {
			return true
		}
		if coord.X >= board.Width || coord.Y >= board.Height {
			return true
		}
		if _, has := obst[coord]; has {
			return true
		}
		return false
	})
	return adj
}
