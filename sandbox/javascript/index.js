const readline = require('readline')
const {CodeFusionIsolate} = require("./src/services/isolate.service");
const {randomUUID} = require("crypto");
const {defaultIsolateConfig, defaultContextConfig, defaultCodeConfig} = require("./src/config/config");
const {preProcessCode, executeCode} = require("./src/services/code.service");
const { stdin, stdout } = require('node:process')

function initReader() {
    let code = '';
    const reader = readline.createInterface({
        input: process.stdin,
        output: process.stdout,
        prompt: 'Enter code (type "END" to execute or "QUIT" to exit):\n'
    });

    reader.prompt(); // Start prompt

    reader.on('line', line => {
        if (line.trim() === 'END') {
            // Code execution logic
            const sessionId = randomUUID();
            const isolate = CodeFusionIsolate.instance({ sessionId, config: defaultIsolateConfig });

            isolate.createContext(defaultContextConfig).then(contextProps => {
                const codeProps = preProcessCode(code); // Process the code if necessary

                executeCode(codeProps, contextProps, defaultCodeConfig).then(result => {
                    console.log(result);
                }).catch(err => {
                    console.error('Error executing code:', err);
                }).finally(() => {
                    isolate.disposeContext(contextProps.contextId)
                    code = ''
                    reader.prompt()
                })
            }).catch(err => {
                console.error('Error creating isolate context:', err);
                code = ''
                reader.prompt()
            });
        } else if (line.trim() === "QUIT") {
            reader.close()
            CodeFusionIsolate.disposeIsolates()
        } else {
            code += `${line}\n`;
        }
    });
}

(function () {
    initReader()
})();
