#include "rg.h"

#define N NORMAL
#define S SPAWN
#define O OBSTACLE

#define ARENA_RADIUS 8
#define GRID_SIZE 19
#define SPAWN_LEN 48

const int GRID[GRID_SIZE][GRID_SIZE] = {
	{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
	{O, O, O, O, O, O, O, S, S, S, S, S, O, O, O, O, O, O, O},
	{O, O, O, O, O, S, S, N, N, N, N, N, S, S, O, O, O, O, O},
	{O, O, O, S, S, N, N, N, N, N, N, N, N, N, S, S, O, O, O},
	{O, O, O, S, N, N, N, N, N, N, N, N, N, N, N, S, O, O, O},
	{O, O, S, N, N, N, N, N, N, N, N, N, N, N, N, N, S, O, O},
	{O, O, S, N, N, N, N, N, N, N, N, N, N, N, N, N, S, O, O},
	{O, S, N, N, N, N, N, N, N, N, N, N, N, N, N, N, N, S, O},
	{O, S, N, N, N, N, N, N, N, N, N, N, N, N, N, N, N, S, O},
	{O, S, N, N, N, N, N, N, N, N, N, N, N, N, N, N, N, S, O},
	{O, S, N, N, N, N, N, N, N, N, N, N, N, N, N, N, N, S, O},
	{O, S, N, N, N, N, N, N, N, N, N, N, N, N, N, N, N, S, O},
	{O, O, S, N, N, N, N, N, N, N, N, N, N, N, N, N, S, O, O},
	{O, O, S, N, N, N, N, N, N, N, N, N, N, N, N, N, S, O, O},
	{O, O, O, S, N, N, N, N, N, N, N, N, N, N, N, S, O, O, O},
	{O, O, O, S, S, N, N, N, N, N, N, N, N, N, S, S, O, O, O},
	{O, O, O, O, O, S, S, N, N, N, N, N, S, S, O, O, O, O, O},
	{O, O, O, O, O, O, O, S, S, S, S, S, O, O, O, O, O, O, O},
	{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O}
};

const Location SPAWN_LOCATIONS[SPAWN_LEN] = {
	(Location){.x = 1, .y = 7},
	(Location){.x = 1, .y = 8},
	(Location){.x = 1, .y = 9},
	(Location){.x = 1, .y = 10},
	(Location){.x = 1, .y = 11},
	(Location){.x = 2, .y = 5},
	(Location){.x = 2, .y = 6},
	(Location){.x = 2, .y = 12},
	(Location){.x = 2, .y = 13},
	(Location){.x = 3, .y = 3},
	(Location){.x = 3, .y = 4},
	(Location){.x = 3, .y = 14},
	(Location){.x = 3, .y = 15},
	(Location){.x = 4, .y = 3},
	(Location){.x = 4, .y = 15},
	(Location){.x = 5, .y = 2},
	(Location){.x = 5, .y = 16},
	(Location){.x = 6, .y = 2},
	(Location){.x = 6, .y = 16},
	(Location){.x = 7, .y = 1},
	(Location){.x = 7, .y = 17},
	(Location){.x = 8, .y = 1},
	(Location){.x = 8, .y = 17},
	(Location){.x = 9, .y = 1},
	(Location){.x = 9, .y = 17},
	(Location){.x = 10, .y = 1},
	(Location){.x = 10, .y = 17},
	(Location){.x = 11, .y = 1},
	(Location){.x = 11, .y = 17},
	(Location){.x = 12, .y = 2},
	(Location){.x = 12, .y = 16},
	(Location){.x = 13, .y = 2},
	(Location){.x = 13, .y = 16},
	(Location){.x = 14, .y = 3},
	(Location){.x = 14, .y = 15},
	(Location){.x = 15, .y = 3},
	(Location){.x = 15, .y = 4},
	(Location){.x = 15, .y = 14},
	(Location){.x = 15, .y = 15},
	(Location){.x = 16, .y = 5},
	(Location){.x = 16, .y = 6},
	(Location){.x = 16, .y = 12},
	(Location){.x = 16, .y = 13},
	(Location){.x = 17, .y = 7},
	(Location){.x = 17, .y = 8},
	(Location){.x = 17, .y = 9},
	(Location){.x = 17, .y = 10},
	(Location){.x = 17, .y = 11}
};
