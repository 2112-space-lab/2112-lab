import React, { useRef } from 'react';
import { useFrame } from '@react-three/fiber';
import { Mesh } from 'three';

const Coverage = ({ position, radius, color }) => {
    const coverageRef = useRef<Mesh>(null);

    useFrame(() => {
        if (coverageRef.current) {
            coverageRef.current.lookAt(0, 0, 0);
            coverageRef.current.rotation.x -= Math.PI / 2;
        }
    });

    return (
        <mesh ref={coverageRef} position={position}>
            <coneGeometry args={[radius, radius * 2, 32]} />
            <meshStandardMaterial color={color} transparent opacity={0.5} />
        </mesh>
    );
};

export default Coverage;
