from webdevice import WebDevice
import requests
import sys

class Client:
    def __init__(self, device):
        self.device = device
        self.flowname = None
        self.session_id = None

    def create_session(self, flowname):
        self.flowname = flowname
        resp = requests.post('http://127.0.0.1:8080/autoflow/create/{}'.format(flowname))
        data = resp.json()
        if data.get('Error', None) != 0:
            printf("failed creating session")
            return False

        self.session_id = data['Result']
        return True


    def execute(self):
        action, params = self.query_next()
        while True:
            if action == "":
                self.stop()
                break

            result = self.perform_action(action, params)
            if result is None:
                action, params = self.send_error()
            else:
                action, params = self.query_next(result)


    def query_next(self, action_result = {}):
        print(action_result)
        url = 'http://127.0.0.1:8080/autoflow/next/{}'.format(self.session_id)
        resp = requests.post(url,json = action_result)
        data = resp.json()
        action = data['Action']
        params = data['Parameters']
        return action, params

    def perform_action(self, action, params):
        try:
            print(action)
            print(params)

            method = getattr(self.device, action)
            result = method(params)

            return result
        except Exception as err:
            print(err)
            return None

    def send_error(self):
        url = 'http://127.0.0.1:8080/autoflow/error/{}'.format(self.session_id)
        data = requests.post(url).json()
        action = data['Action']
        params = data['Parameters']
        return action, params


    def stop(self):
        print("stop")
        url = 'http://127.0.0.1:8080/autoflow/stop/{}'.format(self.session_id)
        requests.post(url)

def main():
    client = Client(WebDevice())
    if not client.create_session(sys.argv[1]):
        exit(1)

    client.execute()


if __name__ == "__main__":
    main()
