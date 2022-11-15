package main

import (
	"go-gin-postgres-reactjs/db"
	// "go-gin-postgres-reactjs/web"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

type App struct {
	d db.DB
}

func main() {
	d, err := sql.Open("postgres", dataSource())
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()
	app := App{
		d: db.NewDB(d),
	}
	technologies, err := app.d.GetTechnologies()
	// CORS is enabled only in prod profile
	prod := os.Getenv("profile") == "prod"
	r := gin.Default()
	if !prod {
		r.Use(cors.Default())
	}
	// app := web.NewApp(db.NewDB(d), prod)
	// err = app.Serve()
	// log.Println("Error", err)
	r.Use(static.Serve("/", static.LocalFile("/webapp", true)))
	r.GET("/api/technologies", func(c *gin.Context) {
		log.Println(technologies)
		c.JSON(http.StatusOK, technologies)
	})
	log.Println("Web server is available on port 8080")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func dataSource() string {
	host := "localhost"
	pass := "pass"
	if os.Getenv("profile") == "prod" {
		host = "db"
		pass = os.Getenv("db_pass")
	}
	return "postgresql://" + host + ":5432/goxygen" +
		"?user=goxygen&sslmode=disable&password=" + pass
}
