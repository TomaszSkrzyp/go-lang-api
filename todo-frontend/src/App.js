import './App.css';
import './main.css';
import  { useContext, useState } from 'react';
import { TodoContext } from './context/TodoContext.js';
import Login from './LoginSite.js';
import { AuthContext } from './context/AuthContext.js'; 

function App() {
  const { token,logout  } = useContext(AuthContext);
  const {
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
  } = useContext(TodoContext);

  const [newTask, setNewTask] = useState('');
  const [newStatus, setNewStatus] = useState('');
  const [newDue, setNewDue] = useState('');

  const [editingTaskId, setEditingTaskId] = useState(null);
  const [editText, setEditText] = useState('');
  const [editStatus, setEditStatus] = useState('');
  const [editDue, setEditDue] = useState('');

  // If no token, show Login component
  
  if (!token) {
    return <Login />;
  }

  const statusOptions = ['Pending', 'In Progress', 'Completed', 'Canceled'];

  const handleAddTask = (e) => {
    e.preventDefault();
    if (!newTask.trim() || !newDue.trim()) return; // basic validation

    addTask({
      task: newTask,
      status: newStatus || undefined,
      due: newDue,
    });
    setNewTask('');
    setNewStatus('');
    setNewDue('');
  };

  const handleUpdateTask = (id) => {
    updateTask(id, {
      task: editText,
      status: editStatus,
      due: editDue,
    });
    setEditingTaskId(null);
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
      <button onClick={logout} style={{ float: 'right', marginBottom: '10px' }}>
        Logout
      </button>
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

      <form onSubmit={handleAddTask} style={{ marginBottom: '20px' }}>
        <input
          type="text"
          placeholder="Task description"
          value={newTask}
          onChange={(e) => setNewTask(e.target.value)}
          required
        />

        <select
          value={newStatus}
          onChange={(e) => setNewStatus(e.target.value)}
        >
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
          {todos.map((task) => (
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

                    <input
                      type="date"
                      value={editDue}
                      onChange={(e) => setEditDue(e.target.value)}
                      required
                    />
                  </div>
                  <div>
                    <button onClick={() => handleUpdateTask(task.id)}>
                      Save
                    </button>
                    <button onClick={() => setEditingTaskId(null)}>
                      Cancel
                    </button>
                  </div>
                </>
              ) : (
                <>
                  <div>
                    {task.task} - <em>{task.status}</em> -{' '}
                    <small>Due: {formatDateForInput(task.due)}</small>
                  </div>
                  <div>
                    <button
                      onClick={() => moveStatusUp(task.id)}
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

      <div className="pagination">
        <button
          onClick={() => setPage((p) => Math.max(p - 1, 1))}
          disabled={page === 1}
          className="pagination-btn"
        >
          Prev
        </button>

        <span className="pagination-info">
          Page {page} of {totalPages}
        </span>

        <button
          onClick={() => setPage((p) => Math.min(p + 1, totalPages))}
          disabled={page === totalPages}
          className="pagination-btn"
        >
          Next
        </button>
      </div>
    </div>
  );
}

export default App;
