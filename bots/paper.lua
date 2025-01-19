function act(self, game)
    toCenter = rg.toward(self.location, rg.CENTER_POINT)
    target = self.location
    for _, loc in pairs(rg.locs_around(self.location, { OBSTACLE })) do
        bot = game.robots[loc]
        if rg.loc_type(loc) == SPAWN or (bot ~= nil and bot.player_id == self.player_id) then
            return { actionType=MOVE, x=toCenter.x, y=toCenter.y }
        elseif rg.toward(loc, rg.CENTER_POINT) == self.location then
            target = loc
        end
    end
    if target == self.location then
        return { actionType=GUARD}
    end
    return { actionType=ATTACK, x=target.x, y=target.y }
end