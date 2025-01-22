import {BounceLoader} from "react-spinners";
import MonacoEditor from "@monaco-editor/react";
import CodeOutput from "./CodeOutput.tsx";

const CodeEditor = ({
    currentSession,
    isEditorLoading,
    showDialog,
    output,
    onMount,
}) => {
    return (
        <div className="flex-grow flex flex-col relative">
            {/* Loader Overlay */}
            {(showDialog || isEditorLoading) && (
                <div
                    className="absolute inset-0 bg-gray-900 bg-opacity-75 flex items-center justify-center z-10">
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
                                    <button className={`${action.color} mr-3 text-white px-4 py-2 rounded-md`}
                                            key={action.label} onClick={action.handler}>
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
                    onMount={onMount}
                    theme={'vs-dark'}
                    options={{
                        autoIndent: 'full',
                        contextmenu: true,
                        fontFamily: 'poppins',
                        fontSize: 15,
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
                />
            </div>

            {/* Output Pane */}
            <CodeOutput output={output}/>
        </div>
    );
};

export default CodeEditor