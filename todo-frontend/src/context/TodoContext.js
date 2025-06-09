import React, { createContext, useState, useEffect, useContext } from 'react';
import { AuthContext } from './AuthContext.js';

export const TodoContext = createContext();

export const TodoProvider = ({ children }) => {
  const { token } = useContext(AuthContext);

  const [todos, setTodos] = useState([]);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const limit = 5;
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const getHeaders = () => ({
    'Content-Type': 'application/json',
    Authorization: `Bearer ${token}`,
  });

  const fetchTodos = async () => {
    if (!token) {
      setTodos([]);
      setTotal(0);
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const res = await fetch(
        `http://localhost:8090/api/todos?page=${page}&limit=${limit}`,
        {
          headers: getHeaders(),
        }
      );

      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || 'Failed to fetch tasks');
      }

      const data = await res.json();
      setTodos(data.tasks);
      setTotal(data.total);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  // Fetch tasks on page or token change
  useEffect(() => {
    fetchTodos();
  }, [page, token]);

  const addTask = async (taskData) => {
    if (!token) return;

    setLoading(true);
    setError(null);

    try {
      const res = await fetch('http://localhost:8090/api/todos', {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify(taskData),
      });
       const text = await res.text();
       if (!res.ok) {
      let errorMessage = 'Failed to add task';
      try {
        const data = JSON.parse(text);
        errorMessage = data.error || errorMessage;
      } catch {
        errorMessage = text;
      }
      throw new Error(errorMessage);
    }

      if (page === 1) {
        await fetchTodos();  // refresh if already at page 1
      } else {
        setPage(1);          // it will fire useEffect and fetchTodos
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const updateTask = async (id, updatedData) => {
    if (!token) return;

    setLoading(true);
    setError(null);

    try {
      const res = await fetch(`http://localhost:8090/api/todos/${id}`, {
        method: 'PATCH',
        headers: getHeaders(),
        body: JSON.stringify(updatedData),
      });
       const text = await res.text();
       if (!res.ok) {
      let errorMessage = 'Failed to update task';
      try {
        const data = JSON.parse(text);
        errorMessage = data.error || errorMessage;
      } catch {
        errorMessage = text;
      }
      throw new Error(errorMessage);
    }

      if (page === 1) {
        await fetchTodos();  // refresh if already at page 1
      } else {
        setPage(1);          // it will fire useEffect and fetchTodos
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const removeTask = async (id) => {
    if (!token) return;

    setLoading(true);
    setError(null);

    try {
      const res = await fetch(`http://localhost:8090/api/todos/${id}`, {
        method: 'DELETE',
        headers: getHeaders(),
      });
       const text = await res.text();
       if (!res.ok) {
      let errorMessage = 'Failed to remove task';
      try {
        const data = JSON.parse(text);
        errorMessage = data.error || errorMessage;
      } catch {
        errorMessage = text;
      }
      throw new Error(errorMessage);
    }

      if (page === 1) {
        await fetchTodos();  // refresh if already at page 1
      } else {
        setPage(1);          // it will fire useEffect and fetchTodos
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const moveStatusUp = (id) => {
    const task = todos.find((t) => t.id === id);
    if (!task) return;

    const statusOrder = ['Pending', 'In Progress', 'Completed', 'Canceled'];
    const currentIndex = statusOrder.indexOf(task.status);
    if (currentIndex < statusOrder.length - 1) {
       updateTask(id, { changeUp: 'true' }); 
    }
  };

  return (
    <TodoContext.Provider
      value={{
        todos,
        page,
        setPage,
        total,
        limit,
        loading,
        error,
        moveStatusUp,
        removeTask,
        updateTask,
        addTask,
      }}
    >
      {children}
    </TodoContext.Provider>
  );
};
