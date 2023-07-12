import React, { useState, useEffect, useContext } from 'react';
import { useParams } from 'react-router';
import { useNavigate } from 'react-router';

import { UserContext } from '../context/UserContext';
import Task from '../components/Task';

import Utils from '../utils/utils';
import { API_URL } from '../utils/globals';

import './styles/VariantView.css';

const VariantView = () => {
  const { id } = useParams();
  const userContext = useContext(UserContext);
  const navigate = useNavigate();

  const [result, setResult] = useState(-1);
  const [availableTasks, setAvailableTasks] = useState([]);
  const [tasks, setTasks] = useState([]);

  const fetchTaskIds = async () => {
    const url = `${API_URL}/variants/${id}`;

    const res = await Utils.fetchApi('GET', url, undefined, userContext.apiKey);

    if (!Utils.arraysEqual(res, availableTasks)) {
      setAvailableTasks(res);
    }
  };

  const getTaskData = async (taskId) => {
    const url = `${API_URL}/variants/${id}/${taskId}`;

    const res = await Utils.fetchApi('GET', url, undefined, userContext.apiKey);
    return res.id === 0 ? null : res;
  };

  const getTestStarted = async () => {
    const url = `${API_URL}/tests/${id}`;

    const res = await Utils.fetchApi(
      'POST',
      url,
      undefined,
      userContext.apiKey
    );

    userContext.setCurrentTestId(res.id);
    alert('test started');
  };

  const sendUserAnswers = async () => {
    const url = `${API_URL}/tests/${id}`;
    const body = {
      test_id: userContext.currentTestId,
      answer: '',
    };

    for (let i = 0; i < tasks.length; i++) {
      const task = tasks[i];
      body.answer = task.userAnswer || '';
      await Utils.fetchApi('PUT', url, body, userContext.apiKey);
    }
  };

  const getTestResult = async () => {
    const url = `${API_URL}/tests/${id}/${userContext.currentTestId}`;
    const res = await Utils.fetchApi('GET', url, undefined, userContext.apiKey);
    return res;
  };

  useEffect(() => {
    if (!userContext.apiKey) navigate('/auth');
    setAvailableTasks([]);
  }, [userContext.apiKey]);

  useEffect(() => {
    if (!userContext.apiKey) navigate('/auth');
    fetchTaskIds();
    getTestStarted();
  }, []);

  useEffect(() => {
    (async () => {
      const tasksFetched = [];
      for (let i = 0; i < availableTasks.length; i++) {
        const taskId = availableTasks[i];
        const res = await getTaskData(taskId);
        tasksFetched.push(res);
      }

      setTasks(
        tasksFetched.map((t) => ({
          ...t,
          variantId: t.variant_id,
          userAnswer: null,
        }))
      );
    })();
  }, [id, availableTasks]);

  const onUserAnswerChanged = (userAnswer, taskId) => {
    const tasksWithUserAnswer = tasks.map((t) => {
      if (taskId === t.id) {
        return { ...t, userAnswer: userAnswer.trim() };
      } else {
        return { ...t };
      }
    });
    setTasks(tasksWithUserAnswer);
  };

  const onCalculateResult = () => {
    (async () => {
      await sendUserAnswers();
      const { result } = await getTestResult();
      console.log(result)
      setResult(result);
    })();
  };

  return (
    <div className='variant-view'>
      <div className='task-feed'>
        {tasks.map((t) => {
          return (
            <Task
              taskData={t}
              key={t.id}
              onUserAnswerChanged={onUserAnswerChanged}
            />
          );
        })}
        <button className='calculate-button' onClick={onCalculateResult}>
          Calculate result
        </button>
        <div className='user-test-result'>
          {result >= 0 && (
            <div className='result-percent'>Your result: {result}%</div>
          )}
        </div>
      </div>
    </div>
  );
};

export default VariantView;
