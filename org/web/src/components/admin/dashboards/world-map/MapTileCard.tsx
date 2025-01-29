import React, { useRef, useState, useEffect, useCallback } from "react";
import { Box } from "@chakra-ui/react";
import Map, { Source, Layer, NavigationControl, GeolocateControl, MapRef } from "react-map-gl";
import "mapbox-gl/dist/mapbox-gl.css";
import Card from "components/card";
import { Tile } from "types/tiles";
import { OrbitDataItem } from "types/satellites";
import { FeatureCollection, Polygon, Position } from "geojson";

const MAPBOX_TOKEN = process.env.MAPBOX_TOKEN;

interface MapTileCardProps {
  tiles: Tile[];
  darkmode: string;
  onLocationChange: (location: { latitude: number; longitude: number }) => void;
  selectedTileIDs?: string[]; // Updated to accept an array of selected tile IDs
  satellitePositionData?: Record<string, OrbitDataItem[]>; // Added satellite position data
  zoomTo: boolean,
}

const generateSquare = (lat: number, lon: number, size: number): Polygon => {
  const earthRadius = 6378137; // Earth's radius in meters for Mercator projection
  const halfSize = size;

  const mercatorX = (lon * Math.PI * earthRadius) / 180;
  const mercatorY =
    earthRadius * Math.log(Math.tan(Math.PI / 4 + (lat * Math.PI) / 360));

  const cornersMercator = [
    [mercatorX - halfSize, mercatorY - halfSize],
    [mercatorX + halfSize, mercatorY - halfSize],
    [mercatorX + halfSize, mercatorY + halfSize],
    [mercatorX - halfSize, mercatorY + halfSize],
    [mercatorX - halfSize, mercatorY - halfSize],
  ];

  const cornersLatLon = cornersMercator.map(([x, y]) => {
    const lonDeg = (x * 180) / (Math.PI * earthRadius);
    const latRad = (2 * Math.atan(Math.exp(y / earthRadius))) - Math.PI / 2;
    const latDeg = (latRad * 180) / Math.PI;
    return [lonDeg, latDeg] as [number, number];
  });
  // Create GeoJSON Polygon
  const coordinates: Position[][] = [cornersLatLon];
  return {
    type: "Polygon",
    coordinates,
  };
};

const generateSatellitePoints = (satellitePositionData: Record<string, OrbitDataItem[]>) => {
  const points: FeatureCollection = {
    type: "FeatureCollection",
    features: Object.keys(satellitePositionData || {}).flatMap((satelliteID) => {
      const positions = satellitePositionData[satelliteID];

      // Map each position to a Point feature
      return positions?.map((pos) => ({
        type: "Feature",
        geometry: {
          type: "Point",
          coordinates: [pos.longitude, pos.latitude],
        },
        properties: {
          satelliteID,
        },
      }));
    }),
  };

  return points;
};


const MapTileCard: React.FC<MapTileCardProps> = ({
  tiles,
  darkmode,
  onLocationChange,
  selectedTileIDs = [], // Default to an empty array
  satellitePositionData,
  zoomTo,
}) => {
  const mapRef = useRef<MapRef | null>(null);
  const [hoveredTile, setHoveredTile] = useState<Tile | null>(null);
  const [clusterZoom, setClusterZoom] = useState<number>(5);

  useEffect(() => {
    if (selectedTileIDs.length > 0) {
      const firstSelectedTile = tiles.find((tile) => tile.ID === selectedTileIDs[0]);
      if (firstSelectedTile && mapRef.current && zoomTo) {
        mapRef.current.flyTo({
          center: [firstSelectedTile.CenterLon, firstSelectedTile.CenterLat],
          zoom: 3,
        });
      }
    }
  }, [selectedTileIDs, tiles]);
  const handleZoomChange = () => {
    const currentZoom = mapRef.current?.getMap()?.getZoom();
    if (currentZoom) {
      setClusterZoom(Math.floor(currentZoom));
    }
  };

  const handleTileHover = (event: any) => {
    const features = event.features;
    if (features && features.length > 0) {
      const hoveredFeature = features[0];
      const properties = hoveredFeature.properties;

      if (properties?.id) {
        const tile = tiles.find((t) => t.ID === properties.id);
        if (tile) {
          setHoveredTile(tile);
        }
      }
    } else {
      setHoveredTile(null);
    }
  };
  const geoJsonSource: FeatureCollection = {
    type: "FeatureCollection",
    features: tiles.map((tile) => ({
      type: "Feature",
      geometry: generateSquare(tile.CenterLat, tile.CenterLon, tile.Radius),
      properties: {
        id: tile.ID,
        quadkey: tile.Quadkey,
        zoomLevel: tile.ZoomLevel,
        nbFaces: tile.NbFaces,
        radius: tile.Radius,
      },
    })),
  };

  return (
    <Card extra={"relative w-full h-full bg-white px-3 py-[18px]"}>
      <Box position="relative" w="100%" h="100%" overflow="hidden" bg="gray.50" borderRadius="md">
        <Map
          ref={mapRef}
          initialViewState={{
            longitude: 0, // Center of the Earth
            latitude: 0,  // Center of the Earth
            zoom: 1,      // Global zoom to show the entire Earth
          }}
          style={{
            borderRadius: "20px",
            width: "100%",
            height: "100%",
          }}
          mapStyle={darkmode}
          mapboxAccessToken={MAPBOX_TOKEN}
          interactiveLayerIds={["tile-boundaries"]}
          onMouseMove={handleTileHover}
          onZoomEnd={handleZoomChange} // Track zoom level changes
        // maxBounds={[-180, -85, 180, 85]} // Longitude range between -180 to 180 and latitude range between -85 to 85
        >
          <Source id="tiles" type="geojson" data={geoJsonSource}>
            <Layer
              id="tile-boundaries"
              type="fill"
              paint={{
                "fill-color": [
                  "case",
                  ["in", ["get", "id"], ["literal", selectedTileIDs]],
                  "#0D47A1", // Highlight color for selected tiles
                  "#888888", // Default color
                ],
                "fill-opacity": 0.4,
              }}
            />
            <Layer
              id="tile-borders"
              type="line"
              paint={{
                "line-color": "#000000",
                "line-width": 2,
              }}
            />
          </Source>

          {/* Render the satellite paths as straight lines */}
          <Source id="satellite-points" type="geojson" data={generateSatellitePoints(satellitePositionData || {})}>
            <Layer
              id="satellite-points-layer"
              type="circle"
              paint={{
                "circle-radius": 5, // Size of the points
                "circle-color": "#FF0000", // Red color for the points
                "circle-opacity": 0.8, // Slightly transparent
              }}
            />
          </Source>

          <GeolocateControl
            position="top-right"
            onGeolocate={(position) => {
              const { latitude, longitude } = position.coords;
              onLocationChange({ latitude, longitude });
            }}
            trackUserLocation={true}
          />
          <NavigationControl position="top-right" />
        </Map>
      </Box>
    </Card>
  );
};

export default MapTileCard;
