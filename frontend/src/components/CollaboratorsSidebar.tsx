const CollaboratorsSidebar = ({ collaborators }) => {
    return (
        <aside className="w-80 bg-gray-800 p-4 border-l border-gray-700 flex flex-col">
            <div className="flex justify-between items-center mb-4">
                <h2 className="text-blue-400 font-bold">Collaborators</h2>
            </div>
            <div className="flex-grow overflow-y-auto">
                {collaborators.length === 0 ? (
                    <div className="flex justify-center items-center">
                        <p>No collaborators</p>
                    </div>
                ) : (
                    collaborators.map((collab) => (
                        <div
                            key={collab.id}
                            className="flex justify-between items-center bg-gray-700 text-gray-300 px-3 py-2 rounded-lg mb-2"
                        >
                            <span>{collab.user.username}</span>
                            {collab.isActive && (
                                <span className="w-[15px] h-[15px] rounded-full bg-green-500"></span>
                            )}
                        </div>
                    ))
                )}
            </div>
        </aside>
    );
};

export default CollaboratorsSidebar