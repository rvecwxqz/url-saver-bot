import grpc
from proto import bert_server_pb2
from proto import bert_server_pb2_grpc
from concurrent import futures


class BertServer(bert_server_pb2_grpc.BertClassifierServicer):
    def __init__(self, clf):
        self._clf = clf

    def Predict(self, request, context):
        pred = self._clf.predict(request.text)
        resp = bert_server_pb2.PredictResponse(prediction=pred)
        return resp


def serve(clf):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    s = BertServer(clf)
    bert_server_pb2_grpc.add_BertClassifierServicer_to_server(s, server)
    server.add_insecure_port("localhost:3233")
    server.start()
    print("classifier server started")
    server.wait_for_termination()
