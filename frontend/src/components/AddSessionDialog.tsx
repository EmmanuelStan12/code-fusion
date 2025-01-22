import {useState, useRef, useEffect} from 'react';
import { BounceLoader } from 'react-spinners';
import {UserActionStatus} from "../features/users/users.actions.ts";
import {CodeSessionActionStatus} from "../features/code-session/session.actions.ts";
import {useAppSelector} from "../redux/hooks.ts";
import {IUser} from "../features/users/users.slice.ts";

const AddSessionDialog = ({ showDialog, setShowDialog, codeSessionState, createNewSession }) => {
    const [selectedCollaborators, setSelectedCollaborators] = useState([]);
    const [isValid, setIsValid] = useState(false);
    const sessionTitleRef = useRef<any>();
    const languageRef = useRef<any>();
    const usersState = useAppSelector(state => state.users)

    const [availableUsers, setAvailableUsers] = useState<IUser[]>(usersState.data?.users || []);

    useEffect(() => {
        setAvailableUsers(usersState.data?.users || [])
    }, [usersState.data?.users]);
    const validateForm = () => {
        setIsValid(sessionTitleRef.current?.value?.trim().length > 0);
    };

    const handleCollaboratorSelect = (e) => {
        const selectedId = parseInt(e.target.value);
        if (selectedId) {
            const selectedUser = availableUsers.find(user => user.id === selectedId);
            setSelectedCollaborators([...selectedCollaborators, selectedUser]);
            setAvailableUsers(availableUsers.filter(user => user.id !== selectedId));
            e.target.value = '';
        }
    };

    const removeCollaborator = (userId) => {
        const userToRemove = selectedCollaborators.find(user => user.id === userId);
        setSelectedCollaborators(selectedCollaborators.filter(user => user.id !== userId));
        setAvailableUsers([...availableUsers, userToRemove]);
    };

    return (
        showDialog && (
            <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
                <div className="relative bg-white rounded-lg p-6 w-[90%] max-w-lg shadow-lg">
                    {usersState.status === UserActionStatus.FETCH_USERS_IN_PROGRESS && (
                        <div className="absolute inset-0 bg-white bg-opacity-70 flex items-center justify-center rounded-lg z-10">
                            <BounceLoader
                                color="rgb(37, 99, 235)"
                                loading
                                size={80}
                            />
                        </div>
                    )}

                    <form onChange={validateForm} onSubmit={(e) => e.preventDefault()}>
                        <h2 className="text-2xl font-semibold text-gray-800 mb-4">
                            Create a New Session
                        </h2>
                        <div className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Session Title
                                </label>
                                <input
                                    type="text"
                                    name="title"
                                    placeholder="Enter session title"
                                    ref={sessionTitleRef}
                                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Programming Language
                                </label>
                                <select
                                    name="language"
                                    ref={languageRef}
                                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                >
                                    <option value="JavaScript">JavaScript</option>
                                    <option value="TypeScript">TypeScript</option>
                                </select>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Add Collaborators
                                </label>
                                <select
                                    onChange={handleCollaboratorSelect}
                                    className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    defaultValue=""
                                >
                                    <option value="">Select a collaborator</option>
                                    {availableUsers.map(user => (
                                        <option key={user.id} value={user.id}>
                                            {user.username}
                                        </option>
                                    ))}
                                </select>
                            </div>
                            {selectedCollaborators.length > 0 && (
                                <div className="flex flex-wrap gap-2 pt-2">
                                    {selectedCollaborators.map(user => (
                                        <div
                                            key={user.id}
                                            className="flex items-center gap-1 bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm"
                                        >
                                            {user.username}
                                            <button
                                                type="button"
                                                onClick={() => removeCollaborator(user.id)}
                                                className="hover:bg-blue-200 rounded-full p-1 leading-none ml-1"
                                                aria-label="Remove collaborator"
                                            >
                                                ×
                                            </button>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                        <div className="flex justify-end mt-6">
                            <button
                                type="button"
                                className="bg-gray-200 text-gray-700 px-4 py-2 rounded-md mr-4 hover:bg-gray-300"
                                onClick={() => setShowDialog(false)}
                                disabled={codeSessionState.status === CodeSessionActionStatus.CREATE_SESSION_IN_PROGRESS}
                            >
                                Cancel
                            </button>
                            <button
                                type="button"
                                className={`${!isValid ? 'bg-blue-100' : 'bg-blue-600'} text-white px-4 py-2 rounded-md`}
                                onClick={() => createNewSession(languageRef.current.value, sessionTitleRef.current.value, selectedCollaborators.map(user => user.id).join(','))}
                                disabled={!isValid}
                            >
                                Create Session
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        )
    );
};

export default AddSessionDialog;