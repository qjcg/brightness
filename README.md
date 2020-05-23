# backlight

A small CLI tool for setting screen backlight brightness on Linux via sysfs.

## Features

- print current brightness
- set brightness as percentage
- increase/decrease brightness by increment

## Install

```sh
go get github.com/qjcg/brightness
```

## Use

```sh
# Print current backlight brightness (as percentage of maximum).
root@demo # brightness
50

# Set backlight brightness to 33%.
root@demo # brightness 33

# Increase backlight brightness by 5%.
root@demo # brightness +5

# Decrease backlight brightness by 10%.
root@demo # brightness -10
```

## Test

```sh
go test -v
```

## Related Projects

- [Hummer12007/brightnessctl](https://github.com/Hummer12007/brightnessctl)
- [multiplexd/brightlight](https://github.com/multiplexd/brightlight)
- [Merovius/backlight](https://github.com/Merovius/brightness)
- [ungerik/go-sysfs](https://github.com/ungerik/go-sysfs)
