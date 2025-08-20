function act(self, game)
    target = rg.toward(self.location, rg.CENTER_POINT)
    if target == self.location then
        return { actionType=GUARD } 
    end
    bot = game[target]
    if bot == nil then 
        return { actionType=MOVE, x=target.x, y=target.y }
    end
    if bot.player_id ~= self.player_id then 
        if self.hp <= rg.SETTINGS.attack_damage.max + rg.SETTINGS.collision_damage then
            return { actionType=SUICIDE } 
        end
    else
        for _, loc in pairs(rg.locs_around(self.location)) do
            bot = game.robots[loc]
            if bot ~= nil and bot.player_id ~= self.player_id then
                return { actionType=ATTACK, x=loc.x, y=loc.y }
            end
        end
        for _, loc in pairs(rg.locs_around(self.location, { OBSTACLE })) do
            for _, loc2 in pairs(rg.locs_around(loc)) do
                bot = game.robots[loc2]
                if bot ~= nil and bot.player_id ~= self.player_id then
                    return { actionType=ATTACK, x=loc.x, y=loc.y }
                end
            end
        end
    end
    return { actionType=GUARD }
end