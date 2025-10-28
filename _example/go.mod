module github.com/denizgursoy/inpu/_example

go 1.25.0

require (
	github.com/denizgursoy/inpu v1.2.1-0.20251028191331-bf83c2e648bd
	github.com/denizgursoy/inpu/loggers/zero v0.0.0-20251028191331-bf83c2e648bd
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
)

replace github.com/denizgursoy/inpu => ..
