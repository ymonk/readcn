package main_test

import (
    "fmt"
    "net/http"
    "testing"

    "github.com/ymonk/readcn/api"
)

func TestVars(t *testing.T) {
    fmt.Println("Begin var test")
    r1 := &http.Request{}
    r2 := &http.Request{}

    var1 := "var1"
    var2 := "var2"

    fmt.Println("Opening Vars")
    main.OpenVars(r1)
    main.OpenVars(r2)

    main.SetVar(r1, "key", var1)
    main.SetVar(r2, "key", var2)

    var t1 = main.GetVar(r1, "key")
    var t2 = main.GetVar(r2, "key")

    if t1 != var1 && t2 != var2 {
        t.Error("Failed to retreive the var")
    }

    main.CloseVars(r1)
    main.CloseVars(r2)
}





