import React, { useEffect, useState } from 'react';

interface AsyncImageProps {
  src: string;
}

const AsyncImage = (props: AsyncImageProps) => {
  const [loadedSrc, setLoadedSrc] = useState<string | null>(null);
  useEffect(() => {
    setLoadedSrc(null);
    if (props.src) {
      const handleLoad = () => {
        setLoadedSrc(props.src);
      };
      const image = new Image();
      image.addEventListener('load', handleLoad);
      image.src = props.src;
      return () => {
        image.removeEventListener('load', handleLoad);
      };
    }
  }, [props.src]);
  if (loadedSrc === props.src) {
    return <img {...props} />;
  }
  return null;
};

export default AsyncImage;
