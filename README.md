# RE Challenge

Application designed for the RE Partners interview.


## Description
This app, simulates a production line shipment software, user can ask for orders that are going to be set depending on the following rules:
1. Only whole packs can be sent. Packs cannot be broken open.
2. Within the constraints of Rule 1 above, send out the least amount of items to fulfil the order.
3. Within the constraints of Rules 1 & 2 above, send out as few packs as possible to fulfil each order.

## Used libraries
I used the next libraries for, Testing, Web Framework and SQL Driver:

```
go get -u github.com/gin-gonic/gin
go get github.com/jackc/pgx/v5
```

## Tests
Tests where made just for the orders business class, logic to generate orders based on the 3 rules given for the challenge
```
go test ./...
```

## Deployed App
The deployment was made using Heroku, can be found here:
- https://re-challenge-28272d1c79fc.herokuapp.com/

## Local Deployment
A docker compose basic configuration was made to set up a workig environment with some dummy data, database configuration can be found on ```/docker/env/Postgres.env``` and dummy generated data can be found on ```/docker/postgres-init/Postgres.sql```

To run the API, just use:
```
docker-compose build
docker-compose up
```

## Endpoints
### Packs
```GET``` ```/api/v1/pack```
Get all available Pack sizes

```
[
    {
        "created": "2024-01-01T00:00:00.000000Z",
        "id": 1,
        "size": 250
    },
    ...
]
```

```POST``` ```/api/v1/pack/create?size=NN```
Create a new pack size for the orders

```
{
    "pack_id": 1
}
```

```DELETE``` ```/api/v1/pack/delete?pack_id=NN```
Delete a pack by it's ID

```
{
    "response": "Done"
}
```

### Orders
```GET``` ```/api/v1/order```
Get all the created orders

```
[
    {
        "created": "2024-01-01T00:00:00.000000Z",
        "deleted": false,
        "id": 1,
        "modified": "2024-01-01T00:00:00.000000Z",
        "packs": [
            {
                "amount": 1,
                "id": 1,
                "size": 250
            }
        ],
        "size": 1
    },
    ...
]
```

```POST``` ```/api/v1/order/create?order_size=NN```
Create a new order

```
{
    "created": "2024-01-01T00:00:00.000000Z",
    "deleted": false,
    "id": 1,
    "modified": "2024-01-01T00:00:00.000000Z",
    "packs": [
        {
            "amount": 1,
            "id": 1,
            "size": 250
        }
    ],
    "size": 1
}
```

```DELETE``` ```/api/v1/order/cancel?order_id=NN```
Delete an order by it's ID

```
{
    "response": "Done"
}
```
