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
	github.com/igadmg/gamemath v0.0.0-20250401194741-5c7111c9b66d
	github.com/igadmg/goex v0.0.0-20250325133153-61aee7990ef8
	github.com/igadmg/gogen v0.0.0-20250327140622-8f1ad806cdc3
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/chewxy/math32 v1.11.1 // indirect
	github.com/igadmg/goel v0.0.0-20250325114530-b9ac5d099412 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/tools v0.31.0 // indirect
)

tool github.com/igadmg/goel
