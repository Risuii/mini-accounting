package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	CountingHandler "mini-accounting/internal/counting/delivery/presenter/http"

	Library "mini-accounting/library"
	// Middleware "mini-accounting/middlewares"
	UtilsPackage "mini-accounting/pkg/utils"
)

type Routes interface {
	Setup()
	GetEngine() *gin.Engine
}

type RoutesImpl struct {
	engine  *gin.Engine
	library Library.Library
	// middleware      Middleware.Middleware
	countingHandler CountingHandler.CountingHandler
}

func New(
	engine *gin.Engine,
	library Library.Library,
	// middleware Middleware.Middleware,
	countingHandler CountingHandler.CountingHandler,
) Routes {
	return &RoutesImpl{
		engine:  engine,
		library: library,
		// middleware:      middleware,
		countingHandler: countingHandler,
	}
}

func (r *RoutesImpl) Setup() {
	path := "Routes:Setup"
	defer UtilsPackage.CatchPanic(path, r.library)
	// SETUP CORS
	r.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, //http or https
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	// LOAD HTML FILE
	r.engine.LoadHTMLGlob("templates/*.html")
	// EMBED ROUTES
	r.SetIndexRoute()
	r.SetCountingRoute()
}

func (r *RoutesImpl) GetEngine() *gin.Engine {
	return r.engine
}

func (r *RoutesImpl) SetIndexRoute() {
	r.engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
}

func (r *RoutesImpl) SetCountingRoute() {
	counting := r.engine.Group("/api/counting")

	counting.POST("/upload-csv", r.countingHandler.ReadCSV)
	counting.GET("/download-csv", r.countingHandler.DownloadCSV)
}
