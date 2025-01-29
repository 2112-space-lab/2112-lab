import { create } from "zustand";
import { useSession } from "@clerk/nextjs";
import apiClient from "../utils/apiClient";

interface User {
    id: string;
    name: string[];
    email: string;
    username: string;
    date: string;
    type: string;
}

interface UserServiceState {
    users: User[];
    totalUsers: number;
    loading: boolean;
    error: string | null;
    fetchPaginatedUsers: (pageIndex: number, pageSize: number) => Promise<void>;
    getUserById: (userId: string) => Promise<User | null>;
    deleteUserById: (userId: string) => Promise<void>;
}

const useUserServiceStore = create<UserServiceState>((set) => {
    const session = useSession();
    const sessionLoaded = session.isLoaded;
    const sessionSignedIn = session.isSignedIn;
    const getToken = session.session?.getToken;

    return {
        users: [],
        totalUsers: 0,
        loading: false,
        error: null,

        fetchPaginatedUsers: async (pageIndex, pageSize) => {
            if (!sessionLoaded || !sessionSignedIn || !getToken) {
                console.error("Session not loaded or user not signed in.");
                set({ error: "User is not authenticated.", loading: false });
                return;
            }

            set({ loading: true, error: null });

            try {
                const token = await getToken(); // Retrieve the token using Clerk session

                const response = await apiClient.get("/users", {
                    headers: {
                        Authorization: `Bearer ${token}`, // Add Authorization header
                    },
                    params: {
                        page: pageIndex + 1,
                        pageSize,
                    },
                });

                const transformedUsers = response.data?.users.map((user: any) => ({
                    id: user.id,
                    name: [
                        user.first_name && user.last_name
                            ? `${user.first_name} ${user.last_name}`
                            : "Unknown User",
                        user.image_url || "https://i.ibb.co/7p0d1Cd/Frame-24.png",
                    ],
                    email: user.email_addresses?.[0]?.email_address || "No Email Available",
                    username: user.username ? `@${user.username}` : "No Username",
                    date: new Date(user.created_at * 1000).toLocaleDateString("en-US"),
                    type: user.type || "Member",
                }));

                set({
                    users: transformedUsers || [],
                    totalUsers: response.data?.total_count || 0,
                    loading: false,
                });
            } catch (err) {
                console.error("Error fetching users:", err);
                set({
                    error: "Failed to load users.",
                    loading: false,
                });
            }
        },

        getUserById: async (userId) => {
            if (!sessionLoaded || !sessionSignedIn || !getToken) {
                console.error("Session not loaded or user not signed in.");
                set({ error: "User is not authenticated.", loading: false });
                return null;
            }

            try {
                set({ loading: true, error: null });
                const token = await getToken(); // Retrieve the token using Clerk session

                const response = await apiClient.get(`/users/${userId}`, {
                    headers: {
                        Authorization: `Bearer ${token}`, // Add Authorization header
                    },
                });

                const user = response.data;
                const transformedUser: User = {
                    id: user.id,
                    name: [
                        user.first_name && user.last_name
                            ? `${user.first_name} ${user.last_name}`
                            : "Unknown User",
                        user.image_url || "https://i.ibb.co/7p0d1Cd/Frame-24.png",
                    ],
                    email: user.email_addresses?.[0]?.email_address || "No Email Available",
                    username: user.username ? `@${user.username}` : "No Username",
                    date: new Date(user.created_at * 1000).toLocaleDateString("en-US"),
                    type: user.type || "Member",
                };

                set({ loading: false });
                return transformedUser;
            } catch (err) {
                console.error("Error fetching user by ID:", err);
                set({
                    error: "Failed to fetch user.",
                    loading: false,
                });
                return null;
            }
        },

        deleteUserById: async (userId) => {
            if (!sessionLoaded || !sessionSignedIn || !getToken) {
                console.error("Session not loaded or user not signed in.");
                set({ error: "User is not authenticated.", loading: false });
                return;
            }

            try {
                set({ loading: true, error: null });
                const token = await getToken(); // Retrieve the token using Clerk session

                await apiClient.delete(`/users/${userId}`, {
                    headers: {
                        Authorization: `Bearer ${token}`, // Add Authorization header
                    },
                });

                set((state) => ({
                    users: state.users.filter((user) => user.id !== userId),
                    loading: false,
                }));
            } catch (err) {
                console.error("Error deleting user:", err);
                set({
                    error: "Failed to delete user.",
                    loading: false,
                });
            }
        },
    };
});

export default useUserServiceStore;
