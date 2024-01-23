!/bin/bash
# chmod +x generate-base.sh
# ./generate-base.sh

# ---------------- CONFIG ---------------- 

echo "Creating config file . . ."
mkdir "config"
cat <<EOL > "config/config.go"
package config

import (
	"log"
	"os"

	godotenv "github.com/joho/godotenv"
)


type (
	Config struct {
		App App
		// Db  Db
		// Jwt Jwt
	}

	App struct {
		Name  string
		Url   string
		Stage string
	}

)

func LoadConfig(path string) Config {
	if err := godotenv.Load(path); err != nil {
		log.Fatal("Error loading .env file : %s", err.Error())
	}

	return Config{
		App: App{
			Name:  os.Getenv("APP_NAME"),
			Url:   os.Getenv("APP_URL"),
			Stage: os.Getenv("APP_STAGE"),
		},
	}

}

EOL

# ---------------- CONFIG ---------------- 


# ---------------- ENV ---------------- 
echo "Creating ENV file . . ."
mkdir "env"
cat <<EOL > "env/.env"
APP_STAGE=dev
APP_NAME=${FILENAME}
APP_URL=:5000
EOL
# ---------------- ENV ---------------- 

# ---------------- SERVER ----------------
echo "Creating SERVER file . . ."

mkdir "server"
cat <<EOL > "server/server.go"
package server

type (
	server struct {
		app    any
		db     any
		cfg    *config.Config
	}
)


func (s *server) gracefulShutdown(pctx context.Context, quit <-chan os.Signal) {

	log.Printf("Starting service: %s", s.cfg.App.Name)

	<-quit

	log.Printf("Shutting down service: %s", s.cfg.App.Name)

	// depend on which library you use to shutdown the app in this case its fiber
	// if err := s.app.Shutdown(); err != nil {
	// 	log.Fatalf("Error: %v", err)
	// }

}

func (s *server) httpListening() {
	// base on library in this case it's fiber
	// if err := s.app.Listen(":5000"); err != nil && err != http.ErrServerClosed {
	// 	log.Fatalf("Error: %v", err)
	// }
}

func Start(pctx context.Context, cfg *config.Config, db any) {
	s := &server{
		db:  db,
		cfg: cfg,
		// app: fiber.New(fiber.Config{
		// 	AppName:      "testing",
		//	BodyLimit:    10 * 1024 * 1024,
		//	ReadTimeout:  10 * time.Second,
		//	WriteTimeout: 20 * time.Second,
		//	JSONEncoder:  json.Marshal,
		//	JSONDecoder:  json.Unmarshal,
		//}),
		app: nil,
	}

	// Body Limit
	// app.Settings.MaxRequestBodySize = 10 * 1024 * 1024 // 10 MB

	// Call the server service here


	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go s.gracefulShutdown(pctx, quit)

	// Listening
	s.httpListening()

}
EOL

# ---------------- SERVER ----------------