package week3

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"os"
	"os/signal"
)

func HelloHandle(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello~")
}

//基于errgroup实现一个http server的启动和关闭 ，以及linux signal信号的注册和处理，要保证能够一个退出，全部注销退出。
func main() {

	http.HandleFunc("/", HelloHandle)
	ctx, cancel := context.WithCancel(context.Background())
	g, errctx := errgroup.WithContext(ctx)
	srv := &http.Server{Addr: ":8080"}
	g.Go(func() error {
		//起一个http服务
		fmt.Println("http server start...")
		return srv.ListenAndServe()
	})
	g.Go(func() error {
		//等待信号
		<-errctx.Done()
		fmt.Println("http server stop")
		//关闭http server
		return srv.Shutdown(errctx)
	})
	//注册linux signal信号
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch)
	//处理linux信号
	g.Go(func() error {
		for {
			select {
			case <-sigch:
				cancel()
			case <-errctx.Done():
				return errctx.Err()
			}
		}
	})

	//退出 http服务
	if err := g.Wait(); err != nil {
		fmt.Println("errgroup err:", err)
	}
	fmt.Println("all serve downr")
}
