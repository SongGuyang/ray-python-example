import argparse
import asyncio
try:
    from grpc import aio as aiogrpc
except ImportError:
    from grpc.experimental import aio as aiogrpc

from generated import (
    automl_service_pb2,
    automl_service_pb2_grpc,
)
import logging
import sys

logging.basicConfig(stream=sys.stdout, format='%(asctime)s %(message)s', level=logging.DEBUG)
logger = logging.getLogger(__name__)

def get_or_create_event_loop() -> asyncio.BaseEventLoop:
    import sys

    vers_info = sys.version_info
    if vers_info.major >= 3 and vers_info.minor >= 10:
        # This follows the implementation of the deprecating `get_event_loop`
        # in python3.10's asyncio. See python3.10/asyncio/events.py
        # _get_event_loop()
        loop = None
        try:
            loop = asyncio.get_running_loop()
            assert loop is not None
            return loop
        except RuntimeError as e:
            # No running loop, relying on the error message as for now to
            # differentiate runtime errors.
            assert "no running event loop" in str(e)
            return asyncio.get_event_loop_policy().get_event_loop()

    return asyncio.get_event_loop()

class Proxy(automl_service_pb2_grpc.AutoMLServiceServicer):
    def __init__(
        self,
        grpc_port: str,
    ):

        self.server = aiogrpc.server(options=(("grpc.so_reuseport", 0),))
        grpc_ip = "0.0.0.0"
        self.grpc_port = self.server.add_insecure_port(f"{grpc_ip}:{grpc_port}")
        logger.info("Proxy grpc address: %s:%s", grpc_ip, self.grpc_port)
    
    async def DoAutoML(self, request, context):
        return automl_service_pb2.DoAutoMLReply(
            success=True,
            task_id=123,
            message="test",
        )
    

    async def GetResult(self, request, context):
        return automl_service_pb2.GetResultReply(
            success=True,
            result="test",
        )


    async def run(self):

        # Start a grpc asyncio server.
        await self.server.start()

        automl_service_pb2_grpc.add_AutoMLServiceServicer_to_server(
            self, self.server
        )

        await self.server.wait_for_termination()

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="AutoML proxy.")
    parser.add_argument(
        "--port", required=True, type=int, help="The grpc port."
    )

    args = parser.parse_args()

    proxy = Proxy(args.port)

    loop = get_or_create_event_loop()

    loop.run_until_complete(proxy.run())
