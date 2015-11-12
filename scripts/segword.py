import pymongo
import nsq
import json
import jieba
import jieba.posseg as posseg

dbconn = pymongo.MongoClient("mongodb://localhost")
db = dbconn.readcn_dev
articles = db.articles

def pseg(itpair):
  lst = []
  for word, flag in itpair:
    phase = u"{0}/{1}".format(word, flag)
    lst.append(phase)
  return lst


def SegWordHandler(message):
  jsonobj = json.loads(message.body)
  permalink = jsonobj[u'permalink']
  for doc in articles.find({u'permalink': permalink}):
    text = doc[u'body']
    doc[u'wseg'] = list(jieba.cut(text, cut_all=True))
    doc[u'pseg'] = pseg(posseg.cut(text))
    articles.replace_one({u'permalink': permalink}, doc)
    print("Successfully Segement document: " + doc[u'title'])
  return True

if __name__ == '__main__':
  r = nsq.Reader(message_handler=SegWordHandler, lookupd_http_addresses=['http://localhost:4161'],
                 topic='readcn-segword', channel='primary', lookupd_poll_interval=15)
  nsq.run()

