const timeRanges = [
  { value: '7d', label: 'Last 7 days' },
  { value: '30d', label: 'Last 30 days' },
  { value: '90d', label: 'Last 90 days' },
  { value: '1y', label: 'Last year' },
];

const regions = [
  { value: 'all', label: 'All Regions' },
  { value: 'north', label: 'North' },
  { value: 'south', label: 'South' },
  { value: 'east', label: 'East' },
  { value: 'west', label: 'West' },
];

const propertyTypes = [
  { value: 'all', label: 'All Types' },
  { value: 'house', label: 'Houses' },
  { value: 'apartment', label: 'Apartments' },
  { value: 'commercial', label: 'Commercial' },
  { value: 'land', label: 'Land' },
];

type FilterBarProps = {
  onFilterChange: (filterType: string, value: string) => void;
};

export function FilterBar({ onFilterChange }: FilterBarProps) {
  return (
    <div>
      <select onChange={(e) => onFilterChange('timeRange', e.target.value)}>
        {timeRanges.map((range) => (
          <option key={range.value} value={range.value}>
            {range.label}
          </option>
        ))}
      </select>
      <select onChange={(e) => onFilterChange('region', e.target.value)}>
        {regions.map((region) => (
          <option key={region.value} value={region.value}>
            {region.label}
          </option>
        ))}
      </select>
      <select onChange={(e) => onFilterChange('propertyType', e.target.value)}>
        {propertyTypes.map((type) => (
          <option key={type.value} value={type.value}>
            {type.label}
          </option>
        ))}
      </select>
    </div>
  );
}
