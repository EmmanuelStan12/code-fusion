const ivm = require("isolated-vm");
const {randomUUID} = require("crypto");

class CodeFusionIsolate {
    static _instances = {}
    constructor(options) {
        const { sessionId, config } = options
        this.sessionId = sessionId
        this.isolateId = randomUUID()
        this.contexts = {}
        this.isolate = new ivm.Isolate(config)
        this._contextCount = 0
    }

    async createContext(contextId, options) {
        this._contextCount++
        const context = await this.isolate.createContext(options)
        this.contexts[contextId] = context

        return {
            context,
            contextId,
        }
    }

    get contextCount() {
        return this._contextCount
    }

    disposeContext(contextId) {
        if (this._contextCount > 0 && this.contexts[contextId]) {
            this._contextCount--
            const context = this.contexts[contextId]
            context.release()
        }
    }

    static instance(options) {
        const { sessionId } = options
        let isolate = this._instances[sessionId]
        if (!isolate) {
            isolate = new CodeFusionIsolate(options)
            this._instances[sessionId] = isolate
        }
        return isolate
    }

    static disposeIsolates() {
        for (const [sessionId, instance] of Object.entries(this._instances)) {
            instance.isolate.dispose()
            delete this._instances[sessionId]
        }
    }

    static dispose(sessionId) {
        const instance = this._instances[sessionId]
        if (instance && !instance.isDisposed) {
            instance.isolate.dispose()
            delete this._instances[sessionId]
        }
    }
}

module.exports = {
    CodeFusionIsolate
}