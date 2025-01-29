import React from 'react';
import Earth from './Earth';
import LandingHeader from './Header';
import Footer from 'components/footer/Footer';
import Logo from './Logo';
import CTA from './Cta';
import Testimonials from './Testimonials';
import ExploreFeature from './FeatureExplore';
import FeatureSatelliteThreat from './FeatureSatelliteThreat';
import { Button } from 'components/shared/ui/button';
import { LandingPrimaryVideoCtaSection } from './SectionEarth';
import { motion } from 'framer-motion';

const Home = () => {
    return (
        <div className="flex flex-col w-full min-h-screen bg-[#001020] text-white relative">
            {/* Header */}
            <header className="relative z-50">
                <LandingHeader
                    className="absolute top-0 w-full"
                    logo={<Logo className="h-9 w-auto" />}
                    logoDark={<Logo className="h-9 w-auto" />}
                />
            </header>

            {/* SVG Decoration at the Top */}
            <div
                className="absolute top-0 left-0 w-full"
                style={{
                    height: '520px', // Match SVG height
                    zIndex: 1,
                    backgroundImage: `url('data:image/svg+xml;utf8,${encodeURIComponent(`
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1440 320">
                            <path fill="url(#gradient)" fill-opacity="1" d="M0,128L60,160C120,192,240,256,360,266.7C480,277,600,235,720,213.3C840,192,960,192,1080,208C1200,224,1320,256,1380,272L1440,288L1440,320L1380,320C1320,320,1200,320,1080,320C960,320,840,320,720,320C600,320,480,320,360,320C240,320,120,320,60,320L0,320Z"></path>
                            <defs>
                                <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="0%">
                                    <stop offset="0%" style="stop-color:#001020;stop-opacity:1" />
                                    <stop offset="50%" style="stop-color:#001530;stop-opacity:1" />
                                    <stop offset="100%" style="stop-color:#000810;stop-opacity:1" />
                                </linearGradient>
                            </defs>
                        </svg>
                    `)}')`,
                    backgroundSize: 'cover',
                    backgroundRepeat: 'no-repeat',
                    backgroundAttachment: 'fixed',
                    transform: 'rotate(180deg)', // Rotate the SVG
                    transformOrigin: 'center', // Center the rotation
                }}
            ></div>

            {/* Main Content */}
            <main className="flex flex-col items-center w-full space-y-16 relative z-10">
                {/* Earth Section with Title */}
                <section className="w-full relative">
                    <div className="absolute inset-0 flex flex-col items-center justify-center">
                        <LandingPrimaryVideoCtaSection
                            title="Master the Art of Space Warfare"
                            description="The Future of Warfare in the Final Frontier"
                            textPosition="center"
                            videoPosition="center"
                            videoSrc=""
                            withBackground
                            variant="secondary"
                            leadingComponent={
                                <p className="z-10 font-cursive font-semibold tracking-wider bg-clip-text bg-gradient-to-r text-transparent from-blue-500 via-blue-400 to-blue-500">
                                    The Future of Warfare in the Final Frontier
                                </p>
                            }
                        >
                            <motion.button className="z-10" whileHover={{ scale: 1.05 }}>
                                <Button size="xl" className="p-7 mt-6 text-xl" variant="outlineSecondary">
                                    <a href="/home/default">Start free today</a>
                                </Button>
                            </motion.button>
                        </LandingPrimaryVideoCtaSection>
                    </div>
                    <Earth />
                </section>

                {/* Testimonials Section */}
                <section className="w-full py-16">
                    <Testimonials />
                </section>

                {/* Explore Features Section */}
                <section className="w-full py-16">
                    <ExploreFeature />
                </section>

                {/* Satellite Threat Feature Section */}
                <section className="w-full py-16">
                    <FeatureSatelliteThreat />
                </section>

                {/* CTA Section */}
                <section className="w-full py-16">
                    <CTA />
                </section>
            </main>

            {/* Footer */}
            <footer className="w-full mt-5 relative z-10">
                <Footer />
            </footer>

            {/* Solid Background Below SVG */}
            <div
                className="absolute top-[520px] left-0 w-full h-[calc(100%-520px)] z-0 from-[#001020] via-[#001530] to-[#000810]"
            ></div>
        </div>
    );
};

export default Home;
