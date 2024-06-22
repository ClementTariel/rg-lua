function act(self, game)
    locs = rg.locs_around(self.location)

    i = 1
    while locs[i] ~= nil do
        loc = locs[i]
        robot = game[loc.x]
        if robot ~= nil and robot[loc.y] ~= nil then
            robot = robot[loc.y]
            if robot.player_id ~= self.player_id then
                return { actionType=ATTACK, x=loc.x, y=loc.y }
            end
        end
        i = i + 1
    end
    possibleMoveCount = i - 1

    d = math.random(1, possibleMoveCount)
    moveTo = locs[d]
    return { actionType=MOVE, x=moveTo.x, y=moveTo.y }
end