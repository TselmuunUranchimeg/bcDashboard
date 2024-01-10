import { useState, useEffect, Dispatch, SetStateAction, Fragment } from "react";
import { Listbox, Transition } from "@headlessui/react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faChevronDown } from "@fortawesome/free-solid-svg-icons";
import axios from "axios";
import { Network } from "../types/general";
import axiosInstance from "../extras/axios";

interface ComponentInterface {
    network: Network | null;
    setNetwork: Dispatch<SetStateAction<Network | null>>;
}

const NetworkComponent = ({ network, setNetwork }: ComponentInterface) => {
    const [networks, setNetworks] = useState<Array<Network>>([]);

    useEffect(() => {
        axiosInstance
            .get("networks")
            .then((res) => {
                let body = res.data as Array<Network>;
                setNetworks(body);
                setNetwork(body[0]);
            })
            .catch((e) => {
                if (axios.isAxiosError(e)) {
                    if (e.code === "ERR_NETWORK") {
                        alert("Either the server is down, or your network is really bad. Please try again after a few minutes.");
                        setNetwork(null);
                        return;
                    }
                    alert(e.response?.data);
                }
            });
    }, []);

    return (
        <Listbox value={network} onChange={setNetwork}>
            <Listbox.Button className="flex items-center justify-between w-full bg-gray-100 box-border px-4 py-2 relative">
                <h1>{network ? network.name : "No data to display"}</h1>
                <FontAwesomeIcon icon={faChevronDown} />
            </Listbox.Button>
            <Transition
                as={Fragment}
                leave="transition ease-in duration-100"
                leaveFrom="opacity-100"
                leaveTo="opacity-0"
            >
                <Listbox.Options className="absolute z-10 w-full shadow-md max-h-[250px] overflow-y-auto">
                    {networks.map((val, ind) => {
                        return (
                            <Listbox.Option
                                key={ind}
                                value={val}
                                className="bg-white hover:bg-gray-200 box-border py-2 px-4"
                            >
                                <h1>{val.name}</h1>
                            </Listbox.Option>
                        );
                    })}
                </Listbox.Options>
            </Transition>
        </Listbox>
    );
};

export default NetworkComponent;