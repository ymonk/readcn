import pymongo

def main():
    connection = pymongo.MongoClient("mongodb://localhost")
    db = connection.readcn_dev
    articles = db.articles
    count = 0
    for doc in articles.find({}):
        doc['permalink'] = str(doc['_id'])
        articles.replace_one({'_id': doc['_id']}, doc)
        count += 1
        print("Add", doc['_id'], ":", count)


if __name__ == '__main__':
    main()

