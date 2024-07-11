const ws = new WebSocket('ws://localhost:8838/gochat?Authorization=asdfawefiqasdfkasdfasdfasfdn');

const login = 1;

ws.onopen = () => {
    console.log('connected');
}

ws.onmessage = (e) => {
    console.log('Received message', e.data);
    const message = JSON.stringify({
        id:1,
        data: {
            userName: "asdfwfs",
            pwd: "123123123123123",
            datetime: (new Date()).getTime(),
        },
    });

    ws.send(message);
}

ws.onclose = () => {
    console.log('closed');
}

ws.onerror = (e) => {
    console.log('error', e);
}