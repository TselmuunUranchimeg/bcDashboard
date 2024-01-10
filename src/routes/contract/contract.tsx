import { useState, FormEvent } from "react";
import { Network } from "../../types/general";
import NetworkComponent from "../../components/network";
import axiosInstance from "../../extras/axios";
import axios from "axios";

const ContractPage = () => {
    const [network, setNetwork] = useState<Network | null>(null);
    const [address, setAddress] = useState("");

    const formSubmittion = async (e: FormEvent<HTMLFormElement>) => {
        try {
            e.preventDefault();
            if (!network) {
                alert("Something went wrong, please try again later.")
                return;
            }
            if (!confirm(`Are you sure you want to proceed with this value?\n${address}`)) {
                return;
            }
            const re = new RegExp("^0x[a-zA-Z0-9]{40}$")
            if (!re.test(address)) {
                alert("Contract address is not valid");
                return;
            }
            let res = await axiosInstance.post("/contract", {
                address,
                networkId: network.id
            });
            alert(res.data);
        } catch (e) {
            if (axios.isAxiosError(e)) {
                alert(e.code === "ERR_NETWORK" ? "Please check your network" : e.response?.data);
            }
        }
    };

    return (
        <div className="w-full h-full flex items-center justify-center">
            <div className="md:w-1/4 w-4/5 md:min-w-[400px] relative">
                <form
                    className="mt-28"
                    onSubmit={async (e) => await formSubmittion(e)}
                >
                    <div className="w-full text-center mb-5">
                        <h1 className = "text-2xl font-semibold">
                            Register a contract in the database
                        </h1>
                        <p className = "italic opacity-60">
                            and track transactions
                        </p>
                    </div>
                    <input
                        type="text"
                        placeholder="Contract address"
                        value={address}
                        onChange={(e) => {
                            setAddress(e.target.value);
                        }}
                    />
                    <NetworkComponent
                        network={network}
                        setNetwork={setNetwork}
                    />
                    <button
                        type="submit"
                        className="from-[#606EFF] to-[#BF66FF] bg-gradient-to-tr py-2 w-full text-white mt-5 font-semibold text-xl"
                    >
                        Register contract
                    </button>
                </form>
            </div>
        </div>
    );
};

export default ContractPage;
