import { useState, useEffect, Fragment, FormEvent } from "react";
import axios from "axios";
import { Listbox, Transition } from "@headlessui/react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faChevronDown, faXmark } from "@fortawesome/free-solid-svg-icons";
import axiosInstance from "../../extras/axios";

type Body = {
    privateKey: string,
    publicKey: string,
    to: string,
    amount: number,
    tokenAddress?: string
};

const TransferPage = () => {
    const [keys, setKeys] = useState<string[] | null>([]);
    const [key, setKey] = useState("");
    const [balance, setBalance] = useState("");
    const [prep, setPrep] = useState(false);
    const [isEth, setIsEth] = useState(true);
    const [to, setTo] = useState("");
    const [contract, setContract] = useState("");
    const [pk, setPk] = useState("");
    const [amount, setAmount] = useState("");

    useEffect(() => {
        window.document.title = "bcDashboard - Transfer";
    }, []);

    useEffect(() => {
        axiosInstance
            .get("/wallet/fetch")
            .then((res) => {
                let body = res.data as string[];
                setKeys(body);
                setKey(body[0]);
            })
            .catch((e) => {
                if (axios.isAxiosError(e)) {
                    alert(e.response?.data);
                }
            });
    }, []);

    useEffect(() => {
        if (key !== "" && key !== undefined) {
            axiosInstance
                .get(`/wallet/${key}`)
                .then(res => {
                    let body = res.data as string;
                    setBalance(body);
                })
                .catch(e => {
                    if (axios.isAxiosError(e)) {
                        alert(e.response?.data);
                    }
                })
        }
    }, [key]);

    if (!keys) {
        return (
            <div className="w-full h-full flex items-center justify-center">
                <h1 className="italic opacity-60">Loading...</h1>
            </div>
        );
    }

    if (keys.length === 0) {
        return (
            <div className="w-full h-full flex items-center justify-center">
                <h1 className="italic opacity-60">
                    You don't have any wallet at the moment.
                </h1>
            </div>
        );
    }

    const handlePrep = (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setPrep(true);
    }

    const Row = ({ keyLeft, value }: { keyLeft: string, value: string }) => {
        return (
            <div className = "w-full flex justify-between border-[1px] mb-4 p-3">
                <h1 className = "w-1/4">{keyLeft}</h1>
                <h1 className = "w-3/4 truncate">{value}</h1>
            </div>
        )
    }

    return (
        <div className="w-full h-full bg-white shadow-md box-border p-5">
            <div className = {`justify-center box-border pt-32 w-screen h-screen absolute top-0 lg:left-0 left-[-60px] ${prep ? "flex" : "hidden"} z-20 bg-black bg-opacity-60`}>
                <div className = {`sm:w-[450px] w-full flex-col`}>
                    <div className = "from-[#606EFF] to-[#BF66FF] bg-gradient-to-tr w-full flex justify-between text-white text-xl items-center p-3 rounded-t-md">
                        <h1 className = "font-semibold">Confirmation</h1>
                        <FontAwesomeIcon 
                            icon = {faXmark}
                            className = "cursor-pointer"
                            onClick = {() => {
                                setPrep(false);
                            }} 
                        />
                    </div>
                    <div className = "bg-white box-border px-5 py-5 rounded-b-md">
                        <Row keyLeft = "From" value = {key} />
                        <Row keyLeft = "To" value = {to} />
                        <Row keyLeft = "Private key" value = {pk} />
                        <Row keyLeft = "Amount" value = {amount} />
                        {
                            !isEth
                            ? <Row keyLeft = "Contract" value = {contract} />
                            : <></>
                        }
                        <button 
                            type = "button" 
                            className = "w-full bg-[#606EFF] text-white py-3 font-semibold"
                            onClick = {async () => {
                                try {
                                    setPrep(false);
                                    let body: Body = {
                                        privateKey: pk,
                                        publicKey: key,
                                        to: to,
                                        amount: isEth ? parseFloat(amount) : parseInt(amount)
                                    };
                                    if (!isEth) {
                                        body.tokenAddress = contract;
                                    }
                                    let res = await axiosInstance.post(isEth ? "/transfer/eth" : "/transfer/tokens", body)
                                    let data = res.data as string;
                                    alert(`Your transaction hash is ${data}. You can check the result on sites like Etherscan.`);
                                } catch (e) {
                                    if (axios.isAxiosError(e)) {
                                        alert(e.response?.data);
                                    }
                                }
                            }}
                        >
                            Confirm
                        </button>
                    </div>
                </div>
            </div>
            <div className="w-full relative mb-8 shadow-md">
                <Listbox value={key} onChange={setKey}>
                    <Listbox.Button className="flex items-center justify-between w-full bg-gray-100 box-border px-4 py-3 relative">
                        <h1 className="truncate">{key}</h1>
                        <FontAwesomeIcon icon={faChevronDown} />
                    </Listbox.Button>
                    <Transition
                        as={Fragment}
                        leave="transition ease-in duration-100"
                        leaveFrom="opacity-100"
                        leaveTo="opacity-0"
                    >
                        <Listbox.Options className="z-10 absolute w-full shadow-md max-h-[250px] overflow-y-auto">
                            {keys?.map((val, ind) => {
                                return (
                                    <Listbox.Option
                                        value={val}
                                        key={ind}
                                        className={`w-full box-border px-4 py-2 hover:bg-gray-200 truncate hover:cursor-pointer ${
                                            val === key
                                                ? "bg-gray-200"
                                                : "bg-[#ffffff]"
                                        }`}
                                    >
                                        {() => (
                                            <>
                                                <span>{val}</span>
                                            </>
                                        )}
                                    </Listbox.Option>
                                );
                            })}
                        </Listbox.Options>
                    </Transition>
                </Listbox>
            </div>
            <div className = "w-full flex md:flex-row flex-col shadow-md box-border px-5">
                <div className = "w-full flex justify-between box-border p-3 pl-0">
                    <h1 className = "font-semibold">Balance</h1>
                    <h1>{balance}</h1>
                </div>
            </div>
            <div className = "w-full box-border px-3 py-5">
                <div className = "md:w-1/3 w-full flex mb-5">
                    <div className = {`w-1/2 flex items-center box-border py-5 duration-100 justify-center text-2xl font-semibold ${
                        isEth ? "bg-gray-200" : "bg-white"
                    } cursor-pointer`}
                        onClick = {() => {
                            if (!isEth) {
                                setIsEth(true);
                            }
                        }}
                    >
                        <h1>Ethereum</h1>
                    </div>
                    <div className = {`w-1/2 flex items-center box-border py-5 duration-100 justify-center text-2xl font-semibold ${
                        isEth ? "bg-white" : "bg-gray-200"
                    } cursor-pointer`}
                        onClick = {() => {
                            if (isEth) {
                                setIsEth(false);
                            }
                        }}
                    >
                        <h1>Tokens</h1>
                    </div>
                </div>
                <form className = "w-full" onSubmit = {e => handlePrep(e)}>
                    <input 
                        required 
                        type = "text" 
                        placeholder = "Recipient" 
                        onChange = {e => {
                            setTo(e.target.value);
                        }}
                        value = {to}
                    />
                    <input 
                        required 
                        type = "password" 
                        placeholder = "Private key" 
                        value = {pk}
                        onChange = {e => {
                            setPk(e.target.value);
                        }}
                    />
                    <input 
                        required 
                        type = "number" 
                        placeholder = "Amount" 
                        min = {isEth ? 0.00001 : 1} 
                        step = {isEth ? 0.00001 : 1} 
                        value = {amount}
                        onChange = {e => {
                            setAmount(e.target.value);
                        }}
                    />
                    <input 
                        required = {isEth ? false : true}
                        type = "text"
                        placeholder = "Contract address"
                        value = {contract}
                        onChange = {(e) => {
                            setContract(e.target.value);
                        }}
                        className = {`${isEth ? "hidden" : ""}`}
                    />
                    <button 
                        type = "submit"
                        className = "from-[#606EFF] to-[#BF66FF] bg-gradient-to-tr text-white px-5 py-3 rounded-md font-semibold text-xl"
                    >
                        Transfer
                    </button>
                </form>
            </div>
        </div>
    );
};

export default TransferPage;
