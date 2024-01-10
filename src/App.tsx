import { useLayoutEffect } from "react";
import { Outlet, useLocation, useNavigate } from "react-router-dom";
import axios from "axios";
import axiosInstance from "./extras/axios";

const App = () => {
    const navigate = useNavigate();
    const location = useLocation();

    useLayoutEffect(() => {
        axiosInstance
            .get("/auth/verify")
            .then(() => {
                navigate(location.pathname === "/" ? "/dashboard/" : location.pathname, {
                    state: {
                        auth: true
                    }
                });
            })
            .catch((e) => {
                if (axios.isAxiosError(e)) {
                    if (e.response) {
                        if (e.response.status === 406) {
                            console.log("Moving to '/auth'")
                            navigate("/auth");
                        }
                    }
                }
            });
    }, []);

    return (
        <div className = "w-screen h-screen">
            <Outlet />
        </div>
    );
};

export default App;
