import React from 'react';
import Card from 'components/card';
import SearchIcon from 'components/icons/SearchIcon';
import { MdChevronRight, MdChevronLeft } from 'react-icons/md';

import {
    PaginationState,
    createColumnHelper,
    useReactTable,
    ColumnFiltersState,
    getCoreRowModel,
    getFilteredRowModel,
    getFacetedRowModel,
    getFacetedUniqueValues,
    getFacetedMinMaxValues,
    getPaginationRowModel,
    getSortedRowModel,
    flexRender,
} from '@tanstack/react-table';

export type RowObj = {
    contextName: string;
    tenantId: string;
    description: string;
    maxSatellite: number;
    maxTiles: number;
    createdAt: string;
    actions: string;
};

export function ContextsOverview(props: { tableData: RowObj[] }) {
    const { tableData } = props;
    const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
        []
    );
    let defaultData = tableData;
    const [globalFilter, setGlobalFilter] = React.useState('');
    const createPages = (count: number) => {
        let arrPageCount = [];
        for (let i = 1; i <= count; i++) {
            arrPageCount.push(i);
        }
        return arrPageCount;
    };

    const columnHelper = createColumnHelper<RowObj>();
    const columns = [
        columnHelper.accessor('contextName', {
            id: 'contextName',
            header: () => (
                <p className="text-sm font-bold text-gray-600 dark:text-white">
                    CONTEXT NAME
                </p>
            ),
            cell: (info) => (
                <p className="text-sm font-bold text-navy-700 dark:text-white">
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
                <p className="text-sm font-bold text-navy-700 dark:text-white">
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
        columnHelper.accessor('maxSatellite', {
            id: 'maxSatellite',
            header: () => (
                <p className="text-sm font-bold text-gray-600 dark:text-white">
                    MAX SATELLITES
                </p>
            ),
            cell: (info) => (
                <p className="text-sm font-bold text-navy-700 dark:text-white">
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
                <p className="text-sm font-bold text-navy-700 dark:text-white">
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
                <p className="text-sm font-bold text-navy-700 dark:text-white">
                    {info.getValue()}
                </p>
            ),
        }),
        columnHelper.accessor('actions', {
            id: 'actions',
            header: () => (
                <p className="text-sm font-bold text-gray-600 dark:text-white">
                    ACTIONS
                </p>
            ),
            cell: (info) => (
                <p
                    className="cursor-pointer font-medium text-brand-500 dark:text-brand-400"
                    onClick={() => {
                        console.log(info.getValue());
                    }}
                >
                    {info.getValue()}
                </p>
            ),
        }),
    ];

    const [data, setData] = React.useState(() => [...defaultData]);
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
            globalFilter,
            pagination,
        },
        onPaginationChange: setPagination,
        onColumnFiltersChange: setColumnFilters,
        onGlobalFilterChange: setGlobalFilter,
        getCoreRowModel: getCoreRowModel(),
        getFilteredRowModel: getFilteredRowModel(),
        getSortedRowModel: getSortedRowModel(),
        getPaginationRowModel: getPaginationRowModel(),
        getFacetedRowModel: getFacetedRowModel(),
        getFacetedUniqueValues: getFacetedUniqueValues(),
        getFacetedMinMaxValues: getFacetedMinMaxValues(),
    });

    return (
        <Card extra={'w-full h-full sm:overflow-auto px-6'}>
            <div className="flex w-[400px] max-w-full items-center rounded-xl pt-[20px]">
                <div className="flex h-[38px] w-[400px] flex-grow items-center rounded-xl bg-lightPrimary text-sm text-gray-600 dark:!bg-navy-900 dark:text-white">
                    <SearchIcon />
                    <input
                        value={globalFilter ?? ''}
                        onChange={(e: any) => setGlobalFilter(e.target.value)}
                        type="text"
                        placeholder="Search...."
                        className="block w-full rounded-full bg-lightPrimary text-base text-navy-700 outline-none dark:!bg-navy-900 dark:text-white"
                    />
                </div>
            </div>

            <div className="mt-8 overflow-x-scroll xl:overflow-x-hidden">
                <table className="w-full">
                    <thead>
                        {table.getHeaderGroups().map((headerGroup) => (
                            <tr key={headerGroup.id} className="!border-px !border-gray-400">
                                {headerGroup.headers.map((header) => (
                                    <th
                                        key={header.id}
                                        colSpan={header.colSpan}
                                        onClick={header.column.getToggleSortingHandler()}
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
                        {table
                            .getRowModel()
                            .rows.slice(0, 7)
                            .map((row) => (
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
                        <p className="text-sm text-gray-700">Showing 6 rows per page</p>
                    </div>
                    <div className="flex items-center gap-2">
                        <button
                            onClick={() => table.previousPage()}
                            disabled={!table.getCanPreviousPage()}
                            className="linear flex h-10 w-10 items-center justify-center rounded-full bg-brand-500 p-2 text-lg text-white transition duration-200 hover:bg-brand-600 active:bg-brand-700 dark:bg-brand-400 dark:text-white dark:hover:bg-brand-300 dark:active:bg-brand-200"
                        >
                            <MdChevronLeft />
                        </button>

                        {createPages(table.getPageCount()).map((pageNumber, index) => (
                            <button
                                className={`linear flex h-10 w-10 items-center justify-center rounded-full p-2 text-sm transition duration-200 ${pageNumber === pageIndex + 1
                                    ? 'bg-brand-500 text-white hover:bg-brand-600 active:bg-brand-700 dark:bg-brand-400 dark:text-white dark:hover:bg-brand-300 dark:active:bg-brand-200'
                                    : 'border-[1px] border-gray-400 bg-[transparent] dark:border-white dark:text-white'
                                    }`}
                                onClick={() => table.setPageIndex(pageNumber - 1)}
                                key={index}
                            >
                                {pageNumber}
                            </button>
                        ))}
                        <button
                            onClick={() => table.nextPage()}
                            disabled={!table.getCanNextPage()}
                            className="linear flex h-10 w-10 items-center justify-center rounded-full bg-brand-500 p-2 text-lg text-white transition duration-200 hover:bg-brand-600 active:bg-brand-700 dark:bg-brand-400 dark:text-white dark:hover:bg-brand-300 dark:active:bg-brand-200"
                        >
                            <MdChevronRight />
                        </button>
                    </div>
                </div>
            </div>
        </Card>
    );
}

export default ContextsOverview;
