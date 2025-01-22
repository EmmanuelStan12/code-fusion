const {CodeFusionIsolate} = require("../isolate.service");
const {defaultIsolateConfig, defaultContextConfig, defaultCodeConfig} = require("../../config/config");
const {preProcessCode, executeCode} = require("../code.service");
const Logger = require('../../utils/logger.utils')
const ts = require('typescript')
const CodeExecutionService = {}

const transpileTypescript = (code) => {
    Logger.instance.info('Transpiling TypeScript code...');
    const transpiledCode = ts.transpileModule(code, {
        compilerOptions: { module: ts.ModuleKind.CommonJS },
    }).outputText
    Logger.instance.info('TypeScript transpiled successfully.');
    return transpiledCode
}

CodeExecutionService.closeSession = function (call) {
    const sessionId = call.request.sessionId;
    try {
        if (!sessionId) {
            throw new Error('Missing sessionId');
        }

        CodeFusionIsolate.dispose(sessionId);
        call.write({ success: true, sessionId });
    } catch (error) {
        Logger.instance.error('Failed to close session:', error);
        call.write({ success: false, sessionId, error: error.message });
    } finally {
        call.end();
    }
}

CodeExecutionService.executeCode = function (call) {
    call.on('data', function (codeRequest) {
        Logger.instance.info('Data request ', codeRequest)
        let { code, sessionId, contextId, language = 'javascript' } = codeRequest;

        if (language.toLowerCase() === 'typescript') {
            try {
                code = transpileTypescript(code)
            } catch (e) {
                call.write({
                    result: null,
                    stdout: '',
                    error: `Cannot transpile Typescript, ${e.message}`,
                    sessionId,
                    contextId,
                })
                return
            }
        }
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
    })
    call.on('end', function () {
        CodeFusionIsolate.disposeIsolates()
        call.end()
    })
}

module.exports = CodeExecutionService