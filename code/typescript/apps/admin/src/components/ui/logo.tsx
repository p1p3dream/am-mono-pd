import { cn } from '@/lib/utils';

type LogoProps = {
  className?: string;
  size?: number;
};

export function Logo({ className, size = 40 }: LogoProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 400 300"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={cn('', className)}
    >
      <path
        d="M100 180 L200 80 L300 180"
        strokeWidth="20"
        strokeLinecap="round"
        strokeLinejoin="round"
        stroke="#0A323C"
      />
      <path
        d="M 150 180 L 280 100 L 400 180"
        strokeWidth="20"
        strokeLinecap="round"
        strokeLinejoin="round"
        stroke="#0A323C"
      />
      <path d="M200 120 L240 160 L160 160 Z" fill="#FF9966" />
    </svg>
  );
}
