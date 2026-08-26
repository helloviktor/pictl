package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/pictl/pictl/internal/services"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: pictl <info|update|restart|shutdown>")
		fmt.Println("\nCommands:")
		fmt.Println("  info      print current system information as JSON")
		fmt.Println("  update    update installed system packages")
		fmt.Println("  restart   restart the device")
		fmt.Println("  shutdown  shut down the device")
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}

	service := services.NewSystemService()
	var err error

	switch flag.Arg(0) {
	case "info":
		var info interface{}
		info, err = service.GetSystemInfo()
		if err == nil {
			err = json.NewEncoder(log.Writer()).Encode(info)
		}
	case "update":
		err = service.UpdateSystem()
	case "restart":
		err = service.RestartSystem()
	case "shutdown":
		err = service.ShutdownSystem()
	default:
		flag.Usage()
		return
	}

	if err != nil {
		log.Fatal(err)
	}
}
