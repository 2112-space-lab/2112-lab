import React, { useEffect } from "react";
import * as Cesium from "cesium";
import 'cesium/Build/Cesium/Widgets/widgets.css';

interface OrbitData {
  latitude: number;
  longitude: number;
  altitude: number;
  time: string;
}

interface CesiumViewerProps {
  orbitData?: OrbitData[];
  noradID?: string;
}

const CesiumViewer: React.FC<CesiumViewerProps> = ({ orbitData = [], noradID = "Unknown" }) => {
  useEffect(() => {
    let viewer: Cesium.Viewer | null = new Cesium.Viewer("cesiumContainer", {
      terrainProvider: new Cesium.EllipsoidTerrainProvider(),
      timeline: true,
      animation: true,
    });

    const plotOrbit = () => {
      if (!viewer || orbitData.length === 0) {
        console.warn("Orbit data not provided. Viewer will display Earth only.");
        return;
      }

      const positionProperty = new Cesium.SampledPositionProperty();

      orbitData.forEach(({ latitude, longitude, altitude, time }) => {
        try {
          const position = Cesium.Cartesian3.fromDegrees(
            longitude,
            latitude,
            altitude * 1000
          );
          const julianTime = Cesium.JulianDate.fromIso8601(time);

          positionProperty.addSample(julianTime, position);
        } catch (error) {
          console.error("Error processing orbit data:", error);
        }
      });

      const orbitPositions = orbitData.map(({ latitude, longitude, altitude }) =>
        Cesium.Cartesian3.fromDegrees(longitude, latitude, altitude * 1000)
      );

      viewer.entities.add({
        name: `Orbit Path for NORAD ID: ${noradID}`,
        polyline: {
          positions: orbitPositions,
          width: 2,
          material: Cesium.Color.YELLOW,
        },
      });

      viewer.entities.add({
        name: `Satellite ${noradID}`,
        position: positionProperty,
        point: {
          pixelSize: 10,
          color: Cesium.Color.RED,
        },
        label: {
          text: `Satellite ${noradID}`,
          font: "14pt sans-serif",
          fillColor: Cesium.Color.WHITE,
          showBackground: true,
          backgroundColor: Cesium.Color.BLACK.withAlpha(0.7),
          verticalOrigin: Cesium.VerticalOrigin.BOTTOM,
          pixelOffset: new Cesium.Cartesian2(0, -15),
        },
        path: {
          material: Cesium.Color.RED.withAlpha(0.5),
          width: 2,
        },
      });

      try {
        const startTime = Cesium.JulianDate.fromIso8601(orbitData[0]?.time ?? "");
        const stopTime = Cesium.JulianDate.fromIso8601(orbitData[orbitData.length - 1]?.time ?? "");

        if (startTime && stopTime) {
          viewer.clock.startTime = startTime.clone();
          viewer.clock.stopTime = stopTime.clone();
          viewer.clock.currentTime = startTime.clone();
          viewer.clock.clockRange = Cesium.ClockRange.LOOP_STOP;
          viewer.clock.multiplier = 60;
          viewer.timeline.zoomTo(startTime, stopTime);

          viewer.scene.camera.flyTo({
            destination: Cesium.Cartesian3.fromDegrees(
              orbitData[0]?.longitude ?? 0,
              orbitData[0]?.latitude ?? 0,
              (orbitData[0]?.altitude ?? 0) * 2000
            ),
          });
        } else {
          console.warn("Invalid orbit data times. Ensure time format is ISO8601.");
        }
      } catch (error) {
        console.error("Error setting clock or camera:", error);
      }
    };

    plotOrbit();

    return () => {
      if (viewer) {
        viewer.destroy();
        viewer = null;
      }
    };
  }, [orbitData, noradID]);

  return <div>
  <style>
    {`
      .cesium-viewer .cesium-widget-credits {
        display: none !important;
      }
    `}
  </style>
  <div id="cesiumContainer" style={{ width: "100%", height: "100%" }} />
</div>
};

export default CesiumViewer;
