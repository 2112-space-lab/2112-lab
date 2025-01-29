import React, { useEffect } from "react";
import { Spinner, Center, Text, Box } from "@chakra-ui/react";
import useTileServiceStore from "services/tileService"; // Import your tile service store
import MapTileCard from "./MapTileCard";
import { OrbitDataItem } from "types/satellites";

interface MapTileViewProps {
  selectedTileIDs?: string[];
  satellitePositionData?: Record<string, OrbitDataItem[]>; // New prop for satellite position data
  zoomTo: boolean;
}

export default function MapTileView({
  selectedTileIDs,
  satellitePositionData,
  zoomTo,
}: MapTileViewProps) {
  const {
    tiles,
    loading,
    error,
    fetchTilesForLocation,
  } = useTileServiceStore();

  useEffect(() => {
    fetchTilesForLocation({ latitude: 0, longitude: 0 });
  }, [fetchTilesForLocation]);

  // If there is an error, show the error message
  if (error) {
    return (
      <Center
        position="absolute"
        top="0"
        left="0"
        right="0"
        bottom="0"
        bg="red.700"
        color="white"
      >
        <Text>{error}</Text>
      </Center>
    );
  }

  return (
    <>
      {/* MapTileCard rendering */}
      <MapTileCard
        tiles={tiles}
        darkmode={
          document.body.classList.contains("dark")
            ? "mapbox://styles/simmmple/cl0qqjr3z000814pq7428ptk5"
            : "mapbox://styles/simmmple/ckwxecg1wapzp14s9qlus38p0"
        }
        onLocationChange={fetchTilesForLocation}
        selectedTileIDs={selectedTileIDs}
        satellitePositionData={satellitePositionData} // Pass position data to MapTileCard
        zoomTo={zoomTo}
      />
    </>
  );

}
