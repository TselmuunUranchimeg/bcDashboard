import { Fragment, useEffect, useState } from "react";
import axios from "axios";
import { Listbox, Transition } from "@headlessui/react";
import { faChevronDown } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
    NoKeysSectionComponent,
    ShowKeyComponent,
    TransactionList,
} from "../../components/wallet/wallet";
import axiosInstance from "../../extras/axios";
import { KeyPair } from "../../types/wallet";
import { Network } from "../../types/general";
import NetworkComponent from "../../components/network";

const WalletPage = () => {
    const [noKeys, setNoKeys] = useState(true);
    const [showKey, setShowKey] = useState<KeyPair | null>(null);
    const [key, setKey] = useState<string>("");
    const [keys, setKeys] = useState<string[] | null>(null);
    const [network, setNetwork] = useState<Network | null>(null);

    useEffect(() => {
        
    }, []);

    useEffect(() => {
        if (network) {
            axiosInstance
                .get(`/wallet/fetch/${network.id}`)
                .then((res) => {
                    let body = res.data as string[];
                    setKeys(body);
                    setNoKeys(body.length > 0 ? false : true);
                    if (body.length > 0) {
                        setKey(body[0]);
                    }
                })
                .catch((e) => {
                    if (axios.isAxiosError(e)) {
                        if (e.response) {
                            alert(e.response.data);
                        }
                    }
                });
        }
    }, [network]);

    const createNewAccount = async () => {
        try {
            let res = await axiosInstance.post("/wallet/create", {
                id: network?.id,
            });
            let body = res.data as KeyPair;
            setNoKeys(false);
            setShowKey(body);
            if (key === "") {
                setKey(body.public);
                setKeys([body.public]);
            } else {
                setKeys((prev) => {
                    let newVal = prev!;
                    if (newVal.findIndex((val) => val === body.public) === -1) {
                        newVal.push(body.public);
                    }
                    return newVal;
                });
            }
        } catch (e) {
            if (axios.isAxiosError(e)) {
                console.log(e.message);
            }
        }
    };

    useEffect(() => {
        if (!showKey && keys) {
            setKey(keys[keys.length - 1]);
        }
    }, [showKey]);

    return (
        <div className="w-full h-full">
            <NoKeysSectionComponent
                visible={noKeys}
                call={createNewAccount}
                network={network}
                setNetwork={setNetwork}
            />
            <ShowKeyComponent state={showKey} setState={setShowKey} />
            <div
                className={`w-full h-full flex flex-col box-border px-5 ${
                    noKeys ? "hidden" : "flex"
                }`}
            >
                <div className="bg-white w-full h-full rounded-md shadow-md">
                    {/* Top section */}
                    <div className="flex justify-between box-border py-5 px-8 sm:items-center items-start sm:flex-row flex-col border-b-[1px] border-opacity-10 border-black">
                        <div className="sm:w-[calc(100%_-_350px)] sm:mb-0 mb-4 w-full relative">
                            <Listbox value={key} onChange={setKey}>
                                <Listbox.Button className="flex items-center justify-between w-full bg-gray-100 box-border px-4 py-2 relative">
                                    <h1 className="truncate">{key}</h1>
                                    <FontAwesomeIcon icon={faChevronDown} />
                                </Listbox.Button>
                                <Transition
                                    as={Fragment}
                                    leave="transition ease-in duration-100"
                                    leaveFrom="opacity-100"
                                    leaveTo="opacity-0"
                                >
                                    <Listbox.Options className="absolute z-10 w-full shadow-md max-h-[250px] overflow-y-auto">
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
                        <div className="flex justify-between w-full sm:w-[325px] items-center">
                            <div className="w-[calc(50%_-_20px)] relative">
                                <NetworkComponent 
                                    setNetwork = {setNetwork}
                                    network = {network}
                                />
                            </div>
                            <button
                                type="button"
                                onClick={async () => {
                                    await createNewAccount();
                                }}
                                className="bg-[#4e52d0] border-[#606EFF] text-white sm:p-3 p-2 w-1/2 rounded-md hover:bg-[#6974ea] duration-100"
                            >
                                Create new wallet
                            </button>
                        </div>
                    </div>
                    <TransactionList address={key} />
                </div>
            </div>
        </div>
    );
};

export default WalletPage;
