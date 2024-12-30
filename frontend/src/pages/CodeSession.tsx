import {useAppDispatch, useAppSelector} from "../redux/hooks.ts";
import MonacoEditor from '@monaco-editor/react';
import {useEffect, useRef, useState} from "react";
import LoadingPage from "../components/LoadingPage.tsx";
import {CodeSessionActionStatus} from "../features/code-session/session.actions.ts";
import {fetchCodeSession} from "../features/code-session/session.slice.ts";
import {Navigate, useParams} from "react-router";
import {toast} from "react-toastify";
import IOSocket, {ExecuteCodeRequestDTO} from "../services/socket.service.ts";

interface CodeExecuteResponse {
    action: string
    data: any
}

const SESSION_INITIALIZED = "SESSION_INITIALIZED"

const CodeSession = () => {
    const codeSessionState = useAppSelector(state => state.codeSession)
    const dispatch = useAppDispatch()
    const { sessionId } = useParams()

    const currentSession = codeSessionState?.data?.currentSession
    const currentStatus = codeSessionState.status || CodeSessionActionStatus.FETCH_CODE_SESSION_IN_PROGRESS
    const collaborators = []
    const [isConnected, setIsConnected] = useState(false)
    const [code, setCode] = useState('')
    const [socket, setSocket] = useState<any>(null)

    const initSocket = () => {
        const socket = new IOSocket(sessionId)
        socket.registerListener('message', (event) => {
            const response = JSON.parse(event.data) as CodeExecuteResponse
            if (response.action === SESSION_INITIALIZED) {
                setIsConnected(true)
                setCode(response.data.code)
            }
        })
        setSocket(socket)
    }

    useEffect(() => {
        console.log('State: ', sessionId)
        if (!currentSession && sessionId) {
            dispatch<any>(fetchCodeSession(sessionId))
        } else {
            initSocket()
        }

        return () => {
            socket?.dispose()
        }
    }, []);

    useEffect(() => {
        if (codeSessionState.status === CodeSessionActionStatus.FETCH_CODE_SESSION_FAILED) {
            toast.error(codeSessionState.message, {
                position: 'top-right'
            })
        }
        if (codeSessionState.status === CodeSessionActionStatus.FETCH_CODE_SESSION_SUCCESSFUL) {
            initSocket()
        }
    }, [codeSessionState.status]);

    const runCode = () => {
        const codeRequest: ExecuteCodeRequestDTO = {
            code,
            sessionId: sessionId as string,
        }
        console.log(codeRequest)
        socket?.send(codeRequest)
    }

    const lockCollaborators = () => {

    }

    const deleteSession = () => {

    }

    const updateCode = (code) => {
        setCode(code)
    }

    if (!sessionId) {
        return <Navigate to="/dashboard" />;
    }

    if (currentStatus === CodeSessionActionStatus.FETCH_CODE_SESSION_IN_PROGRESS || !isConnected) {
        return <LoadingPage />
    }

    return (
        <div className="code-session flex flex-col h-screen bg-gray-900 text-gray-100">
            {/* Header */}
            <header className="flex items-center justify-between p-4 bg-gray-800 shadow-lg">
                <h1 className="text-lg font-bold text-blue-400">{currentSession.title}</h1>
                <div className="flex items-center gap-4">
                    <button
                        className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
                        onClick={runCode}
                    >
                        Run
                    </button>
                    <button
                        className="bg-gray-700 text-gray-300 px-4 py-2 rounded-lg hover:bg-gray-600"
                        onClick={lockCollaborators}
                    >
                        Lock Collaborators
                    </button>
                    <button
                        className="bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700"
                        onClick={deleteSession}
                    >
                        Delete Session
                    </button>
                </div>
            </header>

            {/* Main Content */}
            <div className="flex flex-grow overflow-hidden">
                {/* Code Editor */}
                <div className="flex-grow flex flex-col">
                    <div className="flex-grow">
                        <MonacoEditor
                            height="100%"
                            language={currentSession.language?.toLowerCase()}
                            value={code}
                            options={{
                                theme: 'vs-dark',
                                fontSize: 14,
                                minimap: { enabled: false },
                            }}
                            onChange={updateCode}
                        />
                    </div>
                    {/* Output Pane */}
                    <div className="output-pane bg-gray-800 p-4 text-sm">
                        <h2 className="text-blue-400 font-bold">Output</h2>
                        <pre className="bg-gray-900 text-gray-200 p-4 rounded-lg mt-2 overflow-auto">
                        {'' || 'No output yet. Run your code to see results.'}
                    </pre>
                        <div className="metrics flex gap-5 mt-4">
                            <div className="text-gray-300">Memory: 9 MB</div>
                            <div className="text-gray-300">Duration: 100s</div>
                        </div>
                    </div>
                </div>

                {/* Collaborators Sidebar */}
                <aside className="w-72 bg-gray-800 p-4 border-l border-gray-700">
                    <h2 className="text-blue-400 font-bold">Collaborators</h2>
                    <ul className="mt-4">
                        {collaborators.map((collaborator) => (
                            <li key={collaborator.id} className="flex justify-between items-center p-2 bg-gray-700 rounded-lg mt-2">
                                <span>{collaborator.name}</span>
                                <button
                                    className="bg-red-600 text-white px-2 py-1 rounded hover:bg-red-700"
                                    onClick={() => {}}
                                >
                                    Remove
                                </button>
                            </li>
                        ))}
                    </ul>
                    <button
                        className="bg-blue-600 text-white px-4 py-2 rounded-lg w-full mt-4 hover:bg-blue-700"
                        onClick={() => {}}
                    >
                        Add Collaborator
                    </button>
                </aside>
            </div>
        </div>
    );
}

export default CodeSession