module github.com/igadmg/goecs

go 1.24

replace github.com/igadmg/gogen => ../../pkg/gogen

require (
	deedles.dev/xiter v0.2.1
	github.com/hymkor/go-lazy v0.5.0
	github.com/igadmg/gamemath v0.0.0-20250410222204-28d83654fdf2
	github.com/igadmg/goex v0.0.0-20250502115452-bd40b01ba4eb
	github.com/igadmg/gogen v0.0.0-20250502134002-bf191499d781
	gonum.org/v1/gonum v0.16.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/chewxy/math32 v1.11.1 // indirect
	github.com/igadmg/goel v0.0.0-20250502134036-a60922f656ed // indirect
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/tools v0.32.0 // indirect
)

tool github.com/igadmg/goel
