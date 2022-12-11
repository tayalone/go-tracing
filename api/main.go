package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tayalone/go-trancing/api/trancer"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	// / setup Otel Trancer ---
	tp, err := trancer.JaegertracerProvider(os.Getenv("JEAGER_ENDPOINT"), os.Getenv("SERVICE_NAME"), os.Getenv("ENVIROMENT"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Otel Trancer Loonking Good !!!")
	// / -----------------------

	opts := options.Client()
	opts.ApplyURI(os.Getenv("ME_CONFIG_MONGODB_URL"))
	_, cErr := mongo.Connect(context.Background(), opts)

	fmt.Println(cErr.Error())

	// _, mgErr := mongodb.Connect()
	// if mgErr != nil {
	// 	log.Fatal(mgErr)
	// }
	// log.Println("Mongodb Loonking Good !!!")
	// database := client.Database("go-tracing")

	// / ----------------------
	r := gin.Default()
	r.Use(otelgin.Middleware(os.Getenv("SERVICE_NAME"))) // <- add Otel Middleware

	r.GET("/ping", func(c *gin.Context) {
		time.Sleep(500 * time.Millisecond) /* delay 0.5 secs */
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/test-trancing-1", func(ctx *gin.Context) {
		time.Sleep(50 * time.Millisecond)

		tr := tp.Tracer("/testxtrancing-1")

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

	r.GET("/test-trancing-2", func(ctx *gin.Context) {
		time.Sleep(50 * time.Millisecond)

		tr := tp.Tracer("/testxtrancing-1")

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

	// r.GET("/podcast", func(ctx *gin.Context) {
	// 	var p mongodb.Podcast

	// 	filter := bson.M{}
	// 	err := database.Collection("podcasts").FindOne(context.Background(), filter).Decode(&p)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(p)
	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"message": "ok",
	// 		"podcast":
	// 	})

	// })

	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
