__author__ = 'myan'

import os
from bs4 import BeautifulSoup


def main():
    path = "./html"
    dest = "./text/"
    htmlDocs = os.listdir("./html")
    count = 0
    success = 0
    fail = 0
    for fname in htmlDocs:
        f = open(path + "/" + fname, 'r', encoding="gbk", errors="ignore")
        d = open(dest + fname + ".txt", "w", encoding="utf8", errors="ignore")
        html = f.read()
        f.close()
        soup = BeautifulSoup(html, "html.parser")
        title = soup.find('div', id='doctitle').string
        print(title)
        content = '^' + title + "\n\n"
        for bodyp in soup.find_all('div', id="doccontent"):
            if hasattr(bodyp, "style") and bodyp.style != None:
                bodyp.style.clear()
            if hasattr(bodyp, "body") and hasattr(bodyp.body, "p"):
                if hasattr(bodyp.body.p, "strings"):
                    for string in bodyp.body.p.strings:
                        content += string
                elif hasattr(bodyp.body.p, "string"):
                    content += bodyp.body.p.string
                success += 1
            elif hasattr(bodyp, "strings"):
                for string in bodyp.strings:
                    content += string
                success += 1
            elif hasattr(bodyp, "string"):
                content += bodyp.string
                success += 1
            else:
                fail += 1

        bottom = soup.find('div', id='bottomElements').string
        pos = bottom.find("选自：")
        if pos >= 0:
            source = "\n@"+bottom[pos+3:]+"\n"
            print(source)
            content += source
        d.write(content)
        d.close()
        count += 1
    print("Processed " + str(count) + " documents, success: " + str(success) + ", fail: " + str(fail))

if __name__ == '__main__':
    main()



