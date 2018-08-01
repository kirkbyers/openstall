# Openstall

A distributed system for monitoring when stalls are open.

## Roadmap

- Websockets
    - Server impl
    - Client impl
- Client registration
    - http request
- rPi GPIO
- Bolt DB impl for saving state
- End user UI

## Structs

### WS message

```json
{
    "fromId": string,
    "fromName": string,
    "type": string,
    "payload": {}
}
```

### Registration

```json
{
    "id": string,
    "name": string,
    "type": string,
}
```
