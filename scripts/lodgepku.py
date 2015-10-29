import pymongo
import datetime
import os
import re
import jieba.posseg as pseg

connection = pymongo.MongoClient("mongodb://localhost")
db = connection.readcn_dev
articles = db.articles
pat = re.compile(r'\{\n"author":"([^"]*)",\n"title":"([^"]*)",\n"category":"([^"]+)",\n(?:"author_date":"[^"]+",\n)?(?:"resource":"[^"]+",\n)?"text":"\n([^\}]+)^\}\n',
                 re.MULTILINE | re.DOTALL)
patText = re.compile(r'([^\}]+)\n"$', re.MULTILINE | re.DOTALL)

count = 0

def to_pseg_list(itpair):
    lst = []
    for word, flag in itpair:
        phase = "{0}/{1}".format(word, flag)
        lst.append(phase)
    return lst


def lodgePKU(dirname):
    for fname in os.listdir(dirname):
        print("Processing", fname, "...")
        fpath = os.path.join(dirname, fname)
        with open(fpath, 'r') as f:
            content = f.read()
            lodgeText(content)

def lodgeFile(fpath):
    print("Processing", fpath, "...")
    with open(fpath, 'r') as f:
        lodgeText(f.read())

def is_cnchar(ch):
    return 0x4e00 <= ord(ch) < 0x9fa6

def to_charlist(text):
    return [ch for ch in text if ch is not ' ' and is_cnchar(ch)]


def lodgeText(text):
    global count
    for match in pat.finditer(text):
        author, title, category, body = match.group(1), match.group(2), match.group(3), match.group(4)
        mo = patText.search(body)
        if mo is not None:
            body = mo.group(1)
        post = {'title': title, 'body': body, 'source': '北大语料库',
                'author': author, 'category': category,
                'publishedAt': datetime.datetime.utcnow(),
                'charlist': to_charlist(body),
                'wordsegs': to_pseg_list(pseg.cut(body))}
        articles.insert_one(post)
        count += 1
        print("Lodged ", author, "-", title, ": ", count, "", category)


def main():
    lodgePKU('./pku')


if __name__ == '__main__':
    main()
