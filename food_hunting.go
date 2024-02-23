package main

func stepToClosestFood(head Coord, food Coord, nextMove string, moves []string) string {
	closest := distSqr(head, food)
	for _, step := range moves {
		next := getNextCoord(head, step)
		newDist := distSqr(next, food)
		if newDist < closest {
			closest = newDist
			nextMove = step
		}
	}
	return nextMove
}

func findClosedFood(head Coord, board Board) (Coord, bool) {
	var closest Coord
	var has bool
	var len float64
	for _, food := range board.Food {
		newLen := distSqr(food, head)
		if !has || newLen < len {
			closest = food
			has = true
			len = newLen
		}
	}
	return closest, has
}

func distSqr(a Coord, b Coord) float64 {
	x1, x2, y1, y2 := float64(a.X), float64(b.X), float64(a.Y), float64(b.Y)
	dx := x1 - x2
	dy := y1 - y2
	return dx*dx + dy*dy
}
