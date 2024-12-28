import '../App.css'
import {Link, useNavigate} from "react-router";
import {useEffect, useRef, useState} from "react";
import {useAppDispatch, useAppSelector} from "../redux/hooks.ts";
import {ToastContainer, toast} from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import {RegisterDTO} from "../features/auth/auth.api.ts";
import {AuthActionStatus} from "../features/auth/auth.actions.ts";
import {CircleLoader, BounceLoader} from "react-spinners";
import * as authSlice from '../features/auth/auth.slice.ts'

function Register() {

    const firstNameRef = useRef<any>()
    const lastNameRef = useRef<any>()
    const usernameRef = useRef<any>()
    const passwordRef = useRef<any>()
    const emailRef = useRef<any>()
    const dispatch = useAppDispatch()
    const state = useAppSelector(state => state.auth)
    const [isValid, setIsValid] = useState(false)
    const navigate = useNavigate()

    useEffect(() => {
        if (state.status === AuthActionStatus.REGISTER_FAILED) {
            toast.error(state.message, {
                position: 'top-right',
                className: "text-base"
            })
        }
        if (state.status === AuthActionStatus.REGISTER_SUCCESSFUL) {
            toast.success(state.message, {
                position: 'top-right',
                className: "text-base"
            })
            navigate("/dashboard")
        }
    }, [state.status]);

    const register = () => {
        const data: RegisterDTO = {
            firstName: firstNameRef?.current?.value || '',
            lastName: lastNameRef?.current?.value || '',
            username: usernameRef?.current?.value || '',
            password: passwordRef?.current?.value || '',
            email: emailRef?.current?.value || '',
        }

        dispatch<any>(authSlice.register(data))
    }

    const validateForm = () => {
        setIsValid(firstNameRef?.current?.value &&
            lastNameRef?.current?.value &&
            usernameRef?.current?.value &&
            passwordRef?.current?.value &&
            emailRef?.current?.value
        )
    }

    return (
        <div className="page-login">
            <div className="w-[80%] bg-white p-10 flex items-stretch rounded-[15px] justify-center shadow-gray-500">
                <form className="signup-form flex-1 p-10" onChange={(e) => validateForm()} onSubmit={(e) => e.preventDefault()}>
                    <div>
                        <h2 className="text-2xl font-bold">Join Code Fusion</h2>
                        <p className="text-xs mt-2 text-gray-700">
                            Sign up to start coding, running, and collaborating in real-time.
                        </p>
                    </div>

                    <div className="flex flex-col gap-2 mt-5">
                        <span className="block text-sm text-gray-700">First Name</span>
                        <input
                            className="rounded-md border border-gray-300 focus:outline-none focus:ring-0 focus:border-gray-300 px-3 py-2"
                            type="text"
                            placeholder="Enter your first name"
                            ref={firstNameRef}
                            name="firstname"
                        />
                    </div>

                    <div className="flex flex-col gap-2 mt-5">
                        <span className="block text-sm text-gray-700">Last Name</span>
                        <input
                            className="rounded-md border border-gray-300 focus:outline-none focus:ring-0 focus:border-gray-300 px-3 py-2"
                            type="text"
                            placeholder="Enter your last name"
                            ref={lastNameRef}
                            name="lastname"
                        />
                    </div>

                    <div className="flex flex-col gap-2 mt-5">
                        <span className="block text-sm text-gray-700">Username</span>
                        <input
                            className="rounded-md border border-gray-300 focus:outline-none focus:ring-0 focus:border-gray-300 px-3 py-2"
                            type="text"
                            placeholder="Enter your username"
                            ref={usernameRef}
                            name="username"
                        />
                    </div>

                    <div className="flex flex-col gap-2 mt-5">
                        <span className="block text-sm text-gray-700">Email</span>
                        <input
                            className="rounded-md border border-gray-300 focus:outline-none focus:ring-0 focus:border-gray-300 px-3 py-2"
                            type="email"
                            placeholder="Enter your email"
                            ref={emailRef}
                            name="email"
                        />
                    </div>

                    <div className="flex flex-col gap-2 mt-5">
                        <span className="block text-sm text-gray-700">Password</span>
                        <input
                            className="rounded-md border border-gray-300 focus:outline-none focus:ring-0 focus:border-gray-300 px-3 py-2"
                            type="password"
                            name="password"
                            placeholder="Enter your password"
                            ref={passwordRef}
                        />
                    </div>

                    {state.status === AuthActionStatus.REGISTER_IN_PROGRESS ? (
                        <div className="flex justify-center items-center mt-7">
                            <BounceLoader
                                color="rgb(37, 99, 235)"
                                loading
                                size={40}
                            />
                        </div>
                    ) : (
                        <button
                            onClick={register}
                            disabled={!isValid}
                            className={`${!isValid ? 'bg-blue-100' : 'bg-blue-600'} text-white bold w-[100%] mt-7 p-3 rounded-md`}>Sign Up
                        </button>
                    )}

                    <div className="text-sm text-gray-700 mt-10 text-center">
                        Already have an account? <Link to="/login" className="text-blue-700">Login</Link>
                    </div>
                </form>

                <div
                    className="flex bg-[rgba(0,0,0,0.05)] rounded-[15px] flex-grow box p-10 flex-1 items-start flex-col justify-center">
                    <h2 className="text-2xl font-bold">Code Fusion</h2>
                    <div className="text-gray-700 text-base mt-2">
                        Connect. Create. Collaborate.
                    </div>
                    <p className="text-gray-600 mt-4">
                        Transform your ideas into reality. Code seamlessly with others, run JavaScript, TypeScript, and
                        Python in one platform.
                    </p>
                </div>
            </div>
        </div>
    );

}

export default Register
