package main

import (
	"GolangMessanger/internal/app"
)

func main() {
	c := app.NewContext()
	if err := c.Run(); err != nil {
		panic(err)
	}
}
