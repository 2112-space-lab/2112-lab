import { create } from "zustand";
import apiClient from "../utils/apiClient";
import { Tile, TileSatelliteMapping } from "types/tiles";

interface TileServiceState {
    tiles: Tile[];
    tileMappings: TileSatelliteMapping[];
    satelliteMappingsByNoradID: Record<string, TileSatelliteMapping[]>; // Store mappings per NORAD ID
    totalTileMappings: number;
    totalTiles: number;
    loading: boolean;
    error: string | null;
    fetchTileMappings: (pageIndex: number, pageSize: number, search: string) => Promise<void>;
    fetchTilesForLocation: (location: { latitude: number; longitude: number }) => Promise<void>;
    fetchSatelliteMappingsByNoradID: (noradID: string) => Promise<void>;
    recomputeMappingsByNoradID: (noradID: string, startTime: string, endTime: string) => Promise<void>; // New method
}

const useTileServiceStore = create<TileServiceState>((set) => ({
    tiles: [],
    tileMappings: [],
    satelliteMappingsByNoradID: {}, // Initialize as an empty object
    totalTileMappings: 0,
    totalTiles: 0,
    loading: false,
    error: null,

    fetchTileMappings: async (pageIndex: number, pageSize: number, search: string) => {
        set({ loading: true, error: null });

        try {
            const response = await apiClient.get("/tiles/mappings", {
                params: {
                    page: pageIndex + 1,
                    pageSize,
                    search,
                },
            });
            set({
                tileMappings: response.data?.mappings || [],
                totalTileMappings: response.data?.totalRecords || 0,
                loading: false,
            });
        } catch (err) {
            console.error("Error fetching tile mappings:", err);
            set({
                error: "Failed to load tile mapping data.",
                loading: false,
            });
        }
    },

    fetchTilesForLocation: async ({ latitude, longitude }: { latitude: number; longitude: number }) => {
        set({ loading: true, error: null });

        try {
            const response = await apiClient.get("/tiles/all", {
                params: {
                    minLat: latitude - 1,
                    minLon: longitude - 1,
                    maxLat: latitude + 1,
                    maxLon: longitude + 1,
                },
            });

            set({
                tiles: response.data || [],
                totalTiles: response.data?.length || 0,
                loading: false,
            });
        } catch (err) {
            console.error("Error fetching tile data:", err);
            set({
                error: "Failed to fetch tile data. Please try again later.",
                loading: false,
            });
        }
    },

    fetchSatelliteMappingsByNoradID: async (noradID: string) => {
        set((state) => ({
            ...state,
            loading: true,
            error: null,
        }));

        try {
            const response = await apiClient.get("/tiles/mappings/bynoradID", {
                params: { noradID },
            });

            set((state) => ({
                satelliteMappingsByNoradID: {
                    ...state.satelliteMappingsByNoradID,
                    [noradID]: response.data || [],
                },
                loading: false,
            }));
        } catch (err) {
            console.error("Error fetching satellite mappings by NORAD ID:", err);
            set({
                error: `Failed to fetch satellite mappings for NORAD ID: ${noradID}.`,
                loading: false,
            });
        }
    },

    recomputeMappingsByNoradID: async (noradID: string, startTime: string, endTime: string) => {
        set((state) => ({
            ...state,
            loading: true,
            error: null,
        }));

        try {
            const response = await apiClient.put(
                "/tiles/mappings/recompute/bynoradID",
                {},
                {
                    headers: { Accept: "application/json", "Content-Type": "application/json" },
                    params: {
                        noradID: noradID,
                        startTime: startTime,
                        endTime: endTime,
                    },
                }
            );

            set((state) => ({
                satelliteMappingsByNoradID: {
                    ...state.satelliteMappingsByNoradID,
                    [noradID]: response.data || [],
                },
                loading: false,
            }));
        } catch (err) {
            console.error("Error recomputing satellite mappings by NORAD ID:", err);
            set({
                error: `Failed to recompute satellite mappings for NORAD ID: ${noradID}.`,
                loading: false,
            });
        }
    },
}));

export default useTileServiceStore;
