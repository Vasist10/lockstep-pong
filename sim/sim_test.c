#include "sim.h"
#include <stdio.h>
#include <stdlib.h>

int main(){
    //universe a
    GameState state_a;
    //universe b
    GameState state_b;
    sim_init(&state_a);
    sim_init(&state_b);

    //simulate 500 ticks with the same inputs and compare hashes
    for (int i = 0; i < 500; i++){
        InputSet inputs_a = {1,2}; //fixed bitmask values
        InputSet inputs_b = inputs_a; //same inputs for both universes

        sim_tick(&state_a, inputs_a);
        sim_tick(&state_b, inputs_b);

        uint32_t hash_a = sim_hash(&state_a);
        uint32_t hash_b = sim_hash(&state_b);

        if (hash_a != hash_b){
            printf("Desync at tick %llu: hash_a=%u, hash_b=%u\n", state_a.tick, hash_a, hash_b);
            return 1;
        }
        printf("tick %d: hash_a=%u hash_b=%u\n", i, hash_a, hash_b);
    }
    printf("ok-determinstic\n");
    return 0;
}