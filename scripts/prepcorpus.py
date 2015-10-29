__author__ = 'myan'

import os
import codecs

fileCount = 0


def convertCorpusDir(srcdir, dstdir):
    entries = os.listdir(srcdir)
    for entry in entries:
        spath = os.path.join(srcdir, entry)
        dpath = os.path.join(dstdir, entry)

        if os.path.isdir(spath):
            if not os.path.exists(dpath):
                os.mkdir(dpath)
            convertCorpusDir(spath, dpath)
        elif os.path.isfile(spath):
            convertCorpusFile(spath, dpath)


def convertCorpusFile(srcfile, dstfile):
    global fileCount
    with codecs.open(srcfile, 'r', encoding='gbk') as sf:
        with open(dstfile, 'wt') as df:
            print("Coverting ", srcfile, "...")
            try:
                df.write(sf.read())
                fileCount += 1
                print("Done: ", fileCount)
            except UnicodeDecodeError as e:
                print("Decoding error: ", e.reason)

def main():
    convertCorpusDir(srcdir='./corpus_raw', dstdir='./corpus')

if __name__ == '__main__':
    main()