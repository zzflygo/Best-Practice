package week3

import (
	"io"
	"net/http"
)

func getHandler()

func main() {
	var g errgroup.Group

	g.Go(func() error {
		helloHandler := func(w http.ResponseWriter, req *http.Request) { io.WriteString(w, "Hello, world!\n") }
		if err := http.ListenAndServe(":8080", helloHandler); err != nil {
		}
	}

	return nil
	)
}
