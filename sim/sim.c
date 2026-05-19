#include "sim.h"
#include <string.h>

void sim_init(GameState* state){
    memset(state, 0, sizeof(GameState));
    state->ball.position.x = SCREEN_WIDTH/2 - BALL_SIZE/2;
    state->ball.position.y = SCREEN_HEIGHT/2 - BALL_SIZE/2;
    state->p1.position.x = PADDLE_WIDTH;
    state->p1.position.y = SCREEN_HEIGHT/2 - PADDLE_HEIGHT/2;
    state->p2.position.x = SCREEN_WIDTH - 2*PADDLE_WIDTH;
    state->p2.position.y = SCREEN_HEIGHT/2 - PADDLE_HEIGHT/2;
    state->ball.velocity.x = BALL_SPEED_INITIAL;
    state->ball.velocity.y = BALL_SPEED_INITIAL;
    state->p1.score = 0;
    state->p2.score = 0;
    state->tick = 0;
}

void sim_tick(GameState* state, InputSet inputs){
    
    //paddle movement
    //p1
    if (inputs.p1_keys & 1){
        state->p1.position.y -= PADDLE_SPEED;
        if (state->p1.position.y < 0){
            state->p1.position.y = 0;
        }
        
    }
    if (inputs.p1_keys & 2){
        state->p1.position.y += PADDLE_SPEED;
        if (state->p1.position.y > SCREEN_HEIGHT - PADDLE_HEIGHT){
            state->p1.position.y = SCREEN_HEIGHT - PADDLE_HEIGHT;
        }
        
    }

    //p2
    if (inputs.p2_keys & 1){
        state->p2.position.y -= PADDLE_SPEED;
        if (state->p2.position.y < 0){
            state->p2.position.y = 0;
        }
    }
    if (inputs.p2_keys & 2){
        state->p2.position.y += PADDLE_SPEED;
        if (state->p2.position.y > SCREEN_HEIGHT - PADDLE_HEIGHT){
            state->p2.position.y = SCREEN_HEIGHT - PADDLE_HEIGHT;
        }
    }

    //ball movement
    state->ball.position.x += state->ball.velocity.x;
    state->ball.position.y += state->ball.velocity.y;

    //ball collision with up
    if (state->ball.position.y < 0){
        state->ball.position.y = 0;
        state->ball.velocity.y = -state->ball.velocity.y;
    }
    //bottom
    if (state->ball.position.y > SCREEN_HEIGHT - BALL_SIZE){
        state->ball.position.y = SCREEN_HEIGHT - BALL_SIZE;
        state->ball.velocity.y = -state->ball.velocity.y;
    }
    
    //paddle collision
    int32_t ball_left = state->ball.position.x;
    int32_t ball_right = state->ball.position.x + BALL_SIZE;
    int32_t ball_top = state->ball.position.y;
    int32_t ball_bottom = state->ball.position.y + BALL_SIZE;

    //p1
    if (ball_left < state->p1.position.x + PADDLE_WIDTH &&
        ball_right > state->p1.position.x &&
        ball_top < state->p1.position.y + PADDLE_HEIGHT &&
        ball_bottom > state->p1.position.y) {
        state->ball.position.x = state->p1.position.x + PADDLE_WIDTH;
        state->ball.velocity.x = -state->ball.velocity.x;
    }

    //p2
    if (ball_left < state->p2.position.x + PADDLE_WIDTH &&
        ball_right > state->p2.position.x &&
        ball_top < state->p2.position.y + PADDLE_HEIGHT &&
        ball_bottom > state->p2.position.y) {
        state->ball.position.x = state->p2.position.x - BALL_SIZE;
        state->ball.velocity.x = -state->ball.velocity.x;
    }
    //left
    if (state->ball.position.x < 0){
        state->ball.position.x = SCREEN_WIDTH/2 - BALL_SIZE/2;
        state->ball.position.y = SCREEN_HEIGHT/2 - BALL_SIZE/2;
        state->ball.velocity.x = -BALL_SPEED_INITIAL;
        state->ball.velocity.y = BALL_SPEED_INITIAL;
        state->p2.score++;
    }
    //right
    if (state->ball.position.x > SCREEN_WIDTH - BALL_SIZE){
        state->ball.position.x = SCREEN_WIDTH/2 - BALL_SIZE/2;
        state->ball.position.y = SCREEN_HEIGHT/2 - BALL_SIZE/2;
        state->ball.velocity.x = BALL_SPEED_INITIAL;
        state->ball.velocity.y = BALL_SPEED_INITIAL;
        state->p1.score++;
    }

    state->tick++;
}

uint32_t sim_hash(const GameState* state){
    uint32_t hash = 2166136261u;
    const uint8_t* data = (const uint8_t*)state; //pointer casting
    for (size_t i = 0; i < sizeof(GameState); i++){
        hash ^= data[i];
        hash *= 16777619u;
    }

    return hash;
}

void sim_snapshot(const GameState* src, GameState* dst){
    memcpy(dst, src, sizeof(GameState));
}

void sim_restore(GameState* dst, const GameState* src){
    memcpy(dst, src, sizeof(GameState));
}