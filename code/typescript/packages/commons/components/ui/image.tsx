import React, { useState, useEffect, useRef } from 'react';

interface ImageProps extends React.ImgHTMLAttributes<HTMLImageElement> {
  fallbackClassName?: string;
  wrapperClassName?: string;
  useWrapper?: boolean;
}

export function Image({
  src,
  alt,
  className,
  fallbackClassName,
  wrapperClassName,
  useWrapper = false,
  ...props
}: ImageProps) {
  const [hasError, setHasError] = useState(false);
  const imgRef = useRef<HTMLImageElement>(null);
  const prevSrcRef = useRef<string | undefined>(src);
  const fallbackRef = useRef<HTMLDivElement | null>(null);

  // Handle error by creating a fallback div
  const handleError = (e: React.SyntheticEvent<HTMLImageElement, Event>) => {
    setHasError(true);

    // Hide the image
    if (e.currentTarget) {
      e.currentTarget.style.display = 'none';
    }

    // Create fallback div if image fails to load
    const fallbackDiv = document.createElement('div');
    fallbackDiv.className = fallbackClassName || 'w-full h-full bg-gray-800';

    // If we have a wrapper, append to the wrapper, otherwise use parent node
    const parent =
      useWrapper && imgRef.current?.parentElement
        ? imgRef.current.parentElement
        : e.currentTarget.parentNode;

    parent?.appendChild(fallbackDiv);

    // Store reference to the fallback div
    fallbackRef.current = fallbackDiv;

    // Call the original onError if provided
    if (props.onError) {
      props.onError(e);
    }
  };

  // Reset error state and remove fallback when src changes
  useEffect(() => {
    // If src has changed and we previously had an error
    if (src !== prevSrcRef.current && hasError) {
      setHasError(false);

      // Remove the fallback div if it exists
      if (fallbackRef.current && fallbackRef.current.parentNode) {
        fallbackRef.current.parentNode.removeChild(fallbackRef.current);
        fallbackRef.current = null;
      }

      // Show the image again
      if (imgRef.current) {
        imgRef.current.style.display = '';
      }
    }

    // Update the previous src reference
    prevSrcRef.current = src;
  }, [src, hasError]);

  // Render with or without wrapper based on useWrapper prop
  if (useWrapper) {
    return (
      <div className={wrapperClassName || 'relative flex items-center justify-center'}>
        <img
          ref={imgRef}
          src={src}
          alt={alt || ''}
          className={className}
          {...props}
          onError={handleError}
        />
      </div>
    );
  }

  return (
    <img
      ref={imgRef}
      src={src}
      alt={alt || ''}
      className={className}
      {...props}
      onError={handleError}
    />
  );
}
