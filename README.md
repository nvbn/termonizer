# Simple terminal organizer

![screenshot](https://raw.githubusercontent.com/nvbn/termonizer/master/screenshot.png)

[Download latest release](https://github.com/nvbn/termonizer/releases) (only macos arm)

## Hotkeys

Esc Esc - exit

Navigation:
* ⌥↑ - future/up goal
* ⇧⌥↑ - current/first goal
* ⌥↓ - past/down goal
* ⇧⌥← - longer/left period
* ⇧⌥→ - shorter/right period

Zooming:
* ⌥+ - zoom in / decrease the amount of visible goals
* ⌥- - zoom out / increase the amount of visible goals

Text editing:
* ⌃C - copy
* ⌃X - cut
* ⌃V - paste
* ⌃A - select all
* ⌃Z - undo
* Esc - remove selection

## Development

Prerequisites:
* [go 1.23+](https://go.dev/doc/install) 
* 

Run:
```
make run
```

Test:
```
make test
```

Generate test data:
```
make generate-test-db
```

Run with pre-generated test data:
```
run-test-db
```

Output logs:
```
make log
```

Look into the Makefile for more dev commands.

## [License MIT](https://github.com/nvbn/termonizer/blob/main/LICENSE)
