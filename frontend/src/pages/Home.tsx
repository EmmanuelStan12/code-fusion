function Home() {

    return (
        <div className="home-page w-full bg-gray-900 text-gray-100 p-10">
            <header className="w-full flex justify-between items-center px-10 py-5 bg-gray-800 shadow-lg rounded-lg">
                <h1 className="text-3xl font-bold text-blue-500">Code Fusion</h1>
                <nav className="flex gap-5">
                    <a href="#features" className="text-gray-300 hover:text-white">Features</a>
                    <a href="#contact" className="text-gray-300 hover:text-white">Contact</a>
                </nav>
            </header>

            <section className="hero-section flex flex-col items-center text-center py-20 bg-gradient-to-br from-gray-800 to-gray-700 rounded-lg">
                <h2 className="text-4xl font-bold text-blue-400">Code Fusion</h2>
                <p className="text-lg text-gray-300 mt-5">
                    Empowering developers with real-time code execution, seamless collaboration, and a secure sandbox environment.
                </p>
                <button className="bg-blue-600 hover:bg-blue-700 text-white px-5 py-3 mt-8 rounded-lg shadow-lg">
                    Get Started
                </button>
            </section>

            <section id="features" className="features-section py-20">
                <h2 className="text-3xl font-bold text-center mb-10 text-blue-400">Features</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-10 px-10">
                    <div className="bg-gray-800 p-5 rounded-lg shadow-md">
                        <h3 className="text-xl font-semibold text-blue-300">Run Your Code Instantly</h3>
                        <p className="text-gray-300 mt-3">
                            Execute JavaScript, TypeScript, and Python code in isolated containers with instant feedback.
                        </p>
                    </div>
                    <div className="feature-card bg-gray-800 p-5 rounded-lg shadow-md">
                        <h3 className="text-xl font-semibold text-blue-300">Handle Complex Tasks</h3>
                        <p className="text-gray-300 mt-3">
                            Utilize powerful backend environments to process complex algorithms and computations.
                        </p>
                    </div>
                    <div className="feature-card bg-gray-800 p-5 rounded-lg shadow-md">
                        <h3 className="text-xl font-semibold text-blue-300">Custom Environments</h3>
                        <p className="text-gray-300 mt-3">
                            Set up environment requirements and dependencies tailored to your project's needs.
                        </p>
                    </div>
                    <div className="feature-card bg-gray-800 p-5 rounded-lg shadow-md">
                        <h3 className="text-xl font-semibold text-blue-300">Real-Time Collaboration</h3>
                        <p className="text-gray-300 mt-3">
                            Work with your team in real-time, sharing and editing code seamlessly.
                        </p>
                    </div>
                    <div className="feature-card bg-gray-800 p-5 rounded-lg shadow-md">
                        <h3 className="text-xl font-semibold text-blue-300">Secure Sandbox</h3>
                        <p className="text-gray-300 mt-3">
                            Run code in isolated, containerized environments to ensure security and integrity.
                        </p>
                    </div>
                    <div className="feature-card bg-gray-800 p-5 rounded-lg shadow-md">
                        <h3 className="text-xl font-semibold text-blue-300">Enhanced Productivity</h3>
                        <p className="text-gray-300 mt-3">
                            Streamline your workflow with features designed for efficiency and simplicity.
                        </p>
                    </div>
                </div>
            </section>

            <footer id="contact" className="footer py-10 bg-gray-800 rounded-lg text-center">
                <p className="text-gray-400">Built with ♥ by the Code Fusion Team</p>
            </footer>
        </div>
    );
}

export default Home
