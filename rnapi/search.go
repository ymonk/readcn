package main

import (
    "net/http"
    "strconv"
    "github.com/julienschmidt/httprouter"
)


/*
    GET /bychar?t=中&g=2&p=3&n=10     Get the articles include 中 and char_level == 2
 */
func searchArticlesByChar(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    ch := r.URL.Query().Get("t")
    grade, err := strconv.Atoi(r.URL.Query().Get("g"))
    if err != nil {
        grade = 1
    }
    page, err := strconv.Atoi(r.URL.Query().Get("p"))
    if err != nil {
        page = 0
    }
    num, err := strconv.Atoi(r.URL.Query().Get("n"))
    if err != nil {
        num = 10
    }

    jas, err := daoSearchArticlesByChar(GetDB(r), ch, grade, page * num, num)
    if err != nil {
        respondErr(w, r, http.StatusInternalServerError, "Database error: failed to get latest articles")
    } else {
        respond(w, r, http.StatusOK, jas, true)
    }
}

/*
    GET /byvocabulary?t=中国&g=2&p=3&n=10     Get the articles include 中 and char_level == 2
 */
func searchArticlesByVocabulary(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    word := r.URL.Query().Get("t")
    grade, err := strconv.Atoi(r.URL.Query().Get("g"))
    if err != nil {
        grade = 1
    }
    page, err := strconv.Atoi(r.URL.Query().Get("p"))
    if err != nil {
        page = 0
    }
    num, err := strconv.Atoi(r.URL.Query().Get("n"))
    if err != nil {
        num = 10
    }

    jas, err := daoSearchArticlesByVocabulary(GetDB(r), word, grade, page * num, num)
    if err != nil {
        respondErr(w, r, http.StatusInternalServerError, "Database error: failed to get latest articles")
    } else {
        respond(w, r, http.StatusOK, jas, true)
    }
}


