package main

import (
	"os"

	"github.com/thofisch/ssm2k8s/internal/config"
	"github.com/thofisch/ssm2k8s/internal/logging"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app           = kingpin.New(os.Args[0], "A command-line secret manager")
	globalDebug   = app.Flag("debug", "Enable debug mode.").Envar("DEBUG").Bool()
	globalRegion  = app.Flag("region", "AWS region").Envar("AWS_DEFAULT_REGION").String()
	putCmd        = app.Command("put", "Create/update a secret.")
	putOptions    = NewPutCommand(putCmd)
	listCmd       = app.Command("list", "List secrets")
	listOptions   = NewListCommand(listCmd)
	deleteCmd     = app.Command("delete", "Delete secrets")
	deleteOptions = NewDeleteCommand(deleteCmd)
)

func main() {
	app.Version(config.VersionString)

	command := kingpin.MustParse(app.Parse(os.Args[1:]))
	logger := logging.NewConsoleLogger(*globalDebug)

	//f, err := os.OpenFile("./log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
	//
	//_, err = fmt.Fprintf(f, "%#v\n", os.Args)
	//if err != nil {
	//	panic(err)
	//}

	switch command {
	case putCmd.FullCommand():
		ExecutePut(logger, putOptions)

	case listCmd.FullCommand():
		ExecuteList(logger, listOptions)

	case deleteCmd.FullCommand():
		ExecuteDelete(logger, deleteOptions)

	default:
		kingpin.Usage()
		os.Exit(1)
	}
}
