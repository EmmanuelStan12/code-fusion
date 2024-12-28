const path = require("node:path");
module.exports = {
    defaultIsolateConfig: {
        memoryLimit: 128
    },
    defaultContextConfig: {},
    defaultCodeConfig: {
        timeout: 60 * 1000,
    },
    protoBufConfig: {
        codeExecutionProto: path.resolve(__dirname, '../proto/code_execution.proto')
    }
}
