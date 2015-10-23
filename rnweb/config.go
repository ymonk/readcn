package main

import (
	"flag"
	"os"
)

// AppConfig defines
var AppConfig struct {
    Host                       string
    Port                       string
    CookieSalt                 string
}


// Return string value from environment variable, if not found,
// return the predefined string value
func getEnvOr(env, or string) string {
    if s := os.Getenv(env); len(s) > 0 {
        return s
    }
    return or
}

func init() {
	portPtr := flag.String("port", "5000", "port number")
	flag.Parse()
	AppConfig.Host = getEnvOr("READCN_HOST", "0.0.0.0")
	AppConfig.Port = ":" + *portPtr
	AppConfig.CookieSalt = "fff3a4f5-cd75-46ef-8762-564793d6a3d4"
}
