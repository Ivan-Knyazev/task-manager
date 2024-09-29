package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized.POST("admin", func(c *gin.Context) {
		//user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			User  string `json:"user" binding:"required"`
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			user := json.User
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		}

		fmt.Println("user =", c.MustGet(gin.AuthUserKey).(string))
		fmt.Println("db =", db)
	})

	return r
}

func main() {
	r := setupRouter()
	//gin.SetMode(gin.ReleaseMode)
	// Listen and Server in 0.0.0.0:8080
	r.Run("0.0.0.0:8000")
}

/* example curl for /admin with basicauth header
   Zm9vOmJhcg== is base64("foo:bar")

	curl -X POST \
  	http://localhost:8000/admin \
  	-H 'authorization: Basic Zm9vOmJhcg==' \
  	-H 'content-type: application/json' \
  	-d '{"value":"bar"}'
*/
