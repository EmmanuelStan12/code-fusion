export default function CodeOutput({ output }) {
    return (
        <div className="output-pane bg-gray-800 p-4 text-sm mt-4">
            <h2 className="text-blue-400 font-bold mb-4">Output</h2>
            <div className="bg-gray-900 p-4 rounded-lg text-gray-200">
                <h3 className="text-green-400 font-bold mb-2">Result</h3>
                <pre
                    className="bg-gray-800 p-2 rounded-lg text-gray-200">
                            {output.result || 'nil'}
                        </pre>
            </div>
            <div className="bg-gray-900 p-4 rounded-lg text-gray-200 mt-4">
                <h3 className="text-yellow-400 font-bold mb-2">Stdout</h3>
                <pre className="bg-gray-800 p-2 rounded-lg text-gray-200">
                            {output.stdout || 'nil'}
                        </pre>
            </div>
            <div className="bg-gray-900 p-4 rounded-lg text-gray-200 mt-4">
                <h3 className="text-red-400 font-bold mb-2">Error</h3>
                <pre className="bg-gray-800 p-2 rounded-lg text-gray-200">
                            {output.error || 'nil'}
                        </pre>
            </div>
        </div>
    )
}