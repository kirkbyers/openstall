# Openstall

A distributed system for monitoring when stalls are open.

## Roadmap

- rPi GPIO
- Work on user UI
- Makefile

## Structs

### WS message

```
{
    "fromId": string,
    "fromName": string,
    "type": string,
    "payload": {}
}
```

### Registration

```
{
    "id": string,
    "name": string,
    "type": string,
}
```
