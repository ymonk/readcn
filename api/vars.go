package main

import (
    "net/http"
    "sync"
    "gopkg.in/mgo.v2"
)

var vars map[*http.Request]map[string]interface{}
var varsLock sync.RWMutex
var dbs map[*http.Request]*mgo.Database

func init() {
    dbs = make(map[*http.Request]*mgo.Database, 16)
}

func GetVar(r *http.Request, key string) interface{} {
    varsLock.RLock()
    value := vars[r][key]
    varsLock.RUnlock()
    return value
}

func SetVar(r *http.Request, key string, value interface{}) {
    varsLock.Lock()
    vars[r][key] = value
    varsLock.Unlock()
}

func OpenVars(r *http.Request) {
    varsLock.Lock()
    if vars == nil {
        vars = make(map[*http.Request]map[string]interface{}, 16)
    }
    vars[r] = make(map[string]interface{}, 5)
    varsLock.Unlock()
}

func CloseVars(r *http.Request) {
    varsLock.Lock()
    delete(vars, r)
    varsLock.Unlock()
}

func SetDB(r *http.Request, db *mgo.Database) {
    dbs[r] = db
}

func GetDB(r *http.Request) *mgo.Database {
    return dbs[r]
}
