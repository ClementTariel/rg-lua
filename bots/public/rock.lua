function act(self, game)
    target = rg.toward(self.location, rg.CENTER_POINT)
    if target == self.location then
        return { actionType=GUARD } 
    end
    bot = game[target]
    if bot == nil then 
        return { actionType=MOVE, x=target.x, y=target.y }
    end
    return { actionType=GUARD }
end