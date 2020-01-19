from autoflow import AutoFlow
from webdevice import WebDevice

class Client:
    def __init__(self):
        self.flow_server = AutoFlow()
        self.devices = []

    def init_server(self, test_json_file):
        self.flow_server.load_from_file(test_json_file)

    def add_device(self, device):
        self.devices.append(device)

    def execute(self):
        while True:
            step = self.query_server()
            if step is None:
                print("no step received. test is done")
                break
            self.perform_step(step)


    def query_server(self):
        step = self.flow_server.get_next_step()
        return step

    def perform_step(self, step):
        pass

client = Client()
client.init_server("flow.json")
client.add_device(WebDevice())
client.execute()
