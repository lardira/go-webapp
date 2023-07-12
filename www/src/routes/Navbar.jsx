import React, { Fragment, useContext, useEffect } from 'react';
import { Outlet } from 'react-router';
import { Link, useNavigate } from 'react-router-dom';
import { UserContext } from '../context/UserContext';

import './styles/Navbar.css';

const Navbar = () => {
  const userContext = useContext(UserContext);
  const navigate = useNavigate();

  const onLogout = () => {
    userContext.logOut()
    navigate('/');
  };

  return (
    <Fragment>
      <div className='navbar'>
        <div className='nav-links'>
          <Link className='nav-link' to={'/home'}>
            Home
          </Link>
        </div>
        <div className='user-profile-nav-view'>
          {userContext.apiKey ? (
            <button className='logout-button' onClick={onLogout}>
              Logout
            </button>
          ) : (
            <Link to='/auth'>Login</Link>
          )}
        </div>
      </div>

      <Outlet />
    </Fragment>
  );
};

export default Navbar;
