package main

import (
    "gopkg.in/mgo.v2"
    "testing"
    "fmt"
    "gopkg.in/mgo.v2/bson"
)

func TestMaxSno(t *testing.T) {
    dbconn, _ := mgo.Dial("localhost")
    defer dbconn.Close()
    maxsno := daoMaxSno(dbconn.DB("readcn_dev").C("articles"))
    fmt.Println("Max SNO ==", maxsno)
}

func TestGetArticle(t *testing.T) {
    dbconn, _ := mgo.Dial("localhost")
    defer dbconn.Close()
    db := dbconn.DB("readcn_dev")
    id := "55e660525485092fe0b23a73"
    fmt.Println("Lets's get the article", id)
    jas, err := daoGetArticle(db, id)

    if err != nil {
        t.Error("Failed get article:", err)
    }

    jid := jas["_id"]

    hid := bson.ObjectIdHex(id).Hex()
    if hid != id {
        t.Error("Failed to get the Hex string of an object ID:", hid)
    } else {
        fmt.Println("Restore Hex id from ObjectID:", hid)
    }

    if sid, ok := jid.(bson.ObjectId); !ok {
        t.Error("Failed to cast bson to object")
    } else {
        fmt.Println("Cast bson to object:", sid)
    }

}