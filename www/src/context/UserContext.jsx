import React, { useState } from 'react';

import { API_URL } from '../utils/globals';
import Utils from '../utils/utils';

const defaultValues = {
  apiKey: '',
  currentTestId: 0,
};

export const UserContext = React.createContext();

export const UserProvider = ({ children }) => {
  const [apiKey, setApiKey] = useState(defaultValues.apiKey);
  const [currentTestId, setCurrentTestId] = useState(
    defaultValues.currentTestId
  );

  const authUser = async (login, password) => {
    const url = `${API_URL}/auth`;

    let res = await Utils.fetchApi('POST', url, { password, login });
    if (res && res.key) {
      setApiKey(res.key);
      return true;
    }
  };

  const createUser = async (login, password) => {
    const url = `${API_URL}/users`;

    let success = await Utils.fetchApi('POST', url, { password, login });
    return success;
  };

  const parseLoginPassFromApiKey = (apiKey) => {
    const creds = apiKey.split(':');
    const login = creds[0].slice('Basic '.length);
    const password = creds[1];

    return [login, password];
  };

  const logOut = async () => {
    if (!apiKey) return true;

    const url = `${API_URL}/auth`;

    // const creds = apiKey.split(':');
    // const login = creds[0].slice('Basic '.length);
    // const password = creds[1];

    const [login, password] = parseLoginPassFromApiKey(apiKey);

    await Utils.fetchApi('PUT', url, { password, login }, apiKey);
    setApiKey(defaultValues.apiKey);
    return true;
  };

  const toProvide = {
    apiKey,
    setApiKey,
    currentTestId,
    setCurrentTestId,
    authUser,
    createUser,
    logOut,
    parseLoginPassFromApiKey
  };

  return (
    <UserContext.Provider value={toProvide}>{children}</UserContext.Provider>
  );
};
