import React from "react";

export interface DocumentProps {
  id: string;
  link: string;
  content: string;
  document_id: string;
  className?: string;
}

const DocumentCard: React.FC<DocumentProps> = ({
  link,
  content,
  className,
}) => {
  return (
    <div className="flex flex-wrap w-full">
      <div className={`flex flex-col p-4 ps-0 w-full ${className}`}>
        <div className="flex h-10 bg-destructive-foreground rounded-md items-start mb-2">
          <div className="flex flex-grow items-start ml-2">
            <div className="text-sm text-center pt-2.5 truncate">{link}</div>
          </div>
        </div>
      </div>
      <div className="p-2 -ml-2 w-full">
        <div className="-mt-6 text-muted-foreground break-all">{content}</div>
      </div>
    </div>
  );
};

export default DocumentCard;
