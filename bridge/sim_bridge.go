package bridge

/*
#cgo CFLAGS: -I${SRCDIR}
#include "sim.h"
*/
import "C"

type GameState struct {
	BallPosX int32
	BallPosY int32
	BallVelX int32
	BallVelY int32

	P1PosX  int32
	P1PosY  int32
	P1Score int32

	P2PosX  int32
	P2PosY  int32
	P2Score int32

	Tick uint64
}

type InputSet struct {
	P1Keys uint8
	P2Keys uint8
}

func toCState(g *GameState) C.GameState {
	var c C.GameState

	c.ball.position.x = C.int32_t(g.BallPosX)
	c.ball.position.y = C.int32_t(g.BallPosY)
	c.ball.velocity.x = C.int32_t(g.BallVelX)
	c.ball.velocity.y = C.int32_t(g.BallVelY)

	c.p1.position.x = C.int32_t(g.P1PosX)
	c.p1.position.y = C.int32_t(g.P1PosY)
	c.p1.score = C.int32_t(g.P1Score)

	c.p2.position.x = C.int32_t(g.P2PosX)
	c.p2.position.y = C.int32_t(g.P2PosY)
	c.p2.score = C.int32_t(g.P2Score)

	c.tick = C.uint64_t(g.Tick)

	return c
}

func fromCState(c C.GameState) GameState {
	return GameState{
		BallPosX: int32(c.ball.position.x),
		BallPosY: int32(c.ball.position.y),
		BallVelX: int32(c.ball.velocity.x),
		BallVelY: int32(c.ball.velocity.y),

		P1PosX:  int32(c.p1.position.x),
		P1PosY:  int32(c.p1.position.y),
		P1Score: int32(c.p1.score),

		P2PosX:  int32(c.p2.position.x),
		P2PosY:  int32(c.p2.position.y),
		P2Score: int32(c.p2.score),

		Tick: uint64(c.tick),
	}
}

func SimInit(state *GameState) {
	cState := toCState(state)
	C.sim_init(&cState)
	*state = fromCState(cState)
}

func SimTick(state *GameState, inputs InputSet) {
	cState := toCState(state)
	cInputs := C.InputSet{
		p1_keys: C.uint8_t(inputs.P1Keys),
		p2_keys: C.uint8_t(inputs.P2Keys),
	}
	C.sim_tick(&cState, cInputs)
	*state = fromCState(cState)
}

func SimHash(state *GameState) uint32 {
	cState := toCState(state)
	return uint32(C.sim_hash(&cState))
}

func SimSnapshot(src *GameState, dst *GameState) {
	cSrc := toCState(src)
	cDst := toCState(dst)
	C.sim_snapshot(&cSrc, &cDst)
	*dst = fromCState(cDst)
}

func SimRestore(dst *GameState, src *GameState) {
	cDst := toCState(dst)
	cSrc := toCState(src)
	C.sim_restore(&cDst, &cSrc)
	*dst = fromCState(cDst)
}