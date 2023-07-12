import React from 'react';

import './styles/Task.css';

const Task = ({ taskData, className, onUserAnswerChanged }) => {
  const { id, task, options } = taskData;

  const onChange = (event) => {
    onUserAnswerChanged(event.target.value, id);
  };

  return (
    <div className='task-container'>
      <div className={`task ${className || ''}`}>
        <span className='task-label'>&#11088; {task}</span>
        <div className='options-in-task'>
          {options.map((option) => {
            return (
              <div className='option-in-task'>
                <span className='option-in-task-text'>&#10145; {option}</span>
              </div>
            );
          })}
        </div>
        <div className='user-answer-field'>
          <span className='user-asnwer-label'>Your answer:</span>
          <input
            className='user-asnwer-input'
            type='text'
            name='userAnswer'
            onChange={onChange}
          />
        </div>
      </div>
    </div>
  );
};

export default Task;
