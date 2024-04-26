import { memo } from 'react';
import AddIcon from '@/assets/svgs/add-icon.svg?react';

interface Props {
  title: string;
  handleClick: () => void;
  withBtn: boolean;
}

const SettingHeader = memo(({ title, handleClick, withBtn }: Props) => {
  return (
    <div className="flex flex-row justify-between gap-1 mb-10">
      <h1 className="text-3xl font-bold leading-6">{title}</h1>
      {withBtn && (
        <button
          onClick={handleClick}
          type="button"
          className="text-white gap-3 bg-orange-500 hover:bg-orange-400  focus:outline-none  font-medium rounded-lg text-sm px-5 py-2.5 text-center inline-flex items-center dark:focus:ring-orange-500 me-2 mb-2"
        >
          <AddIcon className="size-5" />
          {title}
        </button>
      )}
    </div>
  );
});

export { SettingHeader };
