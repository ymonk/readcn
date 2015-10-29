import pymongo

connection = pymongo.MongoClient("mongodb://localhost")
db = connection.readcn_dev
articles = db.articles

def main():
    count = 1
    for doc in articles.find({}):
        nchars = len(doc["charlist"])
        char_level = 1
        vocabulary_level = 1
        grammar_level = 1
        if 500 < nchars < 1500:
            char_level, vocabulary_level, grammar_level = 2, 2, 2
        elif 1500 <= nchars < 5000:
            char_level, vocabulary_level, grammar_level = 3, 3, 3
        elif 5000 <= nchars:
            char_level, vocabulary_level, grammar_level = 3, 4, 3

        doc["char_level"], doc["vocabulary_level"], doc["grammar_level"] = char_level, vocabulary_level, grammar_level
        articles.replace_one({'_id': doc['_id']}, doc)
        print(str(count) + " - Add grade to ", doc['title'], ":", char_level, ',', vocabulary_level, ',', grammar_level)
        count += 1


if __name__ == '__main__':
    main()

