package main

import (
	"os"

	"github.com/thofisch/ssm2k8s/internal/config"
	"github.com/thofisch/ssm2k8s/internal/logging"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app           = kingpin.New("mystico", "A command-line secret manager")
	debug         = app.Flag("debug", "Enable debug mode.").Bool()
	putCmd        = app.Command("put", "Create/update a secret.")
	putOptions    = NewPutCommand(putCmd)
	listCmd       = app.Command("list", "List secrets")
	deleteCmd     = app.Command("delete", "Delete secrets")
	deleteOptions = NewDeleteCommand(deleteCmd)
)

func main() {
	logger := logging.NewConsoleLogger()

	app.Version(config.Version + " (" +config.Commit + " " + config.BuildDate + ")")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case putCmd.FullCommand():
		ExecutePut(logger, putOptions)

	case listCmd.FullCommand():
		ExecuteList(logger)

	case deleteCmd.FullCommand():
		ExecuteDelete(logger, deleteOptions)

	default:
		kingpin.Usage()
		os.Exit(1)
	}

}
