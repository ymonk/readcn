import pymongo

connection = pymongo.MongoClient("mongodb://localhost")
db = connection.readcn_dev
articles = db.articles

def wseg(words):
    return [w.split('/')[0] for w in words]


def main():
    count = 1
    for doc in articles.find({}):
        try:
            doc['pseg'] = doc['wordsegs']
            doc['wseg'] = wseg(doc['wordsegs'])
            del(doc['wordsegs'])
            articles.replace_one({'_id': doc['_id']}, doc)
            print("Successfully updated ", doc['title'], ":", count)
            count += 1
        except:
            pass

if __name__ == '__main__':
    main()

