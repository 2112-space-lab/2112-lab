import React, { useState, useEffect } from "react";
import Card from "components/card";
import { IoMdTime } from "react-icons/io";
import { MdChevronLeft, MdChevronRight } from "react-icons/md";

const VisibilityTimeline = (props: { location: { latitude: number; longitude: number }, data: any[] }) => {
  const { data } = props;
  const [selectedItem, setSelectedItem] = useState<number | null>(0); // First item is preselected
  const [canScrollLeft, setCanScrollLeft] = useState(false);
  const [canScrollRight, setCanScrollRight] = useState(false);

  useEffect(() => {
    const container = document.getElementById("timeline-carousel");
    if (container) {
      const updateScrollState = () => {
        setCanScrollLeft(container.scrollLeft > 0);
        setCanScrollRight(container.scrollLeft + container.clientWidth < container.scrollWidth);
      };

      updateScrollState();
      container.addEventListener("scroll", updateScrollState);
      return () => container.removeEventListener("scroll", updateScrollState);
    }
  }, [data]);

  const scrollCarousel = (direction: string) => {
    const container = document.getElementById("timeline-carousel");
    if (container) {
      const scrollAmount = direction === "left" ? -300 : 300;
      container.scrollBy({ left: scrollAmount, behavior: "smooth" });
    }
  };

  return (
    <Card extra={"w-full p-5"}>
      {/* Header */}
      <div className="mb-2">
        <h4 className="text-xl font-bold text-navy-700 dark:text-white">
          Visibility Timeline
        </h4>
      </div>

      {/* Carousel Controls and Items */}
      <div className="relative">
        {/* Scroll Left Button */}
        {canScrollLeft && (
          <button
            onClick={() => scrollCarousel("left")}
            className="absolute left-0 z-10 -ml-4 top-1/2 transform -translate-y-1/2 p-2 text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200"
          >
            <MdChevronLeft size={30} />
          </button>
        )}

        {/* Timeline Items */}
        <div
          id="timeline-carousel"
          className="flex w-full gap-4 overflow-x-auto scroll-smooth whitespace-nowrap custom-scrollbar"
        >
          {data.map((item, index) => (
            <VisibilityTimelineItem
              key={index}
              current={item.current}
              isSelected={selectedItem === index}
              onClick={() => setSelectedItem(index)} // Update selected item
              day={item.day}
              weekday={item.weekday}
              hours={item.hours}
              title={item.title}
              satellite={item.satellite}
              spaceID={item.spaceID}
            />
          ))}
        </div>

        {/* Scroll Right Button */}
        {canScrollRight && (
          <button
            onClick={() => scrollCarousel("right")}
            className="absolute right-0 z-10 -mr-4 top-1/2 transform -translate-y-1/2 p-2 text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200"
          >
            <MdChevronRight size={30} />
          </button>
        )}
      </div>
    </Card>
  );
};

const VisibilityTimelineItem = (props: {
  current?: boolean | string;
  isSelected: boolean;
  onClick: () => void;
  day: string;
  weekday: string;
  hours: string;
  title: string;
  satellite: string;
  spaceID: string;
}) => {
  const {
    current,
    isSelected,
    onClick,
    day,
    weekday,
    hours,
    title,
    satellite,
    spaceID,
  } = props;

  return (
    <div
      onClick={onClick}
      className={`cursor-pointer flex-shrink-0 flex w-[300px] items-end justify-between gap-4 rounded-xl p-1.5 ${isSelected
        ? "bg-brand-900 text-white dark:bg-brand-400 dark:text-white"
        : "bg-white text-navy-700 dark:bg-navy-700 dark:text-white"
        }`}
    >
      {/* Left Side */}
      <div className="flex items-center gap-3">
        <div
          className={`flex h-20 w-20 flex-col items-center justify-center rounded-xl ${isSelected
            ? "bg-brand-900 text-white dark:bg-brand-400"
            : "bg-lightPrimary text-gray-600 dark:bg-navy-900 dark:text-white"
            }`}
        >
          <p className="text-sm font-bold">{weekday}</p>
          <h5 className="text-[34px] font-bold">{day}</h5>
        </div>
        <div className="flex flex-col">
          <h5 className="text-base font-bold leading-6 break-words">
            {title}
          </h5>
          <p className="text-sm font-medium break-words">{satellite}</p>
          <div className="mt-1 flex items-center gap-2">
            <IoMdTime />
            <p className="text-sm font-bold">{hours}</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default VisibilityTimeline;
