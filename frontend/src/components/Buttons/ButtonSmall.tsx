import React from 'react';

import './ButtonSmall.css';

interface ButtonSmallProps extends React.ComponentPropsWithRef<'button'> {
  onClick?: (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => void;
  message: string;
}

const ButtonSmall = (props: ButtonSmallProps) => {
  const { message, onClick } = props;
  return (
    <button className="button--small" onClick={onClick}>
      {message}
    </button>
  );
};

export default ButtonSmall;
