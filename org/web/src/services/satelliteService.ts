import { create } from "zustand";
import apiClient from "../utils/apiClient";
import { SatelliteInfo, OrbitDataItem } from "types/satellites";
import { gql } from "@apollo/client";
import ApolloClient from "lib/ApolloClient";

interface SatelliteServiceState {
    satelliteInfo: SatelliteInfo[];
    totalSatelliteInfo: number;
    orbitData: Record<string, OrbitDataItem[]>; // Orbit data keyed by NORAD ID
    loading: boolean;
    error: string | null;
    fetchPaginatedSatelliteInfo: (pageIndex: number, pageSize: number, search: string) => Promise<void>;
    fetchSatellitePositions: (noradID: string, startTime: string, endTime: string) => Promise<void>;
    fetchSatellitePositionsWithPropagation: (noradID: string, durationHours: number, intervalMinutes: number) => Promise<void>;
}

const GET_SATELLITE_POSITIONS = gql`
  query GetSatellitePositions($id: ID!, $startTime: String!, $endTime: String!) {
    satellitePositionsInRange(id: $id, startTime: $startTime, endTime: $endTime) {
      latitude
      longitude
      altitude
      timestamp
    }
  }
`;

const useSatelliteServiceStore = create<SatelliteServiceState>((set) => ({
    satelliteInfo: [],
    totalSatelliteInfo: 0,
    orbitData: {},
    loading: false,
    error: null,

    fetchPaginatedSatelliteInfo: async (pageIndex: number, pageSize: number, search: string) => {
        set({ loading: true, error: null });

        try {
            const response = await apiClient.get("/satellites/paginated/tles", {
                params: {
                    page: pageIndex + 1,
                    pageSize,
                    search,
                },
            });

            set({
                satelliteInfo: response.data?.satellites || [],
                totalSatelliteInfo: response.data?.totalRecords || 0,
                loading: false,
            });
        } catch (err) {
            console.error("Error fetching satellite info:", err);
            set({
                error: "Failed to load satellite information.",
                loading: false,
            });
        }
    },

    fetchSatellitePositions: async (noradID: string, startTime: string, endTime: string): Promise<void> => {
        set({ loading: true, error: null });

        try {
            const { data } = await ApolloClient.query({
                query: GET_SATELLITE_POSITIONS,
                variables: { id: noradID, startTime, endTime },
            });

            const orbitData: OrbitDataItem[] = data?.satellitePositionsInRange?.map((item: any) => ({
                latitude: item.latitude,
                longitude: item.longitude,
                altitude: item.altitude,
                time: item.timestamp,
            })) || [];

            if (orbitData.length > 0) {
                set((state) => ({
                    orbitData: { ...state.orbitData, [noradID]: orbitData },
                    loading: false,
                }));
            } else {
                set({
                    error: "No satellite position data found.",
                    loading: false,
                });
            }
        } catch (err) {
            console.error("Error fetching satellite positions:", err);
            set({
                error: "Failed to load satellite positions.",
                loading: false,
            });
        }
    },

    fetchSatellitePositionsWithPropagation: async (noradID: string, durationHours: number, intervalMinutes: number): Promise<void> => {
        set({ loading: true, error: null });

        try {
            const response = await apiClient.get(`/satellites/orbit`, {
                params: {
                    noradID,
                    durationHours,
                    intervalMinutes,
                },
            });

            const positions = response.data || [];

            if (positions.length > 0) {
                const orbitData: OrbitDataItem[] = positions.map((position: any) => ({
                    latitude: position.latitude,
                    longitude: position.longitude,
                    altitude: position.altitude,
                    time: position.timestamp,
                }));

                set((state) => ({
                    orbitData: { ...state.orbitData, [noradID]: orbitData },
                    loading: false,
                }));
            } else {
                set({
                    error: "No satellite position data found.",
                    loading: false,
                });
            }
        } catch (err) {
            console.error("Error fetching satellite positions with propagation:", err);
            set({
                error: "Failed to load satellite positions.",
                loading: false,
            });
        }
    },
}));

export default useSatelliteServiceStore;
