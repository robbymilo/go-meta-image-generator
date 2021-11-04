package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/pug"
)

func main() {
	viewsPtr := flag.String("views", "./views", "location of views directory")
	publicPtr := flag.String("public", "./public", "location of public directory")

	flag.Parse()

	engine := pug.New(*viewsPtr, ".pug")
	engine.Debug(true)
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", *publicPtr)

	app.Get(":file.jpg", func(c *fiber.Ctx) error {
		img := screenshot(c.Params("file"), c.Query("background", "/grafana-dashboard.png"))
		c.Type("jpg")
		return c.Send([]byte(img))
	})

	app.Get("/render/:file.jpg", func(c *fiber.Ctx) error {

		message, err := url.QueryUnescape(c.Params("file"))
		if err != nil {
			log.Fatal(err)
		}

		return c.Render("index", fiber.Map{
			"message":    message,
			"background": c.Query("background", "/grafana-dashboard.png"),
		}, "layouts/main")
	})

	log.Fatal(app.Listen(":3000"))
}

func screenshot(file string, background string) []byte {
	name := "cache/" + file + ".jpg"

	if fileExists(name) {

		content, err := ioutil.ReadFile(name)

		if err != nil {
			log.Fatal(err)
		}

		return content

	} else {

		browser := rod.New().MustConnect()
		defer browser.MustClose()
		page := browser.MustPage("http://localhost:3000/render/" + file + ".jpg?background=" + background).MustWaitLoad()

		img, _ := page.Screenshot(true, &proto.PageCaptureScreenshot{
			Format:  proto.PageCaptureScreenshotFormatJpeg,
			Quality: 75,
			Clip: &proto.PageViewport{
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
