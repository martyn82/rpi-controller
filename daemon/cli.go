package daemon

import "flag"

type Arguments struct {
    ConfigFile string
}

var configFile = flag.String("c", "controllerd.conf.json", "Specify a configuration file to load.")

/* Parse command line arguments */
func ParseArguments() Arguments {
    flag.Parse()

    args := Arguments{}
    args.ConfigFile = *configFile

    return args
}
