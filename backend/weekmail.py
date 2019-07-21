from threading import Timer
from datetime import timedelta
from datetime import datetime
from mail import sendMail

class weekmail():
    __timer = []
    __notifyTime = []

    def __init__(self, _dayIndex, _hour, _minute, _second):
        self.__notifyTime = datetime.now()
        self.__notifyTime = self.__notifyTime.replace(hour=_hour,minute=_minute, second=_second)

        while(self.__notifyTime.weekday() != _dayIndex):
            self.__notifyTime += timedelta(days=1)
    
    def updateTime(self):
        while((self.__notifyTime - datetime.now()).total_seconds() < 0):
            self.__notifyTime += timedelta(days=7)

    def start(self):
        self.updateTime()
        print((self.__notifyTime - datetime.today()).total_seconds())
        self.__timer = Timer((self.__notifyTime - datetime.today()).total_seconds(), self.notifyPRIT)
        self.__timer.start()

    def cancel(self):
        self.__timer.cancel()

    def notifyPRIT(self):
        if sendMail("prit","no-reply@chalmers.it","Veckans förslag", "something"):
            print("Email sent to P.R.I.T")
        else:
            print("Unable to send mail to P.R.I.T")
        self.start()