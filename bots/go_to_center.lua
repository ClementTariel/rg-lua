function act(self, game)
    locs = rg.locs_around(self.location)
    loc = rg.toward(self.location, rg.CENTER_POINT)
    if (game[loc] ~= nil) then
        return { actionType=GUARD } 
    end
	return { actionType=MOVE, x=loc.x, y=loc.y } 
end