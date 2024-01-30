package config

import (
	"log"
	"os"
	"strconv"

	godotenv "github.com/joho/godotenv"
)

type (
	Config struct {
		App   App
		Db    Db
		Jwt   Jwt
		Grpc  Grpc
		Redis Redis
	}

	App struct {
		Name  string
		Url   string
		Stage string
	}

	Db struct {
		Url string
	}

	Redis struct {
		Addr     string
		Password string
		DB       int
	}

	Jwt struct {
		AccessSecretKey  string
		RefreshSecretKey string
		ApiSecretKey     string
		AccessDuration   int64
		RefreshDuration  int64
		ApiDuration      int64
		PrivateKeyPem    string
		PublicKeyPem     string
	}

	Grpc struct {
		UserUrl    string
		ProjectUrl string
		FavUrl     string
		CommentUrl string
	}
)

func LoadConfig(path string) Config {
	if err := godotenv.Load(path); err != nil {
		log.Fatalf("Error loading .env file : %s", err.Error())
	}

	// redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	// if err != nil {
	// 	log.Fatalf("Error can't convert REDIS_DB to int")
	// }
	return Config{
		App: App{
			Name:  os.Getenv("APP_NAME"),
			Url:   os.Getenv("APP_URL"),
			Stage: os.Getenv("APP_STAGE"),
		},
		Db: Db{
			Url: os.Getenv("DB_URL"),
		},
		Jwt: Jwt{
			AccessSecretKey:  os.Getenv("JWT_ACCESS_SECRET_KEY"),
			RefreshSecretKey: os.Getenv("JWT_REFRESH_SECRET_KEY"),
			ApiSecretKey:     os.Getenv("JWT_API_SECRET_KEY"),
			PrivateKeyPem:    os.Getenv("PRIVATE_KEY_PEM"),
			PublicKeyPem:     os.Getenv("PUBLIC_KEY_PEM"),
			AccessDuration: func() int64 {
				result, err := strconv.ParseInt(os.Getenv("JWT_ACCESS_DURATION"), 10, 64)
				if err != nil {
					log.Fatal("Error loading access duration file")
				}
				return result
			}(),

			RefreshDuration: func() int64 {
				result, err := strconv.ParseInt(os.Getenv("JWT_REFRESH_DURATION"), 10, 64)
				if err != nil {
					log.Fatal("Error loading refresh duration file")
				}
				return result
			}(),
		},
		Grpc: Grpc{
			UserUrl:    os.Getenv("GRPC_USERS_URL"),
			ProjectUrl: os.Getenv("GRPC_PROJECT_URL"),
			FavUrl:     os.Getenv("GRPC_FAV_URL"),
			CommentUrl: os.Getenv("GRPC_COM_URL"),
		},

		// Redis: Redis{
		// 	Addr:     os.Getenv("REDIS_ADDR"),
		// 	DB:       redisDb,
		// 	Password: os.Getenv("REDIS_PASSWORD"),
		// },
	}

}
