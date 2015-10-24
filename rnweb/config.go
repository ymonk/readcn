package main

import (
	"flag"
	"os"
)

// AppConfig defines
var AppConfig struct {
    Env                     string
    ApiHostAddr             string
    WebHost                 string
    Port                    string
    CookieSalt              string
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
	AppConfig.CookieSalt = "fff3a4f5-cd75-46ef-8762-564793d6a3d4"

    env := os.Getenv("READCN_ENV")
    switch env {
    case "development":
        AppConfig.Env = "development"
        AppConfig.ApiHostAddr = "http://localhost:5050"
        AppConfig.WebHost = "http://localhost"
        AppConfig.Port = ":" + *portPtr
        AppConfig.CookieSalt = "fff3a4f5-cd75-46ef-8762-564793d6a3d4"
    case "production":
        AppConfig.Env = "production"
        AppConfig.ApiHostAddr = "http://writeuptube.com:5050"
        AppConfig.WebHost = "http://writeuptube.com"
        AppConfig.Port = ":80"
        AppConfig.CookieSalt = "9Df3a4f5-cd75-46ef-8762-564793d6a3d4"
    }
}
