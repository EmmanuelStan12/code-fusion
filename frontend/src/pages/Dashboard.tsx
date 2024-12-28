import {useState} from "react";
import '../App.css'

const Dashboard = () => {
    const [showDialog, setShowDialog] = useState(false)

    const collaborators = [1, 2, 3, 4, 5, 6].map(i => ({
        lastActive: '0',
        id: i,
        username: `Random-collaborator-${i}`
    }))

    const sessions = [1, 2, 3, 4, 5, 6].map(i => ({
        isActive: i % 2 === 0,
        id: i,
        title: `Random-${i}`,
        language: 'JavaScript',
        collaborators: collaborators,
        memoryLimit: 256,
    }))

    const user = {
        totalCodeTime: 100,
        name: 'ByteBard',
        activeSessions: sessions,
        collaborators: collaborators,
        sessionsPerLanguage: [1, 2, 4].map(i => ({
            language: 'JavaScript',
            count: i,
        })),
        sessions,
    }

    const navigateToSession = (sessionId) => {

    }

    const createNewSession = () => {

    }

    return (
        <div className="dashboard-container flex flex-col bg-gray-100 p-6">
            {/* Header */}
            <header className="flex items-center justify-between mb-6">
                <h1 className="text-3xl font-bold text-blue-600">Welcome Back, {user.name}!</h1>
                <button
                    className="bg-blue-600 text-white px-4 py-2 rounded-lg shadow-md hover:bg-blue-700"
                    onClick={() => setShowDialog(true)}
                >
                    + Create New Session
                </button>
            </header>

            {/* Main Content */}
            <div className="flex flex-wrap gap-6">
                {/* Overview Section */}
                <section className="flex-grow bg-white rounded-lg shadow-md p-6">
                    <h2 className="text-xl font-semibold text-gray-800 mb-4">Your Coding Activity</h2>
                    <div className="grid grid-cols-2 gap-6">
                        {/* Total Coding Time Card */}
                        <div className="bg-blue-100 rounded-lg p-6 text-center flex items-center justify-center">
                            <div className="flex items-center justify-center">
                                <i className="fas fa-clock text-blue-500 text-3xl mr-3"></i>
                                <div>
                    <span className="block text-2xl font-bold text-blue-500">
                        {user.totalCodeTime} hrs
                    </span>
                                    <span className="text-gray-600">Total Coding Time</span>
                                </div>
                            </div>
                        </div>

                        {/* Active Sessions Card */}
                        <div className="bg-green-100 rounded-lg p-6 text-center flex items-center justify-center">
                            <div className="flex items-center justify-center">
                                <i className="fas fa-play-circle text-green-500 text-3xl mr-3"></i>
                                <div>
                    <span className="block text-2xl font-bold text-green-500">
                        {user.activeSessions.length}
                    </span>
                                    <span className="text-gray-600">Active Sessions</span>
                                </div>
                            </div>
                        </div>

                        {/* Collaborators Card */}
                        <div className="bg-yellow-100 rounded-lg p-6 text-center flex items-center justify-center">
                            <div className="flex items-center justify-center">
                                <i className="fas fa-users text-yellow-500 text-3xl mr-3"></i>
                                <div>
                    <span className="block text-2xl font-bold text-yellow-500">
                        {user.collaborators.length}
                    </span>
                                    <span className="text-gray-600">Collaborators</span>
                                </div>
                            </div>
                        </div>

                        {/* Languages Used Card */}
                        <div className="bg-red-100 rounded-lg p-6 text-center flex items-center justify-center">
                            <div className="flex items-center justify-center">
                                <i className="fas fa-code text-red-500 text-3xl mr-3"></i>
                                <div>
                    <span className="block text-2xl font-bold text-red-500">
                        {user.sessionsPerLanguage.length}
                    </span>
                                    <span className="text-gray-600">Languages Used</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </section>


                {/* Language Sessions Section */}
                <section className="flex-grow bg-white rounded-lg shadow-md p-6">
                    <h2 className="text-xl font-semibold text-gray-800 mb-4">Sessions by Language</h2>
                    <ul className="space-y-2">
                        {user.sessionsPerLanguage.map((lang) => (
                            <li
                                key={lang.language}
                                className="flex items-center justify-between p-3 rounded-lg bg-gray-50 shadow-inner"
                            >
                                <span className="text-gray-800 font-medium">{lang.language}</span>
                                <span className="bg-blue-100 text-blue-600 px-3 py-1 rounded-full text-sm font-bold">
                                {lang.count} Sessions
                            </span>
                            </li>
                        ))}
                    </ul>
                </section>
            </div>

            {/* Sessions Section */}
            <section className="mt-6 bg-white rounded-lg shadow-md p-6">
                <h2 className="text-xl font-semibold text-gray-800 mb-4">Recent Sessions</h2>
                <ul className="space-y-4">
                    {user.sessions.map((session) => (
                        <li
                            key={session.id}
                            className="flex items-center justify-between p-4 bg-gray-50 rounded-lg shadow-sm"
                        >
                            <div>
                                <h3 className="text-lg font-bold text-gray-900">{session.title}</h3>
                                <p className="text-sm text-gray-600">
                                    Language: {session.language} | Collaborators:{" "}
                                    {session.collaborators.length} | Status:{" "}
                                    <span
                                        className={`font-bold ${
                                            session.isActive
                                                ? "text-green-600"
                                                : "text-red-600"
                                        }`}
                                    >
                                    {session.isActive ? "Active" : "Inactive"}
                                </span>
                                </p>
                            </div>
                            <button
                                className="bg-blue-600 text-white px-4 py-2 rounded-lg shadow-md hover:bg-blue-700"
                                onClick={() => navigateToSession(session.id)}
                            >
                                Open Session
                            </button>
                        </li>
                    ))}
                </ul>
            </section>

            {/* Collaborators Section */}
            <section className="mt-6 bg-white rounded-lg shadow-md p-6">
                <h2 className="text-xl font-semibold text-gray-800 mb-4">Recent Collaborators</h2>
                <ul className="space-y-3">
                    {user.collaborators.map((collaborator) => (
                        <li
                            key={collaborator.id}
                            className="flex items-center justify-between p-3 bg-gray-50 rounded-lg shadow-inner"
                        >
                            <span>{collaborator.username}</span>
                            <span className="text-sm text-gray-500">
                            Last Active: {collaborator.lastActive}
                        </span>
                        </li>
                    ))}
                </ul>
            </section>

            {/* Dialog for New Session */}
            {showDialog && (
                <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
                    <div className="bg-white rounded-lg p-6 w-[90%] max-w-lg shadow-lg">
                        <h2 className="text-2xl font-semibold text-gray-800 mb-4">
                            Create a New Session
                        </h2>
                        <div className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Session Title
                                </label>
                                <input
                                    type="text"
                                    name="title"
                                    placeholder="Enter session title"
                                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Programming Language
                                </label>
                                <select
                                    name="language"
                                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                >
                                    <option value="JavaScript">JavaScript</option>
                                    <option value="Python">Python</option>
                                    <option value="TypeScript">TypeScript</option>
                                    <option value="Go">Go</option>
                                </select>
                            </div>
                            <div className="flex gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Memory Limit (MB)
                                    </label>
                                    <input
                                        type="number"
                                        name="memoryLimit"
                                        className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        CPU Limit (Cores)
                                    </label>
                                    <input
                                        type="number"
                                        name="cpuLimit"
                                        className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    />
                                </div>
                            </div>
                        </div>
                        <div className="flex justify-end mt-6">
                            <button
                                className="bg-gray-200 text-gray-700 px-4 py-2 rounded-md mr-4 hover:bg-gray-300"
                                onClick={() => setShowDialog(false)}
                            >
                                Cancel
                            </button>
                            <button
                                className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
                                onClick={() => {
                                }}
                            >
                                Create Session
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

export default Dashboard