package main

import (
	"flag"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "unsafe"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/unix-streamdeck/streamdeckd/streamdeckd"
	"github.com/unix-streamdeck/streamdeckd/streamdeckd/examples"
)

var isRunning = true

func main() {
	initLogger()

	checkDuplicateInstance()

	configPtr := flag.String("config", "", "Path to config file")
	flag.Parse()
	streamdeckd.SetConfigPath(*configPtr)

	go listenForExitSignals()

	go streamdeckd.InitDBUS()
	go handleScreensaver()

	go streamdeckd.UpdateApplication()
	go streamdeckd.EnableVirtualKeyboard()

	examples.RegisterBaseModules()

	streamdeckd.LoadConfig()

	attemptConnection()
}

func initLogger() {
	log.Default().SetFlags(log.Lshortfile | log.Ltime)
	log.Default().SetPrefix("(global) ")
}

func checkDuplicateInstance() {
	processes, err := process.Processes()
	if err != nil {
		log.Println("Could not check for other instances of streamdeckd, assuming no others running")
	}
	for _, proc := range processes {
		name, err := proc.Name()
		if err == nil && name == "streamdeckd" && int(proc.Pid) != os.Getpid() {
			log.Fatalln("Another instance of streamdeckd is already running, exiting...")
		}
	}
}

func handleScreensaver() {
	screensaverDbus, err := streamdeckd.ConnectScreensaver()

	if err != nil {
		log.Println(err)
	} else {
		screensaverDbus.RegisterScreensaverActiveListener()
	}
}

func attemptConnection() {
	streamdeckd.Devs = make(map[string]streamdeckd.IVirtualDev)

	for isRunning {
		streamdeckd.OpenDevice()
		time.Sleep(1 * time.Second)
	}
}

func listenForExitSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGSTOP, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT)
	<-sigs
	shutdown()
}

func shutdown() {
	log.Println("Cleaning up")
	isRunning = false
	streamdeckd.UnmountHandlers()
	for s := range streamdeckd.Devs {
		streamdeckd.Devs[s].Close()
	}
	os.Exit(0)
}
