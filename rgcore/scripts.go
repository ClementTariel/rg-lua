package rgcore

import (
	"fmt"

	"github.com/ClementTariel/rg-lua/rgcore/rgconst"
)

func GetInitialisationScript() string {
	return `MOVE = ` + fmt.Sprintf("%d", rgconst.MOVE) + `
ATTACK = ` + fmt.Sprintf("%d", rgconst.ATTACK) + `
GUARD = ` + fmt.Sprintf("%d", rgconst.GUARD) + `
SUICIDE = ` + fmt.Sprintf("%d", rgconst.SUICIDE) + `
NORMAL = ` + fmt.Sprintf("%d", rgconst.NORMAL) + `
SPAWN = ` + fmt.Sprintf("%d", rgconst.SPAWN) + `
OBSTACLE = ` + fmt.Sprintf("%d", rgconst.OBSTACLE) + `
rg.SETTINGS = {
	spawn_delay = ` + fmt.Sprintf("%d", rgconst.SPAWN_DELAY) + `,
	spawn_count = ` + fmt.Sprintf("%d", rgconst.SPAWN_COUNT) + `,
	robot_hp = ` + fmt.Sprintf("%d", rgconst.MAX_HP) + `,
	attack_range = ` + fmt.Sprintf("%d", rgconst.ATTACK_RANGE) + `,
	attack_damage = { ` +
		`min=` + fmt.Sprintf("%d", rgconst.ATTACK_DAMAGE_MIN) + `, ` +
		`max=` + fmt.Sprintf("%d", rgconst.ATTACK_DAMAGE_MAX) +
		` },
	suicide_damage = ` + fmt.Sprintf("%d", rgconst.SUICIDE_DAMAGE) + `,
	collision_damage = ` + fmt.Sprintf("%d", rgconst.COLLISION_DAMAGE) + `,
	max_turn = ` + fmt.Sprintf("%d", rgconst.MAX_TURN) + `,
}
__RG_CORE_SYSTEM = {
	self = {},
	game = {},
}
`
}

func GetLoadActScript() string {
	return `__RG_CORE_SYSTEM.act = act`
}

func GetResetScript(turn int) string {
	const resetScript = `__RG_CORE_SYSTEM.game.robots = rg.Robots()
for i = 0,%[1]d-1,1 do
	__RG_CORE_SYSTEM.game.robots[i] = {}
end
__RG_CORE_SYSTEM.game.turn = %[2]d
`
	return fmt.Sprintf(resetScript, rgconst.GRID_SIZE, turn)
}

func GetLoadBotScript(botX int, botY int, botHP int, playerId int, botId string) string {
	const loadBotScript = `__RG_CORE_SYSTEM.game.robots[%[1]d][%[2]d] = {
	location = rg.Loc(%[1]d, %[2]d),
	hp = %[3]d,
	player_id = %[4]d,
	id = %[5]s,
}
`
	return fmt.Sprintf(loadBotScript, botX, botY, botHP, playerId, botId)
}

func GetLoadSelfScript(botX int, botY int, botHP int, playerId int, botId int) string {
	const loadSelfScript = `if __RG_CORE_SYSTEM.self[%[5]d] == nil then
	__RG_CORE_SYSTEM.self[%[5]d] = {}
end
__RG_CORE_SYSTEM.self[%[5]d].id = %[5]d
__RG_CORE_SYSTEM.self[%[5]d].location = rg.Loc(%[1]d, %[2]d)
__RG_CORE_SYSTEM.self[%[5]d].hp = %[3]d
__RG_CORE_SYSTEM.self[%[5]d].player_id = %[4]d
__RG_CORE_SYSTEM.self[%[5]d].id = %[5]d
`
	return fmt.Sprintf(loadSelfScript, botX, botY, botHP, playerId, botId)
}
