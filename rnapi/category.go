package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

/*
   GET /category?v=文学>小说     Get the articles include 中 and char_level == 2
*/
func searchArticlesByCategory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cate := r.URL.Query().Get("v")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 0
	}
	num, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err != nil {
		num = 10
	}
	tracer.Trace("Got page=", page, ", and num=", num)
	jas, err := daoSearchArticlesByCategory(GetDB(r), cate, page*num, num)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, "Database error: failed to get latest articles")
	} else {
		respond(w, r, http.StatusOK, jas, true)
	}
}
