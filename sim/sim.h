#if !defined(SIM_H)
#define SIM_H
#include <stdint.h>
#define SCREEN_WIDTH 800
#define SCREEN_HEIGHT 600
#define PADDLE_HEIGHT 100
#define PADDLE_WIDTH 20
#define BALL_SIZE 20
#define TICK_RATE 60
#define PADDLE_SPEED 5
#define BALL_SPEED_INITIAL 5

typedef struct {
    uint8_t p1_keys;
    uint8_t p2_keys;
} InputSet;

typedef struct {
    int32_t x;
    int32_t y;
} Vec2;
typedef struct {
    Vec2 position;
    Vec2 velocity;
} Ball;
typedef struct {
    Vec2 position;
    int32_t score;
} Paddle;

typedef struct {
    Ball ball;
    Paddle p1;
    Paddle p2;
    uint64_t tick;
    
} GameState;

void sim_init(GameState* state);
void sim_tick(GameState* state, InputSet inputs);
uint32_t sim_hash(const GameState* state);
void sim_snapshot(const GameState* src, GameState* dst);
void sim_restore(GameState* dst, const GameState* src);

#endif