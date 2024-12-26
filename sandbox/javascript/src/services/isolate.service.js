const ivm = require("isolated-vm");
const {randomUUID} = require("crypto");
const {ContextOutputStream} = require("../lib/stream.lib");

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

    async createContext(options) {
        this._contextCount++
        const context = await this.isolate.createContext(options)
        const contextId = randomUUID()
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
        for (const isolate of Object.values(this._instances)) {
            isolate.dispose()
        }
    }

    static dispose(sessionId) {
        const isolate = this._instances[sessionId]
        if (isolate) {
            isolate.dispose()
        }
    }
}

module.exports = {
    CodeFusionIsolate
}