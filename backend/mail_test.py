import mail
import unittest
import weekmail
from datetime import datetime
from datetime import timedelta
import io
import sys
import time

class Test_TestAccountValidator(unittest.TestCase):
    def test_sendMail(self):
        self.assertEqual(mail.sendMail('davidm','no-reply@chalmers.it','Test','Test') , True)
    
    def test_weekmail(self):
        
        now = datetime.now()
        now += timedelta(0,1)

        capturedOutput = io.StringIO()
        sys.stdout = capturedOutput

        wm = weekmail.weekmail(now.weekday(),now.hour,now.minute,now.second)
        wm.start()

        time.sleep(2)
        wm.cancel()

        sys.stdout = sys.__stdout__
        self.assertTrue(capturedOutput.getvalue().__contains__("Email sent to P.R.I.T"))