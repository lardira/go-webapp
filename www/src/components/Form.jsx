import React from 'react';

import './styles/Form.css'

const Form = ({ fields, submitField, className }) => {
  const onSubmit = (event) => {
    event.preventDefault();

    const eventFields = {};

    Object.values(event.target).forEach((val) => {
      if (val.name) {
        eventFields[val.name] = val.value;
      }
    });

    submitField.onSubmit(eventFields, event);
    event.target.reset();
  };

  return (
    <form className={`form ${className || ''}`} onSubmit={onSubmit}>
      {fields.map((field) => {
        return (
          <label
            className='form-label'
            key={field.type + field.name + Math.random() * 10000}
          >
            <span className='form-field-name'>{field.label || field.name}</span>
            <input type={field.type} name={field.name} />
          </label>
        );
      })}
      <input type='submit' value={submitField.value} />
    </form>
  );
};

export default Form;
