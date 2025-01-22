import {useEffect, useRef, useState} from "react";
import LoadingPage from "../components/LoadingPage.tsx";
import {Navigate, useNavigate, useParams} from "react-router";
import IOSocket from "../services/socket.service.ts";
import useDebounce from "../components/hooks/debounce.hook.ts";
import CollaboratorsSidebar from "../components/CollaboratorsSidebar.tsx";
import CodeSessionHeader from "../components/CodeSessionHeader.tsx";
import {
    CodeOperation,
    CodeSessionSocketActions,
    ExecuteCodeRequest,
    ExecuteCodeResponse,
    WebSocketRequestMessage, WebSocketResponseMessage
} from "../utils/session.ts";
import CodeEditor from "../components/CodeEditor.tsx";
import {mergeCode} from "../utils/merger.ts";

const CodeSession = () => {
    const {sessionId} = useParams()

    const [currentSession, setCurrentSession] = useState<any>({})
    const [socket, setSocket] = useState<any>(null)
    const [output, setOutput] = useState<ExecuteCodeResponse>({})
    const [isEditorLoading, setIsEditorLoading] = useState(true)
    const [isCodeRunning, setIsCodeRunning] = useState(false)
    const [isConnected, setIsConnected] = useState(false)
    const [showDialog, setShowDialog] = useState<{ message: string, actions: any[] } | null>(null)
    const editorRef = useRef<any>(null)
    const [prevCode, setPrevCode] = useState('')
    const isProgrammaticChange = useRef(false);
    const timerRef = useRef<number | null>(null);

    const [collaborators, setCollaborators] = useState([])
    const navigate = useNavigate()

    useEffect(() => {
        initSocket()

        return () => {
            if (socket) {
                socket.dispose()
            }
        }
    }, []);

    const updateCode = (code: string) => {
        const codeRequest: WebSocketRequestMessage<{ code: string }> = {
            data: { code },
            messageType: CodeSessionSocketActions.ACTION_CODE_UPDATE,
        }
        socket?.send(codeRequest)
        setPrevCode(code)
    }

    const handleChange = (value: string) => {
        console.log(timerRef.current);
        if (timerRef.current) {
            clearTimeout(timerRef.current as number);
        }

        timerRef.current = setTimeout(() => {
            updateCode(value);

            timerRef.current = null;
        }, 3000);
    }

    const closeSession = (route?: string) => {
        socket?.dispose()
        if (route) {
            navigate(route)
        }
    }

    const updateCollaboratorStatus = (updatedCollaborator) => {
        setCollaborators((prevCollaborators) =>
            prevCollaborators.map((collaborator) =>
                collaborator.id === updatedCollaborator.id ? updatedCollaborator : collaborator
            )
        );
    };

    const handleSessionInitialized = (data) => {
        setCollaborators(data.collaborators);
        setCurrentSession(data);
        setPrevCode(data.code)

        handleProgrammaticChange(data.code)
        setIsEditorLoading(false);
        setIsConnected(true);
        setShowDialog(null);
    }

    const handleSessionError = () => {
        setIsConnected(false);
        setShowDialog({
            message: "Initialization failed. Retry or go back to the dashboard.",
            actions: [
                { label: "Retry", color: 'bg-red-600', handler: initSocket },
                { label: "Dashboard", color: 'bg-blue-600', handler: () => closeSession("/dashboard") },
            ],
        });
        setIsEditorLoading(false);
    };

    const handleCodeExecutionSuccess = (data: ExecuteCodeResponse) => {
        setIsCodeRunning(false);
        setOutput(data);
    };

    const handleCodeExecutionFailed = () => {
        setIsCodeRunning(false);
        setShowDialog({
            message: "Code execution failed. Please try again.",
            actions: [
                { label: "Close", color: 'bg-red-600', handler: () => setShowDialog(null) },
            ],
        });
    };

    const handleSessionClosed = () => {
        setIsConnected(false);
        setShowDialog({
            message: "Connection closed unexpectedly. Retry or save session.",
            actions: [
                { label: "Retry", color: 'bg-red-600', handler: initSocket },
                { label: "Dashboard", color: 'bg-blue-600', handler: () => closeSession("/dashboard") },
            ],
        });
    };

    const handleProgrammaticChange = (value: string) => {
        const editor = editorRef.current;
        if (!editor) return;

        const model = editor.getModel();
        isProgrammaticChange.current = true
        if (model) {
            model.pushEditOperations(
                [],
                [
                    {
                        range: model.getFullModelRange(),
                        text: value,
                    },
                ],
                () => null
            );
        }
    }

    const handleCodeUpdate = (data) => {
        const { code } = data
        const editor = editorRef.current;
        if (!editor) return;

        const currentCode = editor.getValue()

        const diffCode = mergeCode(currentCode, code)
        setPrevCode(diffCode);
        handleProgrammaticChange(diffCode)
    };

    const initSocket = () => {
        setIsEditorLoading(true)
        const socket = new IOSocket(sessionId)
        socket.registerListener('message', (event) => {
            const response = JSON.parse(event.data) as WebSocketResponseMessage<any>
            switch (response.action) {
                case CodeSessionSocketActions.SESSION_INITIALIZED:
                    handleSessionInitialized(response.data)
                    break
                case CodeSessionSocketActions.SESSION_ERROR:
                    handleSessionError()
                    break
                case CodeSessionSocketActions.CODE_EXECUTION_SUCCESS:
                    handleCodeExecutionSuccess(response.data)
                    break
                case CodeSessionSocketActions.CODE_EXECUTION_FAILED:
                    handleCodeExecutionFailed()
                    break
                case CodeSessionSocketActions.SESSION_CLOSED:
                    handleSessionClosed()
                    break
                case CodeSessionSocketActions.ACTION_CODE_UPDATE:
                    handleCodeUpdate(response.data)
                    break
                case CodeSessionSocketActions.COLLABORATOR_ACTIVE:
                case CodeSessionSocketActions.COLLABORATOR_INACTIVE:
                    updateCollaboratorStatus(response.data)
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
                    {label: "Dashboard", color: 'bg-blue-600', handler: () => closeSession("/dashboard")},
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
                    {label: "Dashboard", color: 'bg-blue-600', handler: () => closeSession("/dashboard")},
                ]
            })
            setIsEditorLoading(false)
            setIsConnected(false)
        })
        setSocket(socket)
    }

    function handleEditorMount(editor, monaco) {
        editor.onDidChangeModelContent((event) => {
            if (isProgrammaticChange.current) {
                isProgrammaticChange.current = false
                return
            }
            const value = editor.getValue();
            handleChange(value)
        });
        editorRef.current = editor
        const model = editor.getModel()
        if (model) {
            model.setValue(prevCode)
        }
        console.log(editor)
    }

    const runCode = () => {
        const code = editorRef.current?.getValue()
        const codeRequest: WebSocketRequestMessage<ExecuteCodeRequest> = {
            data: {
                code,
                sessionId: sessionId as string,
            },
            messageType: CodeSessionSocketActions.ACTION_CODE_EXECUTION,
        }
        setIsCodeRunning(true)
        socket?.send(codeRequest)
    }

    if (!sessionId) {
        return <Navigate to="/dashboard"/>;
    }

    if (!isConnected && isEditorLoading) {
        return <LoadingPage/>
    }

    return (
        <div className="code-session flex h-screen bg-gray-900 text-gray-100">
            <div className="flex flex-1 flex-col bg-gray-900 text-gray-100">
                {/* Header */}
                <CodeSessionHeader
                    currentSession={currentSession}
                    isCodeRunning={isCodeRunning}
                    isEditorLoading={isEditorLoading}
                    runCode={runCode}
                    closeSession={closeSession}
                />

                {/* Main Content */}
                <CodeEditor
                    currentSession={currentSession}
                    isEditorLoading={isEditorLoading}
                    showDialog={showDialog}
                    output={output}
                    onMount={handleEditorMount}
                />
            </div>

            {/* Sidebar */}
            <CollaboratorsSidebar
                collaborators={collaborators}
            />
        </div>
    );
}

export default CodeSession