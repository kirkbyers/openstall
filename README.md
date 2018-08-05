# Openstall

A distributed system for monitoring when stalls are open.

## Roadmap

- rPi GPIO
- Work on user UI
- Makefile

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
