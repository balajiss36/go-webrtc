package server

import (
	"flag"
	"os"
	"time"

	"github.com/balajiss36/go-webrtc/internal/handlers"
	w "github.com/balajiss36/go-webrtc/pkg/webrtc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

var (
	addr = flag.String("addr", ":8080", os.Getenv("PORT"))
	cert = flag.String("cert", "", "TLS certificate")
	key  = flag.String("key", "", "TLS private key")
)

func Run() error {
	flag.Parse()
	// ...

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{Views: engine})
	app.Use(logger.New()) // middleware for fiber
	app.Use(cors.New())

	app.Get("/", handlers.Welcome)
	app.Get("/room/create", handlers.CreateRoom)
	app.Static("/", "./assets")

	w.Rooms = make(map[string]*w.Room)
	w.Streams = make(map[string]*w.Room)
	go dispatchKeyFrames()

	if *cert != "" {
		return app.ListenTLS(*addr, *cert, *key)
	}
	return app.Listen(*addr)
}

func dispatchKeyFrames() {
	for range time.NewTicker(time.Second * 3).C {
		for _, room := range w.Rooms {
			room.Peers.DispatchKeyFrame()
		}
	}
}