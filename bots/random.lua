function act(self, game)
	d = math.random(1,3)
    if (d == 1) then
        return { actionType=GUARD } 
    end
    locs = rg.locs_around(self.location, { OBSTACLE })
    loc = locs[math.random(1, #locs)]
    if (d == 2) then
        return  { actionType=MOVE, x=loc.x, y=loc.y } 
    end
	return { actionType=ATTACK, x=loc.x, y=loc.y } 
end