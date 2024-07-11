const ws = new WebSocket('ws://localhost:8838/gochat');

ws.onopen = () => {
    console.log('connected');
}

ws.onmessage = (e) => {
    console.log('Received message', e.data);
    const message = JSON.stringify({
        datetime: (new Date()).getTime(),
    });

    ws.send(message);
}

ws.onclose = () => {
    console.log('closed');
}

ws.onerror = (e) => {
    console.log('error', e);
}