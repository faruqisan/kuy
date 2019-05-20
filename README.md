# kuy

match maker queueing system using golang

## example of usage

```go

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/faruqisan/kuy"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

const (
	numberOfPlayerInRoom = 10
)

func main() {

	k := kuy.New(kuy.Option{
		MaxItem:    numberOfPlayerInRoom,
		WaitPeriod: time.Second * 4, // wait for 1 sec
	})

	r := chi.NewRouter()
	// endpoint for find and join pool
	r.Get("/join", func(w http.ResponseWriter, r *http.Request) {
		// mock for user's id
		uID := uuid.New().ID()

		// join to kuy matchmaking by passing user's id
		// this function return channel to notify when pool is full
		c := k.Join(uID)

		// spawn go routine for listen the pool full channel
		go func() {
			select {
			case res := <-c:
				if res.IsFull {
					// when pool is full, do something .. in our case we just print it
					fmt.Println("pool", res.PoolID, " is ready with players : ", res.Items)
				}
				// listen for expired waiting time
				if res.TimeIsUp {
					fmt.Println("pool : ", res.PoolID, " can't find full member, and waiting time is over, members : ", res.Items)
				}
			}
		}()
		w.Write([]byte("success"))
	})

	r.Get("/total", func(w http.ResponseWriter, r *http.Request) {

		pn := k.GetNumberOfPools()
		strpn := fmt.Sprintf("{\"number_of_pools\" : %d}", pn)

		w.Write([]byte(strpn))
	})

	http.ListenAndServe(":3000", r)
}



```
