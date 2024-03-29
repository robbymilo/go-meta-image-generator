package main

import (
	"bytes"
	"crypto/sha256"
	"embed"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/urfave/cli/v2"
)

//go:embed templates
var templatesDir embed.FS

//go:embed public
var publicDir embed.FS

func main() {

	app := &cli.App{
		Name:  "go-meta-image-generator",
		Usage: "An open graph meta image generator for social previews.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "cache",
				Usage: "location of cache directory",
				Value: "./cache",
			},
		},
		Action: func(cCtx *cli.Context) error {

			r := chi.NewRouter()
			r.Use(middleware.Logger)

			var publicFs = fs.FS(publicDir)
			publicContent, err := fs.Sub(publicFs, "public")
			if err != nil {
				log.Fatal(err)
			}
			r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.FS(publicContent))))

			r.Get("/render/{title}.jpg", func(w http.ResponseWriter, r *http.Request) {

				title := chi.URLParam(r, "title")

				t, err := template.ParseFS(templatesDir, "templates/index.html")
				if err != nil {
					fmt.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return

				}

				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				err = t.Execute(w, struct {
					Title string
				}{Title: title})

				if err != nil {
					fmt.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return

				}
			})

			r.Get("/{title}.jpg", func(w http.ResponseWriter, r *http.Request) {
				title := chi.URLParam(r, "title")
				cache := cCtx.String("cache")
				signature := r.URL.Query().Get("signature")

				if !IsSigned(signature, title) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				path, err := filepath.Abs(filepath.Join(cache, url.QueryEscape(title)+".jpg"))
				if err != nil {
					fmt.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if fileExists(path) {
					// serve a cached file
					file, err := os.ReadFile(path)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Println("error reading saved file:", err)
					}
					w.Header().Set("Content-Type", "image/jpeg")
					http.ServeContent(w, r, "thumbnail", time.Now(), bytes.NewReader(file))

				} else {
					// generate, save, and serve a new file
					file, _ := generataeScreenshot(title)

					err := os.WriteFile(path, file, 0644)
					if err != nil {
						fmt.Println("error caching generated file:", err)
					}
					w.Header().Set("Content-Type", "image/jpeg")
					http.ServeContent(w, r, "thumbnail", time.Now(), bytes.NewReader(file))
				}
			})

			r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("ok\n"))
			})

			fmt.Println("listening on :3000")

			err = http.ListenAndServe(":3000", r)
			if err != nil {
				log.Fatal(err)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func generataeScreenshot(title string) ([]byte, error) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	page := browser.MustPage("http://localhost:3000/render/" + title + ".jpg").MustWaitLoad()

	img, _ := page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format:  proto.PageCaptureScreenshotFormatJpeg,
		Quality: createInt(75),
		Clip: &proto.PageViewport{
			Width:  1200,
			Height: 700,
			Scale:  1,
		},
	})

	return img, nil
}

// fileExists checks if a generated image exists on disk.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func IsSigned(signature, title string) bool {

	s := os.Getenv("SIGNATURE")

	if len(s) > 0 {
		// get signature from query param
		sDec, _ := b64.StdEncoding.DecodeString(signature)

		// generate sha to compare
		input := []byte(title + s)
		sCompare := sha256.Sum256(input)

		// test if signatures match
		fmt.Println(hex.EncodeToString(sCompare[:]) == string(sDec))
		return hex.EncodeToString(sCompare[:]) == string(sDec)
	}

	return true

}

func createInt(x int) *int {
	return &x
}
