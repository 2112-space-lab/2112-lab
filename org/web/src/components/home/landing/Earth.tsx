import React from 'react';
import { Canvas, useLoader } from '@react-three/fiber';
import { OrbitControls, Sphere, Stars } from '@react-three/drei';
import { EffectComposer, Bloom, Noise } from '@react-three/postprocessing';
import OrbitingSatellite from './OrbitingSatellite';
import { TextureLoader } from 'three';

const Earth = () => {
    const numberOfSatellites = 20;

    const earthTexture = useLoader(TextureLoader, '/img/satellites/8k_earth_daymap.jpg');

    const satellites = Array.from({ length: numberOfSatellites }, (_, i) => {
        const inclination = (Math.PI / numberOfSatellites) * i;
        const initialPositionAngle = (Math.PI * 2 * i) / numberOfSatellites;
        return (
            <OrbitingSatellite
                key={i}
                radius={2.5}
                color={`hsl(${200}, 50%, 50%)`}
                initialInclination={inclination}
                initialPositionAngle={initialPositionAngle}
                initialDirection={Math.random() * Math.PI * 2}
            />
        );
    });

    return (
        <div style={{ height: '100vh' }}>
            <Canvas camera={{ position: [0, 0, 5] }}>
                <ambientLight intensity={0.2} />
                <directionalLight
                    position={[10, 10, 10]}
                    intensity={1}
                    color="#ffffff"
                />

                <Sphere visible args={[1, 64, 64]} scale={2}>
                    <meshStandardMaterial map={earthTexture} />
                </Sphere>

                {satellites}

                <EffectComposer>
                    <Bloom intensity={1.2} luminanceThreshold={0.3} />
                    <Noise opacity={0.04} />
                </EffectComposer>

                {/* <Stars radius={100} depth={50} count={5000} factor={4} saturation={0} fade speed={1} /> */}
                <OrbitControls enableZoom={false} />
            </Canvas>
        </div>
    );
};

export default Earth;
