import { Dispatch, SetStateAction } from "react";
import { Network } from "./general";

export interface KeyPair {
    private: string;
    public: string;
}
export interface Key {
    text: string | undefined;
    isPublic: boolean;
}
export interface NoKeysSection {
    visible: boolean;
    call: () => Promise<void>;
    network: Network | null;
    setNetwork: Dispatch<SetStateAction<Network | null>>;
}
export interface ShowKey {
    state: KeyPair | null;
    setState: Dispatch<SetStateAction<KeyPair | null>>;
}
export interface Transaction {
    to: string;
    from: string;
    contract: string;	
    value: string;
	hash: string;
	tokens: string;
	createdAt: string;
    block: string;
}
export interface TransactionComponent {
    address: string | null;
}