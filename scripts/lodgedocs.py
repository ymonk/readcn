__author__ = 'myan'

import pymongo
import datetime
import os
import re
import jieba.posseg as pseg


pat = re.compile(r'^\^([^\n]+)$[\n\s]*(.+)', re.DOTALL | re.MULTILINE)
connection = pymongo.MongoClient("mongodb://localhost")
db = connection.readcn_dev
articles = db.articles


def toPsegList(itpair):
    lst = []
    for word, flag in itpair:
        phase = "{0}/{1}".format(word, flag)
        lst.append(phase)
    return lst


def guillotine(text):
    match = pat.search(text)
    return match.group(1), match.group(2)


def main():
    text_dir = "./text/"
    text_docs = os.listdir(text_dir)
    count = 0
    for fname in text_docs:
        with open(os.path.join(text_dir, fname), 'r') as f:
            whole = f.read()
            title, body = guillotine(whole)
            post = {'title': title, 'body': body, 'source': '人教网',
                    'publishedAt': datetime.datetime.utcnow(),
                    'wordsegs': toPsegList(pseg.cut(body))}
            articles.insert_one(post)
        count += 1
        print("Succeeded on ", fname, ": ", count)


if __name__ == '__main__':
    main()

