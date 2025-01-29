import { create } from "zustand";
import axios from "axios";

interface ApiStoreState<T> {
    data: T | null;
    loading: boolean;
    error: string | null;
    fetchApi: (url: string, params?: Record<string, any>, headers?: Record<string, string>) => Promise<void>;
    setData: (data: T) => void;
    setLoading: (loading: boolean) => void;
    setError: (error: string | null) => void;
}

const useApiStore = create<ApiStoreState<any>>((set) => ({
    data: null,
    loading: false,
    error: null,

    setData: (data) => set({ data }),
    setLoading: (loading) => set({ loading }),
    setError: (error) => set({ error }),

    fetchApi: async (url, params = {}, headers = {}) => {
        set({ loading: true, error: null });

        try {
            const response = await axios.get(url, {
                params,
                headers: { Accept: "application/json", ...headers },
            });
            set({ data: response.data, loading: false });
        } catch (err) {
            console.error("Error fetching API data:", err);
            set({ error: "Failed to fetch data.", loading: false });
        }
    },
}));

export default useApiStore;
