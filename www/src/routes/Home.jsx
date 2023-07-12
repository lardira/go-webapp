import React, { useState, useEffect, useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import { UserContext } from '../context/UserContext';

import Variant from '../components/Variant';

import Utils from '../utils/utils';
import { API_URL } from '../utils/globals';

import './styles/Home.css';

const Home = () => {
  const navigate = useNavigate();
  const userContext = useContext(UserContext);

  const [variants, setVariants] = useState([]);

  useEffect(() => {
    if (!userContext.apiKey) {
      console.log('unauthenticated');
      navigate('/auth');
    }

    (async () => {
      const url = `${API_URL}/variants`;
      const res = await Utils.fetchApi('GET', url, undefined, userContext.apiKey);
      setVariants(res)
    })();
  }, [userContext.apiKey]);

  return (
    <div className='home'>
      <div className='variants-feed'>
        {variants.map((v) => {
          const { name, id } = v;
          return <Variant name={name} key={id} id={id} />;
        })}
      </div>
    </div>
  );
};

export default Home;
