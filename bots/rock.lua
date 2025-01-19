function act(self, game)
    target = rg.toward(self.location, rg.CENTER_POINT)
    bot = game[target]
    if target == self.location or (bot ~= nil) then
        return { actionType=GUARD } 
    end
	return { actionType=MOVE, x=target.x, y=target.y } 
end