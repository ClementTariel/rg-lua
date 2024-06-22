function act(self, game)
	d = math.random(1,12)
    if (d <= 4) then
        return { actionType=GUARD } 
    end
    loc = rg.locs_around(self.location)[1+(d%4)]
    if (d <= 8) then
        return  { actionType=MOVE, x=loc.x, y=loc.y } 
    end
	return { actionType=ATTACK, x=loc.x, y=loc.y } 
end