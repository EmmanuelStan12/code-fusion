const {stdout} = require('node:process')

class Logger {
    static instance = new Logger([stdout]);

    constructor(streams = []) {
        this.streams = streams;
    }

    #writeToStreams(level, message, ...args) {
        message = this.#formatMessage(level, message, ...args)
        this.streams.forEach(stream => {
            stream.write(message);
        });
    }

    log(message, ...args) {
        this.#writeToStreams('LOG', message, ...args)
    }

    info(message, ...args) {
        this.#writeToStreams('INFO', message, ...args)
    }

    warn(message, ...args) {
        this.#writeToStreams('WARN', message, ...args)
    }

    error(message, ...args) {
        this.#writeToStreams('ERROR', message, ...args)
    }

    #safeJSONStringify(data) {
        if (!data) return ''
        try {
            return JSON.stringify(data, null, 2)
        } catch {
            return data.toString ? data.toString() : ''
        }
    }

    #formatMessage(level, message, ...args) {
        const timestamp = new Date().toISOString()
        let metaData = ''
        for (const arg of args) {
            if (arg instanceof Error) {
                metaData += `${arg.toString()}\n${arg.stack}`
                continue
            }
            metaData += `${this.#safeJSONStringify(arg)}\n`
        }
        return `[${timestamp}] [${level.toUpperCase()}]: ${message} ${metaData}\n`
    }
}

module.exports = Logger
