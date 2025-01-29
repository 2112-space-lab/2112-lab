import React, { useEffect, useRef, useState } from "react";
import Map, { GeolocateControl, NavigationControl } from "react-map-gl";
import MapboxGeocoder from "@mapbox/mapbox-gl-geocoder"; // Mapbox Geocoder control
import "mapbox-gl/dist/mapbox-gl.css";
import "@mapbox/mapbox-gl-geocoder/dist/mapbox-gl-geocoder.css"; // Geocoder control CSS
import Card from "components/card";

const MAPBOX_TOKEN =
  "pk.eyJ1Ijoic2ltbW1wbGUiLCJhIjoiY2wxeG1hd24xMDEzYzNrbWs5emFkdm16ZiJ9.q9s0sSKQFFaT9fyrC-7--g"; // Replace with your Mapbox token

interface MapCardProps {
  onLocationChange: (location: { latitude: number; longitude: number }) => void; // Callback to pass user location
}

const MapCard: React.FC<MapCardProps> = ({ onLocationChange }) => {
  const [darkmode, setDarkmode] = useState(
    document.body.classList.contains("dark")
      ? "mapbox://styles/simmmple/cl0qqjr3z000814pq7428ptk5"
      : "mapbox://styles/simmmple/ckwxecg1wapzp14s9qlus38p0"
  );

  const mapRef = useRef(null);
  const [userLocation, setUserLocation] = useState<{ latitude: number; longitude: number } | null>(null);

  useEffect(() => {
    const observer = new MutationObserver((mutationsList) => {
      for (const mutation of mutationsList) {
        if (
          mutation.type === "attributes" &&
          mutation.attributeName === "class"
        ) {
          if (document.body.classList.contains("dark")) {
            setDarkmode(
              "mapbox://styles/simmmple/cl0qqjr3z000814pq7428ptk5"
            );
          } else {
            setDarkmode(
              "mapbox://styles/simmmple/ckwxecg1wapzp14s9qlus38p0"
            );
          }
        }
      }
    });
    observer.observe(document.body, { attributes: true });
    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    if (mapRef.current) {
      const map = mapRef.current.getMap();

      // Initialize the Geocoder control
      const geocoder = new MapboxGeocoder({
        accessToken: MAPBOX_TOKEN,
        mapboxgl: map, // Bind to Mapbox GL instance
        marker: true, // Add a marker at the searched location
        // placeholder: "Search for places",
      });

      // Add Geocoder control to the map
      // geocoder.on("result", (e) => {
      //   const { center } = e.result.geometry;
      //   if (center) {
      //     const [longitude, latitude] = center;
      //     setUserLocation({ latitude, longitude });
      //     onLocationChange({ latitude, longitude }); // Notify location change
      //   }
      // });

      // map.addControl(geocoder, "top-left");

      // return () => {
      //   // Remove the geocoder when the component unmounts
      //   map.removeControl(geocoder);
      // };
    }
  }, [onLocationChange]);

  const handleGeolocate = (position: GeolocationPosition) => {
    const { latitude, longitude } = position.coords;
    setUserLocation({ latitude, longitude });
    onLocationChange({ latitude, longitude }); // Notify location change
  };

  const handleMapClick = (e: any) => {
    const { lngLat } = e;
    const { lat, lng } = lngLat;
    setUserLocation({ latitude: lat, longitude: lng });
    onLocationChange({ latitude: lat, longitude: lng }); // Notify location change
  };

  return (
    <Card extra={"relative w-full h-full bg-white px-3 py-[18px]"}>
      <Map
        ref={mapRef}
        initialViewState={{
          latitude: 49.6117, // Latitude for 85 Avenue Guillaume
          longitude: 6.1319, // Longitude for 85 Avenue Guillaume
          zoom: 15, // Adjust zoom level as needed
        }}
        onClick={handleMapClick} // Handle map clicks
        style={{
          borderRadius: "20px",
          width: "100%",
          height: "100%",
        }}
        mapStyle={darkmode}
        mapboxAccessToken={MAPBOX_TOKEN}
      >
        {/* Optional Controls */}
        <GeolocateControl
          position="top-right"
          onGeolocate={handleGeolocate} // Triggered when user's location is determined
          trackUserLocation={true}
        />
        <NavigationControl position="top-right" />
      </Map>
    </Card>
  );
};

export default MapCard;
