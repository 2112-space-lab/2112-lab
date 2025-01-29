import React, { useMemo, useRef } from 'react';
import { Sphere } from '@react-three/drei';
import { useFrame } from '@react-three/fiber';
import * as THREE from 'three';

export const Particles = ({ count = 1000, radius = 3.5, color = '#ffffff' }) => {
    const particlesRef = useRef([]); // Create a mutable ref for the particles

    // Generate particles' positions in a spherical distribution
    const particles = useMemo(() => {
        const positions = [];
        for (let i = 0; i < count; i++) {
            const phi = Math.acos(2 * Math.random() - 1); // Random vertical angle
            const theta = Math.random() * Math.PI * 2; // Random horizontal angle

            // Spherical to Cartesian conversion
            const x = radius * Math.sin(phi) * Math.cos(theta);
            const y = radius * Math.sin(phi) * Math.sin(theta);
            const z = radius * Math.cos(phi);

            positions.push({ position: [x, y, z], speed: Math.random() * 0.01 + 0.9 });
        }
        return positions;
    }, [count, radius]);

    // Animate particles
    useFrame(() => {
        particlesRef.current.forEach((particle, index) => {
            if (particle) {
                const { position } = particles[index];
                const speed = particles[index].speed;

                // Update particle rotation around the Y-axis
                const [x, y, z] = position;
                const angle = performance.now() * 0.0001 * speed; // Vary speed for each particle
                const rotatedX = x * Math.cos(angle) - z * Math.sin(angle);
                const rotatedZ = z * Math.cos(angle) + x * Math.sin(angle);

                particle.position.set(rotatedX, y, rotatedZ);
            }
        });
    });

    return (
        <>
            {particles.map(({ position }, index) => (
                <Sphere
                    args={[0.02, 8, 8]}
                    position={position}
                    key={index}
                    ref={(el) => {
                        // Assign the reference for each particle
                        if (el) particlesRef.current[index] = el;
                    }}
                >
                    <meshBasicMaterial color={color} />
                </Sphere>
            ))}
        </>
    );
};

export default Particles;
