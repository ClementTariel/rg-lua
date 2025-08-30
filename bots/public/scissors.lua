function act(self, game)
    toCenter = rg.toward(self.location, rg.CENTER_POINT)
    shortestDist = rg.GRID_SIZE * 4
    target = toCenter
    for loc, bot in pairs(game.robots) do
        if bot.player_id ~= self.player_id then
            dist = rg.wdist(self.location, loc) + rg.wdist(loc, rg.CENTER_POINT)
            if dist < shortestDist then
                shortestDist = dist
                target = loc
            end
        end
    end
    destination = rg.toward(target, rg.CENTER_POINT)
    if destination == self.location then
        return { actionType=ATTACK, x=target.x, y=target.y }
    end
    next_step = rg.toward(self.location, destination)
    return { actionType=MOVE, x=next_step.x, y=next_step.y }
end