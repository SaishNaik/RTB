package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"ssp/dicontainer"
)

//func GetAdMarkup(rw http.ResponseWriter, r *http.Request) {
//	markup := `<img src="https://upload.wikimedia.org/wikipedia/commons/2/2a/New_Logo_AD.jpg">`
//	rw.Write([]byte(markup))
//}

type IRouter interface {
	InitRoutes(container dicontainer.IDiContainer) // todo check
	GetMux() *gin.Engine
}

type Router struct {
	engine *gin.Engine
}

func NewRouter(ginMode string) IRouter {
	gin.SetMode(ginMode)
	router := &Router{
		engine: gin.Default(),
	}
	//router.engine.Use(SetJSON)
	router.engine.Use(gin.Recovery()) //todo put explanation
	return router
}

func (r *Router) InitRoutes(container dicontainer.IDiContainer) {
	//todo set up middleware
	di := container.GetDiContainer()

	//todo add origins from config(s *SSPSvcHandler)
	r.engine.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	r.engine.OPTIONS("/ad", di.SSPController.AdRequest)
	r.engine.POST("/ad", di.SSPController.AdRequest)
	r.engine.GET("/imp/:bidReqId", di.SSPController.Impression)
}

func (r *Router) GetMux() *gin.Engine {
	return r.engine
}
