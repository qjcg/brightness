// +build packages

//go:generate go build
//go:generate upx brightness
//go:generate holo-build --force --format=pacman holo.toml
//go:generate rm brightness
package main
