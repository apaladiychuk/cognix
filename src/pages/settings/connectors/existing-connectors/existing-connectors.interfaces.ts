import { LABEL_STATUS } from '@/components/ui/label-status';

export namespace Interfaces {
  export interface TableItem {
    id: string;
    connector: string;
    status: LABEL_STATUS | string;
    last_indexed: string;
    docs_indexed: string;
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