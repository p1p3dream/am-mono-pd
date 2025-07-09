import { cn } from '@/lib/utils';

type TSkeletonProps = React.HTMLAttributes<HTMLDivElement>;

const Skeleton = ({ className, ...props }: TSkeletonProps) => {
  return <div className={cn('animate-pulse rounded-md bg-muted', className)} {...props} />;
};

export { Skeleton };
