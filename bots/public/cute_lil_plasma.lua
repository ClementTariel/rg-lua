------------------------------------------------------------------
--
-- converted from python to lua
-- source: 
-- https://github.com/mpeterv/robotgame-bots/blob/master/Cute%20Lil%27%20Plasma.py 
--
------------------------------------------------------------------

function act(self,game)
    function includes(list, element)
        for _, e in pairs(list) do
            if e == element then
                return true
            end
        end
        return false
    end
    min_dist = 99
    goal = rg.CENTER_POINT
    to_center = rg.toward(self.location, rg.CENTER_POINT)
    move_to_center = { actionType=MOVE, x=to_center.x, y=to_center.y }
    enemy_count = 0
    enemies = {}
    if (game.turn % 10) == 0 and rg.loc_type(self.location) == SPAWN then
        return move_to_center
    end
    adjacents = rg.locs_around(self.location, { OBSTACLE, SPAWN })
    for loc, bot in pairs(game.robots) do
        if bot.player_id ~= self.player_id then
            if rg.wdist(loc, self.location) < min_dist then
                min_dist = rg.wdist(loc, self.location)
                goal = loc
            end
            table.insert(enemies, loc)
            if rg.wdist(loc, self.location) == 1 then 
                enemy_count = enemy_count + 1
            end
        end
    end
    toward_goal = rg.toward(self.location, goal)
    dist_to_goal = rg.wdist(self.location, goal)
    goal_at_range = (dist_to_goal == 1)
    if game.robots[self.location].hp < enemy_count * 10 and goal_at_range then
        loc = rg.Loc(2 * self.location.x - toward_goal.x, 2 * self.location.y - toward_goal.y)
        if includes(adjacents, loc) and not game.robots[loc] then 
            return { actionType=MOVE, x=loc.x, y=loc.y } 
        end
    end
    if goal_at_range and game.robots[self.location].hp > 10 * enemy_count or dist_to_goal == 2 and game.robots[self.location].hp < 16 then
        if includes(enemies, toward_goal) and game.robots[toward_goal].hp < 6 and self.hp > 5 then
            return { actionType=MOVE, x=toward_goal.x, y=toward_goal.y }
        end
        return { actionType=ATTACK, x=toward_goal.x, y=toward_goal.y }
    end
    if includes(adjacents, toward_goal) and not game.robots[toward_goal] then
        return { actionType=MOVE, x=toward_goal.x, y=toward_goal.y }
    end
    if includes(game.robots, rg.Loc(move_to_center.x, move_to_center.y)) then
        for _, loc in pairs(adjacents) do
            if not game.robots[loc] then
                return { actionType=MOVE, x=loc.x, y=loc.y } 
            end
        end
    end
    return move_to_center
end