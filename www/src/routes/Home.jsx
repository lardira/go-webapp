import React, { useState, useEffect, useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import { UserContext } from '../context/UserContext';

const Home = () => {
  const navigate = useNavigate();
  const userContext = useContext(UserContext);

  useEffect(() => {
    if (!userContext.isAuth) {
      console.log('unauthenticated');
      navigate('/auth');
    }
  }, [userContext.isAuth]);

  return <div>Home</div>;
};

export default Home;
