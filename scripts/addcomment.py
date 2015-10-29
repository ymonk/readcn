import pymongo

def main():
    connection = pymongo.MongoClient("mongodb://localhost")
    db = connection.readcn_dev
    articles = db.articles
    count = 0
    for doc in articles.find({}):
        doc['num_read'] = 0
        doc['num_comment'] = 0
        doc['comments'] = []
        text = doc['body'].strip()
        doc['body'] = text
        tl = len(text)
        al = tl if tl < 100 else 100
        trail = '' if tl < 100 else '...'
        abstract = text[0:al]
        doc['preview'] = abstract + trail
        articles.replace_one({'_id': doc['_id']}, doc)
        count += 1
        print("Updated ", doc['title'], ":", count)


if __name__ == '__main__':
    main()

