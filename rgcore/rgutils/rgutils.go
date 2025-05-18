package rgutils

import "C"

import (
	"errors"
	"math"

	"github.com/ClementTariel/rg-lua/rgcore/rgconst"
	"github.com/ClementTariel/rg-lua/rgcore/rgentities"
)

func Abs(n int) int {
	if n < 0 {
		return -n
	} else {
		return n
	}
}

func WalkDist(x1, y1, x2, y2 int) int {
	return int(Abs(x1-x2) + Abs(y1-y2))
}

func Dist(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(float64((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)))
}

func Allies(playerId int, bots []rgentities.Bot) []rgentities.Bot {
	allies := []rgentities.Bot{}
	for _, bot := range bots {
		if bot.PlayerId == playerId {
			ally := bot
			allies = append(allies, ally)
		}
	}
	return allies
}

func Enemies(playerId int, bots []rgentities.Bot) []rgentities.Bot {
	enemies := []rgentities.Bot{}
	for _, bot := range bots {
		if bot.PlayerId != playerId {
			enemy := bot
			enemy.Id = 0
			enemies = append(enemies, enemy)
		}
	}
	return enemies
}

func FilterOutDeadBots(bots map[int]rgentities.Bot) []rgentities.Bot {
	filteredBots := []rgentities.Bot{}
	for _, bot := range bots {
		if bot.Hp > 0 {
			filteredBots = append(filteredBots, bot)
		}
	}
	return filteredBots
}

func GetLocationType(x int, y int) rgentities.LocationType {
	return rgentities.LocationType(rgconst.GRID[C.int(x*rgconst.GRID_SIZE+y)])
}

func FilterOutBotsOnSpawn(bots []rgentities.Bot) []rgentities.Bot {
	filteredBots := []rgentities.Bot{}
	for _, bot := range bots {
		if GetLocationType(bot.X, bot.Y) != rgconst.SPAWN {
			filteredBots = append(filteredBots, bot)
		}
	}
	return filteredBots
}

func GetSpawnLocation(i int) (rgentities.Location, error) {
	if i < 0 || i >= len(rgconst.SPAWN_LOCATIONS) {
		return rgentities.Location{X: -1, Y: -1}, errors.New("spawn index out of range in spawn generation")
	}
	return rgentities.Location{X: int(rgconst.SPAWN_LOCATIONS[C.int(i)].X), Y: int(rgconst.SPAWN_LOCATIONS[C.int(i)].Y)}, nil
}
