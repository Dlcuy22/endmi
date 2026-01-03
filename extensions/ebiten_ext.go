package extensions

func init() {
	RegisterTemplate(ebitenTemplate{})
}

type ebitenTemplate struct{}

func (ebitenTemplate) Name() string        { return "ebiten" }
func (ebitenTemplate) Description() string { return "Ebiten game engine template" }
func (ebitenTemplate) RootDir() string     { return "" }
func (ebitenTemplate) Dependencies() []string {
	return []string{"github.com/hajimehoshi/ebiten/v2"}
}
func (ebitenTemplate) Files(projectName string) map[string]string {
	return map[string]string{
		"main.go": `package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game represents the game state
type Game struct{}

// Update updates the game state
func (g *Game) Update() error {
	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(color.RGBA{15, 23, 42, 255})
	
	// Message
	ebitenutil.DebugPrintAt(screen, "It works!", 20, 20)
	ebitenutil.DebugPrintAt(screen, "Remove the boilerplate code to start your project", 20, 40)
	
	// Simple snake drawing
	snakeColor := color.RGBA{34, 197, 94, 255}
	// Snake body
	ebitenutil.DrawRect(screen, 20, 80, 20, 20, snakeColor)
	ebitenutil.DrawRect(screen, 40, 80, 20, 20, snakeColor)
	ebitenutil.DrawRect(screen, 60, 80, 20, 20, snakeColor)
	ebitenutil.DrawRect(screen, 80, 80, 20, 20, snakeColor)
	// Snake eyes
	ebitenutil.DrawRect(screen, 85, 85, 3, 3, color.White)
	ebitenutil.DrawRect(screen, 92, 85, 3, 3, color.White)
}

// Layout returns the game's screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("` + projectName + `")
	
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
`,
	}
}
