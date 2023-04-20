
class OperatorClient:
    def __init__(self, address):
        self._address = address
    
    def start_trainer(self, grpc_port, proxy_address, task_id, spec):
        return f"test id"
    
    def stop_trainer(self, context):
        return True

    def start_worker_group(self, trainer_address, number, specs):
        return "test worker group id"
    
    def stop_worker_group(self, context):
        return True
