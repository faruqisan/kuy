# kuy

match maker queueing system using golang

## example of usage

```go

package main

import (
	"log"

	"github.com/faruqisan/kuy"
)

func main() {
    // create new engine with max item on param
	e := kuy.New(5)

    // simulate user joining
	for i := 0; i < 500; i++ {
		go func(i int) {
			c := e.Join(i)

			select {
			case respChan := <-c:
				if respChan.IsFull {
					log.Println("got full pool : ", respChan.Items)
				} else {
					return
				}
			default:
			}

		}(i)
	}

    // block the app
	for {
	}

}


```
