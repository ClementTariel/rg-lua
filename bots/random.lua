function act(self, game)
	d = math.random(0,11)
    if (d < 4) then
        return { actionType=GUARD } 
    end
	x = self.location.x + (d%2) * ((d-2)%4)
	y = self.location.y + ((d+1)%2) * ((d-1)%4)
    if (d < 8) then
        return  { actionType=MOVE, x=x, y=y } 
    end
	return { actionType=ATTACK, x=x, y=y } 
end