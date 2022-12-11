package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tayalone/go-trancing/api/db"
	"github.com/tayalone/go-trancing/api/trancer"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	// / setup Otel Trancer ---
	tp, err := trancer.JaegertracerProvider(os.Getenv("JEAGER_ENDPOINT"), os.Getenv("SERVICE_NAME"), os.Getenv("ENVIROMENT"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Otel Tracer Loonking Good !!!")
	// / -----------------------

	// // --------- db connection ---------------
	rdb, rdbErr := db.New()
	if rdbErr != nil {
		log.Fatal(err)
	}
	// // ---------------------------------------
	r := gin.Default()
	r.Use(otelgin.Middleware(os.Getenv("SERVICE_NAME"))) // <- add Otel Middleware

	r.GET("/ping", func(c *gin.Context) {
		time.Sleep(500 * time.Millisecond) /* delay 0.5 secs */
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/test-tracing-1", func(ctx *gin.Context) {
		time.Sleep(50 * time.Millisecond)

		tr := tp.Tracer("/testxtracing-1")

		// / do task no1
		_, span1 := tr.Start(ctx.Request.Context(), "task-1")
		time.Sleep(150 * time.Millisecond)
		span1.End()

		// / do task no2
		_, span2 := tr.Start(ctx.Request.Context(), "task-2")
		time.Sleep(150 * time.Millisecond)
		span2.End()

		time.Sleep(100 * time.Millisecond)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.GET("/todo/:id", func(ctx *gin.Context) {
		type getIDUri struct {
			ID uint `uri:"id" binding:"required"`
		}

		var gi getIDUri
		if err := ctx.ShouldBindUri(&gi); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		var td db.Todo
		r := rdb.First(&td, gi.ID)
		if r.RowsAffected != 1 {
			// return emptyBc, errors.New("Barcode Condition Not Found")
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "not found todo",
			})
			return
		}
		// return bc, nil
		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"todo":    td,
		})
	})

	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
