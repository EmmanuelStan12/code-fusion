import {useState} from "react";

const AddCollaboratorDialog = ({ users, setShowDialog, addCollaborator, isLoading }) => {
    const [userId, setUserId] = useState('');

    return (
        <div className="absolute inset-0 bg-gray-900 bg-opacity-75 flex items-center justify-center z-10">
            <div className="bg-white rounded-lg p-6 w-[90%] max-w-lg shadow-lg">
                <h2 className="text-2xl font-semibold text-gray-800 mb-4">
                    Add New Collaborator
                </h2>
                <div className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Users
                        </label>
                        {isLoading ? (
                            <div className="flex justify-center items-center">
                                <BounceLoader color="rgb(37, 99, 235)" loading size={40} />
                            </div>
                        ) : (
                            <select
                                name="language"
                                onChange={(e) => setUserId(e.target.value)}
                                className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            >
                                <option value="">Select a user</option>
                                {users.map(user => (
                                    <option key={user.id} value={user.id}>
                                        {user.username}
                                    </option>
                                ))}
                            </select>
                        )}
                    </div>
                </div>
                <div className="flex justify-end mt-6">
                    <button
                        className="bg-gray-200 text-gray-700 px-4 py-2 rounded-md mr-4 hover:bg-gray-300"
                        onClick={() => setShowDialog(false)}
                        disabled={isLoading}
                    >
                        Cancel
                    </button>
                    {isLoading ? (
                        <div className="flex justify-center items-center mt-7">
                            <BounceLoader color="rgb(37, 99, 235)" loading size={40} />
                        </div>
                    ) : (
                        <button
                            className={`${!userId ? 'bg-blue-100' : 'bg-blue-600'} text-white px-4 py-2 rounded-md`}
                            onClick={addCollaborator}
                            disabled={!userId}
                        >
                            Add Collaborator
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
};

export default AddCollaboratorDialog