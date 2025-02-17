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
    onTargetSatellite: (spaceID: string, positionData: Record<string, OrbitDataItem[]>) => void; // Callback for targeting satellite with position data
    onPropagateSatellite: (spaceID: string) => void; // Callback for targeting satellite
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

    const { fetchSatelliteMappingsBySpaceID, satelliteMappingsBySpaceID, recomputeMappingsBySpaceID } =
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
        const noradId = satellite.Satellite.SpaceID;

        try {
            await fetchSatelliteMappingsBySpaceID(noradId);

            const matchingTileIDs = satelliteMappingsBySpaceID[noradId]?.map((tile) => tile.TileID) || [];
            onSelectSatelliteID(noradId);
            onTilesSelected(matchingTileIDs, false);
        } catch (err) {
            console.error("Error fetching tiles for SPACE ID:", err);
        }
    };

    const handleTargetSatellite = async (spaceID: string) => {
        const startTime = new Date(Date.now()).toISOString(); // UTC format
        const endTime = new Date(Date.now() + 60 * 60 * 1000 * 24).toISOString(); // UTC format

        try {
            await fetchSatellitePositions(spaceID, startTime, endTime);

            localOrbitDataRef.current = { [spaceID]: orbitData[spaceID] || [] };
            onTargetSatellite(spaceID, localOrbitDataRef.current);
        } catch (error) {
            console.error("Error fetching satellite positions:", error);
            localOrbitDataRef.current = {};
        }
    };

    const handlePropagateSatellite = async (spaceID: string) => {
        const durationHours = 24;
        const intervalMinutes = 1;

        try {
            await fetchSatellitePositionsWithPropagation(spaceID, durationHours, intervalMinutes);
            onPropagateSatellite(spaceID);
        } catch (err) {
            console.error("Error propagating satellite:", err);
        }
    };

    const handleRecomputeMapping = async (spaceID: string) => {
        const startTime = new Date(Date.now() - 10 * 60 * 1000).toISOString(); // 10 minutes earlier in UTC
        const endTime = new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString(); // 24 hours ahead in UTC

        try {
            await recomputeMappingsBySpaceID(spaceID, startTime, endTime);
            console.log(`Mappings recomputed successfully for SPACE ID: ${spaceID}`);
        } catch (err) {
            console.error(`Error recomputing mapping for SPACE ID: ${spaceID}`, err);
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
            accessorKey: "Satellite.SpaceID",
            header: "SPACE ID",
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
                getRowId={(row: SatelliteInfo) => row.Satellite.SpaceID}
                columns={columns}
                data={satelliteInfo}
                totalItems={totalSatelliteInfo}
                pageSize={pageSize}
                pageIndex={pageIndex}
                onPageChange={handleOnPaginationChange}
                onRowClick={handleSatelliteSelection}
                actions={(row: SatelliteInfo) => {
                    const spaceID = row.Satellite.SpaceID;
                    const isTargetDisabled = !orbitData[spaceID]; // Disable if no orbit data for the SPACE ID

                    return [
                        {
                            label: "Target",
                            onClick: () => handleTargetSatellite(row.Satellite.SpaceID),
                            icon: <BiTargetLock />,
                            isDisabled: isTargetDisabled,
                        },
                        {
                            label: "Propagate",
                            onClick: () => handlePropagateSatellite(row.Satellite.SpaceID),
                        },
                        {
                            label: "Recompute Mapping",
                            onClick: () => handleRecomputeMapping(row.Satellite.SpaceID),
                            icon: <BiStation />,
                            isDisabled: true
                        },
                    ];
                }}
            />
        </Box>
    );
}
