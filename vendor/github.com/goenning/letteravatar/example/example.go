package main

import (
	"image/png"
	"log"
	"os"

	"github.com/goenning/letteravatar"
)

var names = []string{
	"Jon Snow",
	"Bob Marley",
	"Carol Stark",
	"Dave",
	"Eve",
	"Frank Sinatra",
	"Gloria Pires",
	"Henry",
	"Isabella",
	"Mad Monkey",
	"James",
	"Жозефина",
	"Ярослав",
}

func main() {
	for _, name := range names {

		img, err := letteravatar.Draw(75, letteravatar.Extract(name), nil)
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(name + ".png")
		if err != nil {
			log.Fatal(err)
		}

		err = png.Encode(file, img)
		if err != nil {
			log.Fatal(err)
		}
	}
}
