import React, { useState, useEffect, useRef } from "react";
import { Spinner, Box, Center } from "@chakra-ui/react";
import useSatelliteServiceStore from "services/satelliteService"; // Satellite store
import useTileServiceStore from "services/tileService"; // Tile store
import GenericTableComponent from "components/table";
import { OrbitDataItem, SatelliteInfo } from "types/satellites";
import { BiStation, BiTargetLock } from "react-icons/bi";

interface SatelliteTableViewProps {
    onSelectSatelliteID: (satelliteID: string) => void;
    searchQuery: string;
    onTilesSelected: (tileIDs: string[], zoonmTo: boolean) => void;
    onTargetSatellite: (noradID: string, positionData: Record<string, OrbitDataItem[]>) => void; // Callback for targeting satellite with position data
    onPropagateSatellite: (noradID: string) => void; // Callback for targeting satellite
}

export default function SatelliteTableView({
    onSelectSatelliteID,
    searchQuery,
    onTilesSelected,
    onTargetSatellite,
    onPropagateSatellite,
}: SatelliteTableViewProps) {
    const {
        satelliteInfo,
        totalSatelliteInfo,
        orbitData,
        loading,
        fetchPaginatedSatelliteInfo,
        fetchSatellitePositions,
        fetchSatellitePositionsWithPropagation,
    } = useSatelliteServiceStore();

    const { fetchSatelliteMappingsByNoradID, satelliteMappingsByNoradID, recomputeMappingsByNoradID } =
        useTileServiceStore();

    const [pageIndex, setPageIndex] = useState<number>(0);
    const [pageSize, setPageSize] = useState<number>(20);
    const localOrbitDataRef = useRef<Record<string, OrbitDataItem[]>>({});

    useEffect(() => {
        fetchPaginatedSatelliteInfo(pageIndex, pageSize, searchQuery);
    }, [pageIndex, pageSize, searchQuery, fetchPaginatedSatelliteInfo]);

    const handleOnPaginationChange = (index: number) => {
        setPageIndex(index);
    };

    const handleSatelliteSelection = async (satellite: SatelliteInfo) => {
        const noradId = satellite.Satellite.NoradID;

        try {
            await fetchSatelliteMappingsByNoradID(noradId);

            const matchingTileIDs = satelliteMappingsByNoradID[noradId]?.map((tile) => tile.TileID) || [];
            onSelectSatelliteID(noradId);
            onTilesSelected(matchingTileIDs, false);
        } catch (err) {
            console.error("Error fetching tiles for NORAD ID:", err);
        }
    };

    const handleTargetSatellite = async (noradID: string) => {
        const startTime = new Date(Date.now()).toISOString(); // UTC format
        const endTime = new Date(Date.now() + 60 * 60 * 1000 * 24).toISOString(); // UTC format

        try {
            await fetchSatellitePositions(noradID, startTime, endTime);

            localOrbitDataRef.current = { [noradID]: orbitData[noradID] || [] };
            onTargetSatellite(noradID, localOrbitDataRef.current);
        } catch (error) {
            console.error("Error fetching satellite positions:", error);
            localOrbitDataRef.current = {};
        }
    };

    const handlePropagateSatellite = async (noradID: string) => {
        const durationHours = 24;
        const intervalMinutes = 1;

        try {
            await fetchSatellitePositionsWithPropagation(noradID, durationHours, intervalMinutes);
            onPropagateSatellite(noradID);
        } catch (err) {
            console.error("Error propagating satellite:", err);
        }
    };

    const handleRecomputeMapping = async (noradID: string) => {
        const startTime = new Date(Date.now() - 10 * 60 * 1000).toISOString(); // 10 minutes earlier in UTC
        const endTime = new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString(); // 24 hours ahead in UTC

        try {
            await recomputeMappingsByNoradID(noradID, startTime, endTime);
            console.log(`Mappings recomputed successfully for NORAD ID: ${noradID}`);
        } catch (err) {
            console.error(`Error recomputing mapping for NORAD ID: ${noradID}`, err);
        }
    };

    const columns = [
        {
            accessorKey: "",
            header: "Actions",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "TLEs",
            header: "TLE Epoch",
            cell: (info: any) => {
                const tle = info.row.original.TLEs?.[0];
                return (
                    <p className="text-sm">
                        {tle?.Epoch ? new Date(tle.Epoch).toLocaleString() : "N/A"}
                    </p>
                );
            },
        },
        {
            accessorKey: "Satellite.NoradID",
            header: "NORAD ID",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "Satellite.Name",
            header: "Name",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "Satellite.Owner",
            header: "Owner",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "Satellite.LaunchDate",
            header: "Launch Date",
            cell: (info: any) =>
                info.getValue() ? (
                    <p className="text-sm">{new Date(info.getValue()).toISOString().split("T")[0]}</p> // UTC Date
                ) : (
                    <p className="text-sm">N/A</p>
                ),
        },
        {
            accessorKey: "Satellite.Apogee",
            header: "Apogee (km)",
            cell: (info: any) => <p className="text-sm">{info.getValue() ?? "N/A"}</p>,
        },
        {
            accessorKey: "Satellite.Perigee",
            header: "Perigee (km)",
            cell: (info: any) => <p className="text-sm">{info.getValue() ?? "N/A"}</p>,
        },
    ];

    return (
        <Box className="grid w-full gap-4 rounded-lg shadow-md">
            <GenericTableComponent
                getRowId={(row: SatelliteInfo) => row.Satellite.NoradID}
                columns={columns}
                data={satelliteInfo}
                totalItems={totalSatelliteInfo}
                pageSize={pageSize}
                pageIndex={pageIndex}
                onPageChange={handleOnPaginationChange}
                onRowClick={handleSatelliteSelection}
                actions={(row: SatelliteInfo) => {
                    const noradID = row.Satellite.NoradID;
                    const isTargetDisabled = !orbitData[noradID]; // Disable if no orbit data for the NORAD ID

                    return [
                        {
                            label: "Target",
                            onClick: () => handleTargetSatellite(row.Satellite.NoradID),
                            icon: <BiTargetLock />,
                            isDisabled: isTargetDisabled,
                        },
                        {
                            label: "Propagate",
                            onClick: () => handlePropagateSatellite(row.Satellite.NoradID),
                        },
                        {
                            label: "Recompute Mapping",
                            onClick: () => handleRecomputeMapping(row.Satellite.NoradID),
                            icon: <BiStation />,
                            isDisabled: true
                        },
                    ];
                }}
            />
        </Box>
    );
}
