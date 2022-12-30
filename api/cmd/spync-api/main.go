package main

import (
	"api/internal/applemusic"
	"api/internal/configuration"
	"api/internal/ping"
	"api/internal/spotify"
	"api/internal/syncer"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"strings"
	"time"
)

func main() {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	gin.DefaultWriter = GinStartupLoggingWriter(sugar)
	router := gin.New()
	err := router.SetTrustedProxies([]string{})
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())
	router.Use(LogMiddleware(sugar))

	if err != nil {
		logger.Error(err.Error())
		return
	}

	router.GET("/ping", ping.Ping)

	router.GET("/settings", configuration.GetSettingsEndpoint)
	router.POST("/settings", configuration.SaveSettingsEndpoint)

	router.POST("/apple-music/save-auth", applemusic.SaveAuthEndpoint)
	router.GET("/apple-music/playlist/synced", applemusic.GetSpotifyPlaylistsEndpoint)
	router.GET("/apple-music/tracks/spotify-track/:trackId", applemusic.FindTrackEndpoint)

	router.GET("/spotify/auth", spotify.GetAuthURLEndpoint)
	router.GET("/spotify/save-auth", spotify.SaveAuthEndpoint)
	router.GET("/spotify/me", spotify.GetMeEndpoint)
	router.GET("/spotify/me/playlists", spotify.GetPlaylistsEndpoint)
	router.GET("/spotify/me/playlists/:playlistId/tracks", spotify.GetPlaylistTracksEndpoint)

	router.GET("/sync/status", syncer.StatusSocket)
	router.POST("/sync/playlist/:playlistId", syncer.SyncPlaylistEndpoint)
	router.POST("/sync/all", syncer.SyncPlaylistsEndpoint)

	err = router.Run()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func LogMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UnixNano() / int64(time.Millisecond)

		logger.Debugw("Incoming request",
			"method", c.Request.Method,
			"url", c.Request.URL,
		)

		c.Next()
		end := time.Now().UnixNano() / int64(time.Millisecond)
		diff := end - start

		logger.Debugw("Returning response",
			"status", c.Writer.Status(),
			"duration", fmt.Sprintf("%dms", diff),
			"method", c.Request.Method,
			"url", c.Request.URL,
		)
	}
}

type WriteFunc func([]byte) (int, error)

func (fn WriteFunc) Write(data []byte) (int, error) {
	return fn(data)
}

func GinStartupLoggingWriter(logger *zap.SugaredLogger) io.Writer {
	return WriteFunc(func(data []byte) (int, error) {
		message := strings.TrimSuffix(string(data), "\n")
		message = strings.TrimPrefix(message, "[GIN-debug]")
		message = strings.TrimSpace(message)
		logger.Debug(message)
		return 0, nil
	})
}
