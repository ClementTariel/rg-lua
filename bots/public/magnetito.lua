function act(self, game)
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

    -- increase score when bad tiles are around
    function compute_score(current_loc, game)
        score = rg.wdist(current_loc, rg.CENTER_POINT)
        radius = 2
        for i=-radius, radius do
            for j=-radius, radius do
                loc = rg.Loc(current_loc.x + i, current_loc.y + j)
                weight = 2 * radius + 1 - rg.wdist(current_loc, loc)
                if rg.loc_type(loc) == OBSTACLE then
                    score = score + 3 * weight
                elseif rg.loc_type(loc) == SPAWN then
                    score = score + 2 * weight
                end
                if game.robots[loc] ~= nil then
                    score = score + 1 * weight
                    if game.robots[loc].player_id ~= self.player_id then
                        score = score + 2 * weight
                    end
                end
            end
        end
        return score
    end

    min_score = compute_score(self.location, game)
    target = nil
    for _, loc in pairs(rg.locs_around(self.location, { OBSTACLE })) do
        target_score = compute_score(loc, game)
        if target_score < min_score then
            min_score = target_score
            target = loc
        end
    end

    if target == nil then
        return { actionType=GUARD } 
    end
	return { actionType=MOVE, x=target.x, y=target.y } 
end