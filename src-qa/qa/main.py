import grpc

from qa.protos.ch1_pb2_grpc import InventoryServiceStub
from qa.protos.ch1_pb2 import *

def main():
    channel = grpc.insecure_channel('localhost:8080')
    stub = InventoryServiceStub(channel)


    r = t1(stub)
    if r:
        print(f"QA: {r}")




def t1(s: InventoryServiceStub):

    # given
    resp = s.AddLocation(AddLocationRequest(name="test"))
    l = s.ListLocation(ListLocationRequest(location=resp.id))

    if len(l.items) != 1:
        return "Expected to see one record for that location"












