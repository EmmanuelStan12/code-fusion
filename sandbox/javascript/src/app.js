const grpc = require('@grpc/grpc-js')
const protoLoader = require('@grpc/proto-loader')
const { protoBufConfig } = require('./config/config')
const CodeExecutionService = require('./services/proto/code-execution.service')
const Logger = require('./utils/logger.utils')

const loadProtoFile = (filePath) => {
    const packageDefinition = protoLoader.loadSync(
        filePath, {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        }
    )
    return grpc.loadPackageDefinition(packageDefinition)
}

class App {
    constructor(port) {
        this.port = port
        this.server = new grpc.Server()

        const codeExecutionProto = protoBufConfig.codeExecutionProto

        const descriptor = loadProtoFile(codeExecutionProto)

        this.server.addService(descriptor.CodeExecutionService.service, CodeExecutionService)
    }

    listen() {
        this.server.bindAsync(`0.0.0.0:${this.port}`, grpc.ServerCredentials.createInsecure(), (err, port) => {
            if (err) {
                Logger.instance.error(err)
            } else {
                Logger.instance.info(`gRPC server listening on PORT: ${port}`)
            }
        })
    }
}

module.exports = App