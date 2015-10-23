package main

import (
    "net/http"
    "io/ioutil"
    "github.com/julienschmidt/httprouter"
)

/*
    POST /article/aid             Update the article with id == aid
 */

func editArticle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    permalink := ps.ByName("permalink")
    jso, _ := ioutil.ReadAll(r.Body)
    if err := daoUpdateArticle(GetDB(r), permalink, jso); err != nil {
        respondErr(w, r, http.StatusInternalServerError, "Database error: failed to update article")
        return
    }
    w.Header().Set("Location", "article/" + permalink)
    respond(w, r, http.StatusOK, nil, false)
}

func createArticle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    buf, _ := ioutil.ReadAll(r.Body)
    tracer.Trace(string(buf))
    err := daoNewArticle(GetDB(r), buf)
    if err != nil {
        respondErr(w, r, http.StatusInternalServerError, "Database error: failed to create article")
    } else {
        respond(w, r, http.StatusOK, "", false)
    }
}

func createNewId(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    am := daoNewId(GetDB(r))
    respond(w, r, http.StatusOK, am, false)
}

func deleteArticle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    permalink := ps.ByName("permalink")
    daoDeleteArticle(GetDB(r), permalink)
}


