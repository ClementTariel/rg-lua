package rgcore

import "math"

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

func Allies(playerId int, bots []Bot) []Bot {
	allies := []Bot{}
	for _, bot := range bots {
		if bot.PlayerId == playerId {
			ally := bot
			allies = append(allies, ally)
		}
	}
	return allies
}

func Enemies(playerId int, bots []Bot) []Bot {
	enemies := []Bot{}
	for _, bot := range bots {
		if bot.PlayerId != playerId {
			enemy := bot
			enemy.Id = 0
			enemies = append(enemies, enemy)
		}
	}
	return enemies
}

func FilterOutBotsOnSpawn(bots []Bot) []Bot {
	filteredBots := []Bot{}
	for _, bot := range bots {
		if GetLocationType(bot.X, bot.Y) != SPAWN {
			filteredBots = append(filteredBots, bot)
		}
	}
	return filteredBots
}

func FilterOutDeadBots(bots map[int]Bot) []Bot {
	filteredBots := []Bot{}
	for _, bot := range bots {
		if bot.Hp > 0 {
			filteredBots = append(filteredBots, bot)
		}
	}
	return filteredBots
}
