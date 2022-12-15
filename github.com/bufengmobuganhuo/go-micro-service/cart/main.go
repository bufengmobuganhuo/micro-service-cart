package main

import (
	"github.com/bufengmobuganhuo/go-micro-service/cart/domain/repository"
	service2 "github.com/bufengmobuganhuo/go-micro-service/cart/domain/service"
	"github.com/bufengmobuganhuo/go-micro-service/cart/handler"
	cart "github.com/bufengmobuganhuo/go-micro-service/cart/proto/cart"
	common "github.com/bufengmobuganhuo/micro-service-common"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	consul2 "github.com/micro/go-plugins/registry/consul/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

const QPS = 100

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		log.Fatal(err)
	}
	// 注册中心
	consul := consul2.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})
	t, io, err := common.NewTracer("go.micro.service.cart", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	mysqlConfig := common.GetMysqlFromConsul(consulConfig, "mysql")
	db, err := mysqlConfig.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// 初始化数据库表
	//repo := repository.NewCartRepository(db)
	//repo.InitTable()

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.cart"),
		micro.Version("latest"),
		micro.Address("localhost:8084"),
		micro.Registry(consul),
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
	)

	// Initialise service
	service.Init()

	cartDataService := service2.NewCartDataService(repository.NewCartRepository(db))

	// Register Handler
	cart.RegisterCartHandler(service.Server(), &handler.Cart{cartDataService})

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
