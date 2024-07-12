class Client {
    opts= {
        retryInterval: 30*1000,
    };
    static ON_ERROR="error";
    static ON_OPEN="open";
    static ON_CLOSE="close";
    constructor(baseUrl, authorization, opts) {
        this.baseUrl = baseUrl;
        this.authorization = authorization;
        this.opts = Object.assign({}, this.opts, opts);
        this.events = {};
        this.timer = null;
        this.retryCount = 0;
    }

    on(eventName, listener) {
        if (!this.events[eventName]) {
            this.events[eventName] = [];
        }
        this.events[eventName].push(listener);
    }

    emit(eventName, data) {
        const eventListeners = this.events[eventName];
        if (eventListeners) {
            eventListeners.forEach(listener => {
                listener(data);
            });
        }
    }

    off(eventName, listener) {
        const eventListeners = this.events[eventName];
        if (eventListeners) {
            this.events[eventName] = eventListeners.filter(fn => fn !== listener);
        }
    }

    connect = () => {
        this.ws = null;
        this.ws = new WebSocket(`${this.baseUrl}?Authorization=${this.authorization}`);
        this.ws.onopen = this.onOpen;
        this.ws.onmessage = this.onMessage;
        this.ws.onclose = this.onClose;
        this.ws.onerror = this.onError;
    }

    send = (id, data) => {
        const message = JSON.stringify({
            id:id,
            data: data,
        });

        this.ws.send(this.pack(message));
    }

    getReadyState (){
        return this.ws.readyState
    }

    pack = (data) => {
        return data
    }

    unpack = (data) => {
        return JSON.parse(data)
    }

    onOpen = () => {
        clearInterval(this.timer);
        this.emit(Client.ON_OPEN)
    }

    onMessage = (e) => {
        console.log(`received message: ${e.data}`);
        const msg = this.unpack(e.data);
        this.emit(msg.id, msg.data);
    }

    onError = (e) => {
        this.emit(Client.ON_ERROR, e)
    }

    onClose = () => {
        clearInterval(this.timer)
        this.timer = setInterval(() => {
            this.retryCount++;
            console.log('retry', this.retryCount);

            this.connect()

        }, this.opts.retryInterval);

        this.emit(Client.ON_CLOSE)
    }
}

class V1 {
    constructor(pubKey) {
        this.encrypt = new JSEncrypt();
        this.encrypt.setPublicKey(pubKey);

        const randomKey = () => {
            return (new Date()).getTime().toString()+(Math.floor(Math.random() * 900) + 100).toString();
        }

        this.key = CryptoJS.enc.Utf8.parse(randomKey());
        this.iv = CryptoJS.enc.Utf8.parse(randomKey());
        this.auth = encodeURIComponent(this.replace4Encode(this.encrypt.encrypt(`${CryptoJS.enc.Utf8.stringify(this.key)}${CryptoJS.enc.Utf8.stringify(this.iv)}`)))
    }

    replace4Encode = (data) => {
        return data.replace(/\//g, '_').replace(/\+/g, '-')
    }

    replace4Decode = (data) => {
        return data.replace(/_/g, '/').replace(/-/g, '+')
    }

    pack = (data) => {
        const encrypted = CryptoJS.AES.encrypt(data, this.key, {
            iv: this.iv,
            mode: CryptoJS.mode.CBC,
            padding: CryptoJS.pad.Pkcs7,
        });

        return this.replace4Encode(encrypted.ciphertext.toString(CryptoJS.enc.Base64));
    }

    unpack = (data) => {
        const decrypted = CryptoJS.AES.decrypt(this.replace4Decode(data), this.key, {
            iv: this.iv,
            mode: CryptoJS.mode.CBC,
            padding: CryptoJS.pad.Pkcs7,
        });
        data = decrypted.toString(CryptoJS.enc.Utf8);
        return JSON.parse(data)
    }
}

const publicKeyPem = `-----BEGIN PUBLIC KEY-----
    MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAgEfXy8MpQ7mPt29aF38g
VTaYbUGJMfk7veY/cWKl5XSalroTNfQWidq5oskw83xK1PI15lPlVHdOUqIONHhr
jj8zfaq8i7YjOvNac3UixJfRCbhfaj52NieaHLMdzYHfZZjMKOkDspar1Ur/b5XW
Fc8y5oiFsxj5YSxPc0mgceYCgGh4AUm0/KL+b6XNmID7qWdQPWuP+33k89Fjq7Np
vAku6ZJ94+knXmxbuSp/BcPI8QEYWJerhMA5sNsCPhRlEtBt96J5+cQ0/ABR94FA
zL8dqezepazpqtT1R6eKLN/QYPSpNfK6k5yxlWmtx2p22J6msB5gI+N9RuTN1wTb
hQIDAQAB
-----END PUBLIC KEY-----`;
const v1 = new V1(publicKeyPem);

const client = new Client(`ws://localhost:8838/gochat`, v1.auth);
client.pack = v1.pack

client.unpack = v1.unpack

client.on(Client.ON_ERROR, (e) => {
    console.log(e);
})

client.on(Client.ON_CLOSE, () => {
    console.log('client closed');
})

client.on(Client.ON_OPEN, () => {
    console.log('client opened');
})

const eventLogin = 1;

client.on(eventLogin, (data) => {
    console.log(`receive data: ${JSON.stringify(data)}`)
})

client.connect();

const userId = Math.ceil(Math.random() * 10);

setInterval(() => {
    client.send(eventLogin, {
        userId: userId,
        userName: "我也是",
        pwd: 'sdf2sdfasf',
    })
}, 3000)


