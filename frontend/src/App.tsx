import './App.css'
import {store} from "./redux/store.ts";
import {BrowserRouter, Route, Routes} from "react-router";
import Home from "./pages/Home.tsx";
import Login from "./pages/Login.tsx";
import Register from "./pages/Register.tsx";
import Dashboard from "./pages/Dashboard.tsx";
import {ToastContainer} from "react-toastify";
import ProtectedRoute from "./components/ProtectedRoute.tsx";
import {useEffect} from "react";
import {useAppDispatch, useAppSelector} from "./redux/hooks.ts";
import LocalStorage, {AUTH_TOKEN_KEY} from "./services/storage.service.ts";
import * as authSlice from './features/auth/auth.slice.ts'
import {AuthActionStatus} from "./features/auth/auth.actions.ts";
import LoadingPage from "./components/LoadingPage.tsx";
import AuthRoute from "./components/AuthRoute.tsx";
import CodeSession from "./pages/CodeSession.tsx";

function App() {
    const state = useAppSelector(state => state.auth)
    const dispatch = useAppDispatch()
    const status = state.status || AuthActionStatus.AUTHENTICATION_IN_PROGRESS

    useEffect(() => {
        if (LocalStorage.get(AUTH_TOKEN_KEY) && !state?.data?.user) {
            dispatch<any>(authSlice.getAuthUser())
        }
    }, []);

    if (status === AuthActionStatus.AUTHENTICATION_IN_PROGRESS && LocalStorage.get(AUTH_TOKEN_KEY)) {
        return <LoadingPage />
    }

    return (
        <>
            <BrowserRouter>
                <Routes>
                    <Route index element={<Home/>}/>
                    <Route element={<AuthRoute />}>
                        <Route path="login" element={<Login/>}/>
                        <Route path="register" element={<Register/>}/>
                    </Route>

                    <Route element={<ProtectedRoute/>}>
                        <Route path="dashboard" element={<Dashboard/>}/>
                        <Route path="sessions/:sessionId" element={<CodeSession />} />
                    </Route>

                    {/*<Route path="concerts">
                            <Route index element={<ConcertsHome />} />
                            <Route path=":city" element={<City />} />
                            <Route path="trending" element={<Trending />} />
                        </Route>*/}
                </Routes>
            </BrowserRouter>
            <ToastContainer/>
        </>
    )
}

export default App
