package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("failed to ping redis: %s", err)
	}

	httpServ := http.Server{
		Addr: ":9000",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count, err := rdb.Incr(r.Context(), "visit-count").Result()
			if err != nil {
				count = -1
				log.Printf("Error querying redis: %s:", err)
			}
			resbody := fmt.Sprintf(`
This server has been visited %d times.

%s`, count, r.URL.String())

			w.WriteHeader(200)
			_, _ = w.Write([]byte(resbody))
		}),
	}

	err = httpServ.ListenAndServe()
	if err != nil {
		log.Fatalf("listening failed: %s", err)
	}
}
