import * as React from 'react';
import { Search, X } from 'lucide-react';
import { Input } from './input';
import { Button } from './button';

type TInputSearchProps = {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
};

export function InputSearch({ value, onChange, placeholder }: TInputSearchProps) {
  return (
    <div className="relative">
      <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
      <Input
        placeholder={placeholder}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="pl-8 pr-8"
      />
      {value && (
        <Button
          variant="ghost"
          size="sm"
          className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
          onClick={() => onChange('')}
        >
          <X className="h-4 w-4" />
        </Button>
      )}
    </div>
  );
}
