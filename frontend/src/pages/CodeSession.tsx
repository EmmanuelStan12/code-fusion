import {useAppDispatch, useAppSelector} from "../redux/hooks.ts";
import MonacoEditor from '@monaco-editor/react';
import {useEffect, useRef, useState} from "react";
import LoadingPage from "../components/LoadingPage.tsx";
import {CodeSessionActionStatus} from "../features/code-session/session.actions.ts";
import {fetchCodeSession} from "../features/code-session/session.slice.ts";
import {Navigate, useNavigate, useParams} from "react-router";
import {toast} from "react-toastify";
import IOSocket from "../services/socket.service.ts";
import {BounceLoader} from "react-spinners";
import * as usersSlice from "../features/users/users.slice.ts";
import {UserActionStatus} from "../features/users/users.actions.ts";
import AddCollaboratorDialog from "../components/AddCollaboratorDialog.tsx";

interface ExecuteCodeRequest {
    code: string
    sessionId: string
}

interface ExecuteCodeResponse {
    result?: any
    error?: any
    success?: boolean
    stdout?: string
}

interface WebSocketRequestMessage<T> {
    messageType: string
    userId?: number
    data?: T
}

interface WebSocketResponseMessage<T> {
    action: string
    data: T
}

const SESSION_INITIALIZED = "SESSION_INITIALIZED"
const SESSION_ERROR = "SESSION_ERROR";
const SESSION_CLOSED = "SESSION_CLOSED";
const CODE_EXECUTION_SUCCESS = "CODE_EXECUTION_SUCCESS";
const CODE_EXECUTION_FAILED = "CODE_EXECUTION_FAILED";
const ACTION_CODE_EXECUTION = "CODE_EXECUTION";
const ACTION_ADD_COLLABORATOR = "ADD_COLLABORATOR";
const ACTION_CODE_UPDATE = "CODE_UPDATE";

