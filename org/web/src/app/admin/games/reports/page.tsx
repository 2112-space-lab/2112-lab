'use client';
import Statistics from 'components/admin/contexts/context-reports/Statistics';
// Assets
import WorldMap from '/public/img/satellites/8k_earth_daymap.jpg';
// import ContextTrendGraph from '/public/img/contexts/ContextTrendGraph.png';
import { MdOutlineAddBusiness, MdOutlineMap } from 'react-icons/md';
import Conversion from 'components/admin/contexts/context-reports/Conversion';
import ContextActivity from 'components/admin/contexts/context-reports/ContextActivity';
import tableDataContextReports from 'components/admin/contexts/context-reports/tableDataContextReports';
import ContextReportsTable from 'components/admin/contexts/context-reports/ContextReportsTable';
import Image from 'next/image';

const ContextReport = () => {
  return (
    <div className="mt-3 h-full w-full">
      {/* Statistics */}
      <div className="mb-5 grid w-full grid-cols-1 gap-5 rounded-[20px] md:grid-cols-2 xl:grid-cols-4">
        {/* Total Contexts */}
        <Statistics
          icon={
            <div className="flex h-14 w-14 items-center justify-center rounded-full bg-lightPrimary text-4xl text-brand-500 dark:!bg-navy-700 dark:text-white">
              <MdOutlineAddBusiness />
            </div>
          }
          title="Total Contexts"
          value="125"
        />

        {/* Most Used Context */}
        <Statistics
          endContent={
            <div className="pr-3 text-xs text-gray-600">Last 30 Days</div>
          }
          title="Most Used Context"
          value="Earth Observation"
        />

        {/* Geographical Usage */}
        <Statistics
          endContent={
            <div className="flex items-center">
              <div className="relative flex h-14 w-14 items-center justify-center rounded-full">
                <Image
                  fill
                  style={{ objectFit: 'contain' }}
                  src={WorldMap}
                  alt="World Map"
                />
              </div>
              <select className="text-xs text-gray-600 dark:bg-navy-800">
                <option>Global</option>
                <option>USA</option>
                <option>Europe</option>
                <option>Asia</option>
              </select>
            </div>
          }
          title="Geographical Usage"
          value="Global"
        />

        {/* Trending Context */}
        <Statistics
          icon={
            <div className="flex h-14 w-14 items-center justify-center rounded-full bg-gradient-to-r from-[#4481EB] to-[#04BEFE] text-3xl text-white">
              <MdOutlineMap />
            </div>
          }
          endContent={
            <div className="relative flex items-center justify-center">
            </div>
          }
          title="Trending Context"
          value="Weather Monitoring"
        />
      </div>

      {/* Conversion and Tables */}
      <div className="grid w-full grid-cols-11 gap-5 rounded-[20px]">
        <div className="col-span-11 h-full w-full rounded-[20px] md:col-span-5 lg:col-span-4 xl:col-span-5 2xl:col-span-3">
          {/* Conversion Rate */}
          <div>
            <Conversion />
          </div>

          {/* Context Activity */}
          <div className="mt-3">
            <ContextActivity />
          </div>
        </div>

        {/* Context Table */}
        <div className="col-span-11 h-full w-full rounded-[20px] md:col-span-6 lg:col-span-7 xl:col-span-6 2xl:col-span-8">
          <ContextReportsTable tableData={tableDataContextReports} />
        </div>
      </div>
    </div>
  );
};

export default ContextReport;
