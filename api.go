package main

import (
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
)

func listener() error {

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	p := ginprom.New(
		ginprom.Engine(r),
	)
	r.Use(p.Instrument())

	return r.Run(":9876")
}