const CodeSession = () => {
    const codeSessionState = useAppSelector(state => state.codeSession)
    const dispatch = useAppDispatch()
    const {sessionId} = useParams()

    const [currentSession, setCurrentSession] = useState<any>({})
    const [code, setCode] = useState('')
    const [socket, setSocket] = useState<any>(null)
    const [output, setOutput] = useState<ExecuteCodeResponse>({})
    const [isEditorLoading, setIsEditorLoading] = useState(true)
    const [isCodeRunning, setIsCodeRunning] = useState(false)
    const [isConnected, setIsConnected] = useState(false)
    const [showDialog, setShowDialog] = useState<{ message: string, actions: any[] } | null>(null)

    const [showCollaboratorDialog, setShowCollaboratorDialog] = useState(false)
    const [collaborators, setCollaborators] = useState([])
    const [collaboratorStatus, setCollaboratorStatus] = useState(false)
    const usersState = useAppSelector(state => state.users)
    const navigate = useNavigate()


    const initSocket = () => {
        setIsEditorLoading(true)
        const socket = new IOSocket(sessionId)
        socket.registerListener('message', (event) => {
            const response = JSON.parse(event.data) as WebSocketResponseMessage<any>
            switch (response.action) {
                case SESSION_INITIALIZED:
                    setCode(response.data.code)
                    setCollaborators(response.data.collaborators)
                    setCurrentSession(response.data)
                    setIsEditorLoading(false)
                    setIsConnected(true)
                    setShowDialog(null)
                    break
                case SESSION_ERROR:
                    setIsConnected(false)
                    setShowDialog({
                        message: "Initialization failed. Retry or go back to the dashboard.",
                        actions: [
                            {label: "Retry", color: 'bg-red-600', handler: initSocket},
                            {label: "Dashboard", color: 'bg-blue-600', handler: () => navigate("/dashboard")},
                        ]
                    })
                    setIsEditorLoading(false)
                    break
                case CODE_EXECUTION_SUCCESS:
                    setIsCodeRunning(false)
                    setOutput(response.data as ExecuteCodeResponse)
                    break
                case CODE_EXECUTION_FAILED:
                    setIsCodeRunning(false)
                    setShowDialog({
                        message: "Code execution failed. Please try again.",
                        actions: [{label: "Close", color: 'bg-red-600', handler: () => setShowDialog(null)}]
                    })
                    break
                case SESSION_CLOSED:
                    setIsConnected(false)
                    setShowDialog({
                        message: "Connection closed unexpectedly. Retry or save session.",
                        actions: [
                            {label: "Retry", color: 'bg-red-600', handler: initSocket},
                            {label: "Save Session", color: 'bg-green-600', handler: () => console.log("Session Saved")},
                            {label: "Dashboard", color: 'bg-blue-600', handler: () => navigate("/dashboard")},
                        ]
                    })
                    break
                default:
                    console.warn("Unhandled action:", response.action)
            }
        })
        socket.registerListener('close', () => {
            setShowDialog({
                message: "Connection closed unexpectedly. Retry or save session.",
                actions: [
                    {label: "Retry", color: 'bg-red-600', handler: initSocket},
                    {label: "Save Session", color: 'bg-green-600', handler: () => console.log("Session Saved")},
                    {label: "Dashboard", color: 'bg-blue-600', handler: () => navigate("/dashboard")},
                ]
            })
            setIsEditorLoading(false)
            setIsConnected(false)
        })

        socket.registerListener('error', () => {
            setShowDialog({
                message: "Connection cannot be established. Retry",
                actions: [
                    {label: "Retry", color: 'bg-red-600', handler: initSocket},
                    {label: "Save Session", color: 'bg-green-600', handler: () => console.log("Session Saved")},
                    {label: "Dashboard", color: 'bg-blue-600', handler: () => navigate("/dashboard")},
                ]
            })
            setIsEditorLoading(false)
            setIsConnected(false)
        })
        setSocket(socket)
    }

    useEffect(() => {
        initSocket()
        dispatch<any>(usersSlice.fetchUsers())

        return () => {
            console.log('Disposing...')
            socket?.dispose()
        }
    }, []);

    const runCode = () => {
        const codeRequest: WebSocketRequestMessage<ExecuteCodeRequest> = {
            data: {
                code,
                sessionId: sessionId as string,
            },
            messageType: ACTION_CODE_EXECUTION,
        }
        setIsCodeRunning(true)
        socket?.send(codeRequest)
    }

    const addCollaborator = () => {
        setCollaboratorStatus(true)
    }

    const lockCollaborators = () => {

    }

    const deleteSession = () => {

    }

    const updateCode = (code) => {
        setCode(code)
        const codeRequest: WebSocketRequestMessage<ExecuteCodeRequest> = {
            data: {
                code,
                sessionId: sessionId as string,
            },
            messageType: ACTION_CODE_UPDATE,
        }
        socket?.send(codeRequest)
    }

    if (!sessionId) {
        return <Navigate to="/dashboard"/>;
    }

    if (!isConnected) {
        return <LoadingPage/>
    }

    return (
        <div className="code-session flex h-screen bg-gray-900 text-gray-100">
            <div className="flex flex-1 flex-col bg-gray-900 text-gray-100">
                {/* Header */}
                <header className="flex items-center justify-between p-4 bg-gray-800 shadow-lg">
                    <h1 className="text-lg font-bold text-blue-400">{currentSession?.title}</h1>
                    <div className="flex items-center gap-4">
                        {isCodeRunning || isEditorLoading ? (
                            <div>
                                <BounceLoader
                                    color="rgb(37, 99, 235)"
                                    loading
                                    size={40}
                                />
                            </div>
                        ) : (
                            <button
                                className={`${isCodeRunning || isEditorLoading ? 'bg-blue-100' : 'bg-blue-600'} text-white px-4 py-2 rounded-lg hover:bg-blue-700`}
                                onClick={runCode}
                                disabled={isCodeRunning || isEditorLoading}
                            >
                                Run
                            </button>
                        )}
                        <button
                            className="bg-gray-700 text-gray-300 px-4 py-2 rounded-lg hover:bg-gray-600"
                            onClick={lockCollaborators}
                        >
                            Lock Collaborators
                        </button>
                        <button
                            className={`${isCodeRunning || isEditorLoading ? 'bg-red-100' : 'bg-red-600'} text-white px-4 py-2 rounded-lg hover:bg-red-700`}
                            onClick={deleteSession}
                            disabled={isCodeRunning || isEditorLoading}
                        >
                            Delete Session
                        </button>
                    </div>
                </header>

                {/* Main Content */}
                <div className="flex-grow flex flex-col relative">
                    {/* Loader Overlay */}
                    {(showDialog || isEditorLoading) && (
                        <div className="absolute inset-0 bg-gray-900 bg-opacity-75 flex items-center justify-center z-10">
                            {isEditorLoading ? (
                                <BounceLoader
                                    color="rgb(37, 99, 235)"
                                    loading
                                    size={80}
                                />
                            ) : (
                                <div className={"bg-white rounded-lg p-6 w-[90%] max-w-lg shadow-lg"}>
                                    <p className={'text-md text-gray-800 mb-4'}>{showDialog?.message}</p>
                                    <div className="actions">
                                        {showDialog?.actions.map(action => (
                                            <button className="bg-gray-600 mr-3 text-white px-4 py-2 rounded-md" key={action.label} onClick={action.handler}>
                                                {action.label}
                                            </button>
                                        ))}
                                    </div>
                                </div>
                            )}
                        </div>
                    )}

                    {/* Monaco Editor */}
                    <div className="flex-grow">
                        <MonacoEditor
                            height="100%"
                            language={currentSession.language?.toLowerCase()}
                            value={code}
                            theme={'vs-dark'}
                            options={{
                                autoIndent: 'full',
                                contextmenu: true,
                                fontFamily: 'monospace',
                                fontSize: 13,
                                lineHeight: 24,
                                hideCursorInOverviewRuler: true,
                                matchBrackets: 'always',
                                minimap: {
                                    enabled: true,
                                },
                                scrollbar: {
                                    horizontalSliderSize: 4,
                                    verticalSliderSize: 18,
                                },
                                selectOnLineNumbers: true,
                                roundedSelection: false,
                                readOnly: false,
                                cursorStyle: 'line',
                                automaticLayout: true,
                            }}
                            onChange={updateCode}
                        />
                    </div>

                    {/* Output Pane */}
                    <div className="output-pane bg-gray-800 p-4 text-sm mt-4">
                        <h2 className="text-blue-400 font-bold mb-4">Output</h2>
                        <div className="bg-gray-900 p-4 rounded-lg text-gray-200">
                            <h3 className="text-green-400 font-bold mb-2">Result</h3>
                            <pre
                                className="bg-gray-800 p-2 rounded-lg text-gray-200">
                            {output.result || '--------------------------------'}
                        </pre>
                        </div>
                        <div className="bg-gray-900 p-4 rounded-lg text-gray-200 mt-4">
                            <h3 className="text-yellow-400 font-bold mb-2">Stdout</h3>
                            <pre className="bg-gray-800 p-2 rounded-lg text-gray-200">
                            {output.stdout || '--------------------------------'}
                        </pre>
                        </div>
                        <div className="bg-gray-900 p-4 rounded-lg text-gray-200 mt-4">
                            <h3 className="text-red-400 font-bold mb-2">Error</h3>
                            <pre className="bg-gray-800 p-2 rounded-lg text-gray-200">
                            {output.error || '----------------------------------'}
                        </pre>
                        </div>
                    </div>

                </div>
            </div>

            {/* Sidebar */}
            <aside className="w-80 bg-gray-800 p-4 border-l border-gray-700 flex flex-col">
                <div className="flex justify-between items-center mb-4">
                    <h2 className="text-blue-400 font-bold">Collaborators</h2>
                    <button
                        className="bg-blue-600 text-white px-3 py-2 rounded-lg hover:bg-blue-700"
                        onClick={() => setShowCollaboratorDialog(true)}
                    >
                        + Add
                    </button>
                </div>
                <div className="flex-grow overflow-y-auto">
                    {collaborators.length === 0 ? (
                        <div className="flex justify-center items-center">
                            <p>No collaborators</p>
                        </div>
                    ) : (
                        collaborators.map((collab) => (
                            <div
                                key={collab.id}
                                className="bg-gray-700 text-gray-300 p-2 rounded-lg mb-2"
                            >
                                {collab.user.username}
                            </div>
                        ))
                    )}
                </div>
            </aside>

            {/* Add Collaborator Dialog */}
            {usersState?.data?.users && showCollaboratorDialog && (
                <AddCollaboratorDialog
                    users={usersState.data.users}
                    setShowDialog={setShowCollaboratorDialog}
                    addCollaborator={addCollaborator}
                    isLoading={collaboratorStatus}
                />
            )}
        </div>
    );
}

export default CodeSession