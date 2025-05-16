package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// enum for direction
const (
	Clockwise     = "clockwise"
	AntiClockwise = "anticlockwise"
	Up            = "up"
	Down          = "down"
	Left          = "left"
	Right         = "right"
)

var (
	directions = []string{Up, Right, Down, Left}
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	BorderWidth  = 10 // Thickness of the border
	HeaderHeight = 40
)

const (
	SegmentSize = 10
)

type Segment struct {
	X int
	Y int
}

type Snake struct {
	Length    int
	Direction string
	Segments  []Segment
}

type Game struct {
	snake          Snake
	frameCount     int
	speed          int
	gameOver       bool
	food           Segment
	score          int
	speedIncrement int
	obstacles      []Segment
}

var startingSnake = []Segment{
	{X: 5, Y: 5},
	{X: 4, Y: 5},
	{X: 3, Y: 5},
}

var randomSource *rand.Rand

func init() {
	randomSource = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (g *Game) SpawnObstacles(count int) {
	g.obstacles = []Segment{}
	for i := 0; i < count; i++ {
		for {
			x := rand.Intn((ScreenWidth - 2*BorderWidth) / SegmentSize)
			y := rand.Intn((ScreenHeight - HeaderHeight - 2*BorderWidth) / SegmentSize)

			newObstacle := Segment{X: x, Y: y}

			// Ensure no collisions with snake, food, or existing obstacles
			collision := false
			if newObstacle == g.food {
				collision = true
			}
			for _, segment := range g.snake.Segments {
				if segment == newObstacle {
					collision = true
					break
				}
			}
			for _, obstacle := range g.obstacles {
				if obstacle == newObstacle {
					collision = true
					break
				}
			}

			if !collision {
				g.obstacles = append(g.obstacles, newObstacle)
				break
			}
		}
	}
}

// Spawns food in a random location within the arena boundaries
func (g *Game) SpawnFood() {
	gridWidth := (ScreenWidth - 2*BorderWidth) / SegmentSize
	gridHeight := (ScreenHeight - 2*BorderWidth) / SegmentSize

	g.food = Segment{
		X: rand.Intn(gridWidth) + BorderWidth/SegmentSize,
		Y: rand.Intn(gridHeight) + BorderWidth/SegmentSize,
	}
}

func UpdateSnakeDirection(input string, snake *Snake) {

	index := 0
	for i, dir := range directions {
		if dir == snake.Direction {
			index = i
			break
		}
	}

	if input == Clockwise {
		index = (index + 1) % len(directions) // Move right in the array
	} else if input == AntiClockwise {
		index = (index - 1 + len(directions)) % len(directions) // Move left in the array
	}

	// Update the snake's direction
	snake.Direction = directions[index]
}

func MoveSnake(snake *Snake) {

	head := Segment{X: snake.Segments[0].X, Y: snake.Segments[0].Y}
	switch snake.Direction {
	case Up:
		head.Y--
	case Down:
		head.Y++
	case Left:
		head.X--
	case Right:
		head.X++
	}
	snake.Segments = append([]Segment{head}, snake.Segments[:snake.Length-1]...)
}

func (g *Game) Update() error {

	// Handle restart from game over
	if g.gameOver && inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.snake.Segments = startingSnake
		g.snake.Direction = "right"
		g.snake.Length = 3
		g.frameCount = 0
		g.score = 0
		g.gameOver = false
		g.speed = 5
		g.speedIncrement = 10
	}

	g.frameCount++

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		// Turn anticlockwise
		UpdateSnakeDirection(AntiClockwise, &g.snake)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		// Turn clockwise
		UpdateSnakeDirection(Clockwise, &g.snake)
	}

	if g.frameCount%g.speed == 0 {
		head := g.snake.Segments[0]
		// Collision detection with the border
		if head.X*SegmentSize < BorderWidth ||
			head.Y*SegmentSize < BorderWidth ||
			(head.X+1)*SegmentSize > ScreenWidth-BorderWidth ||
			(head.Y+1)*SegmentSize > (ScreenHeight-HeaderHeight-BorderWidth) {
			g.gameOver = true
		} else if g.food.X == head.X && g.food.Y == head.Y { // collision with food
			g.snake.Length++
			//g.snake.Segments = append([]Segment{head}, g.snake.Segments...)
			//g.snake.Segments = g.snake.Segments[:g.snake.Length]

			newHead := Segment{X: head.X, Y: head.Y}

			switch g.snake.Direction {
			case Up:
				newHead.Y--
			case Down:
				newHead.Y++
			case Left:
				newHead.X--
			case Right:
				newHead.X++
			}

			g.snake.Segments = append([]Segment{newHead}, g.snake.Segments...)

			// If the snake's length increased, we do not slice off the end
			if len(g.snake.Segments) > g.snake.Length {
				g.snake.Segments = g.snake.Segments[:g.snake.Length]
			}

			g.food.X = 0
			g.food.Y = 0
			g.score += 10
			g.SpawnFood()

			if g.score > 0 && g.score%g.speedIncrement == 0 {
				if g.speed > 2 {
					g.speed--
				}
			}
		} else {

			// Check for collision with its own body
			for _, segment := range g.snake.Segments[1:] {
				if head.X == segment.X && head.Y == segment.Y {
					g.gameOver = true
					return nil
				}
			}

			for _, obstacle := range g.obstacles {
				if g.snake.Segments[0] == obstacle {
					g.gameOver = true
					return nil
				}
			}

			MoveSnake(&g.snake)
		}
	}

	// Spawn food every 500 frames
	if g.frameCount%500 == 0 {
		g.SpawnFood()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	borderColor := color.RGBA{255, 0, 0, 255} // Red border for visibility

	// Draw the borders
	vector.DrawFilledRect(screen, 0, HeaderHeight, float32(ScreenWidth), float32(BorderWidth), borderColor, false)                                              // Top
	vector.DrawFilledRect(screen, 0, float32(ScreenHeight-BorderWidth), float32(ScreenWidth), float32(BorderWidth), borderColor, false)                         // Bottom
	vector.DrawFilledRect(screen, 0, HeaderHeight, float32(BorderWidth), float32(ScreenHeight-HeaderHeight), borderColor, false)                                // Left
	vector.DrawFilledRect(screen, float32(ScreenWidth-BorderWidth), HeaderHeight, float32(BorderWidth), float32(ScreenHeight-HeaderHeight), borderColor, false) // Right

	if g.gameOver {
		// Define the font and color
		face := basicfont.Face7x13
		msg := "Game Over! Press R to Restart"
		msgScore := fmt.Sprintf("Final Score: %d", g.score)
		textColor := color.RGBA{255, 255, 255, 255}
		// Measure the text width to center it
		bounds, _ := font.BoundString(face, msg)
		textWidth := (bounds.Max.X - bounds.Min.X).Ceil()
		x := (ScreenWidth - textWidth) / 2
		y := ScreenHeight / 2

		// Draw the text
		text.Draw(screen, msg, face, x, y, textColor)
		text.Draw(screen, msgScore, face, x, y+30, textColor)
		return
	}

	// Draw the food
	if g.food.X != 0 && g.food.Y != 0 {
		foodX := float32(g.food.X * SegmentSize)
		foodY := float32(g.food.Y*SegmentSize) + float32(HeaderHeight) // Add the header height here
		vector.DrawFilledRect(screen, foodX, foodY, float32(SegmentSize), float32(SegmentSize), color.RGBA{255, 165, 0, 255}, false)
	}

	// Draw the score
	scoreText := fmt.Sprintf("Score: %d", g.score)
	text.Draw(screen, scoreText, basicfont.Face7x13, 10, 20, color.RGBA{255, 255, 255, 255})

	// Draw the speed
	speedText := fmt.Sprintf("Speed: %d", g.speed)
	text.Draw(screen, speedText, basicfont.Face7x13, 200, 20, color.RGBA{255, 255, 255, 255})

	// Draw the obstacles
	for _, obstacle := range g.obstacles {
		vector.DrawFilledRect(screen,
			float32(obstacle.X*SegmentSize),
			float32(obstacle.Y*SegmentSize+HeaderHeight),
			float32(SegmentSize),
			float32(SegmentSize),
			color.RGBA{108, 122, 137, 255}, false) // Red colored obstacles
	}

	// Draw the snake
	for _, segment := range g.snake.Segments {
		x := float32(segment.X * SegmentSize)
		y := float32(segment.Y*SegmentSize) + float32(HeaderHeight) // Offset by the header height

		if x >= float32(BorderWidth) &&
			y >= float32(BorderWidth+HeaderHeight) &&
			x+float32(SegmentSize) <= float32(ScreenWidth-BorderWidth) &&
			y+float32(SegmentSize) <= float32(ScreenHeight-BorderWidth) {
			vector.DrawFilledRect(screen, x, y, float32(SegmentSize), float32(SegmentSize), color.RGBA{0, 255, 0, 255}, false)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	snake := Snake{
		Length:    3,
		Direction: "right",
		Segments:  startingSnake,
	}

	game := &Game{snake: snake, speed: 5, speedIncrement: 10, score: 0}
	game.SpawnFood()
	game.SpawnObstacles(5) // Spawn 5 obstacles
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Snake")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
