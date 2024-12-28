import {BounceLoader} from "react-spinners";

export default function LoadingPage() {
    return (
        <div className="flex w-full h-full justify-center items-center">
            <BounceLoader
                color="rgb(37, 99, 235)"
                loading
                size={80}
            />
        </div>
    )
}