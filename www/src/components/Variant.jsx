import React from 'react';
import { useNavigate } from 'react-router';

import "./styles/Variant.css"

const Variant = ({ className, name, id }) => {

  const navigate = useNavigate(); 

  const onClick = () => {
    navigate('/home/variant/' + id)
  }

  return (
    <div className='variant-container'>
      <div className={`variant ${className || ''}`} onClick={onClick}>
        <span className='variant-title'>{name}</span>
      </div>
    </div>
  );
};

export default Variant;
