'use client';
import ContextsOverview, { RowObj } from 'components/admin/contexts/overview/ContextsOverview';
import Card from 'components/card';

const tableDataContexts: RowObj[] = [
  {
    contextName: 'Earth Observation',
    tenantId: 'T001',
    description: 'Monitoring Earth’s surface changes.',
    maxSatellite: 20,
    maxTiles: 500,
    createdAt: 'Jan 15, 2023',
    actions: 'Edit',
  },
  {
    contextName: 'Weather Monitoring',
    tenantId: 'T002',
    description: 'Real-time weather data analysis.',
    maxSatellite: 10,
    maxTiles: 300,
    createdAt: 'Feb 20, 2023',
    actions: 'Edit',
  },
  {
    contextName: 'Agriculture',
    tenantId: 'T003',
    description: 'Optimizing crop yields with satellite data.',
    maxSatellite: 15,
    maxTiles: 400,
    createdAt: 'Mar 10, 2023',
    actions: 'Edit',
  },
  {
    contextName: 'Marine Surveillance',
    tenantId: 'T004',
    description: 'Tracking maritime activities.',
    maxSatellite: 8,
    maxTiles: 200,
    createdAt: 'Apr 5, 2023',
    actions: 'Edit',
  },
  {
    contextName: 'Urban Planning',
    tenantId: 'T005',
    description: 'Improving infrastructure development.',
    maxSatellite: 12,
    maxTiles: 350,
    createdAt: 'May 18, 2023',
    actions: 'Edit',
  },
];


const UserOverview = () => {
  return (
    <Card extra={'w-full h-full mt-3'}>
      <ContextsOverview tableData={tableDataContexts}></ContextsOverview>
    </Card>
  );
};

export default UserOverview;
