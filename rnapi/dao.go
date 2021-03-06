package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func daoMaxSno(articles *mgo.Collection) int {
	result := bson.M{}
	articles.Find(bson.M{}).Select(bson.M{"_id": 0, "sno": 1}).
		Sort("-sno").Limit(1).Iter().Next(&result)
	ret, _ := result["sno"].(int)
	return ret
}

func daoLatestArticles(db *mgo.Database, start, num int, public_only bool) ([]bson.M, error) {
	as := db.C("articles")
	jas := make([]bson.M, num)
	firstSno := daoMaxSno(as) - start

	// visibility == 2: public available
	// visibility == 1: draft
	// visibility == 0: in trash
	query := bson.M{"sno": bson.M{"$lte": firstSno}, "visibility": bson.M{"$gte": 1}}
	if public_only {
		query = bson.M{"sno": bson.M{"$lte": firstSno}, "visibility": bson.M{"$gte": 2}}
	}
	err := as.Find(query).Select(bson.M{"title": 1, "preview": 1, "permalink": 1,
		"source": 1, "publishedAt": 1, "char_level": 1,
		"vocabulary_level": 1, "grammar_level": 1,
		"num_comment": 1, "tags": 1}).
		Sort("-sno").Limit(num).All(&jas)
	if err != nil {
		return nil, err
	}
	return jas, nil
}

func daoGetArticle(db *mgo.Database, permalink string) (bson.M, error) {
	as := db.C("articles")
	jas := bson.M{}
	start := time.Now()
	query := as.Find(bson.M{"permalink": permalink}).
		Select(bson.M{"title": 1, "author": 1, "categories": 1,
		"body": 1, "permalink": 1, "preview": 1,
		"source": 1, "publishedAt": 1, "char_level": 1,
		"vocabulary_level": 1, "grammar_level": 1,
		"num_comment": 1, "tags": 1, "visibility": 1})
	end := time.Now()
	tracer.Tracef("DB query: %v\n", end.Sub(start))
	start = end
	err := query.One(&jas)
	if err != nil {
		return nil, err
	}
	end = time.Now()
	tracer.Tracef("Marshalling: %v\n", end.Sub(start))
	return jas, nil
}

func processAm(am bson.M) {
	categories, _ := am["categories"].(string)
	if len(categories) > 0 {
		am["categories"] = parseCategories(categories)
	}

	scl, svl, sgl := am["char_level"], am["vocabulary_level"], am["grammar_level"]
	var err error
	if cl, ok := scl.(string); ok {
		if am["char_level"], err = strconv.Atoi(cl); err != nil {
			am["char_level"] = 1
		}
	}

	if vl, ok := svl.(string); ok {
		if am["vocabulary_level"], err = strconv.Atoi(vl); err != nil {
			am["vocabulary_level"] = 1
		}
	}

	if gl, ok := sgl.(string); ok {
		if am["grammar_level"], err = strconv.Atoi(gl); err != nil {
			am["grammar_level"] = 1
		}
	}
}

func daoUpdateArticle(db *mgo.Database, permalink string, in []byte) error {
	as := db.C("articles")
	am := bson.M{}
	if err := json.Unmarshal(in, &am); err != nil {
		return err
	}
	if _, ok := am["_id"]; ok {
		delete(am, "_id")
	}

	processAm(am)

	err := as.Update(bson.M{"permalink": permalink}, bson.M{"$set": am})
	if err == nil {
		nsqMessenger.Publish("readcn-segword", []byte(`{"permalink": "`+am["permalink"].(string)+`"}`))
	}

	return err
}

func daoNewArticle(db *mgo.Database, in []byte) error {
	as := db.C("articles")
	am := bson.M{}
	if err := json.Unmarshal(in, &am); err != nil {
		return err
	}
	am["sno"] = daoMaxSno(as) + 1

	processAm(am)

	err := as.Insert(am)
	if err == nil {
		nsqMessenger.Publish("readcn-segword", []byte(`{"permalink": "`+am["permalink"].(string)+`"}`))
	}
	return err
}

func parseCategories(s string) []string {
	cats := parseTags(s)
	result := make([]string, 0, len(cats)*2)
	for _, v := range cats {
		result = append(result, allCategories(v)...)
	}
	sort.Strings(result)
	result = rmdup(result)
	return result
}

