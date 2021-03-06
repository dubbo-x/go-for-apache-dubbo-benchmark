package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

import (
	"github.com/AlexStocks/goext/time"
	"github.com/dubbogo/hessian2"
)

import (
	"github.com/dubbo/go-for-apache-dubbo/common/utils"
	"github.com/dubbo/go-for-apache-dubbo/config"

	"github.com/dubbo/go-for-apache-dubbo/common/logger"
	_ "github.com/dubbo/go-for-apache-dubbo/protocol/dubbo"
	_ "github.com/dubbo/go-for-apache-dubbo/registry/protocol"

	_ "github.com/dubbo/go-for-apache-dubbo/common/proxy/proxy_factory"
	_ "github.com/dubbo/go-for-apache-dubbo/filter/impl"

	_ "github.com/dubbo/go-for-apache-dubbo/cluster/cluster_impl"
	_ "github.com/dubbo/go-for-apache-dubbo/cluster/loadbalance"

	_ "github.com/dubbo/go-for-apache-dubbo/registry/zookeeper"
)

var (
	survivalTimeout = int(3e9)
)

func main() {

	// ------for hessian2------
	hessian.RegisterJavaEnum(Gender(MAN))
	hessian.RegisterJavaEnum(Gender(WOMAN))
	hessian.RegisterPOJO(&User{})
	// ------------

	_, proMap := config.Load()
	if proMap == nil {
		panic("proMap is nil")
	}

	initProfiling()

	initSignal()
}

func initProfiling() {

	ip, err := utils.GetLocalIP()
	if err != nil {
		panic("cat not get local ip!")
	}
	fmt.Println(ip + ":7070")
	go func() {
		logger.Info(http.ListenAndServe(ip+":7070", nil))
	}()
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		logger.Info("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
			// reload()
		default:
			go gxtime.Future(survivalTimeout, func() {
				logger.Warn("app exit now by force...")
				os.Exit(1)
			})

			// 要么fastFailTimeout时间内执行完毕下面的逻辑然后程序退出，要么执行上面的超时函数程序强行退出
			fmt.Println("provider app exit now...")
			return
		}
	}
}
