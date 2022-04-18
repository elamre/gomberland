package game

const AssetsFolder = "assets/"
const FontsFolder = AssetsFolder + "fonts/"
const ImagesFolder = AssetsFolder + "images/"

type Settings struct {
	MapIdentifier string
	PlayerNames   []string
	GameMode      string
}
