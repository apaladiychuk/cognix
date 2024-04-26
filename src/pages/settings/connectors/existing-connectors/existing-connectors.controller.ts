import { useCallback, useState } from 'react';
import { Interfaces } from './existing-connectors.interfaces';
import { LABEL_STATUS } from '@/components/ui/label-status';

const mockData: Interfaces.TableItem[] = [
  {
    id: 1,
    connector: 'Helen',
    status: 'Succeeded',
    last_indexed: '1 day ago',
    docs_indexed: '1 day ago',
  },
  {
    id: 2,
    connector: 'Mike',
    status: 'Super Admin',
    last_indexed: '1 day ago',
    docs_indexed: '1 day ago',
  },
  {
    id: 3,
    connector: 'Greg',
    status: 'User',
    last_indexed: '1 day ago',
    docs_indexed: '1 day ago',
  },
  {
    id: 4,
    connector: 'Ashely',
    status: 'Error',
    last_indexed: '1 day ago',
    docs_indexed: '1 day ago',
  },
  {
    id: 5,
    connector: 'Beyonce',
    status: 'Admin',
    last_indexed: '1 day ago',
    docs_indexed: '1 day ago',
  },
  {
    id: 6,
    connector: 'Slack',
    status: 'Super Admin',
    last_indexed: '1 day ago',
    docs_indexed: '1 day ago',
  },
];

const columns: Interfaces.ColumnItem[] = [
  { label: 'Connector', accessor: 'connector' },
  { label: 'Status', accessor: 'status' },
  { label: 'Last Indexed', accessor: 'last_indexed' },
  { label: 'Docs Indexed', accessor: 'docs_indexed' },
];

export namespace Controller {
  export function useFilterHandler(): Interfaces.UseFilteredHandler {
    const [sortField, setSortField] = useState('');
    const [order, setOrder] = useState('asc');
    const [tableData, setTableData] =
      useState<Interfaces.TableItem[]>(mockData);

    const handleSorting = useCallback(
      (sortField: string, sortOrder: string) => {
        if (sortField) {
          const sorted = [...tableData].sort((a, b) => {
            return (
              (a as any)[sortField]
                .toString()
                .localeCompare((b as any)[sortField].toString(), 'en', {
                  numeric: true,
                }) * (sortOrder === 'asc' ? 1 : -1)
            );
          });
          setTableData(sorted);
        }
      },
      [sortField]
    );

    const handleSortingChange = useCallback(
      (accessor: string): void => {
        const sortOrder =
          accessor === sortField && order === 'asc' ? 'desc' : 'asc';
        setSortField(accessor);
        setOrder(sortOrder);
        handleSorting(accessor, sortOrder);
      },
      [sortField]
    );

    const rebuildData = tableData.map((item) => {
      const statusKey = Object.keys(LABEL_STATUS).find(
        (key) => LABEL_STATUS[key as keyof typeof LABEL_STATUS] === item.status
      );

      const status = statusKey ? statusKey : LABEL_STATUS.ERROR;

      return {
        ...item,
        status: LABEL_STATUS[status as keyof typeof LABEL_STATUS],
      };
    });

    return {
      columns,
      tableData: rebuildData,
      sortField,
      handleSortingChange,
    };
  }
}