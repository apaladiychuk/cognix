import React, { useEffect, useState } from "react";
import { formatDistanceToNow } from "date-fns";

export interface DocumentProps {
  id: string;
  link: string;
  content: string;
  document_id: string;
  className?: string;
  date: string;
}

const DocumentCard: React.FC<DocumentProps> = ({
  link,
  content,
  className,
  date,
}) => {
  const [timeElapsed, setTimeElapsed] = useState("");

  useEffect(() => {
    const parsedDate = new Date(date);

    const updateTimeElapsed = () => {
      setTimeElapsed(formatDistanceToNow(parsedDate, { addSuffix: true }));
    };

    updateTimeElapsed();

    const intervalId = setInterval(updateTimeElapsed, 60000);

    return () => clearInterval(intervalId);
  }, [date]);

  return (
    <div className="flex flex-wrap w-full">
      <div className={`flex flex-col p-4 ps-0 w-full ${className}`}>
        <div className="flex h-10 bg-destructive-foreground rounded-md items-start mb-2">
          <div className="flex flex-grow items-center justify-between ml-2 mb-4">
            <div className="text-sm text-center pt-2.5 truncate">{link}</div>
            <div className="pt-2.5 px-2">
              <p className="text-xs text-[#9299A3] border border-white rounded bg-white px-2">
                Updated {timeElapsed}
              </p>
            </div>
          </div>
        </div>
      </div>
      <div className="p-2 ml-3 w-full">
        <div className="-mt-6 text-muted-foreground break-all">{content}</div>
        <div className="flex mt-5 space-x-3 text-muted"></div>
      </div>
    </div>
  );
};

export default DocumentCard;
