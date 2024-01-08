import { Link } from "react-router-dom";
import img from "../../assets/vector.png";

const Homepage = () => {
    return (
        <div className="w-full h-full flex flex-col justify-between">
            {/* Ad-like tab */}
            <div className="bg-white py-4 px-5 w-full h-[35%] min-h-[270px] rounded-lg flex shadow-md">
                <img alt="make profit" src={img} className="h-full lg:block hidden" />
                <div className="flex flex-col items-start ml-3 py-8">
                    <h1 className="font-semibold md:text-3xl text-2xl mb-7">
                        Trade with other accounts, earn Ethereum and ERC20 tokens and track them
                    </h1>
                    <Link
                        to="/wallet"
                        className="text-2xl from-[#606EFF] to-[#BF66FF] bg-gradient-to-tr text-white px-8 py-2 rounded-md font-semibold"
                    >
                        Go trade
                    </Link>
                </div>
            </div>

            {/* Chart for account balance changes */}
            <div className = "h-[calc(65%_-_50px)] bg-white rounded-md shadow-md flex items-center justify-center">
                <h1 className = "italic opacity-60">
                    For now, this will be empty.
                </h1>
            </div>
        </div>
    );
};

export default Homepage;
