__author__ = 'myan'

import pymongo
import datetime
import os
import re
import jieba.posseg as pseg

connection = pymongo.MongoClient("mongodb://localhost")
db = connection.readcn_dev
articles = db.articles
patStart = re.compile(r'^\^', re.MULTILINE | re.DOTALL)
patArticle = re.compile(r'^\^([^\n]+)$[\n\s]*(.+)', re.DOTALL | re.MULTILINE)
count = 0


class MatchException(Exception):
    pass


def to_pseg_list(itpair):
    lst = []
    for word, flag in itpair:
        phase = "{0}/{1}".format(word, flag)
        lst.append(phase)
    return lst


def lodge_file(path):
    with open(path, 'r') as f:
        text = f.read()
        for passage in patStart.split(text):
            source = os.path.basename(path).split('.')[0]
            try:
                title, body = guillotine('^' + passage)
                post = {'title': title, 'body': body, 'source': source,
                        'publishedAt': datetime.datetime.utcnow(),
                        'wordsegs': to_pseg_list(pseg.cut(body))}
                articles.insert_one(post)
            except MatchException:
                pass


def guillotine(text):
    match = patArticle.search(text)
    if match is None:
        raise MatchException
    return match.group(1), match.group(2)


def lodge_corpus(dir):
    global count
    for entry in os.listdir(dir):
        fpath = os.path.join(dir, entry)
        if os.path.isdir(fpath):
            lodge_corpus(fpath)
        elif os.path.isfile(fpath):
            print("Lodging ", fpath, "...")
            lodge_file(fpath)
            count += 1
            print("Done: ", count)


def main():
    lodge_corpus('./corpus')

if __name__ == '__main__':
    main()




