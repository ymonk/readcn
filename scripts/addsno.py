import pymongo

connection = pymongo.MongoClient("mongodb://localhost")
db = connection.readcn_dev
articles = db.articles

def main():
    global count
    sno = 0
    for doc in articles.find({}):
        doc['sno'] = sno
        articles.replace_one({'_id': doc['_id']}, doc)
        print("Add sno to ", doc['title'], ":", sno)
        sno += 1


if __name__ == '__main__':
    main()

