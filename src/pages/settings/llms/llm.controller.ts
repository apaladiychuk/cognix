import { useCallback, useState } from 'react';
import { Interfaces } from './llm.interfaces';

const columns: Interfaces.ColumnItem[] = [
  { label: 'Name', accessor: 'name' },
  { label: 'Model ID', accessor: 'model_id' },
  { label: 'Endpoint', accessor: 'endpoint' },
];

export namespace Controller {
  export function useFilterHandler(data: Interfaces.TableItem[] | []): Interfaces.UseFilteredHandler {
    const [sortField, setSortField] = useState('');
    const [order, setOrder] = useState('asc');
    const [tableData, setTableData] =
      useState<Interfaces.TableItem[]>(data);
    console.log(data)

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

    return {
      columns,
      tableData,
      sortField,
      handleSortingChange,
    };
  }
}