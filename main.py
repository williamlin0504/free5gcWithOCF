from flask import Flask
import os 


app = Flask(__name__)
@app.route('/')

def testSH():
    headers = flask.request.headers
    os.system('./test.sh {} {}' .format('TestRegistration', str(headers)))   
    return 'Test Started!'

if __name__ == '__main__':
  app.run(host="0.0.0.0",port=8080)