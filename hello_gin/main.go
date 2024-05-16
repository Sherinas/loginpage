package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func ClearCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("X-Accel-Expires", "0")
}

func main() {

	r := gin.Default()
	//	r.Use(gin.Logger())

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.Static("/static", "./static")
	r.LoadHTMLGlob("./template/**")

	r.GET("/", func(ctx *gin.Context) {
		ClearCache(ctx)
		ctx.HTML(http.StatusOK, "index.html", gin.H{})

	})

	r.POST("/", func(ctx *gin.Context) {
		ClearCache(ctx)
		session := sessions.Default(ctx)
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")

		if username == "admin" && password == "12345" {

			session.Set("userID", username)
			session.Save()
			ctx.Redirect(http.StatusSeeOther, "/dashboard")

			return

		}

		ctx.Redirect(http.StatusSeeOther, "/")

	})

	r.GET("/dashboard", func(c *gin.Context) {
		ClearCache(c)
		session := sessions.Default(c)

		userID := session.Get("userID")
		if userID == nil {
			c.Redirect(http.StatusSeeOther, "/")
			return
		}

		c.HTML(http.StatusOK, "dashboard.html", nil)
	})
	r.GET("/logout", func(c *gin.Context) {
		ClearCache(c)
		session := sessions.Default(c)

		session.Clear()
		session.Save()

		c.Redirect(http.StatusSeeOther, "/")
	})

	r.Run(":8000")

}
