# Robot Game - Legendary Ultimate Arena

## Disclaimer

The project is currently Work In Progress, **some security features aren't implemented yet**.

**It is not safe to run untrusted scripts that you don't understand.**

## description

The rebirth of robotgame but in lua instead of python.

## Requirement 

You need [golang](https://go.dev/) 1.21.4 and [lua 5.3](https://www.lua.org/versions.html#5.3).

## Getting started

To run locally a match between PATH/TO/LUA/BOT1 and PATH/TO/LUA/BOT2 use the following command : `go run debug.go lua.go rg.go player.go referee.go main.go PATH/TO/LUA/BOT1 PATH/TO/LUA/BOT2`

*Note: The first time you run it after an update, run go with the `-a` flag to make sure C files are recompiled : `go run -a debug.go lua.go rg.go player.go referee.go main.go PATH/TO/LUA/BOT1 PATH/TO/LUA/BOT2`*

Checkout the [Documentation](#documentation) section to learn how to create your own robot.

Project currently Work In Progress.
Stay tuned.

## Documentation

### Create a robot

To create a bot you just need to create a lua script with a function `act(self, game)` that returns an object with :
- an `actionType` property with a value among `[MOVE, ATTACK, GUARD, SUICIDE]`
- `x` and `y` properties if the `actionType` property is set to `[MOVE, ATTACK]`

*Example: `{actionType=ATTACK, x=8, y=9}` is a valid return value.*

Here is an example of a very simple bot :
```lua
function act(self, game)
    return {actionType=GUARD}
end
```

### Accessing the robot's info

A robot's own info is accessible through the `self` argument. The default properties exposed are :
- `robot_id` : The unique id of the robot.
- `player_id` : An id shared by all the robots of a team.
- `hp` : The number of Health Points of the robot.
- `location` : The location of the robot as an object with the properties `x` and `y`correspondng to its x and y coordinates.

### Accessing any robot's info

**Note: The structure of `game` may change soon**

The `game` argument let you access robots information by coordinate : `game[x][y]` returns the info of the robot placed at `{ x=x, y=y }`. If the tile at `x`, `y` is empty, `nil` is returned. Each robot shares the same info as `self` except that `robot_id` is not accessible for enemy robots.

`game` also lets you access the current turn count with the property `turn`.

### Available interfaces

an object `rg` is accessible by all robots, with the following properties :
- `CENTER_POINT = { x=9, y=9 }` : The location of the center of the arena.
- `GRID_SIZE = 19` : The size of the grid.
- `ARENA_RADIUS = 8` : The radius of the area (The spawn tiles at the border of the arena have a distance to the center between `ARENA_RADIUS - 0.5` and `ARENA_RADIUS + 0.5` tiles).
- `SETTINGS` : an object exposing the main constants of the game :
    - `spawn_delay = 10` : The delay between 2 waves of spawns.
    - `spawn_count = 5` : The number of robots spawned for each team for each wave of spawns.
	- `robot_hp = 50` : The maximum number of Health Points of a robot.
	- `attack_range = 1` : The maximum distance to another bot for an attack to be valid.
	- `attack_damage` : The min and max damages dealt by an attack :
        - `min = 8` : the minimum
        - `max = 10` : the maximum
	- `suicide_damage = 15` : The damages dealt by a suicide.
	- `collision_damage = 5` : The damages dealt by a collision.
    - `max_turn = 100` : The number of turns in a game.
- `wdist(location1, location2)` : A function that returns the walking distance between location1  and location2.
- `locs_around(location)` : A function that returns the list of locations around location.

### Rules

Robots have 10ms to act, they have one action per turn : 
- `MOVE` : moves the robot to an adjacent tile.
- `ATTACK` : attack a tile at range.
- `GUARD` : stay in place, protect against collision damages and take half damage from attacks and suicides.
- `SUICIDE` : dies but deals damages to adjacent tiles.

There are no friendly damages (collison, attack and suicide only deal damage to enemy robots).

They are several waves in which robots spawn randomly, The team that have more robots at the end of the game wins.



