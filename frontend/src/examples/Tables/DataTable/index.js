/**
=========================================================
* Material Dashboard 2 React - v2.2.0
=========================================================

* Product Page: https://www.creative-tim.com/product/material-dashboard-react
* Copyright 2023 Creative Tim (https://www.creative-tim.com)

Coded by www.creative-tim.com

 =========================================================

* The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
*/

import { useMemo, useEffect, useState } from "react";

// prop-types is a library for typechecking of props
import PropTypes from "prop-types";

// react-table components
import {
  useTable,
  usePagination,
  useGlobalFilter,
  useAsyncDebounce,
  useSortBy,
} from "react-table";

// @mui material components
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableContainer from "@mui/material/TableContainer";
import TableRow from "@mui/material/TableRow";
import Icon from "@mui/material/Icon";

// Material Dashboard 2 React components
import MDBox from "components/MDBox";
import MDInput from "components/MDInput";
import MDPagination from "components/MDPagination";

// Material Dashboard 2 React example components
import DataTableHeadCell from "examples/Tables/DataTable/DataTableHeadCell";
import DataTableBodyCell from "examples/Tables/DataTable/DataTableBodyCell";

import Pagination from "@mui/material/Pagination";
import { PaginationItem } from "@mui/material";

function DataTable({
  canSearch,
  columns,
  data,
  isSorted,
  pageNum,
  setPageNum,
  totalPages,
  onRowClick,
}) {
  const memoColumns = useMemo(() => columns || [], [columns]);
  const memoData = useMemo(() => data || [], [data]);

  const tableInstance = useTable(
    {
      columns: memoColumns,
      data: memoData,
    },
    useGlobalFilter,
    useSortBy,
  );

  const { rows, headerGroups, prepareRow, setGlobalFilter } = tableInstance;

  const [search, setSearch] = useState("");

  const onSearchChange = useAsyncDebounce((value) => {
    setGlobalFilter(value || undefined);
  }, 100);

  const setSortedValue = (column) => {
    if (!isSorted) return false;

    if (column.isSorted) {
      return column.isSortedDesc ? "desc" : "asc";
    }

    return "none";
  };

  return (
    <TableContainer
      sx={{
        boxShadow: "none",
        width: "100%",
        overflowX: "auto",
        overflowY: "hidden",
      }}
    >
      {canSearch && (
        <MDBox display="flex" justifyContent="flex-end" p={3}>
          <MDBox width="12rem">
            <MDInput
              placeholder="Search..."
              value={search}
              size="small"
              fullWidth
              onChange={({ currentTarget }) => {
                setSearch(currentTarget.value);
                onSearchChange(currentTarget.value);
              }}
            />
          </MDBox>
        </MDBox>
      )}

      <Table
        sx={{ tableLayout: "auto", width: "max-content", minWidth: "100%" }}
      >
        <thead>
          {headerGroups.map((hg, i) => (
            <TableRow key={i} {...hg.getHeaderGroupProps()}>
              {hg.headers.map((col, idx) => (
                <DataTableHeadCell
                  key={idx}
                  {...col.getHeaderProps(col.getSortByToggleProps())}
                  align={col.align || "left"}
                  width={col.width || "auto"}
                  sorted={setSortedValue(col)}
                >
                  {col.render("Header")}
                </DataTableHeadCell>
              ))}
            </TableRow>
          ))}
        </thead>

        <TableBody>
          {rows.length > 0 ? (
            rows.map((row) => {
              prepareRow(row);
              const { key, ...rowProps } = row.getRowProps();

              return (
                <TableRow
                  key={key}
                  {...rowProps}
                  onClick={() => onRowClick?.(row.original)}
                  sx={{
                    cursor: onRowClick ? "pointer" : "default",
                    "&:hover": onRowClick ? { backgroundColor: "#f5f5f5" } : {},
                  }}
                >
                  {row.cells.map((cell) => {
                    const { key: cellKey, ...cellProps } = cell.getCellProps();

                    return (
                      <DataTableBodyCell
                        key={cellKey}
                        {...cellProps}
                        align={cell.column.align || "left"}
                      >
                        {cell.render("Cell")}
                      </DataTableBodyCell>
                    );
                  })}
                </TableRow>
              );
            })
          ) : (
            <TableRow>
              <DataTableBodyCell
                colSpan={memoColumns?.length || 1}
                align="center"
              >
                No records found
              </DataTableBodyCell>
            </TableRow>
          )}
        </TableBody>
      </Table>

      {totalPages > 1 && (
        <MDBox
          display="flex"
          justifyContent="flex-end"
          alignItems="center"
          mt={6}
          pr={2}
          mb={5}
          gap={0.4}
        >
          {pageNum > 1 && (
            <PaginationItem
              type="previous"
              onClick={() => setPageNum((p) => Math.max(p - 1, 1))}
              sx={paginationStyle}
            />
          )}

          {(() => {
            const pages = [];
            const addPage = (p) => {
              pages.push(
                <PaginationItem
                  key={p}
                  page={p}
                  selected={pageNum === p}
                  onClick={() => setPageNum(p)}
                  sx={paginationStyle}
                />,
              );
            };

            const siblingCount = 2;

            const startPage = Math.max(2, pageNum - siblingCount);
            const endPage = Math.min(totalPages - 1, pageNum + siblingCount);
            // First page
            addPage(1);

            if (startPage > 2) {
              pages.push(
                <PaginationItem key="start" type="start-ellipsis" disabled />,
              );
            }

            for (let i = startPage; i <= endPage; i++) {
              addPage(i);
            }

            if (endPage < totalPages - 1) {
              pages.push(
                <PaginationItem key="end" type="end-ellipsis" disabled />,
              );
            }

            if (totalPages > 1) {
              addPage(totalPages);
            }
            return pages;
          })()}

          {pageNum < totalPages && (
            <PaginationItem
              type="next"
              onClick={() => {
                setPageNum((prev) => Math.min(prev + 1, totalPages));
              }}
              sx={paginationStyle}
            />
          )}
        </MDBox>
      )}
    </TableContainer>
  );
}

const paginationStyle = {
  borderRadius: "50%",
  width: 45,
  height: 45,
  minWidth: 45,
  margin: "0 !important",
  padding: 0,
  fontWeight: 600,
  fontSize: "1rem",
  boxSizing: "border-box",
  border: "1px solid #003D7F",
  backgroundColor: "transparent",
  "&:hover": {
    backgroundColor: "transparent !important",
  },
  "&.Mui-selected": {
    backgroundColor: "#003D7F !important",
    color: "#fff",
    "&:hover": {
      backgroundColor: "#003D7F !important",
    },
  },
};

DataTable.defaultProps = {
  canSearch: false,
  isSorted: true,
  pageNum: 1,
  totalPages: 1,
};

DataTable.propTypes = {
  canSearch: PropTypes.bool,
  columns: PropTypes.array.isRequired,
  data: PropTypes.array.isRequired,
  isSorted: PropTypes.bool,
  pageNum: PropTypes.number,
  setPageNum: PropTypes.func,
  totalPages: PropTypes.number,
  onRowClick: PropTypes.func,
};

export default DataTable;
