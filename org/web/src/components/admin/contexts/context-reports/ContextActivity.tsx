import BarChart from 'components/charts/BarChart';
import {
  barChartDataContextActivity,
  barChartOptionsContextActivity,
} from 'variables/charts';

export const ContextActivity = () => {
  return (
    <div className="relative flex h-[355px] w-full flex-col rounded-[20px] bg-white bg-clip-border px-[25px] py-[29px] shadow-3xl shadow-shadow-500 dark:!bg-navy-800 dark:shadow-none">
      {/* Header */}
      <div className="flex w-full justify-between px-[8px]">
        <h4 className="text-lg font-bold text-navy-700 dark:text-white">
          Context Activity
        </h4>
        <select
          className="text-sm font-medium text-gray-600 dark:!bg-navy-800 dark:text-white"
          name="timeframe"
          id="timeframe"
          onChange={(e) => console.log(`Selected timeframe: ${e.target.value}`)} // Placeholder for timeframe change logic
        >
          <option value="weekly">Weekly</option>
          <option value="monthly">Monthly</option>
          <option value="yearly">Yearly</option>
        </select>
      </div>

      {/* Chart Section */}
      <div className="h-full w-full">
        <BarChart
          chartData={barChartDataContextActivity}
          chartOptions={barChartOptionsContextActivity}
        />
      </div>
    </div>
  );
};

export default ContextActivity;
