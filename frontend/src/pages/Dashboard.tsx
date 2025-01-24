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
import AddSessionDialog from "../components/AddSessionDialog.tsx";
import * as usersSlice from "../features/users/users.slice.ts";
import LocalStorage from "../services/storage.service.ts";
import * as authSlice from "../features/auth/auth.slice.ts";

const Dashboard = () => {
    const [showDialog, setShowDialog] = useState(false)

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
            toast.success(codeSessionState.message || 'Created session successfully', {
                position: 'top-right',
                className: "text-base"
            })
            const { currentSession } = codeSessionState?.data
            dispatch(sessionSlice.clearState())
            navigate(`/sessions/${currentSession?.sessionId}?new=true`)
        }
    }, [codeSessionState.status]);

    useEffect(() => {
        dispatch<any>(dashboardSlice.fetchDashboard())
        dispatch<any>(usersSlice.fetchUsers())
    }, []);

    const createNewSession = (language, title, collaborators) => {
        console.log(language, title, collaborators)
        const data: CreateSessionDTO = {
            language,
            title,
            collaboratorIds: collaborators,
        }

        dispatch<any>(sessionSlice.createCodeSession(data))
    }

    const navigateToSession = (sessionId) => {
        navigate(`/sessions/${sessionId}`)
    }

    const handleLogout = () => {
        dispatch<any>(authSlice.clear())
        navigate('/login')
    }

    if (!dashboardAnalytics.status || dashboardAnalytics.status === DashboardActionStatus.FETCH_DASHBOARD_IN_PROGRESS) {
        return <LoadingPage />
    }

    if (dashboardAnalytics?.status === DashboardActionStatus.FETCH_DASHBOARD_FAILED) {
        return <ErrorDialog error={dashboardAnalytics.message} retryHandler={() => {}} />
    }

    return (
        <div className="dashboard-container flex flex-col bg-gray-100 p-6">
            {/* Header */}
            <header className="flex items-center justify-between mb-6">
                <h1 className="text-3xl font-bold text-blue-600">Welcome Back, {user?.username}!</h1>
                <div className="flex items-center gap-2">
                    <button
                        className="bg-blue-600 text-white px-4 py-2 rounded-lg shadow-md hover:bg-blue-700"
                        onClick={() => setShowDialog(true)}
                    >
                        + Create New Session
                    </button>

                    <button
                        className="bg-red-600 text-white px-4 py-2 rounded-lg shadow-md hover:bg-red-700"
                        onClick={handleLogout}
                    >
                        Logout
                    </button>
                </div>
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
                                    {session?.collaborators?.length || 0}
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
                            <span>{collaborator.user.username}</span>
                            <span className="text-sm text-gray-500">
                            Last Active: {new Date(collaborator.lastActive).toLocaleString()}
                        </span>
                        </li>
                    ))}
                </ul>
            </section>

            {/* Dialog for New Session */}
            <AddSessionDialog
                showDialog={showDialog}
                setShowDialog={setShowDialog}
                codeSessionState={codeSessionState}
                createNewSession={createNewSession}
            />
        </div>
    );
}

export default Dashboard