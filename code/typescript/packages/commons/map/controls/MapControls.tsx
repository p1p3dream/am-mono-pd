import React from 'react';
import { MapPin, Plus, Minus, Edit, Move, Square } from 'lucide-react';

interface MapControlsProps {
  onZoomIn?: () => void;
  onZoomOut?: () => void;
  onPinClick?: () => void;
  onPolygonDraw?: () => void;
  isDrawingPolygon?: boolean;
  className?: string;
}

export function MapControls({
  onZoomIn,
  onZoomOut,
  onPinClick,
  onPolygonDraw,
  isDrawingPolygon = false,
  className = '',
}: MapControlsProps) {
  return (
    <div className={`absolute top-4 right-4 flex flex-col space-y-2 z-10 ${className}`}>
      <button
        onClick={onPinClick}
        className="w-10 h-10 bg-background opacity-50 text-white flex items-center justify-center rounded-md hover:opacity-100 transition-opacity cursor-pointer hover:scale-110"
        type="button"
        title="Center Map"
      >
        <MapPin size={20} />
      </button>
      <button
        onClick={onPolygonDraw}
        className={`w-10 h-10 flex items-center justify-center rounded-md transition-colors cursor-pointer hover:scale-110 ${
          isDrawingPolygon
            ? 'bg-primary text-background'
            : 'bg-background opacity-50 text-white hover:opacity-100'
        }`}
        type="button"
        title="Draw Rectangle Selection"
      >
        <Square size={20} />
      </button>
      <button
        onClick={onZoomIn}
        className="w-10 h-10 bg-background opacity-50 text-white flex items-center justify-center rounded-md hover:opacity-100 transition-opacity cursor-pointer hover:scale-110"
        type="button"
        title="Zoom In"
      >
        <Plus size={20} />
      </button>
      <button
        onClick={onZoomOut}
        className="w-10 h-10 bg-background opacity-50 text-white flex items-center justify-center rounded-md hover:opacity-100 transition-opacity cursor-pointer hover:scale-110"
        type="button"
        title="Zoom Out"
      >
        <Minus size={20} />
      </button>
    </div>
  );
}

// Simple toggle button component
interface ToggleButtonProps {
  icon: React.ReactNode;
  onClick?: () => void;
  title?: string;
  isActive?: boolean;
  className?: string;
}

export function ToggleButton({
  icon,
  onClick,
  title,
  isActive = false,
  className = '',
}: ToggleButtonProps) {
  return (
    <button
      onClick={onClick}
      className={`w-10 h-10 flex items-center justify-center rounded-md transition-colors ${
        isActive
          ? 'bg-primary text-background'
          : 'bg-background bg-opacity-80 text-white hover:bg-opacity-100'
      } ${className}`}
      type="button"
      title={title}
      aria-pressed={isActive}
    >
      {icon}
    </button>
  );
}
