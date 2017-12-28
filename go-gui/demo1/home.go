package main

import "github.com/murlokswarm/app"

type Home struct{}

func (h *Home) Render() string {
	return `<div class="Home">
  	<div class="Example">
  		<h1>Copy/paste</h1>
  		<ul oncontextmenu="OnContextMenu">
  			<li>Select me</li>
  			<li>Right click</li>
  			<li>Copy</li>
  		</ul>
      <a href="EditMenu">Go to San Francisco</a>
  	</div>

  	<div class="Example">
  		<h1>Custom menu</h1>
  		<button onclick="OnButtonClick">Show</button>
  	</div>
  </div>`
}

func init() {
	app.RegisterComponent(&Home{})
}
