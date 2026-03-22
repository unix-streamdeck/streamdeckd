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
	streamdeckd "github.com/unix-streamdeck/streamdeckd/streamdeckd"
	"github.com/unix-streamdeck/streamdeckd/streamdeckd/examples"
)

var isRunning = true

func main() {
	log.Default().SetFlags(log.Lshortfile | log.Ltime)
	log.Default().SetPrefix("(global) ")
	checkDuplicateStreamdeckdInstance()
	configPtr := flag.String("config", "", "Path to config file")
	flag.Parse()
	streamdeckd.SetConfigPath(*configPtr)
	cleanupHook()
	go streamdeckd.InitDBUS()
	go streamdeckd.UpdateApplication()
	go streamdeckd.EnableVirtualKeyboard()
	examples.RegisterBaseModules()
	streamdeckd.LoadConfig()
	streamdeckd.Devs = make(map[string]streamdeckd.IVirtualDev)
	screensaverDbus, err := streamdeckd.ConnectScreensaver()
	if err != nil {
		log.Println(err)
	} else {
		go screensaverDbus.RegisterScreensaverActiveListener()
	}
	attemptConnection()
}

func checkDuplicateStreamdeckdInstance() {
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

func attemptConnection() {
	for isRunning {
		streamdeckd.OpenDevice()
		time.Sleep(1 * time.Second)
	}
}

func cleanupHook() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGSTOP, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT)
	go func() {
		<-sigs
		shutdown()
	}()
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
