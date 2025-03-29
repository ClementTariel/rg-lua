# Robot Game - Legendary Ultimate Arena

## Disclaimer

The project is currently Work In Progress, **some security features aren't implemented yet**.

**It is not safe to run untrusted scripts that you don't understand.**

**This project is not ready for production. Do not deploy it on a machine exposed to external connections**

## description

The rebirth of robotgame but in lua instead of python.

## Requirement 

You need:
- [docker](https://www.docker.com/) 
- an API client such as [bruno](https://www.usebruno.com/) (a bruno collection is provided in this project under the flolder `bruno`)

## Quick start

### Build

You need to build the first time you clone the repository and everytime you pull some changes.

To build you can simply run `make build`.

To clear the database if you have messed it up and want to reset it run `make clear-db`

### Run

Once the images are built, you can run matches locally.

First you need to run the images, open a terminal and run `make run`.
The output of a match will appear in this terminal when it ends.
You can stop the images at anytime by running `make stop` in another tab.

Then open your API Client.

- With an API Client other than Bruno

Make a `POST` request to `localhost:4444/request-match`, the body should have the following structure:
```json
{
  "blueName": "{{blueName}}",
  "redName": "{{redName}}"
}
```
With `{{blueName}}` and `{{redName}}` being the names of the bots.

- With Bruno

Open the rglua collection (under the folder `bruno` in this project) and run the `request-match` request. You can adjust the name of the bots to use in the `Vars` tab. 

The names must correspond to names of bots in the database.  

Checkout the [Default robots](#default-robots) section for more information on the robots available.

### Troubleshooting

Sometimes a `Network deployment_default  Error` can occur. I did not find any other solution than a complete reboot of the host to fix this issue.

you can use the adminer client on a browser at `http://localhost:8888/?pgsql=rglua_db&db=rglua&ns=public` to inspect the database with username `rglua_user` an password `rglua_temporary_password`.

### Default robots

The first time you launch the database (or the first time after a `make clear-db` is run), it is populated with bots under the `bots`. For each file, a bot entity is created, with the field `script` set to the content of the file and the field `name` set to the name of the file **without its extension**. For example to run the script under `random.lua`, the name to use is `random`.

If you want to add a robot to the defaults robots you can put your robot script in a lua file under `bots` and run `make clear-db` to make sure that it will be loaded the next time your launch the database.

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

Some constants are defined for the possible actions: `MOVE`, `ATTACK`, `GUARD`, `SUICIDE`.

Some constants are defined for the possible locations types: `NORMAL`, `SPAWN`, `OBSTACLE`.

an object `rg` is accessible by all robots, with the following functions and properties :
- `Loc(x, y)` : A function that returns a Location `{ x=x, y=y }` that can be directly compared to another Location. For example `rg.Loc(9, 9) == rg.Loc(9, 9)` and `rg.Loc(9, 9) == { x=9, y=9 }` will both return `true` when `{ x=9, y=9 } == { x=9, y=9 }` will return `false`.
- `wdist(loc1, loc2)` : A function that returns the walking distance (manhattan distance) between Locations `loc1`  and `loc2`.
- `dist(loc1, loc2)` : A function that returns the euclidean distance between Locations `loc1`  and `loc2`.
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
	- `attack_damage` : The min and max damage dealt by an attack :
        - `min = 8` : the minimum
        - `max = 10` : the maximum
	- `suicide_damage = 15` : The damage dealt by a suicide.
	- `collision_damage = 5` : The damage dealt by a collision.
    - `max_turn = 100` : The number of turns in a game.

### Rules

Robots have 10ms to act, they have one action per turn: 
- `MOVE` : moves the robot to an adjacent tile.
- `ATTACK` : attack a tile at range.
- `GUARD` : stay in place, protect against collision damage and take half damage from attacks and suicides.
- `SUICIDE` : dies but deals damage to adjacent tiles.

There are no friendly damage (collison, attack and suicide only deal damage to enemy robots).

They are several waves in which robots spawn randomly, The team that have more robots at the end of the game wins.

## custom maps

You can create your own custom map in `grid.c`. The map must be a square with all its edges set to `OBSTACLE` to avoid unexpected behaviours. If you just want to play with the size of the arena you can use `generate_grid.py`, otherwise you can write your own script to generate any map you want. Make sure to keep consistency between`GRID` and `GRID_SIZE`, and between `SPAWN_LOCATIONS` and `SPAWN_LEN`.