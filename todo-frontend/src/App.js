import './App.css';
import './main.css';
import React, { useEffect, useState } from 'react';

function App() {
  const [tasks, setTasks] = useState([]);
  const [page, setPage] = useState(1);
  const [limit] = useState(5);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);

  const [newTask, setNewTask] = useState('');
  const [newStatus, setNewStatus] = useState('');
  const [newDue, setNewDue] = useState(''); 

  const [editingTaskId, setEditingTaskId] = useState(null);
  const [editText, setEditText] = useState('');
  const [editStatus, setEditStatus] = useState('');
  const [editDue, setEditDue] = useState('');

  const [error, setError] = useState(null);

  const statusOptions = ['Pending', 'In Progress', 'Completed', 'Canceled'];

 
  async function fetchWithErrorHandling(url, options) {
    const res = await fetch(url, options);
    if (!res.ok) {
      let errorMsg = `Error: ${res.status}`;
      try {
        const data = await res.json();
        if (data.error) errorMsg = data.error;
        else if (data.message) errorMsg = data.message;
      } catch {
      }
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
      setTasks(data.tasks);
      setTotal(data.total);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTasks(page);
  }, [page]);

  const moveStatusUp = async (task) => {
    setError(null);
    try {
      await fetchWithErrorHandling(`http://localhost:8090/todos/${task.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ changeUp: 'true' }),
      });
      fetchTasks(page);
    } catch (err) {
      setError(err.message);
    }
  };

  const removeTask = async (taskId) => {
    setError(null);
    try {
      await fetchWithErrorHandling(`http://localhost:8090/todos/${taskId}`, {
        method: 'DELETE',
      });
      if (tasks.length === 1 && page > 1) {
        setPage(page - 1);
      } else {
        fetchTasks(page);
      }
    } catch (err) {
      setError(err.message);
    }
  };

  const updateTask = async (id) => {
    setError(null);
    try {
      await fetchWithErrorHandling(`http://localhost:8090/todos/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          task: editText,
          status: editStatus,
          due: editDue,  
        }),
      });
      await fetchTasks(page);
      setEditingTaskId(null);
    } catch (err) {
      setError(err.message);
    }
  };

  const addTask = async (e) => {
    e.preventDefault();
    if (!newTask.trim()) {
      setError('Task name is required.');
      return;
    }
    if (!newDue.trim()) {
      setError('Due date is required.');
      return;
    }
    setError(null);
    try {
      await fetchWithErrorHandling('http://localhost:8090/todos', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          task: newTask, 
          status: newStatus || undefined,
          due: newDue,     
        }),
      });
      setNewTask('');
      setNewStatus('');
      setNewDue(''); 
      fetchTasks(page);
    } catch (err) {
      setError(err.message);
    }
  };
  
  const formatDateForInput = (isoString) => {
    if (!isoString) return '';
    const datePart = isoString.split('T')[0];
    const [year, month, day] = datePart.split('-');
    return `${day}.${month}.${year}`;
  };

  const totalPages = Math.ceil(total / limit);

  return (
    <div style={{ padding: '20px' }}>
      <h1>TODO List (Page {page})</h1>

      {error && (
        <div
          style={{
            marginBottom: '10px',
            color: 'red',
            fontWeight: 'bold',
          }}
          role="alert"
        >
          {error}
        </div>
      )}

      <form onSubmit={addTask} style={{ marginBottom: '20px' }}>
        <input
          type="text"
          placeholder="Task description"
          value={newTask}
          onChange={(e) => setNewTask(e.target.value)}
          required
        />

        <select value={newStatus} onChange={(e) => setNewStatus(e.target.value)}>
          <option value="">Select Status (optional)</option>
          {statusOptions.map((status) => (
            <option key={status} value={status}>
              {status}
            </option>
          ))}
        </select>

        <input
          type="date"
          value={newDue}
          onChange={(e) => setNewDue(e.target.value)}
          required
        />

        <button type="submit">Add Task</button>
      </form>

      {loading && <p>Loading...</p>}

      {!loading && (
        <ul>
          {tasks.map((task) => (
            <li
              key={task.id}
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
              }}
            >
              {editingTaskId === task.id ? (
                <>
                  <div>
                    <input
                      type="text"
                      value={editText}
                      onChange={(e) => setEditText(e.target.value)}
                    />
                    <select
                      value={editStatus}
                      onChange={(e) => setEditStatus(e.target.value)}
                    >
                      {statusOptions.map((s) => (
                        <option key={s} value={s}>
                          {s}
                        </option>
                      ))}
                    </select>

                    {/* Edit due input */}
                    <input
                      type="date"
                      value={editDue}
                      onChange={(e) => setEditDue(e.target.value)}
                      required
                    />
                  </div>
                  <div>
                    <button onClick={() => updateTask(task.id)}>Save</button>
                    <button onClick={() => setEditingTaskId(null)}>Cancel</button>
                  </div>
                </>
              ) : (
                <>
                  <div>
                    {task.task} - <em>{task.status}</em> - <small>Due: {formatDateForInput(task.due)}</small>
                  </div>
                  <div>
                    <button
                      onClick={() => moveStatusUp(task)}
                      disabled={task.status === 'Completed'}
                    >
                      Status Up
                    </button>
                    <button
                      onClick={() => {
                        setEditingTaskId(task.id);
                        setEditText(task.task);
                        setEditStatus(task.status);
                        setEditDue(task.due); 
                      }}
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => removeTask(task.id)}
                      style={{ backgroundColor: 'red', color: 'white' }}
                    >
                      Remove
                    </button>
                  </div>
                </>
              )}
            </li>
          ))}
        </ul>
      )}

      <div style={{ marginTop: 20 }}>
        <button onClick={() => setPage((p) => Math.max(p - 1, 1))} disabled={page === 1}>
          Prev
        </button>

        <span style={{ margin: '0 10px' }}>
          Page {page} of {totalPages}
        </span>

        <button
          onClick={() => setPage((p) => Math.min(p + 1, totalPages))}
          disabled={page === totalPages}
        >
          Next
        </button>
      </div>
    </div>
  );
}

export default App;
