const {randomUUID} = require("crypto");
const {ContextOutputStream} = require("../lib/stream.lib");
const Logger = require('../utils/logger.utils')

const startupScript = `
    global.console = {
        log: function(...args) {
            $0.applyIgnored(undefined, args, { arguments: { copy: true } });
        }
    }
`

function preProcessCode(input) {
    const codeId = randomUUID()
    return {
        codeId,
        code: input,
    }
}

async function executeCode(codeProps, contextProps, options = { timeout: 60 * 1000 }) {
    const { context, contextId } = contextProps
    const outputStream = new ContextOutputStream(contextId)
    const global = context.global

    const consoleRef = function (...args) {
        const message = args.join(' ')
        outputStream.write(message)
    }

    await global.set('global', global.derefInto())

    await context.evalClosure(
        startupScript,
        [consoleRef],
        { arguments: { reference: true } }
    )

    const { code } = codeProps
    let result, error = ''
    try {
        result = await context.eval(code, options)
    } catch (e) {
        error = e.message || e.stack
    }
    return {
        success: !error,
        result,
        stdout: outputStream.output,
        error,
    }
}

module.exports = {
    preProcessCode,
    executeCode,
}
