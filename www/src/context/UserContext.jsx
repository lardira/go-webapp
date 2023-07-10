import React, { useState } from 'react';

const defaultValues = {
  isAuth: false,
  login: '',
  password: '',
};

export const UserContext = React.createContext();

export const UserProvider = ({ children }) => {
  const [isAuth, setIsAuth] = useState(defaultValues.isAuth);
  const [login, setLogin] = useState(defaultValues.login);
  const [password, setPassword] = useState(defaultValues.password);

  const toProvide = {
    isAuth,
    setIsAuth,
    login,
    setLogin,
    password,
    setPassword,
    clear: () => {
      setIsAuth(defaultValues.isAuth)
      setLogin(defaultValues.login)
      setPassword(defaultValues.password)
    }
  };

  return (
    <UserContext.Provider value={toProvide}>{children}</UserContext.Provider>
  );
};
