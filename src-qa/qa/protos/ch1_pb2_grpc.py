# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from . import ch1_pb2 as ch1__pb2


class InventoryServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.AddLocation = channel.unary_unary(
                '/protos.InventoryService/AddLocation',
                request_serializer=ch1__pb2.AddLocationRequest.SerializeToString,
                response_deserializer=ch1__pb2.AddLocationResponse.FromString,
                )
        self.AddProduct = channel.unary_unary(
                '/protos.InventoryService/AddProduct',
                request_serializer=ch1__pb2.AddProductRequest.SerializeToString,
                response_deserializer=ch1__pb2.AddProductResponse.FromString,
                )
        self.ChangeQuantity = channel.unary_unary(
                '/protos.InventoryService/ChangeQuantity',
                request_serializer=ch1__pb2.ChangeQuantityRequest.SerializeToString,
                response_deserializer=ch1__pb2.ChangeQuantityResponse.FromString,
                )
        self.ListLocation = channel.unary_unary(
                '/protos.InventoryService/ListLocation',
                request_serializer=ch1__pb2.ListLocationRequest.SerializeToString,
                response_deserializer=ch1__pb2.ListLocationResponse.FromString,
                )


class InventoryServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def AddLocation(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def AddProduct(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ChangeQuantity(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ListLocation(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_InventoryServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'AddLocation': grpc.unary_unary_rpc_method_handler(
                    servicer.AddLocation,
                    request_deserializer=ch1__pb2.AddLocationRequest.FromString,
                    response_serializer=ch1__pb2.AddLocationResponse.SerializeToString,
            ),
            'AddProduct': grpc.unary_unary_rpc_method_handler(
                    servicer.AddProduct,
                    request_deserializer=ch1__pb2.AddProductRequest.FromString,
                    response_serializer=ch1__pb2.AddProductResponse.SerializeToString,
            ),
            'ChangeQuantity': grpc.unary_unary_rpc_method_handler(
                    servicer.ChangeQuantity,
                    request_deserializer=ch1__pb2.ChangeQuantityRequest.FromString,
                    response_serializer=ch1__pb2.ChangeQuantityResponse.SerializeToString,
            ),
            'ListLocation': grpc.unary_unary_rpc_method_handler(
                    servicer.ListLocation,
                    request_deserializer=ch1__pb2.ListLocationRequest.FromString,
                    response_serializer=ch1__pb2.ListLocationResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'protos.InventoryService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class InventoryService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def AddLocation(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/protos.InventoryService/AddLocation',
            ch1__pb2.AddLocationRequest.SerializeToString,
            ch1__pb2.AddLocationResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def AddProduct(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/protos.InventoryService/AddProduct',
            ch1__pb2.AddProductRequest.SerializeToString,
            ch1__pb2.AddProductResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ChangeQuantity(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/protos.InventoryService/ChangeQuantity',
            ch1__pb2.ChangeQuantityRequest.SerializeToString,
            ch1__pb2.ChangeQuantityResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ListLocation(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/protos.InventoryService/ListLocation',
            ch1__pb2.ListLocationRequest.SerializeToString,
            ch1__pb2.ListLocationResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
