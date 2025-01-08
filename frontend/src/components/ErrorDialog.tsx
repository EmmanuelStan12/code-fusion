const ErrorDialog = ({ error, retryHandler }) => {
    return (
        <div>
            {error ? (
                <ErrorPage
                    title="Unexpected Error"
                    message="We encountered an issue processing your request. Please try again."
                    onRetry={retryHandler}
                />
            ) : (
                <div className="text-gray-100 bg-gray-900 h-screen flex justify-center items-center">
                    <h1 className="text-3xl">{error}</h1>
                </div>
            )}
        </div>
    );
}

export default ErrorDialog