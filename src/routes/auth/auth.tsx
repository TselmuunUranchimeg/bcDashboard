import { useState, useEffect, FormEvent } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import axiosInstance from "../../extras/axios";
import "./auth.css";

interface User {
    username: string;
    password: string;
}

const AuthPage = () => {
    const [willLogIn, setLogin] = useState(false);
    const [state, setState] = useState<User>({
        username: "",
        password: "",
    });
    const navigate = useNavigate();

    useEffect(() => {
        window.document.title = `${
            willLogIn ? "Log in" : "Sign up"
        } - bcDashboard`;
    }, [willLogIn]);

    const handleForm = async (e: FormEvent<HTMLFormElement>) => {
        try {
            e.preventDefault();
            let res = await axiosInstance.post(`/auth/${willLogIn ? "login" : "signup"}`, state, {
                withCredentials: true
            });
            if (res.status === 200) {
                alert(res.data);
                navigate("/dashboard");
            }
        } catch (e) {
            if (axios.isAxiosError(e)) {
                alert(e.response?.data);
                return;
            }
            console.log(e);
        }
    };

    return (
        <div className="w-full h-full flex items-center justify-center">
            <div className="flex flex-col sm:shadow-md sm:w-1/4 min-w-[400px] w-full rounded-md">
                <div className="flex cursor-pointer border-t-[1px]">
                    <div
                        className={`${
                            !willLogIn
                                ? "border-white"
                                : "bg-gray-200 border-black border-r-[1px] border-opacity-10"
                        } auth-header`}
                        onClick={() => {
                            setLogin(false);
                        }}
                    >
                        <h1>Sign up</h1>
                    </div>
                    <div
                        className={`${
                            !willLogIn
                                ? "bg-gray-200 border-black border-l-[1px] border-opacity-10"
                                : "border-white"
                        } auth-header`}
                        onClick={() => {
                            setLogin(true);
                        }}
                    >
                        <h1>Log in</h1>
                    </div>
                </div>
                <div className="h-[calc(100%_-_52px)]">
                    <form
                        className="flex flex-col items-center h-full box-border p-5"
                        onSubmit={async (e) => await handleForm(e)}
                    >
                        <input
                            required
                            value={state.username}
                            onChange={(e) => {
                                setState((prev) => {
                                    return {
                                        ...prev,
                                        username: e.target.value,
                                    };
                                });
                            }}
                            type="text"
                            placeholder="Username"
                        />
                        <input
                            required
                            value={state.password}
                            onChange={(e) => {
                                setState((prev) => {
                                    return {
                                        ...prev,
                                        password: e.target.value,
                                    };
                                });
                            }}
                            type="password"
                            placeholder="Password"
                        />
                        <button
                            type="submit"
                            className="bg-[#4e52d0] border-[#4a4dc4] text-white w-full py-3 font-semibold text-xl"
                        >
                            {`${willLogIn ? "Log in" : "Sign up"}`}
                        </button>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default AuthPage;
