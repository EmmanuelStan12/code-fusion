const CodeSession = () => {
    return (
        <div class="code-session flex flex-col h-screen bg-gray-900 text-gray-100">
            {/* Header */}
            <header class="flex items-center justify-between p-4 bg-gray-800 shadow-lg">
                <h1 class="text-lg font-bold text-blue-400">{session.title}</h1>
                <div class="flex items-center gap-4">
                    <button
                        class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
                        onClick={runCode}
                    >
                        Run
                    </button>
                    <button
                        class="bg-gray-700 text-gray-300 px-4 py-2 rounded-lg hover:bg-gray-600"
                        onClick={lockCollaborators}
                    >
                        Lock Collaborators
                    </button>
                    <button
                        class="bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700"
                        onClick={deleteSession}
                    >
                        Delete Session
                    </button>
                </div>
            </header>

            {/* Main Content */}
            <div class="flex flex-grow overflow-hidden">
                {/* Code Editor */}
                <div class="flex-grow flex flex-col">
                    <div class="flex-grow">
                        <MonacoEditor
                            height="100%"
                            language={session.language}
                            value={session.code}
                            options={{
                                theme: 'vs-dark',
                                fontSize: 14,
                                minimap: { enabled: false },
                            }}
                            onChange={updateCode}
                        />
                    </div>
                    {/* Output Pane */}
                    <div class="output-pane bg-gray-800 p-4 text-sm">
                        <h2 class="text-blue-400 font-bold">Output</h2>
                        <pre class="bg-gray-900 text-gray-200 p-4 rounded-lg mt-2 overflow-auto">
                        {output || 'No output yet. Run your code to see results.'}
                    </pre>
                        <div class="metrics flex gap-5 mt-4">
                            <div class="text-gray-300">Memory: {metrics.memory} MB</div>
                            <div class="text-gray-300">CPU: {metrics.cpu}%</div>
                            <div class="text-gray-300">Time: {metrics.time}ms</div>
                        </div>
                    </div>
                </div>

                {/* Collaborators Sidebar */}
                <aside class="w-72 bg-gray-800 p-4 border-l border-gray-700">
                    <h2 class="text-blue-400 font-bold">Collaborators</h2>
                    <ul class="mt-4">
                        {collaborators.map((collab) => (
                            <li key={collab.id} class="flex justify-between items-center p-2 bg-gray-700 rounded-lg mt-2">
                                <span>{collab.name}</span>
                                <button
                                    class="bg-red-600 text-white px-2 py-1 rounded hover:bg-red-700"
                                    onClick={() => removeCollaborator(collab.id)}
                                >
                                    Remove
                                </button>
                            </li>
                        ))}
                    </ul>
                    <button
                        class="bg-blue-600 text-white px-4 py-2 rounded-lg w-full mt-4 hover:bg-blue-700"
                        onClick={addCollaborator}
                    >
                        Add Collaborator
                    </button>
                </aside>
            </div>
        </div>
    );
}

export default CodeSession