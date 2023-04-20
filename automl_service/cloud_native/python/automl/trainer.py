import argparse
import asyncio
try:
    from grpc import aio as aiogrpc
except ImportError:
    from grpc.experimental import aio as aiogrpc

from automl.generated import (
    automl_service_pb2,
    automl_service_pb2_grpc,
)

from automl.operator import OperatorClient
from automl.utils import get_or_create_event_loop

import logging
import sys

logging.basicConfig(stream=sys.stdout, format='%(asctime)s %(message)s', level=logging.DEBUG)
logger = logging.getLogger(__name__)

class Proxy(automl_service_pb2_grpc.WorkerRegisterService):
    class Context:
        def __init__(self, result = None):
            self.result = result

    def __init__(
        self,
        args,
    ):

        self.server = aiogrpc.server(options=(("grpc.so_reuseport", 0),))
        grpc_ip = "0.0.0.0"
        self.grpc_port = self.server.add_insecure_port(f"{grpc_ip}:{args.grpc_port}")
        logger.info("Proxy grpc address: %s:%s", grpc_ip, self.grpc_port)
        self._host_name = args.host_name
        self._operator_client = OperatorClient(args.operator_address)

    async def WorkerRegister(self, request, context):
        if request.worker_id not in self._workers:
            return automl_service_pb2.WorkerRegisterReply(
                success=False,
                message=f"Task id {request.worker_id} not found.",
            )
        context = self._workers[request.worker_id]
        return automl_service_pb2.WorkerRegisterReply(
            success=True,
            model=,
            df=,
            train_indices=,
            test_indices=,
            label_column=,
            metrics=,
            freq=,
        )

    async def WorkerReportResult(self, request, context):
        if request.worker_id not in self._workers:
            return automl_service_pb2.WorkerReportResultReply(
                success=False,
                message=f"Task id {request.worker_id} not found.",
            )
        self._workers[request.worker_id].result = request.result
        return automl_service_pb2.WorkerReportResultReply(
            success=True,
        )

    async def run(self):

        # Start a grpc asyncio server.
        await self.server.start()

        automl_service_pb2_grpc.add_WorkerRegisterService_to_server(
            self, self.server
        )

        await self.server.wait_for_termination()

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="AutoML proxy.")
    parser.add_argument(
        "--grpc-port", required=True, type=int, help="The grpc port."
    )
    parser.add_argument(
        "--host-name", required=True, type=str, help="The current host name."
    )
    parser.add_argument(
        "--proxy-address", required=True, type=str, help="The grpc address of automl proxy."
    )
    parser.add_argument(
        "--trainer-id", required=True, type=str, help="The current trainer id."
    )
    parser.add_argument(
        "--task-id", required=True, type=str, help="The automl task id."
    )
    parser.add_argument(
        "--operator-address", required=True, type=str, help="The automl operator address."
    )

    args = parser.parse_args()

    proxy = Proxy(args)

    loop = get_or_create_event_loop()

    loop.run_until_complete(proxy.run())
