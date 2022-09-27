import grpc

from qa.protos.ch1_pb2_grpc import InventoryServiceStub
from qa.protos.ch1_pb2 import *

def main():
    channel = grpc.insecure_channel('localhost:8080')
    stub = InventoryServiceStub(channel)


    t1(stub)

    blah: type[str]



def t1(s: InventoryServiceStub):

    resp = s.AddLocation(AddLocationRequest())
    









