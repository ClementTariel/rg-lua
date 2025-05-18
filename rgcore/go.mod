module github.com/ClementTariel/rg-lua/rgcore

go 1.21.4

require (
	github.com/ClementTariel/rg-lua/rgcore/rgconst v0.0.0
	github.com/ClementTariel/rg-lua/rgcore/rgdebug v0.0.0
	github.com/ClementTariel/rg-lua/rgcore/rgentities v0.0.0
	github.com/ClementTariel/rg-lua/rgcore/rgerrors v0.0.0
	github.com/ClementTariel/rg-lua/rgcore/rgutils v0.0.0
)

replace (
	github.com/ClementTariel/rg-lua/rgcore/rgconst v0.0.0 => ./rgconst
	github.com/ClementTariel/rg-lua/rgcore/rgdebug v0.0.0 => ./rgdebug
	github.com/ClementTariel/rg-lua/rgcore/rgentities v0.0.0 => ./rgentities
	github.com/ClementTariel/rg-lua/rgcore/rgerrors v0.0.0 => ./rgerrors
	github.com/ClementTariel/rg-lua/rgcore/rgutils v0.0.0 => ./rgutils
)
