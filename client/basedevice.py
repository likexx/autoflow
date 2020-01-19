from abc import ABC, abstractmethod

class BaseDevice(ABC):
    def __init__(self):
        pass

    @abstractmethod
    def create_device(self, device_name):
        pass

    @abstractmethod
    def open_url(self, url):
        pass

    @abstractmethod
    def find_by_id(self, id):
        pass

    @abstractmethod
    def send_keys_to_id(self, id, value):
        pass

    @abstractmethod
    def click_on_id(self, id):
        pass




