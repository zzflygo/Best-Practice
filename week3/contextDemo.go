package week3

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
)

func f2(ctx context.Context) error {
	for {
		fmt.Println("zzzzzz")
		select {
		case <-ctx.Done():
			return errors.New("f2 out")
		default:
		}
	}
}
func f1(ctx context.Context) error {
	for {
		fmt.Println("aaaaaa")
		select {
		case <-ctx.Done():
			return errors.New("f1 out")
		default:
		}
	}
}

func contex1() {
	ctx, cancel := context.WithCancel(context.Background())
	group, errctx := errgroup.WithContext(ctx)
	defer cancel()
	group.Go(func() error {
		return f1(errctx)
	})
	group.Go(func() error {
		return f2(errctx)
	})
	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}

}
