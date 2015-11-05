# -*- encoding: utf-8 -*-

import pymongo

connection = pymongo.MongoClient("mongodb://localhost")
db = connection.readcn_dev
articles = db.articles
count = 0

xiaoshuo_authors = {u'余华', u'刘心武', u'方芳', u'毕淑敏', u'池莉', u'徐坤', u'王朔', u'王蒙', u'老舍',
           u'浩然', u'刘连群', u'刘震云', u'田小菲', u'张郎郎', u'阿城', u'尤凤伟', u'乔典运',
           u'廉声', u'礼平', u'何玉茹', u'张洁', u'冯向光', u'简平', u'赵琪', u'陆文夫', u'汪曾祺'}

jishi_authors = {u'冯骥才'}
xinli_authors = {u'方富熹', u'方格'}
zhexue_authors = {u'冯友兰', u'涂又光'}

sourcecat = {
    u"百科 幽默智慧故事": [u'故事', u"故事> 幽默"],
    u'百科 世界宗教': [u'百科', u'百科>社会', u'百科>社会>宗教'],
    u'百科 体育天地': [u'百科', u'百科>社会', u'百科>社会>体育'],
    u'百科 医学史话': [u'百科', u'百科>自然', u'百科>自然>医学'],
    u'百科 控制论': [u'百科', u'百科>自然', u'百科>自然>工程'],
    u'百科 海洋奥秘': [u'百科', u'百科>自然', u'百科>自然>地理'],
    u'百科 环球揽胜': [u'百科', u'百科>自然', u'百科>自然>地理'],
    u'百科 系统工程': [u'百科', u'百科>自然', u'百科>自然>工程'],
    u'百科 认识世界': [u'百科'],
    u'百科 认识中国': [u'百科'],
    u'丁玲散文集': [u'文学', u'文学>散文'],
    u'巴金散文集': [u'文学', u'文学>散文'],
    u'柏杨散文集': [u'文学', u'文学>散文'],
    u'老舍散文集': [u'文学', u'文学>散文'],
    u'教材 格言': [u'文学', u'文学>其他', u'文学>其他>格言'],
    u'教材 邮件要点': [u'应用文', u'应用文>日常文书', u'应用文>日常文书>邮件'],
    u'中国古代寓言': [u'故事', u'故事>寓言'],
    u'伊索寓言': [u'故事', u'故事>寓言'],
    u'巧舌胜似强兵 口才故事250例': [u'故事', u'故事>智慧'],
    u'语文故事': [u'故事', u'故事>智慧'],
    u'趣味历史故事': [u'故事', u'故事>智慧'],
    u'鲁迅文集': [u'文学', u'文学>小说'],
    u'评论': [u'文学', u'文学>其他', u'文学>其他>评论'],
    u'快乐汉语 第二册': [u'课文', u'课文>第二语言用', u'课文>第二语言用>实用汉语'],
    u'广告语': [u'应用文', u'应用文>行业文书', u'应用文>行业文书>广告文案'],
    u'幼儿童谣': [u'文学', u'文学>诗歌', u'文学>诗歌>童谣'],
    u'多音字一篇通': [u'课文', u'课文>其他'],
    u'传记 摩托罗拉创业者': [u'文学', u'文学>纪实文学', u'文学>纪实文学>传记'],
    u'传记 福特': [u'文学', u'文学>纪实文学', u'文学>纪实文学>传记'],
    u'传记 美国酒王': [u'文学', u'文学>纪实文学', u'文学>纪实文学>传记'],
    u'传记 老洛克菲勒': [u'文学', u'文学>纪实文学', u'文学>纪实文学>传记'],
    u'传记 苦心经营': [u'文学', u'文学>纪实文学', u'文学>纪实文学>传记'],
    u'人教网': [u'故事', u'故事>成语故事']
}



def main():
    global count
    for doc in articles.find({}):
        try:
            del(doc['categories'])
            author = doc['author']
            author = author.split(';')[0].strip()
            if author in xiaoshuo_authors:
                print 'Got ' + author
                doc['categories'] = [u'文学', u'文学>小说']
            elif author in jishi_authors:
                doc['categories'] = [u'文学', u'文学>纪实文学', u'文学>纪实文学>传记', u'文学>纪实文学>报告文学', ]
            elif author in xinli_authors:
                doc['categories'] = [u'百科', u'百科>社会',u'百科>社会>心理']
            elif author in zhexue_authors:
                doc['categories'] = [u'百科', u'百科>社会',u'百科>社会>哲学']
            articles.replace_one({'_id': doc[u'_id']}, doc)
            count += 1
            print "Categorized:", doc[u'title'], ' - ', doc['categories'], ' : ', count
        except:
            try:
                source = doc['source']
                title = doc['title']
                if source in sourcecat:
                    #print source, u':', doc[u'title']
                    doc['categories'] = sourcecat[source]
                if source.startswith(u'人教版小学语文'):
                    doc['categories'] = [u'课文', u'课文>第一语言用', u'课文>第一语言用>小学语文']
                if source.startswith(u'上3-'):
                    #print source, u':', doc[u'title']
                    doc['categories'] = [u'课文', u'课文>第二语言用', u'课文>第二语言用>商务汉语']
                if source.startswith(u'幽默故事'):
                    #print source, u':', doc[u'title']
                    doc['categories'] = [u'故事', u'故事>幽默']

                if title.startswith(u'编辑部'):
                    doc['categories'] = [u'文学', u'文学>剧本']
                elif title.startswith(u'朱与中国'):
                    doc['categories'] = [u'文学', u'文学>其他', u'文学>其他>评论']


                articles.replace_one({'_id': doc[u'_id']}, doc)
                count += 1
                print "Categorized:", doc[u'title'], ' - ', doc['categories'], ' : ', count
            except:
                pass


    print("Got", count, " documents")


if __name__ == '__main__':
    main()


