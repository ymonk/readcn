
import pymongo

connection = pymongo.MongoClient("mongodb://localhost")
db = connection.readcn_dev
articles = db.articles
count = 0


def main():
    global count
    for doc in articles.find({}):
        text = doc['body']
        tl = len(text)
        al = tl if tl < 100 else 100
        abstract = text[0:al]
        doc['abstract'] = abstract + "..."
        articles.replace_one({'_id': doc['_id']}, doc)
        count += 1
        print("Updated", doc['title'], ":", count)


if __name__ == '__main__':
    main()

