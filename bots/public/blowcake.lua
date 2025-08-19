------------------------------------------------------------------
--
-- converted from python to lua
-- source: 
-- https://github.com/mpeterv/robotgame-bots/blob/master/blowcake.py 
--
------------------------------------------------------------------

-- blowcake by sne11ius
-- http://robotgame.org/viewrobot/3666

function act(self, game)
    function get_center(self)
        centers = { rg.Loc(9, 4), rg.Loc(14, 9), rg.Loc(9, 14), rg.Loc(4, 9) }
        -- centers = { rg.Loc(9, 4), rg.Loc(9, 14) }
        min_distance = math.maxinteger
        center = centers[1]
        for _, c in pairs(centers) do
            dist = rg.wdist(self.location, c)
            if dist < min_distance then
                min_distance = dist 
                center = c
            end
        end
        return center
    end
    
    function num_enemies(self, game)
        enemies = 0
        for location, bot in pairs(game.robots) do
            if bot.player_id ~= self.player_id then
                print("test")
                if 1 == rg.wdist(self.location, location) then
                    enemies = enemies + 1
                end
            end
        end
        return enemies
    end
    
    function num_frieds(self, game)
        friends = 0
        for location, bot in pairs(game.robots) do
            if bot.player_id == self.player_id then
                if 1 == rg.wdist(self.location, location) then
                    friends = friends + 1
                end
            end
        end
        return friends
    end

    num_enemies = num_enemies(self, game)
    if num_enemies * 9 > self.hp then
        return { actionType=SUICIDE }
    end

    min_distance = math.maxinteger
    move_to = get_center(self)
    for location, bot in pairs(game.robots) do
        if bot.player_id ~= self.player_id then
            if rg.dist(location, self.location) <= 1 then
                return { actionType=ATTACK, x=location.x, y=location.y }
            end
        end
        if bot.player_id == self.player_id then
            if rg.wdist(location, self.location) < min_distance then
                min_distance = rg.wdist(location, self.location)
                move_to = location
            end
        end
    end
    if min_distance < 2 then
        move_to = get_center(self)
    end
    
    if self.location == get_center(self) then
        return { actionType=GUARD }
    end
    
    if num_frieds(self, game) > 1 then 
        return { actionType=GUARD }
    end
    
    res = rg.toward(self.location, move_to)
    return { actionType=MOVE, x=res.x, y=res.y }
end