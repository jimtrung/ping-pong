package main

import (
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"

	"fmt"

	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const (
    screenWidth = 640
    screenHeight = 480
    ballSpeed = 3
    paddleSpeed = 6
)

type Object struct {
    X, Y, W, H int
}

type Paddle struct {
    Object
}

type Ball struct {
    Object
    dxdt  int
    dydt  int
}

type Game struct {
    paddle    Paddle
    ball      Ball
    score     int
    highScore int
    extra     int
}

func main() {
    ebiten.SetWindowTitle("Ping Pong")
    ebiten.SetWindowSize(screenWidth, screenHeight)
    paddle := Paddle {
        Object: Object{
            X: 600,
            Y: 200,
            W: 15,
            H: 75,
        },
    }
    ball := Ball {
        Object: Object{
            X: 0,
            Y: 0,
            W: 15,
            H: 15,
        },
        dxdt: ballSpeed,
        dydt: ballSpeed,
    }
    g := &Game {
        paddle: paddle,
        ball: ball,
    }

    err := ebiten.RunGame(g)
    if err != nil {
        log.Fatal(err)
    }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
    vector.DrawFilledRect(screen,
        float32(g.paddle.X), float32(g.paddle.Y),
        float32(g.paddle.W), float32(g.paddle.H),
        color.White, false,
    )

    vector.DrawFilledRect(screen,
        float32(g.ball.X), float32(g.ball.Y),
        float32(g.ball.W), float32(g.ball.H),
        color.White, false,
    )

    scoreString := fmt.Sprintf("Score: %v", g.score)
    text.Draw(screen, scoreString, basicfont.Face7x13, 10, 10, color.White)

    highscoreString := fmt.Sprintf("HighScore: %v", g.highScore)
    text.Draw(screen, highscoreString, basicfont.Face7x13, 10, 30, color.White)
}

func (g *Game) Update() error {
    g.paddle.MoveOnKeyPress()
    g.ball.Move()
    g.CollideWithWall()
    g.CollideWithPaddle()
    g.Speed()
    g.GetHighScore()
    return nil
}

func (p *Paddle) MoveOnKeyPress() {
    if ebiten.IsKeyPressed(ebiten.KeyS) {
        p.Y += paddleSpeed
    }
    if ebiten.IsKeyPressed(ebiten.KeyW) {
        p.Y -= paddleSpeed
    }
    if ebiten.IsKeyPressed(ebiten.KeyEscape) {
        os.Exit(0)
    }
}

func (b *Ball) Move() {
    b.X += b.dxdt
    b.Y += b.dydt
}

func (g *Game) Reset() {
    g.ball.X = 0
    g.ball.Y = 0
    g.extra = 0

    g.score = 0
}

func (g *Game) CollideWithWall() {
    if g.ball.X >= screenWidth {
        g.Reset()
    } else if g.ball.X <= 0 {
        g.ball.dxdt = ballSpeed + g.extra
    } else if g.ball.Y >= screenHeight {
        g.ball.dydt = -ballSpeed -g.extra
    } else if g.ball.Y <= 0 {
        g.ball.dydt = ballSpeed + g.extra
    }
}

func (g *Game) CollideWithPaddle() {
    if g.ball.X >= g.paddle.X && g.ball.Y >= g.paddle.Y && g.ball.Y <= g.paddle.Y + g.paddle.H {
        g.ball.dxdt = -g.ball.dxdt
        g.score++
        if g.score > g.highScore {
            g.highScore = g.score
        }
    }
}

func (g *Game) Speed() {
    g.extra = g.score / 10
}
