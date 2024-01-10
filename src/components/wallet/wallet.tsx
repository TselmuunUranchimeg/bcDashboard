import { useState, useEffect } from "react";
import axios from "axios";
import {
    Key,
    NoKeysSection,
    ShowKey,
    Transaction,
    TransactionComponent,
} from "../../types/wallet";
import "./wallet.css";
import axiosInstance from "../../extras/axios";
import NetworkComponent from "../network";

const KeyComponent = ({ text, isPublic }: Key) => {
    const [copied, setCopied] = useState(false);

    const copyToClipboard = () => {
        navigator.clipboard.writeText(text!);
        setCopied(true);
    };

    const mouseOut = () => {
        if (copied) {
            setCopied(false);
        }
    };

    return (
        <div className="flex items-center justify-between w-full box-border px-5 py-1 mb-2">
            <h1 className="w-1/5 text-left">
                {isPublic ? "Public key" : "Private key"}
            </h1>
            <div
                className={`cursor-pointer w-4/5 overflow-x-auto flex items-center text-left rounded-md box-border p-2 ${
                    copied
                        ? "bg-[#696969] text-white"
                        : "bg-[#c0c0c0] text-black"
                }`}
                title="Copy to clipboard"
                onClick={() => copyToClipboard()}
                onMouseOut={() => mouseOut()}
            >
                <p className="pl-2">
                    {text && text!.length > 42
                        ? `${text.substring(0, 42)}...`
                        : text}
                </p>
            </div>
        </div>
    );
};

const ShowKeyComponent = ({ state, setState }: ShowKey) => {
    return (
        <div
            className={`bg-black bg-opacity-60 absolute lg:left-0 left-[-60px] top-0 w-screen h-screen z-30 justify-center items-center ${
                state ? "flex" : "hidden"
            }`}
        >
            <div className="flex flex-col w-1/3 min-w-[550px] bg-white">
                <div className="bg-[#5C60F5] w-full text-center text-white">
                    <h1 className="font-bold text-xl py-3">
                        Successfully created a wallet
                    </h1>
                </div>
                <div className="w-full flex flex-col text-center">
                    <h1 className="text-red-600 font-semibold my-3">
                        Keep your private key a secret. (This is your only
                        chance.)
                    </h1>
                    <KeyComponent text={state?.public} isPublic={true} />
                    <KeyComponent text={state?.private} isPublic={false} />
                </div>
                <div className="flex justify-center w-full mb-4 mt-2">
                    <button
                        className="bg-[#4e52d0] border-[#4a4dc4] text-white px-7 py-2 rounded-md font-semibold"
                        onClick={() => {
                            setState(null);
                        }}
                    >
                        Close
                    </button>
                </div>
            </div>
        </div>
    );
};

const NoKeysSectionComponent = ({
    visible,
    call,
    setNetwork,
    network
}: NoKeysSection) => {
    return (
        <div
            className={`w-full h-full flex-col items-center justify-center ${
                visible ? "flex" : "hidden"
            }`}
        >
            <div className="flex flex-col items-center">
                <p className="italic opacity-60 mb-5 text-xl">
                    This device currently doesn't have any registered wallets
                </p>
                <div className="sm:w-1/4 w-full mb-5 shadow-md relative">
                    <NetworkComponent setNetwork = {setNetwork} network = {network} />
                </div>
                <button
                    className="bg-[#4e52d0] border-[#4a4dc4] text-white w-40 font-bold text-lg py-2 rounded-md"
                    type="button"
                    onClick={async () => await call()}
                >
                    Create a wallet
                </button>
            </div>
        </div>
    );
};

const TransactionList = ({ address }: TransactionComponent) => {
    const [message, setMessage] = useState("Loading...");
    const [transactions, setTransactions] = useState<Transaction[]>([]);

    useEffect(() => {
        if (address) {
            axiosInstance
                .get(`/address/${address}`)
                .then((res) => {
                    let body = res.data as Transaction[];
                    setTransactions(body);
                    if (body.length === 0) {
                        setMessage(
                            "There are no transactions to display at the moment."
                        );
                    }
                })
                .catch((e) => {
                    if (axios.isAxiosError(e)) {
                        setTransactions([]);
                        setMessage(e.response?.data);
                    }
                });
        }
    }, [address]);

    const AddressHash = ({
        isTx,
        address,
    }: {
        isTx: boolean;
        address: string | null;
    }) => {
        if (!address) {
            return (
                <div className="column">
                    <h1>NULL</h1>
                </div>
            );
        }
        return (
            <div className="column">
                <a
                    href={`https://sepolia.etherscan.io/${
                        isTx ? "tx/" : "address/"
                    }${address}`}
                    target="_blank"
                    className="hover:underline duration-100"
                >
                    {address.substring(0, 7)}...
                </a>
            </div>
        );
    };

    const getAge = (createdAt: string): string => {
        let gap = (Date.now() - new Date(createdAt).getTime()) / 1000; // In seconds
        if (gap >= 60 * 60 * 24) {
            let v = Math.round(gap / (60 * 60 * 24));
            return `${v} ${v === 1 ? "day" : "days"} ago`;
        }
        if (gap >= 60 * 60) {
            let v = Math.round(gap / (60 * 60));
            return `${v} ${v === 1 ? "hour" : "hours"} ago`;
        }
        let v = Math.round(gap / 60);
        return `${v} ${v === 1 ? "minute" : "minutes"} ago`;
    };

    if (!address || transactions.length === 0) {
        return (
            <div className="w-full sm:h-[calc(100%_-_89px)] h-[calc(100%_-_137px)] flex justify-center items-center">
                <h1 className="italic opacity-70">{message}</h1>
            </div>
        );
    }
    return (
        <div className="w-full sm:h-[calc(100%_-_89px)] h-[calc(100%_-_137px)] xl:overflow-x-hidden overflow-x-scroll">
            <div className="row header">
                <div className="column">
                    <h1>Transaction hash</h1>
                </div>
                <div className="column">
                    <h1>From</h1>
                </div>
                <div className="column">
                    <h1>Value</h1>
                </div>
                <div className="column">
                    <h1>Contract</h1>
                </div>
                <div className="column">
                    <h1>Tokens</h1>
                </div>
                <div className="column">
                    <h1>Age</h1>
                </div>
                <div className="column">
                    <h1>Block number</h1>
                </div>
            </div>
            {transactions.map((val, ind) => {
                return (
                    <div
                        key={ind}
                        className={`row ${ind % 2 === 0 ? "" : "bg-[#f5f5f5]"}`}
                    >
                        <AddressHash isTx={true} address={val.hash} />
                        <AddressHash isTx={false} address={val.from} />
                        <div className="column">
                            <h1>{val.value}</h1>
                        </div>
                        <AddressHash isTx={false} address={val.contract} />
                        <div className="column">
                            <h1>
                                {val.tokens.length === 0 ? "NULL" : val.tokens}
                            </h1>
                        </div>
                        <div className="column">
                            <h1>{getAge(val.createdAt)}</h1>
                        </div>
                        <div className="column">
                            <h1>{val.block}</h1>
                        </div>
                    </div>
                );
            })}
        </div>
    );
};

export { NoKeysSectionComponent, ShowKeyComponent, TransactionList };
