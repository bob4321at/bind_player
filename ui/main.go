package main

import (
	"bytes"
	"image/color"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.design/x/clipboard"
	"golang.org/x/text/language"
)

type Game struct{}

var typed = ""

var key_to_string = map[ebiten.Key]string{
	ebiten.KeyQ:      "Q",
	ebiten.KeyW:      "W",
	ebiten.KeyE:      "E",
	ebiten.KeyR:      "R",
	ebiten.KeyT:      "T",
	ebiten.KeyY:      "Y",
	ebiten.KeyU:      "U",
	ebiten.KeyI:      "I",
	ebiten.KeyO:      "O",
	ebiten.KeyP:      "P",
	ebiten.KeyA:      "A",
	ebiten.KeyS:      "S",
	ebiten.KeyD:      "D",
	ebiten.KeyF:      "F",
	ebiten.KeyG:      "G",
	ebiten.KeyH:      "H",
	ebiten.KeyJ:      "J",
	ebiten.KeyK:      "K",
	ebiten.KeyL:      "L",
	ebiten.KeyZ:      "Z",
	ebiten.KeyX:      "X",
	ebiten.KeyC:      "C",
	ebiten.KeyV:      "V",
	ebiten.KeyB:      "B",
	ebiten.KeyN:      "N",
	ebiten.KeyM:      "M",
	ebiten.Key1:      "1",
	ebiten.Key2:      "2",
	ebiten.Key3:      "3",
	ebiten.Key4:      "4",
	ebiten.Key5:      "5",
	ebiten.Key6:      "6",
	ebiten.Key7:      "7",
	ebiten.Key8:      "8",
	ebiten.Key9:      "9",
	ebiten.Key0:      "0",
	ebiten.KeySpace:  " ",
	ebiten.KeyPeriod: ".",
}

func (g *Game) Update() error {
	if !Downloading {
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			os.Exit(0)
		}

		var hit_keys []ebiten.Key

		if !ebiten.IsKeyPressed(ebiten.KeyControl) {
			hit_keys = inpututil.AppendJustPressedKeys(nil)
		}

		if hit_keys != nil {
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				typed += key_to_string[hit_keys[0]]
			} else {
				typed += strings.ToLower(key_to_string[hit_keys[0]])
			}
		}

		if len(typed) != 0 {
			if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
				typed = typed[:len(typed)-1]
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyDelete) {
			typed = ""
		}

		list_to_render = nil
		for _, song := range songs {
			if strings.Contains(strings.ToUpper(song), strings.ToUpper(typed)) {
				list_to_render = append(list_to_render, song)
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
			if len(list_to_render) > scroll+1 {
				scroll += 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				if len(list_to_render) > scroll+1 {
					scroll -= 2
				}
			}
		}

		if scroll > len(list_to_render) {
			scroll = 0
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			home_dir, err := os.UserHomeDir()
			if err != nil {
				panic(err)
			}
			err = os.Truncate(home_dir+"/Documents/current_song", 0)
			if err != nil {
				panic(err)
			}
			f, err := os.OpenFile(home_dir+"/Documents/current_song", os.O_WRONLY, 0644)
			f.WriteString(list_to_render[scroll] + "^")
			f.Close()

			os.Exit(0)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			if ebiten.IsKeyPressed(ebiten.KeyControl) {
				clipboardText := string(clipboard.Read(clipboard.FmtText))

				if strings.Contains(clipboardText, "youtube.com") {
					go DownloadSong(clipboardText, typed)
					Downloading = true
				}
			}
		}
	}

	return nil
}

var Downloading = false

var scroll = 0

var list_to_render []string

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 175})
	if !Downloading {
		text_op := text.DrawOptions{}
		text_op.GeoM.Translate(640, 16)
		text_op.PrimaryAlign = text.AlignCenter
		text.Draw(screen, typed, font_face, &text_op)

		for i, song_element := range list_to_render {
			if i == scroll {
				text_op.GeoM.Reset()
				text_op.GeoM.Translate(640, 16+32*float64(i)+64-float64(scroll*32))
				text.Draw(screen, ">"+song_element, list_font_face, &text_op)
			} else {
				text_op.GeoM.Reset()
				text_op.GeoM.Translate(640, 16+32*float64(i)+64-float64(scroll*32))
				text.Draw(screen, song_element, list_font_face, &text_op)
			}
		}
	} else {
		text_op := text.DrawOptions{}
		text_op.GeoM.Translate(640, 360/2)
		text_op.PrimaryAlign = text.AlignCenter
		text.Draw(screen, "Downloading", font_face, &text_op)
	}
}

func (g *Game) Layout(ow, oh int) (sw, sh int) {
	return 1280, 360
}

var font_face *text.GoTextFace
var list_font_face *text.GoTextFace

var songs []string

var selector_img, _, _ = ebitenutil.NewImageFromFile("./selector.png")

func DownloadSong(SongUrl, FileName string) {
	homepath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	FileName = homepath + "/Music/" + FileName
	download := exec.Command("yt-dlp", "-o", FileName, "-x", "--audio-format", "mp3", SongUrl)

	if err := download.Run(); err != nil {
		log.Fatal(err)
	}

	Downloading = false
}

func main() {
	home_path, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	ebiten.SetWindowSize(1280, 360)
	size_x, size_y := ebiten.Monitor().Size()
	ebiten.SetWindowPosition(size_x/4, size_y/3)
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)

	font_bytes, err := os.ReadFile(home_path + "/Documents/ComicRelief-Regular.ttf")
	if err != nil {
		panic(err)
	}

	font, err := text.NewGoTextFaceSource(bytes.NewBuffer(font_bytes))

	font_face = &text.GoTextFace{
		Source:    font,
		Direction: text.DirectionLeftToRight,
		Size:      48,
		Language:  language.AmericanEnglish,
	}

	list_font_face = &text.GoTextFace{
		Source:    font,
		Direction: text.DirectionLeftToRight,
		Size:      32,
		Language:  language.AmericanEnglish,
	}

	music_dir, err := os.ReadDir(home_path + "/Music")
	if err != nil {
		panic(err)
	}

	for _, f := range music_dir {
		if !f.Type().IsDir() {
			songs = append(songs, f.Name())
		}
	}

	if err := ebiten.RunGameWithOptions(&Game{}, &ebiten.RunGameOptions{ScreenTransparent: true}); err != nil {
		panic(err)
	}
}
