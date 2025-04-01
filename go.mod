module github.com/igadmg/goecs

go 1.24

replace (
	deedles.dev/xiter => ../../pkg/xiter
	github.com/hymkor/go-lazy => ../../pkg/go-lazy
	github.com/igadmg/goel => ../../cmd/goel
	github.com/igadmg/goex => ../../pkg/goex
	github.com/igadmg/gogen => ../../pkg/gogen
	github.com/igadmg/gamemath => ../../pkg/gamemath
)

require (
	deedles.dev/xiter v0.2.1
	github.com/hymkor/go-lazy v0.5.0
	github.com/igadmg/goex v0.0.0-20250321131421-ccb743b21181
	github.com/igadmg/gogen v0.0.0-20250318100828-211ca23c6b9f
	github.com/igadmg/gamemath v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/chewxy/math32 v1.11.1 // indirect
	github.com/igadmg/goel v0.0.0-20250302140239-96fa936747cc // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/tools v0.31.0 // indirect
)

tool github.com/igadmg/goel
