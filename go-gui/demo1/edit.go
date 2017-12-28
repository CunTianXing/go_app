package main

import "github.com/murlokswarm/app"

type EditMenu struct{}

func (e *EditMenu) Render() string {
	return `<div class="Content Gary">Hello world</div>`
}

func init() {
	app.RegisterComponent(&EditMenu{})
}
