# RE Challenge

Application designed for the RE Partners interview.


## Description
This app, simulates a production line shipment software, user can ask for orders that are going to be set depending on the following rules:
1. Only whole packs can be sent. Packs cannot be broken open.
2. Within the constraints of Rule 1 above, send out the least amount of items to fulfil the order.
3. Within the constraints of Rules 1 & 2 above, send out as few packs as possible to fulfil each order.

## Local Deployment

```
docker-compose build
docker-compose up
```

## Deployed App
The deployment was made using Heroku, can be found here:
- https://re-challenge-28272d1c79fc.herokuapp.com/