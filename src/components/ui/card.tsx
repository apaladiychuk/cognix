import React from "react";

// Define props interface
interface CardProps {
  title: string;
  text: string;
}

const Card: React.FC<CardProps> = ({ title, text }) => {
  return (
    <div className="bg-white p-4 rounded-lg shadow-md">
      <h2 className="text-xl font-bold mb-2">{title}</h2>
      <p className="text-gray-600">{text}</p>
    </div>
  );
};

export { Card };
