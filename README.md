# Webservice allowing for synchronisation.

## Usage
Run the cloud sync-point with:
```sh
cargo run
```

If you have Golang installed, you can use the test.go script to send POST requests to the cloud sync-point
I'm not sure how the google/uuid package will be shared so if it doesn't get automatically imported please run:
```sh
go get github.com/google/uuid
```
To launch the test script run:
```sh
go run test.go
```
The script will run 3 tests, one positive test, one test that triggers the timeout and one tests will different uuids which will also trigger the timeout.

If you do not have Golang installed, you can use a tool like postman or a terminal with curl commands.
All you will need is a way tocreate some uuids.
To create a UUID, on Unix based systems you can just run 
```sh
uuidgen
```
Otherwise there are online tools to do it such as https://www.uuidgenerator.net/

Once you have a uuid, you can use curl to send requests, the easiest way is to have 2 terminals and send a request in each of them like so:

First terminal:
```sh
curl -X POST http://127.0.0.1:3030/wait-for-second-party/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
```
Second terminal:
```sh
curl -X POST http://127.0.0.1:3030/wait-for-second-party/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
```
You should see output in the terminal running the cloud sync-point like that:
```sh
cloud sync-point running at: 127.0.0.1:3030

request id: a3fd8c55-2ffb-4ec8-9aae-01ae2af46d03
the first party requested the URI
request id: a3fd8c55-2ffb-4ec8-9aae-01ae2af46d03
the second party requested the URI

request id: 5ed01c5b-44c4-4973-a146-a36086171250
the first party requested the URI
the second party did not request the URI, reponding with TIMEOUT
```