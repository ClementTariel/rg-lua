# Robot Game - Legendary Ultimate Arena

## Disclaimer

The project is currently Work In Progress, **some security features aren't implemented yet**.

**It is not safe to run untrusted scripts that you don't understand.**

## description

The rebirth of robotgame but in lua instead of python.

## Requirement 

You need [golang](https://go.dev/) 1.21.4 and [lua 5.3](https://www.lua.org/versions.html#5.3).

## Getting started

### Build

You need to build the first time you clone the repository and everytime you pull some changes.

To build you can simply run `make build`.

Or you can manually run the same commands individually:

```shell
cd referee
go install ../rgcore
go install ../player
go build -o ../rg.exe -a referee.go main.go
cd ..
```
### Run

Once the scripts are built, you can run matches locally.
To run locally a match between PATH/TO/LUA/BOT1 and PATH/TO/LUA/BOT2 use the following command: `./rg.exe PATH/TO/LUA/BOT1 PATH/TO/LUA/BOT2`.

exemple
```shell
./rg.exe bots/random.lua bots/random.lua
```


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
- `location` : The location of the robot as a Location with the properties `x` and `y`corresponding to its x and y coordinates. *Note: a Location can be directly compared to another Location, for example `self.location == rg.CENTER_POINT` will return `true` if the bot is at the center of the arena.*

### Accessing any robot's info

The `game` argument let you access robots information by Location  with its property `robots`: `game.robots[location]` returns the info of the robot placed at `location`. If the tile at `location.x`, `location.y` is empty, `nil` is returned. Each robot shares the same info as `self` except that `robot_id` is not accessible for enemy robots. You can loop over `game.robots` with `pairs`:
```lua
for loc, bot in pairs(game.robots) do
    -- Do something with loc and/or bot.
    -- Note that you do not need loc,
    -- because loc == bot.location is true
end
```

`game` also lets you access the current turn count with the property `turn`.

### Available interfaces

an object `rg` is accessible by all robots, with the following functions and properties :
- `Loc(x, y)` : A function that returns a Location `{ x=x, y=y }` that can be directly compared to another Location. For example `rg.Loc(9, 9) == rg.Loc(9, 9)` and `rg.Loc(9, 9) == { x=9, y=9 }` will both return `true` when `{ x=9, y=9 } == { x=9, y=9 }` will return `false`.
- `wdist(loc1, loc2)` : A function that returns the walking distance between Locations `loc1`  and `loc2`.
- `loc_type(loc)` : A function hat returns the type of a Location. a Location types are `OBSTACLE` for tiles outside of the arena, `SPAWN` for tiles where robots can spawn (border of the arena) and `NORMAL` for other tiles of the arena. The Location type of a tile stay the same during the game (a `NORMAL` tile will stay of type `NORMAL` even if occupied by a robot).
- `locs_around(loc)` : A function that returns the list of Locations around the Location `loc`. An optional list of location types to filter out can also be passed: `locs_around(loc, { SPAWN, OBSTACLE })` will remove from its output all tiles of type `SPAWN` or `OBSTACLE`.
- `toward(loc1, loc2)` : A function that returns the Location to move to from `loc1` in order to go to `loc2`.
- `CENTER_POINT = rg.Loc(9, 9)` : A Location corresponding to the center of the arena.
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

### Rules

Robots have 10ms to act, they have one action per turn: 
- `MOVE` : moves the robot to an adjacent tile.
- `ATTACK` : attack a tile at range.
- `GUARD` : stay in place, protect against collision damages and take half damage from attacks and suicides.
- `SUICIDE` : dies but deals damages to adjacent tiles.

There are no friendly damages (collison, attack and suicide only deal damage to enemy robots).

They are several waves in which robots spawn randomly, The team that have more robots at the end of the game wins.

## custom maps

You can create your own custom map in `grid.c`. The map must be a square with all its edges set to `OBSTACLE` to avoid unexpected behaviours. If you just want to play with the size of the arena you can use `generate_grid.py`, otherwise you can write your own script to generate any map you want. Make sure to keep consistency between`GRID` and `GRID_SIZE`, and between `SPAWN_LOCATIONS` and `SPAWN_LEN`.