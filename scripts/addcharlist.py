import pymongo

connection = pymongo.MongoClient("mongodb://localhost")
db = connection.readcn_dev
articles = db.articles
count = 0

def is_cnchar(ch):
    return 0x4e00 <= ord(ch) < 0x9fa6

def to_charlist(text):
    return [ch for ch in text if ch is not ' ' and is_cnchar(ch)]


def add_charlist():
    global count
    for doc in articles.find({}):
        text = doc['body']
        chars = to_charlist(text)
        doc['charlist'] = chars
        articles.replace_one({'_id': doc['_id']}, doc)
        count += 1
        print("Updated", doc['title'], ":", count)

def main():
    global count
    for doc in articles.find({}):
        count += 1
        print("Found", doc['title'], ":", count)


if __name__ == '__main__':
    main()

