from model.BertClassifier import BertClassifier
from server import server


if __name__ == '__main__':
    clf = BertClassifier()
    server.serve(clf)
