from selenium import webdriver
from selenium.webdriver.common.keys import Keys

class WebDevice:
    def __init__(self):
        self.driver = None

    def create_browser(self, params):
        path = "./{}driver".format(params['value'])
        self.driver = webdriver.Chrome(path)
        return {}

    def open_url(self, params):
        self.driver.get(params['value'])
        return {}

    def find_by_id(self, params):
        self.driver.find_element_by_id(params['id'])
        return {}

    def send_keys_to_id(self, params):
        element = self.driver.find_element_by_id(params['id'])
        element.send_keys(params['value'])
        return {}

    def click_on_id(self, params):
        element = self.driver.find_element_by_id(params['id'])
        element.click()
        return {}

    def send_value_by_classname(self, params):
        element = self.driver.find_element_by_class_name(params['classname'])
        result = {}
        result[params['resultKey']]=element.text
        return result

