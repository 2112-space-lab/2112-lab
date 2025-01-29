import React, { useState, useEffect } from "react";
import { Spinner, Box, Text, Center } from "@chakra-ui/react";
import SatelliteTableComponent from "./SatelliteTable";
import axios from "axios";

interface SatelliteTableWithDataProps {
  onSelectNoradID: (noradID: string) => void; // Callback for selected NORAD ID
}

export default function SatelliteTableWithData({
  onSelectNoradID,
}: SatelliteTableWithDataProps) {
  const [satellites, setSatellites] = useState([]); // Satellite data
  const [loading, setLoading] = useState(true); // Loading state
  const [error, setError] = useState<string | null>(null); // Error state
  const [search, setSearch] = useState<string>(""); // Search term for filtering satellites
  const [pageIndex, setPageIndex] = useState<number>(0); // Current page index
  const [pageSize, setPageSize] = useState<number>(5); // Rows per page
  const [totalPages, setTotalPages] = useState<number>(0); // Total pages from the API
  const [totalItems, setTotalItems] = useState<number>(0); // Total items from the API
  const [submit, setSubmit] = useState<boolean>(true); // Trigger fetch for data

  useEffect(() => {
    const fetchSatellites = async () => {
      if (!submit) return;

      setLoading(true);
      setError(null);
      try {
        const response = await axios.get("http://localhost:8081/satellites/paginated", {
          params: { page: pageIndex + 1, pageSize, search },
          headers: { Accept: "application/json" },
        });

        setSatellites(response.data.satellites || []);
        setTotalItems(response.data.totalRecords || 0);
        setTotalPages(Math.ceil((response.data.totalRecords || 0) / pageSize));
      } catch (err) {
        console.error("Error fetching satellites:", err);
        setError("Failed to load satellite data.");
      } finally {
        setLoading(false);
        setSubmit(false);
      }
    };

    fetchSatellites();
  }, [pageIndex, pageSize, submit]);

  const handleSearchChange = (value: string) => {
    setSearch(value);
  };

  const handleSearchSubmit = () => {
    setPageIndex(0);
    setSubmit(true);
  };

  const handleOnPaginationChange = (index: number) => {
    setPageIndex(index);
    setSubmit(true);
  };

  const handleSatelliteSelection = (noradID: string) => {
    onSelectNoradID(noradID); // Pass the NORAD ID to the parent component
  };

  if (loading) {
    return (
      <Center className="grid h-full w-full place-items-center">
        <Spinner size="xl" color="blue.500" />
      </Center>
    );
  }

  if (error) {
    return (
      <Center className="grid h-full w-full place-items-center">
        <Box className="text-center">
          <Text className="text-lg font-bold text-red-500">{error}</Text>
        </Box>
      </Center>
    );
  }

  return (
    <Box className="grid w-full gap-4 rounded-lg shadow-md">
      <SatelliteTableComponent
        tableData={satellites}
        pageIndex={pageIndex}
        pageSize={pageSize}
        onPageChange={handleOnPaginationChange}
        onPageSizeChange={setPageSize}
        searchValue={search}
        onSearchChange={handleSearchChange}
        onSearchSubmit={handleSearchSubmit}
        totalItems={totalItems}
        onRowClick={handleSatelliteSelection} // Handle row click to get NORAD ID
      />
    </Box>
  );
}
