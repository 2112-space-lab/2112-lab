import FeatureSection from './FeatureSection';

export const SatelliteThreatFeature = () => {
    const keyPoints = [
        {
            title: 'Threat Detection',
            description: 'Monitor real-time alerts for potential satellite collisions or hazardous events.',
        },
        {
            title: 'Global Activity',
            description: 'Analyze satellite operations and risks across countries.',
        },
        {
            title: 'Threat Leaderboard',
            description: 'Track countries generating the most alerts and identify patterns in space activities.',
        },
    ];

    return (
        <FeatureSection
            title="Real-time Satellite Threat Alerts and Leaderboard"
            subtitle="Gain critical insights into satellite activities and potential threats in space. Stay ahead of emerging trends with global analytics."
            keyPoints={keyPoints}
            buttonText="Discover Our Offers"
            buttonColor="bg-red-600 text-white hover:bg-red-700"
        />
    );
};

export default SatelliteThreatFeature;
