import {useEffect, useRef, useState} from "react";
import '../App.css'
import {useAppDispatch, useAppSelector} from "../redux/hooks.ts";
import {useNavigate} from "react-router";
import {CodeSessionActionStatus} from "../features/code-session/session.actions.ts";
import {toast} from "react-toastify";
import {BounceLoader} from "react-spinners";
import * as sessionSlice from "../features/code-session/session.slice.ts";
import * as dashboardSlice from "../features/dashboard/dashboard.slice.ts";
import {CreateSessionDTO} from "../features/code-session/session.api.ts";
import {DashboardActionStatus} from "../features/dashboard/dashboard.actions.ts";
import LoadingPage from "../components/LoadingPage.tsx";
import ErrorDialog from "../components/ErrorDialog.tsx";

const Dashboard = () => {
    const [showDialog, setShowDialog] = useState(false)
    const sessionTitleRef = useRef<any>()
    const languageRef = useRef<any>()
    const memoryLimitRef = useRef<any>()
    const timeoutRef = useRef<any>()

    const [isValid, setIsValid] = useState(false)
    const dispatch = useAppDispatch()
    const codeSessionState = useAppSelector(state => state.codeSession)
    const auth = useAppSelector(state => state.auth)
    const user = auth?.data?.user
    const dashboardAnalytics = useAppSelector(state => state.dashboard)
    const dashboardData = dashboardAnalytics?.data
    const navigate = useNavigate()

    useEffect(() => {
        if (codeSessionState.status === CodeSessionActionStatus.CREATE_SESSION_FAILED) {
            toast.error(codeSessionState.message, {
                position: 'top-right',
                className: "text-base"
            })
        }
        if (codeSessionState.status === CodeSessionActionStatus.CREATE_SESSION_SUCCESSFUL) {
            toast.success(codeSessionState.message, {
                position: 'top-right',
                className: "text-base"
            })
            const { currentSession, collaborator } = codeSessionState?.data
            navigate(`/sessions/${currentSession?.sessionId}`)
        }
    }, [codeSessionState.status]);

    useEffect(() => {
        dispatch<any>(dashboardSlice.fetchDashboard())
    }, []);

    const validateForm = () => {
        setIsValid(sessionTitleRef?.current?.value && languageRef?.current?.value && memoryLimitRef?.current?.value && timeoutRef?.current?.value)
    }

    const createNewSession = () => {
        const data: CreateSessionDTO = {
            timeout: Number(timeoutRef.current.value || '0'),
            memoryLimit: Number(memoryLimitRef.current.value || '0'),
            language: languageRef.current?.value,
            title: sessionTitleRef.current?.value
        }

        dispatch<any>(sessionSlice.createCodeSession(data))
    }

    const navigateToSession = (sessionId) => {
        navigate(`/sessions/${sessionId}`)
    }

    if (!dashboardAnalytics.status || dashboardAnalytics.status === DashboardActionStatus.FETCH_DASHBOARD_IN_PROGRESS) {
        return <LoadingPage />
    }

    console.log(dashboardAnalytics, dashboardData)
    if (dashboardAnalytics?.status === DashboardActionStatus.FETCH_DASHBOARD_FAILED) {
        return <ErrorDialog error={dashboardAnalytics.message} retryHandler={() => {}} />
    }

    return (
        <div className="dashboard-container flex flex-col bg-gray-100 p-6">
            {/* Header */}
            <header className="flex items-center justify-between mb-6">
                <h1 className="text-3xl font-bold text-blue-600">Welcome Back, {user?.username}!</h1>
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
                                    <span className="block text-2xl font-bold text-blue-500">{dashboardData.analytics.totalMinutes} Minutes</span>
                                    <span className="text-gray-600">Total Coding Time</span>
                                </div>
                            </div>
                        </div>

                        {/* Active Sessions Card */}
                        <div className="bg-green-100 rounded-lg p-6 text-center flex items-center justify-center">
                            <div className="flex items-center justify-center">
                                <i className="fas fa-play-circle text-green-500 text-3xl mr-3"></i>
                                <div>
                                    <span className="block text-2xl font-bold text-green-500">{dashboardData.analytics.totalSessions}</span>
                                    <span className="text-gray-600">Active Sessions</span>
                                </div>
                            </div>
                        </div>

                        {/* Collaborators Card */}
                        <div className="bg-yellow-100 rounded-lg p-6 text-center flex items-center justify-center">
                            <div className="flex items-center justify-center">
                                <i className="fas fa-users text-yellow-500 text-3xl mr-3"></i>
                                <div>
                                    <span className="block text-2xl font-bold text-yellow-500">{dashboardData.recentCollaborators?.length}</span>
                                    <span className="text-gray-600">Collaborators</span>
                                </div>
                            </div>
                        </div>

                        {/* Languages Used Card */}
                        <div className="bg-red-100 rounded-lg p-6 text-center flex items-center justify-center">
                            <div className="flex items-center justify-center">
                                <i className="fas fa-code text-red-500 text-3xl mr-3"></i>
                                <div>
                                    <span className="block text-2xl font-bold text-red-500">{dashboardData.analytics.totalLanguagesUsed}</span>
                                    <span className="text-gray-600">Languages Used</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </section>
            </div>

            {/* Sessions Section */}
            <section className="mt-6 bg-white rounded-lg shadow-md p-6">
                <h2 className="text-xl font-semibold text-gray-800 mb-4">Recent Sessions</h2>
                <ul className="space-y-4">
                    {dashboardData?.recentSessions?.map((session) => (
                        <li
                            key={session.id}
                            className="flex items-center justify-between p-4 bg-gray-50 rounded-lg shadow-sm"
                        >
                            <div>
                                <h3 className="text-lg font-bold text-gray-900">{session.title}</h3>
                                <p className="text-sm text-gray-600">
                                    Language: {session.language} | Collaborators:{" "}
                                    {session?.collaborators?.length || 0} | Status:{" "}
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
                                onClick={() => navigateToSession(session.sessionId)}
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
                    {dashboardData?.recentCollaborators?.map((collaborator) => (
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
                    <form onChange={validateForm} onSubmit={(e) => e.preventDefault()} className="bg-white rounded-lg p-6 w-[90%] max-w-lg shadow-lg">
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
                                    ref={sessionTitleRef}
                                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Programming Language
                                </label>
                                <select
                                    name="language"
                                    ref={languageRef}
                                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                >
                                    <option value="JavaScript">JavaScript</option>
                                    <option value="Python">Python</option>
                                    <option value="TypeScript">TypeScript</option>
                                    <option value="Go">Go</option>
                                </select>
                            </div>
                            <div className="flex gap-4">
                                <div className="flex-1">
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Memory Limit (MB)
                                    </label>
                                    <select
                                        name="memoryLimit"
                                        ref={memoryLimitRef}
                                        className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    >
                                        <option value="8">8</option>
                                        <option value="16">16</option>
                                        <option value="32">32</option>
                                        <option value="64">64</option>
                                    </select>
                                </div>
                                <div className="flex-1">
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Timeout (Seconds)
                                    </label>
                                    <select
                                        name="timeout"
                                        ref={timeoutRef}
                                        className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    >
                                        <option value="30">30</option>
                                        <option value="60">60</option>
                                        <option value="120">120</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div className="flex justify-end mt-6">
                            <button
                                className="bg-gray-200 text-gray-700 px-4 py-2 rounded-md mr-4 hover:bg-gray-300"
                                onClick={() => setShowDialog(false)}
                                disabled={codeSessionState.status === CodeSessionActionStatus.CREATE_SESSION_IN_PROGRESS}
                            >
                                Cancel
                            </button>

                            {codeSessionState.status === CodeSessionActionStatus.CREATE_SESSION_IN_PROGRESS ? (
                                <div className="flex justify-center items-center mt-7">
                                    <BounceLoader
                                        color="rgb(37, 99, 235)"
                                        loading
                                        size={40}
                                    />
                                </div>
                            ) : (
                                <button
                                    className={`${!isValid ? 'bg-blue-100' : 'bg-blue-600'} text-white px-4 py-2 rounded-md`}
                                    onClick={createNewSession}
                                    disabled={!isValid}
                                >
                                    Create Session
                                </button>
                            )}
                        </div>
                    </form>
                </div>
            )}
        </div>
    );
}

export default Dashboard