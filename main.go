package main

import (
	// "fmt"
	"net/http"
	"runtime"

	"api-gin/controllers"
	"api-gin/models"
	"github.com/gin-gonic/gin"
)

var db_redis = new(models.RedisPool) //nama-package.nama-type-data-file-package

//Enable CORS (Cross-Origin Resource Sharing) - using middleware
func CORS_Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

func main() {
	router := gin.Default()
	router.Use(CORS_Middleware())
	db_redis.Init()

	okz := router.Group("/okz")
	{
		content := new(controllers.ContentController)

		okz.POST("/create/", content.Create)                           //make sure func name in folder we want to imported is uppercase
		okz.GET("/headline/:type/:id/:limit/:start", content.Headline) //listContent
		okz.GET("/breaking/:type/:id/:limit/:start", content.Breaking) //listContent
		okz.GET("/read/:type/:id", content.Read)                       //detailContent
		okz.PUT("/:id", content.Update)
		okz.DELETE("/:id", content.Drop)

		okz.POST("/update/", content.Update)
		okz.POST("/drop/", content.Drop)
	}

	gui := router.Group("/gui")
	{
		gui.GET("/path/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "path.html", gin.H{"docNecessity": "path for accesing API"})
		})
		gui.GET("/action/", func(c *gin.Context) { c.HTML(http.StatusOK, "action_list.html", false) })
		gui.GET("/insert/", func(c *gin.Context) { c.HTML(http.StatusOK, "insert.html", false) })
		gui.GET("/update/", func(c *gin.Context) { c.HTML(http.StatusOK, "update.html", false) })
		gui.GET("/delete/", func(c *gin.Context) { c.HTML(http.StatusOK, "delete.html", false) })
		gui.GET("/list/", func(c *gin.Context) { c.HTML(http.StatusOK, "test.html", false) })
	}

	router.LoadHTMLGlob("./public/html/*")
	router.Static("/public", "./public") //serve static files
	router.Static("/assets", "./assets")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"ginBoilerplateVersion": "v0.03",
			"goVersion":             runtime.Version(),
		})
	})
	router.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})

	router.Run(":8080")
}
