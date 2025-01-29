import React, { useMemo } from "react";
import { useQuery, gql } from "@apollo/client";
import { Spinner, Box, Text, Center } from "@chakra-ui/react";
import Card from "components/card";
import MapSatellite from "./MapSatellite";

type OrbitDataItem = {
  latitude: number;
  longitude: number;
  altitude: number;
  time: string;
};

interface MapSatelliteWithDataProps {
  noradID: string; // NORAD ID passed as a prop
  userLocation?: { latitude: number; longitude: number }; // Optional user location
}

// GraphQL Query to fetch satellite positions
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

const MapSatelliteWithData: React.FC<MapSatelliteWithDataProps> = ({
  noradID,
  userLocation,
}) => {
  // Set up query parameters
  const startTime = useMemo(() => new Date(Date.now()).toISOString(), []);
  const endTime = useMemo(() => new Date(Date.now() + 60 * 60 * 1000 * 24).toISOString(), []); // 24 hours ahead

  const { data, loading, error } = useQuery(GET_SATELLITE_POSITIONS, {
    variables: { id: noradID, startTime, endTime },
    skip: !noradID,
  });


  let orbitData: OrbitDataItem[] = [];
  if (data && Array.isArray(data.satellitePositionsInRange)) {
    orbitData = data.satellitePositionsInRange.map((item: any) => ({
      latitude: item.latitude,
      longitude: item.longitude,
      altitude: item.altitude,
      time: item.timestamp,
    }));
  } else {
    console.warn(
      "Unexpected data structure:",
      data?.satellitePositionsInRange
    );
  }

  return (
    <Card extra={"w-full h-full bg-white px-3 py-[18px]"}>
      <Box className="grid h-[60vh] grid-cols-1 grid-rows-1 rounded-md relative">
        {/* MapSatellite Component */}
        <MapSatellite
          orbitData={orbitData}
          noradID={noradID}
          userLocation={userLocation}
        />

        {/* Loading Indicator */}
        {loading && (
          <Center className="absolute inset-0 bg-black/70 z-10 flex flex-col">
            <Spinner thickness="4px" speed="0.65s" color="blue.500" size="xl" />
            <Text mt={4} color="white">
              Fetching Orbit Data...
            </Text>
          </Center>
        )}

        {/* Error Message */}
        {error && (
          <Center className="absolute inset-0 bg-red-700 text-white z-10">
            <Text>
              {error.message ||
                "Failed to fetch orbit data. Please try again later."}
            </Text>
          </Center>
        )}
      </Box>
    </Card>
  );
};

export default MapSatelliteWithData;
