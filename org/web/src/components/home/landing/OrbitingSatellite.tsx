import React, { useRef } from 'react';
import { useFrame } from '@react-three/fiber';
import { Line } from '@react-three/drei';
import { Mesh } from 'three';
import Satellite from './Satellite';

const OrbitingSatellite = ({
    radius,
    color = '#ff9900',
    satelliteColor = '#ffffff',
    showTrajectory = false,
    initialInclination = Math.random() * Math.PI,
    initialPositionAngle = Math.random() * Math.PI * 2,
    initialDirection = Math.random() * Math.PI * 2,
}) => {
    const satelliteRef = useRef<Mesh>(null);

    const orbitPoints: [number, number, number][] = Array.from({ length: 100 }, (_, i) => {
        const angle = (i / 100) * Math.PI * 2.5;
        const x = Math.cos(angle) * radius;
        const y = Math.sin(angle) * radius;
        const z = 0;
        return [x, y, z];
    });

    useFrame(({ clock }) => {
        const elapsed = clock.getElapsedTime();
        const angle = elapsed * 0.5 + initialPositionAngle + initialDirection;
        const x = Math.cos(angle) * radius;
        const y = Math.sin(angle) * radius;
        const z = 0;
        if (satelliteRef.current) {
            satelliteRef.current.position.set(
                x * Math.cos(initialInclination) - z * Math.sin(initialInclination),
                y,
                x * Math.sin(initialInclination) + z * Math.cos(initialInclination)
            );
            satelliteRef.current.lookAt(0, 0, 0);
        }
    });

    return (
        <>
            {showTrajectory && (
                <Line
                    points={orbitPoints.map(([x, y, z]) => [
                        x * Math.cos(initialInclination) - z * Math.sin(initialInclination),
                        y,
                        x * Math.sin(initialInclination) + z * Math.cos(initialInclination),
                    ])}
                    color={color}
                    lineWidth={1}
                />
            )}

            <mesh ref={satelliteRef}>
                <Satellite />
            </mesh>
        </>
    );
};

export default OrbitingSatellite;
