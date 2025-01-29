'use client';
import { useUser } from '@clerk/nextjs';
import { MonitoringGrafanaLogsPage } from 'components/admin/monitoring/GrafanaLogs';


const MonitoringPage = () => {
    const { isLoaded, isSignedIn, user } = useUser();

    if (!isLoaded) {
        return <div>Loading...</div>;
    }

    return (
        <MonitoringGrafanaLogsPage></MonitoringGrafanaLogsPage>
    );
};

export default MonitoringPage;
