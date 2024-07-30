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
	(Location){.X = 1, .Y = 7},
	(Location){.X = 1, .Y = 8},
	(Location){.X = 1, .Y = 9},
	(Location){.X = 1, .Y = 10},
	(Location){.X = 1, .Y = 11},
	(Location){.X = 2, .Y = 5},
	(Location){.X = 2, .Y = 6},
	(Location){.X = 2, .Y = 12},
	(Location){.X = 2, .Y = 13},
	(Location){.X = 3, .Y = 3},
	(Location){.X = 3, .Y = 4},
	(Location){.X = 3, .Y = 14},
	(Location){.X = 3, .Y = 15},
	(Location){.X = 4, .Y = 3},
	(Location){.X = 4, .Y = 15},
	(Location){.X = 5, .Y = 2},
	(Location){.X = 5, .Y = 16},
	(Location){.X = 6, .Y = 2},
	(Location){.X = 6, .Y = 16},
	(Location){.X = 7, .Y = 1},
	(Location){.X = 7, .Y = 17},
	(Location){.X = 8, .Y = 1},
	(Location){.X = 8, .Y = 17},
	(Location){.X = 9, .Y = 1},
	(Location){.X = 9, .Y = 17},
	(Location){.X = 10, .Y = 1},
	(Location){.X = 10, .Y = 17},
	(Location){.X = 11, .Y = 1},
	(Location){.X = 11, .Y = 17},
	(Location){.X = 12, .Y = 2},
	(Location){.X = 12, .Y = 16},
	(Location){.X = 13, .Y = 2},
	(Location){.X = 13, .Y = 16},
	(Location){.X = 14, .Y = 3},
	(Location){.X = 14, .Y = 15},
	(Location){.X = 15, .Y = 3},
	(Location){.X = 15, .Y = 4},
	(Location){.X = 15, .Y = 14},
	(Location){.X = 15, .Y = 15},
	(Location){.X = 16, .Y = 5},
	(Location){.X = 16, .Y = 6},
	(Location){.X = 16, .Y = 12},
	(Location){.X = 16, .Y = 13},
	(Location){.X = 17, .Y = 7},
	(Location){.X = 17, .Y = 8},
	(Location){.X = 17, .Y = 9},
	(Location){.X = 17, .Y = 10},
	(Location){.X = 17, .Y = 11}
};
