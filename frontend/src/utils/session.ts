export interface ExecuteCodeRequest {
    code: string
    sessionId: string
}

export interface ExecuteCodeResponse {
    result?: any
    error?: any
    success?: boolean
    stdout?: string
}

export interface WebSocketRequestMessage<T> {
    messageType: string
    data?: T
}

export interface WebSocketResponseMessage<T> {
    action: string
    data: T
}

export interface CodeOperation {
    type?: 'insert' | 'delete'
    position: number
    text?: string
    length?: number
}

export enum CodeSessionSocketActions {
    SESSION_INITIALIZED = "SESSION_INITIALIZED",
    SESSION_ERROR = "SESSION_ERROR",
    SESSION_CLOSED = "SESSION_CLOSED",
    CODE_EXECUTION_SUCCESS = "CODE_EXECUTION_SUCCESS",
    CODE_EXECUTION_FAILED = "CODE_EXECUTION_FAILED",
    ACTION_CODE_EXECUTION = "CODE_EXECUTION",
    COLLABORATOR_INACTIVE = "COLLABORATOR_INACTIVE",
    COLLABORATOR_ACTIVE = "COLLABORATOR_ACTIVE",
    ACTION_CODE_UPDATE = "CODE_UPDATE",
}
