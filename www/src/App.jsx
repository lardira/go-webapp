import React from 'react';
import { Routes, Route } from 'react-router';

import Home from './routes/Home';
import Navbar from './routes/Navbar';
import Auth from './routes/Auth';
import Error from './routes/Error';

const App = () => {
  return (
    <Routes>
      <Route path='/' element={<Navbar />}>
        <Route index element={<Home />} />
        <Route path='/home' element={<Home />} />
        <Route path='/auth' element={<Auth />} />
        <Route path='*' exact={true} element={<Error />} />
      </Route>
    </Routes>
  );
};

export default App;
