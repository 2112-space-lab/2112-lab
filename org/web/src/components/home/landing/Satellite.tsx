import React, { useRef } from 'react';
import { useFrame } from '@react-three/fiber';
import { Mesh } from 'three';
import Coverage from './Coverage';

const Satellite = () => {
    const satelliteRef = useRef<Mesh>(null);
    const solarPanelRef = useRef<Mesh>(null);

    const satelliteColor = '#ffffff';
    const solarPanelColor = '#0000ff';
    const coverageColor = '#ff0000'; // Red color for coverage

    useFrame(() => {
        if (satelliteRef.current) {
            satelliteRef.current.rotation.y += 0.00;
        }
    });

    return (
        <group>
            {/* Satellite body */}
            <mesh ref={satelliteRef} position={[0, 0, 0]}>
                <boxGeometry args={[0.1, 0.1, 0.01]} />
                <meshStandardMaterial color={satelliteColor} />
            </mesh>

            {/* Solar Panels */}
            <mesh ref={solarPanelRef} position={[0.1, 0, 0]}>
                <boxGeometry args={[0.1, 0.2, 0.02]} />
                <meshStandardMaterial color={solarPanelColor} />
            </mesh>
            <mesh ref={solarPanelRef} position={[-0.1, 0, 0]}>
                <boxGeometry args={[0.1, 0.2, 0.02]} />
                <meshStandardMaterial color={solarPanelColor} />
            </mesh>

            {/* Coverage */}
            <Coverage position={[-0.0, -0.0, 0.8]} radius={0.8} color={coverageColor} />
        </group>
    );
};

export default Satellite;
