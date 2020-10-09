package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"redisCluster/com"
)

var app = cli.NewApp()
var conf = &com.Conf{}

func init() {
	conf.Restore()
	com.ConfigAll(app)
}

func main() {

	log.Println(`
	 ____  ____  ____  ____  ___     ___  __    __  __  ___  ____  ____  ____    ___  ____   __    ____  ____         
	(  _ \( ___)(  _ \(_  _)/ __)   / __)(  )  (  )(  )/ __)(_  _)( ___)(  _ \  / __)(_  _) /__\  (  _ \(_  _)        
	 )   / )__)  )(_) )_)(_ \__ \  ( (__  )(__  )(__)( \__ \  )(   )__)  )   /  \__ \  )(  /(__)\  )   /  )(          
	(_)\_)(____)(____/(____)(___/   \___)(____)(______)(___/ (__) (____)(_)\_)  (___/ (__)(__)(__)(_)\_) (__) `)

	app.Run(os.Args)

}
