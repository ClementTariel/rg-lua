function act(self, game)
    enemyDamage = 0
    for _, bot in pairs(rg.locs_around(self.location)) do
        if bot.player_id ~= self.player_id then
            enemyDamage = enemyDamage + rg.SETTINGS.attack_damage.min
        end
    end
    if self.hp <= enemyDamage or self.hp <= rg.SETTINGS.attack_damage.min then
        return { actionType=SUICIDE }
    end
    toCenter = rg.toward(self.location, rg.CENTER_POINT)
    shortestDist = rg.GRID_SIZE * 2
    target = toCenter
    for _, bot in pairs(game.robots) do
        if bot.player_id ~= self.player_id then
            loc = rg.toward(bot.location, rg.CENTER_POINT)
            if loc == self.location then
                return { actionType=ATTACK, x=bot.location.x, y=bot.location.y }
            end
            dist = rg.wdist(self.location, loc)
            if dist < shortestDist then
                shortestDist = dist
                target = rg.toward(self.location, loc)
            end
        end
    end
    return { actionType=MOVE, x=target.x, y=target.y }
end