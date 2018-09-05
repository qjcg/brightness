// +build packages

//go:generate go build
//go:generate holo-build --force --format=rpm holo.toml
//go:generate holo-build --force --format=pacman holo.toml
//go:generate holo-build --force --format=debian holo.toml
//go:generate rm brightness
package main
