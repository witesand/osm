package main


var (
	osmControllerName   string
	osmInjectorName   string
)

func wsInjectorInit() {
	flags.StringVar(&osmControllerName, "osm-controller-name", "osm-controller", "Service name of osm-controller.")
	flags.StringVar(&osmInjectorName, "osm-injector-name", "osm-injector", "Service name of osm-injector.")
}
