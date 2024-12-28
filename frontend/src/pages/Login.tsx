import '../App.css'
import {Link, useNavigate} from "react-router";
import {useEffect, useRef, useState} from "react";
import {useAppDispatch, useAppSelector} from "../redux/hooks.ts";
import {AuthActionStatus} from "../features/auth/auth.actions.ts";
import {toast} from "react-toastify";
import {LoginDTO, RegisterDTO} from "../features/auth/auth.api.ts";
import * as authSlice from "../features/auth/auth.slice.ts";
import {BounceLoader} from "react-spinners";

function Login() {

    const emailRef = useRef<any>()
    const passwordRef = useRef<any>()
    const dispatch = useAppDispatch()
    const state = useAppSelector(state => state.auth)
    const [isValid, setIsValid] = useState(false)
    const navigate = useNavigate()

    useEffect(() => {
        if (state.status === AuthActionStatus.LOGIN_FAILED) {
            toast.error(state.message, {
                position: 'top-right',
                className: "text-base"
            })
        }
        if (state.status === AuthActionStatus.LOGIN_SUCCESSFUL) {
            toast.success(state.message, {
                position: 'top-right',
                className: "text-base"
            })
            navigate("/dashboard")
        }
    }, [state.status]);

    const login = () => {
        const data: LoginDTO = {
            password: passwordRef?.current?.value || '',
            email: emailRef?.current?.value || '',
        }

        dispatch<any>(authSlice.login(data))
    }

    const validateForm = () => {
        setIsValid(passwordRef?.current?.value &&
            emailRef?.current?.value
        )
    }

    return (
        <div className="page-login">
            <div className="w-[80%] bg-white p-10 flex items-stretch rounded-[15px] justify-center shadow-gray-500">
                <div className="flex bg-[rgba(0,0,0,0.05)] rounded-[15px] flex-grow box p-10 flex-1 items-start flex-col justify-center">
                    <h2 className="text-2xl font-bold">Code Fusion</h2>
                    <div className="text-gray-700 text-base mt-2">
                        Write. Run. Collaborate.
                    </div>
                    <p className="text-gray-600 mt-4">
                        Empower your coding journey with real-time collaboration and seamless execution of JavaScript, TypeScript, and Python.
                    </p>
                </div>
                <form className="login-form flex-1 p-10" onChange={(e) => validateForm()} onSubmit={(e) => e.preventDefault()}>
                    <div>
                        <h2 className="text-2xl font-bold">Welcome Back</h2>
                        <p className="text-xs mt-2 text-gray-700">
                            Log in to unleash your creativity and build amazing things.
                        </p>
                    </div>

                    <div className="flex flex-col gap-2 mt-5">
                        <span className="block text-sm text-gray-700">Email</span>
                        <input
                            className="rounded-md border border-gray-300 focus:outline-none focus:ring-0 focus:border-gray-300 px-3 py-2"
                            type="email"
                            placeholder="Enter your email"
                            name="email"
                            ref={emailRef}
                        />
                    </div>

                    <div className="flex flex-col gap-2 mt-5">
                        <span className="block text-sm text-gray-700">Password</span>
                        <input
                            className="rounded-md border border-gray-300 focus:outline-none focus:ring-0 focus:border-gray-300 px-3 py-2"
                            type="password"
                            placeholder="Enter your password"
                            name="password"
                            ref={passwordRef}
                        />
                    </div>

                    {state.status === AuthActionStatus.LOGIN_IN_PROGRESS ? (
                        <div className="flex justify-center items-center mt-7">
                            <BounceLoader
                                color="rgb(37, 99, 235)"
                                loading
                                size={40}
                            />
                        </div>
                    ) : (
                        <button
                            onClick={login}
                            disabled={!isValid}
                            className={`${!isValid ? 'bg-blue-100' : 'bg-blue-600'} text-white bold w-[100%] mt-7 p-3 rounded-md`}>Login
                        </button>
                    )}

                    <div className="text-sm text-gray-700 mt-10 text-center">
                        Don't have an account? <Link className="text-blue-700" to="/register">Sign up</Link>
                    </div>
                </form>
            </div>
        </div>
    )
}

export default Login
