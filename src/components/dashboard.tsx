import { useState, useEffect } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
    faCube,
    faHouse,
    IconDefinition,
    faWallet,
    faUser,
    faArrowRightFromBracket,
    faBars,
} from "@fortawesome/free-solid-svg-icons";
import { faEthereum } from "@fortawesome/free-brands-svg-icons";
import { Link, Outlet, useNavigate } from "react-router-dom";

interface CustomLink {
    icon: IconDefinition;
    text: string;
    path: string;
    hidden: boolean;
}

function Dashboard() {
    const [hidden, setHidden] = useState(false);
    const navigate = useNavigate();

    useEffect(() => {
        window.addEventListener("resize", () => {
            if (window.innerWidth <= 1024) {
                setHidden(true);
                return;
            }
            setHidden(false);
        });
    }, []);

    useEffect(() => {
        window.document.title = "bcDashboard";
    }, []);

    const CustomLink = ({ icon, text, path, hidden }: CustomLink) => {
        return (
            <Link
                to={`/dashboard${path}`}
                className="flex w-full justify-start items-center box-border pl-6 hover:bg-gray-100 py-3"
                title = {text}
                onClick = {() => {
                    if (!hidden && window.innerWidth <= 1024) {
                        setHidden(true);
                    }
                }}
            >
                <FontAwesomeIcon icon={icon} className="text-lg mr-4" />
                <h1 className={`${hidden ? "hidden" : "block"}`}>{text}</h1>
            </Link>
        );
    };

    return (
        <div className="w-full h-full flex">
            {/* Sidebar */}
            <div
                className={`h-full flex-col flex border-r-[1px] bg-white ${
                    hidden ? "w-[60px]" : "w-1/5 min-w-[300px]"
                } duration-200 lg:static lg:z-0 z-20 fixed left-0`}
            >
                <div className={`flex w-full items-center p-4 pl-5 mb-4 ${
                    hidden ? "justify-center" : "justify-between"
                }`}>
                    <div
                        className={`flex items-center ${
                            hidden ? "hidden" : "block"
                        } cursor-pointer duration-100`}
                        onClick={() => {
                            navigate("/dashboard");
                        }}
                    >
                        <FontAwesomeIcon
                            icon={faCube}
                            className="text-3xl from-[#606EFF] to-[#BF66FF] bg-gradient-to-tr text-white p-2 rounded-lg"
                        />
                        <h1 className="text-xl font-semibold ml-4">
                            bcDashboard
                        </h1>
                    </div>
                    <FontAwesomeIcon
                        className={`cursor-pointer text-xl`}
                        icon={faBars}
                        onClick={() => {
                            setHidden(prev => !prev);
                        }}
                    />
                </div>
                {/* Menu bar */}
                <div className="flex w-full flex-col items-center justify-between h-full">
                    <div className="flex w-full flex-col">
                        <CustomLink
                            hidden={hidden}
                            path="/"
                            icon={faHouse}
                            text="Dashboard"
                        />
                        <CustomLink
                            hidden={hidden}
                            path="/wallet"
                            icon={faWallet}
                            text="Wallet"
                        />
                        <CustomLink
                            hidden={hidden}
                            path="/transfer"
                            icon={faEthereum}
                            text="Transfer"
                        />
                    </div>
                    <div className="w-full flex flex-col pb-8">
                        <CustomLink
                            hidden={hidden}
                            path="/account"
                            icon={faUser}
                            text="Account"
                        />
                        <div className="flex w-full justify-start items-center pl-6 py-3 cursor-pointer hover:bg-gray-100">
                            <FontAwesomeIcon
                                icon={faArrowRightFromBracket}
                                className="mr-4"
                            />
                            <h1 className = {`${hidden ? "hidden" : ""}`}>Logout</h1>
                        </div>
                    </div>
                </div>
            </div>
            <div className = {`w-screen h-screen lg:!hidden ${
                !hidden ? "z-10 block" : "hidden"
            } bg-black opacity-60 duration-100`}></div>
            <div
                className={`h-full bg-[#F5F9FC] box-border p-5 ${
                    !hidden ? "lg:w-4/5 lg:max-w-[calc(100%_-_300px)]" : "lg:w-[calc(100%_-_60px)]"
                } duration-200 w-[calc(100%_-_60px)] lg:static absolute right-0`}
            >
                <Outlet />
            </div>
        </div>
    );
}

export default Dashboard;
