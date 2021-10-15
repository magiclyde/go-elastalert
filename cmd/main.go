/**
 * Created by GoLand.
 * @author: clyde
 * @date: 2021/10/12 下午4:21
 * @note:
 */

package main

import (
	"context"
	. "github.com/magiclyde/go-elastalert"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[elastalert] ")

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	NewElasticAlerter(NewConfig()).Run(ctx)
}
