package android

import (
	"github.com/CunTianXing/go_app/go-mobile/view/ios"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("github.com/CunTianXing/go_app/go-mobile/view/android NewStackView", func() view.View {
		return NewStackView()
	})
}

func NewStackView() view.View {
	return ios.NewStackView()
}
