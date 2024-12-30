const {CodeFusionIsolate} = require("../isolate.service");
const {defaultIsolateConfig, defaultContextConfig, defaultCodeConfig} = require("../../config/config");
const {preProcessCode, executeCode} = require("../code.service");
const Logger = require('../../utils/logger.utils')
const CodeExecutionService = {}

CodeExecutionService.executeCode = function (call) {
    call.on('data', function (codeRequest) {
        Logger.instance.info('Data request ', codeRequest)
        const { code, sessionId, contextId } = codeRequest;

        const isolate = CodeFusionIsolate.instance({ sessionId, config: defaultIsolateConfig });

        isolate.createContext(contextId, defaultContextConfig).then(contextProps => {
            const codeProps = preProcessCode(code);

            executeCode(codeProps, contextProps, defaultCodeConfig).then(result => {
                const data = {
                    ...result,
                    contextId: contextProps.contextId,
                    sessionId,
                }
                Logger.instance.info(`Executed code with result: `, data)
                call.write(data)
            }).catch(err => {
                Logger.instance.error('Error executing code:', err);
                call.write({
                    result: null,
                    stdout: '',
                    error: 'Something went wrong..., Please try again later.',
                    sessionId,
                })
            }).finally(() => {
                isolate.disposeContext(contextProps.contextId)
            })
        }).catch(err => {
            Logger.instance.error('Error creating isolate context:', err);
        });
        call.on('end', function () {
            CodeFusionIsolate.dispose(sessionId)
            call.end();
        })
    })
}

module.exports = CodeExecutionService