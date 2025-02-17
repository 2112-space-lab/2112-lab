import React, { useEffect, useState } from "react";
import { Spinner, Box, Text, Center } from "@chakra-ui/react";
import useTileServiceStore from "services/tileService";
import GenericTableComponent from "components/table";
import { TileSatelliteMapping } from "types/tiles";

interface MappingTableViewProps {
    onSelectTileID: (tileID: string) => void;
    searchQuery: string;
}

export default function MappingTableView({
    onSelectTileID,
    searchQuery,
}: MappingTableViewProps) {
    const { tileMappings, totalTileMappings, loading, error, fetchTileMappings } = useTileServiceStore();

    const [pageIndex, setPageIndex] = useState<number>(0);
    const [pageSize, setPageSize] = useState<number>(20);

    useEffect(() => {
        fetchTileMappings(pageIndex, pageSize, searchQuery);
    }, [pageIndex, pageSize, searchQuery]);

    const handleOnPaginationChange = (index: number) => {
        setPageIndex(index);
    };

    const handleMappingSelection = (tile: TileSatelliteMapping) => {
        onSelectTileID(tile.TileID);
    };

    const columns = [
        {
            accessorKey: "SpaceID",
            header: "SPACE ID",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "TileID",
            header: "Tile ID",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "Intersection.Longitude",
            header: "I. Longitude",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "Intersection.Latitude",
            header: "I. Latitude",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "TileCenterLat",
            header: "T. Center Lat.",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "TileCenterLon",
            header: "T. Center Lon.",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "TileZoomLevel",
            header: "T. Zoom Level",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
    ];

    return (
        <Box className="grid w-full gap-4 rounded-lg shadow-md relative">
            {/* Spinner at top-right */}

            {/* Error message */}
            {error && (
                <Center className="grid h-full w-full place-items-center">
                    <Box className="text-center">
                        <Text className="text-lg font-bold text-red-500">{error}</Text>
                    </Box>
                </Center>
            )}

            {/* Table */}
            {!error && (
                <GenericTableComponent
                    getRowId={(row: TileSatelliteMapping) => row.MappingID}
                    columns={columns}
                    data={tileMappings}
                    totalItems={totalTileMappings}
                    pageSize={pageSize}
                    pageIndex={pageIndex}
                    onPageChange={handleOnPaginationChange}
                    onRowClick={handleMappingSelection}
                />
            )}
        </Box>
    );
}
