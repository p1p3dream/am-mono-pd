import React from 'react';
import { Eye, EyeOff } from 'lucide-react';

interface LayerToggleProps {
  layerId: string;
  visible: boolean;
  onToggle: (layerId: string, visible: boolean) => void;
  label: string;
  icon?: React.ReactNode;
  className?: string;
}

export function LayerToggle({
  layerId,
  visible,
  onToggle,
  label,
  icon,
  className = '',
}: LayerToggleProps) {
  return (
    <button
      onClick={() => onToggle(layerId, !visible)}
      className={`flex items-center gap-2 px-3 py-2 rounded-md transition-colors ${
        visible
          ? 'bg-primary text-background hover:bg-primary/90'
          : 'bg-background text-foreground hover:bg-background/90'
      } ${className}`}
      type="button"
      title={`${visible ? 'Hide' : 'Show'} ${label}`}
      aria-pressed={visible}
    >
      {icon || (visible ? <Eye size={18} /> : <EyeOff size={18} />)}
      <span>{label}</span>
    </button>
  );
}

// Quick array of toggles
interface LayerTogglesProps {
  layers: Array<{
    id: string;
    label: string;
    visible: boolean;
    icon?: React.ReactNode;
  }>;
  onToggle: (layerId: string, visible: boolean) => void;
  className?: string;
}

export function LayerToggles({ layers, onToggle, className = '' }: LayerTogglesProps) {
  return (
    <div className={`flex flex-wrap gap-2 ${className}`}>
      {layers.map((layer) => (
        <LayerToggle
          key={layer.id}
          layerId={layer.id}
          visible={layer.visible}
          onToggle={onToggle}
          label={layer.label}
          icon={layer.icon}
        />
      ))}
    </div>
  );
}
