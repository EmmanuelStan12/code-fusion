import {io} from "socket.io-client";
import {WS_API_BASE_URL} from "../config/env.config.ts";
import LocalStorage, {AUTH_TOKEN_KEY} from "./storage.service.ts";

class IOSocket {
    private socket: any;
    private eventListeners: { [key: string]: Function[] } = {
        'open': [() => console.log('WebSocket connection established')],
        'close': [(event) => console.log('WebSocket connection closed:', event)],
        'error': [(event) => console.error('WebSocket error: ', event)],
        'message': [(event) => console.log('Message received:', event.data)],
    };
    constructor(sessionId) {
        this.socket = new WebSocket(`${WS_API_BASE_URL}/api/v1/sessions/init/${sessionId}?token=${LocalStorage.get(AUTH_TOKEN_KEY)}`)

        this.socket.onmessage = (event) => this.dispatchEvent('message', event);
        this.socket.onopen = (event) => this.dispatchEvent('open', event);
        this.socket.onerror = (event) => this.dispatchEvent('error', event);
        this.socket.onclose = (event) => this.dispatchEvent('close', event);
    }

    registerListener(event: string, fn: Function) {
        if (!this.eventListeners[event]) {
            this.eventListeners[event] = [];
        }
        this.eventListeners[event].push(fn);
    }

    private dispatchEvent(event: string, eventObject: Event) {
        const listeners = this.eventListeners[event];
        if (listeners) {
            listeners.forEach((listener) => listener(eventObject));
        }
    }

    send(data) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify(data));
        } else {
            console.error('WebSocket is not open');
        }
    }

    dispose() {
        if (this.socket) {
            this.socket.close();
        }
    }
}

export default IOSocket