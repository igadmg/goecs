module github.com/igadmg/goecs

go 1.24

replace (
	deedles.dev/xiter => ../../pkg/xiter
	github.com/hymkor/go-lazy => ../../pkg/go-lazy
	github.com/igadmg/gamemath => ../../pkg/gamemath
	github.com/igadmg/goel => ../../cmd/goel
	github.com/igadmg/goex => ../../pkg/goex
	github.com/igadmg/gogen => ../../pkg/gogen
)

require (
	deedles.dev/xiter v0.2.1
	github.com/hymkor/go-lazy v0.5.0
	github.com/igadmg/gamemath v0.0.0-20250410220553-87672598e3a6
	github.com/igadmg/goex v0.0.0-20250407220752-712c023573b8
	github.com/igadmg/gogen v0.0.0-20250410220610-3903be0a0ba3
	gopkg.in/yaml.v3 v3.0.1
	gonum.org/v1/gonum v0.16.0
)

require (
	github.com/chewxy/math32 v1.11.1 // indirect
	github.com/igadmg/goel v0.0.0-20250410203636-1e64bb5aa9ed // indirect
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/tools v0.32.0 // indirect
)

tool github.com/igadmg/goel
