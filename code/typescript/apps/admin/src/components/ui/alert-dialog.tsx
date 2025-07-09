import * as React from 'react';
import * as AlertDialogPrimitive from '@radix-ui/react-alert-dialog';
import { cn } from '@/lib/utils';
import { buttonVariants } from '@/components/ui/button';

const AlertDialog = AlertDialogPrimitive.Root;

type TAlertDialogTriggerElement = React.ElementRef<typeof AlertDialogPrimitive.Trigger>;
type TAlertDialogTriggerProps = React.ComponentPropsWithoutRef<typeof AlertDialogPrimitive.Trigger>;

const AlertDialogTrigger = React.forwardRef<TAlertDialogTriggerElement, TAlertDialogTriggerProps>(
  (props, ref) => <AlertDialogPrimitive.Trigger ref={ref} {...props} />
);

type TAlertDialogPortalProps = React.ComponentPropsWithoutRef<typeof AlertDialogPrimitive.Portal>;

function AlertDialogPortal({ ...props }: TAlertDialogPortalProps) {
  return <AlertDialogPrimitive.Portal {...props} />;
}

type TAlertDialogOverlayElement = React.ElementRef<typeof AlertDialogPrimitive.Overlay>;
type TAlertDialogOverlayProps = React.ComponentPropsWithoutRef<typeof AlertDialogPrimitive.Overlay>;

const AlertDialogOverlay = React.forwardRef<TAlertDialogOverlayElement, TAlertDialogOverlayProps>(
  ({ className, ...props }, ref) => (
    <AlertDialogPrimitive.Overlay
      className={cn(
        'fixed inset-0 z-50 bg-black/80 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0',
        className
      )}
      {...props}
      ref={ref}
    />
  )
);

type TAlertDialogContentElement = React.ElementRef<typeof AlertDialogPrimitive.Content>;
type TAlertDialogContentProps = React.ComponentPropsWithoutRef<typeof AlertDialogPrimitive.Content>;

const AlertDialogContent = React.forwardRef<TAlertDialogContentElement, TAlertDialogContentProps>(
  ({ className, ...props }, ref) => (
    <AlertDialogPortal>
      <AlertDialogOverlay />
      <AlertDialogPrimitive.Content
        ref={ref}
        className={cn(
          'fixed left-[50%] top-[50%] z-50 grid w-full max-w-lg translate-x-[-50%] translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg duration-200 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[state=closed]:slide-out-to-left-1/2 data-[state=closed]:slide-out-to-top-[48%] data-[state=open]:slide-in-from-left-1/2 data-[state=open]:slide-in-from-top-[48%] sm:rounded-lg',
          className
        )}
        {...props}
      />
    </AlertDialogPortal>
  )
);

type TAlertDialogHeaderProps = React.HTMLAttributes<HTMLDivElement>;

function AlertDialogHeader({ className, ...props }: TAlertDialogHeaderProps) {
  return (
    <div className={cn('flex flex-col space-y-2 text-center sm:text-left', className)} {...props} />
  );
}

type TAlertDialogFooterProps = React.HTMLAttributes<HTMLDivElement>;

function AlertDialogFooter({ className, ...props }: TAlertDialogFooterProps) {
  return (
    <div
      className={cn('flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2', className)}
      {...props}
    />
  );
}

type TAlertDialogTitleElement = React.ElementRef<typeof AlertDialogPrimitive.Title>;
type TAlertDialogTitleProps = React.ComponentPropsWithoutRef<typeof AlertDialogPrimitive.Title>;

const AlertDialogTitle = React.forwardRef<TAlertDialogTitleElement, TAlertDialogTitleProps>(
  ({ className, ...props }, ref) => (
    <AlertDialogPrimitive.Title
      ref={ref}
      className={cn('text-lg font-semibold', className)}
      {...props}
    />
  )
);

type TAlertDialogDescriptionElement = React.ElementRef<typeof AlertDialogPrimitive.Description>;
type TAlertDialogDescriptionProps = React.ComponentPropsWithoutRef<
  typeof AlertDialogPrimitive.Description
>;

const AlertDialogDescription = React.forwardRef<
  TAlertDialogDescriptionElement,
  TAlertDialogDescriptionProps
>(({ className, ...props }, ref) => (
  <AlertDialogPrimitive.Description
    ref={ref}
    className={cn('text-sm text-muted-foreground', className)}
    {...props}
  />
));

type TAlertDialogActionElement = React.ElementRef<typeof AlertDialogPrimitive.Action>;
type TAlertDialogActionProps = React.ComponentPropsWithoutRef<typeof AlertDialogPrimitive.Action>;

const AlertDialogAction = React.forwardRef<TAlertDialogActionElement, TAlertDialogActionProps>(
  ({ className, ...props }, ref) => (
    <AlertDialogPrimitive.Action ref={ref} className={cn(buttonVariants(), className)} {...props} />
  )
);

type TAlertDialogCancelElement = React.ElementRef<typeof AlertDialogPrimitive.Cancel>;
type TAlertDialogCancelProps = React.ComponentPropsWithoutRef<typeof AlertDialogPrimitive.Cancel>;

const AlertDialogCancel = React.forwardRef<TAlertDialogCancelElement, TAlertDialogCancelProps>(
  ({ className, ...props }, ref) => (
    <AlertDialogPrimitive.Cancel
      ref={ref}
      className={cn(buttonVariants({ variant: 'outline' }), 'mt-2 sm:mt-0', className)}
      {...props}
    />
  )
);

export {
  AlertDialog,
  AlertDialogTrigger,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogFooter,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogAction,
  AlertDialogCancel,
};
