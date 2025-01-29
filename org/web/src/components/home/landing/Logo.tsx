import React from 'react';
import logo from '/public/img/auth/auth.png';

const Logo = ({ className }: { className?: string }) => {
  return (
    <div
      style={{ backgroundImage: `url(${logo.src})` }
      }
      className={className}
    />
  );
};

export default Logo;
