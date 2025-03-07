module github.com/igadmg/goecs

go 1.24

replace deedles.dev/xiter => ../../pkg/xiter

replace github.com/hymkor/go-lazy => ../../pkg/go-lazy

replace github.com/igadmg/goel => ../../cmd/goel

replace github.com/igadmg/goex => ../../pkg/goex

replace github.com/igadmg/gogen => ../../pkg/gogen

replace github.com/igadmg/raylib-go/raymath => ../../pkg/raylib-go/raymath

tool github.com/igadmg/goel

require (
	deedles.dev/xiter v0.1.1
	github.com/hymkor/go-lazy v0.0.0-00010101000000-000000000000
	github.com/igadmg/goex v0.0.0-20250226161117-f8fd602fe0c7
	github.com/igadmg/gogen v0.0.0-00010101000000-000000000000
	github.com/igadmg/raylib-go/raymath v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/igadmg/goel v0.0.0-20250123180020-ee1a98205fb0 // indirect
	golang.org/x/exp v0.0.0-20250228200357-dead58393ab7 // indirect
	golang.org/x/mod v0.23.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/tools v0.30.0 // indirect
)
