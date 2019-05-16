# kuy
[![Documentation](https://godoc.org/github.com/faruqisan/kuy?status.svg)](https://godoc.org/github.com/faruqisan/kuy)
[![Go Report Card](https://goreportcard.com/badge/github.com/faruqisan/kuy)](https://goreportcard.com/report/github.com/faruqisan/kuy)

match maker queueing system using golang

## example of usage

for example we create a system for matchmaking in dota2 where's the game required 10 player

```go

package main

import (
	"fmt"
	"net/http"

	"github.com/faruqisan/kuy"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

const (
	numberOfPlayerInRoom = 10
)

func main() {

	k := kuy.New(numberOfPlayerInRoom)

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
			for res := range c {
				if res.IsFull {
					// when pool is full, do something .. in our case we just print it
					fmt.Println("pool is ready with players : ", res.Items)
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

test the API using [hey](https://github.com/rakyll/hey) to simulate 100 user do the join request simultaneously

```bash 
 hey -n=100 -c=100 -m=GET http://localhost:3000/join 
```

on binary's console we should see something like this

```
PS C:\Users\faruqisan\go\src\github.com\faruqisan\testkuy> go run .\main.go
pool is ready with players :  [143307675 96363517 2719185236 1212008915 4066937870 85134676 651989607 2481560572 3491610833 2539808567]
pool is ready with players :  [4251982938 375222414 2260111718 1543982294 3779365849 4220720617 4260947384 2596047906 391349991 227458415]
pool is ready with players :  [278406534 292084715 3672121491 3855961325 3996627761 1961200593 1952689873 35570535 1774843344 2292847726]
pool is ready with players :  [2445225823 1075815861 2639597727 635431014 361501496 4058828195 1771799564 1043404946 2471276338 1313410725]
pool is ready with players :  [76016995 1708349512 1658771148 1703825659 3263739427 1015580438 1568046919 2844632825 4062717031 514537656]
pool is ready with players :  [3981065987 838893337 2602906694 1937865972 3495382920 3117318946 443082388 355671035 1598443688 750210196]
pool is ready with players :  [782453394 2171594867 1557878423 98355756 3946009775 369563420 1973375133 3682175270 798391453 395243692]
pool is ready with players :  [815681643 725905823 2790529077 3239416422 3381845128 3290407237 895526513 3516610004 1159466696 3258242278]
pool is ready with players :  [1091183992 1644771226 3923443848 3297405419 752093695 1302577618 4139572002 2691586518 2716704208 2737147699]
pool is ready with players :  [429182431 30544284 2964614348 583734382 3981060614 998643018 3453168831 625346978 1325249995 388061744]

```

and `/total` endpoint should show response like this

```json
{
	"number_of_pools": 10
}
```
