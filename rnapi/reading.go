package main

import (
    "net/http"
    "strconv"
    "github.com/julienschmidt/httprouter"
)


/**
    API URL mapping

    GET /
    GET /articles?page=1&ipp=10     Obtain the articles range from [page * ipp, page * ipp + ipp),
                                    latest added article show first, default: page == 0, ipp == 15
 */


func latestArticles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil {
        page = 0
    }

    ipp, err := strconv.Atoi(r.URL.Query().Get("ipp"))
    if err != nil {
        ipp = 10
    }

    db := GetDB(r)
    jas, err := daoLatestArticles(db, page * ipp, ipp, true)
    if err != nil {
        respondErr(w, r, http.StatusInternalServerError, "Database error: failed to get latest articles")
    } else {
        respond(w, r, http.StatusOK, jas, false)
    }
}


/*
    GET /article/[aid]             Get the article with id == aid
 */

func readArticle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    link := ps.ByName("permalink")
    jas, err := daoGetArticle(GetDB(r), link)
    if err != nil {
        respondErr(w, r, http.StatusInternalServerError, "Database error: failed to get article")
    } else {
        respond(w, r, http.StatusOK, jas, true)
    }
}







