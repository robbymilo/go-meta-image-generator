package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/pug"
)

func main() {
	engine := pug.New("views", ".pug")
	engine.Debug(true)
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get(":file.jpg", func(c *fiber.Ctx) error {
		img := screenshot(c.Params("file"))
		c.Type("jpg")
		return c.Send([]byte(img))
	})

	app.Get("/render/:file.jpg", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		return c.Render("index", fiber.Map{
			"message": c.Params("file"),
			"background": "https://grafana.com/products/assets/cloud-grafana-0.png",
		}, "layouts/main")
	})

	log.Fatal(app.Listen(":3000"))
}

func screenshot(file string) []byte {
	name := "cache/" + file + ".jpg"

	if fileExists(name) {

		content, err := ioutil.ReadFile(name)

		if err != nil {
			log.Fatal(err)
		}

		return content

	} else {
		page := rod.New().MustConnect().MustPage("http://localhost:3000/render/" + file + ".jpg").MustWaitLoad()

		img, _ := page.Screenshot(true, &proto.PageCaptureScreenshot{
			Format:  proto.PageCaptureScreenshotFormatJpeg,
			Quality: 75,
			Clip: &proto.PageViewport{
				X:      0,
				Y:      0,
				Width:  1200,
				Height: 700,
				Scale:  1,
			},
		})

		_ = utils.OutputFile(name, img)
		return img

	}
}

// https://golangcode.com/check-if-a-file-exists/
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
