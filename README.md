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
Otherwise there are online tools to do it such as [https://www.uuidgenerator.net/] https://www.uuidgenerator.net/

Once you have a uuid, you can use curl to send requests, the easiest way is to have 2 terminals and send a request in each of them like so:

First terminal:
```sh
curl -X POST http://127.0.0.1:3030/wait-for-second-party/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
```
Second terminal:
```sh
curl -X POST http://127.0.0.1:3030/wait-for-second-party/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
```
