package menu_state

import (
	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/elamre/gomberman/assets"
	"github.com/elamre/tentsuyu"
	"golang.org/x/image/font"
	"image/color"
	"log"
)

type mainMenu struct {
	rootContainer *widget.Container
	ui            ebitenui.UI
	curState      string
	buttons       []*widget.Button
}

func newMainMenu() *mainMenu {
	m := mainMenu{curState: mainMenuState}
	buttonAsset := assets.GetManager().AssetMap[assets.GenericButtonAsset].(*widget.ButtonImage)
	widgetOptions := widget.ButtonOpts.WidgetOpts(
		// instruct the container's anchor layout to center the button both horizontally and vertically
		widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		}),
		widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		}),
	)
	padding := widget.ButtonOpts.TextPadding(widget.Insets{
		Left:  30,
		Right: 30,
	})
	log.Printf("%+v", buttonAsset)
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.RGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Padding(widget.Insets{
				Left:  25,
				Right: 25,
			}),
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{false, true}, []bool{true}),
			widget.GridLayoutOpts.Spacing(20, 0))),
	)
	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(50),
			widget.RowLayoutOpts.Padding(widget.Insets{Top: 50}),
		)))

	// construct a button
	singlePlayer := widget.NewButton(
		widgetOptions,
		widget.ButtonOpts.Image(buttonAsset),
		widget.ButtonOpts.Text("SinglePlayer", assets.GetManager().AssetMap[assets.GenericButtonFont].(font.Face), &widget.ButtonTextColor{
			Idle: color.RGBA{R: 0xdf, G: 0xf4, B: 0xff, A: 0xff},
		}),
		padding,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			m.curState = singlePlayerMenuState
		}),
	)
	multiplayer := widget.NewButton(
		widgetOptions,
		widget.ButtonOpts.Image(buttonAsset),
		widget.ButtonOpts.Text("Multiplayer", assets.GetManager().AssetMap[assets.GenericButtonFont].(font.Face), &widget.ButtonTextColor{
			Idle: color.RGBA{R: 0xdf, G: 0xf4, B: 0xff, A: 0xff},
		}),
		padding,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			m.curState = multiPlayerMenuState
		}),
	)
	options := widget.NewButton(
		widgetOptions,
		widget.ButtonOpts.Image(buttonAsset),
		widget.ButtonOpts.Text("Options", assets.GetManager().AssetMap[assets.GenericButtonFont].(font.Face), &widget.ButtonTextColor{
			Idle: color.RGBA{R: 0xdf, G: 0xf4, B: 0xff, A: 0xff},
		}),
		padding,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			m.curState = optionsMenuState
		}),
	)
	exit := widget.NewButton(
		widgetOptions,
		widget.ButtonOpts.Image(buttonAsset),
		widget.ButtonOpts.Text("Exit", assets.GetManager().AssetMap[assets.GenericButtonFont].(font.Face), &widget.ButtonTextColor{
			Idle: color.RGBA{R: 0xdf, G: 0xf4, B: 0xff, A: 0xff},
		}),
		padding,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			m.curState = BackState
		}),
	)
	buttonsContainer.AddChild(singlePlayer)
	buttonsContainer.AddChild(multiplayer)
	buttonsContainer.AddChild(options)
	buttonsContainer.AddChild(exit)
	// add the button as a child of the container
	rootContainer.AddChild(buttonsContainer)

	// construct the UI
	m.ui = ebitenui.UI{
		Container: rootContainer,
	}
	return &m
}

func (m *mainMenu) GetState() string {
	return m.curState
}

func (m *mainMenu) Init() any {
	return nil
}

func (m *mainMenu) Update() {
	m.ui.Update()
}

func (m *mainMenu) Draw(game *tentsuyu.Game) {
	m.ui.Draw(game.Screen)
}

func (m *mainMenu) DeInit() any {
	return nil
}
