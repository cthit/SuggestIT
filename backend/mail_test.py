import mail
import unittest


class Test_TestAccountValidator(unittest.TestCase):
    def test_sendMail(self):
        self.assertEqual(mail.sendMail('davidm','no-reply@chalmers.it','Test','Test') , True)