const {Writable} = require("stream");

class ContextOutputStream extends Writable {
    constructor(name, options) {
        super(options);
        this.name = name;
        this._output = '';
    }

    append(chunk) {
        this._output += `${chunk}\n`
    }

    _write(chunk, encoding, callback) {
        this.append(chunk.toString().trim())
        callback()
    }

    get output() {
        return this._output
    }
}

module.exports = {
    ContextOutputStream,
}