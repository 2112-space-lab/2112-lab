const Statistics = (props: {
  icon?: JSX.Element;
  title: string;
  value: number | string;
  subtitle?: string; // Optional subtitle for additional information
  endContent?: JSX.Element;
  extraClasses?: string; // Allow extra custom styles or classes
}) => {
  const { icon, title, value, subtitle, endContent, extraClasses } = props;
  return (
    <div
      className={`flex h-[88px] w-full justify-between rounded-[20px] bg-white bg-clip-border px-4 py-3 shadow-3xl shadow-shadow-500 dark:!bg-navy-800 dark:shadow-none ${extraClasses}`}
    >
      <div className="flex items-center gap-3">
        {icon && (
          <div className="flex h-14 w-14 items-center justify-center rounded-full bg-lightPrimary dark:!bg-navy-700">
            {icon}
          </div>
        )}
        <div>
          <h5 className="text-sm font-medium leading-5 text-gray-600 dark:text-gray-300">
            {title}
          </h5>
          <p className="mt-1 text-2xl font-bold leading-6 text-navy-700 dark:text-white">
            {value}
          </p>
          {subtitle && (
            <p className="text-xs font-medium text-gray-500 dark:text-gray-400">
              {subtitle}
            </p>
          )}
        </div>
      </div>

      {endContent && <div className="flex items-center">{endContent}</div>}
    </div>
  );
};

export default Statistics;
