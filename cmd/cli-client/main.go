package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/MaximBayurov/rate-limiter/internal/client"
	"github.com/MaximBayurov/rate-limiter/internal/configuration"
)

var (
	configFile string
	login      string
	ip         string
	listType   string
	overwrite  bool
)

func init() {
	flag.StringVar(&configFile, "config", "./config/cli/config.yaml", "Path to configuration file")
	flag.StringVar(&login, "login", "", "логин пользователя")
	flag.StringVar(&ip, "ip", "", "IP пользователя")
	flag.StringVar(&listType, "listType", "", "тип списка")
	flag.BoolVar(&overwrite, "overwrite", false, "перезаписывать IP")
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config, err := configuration.New(configFile)
	if err != nil {
		panic(err)
	}

	app := client.New(config.Client)

	switch flag.Arg(0) {
	case "add-ip":
		if len(ip) == 0 || len(listType) == 0 {
			fmt.Println("Не передан IP или тип списка")
			break
		}

		resp, _ := app.AddIP(ctx, ip, listType, overwrite)
		fmt.Println(resp.Message)
	case "delete-ip":
		if len(ip) == 0 || len(listType) == 0 {
			fmt.Println("Не передан IP или тип списка")
			break
		}

		resp, _ := app.DeleteIP(ctx, ip, listType)
		fmt.Println(resp.Message)
	case "clear-bucket":
		if len(ip) == 0 || len(login) == 0 {
			fmt.Println("Не передан IP или логин")
			break
		}

		resp, _ := app.ClearBucket(ctx, ip, login)
		fmt.Println(resp.Message)
	case "help":
		fmt.Println("Доступные команды: add-ip, delete-ip, clear-bucket")
	default:
		fmt.Println("Неподдерживаемая команда. Используйте help для получения информации о доступных командах")
	}
}
