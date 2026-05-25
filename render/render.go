package render

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "lockstep-pong/bridge"
    "fmt"
	"image/color"
)

type PongGame struct {
	State *bridge.GameState
	OnUpdate func()
}
const screenWidth, screenHeight = 800, 600

func (pg *PongGame) Update() error {
    if pg.OnUpdate != nil {
        pg.OnUpdate()
    }
    return nil
}
func NewPongGame(state *bridge.GameState, onUpdate func()) *PongGame {
	return &PongGame{
		State: state,
		OnUpdate: onUpdate,
	}
}

func (pg *PongGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	paddle := ebiten.NewImage(20, 100)
	paddle.Fill(color.White)

	ball := ebiten.NewImage(20, 20)
	ball.Fill(color.White)

	p1Op := &ebiten.DrawImageOptions{}
	p1Op.GeoM.Translate(float64(pg.State.P1PosX), float64(pg.State.P1PosY))
	screen.DrawImage(paddle, p1Op)

	p2Op := &ebiten.DrawImageOptions{}
	p2Op.GeoM.Translate(float64(pg.State.P2PosX), float64(pg.State.P2PosY))
	screen.DrawImage(paddle, p2Op)

	ballOp := &ebiten.DrawImageOptions{}
	ballOp.GeoM.Translate(float64(pg.State.BallPosX), float64(pg.State.BallPosY))
	screen.DrawImage(ball, ballOp)

	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf(
			"P1: %d P2: %d",
			pg.State.P1Score,
			pg.State.P2Score,
		),
	)
}

func (pg *PongGame) Layout(outsideWidth, outsideHeight int) (int,int) {
	return 800, 600
}

