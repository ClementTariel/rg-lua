module github.com/ClementTariel/rg-lua/player

go 1.23.9

require (
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

require (
	github.com/ClementTariel/rg-lua/rgcore v0.0.0
	github.com/ClementTariel/rg-lua/rgcore/rgconst v0.0.0
	github.com/ClementTariel/rg-lua/rgcore/rgdebug v0.0.0
	github.com/ClementTariel/rg-lua/rgcore/rgentities v0.0.0
	github.com/ClementTariel/rg-lua/rgcore/rgerrors v0.0.0
)

require github.com/ClementTariel/rg-lua/rgcore/rgutils v0.0.0 // indirect

replace (
	github.com/ClementTariel/rg-lua/rgcore v0.0.0 => ../rgcore
	github.com/ClementTariel/rg-lua/rgcore/rgconst v0.0.0 => ../rgcore/rgconst
	github.com/ClementTariel/rg-lua/rgcore/rgdebug v0.0.0 => ../rgcore/rgdebug
	github.com/ClementTariel/rg-lua/rgcore/rgentities v0.0.0 => ../rgcore/rgentities
	github.com/ClementTariel/rg-lua/rgcore/rgerrors v0.0.0 => ../rgcore/rgerrors
	github.com/ClementTariel/rg-lua/rgcore/rgutils v0.0.0 => ../rgcore/rgutils
)
