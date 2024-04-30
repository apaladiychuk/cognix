export namespace Interfaces {
  export interface TableItem {
    id: number;
    name: string;
    model_id: string;
    endpoint: string;
  }

  export interface UseFilteredHandler {
    columns: ColumnItem[];
    tableData: TableItem[];
    sortField: string;
    handleSortingChange: (accessor: string) => void;
  }

  export interface ColumnItem {
    label: string;
    accessor: string;
  }

  export interface Tabs {
    key: string;
    label: string;
  }
}