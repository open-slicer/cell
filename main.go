package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("Couldn't read config")
	}

	r := gin.Default()
	// Gin defaults to debug mode.
	if !viper.GetBool("debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	addr := viper.GetString("http.address")
	if certFile := viper.GetString("security.cert_file"); certFile != "" {
		// Let's assume key_file is present.
		keyFile := viper.GetString("security.key_file")
		log.Err(r.RunTLS(addr, certFile, keyFile)).Msg("Running HTTP server (tls)")
	}
	log.Err(r.Run(addr)).Msg("Running HTTP server")
}
