import React, { createContext, useState, useEffect } from "react";

export const TodoContext = createContext();

export const TodoProvider = ({ children }) => {
  const [page, setPage] = useState(1);
  const [limit] = useState(5);
  const [total, setTotal] = useState(0);
  const [todos, setTasks] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  async function fetchWithErrorHandling(url, options) {
    const res = await fetch(url, options);
    if (!res.ok) {
      let errorMsg = `Error: ${res.status}`;
      try {
        const data = await res.json();
        if (data.error) errorMsg = data.error;
        else if (data.message) errorMsg = data.message;
      } catch {}
      throw new Error(errorMsg);
    }
    return res.json();
  }

  const fetchTasks = async (page) => {
  setLoading(true);
  setError(null);
  try {
    const data = await fetchWithErrorHandling(
      `http://localhost:8090/todos?page=${page}&limit=${limit}`
    );
    console.log('Backend response:', data);

    const tasks = Array.isArray(data.tasks)
      ? data.tasks
      : Array.isArray(data)
      ? data
      : [];
    setTasks(tasks);

    const totalCount = data.total || (Array.isArray(tasks) ? tasks.length : 0);
    setTotal(totalCount);
  } catch (err) {
    setError(err.message);
  } finally {
    setLoading(false);
  }
};

  

  const removeTask = async (taskId) => {
    setError(null);
    try {
      await fetchWithErrorHandling(
        `http://localhost:8090/todos/${taskId}`,
        { method: "DELETE" }
      );
      if (todos.length <= 1 && page !== 1) {
        setPage(page - 1);
      } else {
        fetchTasks(page);
      }
    } catch (err) {
      setError(err.message);
    }
  };
  const updateTask = async (id, { task, status, due }) => {
    setError(null);
    try {
      await fetchWithErrorHandling(`http://localhost:8090/todos/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ task, status, due }),
      });
      fetchTasks(page);
    } catch (err) {
      setError(err.message);
    }
  };

  const moveStatusUp = async (id) => {
    setError(null);
    try {
      await fetchWithErrorHandling(`http://localhost:8090/todos/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ changeUp: 'true' }),
      });
      fetchTasks(page);
      
    } catch (err) {
      setError(err.message);
    }
  };

  const addTask = async ({ task, status, due }) => {
    setError(null);
    if (!task?.trim()) {
      setError("Task name is required.");
      return;
    }
    if (!due?.trim()) {
      setError("Due date is required.");
      return;
    }
    try {
      await fetchWithErrorHandling("http://localhost:8090/todos", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ task, status: status || undefined, due }),
      });
      fetchTasks(page);
    } catch (err) {
      setError(err.message);
    }
  };
  useEffect(() => {
    fetchTasks(page);
  }, [page]);

  return (
    <TodoContext.Provider
      value={{
        todos,
        loading,
        error,
        page,
        setPage,
        total,
        limit,
        addTask,
        updateTask,
        removeTask,
        moveStatusUp,
      }}
    >
      {children}
    </TodoContext.Provider>
  );
};
