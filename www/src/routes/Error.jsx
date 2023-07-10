import React, {useState} from 'react';

const Error = () => {
  const [errorCode, setErrorCode] = useState(404)
  const errorMessage = errorCode + ' error :(';

  return (
    <div className='error-route'>
      <h1 className='error message'>{errorMessage}</h1>
    </div>
  );
};

export default Error;
