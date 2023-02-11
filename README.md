## Example GO Gin MongoDB

list api in this repo


POST: `/player`
```json
{
    "name": "Lionel Messi",
    "region": "Argentina",
    "position": "FW"
}
```

GET: `/player/:playerId`
response example: 
```json
{
    "status": 200,
    "message": "success",
    "data": {
        "data": {
            "id": "63e7550798c6c4a177c9bb69",
            "name": "Lionel Messi",
            "region": "Argentina",
            "position": "FW"
        }
    }
}
```

PUT: `/player/:playerId`
```json
{
    "name": "Lionel Messi",
    "region": "Argentina",
    "position": "FW"
}
```

DELETE: `/player/:playerId`
response example: 
```json
{
    "status": 200,
    "message": "success",
    "data": {
        "data": "Player successfully deleted!"
    }
}
```

GET: `/players`
response example: 
```json
{
    "status": 200,
    "message": "players",
    "data": {
        "data": [
            {
                "id": "63e754db98c6c4a177c9bb66",
                "name": "Enzo Fernández",
                "region": "Argentina",
                "position": "MF"
            },
            {
                "id": "63e754f098c6c4a177c9bb67",
                "name": "Julián Álvarez",
                "region": "Argentina",
                "position": "FW"
            },
            {
                "id": "63e754fd98c6c4a177c9bb68",
                "name": "Emiliano Martínez",
                "region": "Argentina",
                "position": "GK"
            },
            {
                "id": "63e7550798c6c4a177c9bb69",
                "name": "Lionel Messi",
                "region": "Argentina",
                "position": "FW"
            }
        ]
    }
}
```
