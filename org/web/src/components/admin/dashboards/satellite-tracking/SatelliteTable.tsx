import React, { useEffect, useState } from "react";
import Card from "components/card";
import { MdChevronRight, MdChevronLeft } from "react-icons/md";
import { FiSearch } from "react-icons/fi";
import { BiTargetLock } from "react-icons/bi";

import {
  createColumnHelper,
  useReactTable,
  getCoreRowModel,
  flexRender,
} from "@tanstack/react-table";

type Satellite = {
  Id: string;
  Name: string;
  NoradID: string;
  Owner: string;
  LaunchDate: string | null;
  Apogee: number | null;
  Perigee: number | null;
  TleUpdatedAt: string | null;
  Action: string | null;
};

const columnHelper = createColumnHelper<Satellite>();

// SearchBar Component
function SearchBar({
  searchValue,
  onSearchChange,
  onSearchSubmit,
}: {
  searchValue: string;
  onSearchChange: (value: string) => void;
  onSearchSubmit: () => void;
}) {
  return (
    <div className="relative flex h-[61px] w-full items-center rounded-full bg-white px-2 py-2 shadow-xl dark:bg-navy-800">
      <button
        onClick={onSearchSubmit}
        className="flex h-full items-center justify-center rounded-l-full bg-lightPrimary text-navy-700 dark:bg-navy-900 px-4"
      >
        <FiSearch className="h-4 w-4 text-gray-400 dark:text-white" />
      </button>
      <input
        type="text"
        value={searchValue}
        onChange={(e) => onSearchChange(e.target.value)}
        onKeyDown={(e) => e.key === "Enter" && onSearchSubmit()}
        placeholder="Search..."
        className="block h-full w-full rounded-r-full bg-lightPrimary text-sm font-medium text-navy-700 outline-none dark:bg-navy-900 dark:text-white"
      />
    </div>
  );
}

// TableHeader Component
function TableHeader({ headers }: { headers: any[] }) {
  return (
    <thead className="sticky top-0 bg-white dark:bg-navy-800">
      {headers.map((headerGroup) => (
        <tr key={headerGroup.id}>
          {headerGroup.headers.map((header) => (
            <th
              key={header.id}
              colSpan={header.colSpan}
              className="py-2 text-left text-xs font-bold text-gray-600 dark:text-white"
            >
              {flexRender(
                header.column.columnDef.header,
                header.getContext()
              )}
            </th>
          ))}
        </tr>
      ))}
    </thead>
  );
}

// TableBody Component
function TableBody({
  rows,
  onRowClick,
  selectedRowId,
  setSelectedRowId,
}: {
  rows: any[];
  onRowClick: (noradID: string) => void;
  selectedRowId: string | null;
  setSelectedRowId: (id: string | null) => void;
}) {
  return (
    <tbody>
      {rows.map((row) => (
        <tr
          key={row.id}
          onClick={() => setSelectedRowId(row.id)}
          onDoubleClick={() => onRowClick(row.original.NoradID)}
          className={`cursor-pointer ${
            selectedRowId === row.id
              ? "bg-brand-900 text-white dark:bg-brand-400"
              : ""
          }`}
        >
          {row.getVisibleCells().map((cell) => (
            <td key={cell.id} className="py-2 px-4 text-sm">
              {flexRender(cell.column.columnDef.cell, cell.getContext())}
            </td>
          ))}
        </tr>
      ))}
    </tbody>
  );
}

// Pagination Component
function Pagination({
  pageIndex,
  pageCount,
  onPageChange,
}: {
  pageIndex: number;
  pageCount: number;
  onPageChange: (pageIndex: number) => void;
}) {
  const handleNextPage = () => {
    if (pageIndex + 1 < pageCount) onPageChange(pageIndex + 1);
  };

  const handlePreviousPage = () => {
    if (pageIndex > 0) onPageChange(pageIndex - 1);
  };

  return (
    <div className="flex justify-between items-center py-4">
      <button
        onClick={handlePreviousPage}
        disabled={pageIndex === 0}
        className="px-4 py-2 bg-gray-200 dark:bg-gray-700 rounded-lg"
      >
        <MdChevronLeft />
      </button>
      <p className="text-sm">
        Page {pageIndex + 1} of {pageCount}
      </p>
      <button
        onClick={handleNextPage}
        disabled={pageIndex + 1 >= pageCount}
        className="px-4 py-2 bg-gray-200 dark:bg-gray-700 rounded-lg"
      >
        <MdChevronRight />
      </button>
    </div>
  );
}

// Main SatelliteTableComponent
export default function SatelliteTableComponent(props: {
  tableData: Satellite[];
  totalItems: number;
  pageSize: number;
  pageIndex: number;
  onPageChange: (pageIndex: number) => void;
  onPageSizeChange: (pageSize: number) => void;
  searchValue: string;
  onSearchChange: (search: string) => void;
  onSearchSubmit: () => void;
  onRowClick: (noradID: string) => void;
}) {
  const {
    tableData,
    totalItems,
    pageSize,
    pageIndex,
    onPageChange,
    searchValue,
    onSearchChange,
    onSearchSubmit,
    onRowClick,
  } = props;

  const [selectedRowId, setSelectedRowId] = useState<string | null>(null);

  const columns = [
    columnHelper.accessor("Action", {
      header: () => null,
      cell: (info) => {
        const row = info.row.original;
        return row.TleUpdatedAt ? (
          <button
            className="text-blue-500 hover:text-blue-700"
            onClick={() => onRowClick(row.NoradID)}
          >
            <BiTargetLock size={20} />
          </button>
        ) : (
          <button
            className="text-gray-400 cursor-not-allowed"
            disabled
            title="TLE data not available"
          >
            <BiTargetLock size={20} />
          </button>
        );
      },
    }),
    columnHelper.accessor("Name", {
      header: "NAME",
      cell: (info) => <p className="text-sm">{info.getValue()}</p>,
    }),
    columnHelper.accessor("NoradID", {
      header: "NORAD ID",
      cell: (info) => <p className="text-sm">{info.getValue()}</p>,
    }),
    columnHelper.accessor("Owner", {
      header: "OWNER",
      cell: (info) => <p className="text-sm">{info.getValue()}</p>,
    }),
    columnHelper.accessor("LaunchDate", {
      header: "LAUNCH DATE",
      cell: (info) =>
        info.getValue()
          ? new Date(info.getValue() as string).toLocaleDateString()
          : "N/A",
    }),
  ];

  const table = useReactTable({
    data: tableData,
    columns,
    pageCount: Math.ceil(totalItems / pageSize),
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <Card extra="h-full w-full pb-4 px-4">
      {/* Search Bar */}
      <SearchBar
        searchValue={searchValue}
        onSearchChange={onSearchChange}
        onSearchSubmit={onSearchSubmit}
      />
      <div className="mt-4 overflow-auto">
        <table className="w-full">
          {/* Table Header */}
          <TableHeader headers={table.getHeaderGroups()} />
          {/* Table Body */}
          <TableBody
            rows={table.getRowModel().rows}
            onRowClick={onRowClick}
            selectedRowId={selectedRowId}
            setSelectedRowId={setSelectedRowId}
          />
        </table>
      </div>
      {/* Pagination */}
      <Pagination
        pageIndex={pageIndex}
        pageCount={Math.ceil(totalItems / pageSize)}
        onPageChange={onPageChange}
      />
    </Card>
  );
}
