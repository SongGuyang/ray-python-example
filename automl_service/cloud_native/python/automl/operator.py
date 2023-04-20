
class OperatorClient:
    def __init__(self, address):
        self._address = address
    
    def start_trainer(self, grpc_port, proxy_address, task_id, spec):
        return f"test id"
    
    def stop_trainer(self, trainer_id):
        return True

    def start_worker_group(self, trainer_address, number, specs):
        return "group_id", ["test worker id 1", "test worker id 2"]
    
    def stop_worker_group(self, "group_id"):
        return True
