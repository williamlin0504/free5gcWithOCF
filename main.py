from flask import Flask, request
import os 


app = Flask(__name__)
@app.route('/')

def testSH():
    ueID = request.headers.get("Ueid")
    os.system('./test.sh {} {}' .format('TestRegistration', ueID))   
    return 'Test Started!'

if __name__ == '__main__':
  app.run(host="0.0.0.0",port=8080)