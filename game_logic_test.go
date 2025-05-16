package main

import (
	"testing"
)

func TestMoveSnakeUp(t *testing.T) {
	snake := Snake{
		Length:    3,
		Direction: Up,
		Segments: []Segment{
			{X: 10, Y: 10},
			{X: 10, Y: 11},
			{X: 10, Y: 12},
		},
	}

	MoveSnake(&snake)

	if snake.Segments[0].Y != 9 && snake.Segments[0].X != 10 {
		t.Errorf("Expected snake head to move up to (10, 9), got (%d, %d)", snake.Segments[0].X, snake.Segments[0].Y)
	}

	if snake.Segments[1].Y != 10 && snake.Segments[1].X != 10 {
		t.Errorf("Expected second segment to move to (10, 10), got (%d, %d)", snake.Segments[1].X, snake.Segments[1].Y)
	}

	if snake.Segments[2].Y != 11 && snake.Segments[2].X != 10 {
		t.Errorf("Expected third segment to move to (10, 11), got (%d, %d)", snake.Segments[2].X, snake.Segments[2].Y)
	}
}

func TestMoveSnakeRight(t *testing.T) {
	snake := Snake{
		Length:    3,
		Direction: Right,
		Segments: []Segment{
			{X: 10, Y: 10},
			{X: 10, Y: 11},
			{X: 10, Y: 12},
		},
	}

	MoveSnake(&snake)

	if snake.Segments[0].Y != 10 && snake.Segments[0].X != 11 {
		t.Errorf("Expected snake head to move up to (10, 11), got (%d, %d)", snake.Segments[0].X, snake.Segments[0].Y)
	}

	if snake.Segments[1].Y != 10 && snake.Segments[1].X != 10 {
		t.Errorf("Expected second segment to move to (10, 10), got (%d, %d)", snake.Segments[1].X, snake.Segments[1].Y)
	}

	if snake.Segments[2].Y != 10 && snake.Segments[2].X != 11 {
		t.Errorf("Expected third segment to move to (10, 11), got (%d, %d)", snake.Segments[2].X, snake.Segments[2].Y)
	}
}
