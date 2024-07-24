### peer-programming session
* Create users via API and asynchronously process user verification
* Transfer money between users asynchronously
* Supports graceful shutdown of server

### Running
```shell
export WORKER_COUNT=5 PORT=9001 && go run main.go
```

### Endpoints
Create a user
```json
POST /users
{
  "name": "Hammed"
}

{
  "message": "user created",
  "error": false,
  "data": {
    "ID": 4,
    "Name": "Hammed2",
    "Balance": 29,
    "Verified": false
  }
}
```
Get all users
```json
GET /users
{
  "message": "users",
  "error": false,
  "data": [
    {
      "ID": 1,
      "Name": "Hammed2",
      "Balance": 30,
      "Verified": true
    }
  ]
}
```
Create a Transaction/send money from one user to another
```json
POST /transactions
{
  "userId": 2,
  "receiverId": 1,
  "amount": 10
}

{
  "message": "transaction created",
  "error": false,
  "data": {
    "ID": 4,
    "Amount": 10,
    "UserId": 2,
    "ReceiverId": 1,
    "Status": "CREATED"
  }
}
```