import {useEffect} from "react";
import {CodeOperation} from "../../utils/session.ts";
import {getOperations} from "../../utils/merger.ts";

function useDebounce(func: (ops: CodeOperation[]) => any, delay: number) {
    let timer = null

    useEffect(() => {
        return () => clearTimeout(timer)
    }, []);

    return (prevCode: string, currentCode: string) => {
        console.log(timer, prevCode, currentCode)
        clearTimeout(timer)
        timer = setTimeout(() => {
            const debouncedValues = getOperations(prevCode, currentCode)
            console.log(debouncedValues)
            func(debouncedValues)
        }, delay)
    }
}

export default useDebounce