package main

import (
	"github.com/joho/godotenv"
	"github.com/lucianboboc/goBackendEngineering/internal/auth"
	"github.com/lucianboboc/goBackendEngineering/internal/db"
	"github.com/lucianboboc/goBackendEngineering/internal/env"
	"github.com/lucianboboc/goBackendEngineering/internal/mailer"
	"github.com/lucianboboc/goBackendEngineering/internal/ratelimiter"
	"github.com/lucianboboc/goBackendEngineering/internal/store"
	"github.com/lucianboboc/goBackendEngineering/internal/store/cache"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"os"
	"time"
)

const version = "0.0.1"

type Post struct {
	Title string `json:"title" validate:"required,max=100"`
}

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gophers.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @SecurityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
		os.Exit(1)
	}

	cfg := config{
		addr:        env.GetString("ADDR", "8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "localhost:5173"),
		db: dbConfig{
			dsn:          env.GetString("DB_DSN", "postgres://postgres:postgres@localhost/postgres?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConn:  env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pass:    env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3,
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			maiLTrap: mailTrapConfig{
				apiKey: env.GetString("MAILTRAP_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", ""),
				exp:    time.Hour * 24 * 3,
				iss:    "gopherSocial",
			},
		},
		ratelimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATE_LIMITER_REQUESTS_COUNT", 20),
			TimeFrame:            time.Second * 5,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true),
		},
	}

	// Database
	db, err := db.New(
		cfg.db.dsn,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConn,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()
	logger.Info("database connection established...")

	// Cache
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pass, cfg.redisCfg.db)
		logger.Info("redis cache connection established")
	}

	// Rate limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.ratelimiter.RequestsPerTimeFrame,
		cfg.ratelimiter.TimeFrame,
	)

	storage := store.NewPostgresStorage(db)
	cacheStore := cache.NewRedisStorage(rdb)

	sendGridMailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	nwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	app := &application{
		config:        cfg,
		store:         storage,
		cacheStorage:  cacheStore,
		logger:        logger,
		mailer:        sendGridMailer,
		authenticator: nwtAuthenticator,
		rateLimiter:   rateLimiter,
	}

	mux := app.mount()
	if err := app.run(mux); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
