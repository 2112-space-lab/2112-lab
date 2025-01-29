import React from 'react';
import Logo from './Logo';
import LandingHeader from './Header';
import Earth from './Earth';

export const LandingScene3D = () => {
    return (
        <div className="flex flex-col w-full min-h-screen">
            <header className="relative z-50">
                <LandingHeader
                    className="absolute top-0 w-full"
                    logo={<Logo className="h-9 w-auto" />}
                    logoDark={<Logo className="h-9 w-auto" />}
                />
            </header>
            <div style={{ height: '100vh', background: '#001020' }}>
                <Earth></Earth>
            </div>
        </div>
    );
};

export default LandingScene3D;
