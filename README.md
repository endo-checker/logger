# Logger gRPC microserverice #

This application uses `gRPC` and `proto` to send user data to a database (mongoDB). <br/>
This uses CRUD functions. <br />
For now I haven't added any securtity features, or protected my db connection as this is supposed to only showcase the very basics of CRUD functions in gRPC. 

## Get Started ##
- Install/update dependencies `go get -u`
- Initiate gRPC: `protoc --proto_path=proto patient/v1/patient.proto --go_out=. --go-grpc_out=.` This is the standard way of generating gRPC, however with this app I am using a makefile configured with dapr to manage my commands. So instead run `make proto`, this will also send the proto definitions to `https://buf.build` as a commit. 
- Run the application `go run .` or `make run`

### Dapr ###

Install and configure the dapr CLI:

```bash
brew install dapr/tap/dapr-cli
dapr init
```

### Buf ###

Add Buf API key to `.netrc` file

```bash
machine buf.build password <your Buf API key>
machine go.buf.build login <your Buf username> password <your Buf API key>
```

## Set up Docker ## 
* To make a build on `docker` run `docker build --tag patient .`
* To run on docker `docker run -d -p 8080:8080 patient`

## Postman Setup ##
- Go to New, select gRPC Request
- Call `localhost:8080`
- Select the proto definitions file to put into the methods input
- Make sure server reflection is enabled in methods
- You can generate example JSON Messages (this is particularly useful when using the CREATE method, make sure the ID field isn't present as this is being     filled in automatically when an entry is created) 
- When making a GET Request make sure the JSON message only includes existing ID's <br />
 
```bash 
{ "id": "4942b3e9-3699-412a-9f16-01503113178f"}
```

# Setup Heroku 

## Install the Heroku CLI

* To start you need install Heroku CLI `brew tap heroku/brew && brew install heroku`
* Login into heroku from your terminal with `heroku login`
* Create a new application using either Heroku dashboard - from the 'New' tab OR
  in your terminal run command `heroku create -a <app-name>` . If you dont specify an app name, one will be assigned.
  This will return the heroku url address for your app and a git url.
* `heroku git:remote -a <app-name>` | This command adds a git remote to the app repo
* The default stack that Heroku assigns is 'heroku-20'. You can check this under the settings on the
  dashboard. We need to set the stack to 'container'. Run this command `heroku stack:set container`
  This will also automatically set the Framework type to Container as well.
* Setup a enviroment variable on Heroku dashboard. Settings -> Config Vars -> "PORT" : "8080" 
  This is for the first initialization.
* Then run `git push heroku main` to trigger deployment from your local machine to Heroku remote.
* Once this has finished building, you will find your URL for the container under the settings tab.
* heroku.yaml | Add this file to your code repo. This tells Heroku to build a Docker container according to the
  Docker file located in the root directory of the application.
* https://devcenter.heroku.com/articles/heroku-cli-commands | Heroku CLI commands# logger
