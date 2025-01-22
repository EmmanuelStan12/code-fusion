import {BounceLoader} from "react-spinners";

const CodeSessionHeader = ({
    currentSession,
    isCodeRunning,
    isEditorLoading,
    runCode,
    closeSession
}) => {
    return <header className="flex items-center justify-between p-4 bg-gray-800 shadow-lg">
        <div>
            <h1 className="text-lg font-bold text-blue-400">{currentSession?.title}</h1>
            <span className={"text-xs poppins-medium-italic text-gray-300"}>({currentSession?.language})</span>
        </div>
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
                className={`${isCodeRunning || isEditorLoading ? 'bg-green-100' : 'bg-green-600'} text-white px-4 py-2 rounded-lg hover:bg-green-700`}
                onClick={() => closeSession('/dashboard')}
                disabled={isCodeRunning || isEditorLoading}
            >
                Go to Dashboard
            </button>
        </div>
    </header>;
}

export default CodeSessionHeader