func rmdup(ss []string) []string {
	result := make([]string, 0, len(ss))
	cv := ss[0]
	result = append(result, cv)
	for _, v := range ss {
		if v != cv {
			cv = v
			result = append(result, v)
		}
	}
	return result
}

func allCategories(s string) []string {
	result := make([]string, 0, 3)
	if strings.Contains(s, ">") {
		result = append(result, allCategories(stripLastCategory(s))...)
	}
	return append(result, s)
}

func stripLastCategory(s string) string {
	i := strings.LastIndex(s, ">")
	return s[:i]
}

func parseTags(s string) []string {
	return func(ss []string) []string {
		for i := 0; i < len(ss); i++ {
			ss[i] = strings.TrimSpace(ss[i])
		}
		return ss
	}(regexp.MustCompile("[,;，；]+").Split(s, -1))
}

func daoNewId(db *mgo.Database) bson.M {
	am := bson.M{}
	id := bson.NewObjectId()
	am["_id"] = id
	if _, ok := am["permalink"]; !ok {
		am["permalink"] = id.Hex()
	}
	return am
}

func daoDeleteArticle(db *mgo.Database, permalink string) error {
	as := db.C("articles")
	return as.Remove(bson.M{"permalink": permalink})
}

func daoSearchArticlesByChar(db *mgo.Database, ch string, grade int, start, num int) ([]bson.M, error) {
	as := db.C("articles")
	jas := make([]bson.M, num)
	query := bson.M{"charlist": ch, "char_level": grade, "visibility": bson.M{"$gte": 1}}
	err := as.Find(query).
		Select(bson.M{"title": 1, "preview": 1, "permalink": 1,
		"source": 1, "publishedAt": 1, "char_level": 1,
		"vocabulary_level": 1, "grammar_level": 1,
		"num_comment": 1, "tags": 1}).
		Sort("-sno").Skip(start).Limit(num).All(&jas)
	if err != nil {
		return nil, err
	}
	return jas, nil
}

func daoSearchArticlesByVocabulary(db *mgo.Database, word string, grade int, start, num int) ([]bson.M, error) {
	as := db.C("articles")
	jas := make([]bson.M, num)
	query := bson.M{"wseg": word, "vocabulary_level": grade, "visibility": bson.M{"$gte": 1}}
	err := as.Find(query).
		Select(bson.M{"title": 1, "preview": 1, "permalink": 1,
		"source": 1, "publishedAt": 1, "char_level": 1,
		"vocabulary_level": 1, "grammar_level": 1,
		"num_comment": 1, "tags": 1}).
		Sort("-sno").Skip(start).Limit(num).All(&jas)
	if err != nil {
		return nil, err
	}
	return jas, nil
}

func daoSearchArticlesByCategory(db *mgo.Database, category string, start, num int) ([]bson.M, error) {
	as := db.C("articles")
	jas := make([]bson.M, num)
	query := bson.M{"categories": category, "visibility": bson.M{"$gte": 1}}
	err := as.Find(query).
		Select(bson.M{"title": 1, "preview": 1, "permalink": 1,
		"source": 1, "publishedAt": 1, "char_level": 1,
		"vocabulary_level": 1, "grammar_level": 1,
		"num_comment": 1, "tags": 1}).
		Sort("-sno").Skip(start).Limit(num).All(&jas)

	if err != nil {
		return nil, err
	}
	return jas, nil
}

func daoSearchArticlesByGrammar(db *mgo.Database, grammar string, grade int, start, num int) ([]bson.M, error) {
	as := db.C("articles")
	jas := make([]bson.M, num)
	// query := bson.M{"pseg": grammar, "vocabulary_level": grade, "visibility": bson.M{"$gte": 1}}
	query := bson.M{
		"vocabulary_level": grade, 
		"visibility": bson.M{"$gte": 1},
		"pseg": bson.M{"$regex": grammar},
	}


	err := as.Find(query).
		Select(bson.M{"title": 1, "preview": 1, "permalink": 1,
		"source": 1, "publishedAt": 1, "char_level": 1,
		"vocabulary_level": 1, "grammar_level": 1,
		"num_comment": 1, "tags": 1}).
		Sort("-sno").Skip(start).Limit(num).All(&jas)
	if err != nil {
		return nil, err
	}
	return jas, nil
}
