'use client';
import { useUser } from '@clerk/nextjs';
import { MonitoringGrafanaHealthServicePage } from 'components/admin/monitoring/GrafanaHealth';


const MonitoringPage = () => {
    const { isLoaded, isSignedIn, user } = useUser();

    if (!isLoaded) {
        return <div>Loading...</div>;
    }

    return (
        <MonitoringGrafanaHealthServicePage></MonitoringGrafanaHealthServicePage>
    );
};

export default MonitoringPage;
