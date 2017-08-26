package runner

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/fatih/color"
)

var v string

func init() {
	v = "0.6.1"
}

/*
Logger - logger interface. Logger could be passed
*/
type Logger interface {
	Log(...interface{}) error
	Error(...interface{}) error
	Fatal(...interface{}) error
}

/*
Runner — is a Main service runner. Implements IRunner interface
Collects all services to handle.
Before Init every service initializes:
— Logger handler for Loggin/Monitoring/Alerting

Will:
— Initialize (Init)
— Run (Run)
*/
type Runner struct {
	Services []interface{}
	Logger   Logger
}

/*
IRunner interface that implements composite Interface for running services
*/
type IRunner interface {
	Exit() error
	Run()
}

/*
INamer - Namer interface
*/
type INamer interface {
	Name() string
}

/*
Initiator - Init interface
*/
type Initiator interface {
	Init()
}

var cyan *color.Color

func init() {
	cyan = color.New(color.FgCyan).Add(color.Underline)
}

/*
New — service constructor
*/
func New() Runner {
	runner := Runner{}
	runner.Intro()
	return runner
}

/*
Intro — starting app
*/
func (runner *Runner) Intro() {
	fmt.Println("\n=========================================================")
	fmt.Println("\n\t ", color.HiCyanString("Runner service"), " ", v)
	fmt.Println("\n=========================================================")
}

/*
Name - getting service name
*/
func (runner *Runner) Name(service interface{}) string {
	if s, ok := service.(INamer); ok {
		return s.Name()
	}
	return reflect.TypeOf(service).String()
}

/*
Add — adding service to map
*/
func (runner *Runner) Add(service IRunner) {
	runner.Log("Adding service " + color.HiMagentaString(runner.Name(service)))
	runner.Services = append(runner.Services, service)
}

/*
InitServices — Initializing serices
*/
func (runner *Runner) InitServices() {

	for _, v := range runner.Services {
		if s, ok := v.(Initiator); ok {
			s.Init()
		}
		runner.Log("Service " + color.CyanString(runner.Name(v)) + " initialized")
	}
}

/*
Run — Launching added services
*/
func (runner *Runner) Run() {
	runner.InitServices()
	runner.Log("Starting to run services")
	for _, v := range runner.Services {
		runner.Log("Service " + color.HiYellowString(runner.Name(v)) + " launching...")
		if s, ok := v.(IRunner); ok {
			go s.Run()
		}
		runner.Log("Service " + color.HiGreenString(runner.Name(v)) + " launched")
	}
	fmt.Println(".........................................................")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-interrupt:
		runner.Exit()
		runner.Log(color.HiGreenString("Exiting") + " OK")
		os.Exit(0)
	}
	time.Sleep(time.Second)

}

/*
Exit — calling exit methods for all services for graceful exit
*/
func (runner *Runner) Exit() {
	for _, v := range runner.Services {
		if s, ok := v.(IRunner); ok {
			s.Exit()
		}
	}
}

/*
Log — logging output
*/
func (runner *Runner) Log(message string) {
	fmt.Println(color.HiYellowString("Runner: "), message)
}
