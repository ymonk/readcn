package main

import (
    "encoding/json"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "time"
"strings"
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
                Select(bson.M{"title": 1, "author": 1,
                              "body": 1, "permalink": 1, "preview": 1,
                              "source": 1, "publishedAt": 1, "char_level": 1,
                              "vocabulary_level": 1, "grammar_level": 1,
                              "num_comment": 1, "tags": 1})
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

func daoUpdateArticle(db *mgo.Database, permalink string, in []byte) error {
    as := db.C("articles")
    am := bson.M{}
    if err := json.Unmarshal(in, &am); err != nil {
        return err
    }
    if _, ok := am["_id"]; ok {
        delete(am, "_id")
    }

    categories, _ := am["categories"].(string)
    if len(categories) > 0 {
        am["categories"] = parseTags(categories)
    }
    err := as.Update(bson.M{"permalink": permalink}, bson.M{"$set": am})
    return err
}

func daoNewArticle(db *mgo.Database, in []byte) error {
    as := db.C("articles")
    am := bson.M{}
    if err := json.Unmarshal(in, &am); err != nil {
        return err
    }
    am["sno"] = daoMaxSno(as) + 1
    categories, _ := am["categories"].(string)
    if len(categories) > 0 {
        am["categories"] = parseTags(categories)
    }
    tracer.Trace("permalink=", am["permalink"])
    err := as.Insert(am)
    return err
}

func parseTags(s string) []string {
    strip := func (ss []string) {
        for i := 0; i < len(ss); i++ {
            ss[i] = strings.TrimSpace(ss[i])
        }
    }

    // try , first
    tags := strings.Split(s, ",")
    if len(tags) <= 1 {
        tags = strings.Split(s, ";")
    }
    if len(tags) <= 1 {
        tags = strings.Split(s, "，")
    }
    if len(tags) <= 1 {
        tags = strings.Split(s, "；")
    }

    strip(tags)
    return tags
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
    err := as.Find(bson.M{"charlist": ch, "char_level": grade}).
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
    err := as.Find(bson.M{"wseg": word, "vocabulary_level": grade}).
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