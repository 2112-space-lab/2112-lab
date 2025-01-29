import FeatureSection from './FeatureSection';

export const ExploreFeature = () => {
    const keyPoints = [
        {
            title: 'Real-Time Satellite Visibility',
            description: 'Easily check which satellites are visible above any location, at any time, with a single click.',
        },
        {
            title: 'Global Coverage',
            description: 'Access satellite data for any point on Earth, whether you’re planning a mission or stargazing.',
        },
        {
            title: 'Interactive Visualization',
            description: 'Explore satellite trajectories in an intuitive, visually rich interface for better decision-making.',
        },
    ];

    return (
        <FeatureSection
            title="Discover Visible Satellites"
            subtitle="Unlock the power of real-time satellite tracking, empowering you to explore the world above us with ease and precision."
            keyPoints={keyPoints}
            buttonText="Start Exploring for Free"
            buttonColor="bg-brand-500 text-white hover:bg-brand-600"
        // containerBg="bg-navy-900 text-white"
        // cardBg="bg-navy-800"
        />
    );
};

export default ExploreFeature;
