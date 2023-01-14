import time
import multiprocessing
import threading
import functools
from pynput import keyboard
from playsound import playsound

class Telephone:
	def __init__(self):
		self.audio_thread = None
		self.input_thread = None
		self.monitor_thread = None
		self.number = ""
		self.timeout = 0
	
	def off_hook(self):
		self.number = ""

		self.audio_thread = multiprocessing.Process(target=functools.partial(playsound, 'dialtone.mp3'))
		self.audio_thread.start()

		self.input_thread = keyboard.Listener(on_press = lambda key: self.on_press(key))
		self.input_thread.start()

		self.monitor_thread = threading.Thread(target=self.monitor_input)
		self.monitor_thread.start()
		self.monitor_thread.join()

		self.input_thread.stop()
		self.input_thread.join()

		print(f"Received number: {self.number}")

	def monitor_input(self):
		while self.timeout < 1 or self.number == "":
			self.timeout += 1
			time.sleep(1)

	def on_press(self, key):
		self.timeout = 0
		self.audio_thread.terminate()
		try:
			if key.char.isnumeric():
				self.number += key.char
				playsound(f"{key.char}.wav")
		except AttributeError:
			return

if __name__ == '__main__':
	telephone = Telephone()
	telephone.off_hook()
