# Share & Charge MSP Backoffice API

For any questions, please ask :)

## Usage Guide :crystal_ball:

#### The API is available at: 18.195.223.26:9090
located on AWS with the instance-id: i-0c20e24f744da7c2a and the name: cpo backend api

Quick test: <a href="http://18.195.223.26:9090/api/v1/">Click ME!</a>

if you get "This site can’t be reached" this means, probably, I stopped it for some reason. Ask me on slack why it's stopped. (Andy) :palm_tree:


~~~~
[GIN-debug] GET    /api/v1/                  --> github.com/motionwerkGmbH/msp-backend-api/handlers.Index (3 handlers)
[GIN-debug] DELETE /api/v1/s3cr3tReinitf32fdsfsdf98yu32jlkjfsd89yaf98j320j --> github.com/motionwerkGmbH/msp-backend-api/handlers.Reinit (3 handlers)
[GIN-debug] GET    /api/v1/token/info        --> github.com/motionwerkGmbH/msp-backend-api/handlers.TokenInfo (3 handlers)
[GIN-debug] GET    /api/v1/token/balance/:addr --> github.com/motionwerkGmbH/msp-backend-api/handlers.TokenBalance (3 handlers)
[GIN-debug] POST   /api/v1/token/mint/:addr  --> github.com/motionwerkGmbH/msp-backend-api/handlers.TokenMint (3 handlers)
[GIN-debug] GET    /api/v1/wallet/:addr      --> github.com/motionwerkGmbH/msp-backend-api/handlers.GetWalletBalance (3 handlers)
[GIN-debug] GET    /api/v1/drivers           --> github.com/motionwerkGmbH/msp-backend-api/handlers.GetAllDrivers (3 handlers)
[GIN-debug] GET    /api/v1/msp               --> github.com/motionwerkGmbH/msp-backend-api/handlers.MspInfo (3 handlers)
[GIN-debug] POST   /api/v1/msp               --> github.com/motionwerkGmbH/msp-backend-api/handlers.MspCreate (3 handlers)
[GIN-debug] POST   /api/v1/msp/wallet/generate --> github.com/motionwerkGmbH/msp-backend-api/handlers.MspGenerateWallet (3 handlers)
[GIN-debug] GET    /api/v1/msp/wallet/seed   --> github.com/motionwerkGmbH/msp-backend-api/handlers.MspGetSeed (3 handlers)
[GIN-debug] GET    /api/v1/msp/history       --> github.com/motionwerkGmbH/msp-backend-api/handlers.MSPHistory (3 handlers)
[GIN-debug] GET    /api/v1/msp/locations     --> github.com/motionwerkGmbH/msp-backend-api/handlers.MspGetLocations (3 handlers)
~~~~


## Install Guide :sun_with_face:

### Attention :fire: if you are trying to follow this steps and get stuck at something, it's very important that you update this readme with the fix, so other developers will not encounter the same problem.


1. Get an Ubuntu Instance
2. Install Golang. Configure Golang's GOROOT, GOPATH.

~~~~
cd /tmp
wget -q https://storage.googleapis.com/golang/getgo/installer_linux
chmod +x installer_linux 
./installer_linux 
source $HOME/.bash_profile

echo 'export GOPATH=$HOME/go' >> ~/.bashrc 
echo 'export PATH=${PATH}:${GOPATH}/bin' >> ~/.bashrc 
source ~/.bashrc 

go get github.com/golang/example/hello
test it: ~/go/src/github.com/golang/example/hello$ go run hello.go
~~~~

3. Under your GOPATH (ex: /home/you/go/)

create the directory ~/go/src/github.com/motionwerkGmbH/

into that directory run: git clone git@github.com:motionwerkGmbH/msp-backend-api.git (remember to have this command work, you need to add your ssh key into github)

4. the share & charge config files are under configs/sc_configs. Also there you'll find a script called copy.sh that will copy this configs to ~/.sharecharge folder!
5. chmod +x copy.sh then ./copy.sh
6. Install all the dependencies of this app with: go get ./...  (it will take ~1 min)

## Configure Share & Charge API

this api is based on share & charge api :), so make sure you have it running on localhost:3000

~~~~
cd sharecharge-api
git branch
npm install
npm run install
npm run start
~~~~

#### Running the API Server

Under the msp-backend-api folder

~~~~
go run *.go
~~~~


## FAQ :question:

1. I want to run it in the background

Create the file /var/log/msp_backend.log and give it appropriate permissions
Supervisor. Here's a config file:

~~~~
[program:mspbackendapi]
user=ubuntu
numprocs=1
command=/home/ubuntu/go/src/github.com/motionwerkGmbH/msp-backend-api/backend
directory=/home/ubuntu/go/src/github.com/motionwerkGmbH/msp-backend-api/
autostart=true
autorestart=true
redirect_stderr=true
stdout_logfile=/var/log/msp_backend.log
stdout_logfile_maxbytes=10MB
stdout_logfile_backups=1
~~~~


#### Licence Mozilla Public License Version 2.0

why this license ? see https://christoph-conrads.name/why-i-chose-the-mozilla-public-license-2-0/
