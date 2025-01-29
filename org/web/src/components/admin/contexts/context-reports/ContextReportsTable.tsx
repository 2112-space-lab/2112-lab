import React from 'react';
import Card from 'components/card';
import { MdChevronRight, MdChevronLeft } from 'react-icons/md';

import {
  PaginationState,
  createColumnHelper,
  useReactTable,
  ColumnFiltersState,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  flexRender,
} from '@tanstack/react-table';

type ContextRow = {
  contextName: string;
  tenantId: string;
  description: string;
  createdAt: string;
  maxSatellites: number;
  maxTiles: number;
};

export function ContextReportsTable({ tableData }: { tableData: ContextRow[] }) {
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
    []
  );
  const columnHelper = createColumnHelper<ContextRow>();

  const columns = [
    columnHelper.accessor('contextName', {
      id: 'contextName',
      header: () => (
        <p className="text-sm font-bold text-gray-600 dark:text-white">
          CONTEXT NAME
        </p>
      ),
      cell: (info) => (
        <p className="text-sm font-medium text-navy-700 dark:text-white">
          {info.getValue()}
        </p>
      ),
    }),
    columnHelper.accessor('tenantId', {
      id: 'tenantId',
      header: () => (
        <p className="text-sm font-bold text-gray-600 dark:text-white">
          TENANT ID
        </p>
      ),
      cell: (info) => (
        <p className="text-sm font-medium text-navy-700 dark:text-white">
          {info.getValue()}
        </p>
      ),
    }),
    columnHelper.accessor('description', {
      id: 'description',
      header: () => (
        <p className="text-sm font-bold text-gray-600 dark:text-white">
          DESCRIPTION
        </p>
      ),
      cell: (info) => (
        <p className="text-sm font-medium text-navy-700 dark:text-white">
          {info.getValue()}
        </p>
      ),
    }),
    columnHelper.accessor('createdAt', {
      id: 'createdAt',
      header: () => (
        <p className="text-sm font-bold text-gray-600 dark:text-white">
          CREATED AT
        </p>
      ),
      cell: (info) => (
        <p className="text-sm font-medium text-navy-700 dark:text-white">
          {info.getValue()}
        </p>
      ),
    }),
    columnHelper.accessor('maxSatellites', {
      id: 'maxSatellites',
      header: () => (
        <p className="text-sm font-bold text-gray-600 dark:text-white">
          MAX SATELLITES
        </p>
      ),
      cell: (info) => (
        <p className="text-sm font-medium text-navy-700 dark:text-white">
          {info.getValue()}
        </p>
      ),
    }),
    columnHelper.accessor('maxTiles', {
      id: 'maxTiles',
      header: () => (
        <p className="text-sm font-bold text-gray-600 dark:text-white">
          MAX TILES
        </p>
      ),
      cell: (info) => (
        <p className="text-sm font-medium text-navy-700 dark:text-white">
          {info.getValue()}
        </p>
      ),
    }),
  ];

  const [data, setData] = React.useState(() => [...tableData]);
  const [{ pageIndex, pageSize }, setPagination] =
    React.useState<PaginationState>({
      pageIndex: 0,
      pageSize: 6,
    });

  const pagination = React.useMemo(
    () => ({
      pageIndex,
      pageSize,
    }),
    [pageIndex, pageSize]
  );

  const table = useReactTable({
    data,
    columns,
    state: {
      columnFilters,
      pagination,
    },
    onPaginationChange: setPagination,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
  });

  const createPages = (count: number) => {
    return Array.from({ length: count }, (_, i) => i + 1);
  };

  return (
    <Card extra="w-full h-full sm:overflow-auto px-6">
      <div className="mt-8 overflow-x-scroll xl:overflow-x-hidden">
        <table className="w-full">
          <thead>
            {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id} className="!border-px !border-gray-400">
                {headerGroup.headers.map((header) => (
                  <th
                    key={header.id}
                    colSpan={header.colSpan}
                    className="cursor-pointer border-b border-gray-200 pb-2 pr-4 pt-4 text-start dark:border-white/30"
                  >
                    <div className="items-center justify-between text-xs text-gray-200">
                      {flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                    </div>
                  </th>
                ))}
              </tr>
            ))}
          </thead>
          <tbody>
            {table.getRowModel().rows.map((row) => (
              <tr key={row.id}>
                {row.getVisibleCells().map((cell) => (
                  <td
                    key={cell.id}
                    className="min-w-[150px] border-white/0 py-3 pr-4"
                  >
                    {flexRender(
                      cell.column.columnDef.cell,
                      cell.getContext()
                    )}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
        {/* Pagination */}
        <div className="mt-2 flex h-20 w-full items-center justify-between px-6">
          <div className="flex items-center gap-3">
            <p className="text-sm text-gray-700">
              Showing {pageSize} rows per page
            </p>
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={() => table.previousPage()}
              disabled={!table.getCanPreviousPage()}
              className="linear flex h-10 w-10 items-center justify-center rounded-full bg-brand-500 p-2 text-lg text-white transition duration-200 hover:bg-brand-600 active:bg-brand-700 dark:bg-brand-400 dark:hover:bg-brand-300"
            >
              <MdChevronLeft />
            </button>

            {createPages(table.getPageCount()).map((pageNumber) => (
              <button
                key={pageNumber}
                onClick={() => table.setPageIndex(pageNumber - 1)}
                className={`linear flex h-10 w-10 items-center justify-center rounded-full p-2 text-sm transition duration-200 ${pageNumber === pageIndex + 1
                  ? 'bg-brand-500 text-white'
                  : 'border-[1px] border-gray-400 bg-transparent dark:border-white dark:text-white'
                  }`}
              >
                {pageNumber}
              </button>
            ))}

            <button
              onClick={() => table.nextPage()}
              disabled={!table.getCanNextPage()}
              className="linear flex h-10 w-10 items-center justify-center rounded-full bg-brand-500 p-2 text-lg text-white transition duration-200 hover:bg-brand-600 active:bg-brand-700 dark:bg-brand-400 dark:hover:bg-brand-300"
            >
              <MdChevronRight />
            </button>
          </div>
        </div>
      </div>
    </Card>
  );
}

export default ContextReportsTable;
