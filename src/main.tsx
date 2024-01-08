import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import App from "./App";
import Dashboard from "./components/dashboard";
import ErrorPage from "./ErrorPage";
import Homepage from "./routes/dashboard/homepage";
import AccountPage from "./routes/dashboard/account";
import WalletPage from "./routes/dashboard/wallet";
import TransferPage from "./routes/dashboard/transfer";
import AuthPage from "./routes/auth/auth";
import "./index.css";

const router = createBrowserRouter([
    {
        path: "/",
        element: <App />,
        errorElement: <ErrorPage />,
        children: [
            {
                path: "/auth",
                element: <AuthPage />,
            },
            {
                path: "/dashboard",
                element: <Dashboard />,
                children: [
                    {
                        path: "/dashboard",
                        element: <Homepage />,
                    },
                    {
                        path: "/dashboard/wallet",
                        element: <WalletPage />,
                    },
                    {
                        path: "/dashboard/account",
                        element: <AccountPage />,
                    },
                    {
                        path: "/dashboard/transfer",
                        element: <TransferPage />,
                    },
                ],
            },
        ],
    },
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
    <React.StrictMode>
        <RouterProvider router={router} />
    </React.StrictMode>
);
