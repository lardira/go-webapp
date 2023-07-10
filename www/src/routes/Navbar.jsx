import React, { Fragment, useContext, useEffect } from 'react';
import { Outlet } from 'react-router';
import { Link, useNavigate } from 'react-router-dom';
import { UserContext } from '../context/UserContext';

import './styles/Navbar.css';

const Navbar = () => {
  const userContext = useContext(UserContext);
  const navigate = useNavigate();

  const onLogout = () => {
    userContext.clear();
    //TODO: send request to logout
    navigate('/')
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
          {userContext.isAuth ? (
            <div>
              {userContext.login}
              <button className='logout-button' onClick={onLogout}>
                Logout
              </button>
            </div>
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
