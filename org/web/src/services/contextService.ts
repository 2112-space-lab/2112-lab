import { create } from "zustand";
import apiClient from "../utils/apiClient";

interface Context {
    id: string;
    name: string;
    tenantID: string;
    description?: string;
    maxSatellite: number;
    maxTiles: number;
    activatedAt?: string;
    desactivatedAt?: string;
}

interface ContextServiceState {
    contexts: Context[];
    totalContexts: number;
    loading: boolean;
    error: string | null;
    fetchPaginatedContexts: (pageIndex: number, pageSize: number, search?: string) => Promise<void>;
    createContext: (context: Omit<Context, "id">) => Promise<void>;
    updateContext: (name: string, updates: Partial<Omit<Context, "id" | "tenantID">>) => Promise<void>;
    getContextByName: (name: string) => Promise<Context | null>;
    deleteContextByName: (name: string) => Promise<void>;
    activateContext: (name: string) => Promise<void>;
    deactivateContext: (name: string) => Promise<void>;
}

const useContextServiceStore = create<ContextServiceState>((set) => ({
    contexts: [],
    totalContexts: 0,
    loading: false,
    error: null,

    fetchPaginatedContexts: async (pageIndex: number, pageSize: number, search?: string) => {
        set({ loading: true, error: null });

        try {
            const response = await apiClient.get("/contexts/all", {
                params: {
                    page: pageIndex + 1,
                    pageSize,
                    search,
                },
            });

            set({
                contexts: response.data?.contexts || [],
                totalContexts: response.data?.totalCount || 0,
                loading: false,
            });
        } catch (err) {
            console.error("Error fetching contexts:", err);
            set({
                error: "Failed to load contexts.",
                loading: false,
            });
        }
    },

    createContext: async (context) => {
        set({ loading: true, error: null });

        try {
            await apiClient.post("/contexts", context);
            set({ loading: false });
        } catch (err) {
            console.error("Error creating context:", err);
            set({
                error: "Failed to create context.",
                loading: false,
            });
        }
    },

    updateContext: async (name, updates) => {
        set({ loading: true, error: null });

        try {
            await apiClient.put(`/contexts/${name}`, updates);
            set({ loading: false });
        } catch (err) {
            console.error("Error updating context:", err);
            set({
                error: "Failed to update context.",
                loading: false,
            });
        }
    },

    getContextByName: async (name) => {
        set({ loading: true, error: null });

        try {
            const response = await apiClient.get(`/contexts/${name}`);
            set({ loading: false });
            return response.data;
        } catch (err) {
            console.error("Error fetching context by name:", err);
            set({
                error: "Failed to fetch context.",
                loading: false,
            });
            return null;
        }
    },

    deleteContextByName: async (name) => {
        set({ loading: true, error: null });

        try {
            await apiClient.delete(`/contexts/${name}`);
            set({ loading: false });
        } catch (err) {
            console.error("Error deleting context:", err);
            set({
                error: "Failed to delete context.",
                loading: false,
            });
        }
    },

    activateContext: async (name) => {
        set({ loading: true, error: null });

        try {
            await apiClient.put(`/contexts/${name}/activate`);
            set({ loading: false });
        } catch (err) {
            console.error("Error activating context:", err);
            set({
                error: "Failed to activate context.",
                loading: false,
            });
        }
    },

    deactivateContext: async (name) => {
        set({ loading: true, error: null });

        try {
            await apiClient.put(`/contexts/${name}/deactivate`);
            set({ loading: false });
        } catch (err) {
            console.error("Error deactivating context:", err);
            set({
                error: "Failed to deactivate context.",
                loading: false,
            });
        }
    },
}));

export default useContextServiceStore;
