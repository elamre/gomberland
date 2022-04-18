package assets

import (
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/elamre/go_helpers/pkg/misc"
	"github.com/elamre/tentsuyu"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"image/color"
	"io/ioutil"
	"log"
	"strconv"
)

const (
	// Assetnames here
	GenericButtonAsset = "GenericButtonAsset"
	GenericButtonFont  = "GenericButtonFont"
	GenericInputImage  = "GenericInputAsset"
	GenericInputColor  = "GenericInputColor"
)

const (
	ImageFolder = "assets/images/"
	FontFolder  = "assets/fonts/"
)

const (
	textIdleColor               = "dff4ff"
	textDisabledColor           = "5a7a91"
	textInputCaretColor         = "e7c34b"
	textInputDisabledCaretColor = "766326"
)

var manager *tentsuyu.AssetsManager
var deferList []func()

func loadNineSlice(path string, w [3]int, h [3]int) (*image.NineSlice, error) {
	i, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, err
	}

	return image.NewNineSlice(i, w, h), nil
}

func textInputImages() (*widget.TextInputImage, error) {
	idle, _, err := ebitenutil.NewImageFromFile(ImageFolder + "text-input-idle.png")
	if err != nil {
		return nil, err
	}

	disabled, _, err := ebitenutil.NewImageFromFile(ImageFolder + "text-input-disabled.png")
	if err != nil {
		return nil, err
	}
	return &widget.TextInputImage{
		Idle:     image.NewNineSlice(idle, [3]int{9, 14, 6}, [3]int{9, 14, 6}),
		Disabled: image.NewNineSlice(disabled, [3]int{9, 14, 6}, [3]int{9, 14, 6}),
	}, nil
}

func hexToColor(h string) color.Color {
	u, err := strconv.ParseUint(h, 16, 0)
	if err != nil {
		panic(err)
	}

	return color.RGBA{
		R: uint8(u & 0xff0000 >> 16),
		G: uint8(u & 0xff00 >> 8),
		B: uint8(u & 0xff),
		A: 255,
	}
}

func textInputColor() *widget.TextInputColor {
	return &widget.TextInputColor{
		Idle:          hexToColor(textIdleColor),
		Disabled:      hexToColor(textDisabledColor),
		Caret:         hexToColor(textInputCaretColor),
		DisabledCaret: hexToColor(textInputDisabledCaretColor),
	}
}

func loadFont(path string, size float64) (font.Face, error) {
	fontData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ttfFont, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}

func LoadButtonImages() (*widget.ButtonImage, error) {
	idle, err := loadNineSlice(ImageFolder+"button-idle.png", [3]int{25, 12, 25}, [3]int{21, 0, 20})
	if err != nil {
		return nil, err
	}

	hover, err := loadNineSlice(ImageFolder+"button-hover.png", [3]int{25, 12, 25}, [3]int{21, 0, 20})
	if err != nil {
		return nil, err
	}

	pressed, err := loadNineSlice(ImageFolder+"button-pressed.png", [3]int{25, 12, 25}, [3]int{21, 0, 20})
	if err != nil {
		return nil, err
	}
	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func GetManager() *tentsuyu.AssetsManager {
	if manager == nil {
		deferList = make([]func(), 0)
		manager = tentsuyu.NewAssetsManager()
		manager.AssetMap[GenericButtonAsset] = misc.CheckErrorRetVal[*widget.ButtonImage](LoadButtonImages())
		manager.AssetMap[GenericInputColor] = textInputColor()
		manager.AssetMap[GenericInputImage] = misc.CheckErrorRetVal[*widget.TextInputImage](textInputImages())
		face := misc.CheckErrorRetVal[font.Face](loadFont(FontFolder+"NotoSans-Regular.ttf", 20))
		deferList = append(deferList, func() {
			if err := face.Close(); err != nil {
				log.Printf("err closing font: %s", err.Error())
			}
		})

		manager.AssetMap[GenericButtonFont] = face
	}
	return manager
}

func CleanUp() {
	mng := GetManager()
	for _, f := range deferList {
		f()
	}
	_ = mng
}
