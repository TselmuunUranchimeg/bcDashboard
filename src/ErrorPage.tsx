import { useRouteError } from "react-router-dom";

const ErrorPage = () => {
    const error = useRouteError() as any;

    return (
        <div className = "w-screen h-screen flex items-center justify-center">
            <div className = "flex flex-col items-center">
                <h1 className = "font-bold text-3xl mb-3">Oops!</h1>
                <p>Sorry, there was an error. Error message is "{error.statusText || error.message}".</p>
            </div>
        </div>
    )
}

export default ErrorPage;