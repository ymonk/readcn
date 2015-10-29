__author__ = 'myan'

import pymongo

connection = pymongo.MongoClient("mongodb://localhost")
db = connection.readcn_dev
articles = db.articles
count = 0

xiaoshuo_authors = {'余华', '刘心武', '方芳', '毕淑敏', '池莉', '徐坤', '王朔', '王蒙', '老舍',
           '浩然', '刘连群', '刘震云', '田小菲', '张郎郎', '阿城', '尤凤伟', '乔典运',
           '廉声', '礼平', '何玉茹', '张洁', '冯向光', '简平', '赵琪', '陆文夫', '汪曾祺'}

zhuanji_authors = {'冯骥才'}
xinli_authors = {'方富熹', '方格'}
zhexue_authors = {'冯友兰', '涂又光'}

sourcecat = {
    "百科 幽默智慧故事": ['故事', "故事> 幽默"],
    '百科 世界宗教': ['百科', '百科>社会', '百科>社会>宗教'],
    '百科 体育天地': ['百科', '百科>社会', '百科>社会>体育'],
    '百科 医学史话': ['百科', '百科>自然', '百科>自然>医学'],
    '百科 控制论': ['百科', '百科>自然', '百科>自然>工程'],
    '百科 海洋奥秘': ['百科', '百科>自然', '百科>自然>地理'],
    '百科 环球揽胜': ['百科', '百科>自然', '百科>自然>地理'],
    '百科 系统工程': ['百科', '百科>自然', '百科>自然>工程'],
    '百科 认识世界': ['百科'],
    '百科 认识中国': ['百科'],
    '丁玲散文集': ['文学', '文学>散文'],
    '巴金散文集': ['文学', '文学>散文'],
    '柏杨散文集': ['文学', '文学>散文'],
    '老舍散文集': ['文学', '文学>散文'],
    '教材 格言': ['文学', '文学>其他', '文学>其他>格言'],
    '教材 邮件要点': ['应用文', '应用文>日常文书', '应用文>日常文书>邮件'],
    '人教版小学语文*': ['课文', '课文>第一语言用', '课文>第一语言用>小学语文'],
    '上3-*': ['课文', '课文>第二语言用', '课文>第二语言用>商务汉语'],
    '中国古代寓言': ['故事', '故事>寓言'],
    '伊索寓言': ['故事', '故事>寓言'],
    '巧舌胜似强兵 口才故事250例': ['故事', '故事>智慧'],
    '幽默故事*': ['故事', '故事>幽默'],
    '语文故事': ['故事', '故事>智慧'],
    '趣味历史故事': ['故事', '故事>智慧'],
    '鲁迅文集': ['文学', '文学>小说'],
    '评论': ['文学', '文学>其他', '文学>其他>评论'],
    '快乐汉语 第二册': ['课文', '课文>第二语言用', '课文>第二语言用>实用汉语'],
    '广告语': ['应用文', '应用文>行业文书', '应用文>行业文书>广告文案'],
    '幼儿童谣': ['文学', '文学>诗歌', '文学>诗歌>童谣'],
    '多音字一篇通': ['课文', '课文>其他'],
    '传记 摩托罗拉创业者': ['文学', '文学>传记'],
    '传记 福特': ['文学', '文学>传记'],
    '传记 美国酒王': ['文学', '文学>传记'],
    '传记 老洛克菲勒': ['文学', '文学>传记'],
    '传记 苦心经营': ['文学', '文学>传记'],
    '人教网': ['故事', '故事>成语故事']
}



def main():
    global count
    for doc in articles.find({}):
        try:
            author = doc['author']
            author = author.split(';')[0].strip()
            if author in xiaoshuo_authors:
                doc['categories'] = ['文学', '文学>小说']
            elif author in zhuanji_authors:
                doc['categories'] = ['文学', '文学>传记']
            elif author in xinli_authors:
                doc['categories'] = ['文学', '社会>心理']
            elif author in zhexue_authors:
                doc['categories'] = ['文学', '社会>哲学']
            count += 1
            print("Categorized:", doc['title'], ' - ', count)
        except:
            try:
                source = doc['source']
                if source in sourcecat:
                    doc['categories'] = sourcecat[source]
                if source[:3] == '上3-':
                    doc['categories'] = ['课文', '课文>第二语言用', '课文>第二语言用>商务汉语']
                if source[:4] == '幽默故事':
                    doc['categories'] = ['故事', '故事>幽默']
                articles.replace_one({'_id': doc['_id']}, doc)
                count += 1
                print("Categorized:", doc['title'], ' - ', count)
            except:
                pass


    print("Got", count, " documents")


if __name__ == '__main__':
    main()


