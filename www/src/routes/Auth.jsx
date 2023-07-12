import React, { useContext, useEffect } from 'react';
import { useNavigate } from 'react-router';
import { UserContext } from '../context/UserContext';
import Form from '../components/Form';
import Utils from '../utils/utils';

import './styles/Auth.css';

const Auth = () => {
  const MIN_PASS_LENGTH = 5;

  const userContext = useContext(UserContext);
  const navigate = useNavigate();

  const onLoginSubmit = ({ login, password }) => {
    const aggregator = Utils.newAggregator();

    if (!login || !password) aggregator.add('Не все поля заполнены');
    else if (password.length < MIN_PASS_LENGTH)
      aggregator.add('Пароль не может быть меньше ' + MIN_PASS_LENGTH);

    if (!aggregator.empty()) {
      alert(aggregator.aggregate());
    } else {
      userContext.authUser(login, password).then(() => {
        navigate('/home');
      });
    }
  };

  const onSignUpSubmit = ({ login, password, confirmPassword }) => {
    const aggregator = Utils.newAggregator();

    if (!login || !password || !confirmPassword)
      aggregator.add('Не все поля заполнены');
    else if (login.length < 2) aggregator.add('Неправильный логин ' + login);
    else if (password < MIN_PASS_LENGTH)
      aggregator.add('Пароль не может быть меньше ' + MIN_PASS_LENGTH);
    else if (password !== confirmPassword)
      aggregator.add('Введённые пароли не совпадают');

    if (!aggregator.empty()) alert(aggregator.aggregate());
    else {
      userContext.createUser(login, password).then(() => {
        alert(
          'Вы успешно зарегистрированы в системе!\n' +
            'Введите ваш логин и пароль для входа в систему'
        );
      });
    }
  };

  return (
    <div id='auth'>
      <div className='login-form-container form-container'>
        <span className='login-form-message'>Уже зарегистрированы?</span>
        <Form
          className='login-form'
          fields={[
            {
              label: 'логин',
              type: 'text',
              name: 'login',
            },
            {
              label: 'пароль',
              type: 'password',
              name: 'password',
            },
          ]}
          submitField={{
            value: 'войти',
            onSubmit: onLoginSubmit,
          }}
        />
      </div>
      <div className='sign-up-form-container form-container'>
        <span className='sign-up-form-message'>Нет аккаунта?</span>
        <Form
          className='sign-up-form'
          fields={[
            {
              label: 'логин',
              type: 'text',
              name: 'login',
            },
            {
              label: 'пароль',
              type: 'password',
              name: 'password',
            },
            {
              label: 'повторите пароль',
              type: 'password',
              name: 'confirmPassword',
            },
          ]}
          submitField={{
            value: 'зарегистрироваться',
            onSubmit: onSignUpSubmit,
          }}
        />
      </div>
    </div>
  );
};

export default Auth;
