# Blockchain dashboard

A dashboard that lets the user create a wallet address, track all deposits, and transfer Ethereum and ERC20 tokens to other addressses, while also providing JWT authentication. 

# Documentation

- [Front-end](#front-end)
- [Back-end](/server/README.md#project-structure)

# Tech stack

- Front-end: [React.js](https://react.dev/) with [Typescript](https://www.typescriptlang.org/)
    - Styling: [TailwindCSS](https://tailwindcss.com/)
    - Icons: [Font awesome](https://fontawesome.com/)
    - Components: [Headless UI](https://headlessui.com/)
    - Routing: [React Router](https://reactrouter.com/)
- Back-end: [Golang](https://go.dev) with [go-chi](https://github.com/go-chi/chi)
- Database: [PostgreSQL](https://www.postgresql.org/)

# Front-end

## Project structure

|  Folders  |  Description |
| :-------  | :----------- |
| `assets`  | Images and other files used throughout the application |
| `components` | React components. The folder suggests the ones associated with specific pages, while files are general components that are used in multiple components and pages. |
| `extras` | Services that can be configured with default values |
| `routes` | Pages for the applications. Folder name suggests the parent route, and the file names are the child routes. |
| `types` | Typescript interfaces


## Routing

As shown in the code snippet below, there are only 2 main routes. 

- /auth - For authentication (signing in and logging in)
- /dashboard - To manage wallet addresses and other key information (Only available for authenticated users only)

```
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
                        path: "/dashboard/",
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
                    }
                ],
            },
        ],
    }
```

To make sure there are only 2 main routes, at the root of the `App` component, there is an `useEffect` hook that redirects all requests to these 2 paths. 

Basically, if the client has a valid token, they are considered to be authenticated. Otherwise, they will be redirected to `/auth`.

```
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
```

Additionally, for the one and only admin, there is an another page. The purpose of this path `/contracts` is to provide the admin a simple form to register contracts to track in the database.

```
{
    path: "/contracts",
    element: <ContractPage />
}
```